package controllers

import (
	"database/sql"
	"fmt"
	"go-auth/src/models"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

func (me *Controller) getUsers(w http.ResponseWriter, _ *http.Request) {
	var users []models.User
	err := me.db.Select(&users, `
	SELECT * FROM users
	`)
	if err != nil {
		fmt.Println("ERR", err)
		sendDefaultInternalErrorResponse(w)
		return
	}

	sendSuccessResponse(w, users)
}

func (me *Controller) getUserById(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		fmt.Println("ERR", err)
		sendBadRequestResponse(w, "Payload is not valid")
		return
	}

	var user models.User
	err = me.db.Get(&user, `
		SELECT * FROM users WHERE id = ?
	`, id)

	if err != nil {
		if err == sql.ErrNoRows {
			sendBadRequestResponse(w, "User is not found")
			return
		}
		fmt.Println("ERR", err)
		sendDefaultInternalErrorResponse(w)
		return
	}

	sendSuccessResponse(w, user)
}
