//go:generate go run github.com/99designs/gqlgen generate

package graph

import (
	"github.com/lf-hernandez/go-rss-aggregator/graph/model"
	"github.com/lf-hernandez/go-rss-aggregator/internal/database"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	UserStore map[string]model.User
	Database  *database.Queries
}
