package auth

import (
	"go-auth/src/db"
	"go-auth/src/dtos"
	"go-auth/src/models"
	"go-auth/src/utils"
	"net/http"
)

func getUsers() dtos.Response {
	conn, err := db.BorrowDbConnection()
	if err != nil {
		return utils.GetErrorResponse(http.StatusInternalServerError, "Something wrong!")
	}
	defer db.ReturnDbConnection(conn)

	users, err := db.Get[models.User](conn, `
	SELECT * FROM users
	`)
	if err != nil {
		return utils.GetErrorResponse(http.StatusInternalServerError, "Something Wrong!")
	}

	return utils.GetSuccessResponse(users)
}


func getUserById(id string) dtos.Response{
	conn, err := db.BorrowDbConnection()
	if err != nil {
		return utils.GetErrorResponse(http.StatusInternalServerError, "Something wrong!")
	}
	defer db.ReturnDbConnection(conn)
	
	user, err := db.GetFirst[models.User](conn, `
		SELECT * FROM users WHERE id = $1
	`, id)

	if err != nil{
		return utils.GetBadRequestResponse("Payload is not valid")
	}

	if (user == nil){
		return utils.GetErrorResponse(http.StatusNotFound, "User not found")
	}

	return utils.GetSuccessResponse(user)
}