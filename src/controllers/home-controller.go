package controllers

import "net/http"

func (me *ChiController) home(w http.ResponseWriter, r *http.Request) {
	sendSuccessResponse(w, "hello, world!")
}

func (me *ChiController) ping(w http.ResponseWriter, _ *http.Request) {
	sendSuccessResponse(w, "pong")
}
