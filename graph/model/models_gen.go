// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type CreateFeedFollowInput struct {
	FeedID string `json:"feedId"`
	UserID string `json:"userId"`
}

type CreateFeedInput struct {
	Name   string `json:"name"`
	URL    string `json:"url"`
	UserID string `json:"userId"`
}

type CreateUserInput struct {
	Name string `json:"name"`
}

type Feed struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	URL  string `json:"url"`
}

type FeedFollow struct {
	ID   string `json:"id"`
	Feed *Feed  `json:"feed"`
	User *User  `json:"user"`
}

type Post struct {
	ID          string  `json:"id"`
	Title       string  `json:"title"`
	Description *string `json:"description,omitempty"`
	PublishedAt *string `json:"publishedAt,omitempty"`
	URL         string  `json:"url"`
}

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
