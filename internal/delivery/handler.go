package delivery

import (
	"fmt"
	"net/http"

	"github.com/jamascrorpJS/auth/internal/service"
)

type Handler struct {
	Mux     *http.ServeMux
	Service *service.ServiceRepository
}

func NewHandler(service *service.ServiceRepository) *Handler {
	mux := http.NewServeMux()
	return &Handler{
		Mux:     mux,
		Service: service,
	}
}

func (h *Handler) Start() {
	h.Mux.HandleFunc("/", h.Router)
}

func (h *Handler) Router(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		{
			if r.URL.Path == "/refresh" {
				AccessToken(h.RefreshPairTokens)(w, r)
			} else if r.URL.Path == "/create" {
				h.TokenQuery(w, r)
			} else {
				ifNotExist(w)
				return
			}
		}
	case http.MethodPost:
		{
			if r.URL.Path == "/register" {
				h.Register(w, r)
			} else {
				ifNotExist(w)
				return
			}
		}
	default:
		{
			ifNotExist(w)
			return
		}
	}
}
func ifNotExist(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, "Такого пути нет")
}
