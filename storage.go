package main

import (
	"time"
)

func readPosts() ([]Post, error) {

	rows, err := DB.Query(`
		SELECT id, title, content, author, created_at, updated_at
		FROM posts
	`)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var posts []Post

	for rows.Next() {

		var p Post
		var created time.Time
		var updated time.Time

		err := rows.Scan(
			&p.ID,
			&p.Title,
			&p.Content,
			&p.Author,
			&created,
			&updated,
		)

		if err != nil {
			return nil, err
		}

		p.CreatedAt = created.Format(time.RFC3339)
		p.UpdatedAt = updated.Format(time.RFC3339)

		posts = append(posts, p)
	}

	return posts, nil
}

func writePosts(posts []Post) error {

	tx, err := DB.Begin()

	if err != nil {
		return err
	}

	_, err = tx.Exec("DELETE FROM posts")

	if err != nil {
		tx.Rollback()
		return err
	}

	query := `
	INSERT INTO posts (id, title, content, author, created_at, updated_at)
	VALUES ($1,$2,$3,$4,$5,$6)
	`

	for _, post := range posts {

		created, _ := time.Parse(time.RFC3339, post.CreatedAt)
		updated, _ := time.Parse(time.RFC3339, post.UpdatedAt)

		_, err := tx.Exec(
			query,
			post.ID,
			post.Title,
			post.Content,
			post.Author,
			created,
			updated,
		)

		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}
