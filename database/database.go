package database

import (
	"blog/model"
	"database/sql"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func GetDatabaseConnection() *sql.DB {
	db, err := sql.Open("sqlite3", "./database/blog.db")
	if err != nil {
		panic(err)
	}
	createTables(db)
	return db
}

func GetPosts(db *sql.DB) ([]model.Post, error) {
	rows, err := db.Query(`
	SELECT p.ID, p.Title, p.Subtitle, p.Description, p.Image, p.Slug, GROUP_CONCAT(tp.tag_title) AS TagList
	FROM posts AS p
	LEFT JOIN tags_posts AS tp ON p.ID = tp.post_id
	GROUP BY p.ID
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []model.Post
	for rows.Next() {
		var post model.Post
		var tagList string
		err := rows.Scan(&post.ID, &post.Title, &post.Subtitle, &post.Description, &post.Image, &post.Slug, &tagList)
		if err != nil {
			return nil, err
		}
		post.Tags = strings.Split(tagList, ",")

		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return posts, nil
}

func SavePost(db *sql.DB, post model.Post) error {
	var err error
	var id int
	if post.ID < 0 {
		// Insert a new post
		result, err := db.Exec("INSERT INTO posts (title, subtitle, description, image, slug, publication_date, modified_date) VALUES (?, ?, ?, ?, ?, ?, ?)",
			post.Title, post.Subtitle, post.Description, post.Image, post.Slug, time.Now(), time.Now())
		if err != nil {
			return err
		}
		lastId, err := result.LastInsertId()
		if err != nil {
			return err
		}
		id = int(lastId)
	} else {
		// Update post
		_, err = db.Exec("UPDATE posts SET title = ?, subtitle = ?, description = ?, image = ?, slug = ?, modified_date = ? WHERE id = ?",
			post.Title, post.Subtitle, post.Description, post.Image, post.Slug, time.Now(), post.ID)
		id = post.ID
	}
	if err != nil {
		return err
	}

	err = linkTags(db, post.Tags, int(id))
	return err

}

func DeletePost(db *sql.DB, id string) error {
	_, err := db.Exec("DELETE FROM posts WHERE id = ?", id)
	return err
}

func GetTags(db *sql.DB) ([]string, error) {
	rows, err := db.Query("SELECT DISTINCT title FROM tags")
	if err != nil {
		return nil, err
	}
	var tags []string
	for rows.Next() {
		var tag string
		if err := rows.Scan(&tag); err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tags, nil

}

func linkTags(db *sql.DB, tags []string, postId int) error {
	_, err := db.Exec("DELETE FROM tags_posts WHERE post_id = ?", postId)
	if err != nil {
		return err
	}
	for _, tag := range tags {
		_, err := db.Exec("INSERT OR IGNORE INTO tags (title) VALUES (?)", tag)
		if err != nil {
			return err
		}
		_, err = db.Exec("INSERT INTO tags_posts (tag_title, post_id) VALUES (?, ?)", tag, postId)
		if err != nil {
			return err
		}
	}
	return nil
}

func createTables(db *sql.DB) {
	createPostTableSQL := `CREATE TABLE IF NOT EXISTS posts (
        "id" INTEGER PRIMARY KEY AUTOINCREMENT,
        "title" TEXT,
        "subtitle" TEXT,
        "description" TEXT,
        "image" TEXT,
        "slug" TEXT,
		"publication_date" TIMESTAMP,
        "modified_date" TIMESTAMP
    );`

	creatTagTableSQL := `CREATE TABLE IF NOT EXISTS tags (
        "title" TEXT PRIMARY KEY
    );`

	creatTagPostTableSQL := `CREATE TABLE IF NOT EXISTS tags_posts (
        "post_id" INTEGER,
        "tag_title" TEXT,
        FOREIGN KEY (post_id) REFERENCES posts (id),
        FOREIGN KEY (tag_title) REFERENCES tags (title),
        PRIMARY KEY (post_id, tag_title)
    );`

	_, err := db.Exec(createPostTableSQL)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(creatTagTableSQL)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(creatTagPostTableSQL)
	if err != nil {
		panic(err)
	}
}
