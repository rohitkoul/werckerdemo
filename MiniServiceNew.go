package main

import (
	"net/http"
       // "net/url"
        "log"
	"fmt"
	"net/http/httputil"
	"io/ioutil"
	"io"
	"os"
        "bytes"
)

var sam User
func main() {
     //   proxy := httputil.NewSingleHostReverseProxy(&url.URL{
     //           Scheme: "http",
     //           Host:   "adc6170676.us.oracle.com:9443",
     //   })

	sam = User{"Sam", 12345, "sam@gmai.com", "12345678"}
	server := http.Server{
		Addr:    ":8443",
		Handler: nil,
	}
    
	http.DefaultServeMux.HandleFunc("/name", nameFunc)
	http.DefaultServeMux.HandleFunc("/id", idFunc)
	http.DefaultServeMux.HandleFunc("/authorize", authzFunc)
	http.DefaultServeMux.HandleFunc("/authenticate", authNFunc)
	http.DefaultServeMux.HandleFunc("/", allFunc)

        err := server.ListenAndServeTLS("/scratch/openssl/server/newcert.pem", "/scratch/openssl/server/key.pem")
        if err != nil {
           log.Fatal(err)
        }
}


func authNFunc(writer http.ResponseWriter, request *http.Request) {
	log.Printf("AUTHN REQUEST:\n")
	dump, err := httputil.DumpRequest(request, true)
	if err != nil {
		log.Print(err)
	}
	log.Print(string(dump))
        log.Printf("#################################\n");
        defer request.Body.Close()
        rbody, _ := ioutil.ReadAll(request.Body)

         f_log,_:=os.OpenFile("authN.txt", 
                              os.O_APPEND|os.O_RDWR|os.O_CREATE, 
                              0666)
          defer f_log.Close()

       //file
       io.WriteString(f_log,string(dump))

        // to authN proxy now
        url := "http://adc6170676.us.oracle.com:8080/authenticate"
        req, err := http.NewRequest("POST", url, bytes.NewBuffer(rbody))
        req.Header.Set("Content-Type", "application/json")

       client := &http.Client{}
       resp, err := client.Do(req)
       if err != nil {
          panic(err)
       }
       defer resp.Body.Close()
      
        body, _ := ioutil.ReadAll(resp.Body)
    
       //console
       fmt.Println("AUTHN RESPONSE:\n", string(body))
       
       //writer
       fmt.Fprint(writer, string(body))
}



func authzFunc(writer http.ResponseWriter, request *http.Request) {
	log.Printf("REQUEST:\n")
	dump, err := httputil.DumpRequest(request, true)
	if err != nil {
		log.Print(err)
	}
	log.Print(string(dump))
        log.Printf("=============================\n");
        defer request.Body.Close()
        
         rbody, _ := ioutil.ReadAll(request.Body)

         f_log,_:=os.OpenFile("history.txt", 
                              os.O_APPEND|os.O_RDWR|os.O_CREATE, 
                              0666)
          defer f_log.Close()

       //file
       io.WriteString(f_log,string(dump))
       //io.WriteString(f_log,string(rbody))

        // to authz proxy now
        url := "http://adc6170676.us.oracle.com:9999/opss/v1/authorize"
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
    
}


func allFunc(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprint(writer, request.URL.Path)
}

type User struct {
	name string
	id int
	email string
	phone string
}

func idFunc(writer http.ResponseWriter, request *http.Request) {

	fmt.Fprint(writer, sam.id)
}


func nameFunc(w http.ResponseWriter, r *http.Request) {

	fmt.Fprint(w, sam.name)
}

