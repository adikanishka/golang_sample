package main

import (
	"fmt"
	"log"
	"net/http"
)

func main(){
	http.HandleFunc("/",handler)
	http.HandleFunc("/get",getHandler)
	http.HandleFunc("/get/",getByIdHandler)
	http.HandleFunc("/post",postHandler)
	http.HandleFunc("/delete/",deleteHandler)
	http.HandleFunc("/put/",putHandler)
	fmt.Println("Server started....")
	log.Fatal(http.ListenAndServe(":8080",nil))

}