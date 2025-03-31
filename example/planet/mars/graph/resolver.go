package graph

import (
	"encoding/json"
	"example.mars/graph/model"
	"fmt"
	"strings"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	repository *Repository
}

func NewResolver() *Resolver {
	return &Resolver{repository: &Repository{}}
}

type UserRepository interface {
	findUserByID(id string) (*model.User, error)
}

type PostRepository interface {
	findPostByID(id string) (*model.Post, error)
}

type CommentRepository interface {
	findCommentByID(id string) (*model.Comment, error)
}

type Repository struct{}

var postsJson = `[
{
"id": "1",
"title": "1: Hello World",
"body": "1: Hello World",
"author": {
"id": "1",
"name": "John",
"email": "<EMAIL>"
},
"comments": []
},
{
"id": "2",
"title": "2: Hello World",
"body": "2: Hello World",
"author": {
"id": "2",
"name": "Jane",
"email": "<EMAIL>"
},
"comments": []
},
{
"id": "3",
"title": "3: Hello World",
"body": "3: Hello World",
"author": {
"id": "3",
"name": "Bob",
"email": "<EMAIL>"
},
"comments": [
{
"id": "1",
"body": "bar",
"author": {
"id": "1",
"name": "John",
"email": "<EMAIL>"
}
},
{
"id": "2",
"body": "baz",
"author": {
"id": "2",
"name": "Jane",
"email": "<EMAIL>"
}
}
]
}
]`

func (r *Repository) findUserByID(id string) (*model.User, error) {
	posts := []model.Post{}
	err := json.Unmarshal([]byte(postsJson), &posts)
	if err != nil {
		return nil, err
	}
	for _, post := range posts {
		if strings.EqualFold(*post.Author.ID, id) {
			return post.Author, nil
		}
	}
	return nil, fmt.Errorf("not found: %v", id)
}

func (r *Repository) findPostByID(id string) (*model.Post, error) {
	posts := []model.Post{}
	err := json.Unmarshal([]byte(postsJson), &posts)
	if err != nil {
		return nil, err
	}
	for _, post := range posts {
		if strings.EqualFold(*post.ID, id) {
			return &post, nil
		}
	}
	return nil, nil
}

func (r *Repository) findCommentByID(id string) (*model.Comment, error) {
	posts := []model.Post{}
	err := json.Unmarshal([]byte(postsJson), &posts)
	if err != nil {
		return nil, err
	}
	for _, post := range posts {
		for _, comment := range post.Comment {
			if strings.EqualFold(*comment.ID, id) {
				return comment, nil
			}
		}
	}
	return nil, nil
}
