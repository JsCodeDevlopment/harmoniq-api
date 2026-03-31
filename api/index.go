package handler

import (
	"api/src/app"
	"net/http"
)

var engine = app.Bootstrap()

func Handler(w http.ResponseWriter, r *http.Request) {
	engine.ServeHTTP(w, r)
}
