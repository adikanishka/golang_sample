package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
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

func getByIdHandler(w http.ResponseWriter, r *http.Request){
	parts:=strings.Split(r.URL.Path,"/")
	id:=parts[2]
	data,err :=os.ReadFile("blog.json")
	if err!=nil{
		http.Error(w,"Failed",http.StatusInternalServerError)
		return
	}
	var posts []Post
	json.Unmarshal(data,&posts)
	for _,p := range posts{
		if p.ID == id{
			w.Header().Set("Content-Tyoe","application/json")
			json.NewEncoder(w).Encode(p)
			return
		}
	}
	http.Error(w, "Post not found", http.StatusNotFound)

}


func main(){
	http.HandleFunc("/",handler)
	http.HandleFunc("/get",getHandler)
	http.HandleFunc("/get/",getByIdHandler)
	fmt.Println("Server started....")
	log.Fatal(http.ListenAndServe(":8080",nil))

}