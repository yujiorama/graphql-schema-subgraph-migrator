package graph

import (
	"encoding/json"
	"example.venus/graph/model"
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

type Repository struct{}

var usersJson = `[
{
"id": "1",
"name": "John",
"email": "<EMAIL>"
},
{
"id": "2",
"name": "Jane",
"email": "<EMAIL>"
},
{
"id": "3",
"name": "Bob",
"email": "<EMAIL>"
}
]`

func (r *Repository) findUserByID(id string) (*model.User, error) {
	users := []model.User{}
	err := json.Unmarshal([]byte(usersJson), &users)
	if err != nil {
		return nil, err
	}
	for _, user := range users {
		if strings.EqualFold(*user.ID, id) {
			return &user, nil
		}
	}

	return nil, fmt.Errorf("not found: %v", id)
}
