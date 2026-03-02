package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
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

	newPost.ID = uuid.NewString()
	now := time.Now().Format(time.RFC3339)

	newPost.CreatedAt = now
	newPost.UpdatedAt = now

	posts = append(posts, newPost)
	updatedBlog, err := json.MarshalIndent(posts, "", " ") //for readability in .json
	if err != nil {
		http.Error(w, "Failed to convert JSON", http.StatusInternalServerError)
		return
	}
	os.WriteFile("blog.json", updatedBlog, 0644)
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
		http.Error(w, "Post deleted successfully", http.StatusOK)
	}
	updatedData, err := json.MarshalIndent(updatedPosts, "", "  ")
	if err != nil {
		http.Error(w, "Failed to convert JSON", http.StatusInternalServerError)
		return
	}
	err = os.WriteFile("blog.json", updatedData, 0644)
	w.WriteHeader(http.StatusNoContent)
}

func putHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
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
