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

func postHandler(w http.ResponseWriter, r *http.Request){
	data,err:=os.ReadFile("blog.json")
	if err!=nil{
		http.Error(w,"Failed",http.StatusInternalServerError)
	}
	var posts []Post
	json.Unmarshal(data,&posts)

	var newPost Post
	json.NewDecoder(r.Body).Decode(&newPost)

	posts=append(posts,newPost)
	updatedBlog,err:=json.MarshalIndent(posts,""," ") //for readability in .json
	if err != nil { 
		http.Error(w, "Failed to convert JSON", http.StatusInternalServerError) 
		return 
	}
	os.WriteFile("blog.json",updatedBlog,0644)
	w.Header().Set("Content-Type","application/json")
	json.NewEncoder(w).Encode(newPost)
}
func deleteHandler(w http.ResponseWriter, r *http.Request){
	parts := strings.Split(r.URL.Path, "/")
	id := parts[2]
	data,err:=os.ReadFile("blog.json")
	if err!=nil{
		http.Error(w,"Failed",http.StatusInternalServerError)
	}
	var posts []Post
	json.Unmarshal(data,&posts)

	var updatedPosts []Post
	found := false
	for _, post := range posts {
		if post.ID == id {
			found = true
			continue // to skip this post
		}
		updatedPosts = append(updatedPosts, post)
	}
	if !found {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}
	updatedData, err := json.MarshalIndent(updatedPosts, "", "  ")
	err = os.WriteFile("blog.json", updatedData, 0644)
	w.WriteHeader(http.StatusNoContent)
}

func putHandler(w http.ResponseWriter, r *http.Request){
	parts := strings.Split(r.URL.Path, "/")
	id := parts[2]
	data, err := os.ReadFile("blog.json")
	if err != nil {
		http.Error(w, "Failed", http.StatusInternalServerError)
		return
	}
	var posts []Post
	json.Unmarshal(data, &posts)

	var updatedPost Post
	json.NewDecoder(r.Body).Decode(&updatedPost)

	found := false
	for i, post := range posts {
		if post.ID == id {
			posts[i] = updatedPost
			found = true
			break
		}
	}
	if !found {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}
	updatedData, err := json.MarshalIndent(posts, "", "  ")
	if err != nil {
		http.Error(w, "Failed to convert JSON", http.StatusInternalServerError)
		return
	}
	os.WriteFile("blog.json", updatedData, 0644)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedPost)
}

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