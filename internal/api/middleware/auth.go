package middleware

import (
	"fmt"
	"net/http"

	"github.com/PedroMartini98/rss_aggregator_go/internal/database"
	"github.com/PedroMartini98/rss_aggregator_go/internal/response"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (h *middlewareHandler) Auth(calledHandler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("ApiKey")
		if apiKey == "" {
			response.WithError(w, http.StatusForbidden, "Did not recieve ApiKey from the header")
			return
		}

		user, err := h.dbQueries.GetUserByApiKey(r.Context(), apiKey)
		if err != nil {
			response.WithError(w, http.StatusForbidden, fmt.Sprintf("Error getting user with ApiKey: %v", err))
			return
		}

		calledHandler(w, r, user)

	}
}
