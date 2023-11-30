package model

import "time"

type Post struct {
	ID              int
	Title           string
	Subtitle        string
	Description     string
	Image           string
	Tags            []string
	Slug            string
	PublicationDate time.Time
	ModifiedDate    time.Time
}
