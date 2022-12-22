package main

import (
	"encoding/json"
	"fmt"
	"go-auth/src/dtos"
	"reflect"
	// "net/http"
	// "reflect"
	// "go-auth/src/config"
	// "go-auth/src/controllers"
)

func updateProduct(id string, data dtos.CreateProduct) {
	fmt.Println("Id:", id)
	fmt.Println("Data:", data)
}

func main() {
	var ref interface{} = updateProduct
	refType := reflect.TypeOf(ref)
	var argTypes []reflect.Type
	for i := 0; i < refType.NumIn(); i++ {
		argType := refType.In(i)
		argTypes = append(argTypes, argType)
	}

	jsonString := `
				{
					"name" : "Keyboard",
					"description": "Keyboard baru dari paris",
					"price" : 10000
				}
			`
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
	// fmt.Printf("Data: %+v\n", callParams[1])
	reflect.ValueOf(ref).Call(callParams)

	// router := controllers.NewRouter()

	// fmt.Println("Server is listening on port", config.LISTEN_PORT)
	// http.ListenAndServe("127.0.0.1:"+config.LISTEN_PORT, router)
}
