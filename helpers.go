package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

const DefaultLimit = 20
const DefaultOffset = 0

func GetIdFromPath(pathVars map[string]string, varName string) int {
	Id, _ := strconv.ParseInt(pathVars[varName], 10, 64)
	return int(Id)
}

func GetBodyIntVariable(request *http.Request, varName string, defaultValue int) int {
	variable, err := strconv.ParseInt(request.URL.Query().Get(varName), 10, 64)
	if err != nil {
		return defaultValue
	}
	return int(variable)
}

func GetBodyStringVariable(request *http.Request, varName string) string {
	variable := request.URL.Query().Get(varName)
	return variable
}

func ReturnJson(writer http.ResponseWriter, response interface{}, statusCode *int) {
	if *statusCode > 299 || *statusCode == 204 {
		writer.WriteHeader(*statusCode)
		return
	}
	writer.WriteHeader(*statusCode)
	if reflect.TypeOf(response).Elem().Kind() == reflect.Slice && reflect.ValueOf(response).Elem().Len() < 1 {
		response = make([]int, 0)
	}
	if err := json.NewEncoder(writer).Encode(response); err != nil {
		panic(err)
	}
}

func GetDefaultStatus(method string) *int {
	statusCode := new(int)
	switch strings.ToUpper(method) {
	case http.MethodGet:
		*statusCode = http.StatusOK
	case http.MethodPost:
		*statusCode = http.StatusCreated
	case http.MethodPut:
		*statusCode = http.StatusNoContent
	case http.MethodDelete:
		*statusCode = http.StatusNoContent
	default:
		*statusCode = http.StatusNotImplemented
	}
	return statusCode
}

func GetMD5Hash(text string) string {
	hasher := md5.New()
	_, _ = hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}
