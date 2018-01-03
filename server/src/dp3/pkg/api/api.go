package api

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"goji.io"
	"goji.io/pat"
	"io/ioutil"
	"log"
	"net/http"
)

// Mux creates the API router and returns it for inclusion in the app router
func Mux() *goji.Mux {

	apiMux := goji.SubMux()

	version1Mux := goji.SubMux()
	version1Mux.HandleFunc(pat.Post("/issues"), submitIssueHandler)
	version1Mux.HandleFunc(pat.Get("/swagger.yml"), swaggerYAMLHandler)
	apiMux.Handle(pat.New("/v1/*"), version1Mux)

	return apiMux
}

// Incoming body for POST /issues
type issue struct {
	Issue string `json:"issue"`
}

// Response to POST /issues
type newIssueResponse struct {
	ID int64 `json:"id"`
}

func swaggerYAMLHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("WEOIFNWEOIFNWIOEFNOWINF")
	files, err := ioutil.ReadDir("./")
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		fmt.Println(f.Name())
	}
	http.ServeFile(w, r, "./swagger.yml")
}

func submitIssueHandler(w http.ResponseWriter, r *http.Request) {
	var newIssue issue

	if err := json.NewDecoder(r.Body).Decode(&newIssue); err != nil {
		zap.L().Error("Json decode", zap.Error(err))
		http.Error(w, http.StatusText(400), http.StatusBadRequest)
	} else {
		resp := newIssueResponse{1}
		responseJSON, err := json.Marshal(resp)

		if err != nil {
			zap.L().Error("Encode message", zap.Error(err))
			http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		} else {
			w.Write(responseJSON)
		}
	}
}
