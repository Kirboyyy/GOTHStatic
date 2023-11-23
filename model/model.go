package model

type BlogPost struct {
	Title        string `json:"title"`
	HeroImageURL string `json:"heroImageURL"`
	Author       Author `json:"author"`
	HTML         string `json:"html"`
}

type Author struct {
	Name     string   `json:"name"`
	ImageURL string   `json:"imageURL"`
	Socials  []Social `json:"socials"`
}

type Social struct {
	Name  string `json:"name"`
	Image string `json:"image"`
	URL   string `json:"url"`
}
