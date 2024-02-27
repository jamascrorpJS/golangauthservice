package delivery

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Microsoft/go-winio/pkg/guid"
	"github.com/jamascrorpJS/auth/internal/core"
	"github.com/jamascrorpJS/auth/pkg/hash"
)

var (
	ErrBody        = ""
	ErrGuid        = ""
	ErrCreateUser  = ""
	ErrCreateToken = ""
	ErrCookie      = ""
	ErrToken       = ""
)

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {

	user := core.User{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "error:%v", errors.New(ErrBody))
		return
	}
	id, err := guid.NewV4()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "error:%v", errors.New(ErrGuid))
		return
	}
	user.ID = id.String()
	userId, err := h.Service.ServiceUser.CreateUser(r.Context(), user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "error:%v", errors.New(ErrCreateUser))
		return
	}

	tokenPair, err := h.GetNewToken(userId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "error:%v", err)
		return
	}
	_, err = h.Service.ServiceTokens.CreateTokens(r.Context(), tokenPair.RefreshToken)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "error:%v", errors.New(ErrCreateToken))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	body, err := json.Marshal(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "error:%v", errors.New(ErrBody))
		return
	}
	http.SetCookie(w, &http.Cookie{Name: "refresh_token", Value: tokenPair.RefreshToken.Token, HttpOnly: true})
	http.SetCookie(w, &http.Cookie{Name: "access_token", Value: tokenPair.AccessToken, HttpOnly: true})
	w.Write(body)
}

func (h *Handler) RefreshPairTokens(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "error:%v", errors.New(ErrCookie))
		return
	}
	id := r.Context().Value("id")
	ids, ok := id.(ID)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	ok, err = h.Validator(r.Context(), ids.TokenID, cookie.Value)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "error:%v", err)
		return
	}
	if !ok {
		http.SetCookie(w, &http.Cookie{Name: "refresh_token", Value: "", Expires: time.Now()})
		http.SetCookie(w, &http.Cookie{Name: "access_token", Value: "", Expires: time.Now()})

		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintf(w, "error:%v", errors.New(ErrToken))
		return
	}
	tokenPair, err := h.GetNewToken(ids.UserID)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "error:%v", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	http.SetCookie(w, &http.Cookie{Name: "refresh_token", Value: tokenPair.RefreshToken.Token, HttpOnly: true})
	http.SetCookie(w, &http.Cookie{Name: "access_token", Value: tokenPair.AccessToken, HttpOnly: true})
	_, err = h.Service.ServiceTokens.CreateTokens(r.Context(), tokenPair.RefreshToken)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "error:%v", errors.New(ErrCreateToken))
		return
	}
	h.Service.ServiceTokens.UpdateToken(r.Context(), ids.TokenID)
	body, err := json.Marshal(tokenPair)
	if err != nil {
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "error:%v", errors.New(ErrBody))
			return
		}
	}
	w.Write(body)
}

func (h *Handler) GetNewToken(id string) (core.TokenPair, error) {
	guid, err := guid.FromString(id)
	if err != nil {
		return core.TokenPair{}, errors.New("Проверьте id")
	}
	return h.Service.ServiceTokens.CreateTokenPairs(guid)
}

func (h *Handler) Validator(ctx context.Context, id string, token string) (bool, error) {
	rt, err := h.Service.ServiceTokens.GetTokens(ctx, id)
	if err != nil {
		return false, err
	}
	if hash.Check(rt.Token, token) {
		if !rt.IsUse && rt.ExpiresAt.After(time.Now().UTC()) {
			return true, nil
		}
	}
	return false, nil
}

func (h *Handler) TokenQuery(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	tokenPair, err := h.GetNewToken(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "error:%v", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	http.SetCookie(w, &http.Cookie{Name: "refresh_token", Value: tokenPair.RefreshToken.Token, HttpOnly: true})
	http.SetCookie(w, &http.Cookie{Name: "access_token", Value: tokenPair.AccessToken, HttpOnly: true})
	_, err = h.Service.ServiceTokens.CreateTokens(r.Context(), tokenPair.RefreshToken)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "error:%v", errors.New(ErrCreateToken))
		return
	}
	body, err := json.Marshal(tokenPair)
	w.Write(body)
}
