package main

import (
	"encoding/json"
	"os"
)


func readPosts() ([]Post, error) {
	data, err := os.ReadFile("blog.json")
	if err != nil {
		return nil, err
	}

	var posts []Post
	err = json.Unmarshal(data, &posts)
	return posts, err
}

// func writePosts(posts []Post) error {
// 	data, err := json.MarshalIndent(posts, "", "  ")
// 	if err != nil {
// 		return err
// 	}
// 	return os.WriteFile("blog.json", data, 0644)
// }