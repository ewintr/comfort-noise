package player

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type API struct {
	radio *Player
}

func NewAPI(radio *Player) *API {
	return &API{
		radio: radio,
	}
}

func (a *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		a.Status(w, r)
	case http.MethodPost:
		a.Select(w, r)
	default:
		http.NotFound(w, r)
	}

}

func (a *API) Status(w http.ResponseWriter, r *http.Request) {
	body, err := json.Marshal(a.radio.Status())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	fmt.Fprintf(w, string(body))
}

func (a *API) Select(w http.ResponseWriter, r *http.Request) {
	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	command := struct {
		Action  string `json:"action"`
		Channel int    `json:"channel"`
	}{}
	if err := json.Unmarshal(reqBody, &command); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if command.Action == "play" {
		if err := a.radio.Select(command.Channel); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	} else {
		a.radio.Stop()
	}

	a.Status(w, r)
}
