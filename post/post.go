package post

import (
	"blog/database"
	"blog/generator"
	"blog/model"
	"blog/parsing"
	"blog/web/pages"
	"database/sql"
	"errors"
)

func GeneratePost(content string, db *sql.DB) (string, error) {
	html, post := parsing.ConvertMarkdownToHTML(content)
	file, err := generator.SaveTemplComponent(post.Slug, pages.BlogPage(post, html))
	if err != nil {
		//slog.Error("Can't generate TemplComponent", "post", post, "error", err)
		return "", err
	}
	err = database.SavePost(db, post)
	if err != nil {
		//slog.Error("Can't save post to database", "post", post, "error", err)
		return "", err
	}
	return file.Name(), nil
}

func GetPosts(db *sql.DB) ([]model.Post, error) {
	return database.GetPosts(db)
}

func DeletePost(id string, db *sql.DB) error {
	if id == "" {
		return errors.New("Can't delete post without id given!")
	}
	return database.DeletePost(db, id)
}
