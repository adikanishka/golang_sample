package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
	"github.com/google/uuid"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Runing...")
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	posts, err := readPosts()
	if err != nil {
		http.Error(w, "Failed to read posts", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

func getByIdHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	parts := strings.Split(r.URL.Path, "/")
	id := parts[2]
	posts, err := readPosts()
	if err != nil {
		http.Error(w, "Failed to read posts", http.StatusInternalServerError)
		return
	}
	for _, p := range posts {
		if p.ID == id {
			w.Header().Set("Content-Tyoe", "application/json")
			json.NewEncoder(w).Encode(p)
			return
		}
	}
	http.Error(w, "Post not found", http.StatusNotFound)

}

func postHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	posts, err := readPosts()
	if err != nil {
		http.Error(w, "Failed to read posts", http.StatusInternalServerError)
		return
	}

	var newPost Post
	json.NewDecoder(r.Body).Decode(&newPost)

	for _, post := range posts {
		if post.Title == newPost.Title && post.Author == newPost.Author {
			http.Error(w, "Post already exists", http.StatusConflict)
			return
		}
	}
	newPost.ID = uuid.NewString()
	now := time.Now().Format(time.RFC3339)

	newPost.CreatedAt = now
	newPost.UpdatedAt = now

	posts = append(posts, newPost)
	err = writePosts(posts)
	if err != nil {
		http.Error(w, "Failed to save posts", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newPost)
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	parts := strings.Split(r.URL.Path, "/")
	id := parts[2]
	posts, err := readPosts()
	if err != nil {
		http.Error(w, "Failed to read posts", http.StatusInternalServerError)
		return
	}

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
	}else{
		w.WriteHeader(http.StatusNoContent)
	}
	err = writePosts(updatedPosts)
	if err != nil {
		http.Error(w, "Failed to save posts", http.StatusInternalServerError)
		return
	}

}

func putHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	parts := strings.Split(r.URL.Path, "/")
	id := parts[2]
	posts, err := readPosts()
	if err != nil {
		http.Error(w, "Failed to read file", http.StatusInternalServerError)
		return
	}

	var updatedPost Post
	err = json.NewDecoder(r.Body).Decode(&updatedPost)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	found := false
	var result Post
	for i, post := range posts {
		if post.ID == id {
			if updatedPost.Title != "" {
				posts[i].Title = updatedPost.Title
			}
			if updatedPost.Content != "" {
				posts[i].Content = updatedPost.Content
			}
			if updatedPost.Author != "" {
				posts[i].Author = updatedPost.Author
			}
			posts[i].UpdatedAt = time.Now().Format(time.RFC3339)
			result = posts[i]
			found = true
			break
		}
	}
	if !found {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}
	err = writePosts(posts)
	if err != nil {
		http.Error(w, "Failed to save posts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func blogHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {

	case http.MethodGet:
		getHandler(w, r)

	case http.MethodPost:
		postHandler(w, r)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func blogByIDHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {

	case http.MethodGet:
		getByIdHandler(w, r)

	case http.MethodPut:
		putHandler(w, r)

	case http.MethodDelete:
		deleteHandler(w, r)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

