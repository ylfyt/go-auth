package controllers

import (
	"fmt"
	"go-auth/src/models"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

func (me *Controller) getUsers(w http.ResponseWriter, _ *http.Request) {
	var users []models.User
	err := me.db.Get(&users, `
	SELECT * FROM users
	`)
	if err != nil {
		sendDefaultInternalErrorResponse(w)
		return
	}

	sendSuccessResponse(w, users)
}

func (me *Controller) getUserById(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		fmt.Println("ERR", err)
		sendBadRequestResponse(w, "Payload is not valid")
		return
	}

	var user *models.User
	err = me.db.GetFirst(&user, `
		SELECT * FROM users WHERE id = $1
	`, id)

	if err != nil {
		sendDefaultInternalErrorResponse(w)
		return
	}

	if user == nil {
		sendBadRequestResponse(w, "User not found")
		return
	}

	sendSuccessResponse(w, user)
}
