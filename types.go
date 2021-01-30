package main

type OriginProject struct {
	Repo   string   `json:"repo"`
	Org    string   `json:"org"`
	Tags   []string `json:"tags"`
	Name   string   `json:"name"`
	Branch string   `json:"branch"`
}

type Tag struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

type JSONFile struct {
	Date     string    `json:"date"`
	Tags     []Tag     `json:"tags"`
	Projects []Project `json:"projects"`
}

// Project -
type Project struct {
	Name        string `json:"name"`
	FullName    string `json:"full_name"`
	Description string `json:"description"`
	Stars       int    `json:"stars"`
	Trends      struct {
		Daily   int `json:"daily"`
		Weekly  int `json:"weekly"`
		Monthly int `json:"monthly"`
		Yearly  int `json:"yearly"`
	} `json:"trends"`
	Tags             []string `json:"tags"`
	ContributorCount int      `json:"contributor_count"`
	PushedAt         string   `json:"pushed_at"`
	OwnerID          int64    `json:"owner_id"`
	CreatedAt        string   `json:"created_at"`

	// below special to npm
	Npm string `json:"npm"`
	// origin is npm, now show be go.dev
	Downloads int `json:"downloads"`

	Icon   string `json:"icon"`
	Branch string `json:"branch"`
	URL    string `json:"url"`
}
