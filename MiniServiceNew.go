package werckerdemo

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
)

func main() {

	//   proxy := httputil.NewSingleHostReverseProxy(&url.URL{
	//           Scheme: "http",
	//           Host:   "localhost:9443",
	//   })

	server := http.Server{
		Addr:    ":8443",
		Handler: nil,
	}

	http.DefaultServeMux.HandleFunc("/authorize", authzFunc)
	http.DefaultServeMux.HandleFunc("/", allFunc)

	err := server.ListenAndServeTLS("/scratch/openssl/server/newcert.pem", "/scratch/openssl/server/key.pem")
	if err != nil {
		log.Fatal(err)
	}
}


func authzFunc(writer http.ResponseWriter, request *http.Request) {
	log.Printf("REQUEST:\n")
	dump, err := httputil.DumpRequest(request, true)
	if err != nil {
		log.Print(err)
	}
	log.Print(string(dump))
	log.Printf("=============================\n")
	defer request.Body.Close()

	rbody, _ := ioutil.ReadAll(request.Body)

	f_log, _ := os.OpenFile("history.txt",os.O_APPEND|os.O_RDWR|os.O_CREATE,0666)
	defer f_log.Close()

	//file
	io.WriteString(f_log, string(dump))
	//io.WriteString(f_log,string(rbody))

	// to authz proxy now
	url := "http://localhost:9999/opss/v1/authorize"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(rbody))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// fmt.Println("response Headers:", resp.Header)

	body, _ := ioutil.ReadAll(resp.Body)

	//console
	fmt.Println("RESPONSE:\n", string(body))

	//writer
	fmt.Fprint(writer, string(body))
	
	//fmt.Fprint(writer, "{\"apiVersion\": \"authorization.k8s.io/v1beta1\",\"kind\": \"SubjectAccessReview\",\"status\": {\"allowed\":true}}")

}

func allFunc(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprint(writer, request.URL.Path)
}
