package main

import (
	"fmt"
	"log"
	"net/http"
)

func main(){
	InitDB()
	http.HandleFunc("/",handler)
	http.HandleFunc("/blog",blogHandler)
	http.HandleFunc("/blog/",blogByIDHandler)
	
	fmt.Println("Server started....")
	log.Fatal(http.ListenAndServe(":8080",nil))

}