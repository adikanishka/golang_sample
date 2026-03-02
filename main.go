package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"log"
)

type Post struct{
	ID string `json:"id"`
	Title string `json:"title"`
	Content string `json:"content"`
	Author string `json:"author"`
}

func handler(w http.ResponseWriter, r *http.Request){
	fmt.Fprintln(w,"Runing...")
}

func getHandler(w http.ResponseWriter, r *http.Request){
	data,err := os.ReadFile("blog.json")
	if err!=nil{
		http.Error(w,"Failed",http.StatusInternalServerError)
		return
	}
	
	var posts []Post
	json.Unmarshal(data, &posts)
	w.Header().Set("Content-Type","application/json")
	json.NewEncoder(w).Encode(posts)
}


func main(){
	http.HandleFunc("/",handler)
	http.HandleFunc("/get",getHandler)
	fmt.Println("Server started....")
	log.Fatal(http.ListenAndServe(":8080",nil))

}