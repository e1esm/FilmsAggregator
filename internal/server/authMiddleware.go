package server

import (
	"fmt"
	"github.com/e1esm/FilmsAggregator/internal/auth"
	"github.com/e1esm/FilmsAggregator/utils/logger"
	"net/http"
)

func (s *AggregatorServer) UserIdentity(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("jwt_token")

		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if cookie.Value == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		userId, userRole, err := s.AuthService.ParseToken(r.Context(), cookie.Value)
		if err != nil {
			logger.Logger.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if auth.Role(userRole) == auth.GUEST {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		w.Header().Set("userID", fmt.Sprintf("%v", userId))
		w.Header().Set("userRole", userRole)
		next.ServeHTTP(w, r)
	})
}
