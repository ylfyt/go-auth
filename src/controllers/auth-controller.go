package controllers

import (
	"encoding/json"
	"fmt"
	"go-auth/src/db"
	"go-auth/src/dtos"
	"go-auth/src/models"
	"go-auth/src/services"
	"io"
	"net/http"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func Register(r *http.Request) dtos.Response {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return getErrorResponse(http.StatusBadRequest, "Failed to get request body")
	}

	var data dtos.Register
	err = json.Unmarshal(body, &data)
	if err != nil {
		return getErrorResponse(http.StatusBadRequest, "Failed to get payload")
	}

	if len(data.Username) < 4 || len(data.Username) > 20 {
		return getErrorResponse(http.StatusBadRequest, "Username should be > 4 and < 20")
	}

	if len(data.Password) < 4 {
		return getErrorResponse(http.StatusBadRequest, "Password should be > 4")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
	if err != nil {
		return getErrorResponse(http.StatusInternalServerError, "Something wrong!")
	}

	conn, err := db.BorrowDbConnection()
	if err != nil {
		return getErrorResponse(http.StatusInternalServerError, "Something wrong!")
	}
	defer db.ReturnDbConnection(conn)

	exists, err := db.GetRowCount(*conn, `
		SELECT count(*) FROM users WHERE username = $1
	`, data.Username)
	if err != nil {
		fmt.Println("Error", err)
		return getErrorResponse(http.StatusInternalServerError, "Something wrong")
	}

	if exists != 0 {
		return getErrorResponse(http.StatusBadRequest, "Username already exist")
	}

	newId := uuid.New()

	inserted, err := db.Write(*conn, `
		INSERT INTO users VALUES($1, $2, $3, NOW())
	`, newId, data.Username, string(hashedPassword))

	if err != nil {
		fmt.Println("Error:", err)
		return getErrorResponse(http.StatusInternalServerError, "Failed to insert new user")
	}

	if inserted == 0 {
		return getErrorResponse(http.StatusInternalServerError, "Failed to insert new user")
	}

	return getSuccessResponse(newId)
}

func Login(r *http.Request) dtos.Response {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return getErrorResponse(http.StatusBadRequest, "Failed to get request body")
	}

	var data dtos.Register
	err = json.Unmarshal(body, &data)
	if err != nil {
		return getErrorResponse(http.StatusBadRequest, "Failed to get payload")
	}

	if len(data.Username) < 4 || len(data.Username) > 20 {
		return getErrorResponse(http.StatusBadRequest, "Username should be > 4 and < 20")
	}

	if len(data.Password) < 4 {
		return getErrorResponse(http.StatusBadRequest, "Password should be > 4")
	}

	conn, err := db.BorrowDbConnection()
	if err != nil {
		return getErrorResponse(http.StatusInternalServerError, "Something wrong!")
	}
	defer db.ReturnDbConnection(conn)

	user, err := db.GetFirst[models.User](*conn, `
	SELECT * FROM users WHERE username = $1
	`, data.Username)
	if err != nil {
		return getErrorResponse(http.StatusInternalServerError, "Something wrong!")
	}

	if user == nil {
		return getErrorResponse(http.StatusBadRequest, "Username or password is wrong")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.Password))
	if err != nil {
		return getErrorResponse(http.StatusBadRequest, "Username or password is wrong")
	}

	refresh, access, jid, err := services.CreateJwtToken(*user)
	if err != nil {
		return getErrorResponse(http.StatusInternalServerError, "Something wrong!")
	}
	inserted, _ := db.Write(*conn, `
		INSERT INTO jwt_tokens VALUES($1, $2, NOW())
	`, jid, user.Id)
	if inserted == 0 {
		return getErrorResponse(http.StatusInternalServerError, "Something wrong!")
	}

	return getSuccessResponse(dtos.LoginResponse{
		User: *user,
		Token: dtos.TokenPayload{
			RefreshToken: refresh,
			AccessToken:  access,
		},
	})
}

func RefreshToken(r *http.Request) dtos.Response {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return getErrorResponse(http.StatusBadRequest, "Failed to get request body")
	}

	var data dtos.RefreshPayload
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println("Error", err)
		return getErrorResponse(http.StatusBadRequest, "Failed to get payload")
	}

	valid, jid := services.VerifyRefreshToken(data.RefreshToken)
	if !valid {
		return getErrorResponse(http.StatusBadRequest, "Token is not valid")
	}

	conn, err := db.BorrowDbConnection()
	if err != nil {
		return getErrorResponse(http.StatusInternalServerError, "Something wrong!")
	}
	defer db.ReturnDbConnection(conn)

	token, err := db.GetFirst[models.JwtToken](*conn, `
		SELECT * FROM jwt_tokens WHERE id = $1
	`, jid)
	if err != nil {
		return getErrorResponse(http.StatusInternalServerError, "Something wrong!")
	}
	if token == nil {
		return getErrorResponse(http.StatusBadRequest, "Token is not found")
	}

	user, err := db.GetFirst[models.User](*conn, `
	SELECT * FROM users WHERE id = $1
	`, token.UserId)
	if err != nil {
		return getErrorResponse(http.StatusInternalServerError, "Something wrong!")
	}

	refresh, access, newJid, err := services.CreateJwtToken(*user)
	if err != nil {
		return getErrorResponse(http.StatusInternalServerError, "Something wrong!")
	}
	inserted, _ := db.Write(*conn, `
		INSERT INTO jwt_tokens VALUES($1, $2, NOW())
	`, newJid, user.Id)
	if inserted == 0 {
		return getErrorResponse(http.StatusInternalServerError, "Something wrong!")
	}

	deleted, _ := db.Write(*conn, `
	DELETE FROM jwt_tokens WHERE id = $1
	`, jid)
	if deleted == 0 {
		return getErrorResponse(http.StatusInternalServerError, "Something wrong!")
	}

	return getSuccessResponse(dtos.TokenPayload{
		RefreshToken: refresh,
		AccessToken:  access,
	})
}

func Logout(r *http.Request) dtos.Response {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return getErrorResponse(http.StatusBadRequest, "Failed to get request body")
	}

	var data dtos.RefreshPayload
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println("Error", err)
		return getErrorResponse(http.StatusBadRequest, "Failed to get payload")
	}

	valid, jid := services.VerifyRefreshToken(data.RefreshToken)
	if !valid {
		return getErrorResponse(http.StatusBadRequest, "Token is not valid")
	}

	conn, err := db.BorrowDbConnection()
	if err != nil {
		return getErrorResponse(http.StatusInternalServerError, "Something wrong!")
	}
	defer db.ReturnDbConnection(conn)

	deleted, err := db.Write(*conn, `
		DELETE FROM jwt_tokens WHERE id = $1
	`, jid)
	if err != nil {
		fmt.Println("Err", err)
		return getErrorResponse(http.StatusInternalServerError, "Something wrong!")
	}

	if deleted == 0 {
		return getErrorResponse(http.StatusBadRequest, "Token is not found")
	}

	// Implement Blacklist Token Mechanism

	return getSuccessResponse(true)
}
