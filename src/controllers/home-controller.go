package controllers

import "net/http"

func (me *Controller) home(w http.ResponseWriter, r *http.Request) {
	sendSuccessResponse(w, "hello, world!")
}

func (me *Controller) ping(w http.ResponseWriter, _ *http.Request) {
	sendSuccessResponse(w, "pong")
}
