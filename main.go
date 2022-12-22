package main

import (
	"encoding/json"
	"fmt"
	"go-auth/src/config"
	"go-auth/src/controllers"
	"go-auth/src/dtos"
	"go-auth/src/utils"
	"io"
	"net/http"
	"reflect"
)

func updateProduct(id string, data dtos.CreateProduct) dtos.Response {
	fmt.Println("Id:", id)
	fmt.Println("Data:", data)
	return utils.GetSuccessResponse(data)
}

func getCallParams(r *http.Request, refFunc interface{}) []reflect.Value {
	refType := reflect.TypeOf(refFunc)
	var argTypes []reflect.Type
	for i := 0; i < refType.NumIn(); i++ {
		argType := refType.In(i)
		argTypes = append(argTypes, argType)
	}

	jsonString, _ := io.ReadAll(r.Body)
	idx := 0
	urlParams := []string{}
	var callParams []reflect.Value
	for _, v := range argTypes {
		if v.Kind() != reflect.Struct {
			if idx < len(urlParams) {
				callParams = append(callParams, reflect.ValueOf(urlParams[idx]))
				idx++
			} else {
				callParams = append(callParams, reflect.ValueOf(""))
			}
			continue
		}

		temp := reflect.New(v).Interface()
		_ = json.Unmarshal([]byte(jsonString), &temp)
		callParams = append(callParams, reflect.ValueOf(temp).Elem())
	}

	return callParams
}

func main() {
	router := controllers.NewRouter()
	router.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		params := getCallParams(r, updateProduct)
		response := reflect.ValueOf(updateProduct).Call(params)[0].Interface().(dtos.Response)

		w.Header().Add("content-type", "application/json")
		w.WriteHeader(response.Status)
		err := json.NewEncoder(w).Encode(response)
		if err != nil {
			fmt.Println("Failed to send response", err)
		}
	}).Methods("POST")

	fmt.Println("Server is listening on port", config.LISTEN_PORT)
	http.ListenAndServe("127.0.0.1:"+config.LISTEN_PORT, router)
}
