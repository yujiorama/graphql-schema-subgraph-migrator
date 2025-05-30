package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.69

import (
	"context"

	"example.mars/graph/model"
)

// GetPost is the resolver for the getPost field.
func (r *queryResolver) GetPost(ctx context.Context, postID *string) (*model.Post, error) {
	post, err := r.repository.findPostByID(*postID)
	return post, err
}

// GetComment is the resolver for the getComment field.
func (r *queryResolver) GetComment(ctx context.Context, commentID *string) (*model.Comment, error) {
	comment, err := r.repository.findCommentByID(*commentID)
	return comment, err
}

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
