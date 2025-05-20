package middleware

import (
	"github.com/PedroMartini98/rss_aggregator_go/internal/database"
)

type middlewareHandler struct {
	dbQueries *database.Queries
}

func NewMiddlewareHandler(dbQueries *database.Queries) *middlewareHandler {
	return &middlewareHandler{
		dbQueries: dbQueries,
	}
}
