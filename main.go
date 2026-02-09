package main

import "fmt"
import "net/http"

func handler(w http.ResponseWriter, r *http.Request){
	fmt.Println("Runing...")
}

func main(){
	http.HandleFunc("/",handler)
	fmt.Println("Server started....")
	http.ListenAndServe(":8080",nil)

}