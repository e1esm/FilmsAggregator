package server

import (
	"context"
	"encoding/json"
	"github.com/e1esm/FilmsAggregator/internal/auth"
	"github.com/e1esm/FilmsAggregator/utils/logger"
	"io"
	"net/http"
	"time"
)

func (s *AggregatorServer) SignUp(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var inputUser auth.User

	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(bytes, &inputUser)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	user, err := s.AuthService.CreateUser(context.Background(), inputUser)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	bytes, err = json.Marshal(&user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(bytes)

}

func (s *AggregatorServer) SignIn(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var signInReqest auth.SignInRequest
	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(bytes, &signInReqest)
	if err != nil {
		logger.Logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	token, err := s.AuthService.GenerateToken(r.Context(), signInReqest.Username, signInReqest.Password)
	if err != nil {
		logger.Logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	bytes, err = json.Marshal(token)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	cookie := http.Cookie{Name: "jwt_token",
		Value: token, Expires: time.Now().Add(12 * time.Hour), HttpOnly: true, Path: "/"}

	http.SetCookie(w, &cookie)
	w.Header().Set("Authorization", token)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
}
