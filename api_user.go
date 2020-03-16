package main

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strings"

	"github.com/gorilla/mux"
)

const (
	UserApiPath       = "/users"
	UserApiPathWithId = "/users/{userId}"
	FormatPassword    = `^\S{4,45}$` // без пробелов, от 4 до 45 символов
)

type UserApi struct{}

func (api UserApi) Path() string {
	return UserApiPath
}

func (api UserApi) PathWithId() string {
	return UserApiPathWithId
}

func (api UserApi) One(writer http.ResponseWriter, request *http.Request) {
	statusCode := GetDefaultStatus(request.Method)
	userToResponse := new(User)
	defer ReturnJson(writer, userToResponse, statusCode)

	pathVars := mux.Vars(request)
	userToResponse.Id = GetIdFromPath(pathVars, "userId")

	err := userToResponse.Get()
	userToResponse.Password = nil
	if err != nil {
		*statusCode = http.StatusNotFound
		return
	}
}

func (api UserApi) Many(writer http.ResponseWriter, request *http.Request) {
	statusCode := GetDefaultStatus(request.Method)
	type UsersResponse []User
	response := new(UsersResponse)
	defer ReturnJson(writer, response, statusCode)

	limit := GetBodyIntVariable(request, "limit", DefaultLimit)
	offset := GetBodyIntVariable(request, "offset", DefaultOffset)
	order := GetBodyStringVariable(request, "order")

	where := map[string]string{}

	err := GetAllEntities(limit, offset, "", where, order, response)

	if err != nil {
		return
	}
}

func (api UserApi) Create(writer http.ResponseWriter, request *http.Request) {
	statusCode := GetDefaultStatus(request.Method)
	user := new(User)
	defer ReturnJson(writer, user, statusCode)

	err := json.NewDecoder(request.Body).Decode(user)
	if err != nil {
		*statusCode = http.StatusBadRequest
		return
	}
	user.Email = strings.ToLower(user.Email)
	if user.Login == "" && user.Email != "" {
		user.Login = user.Email
	}

	isValidPassword, _ := regexp.MatchString(FormatPassword, *user.Password)
	if !isValidPassword {
		*statusCode = http.StatusBadRequest
		return
	}

	err = user.Create()
	if err != nil {
		*statusCode = http.StatusBadRequest
		return
	}
}

func (api UserApi) Update(writer http.ResponseWriter, request *http.Request) {
	statusCode := GetDefaultStatus(request.Method)
	defer ReturnJson(writer, nil, statusCode)

	user := new(User)
	err := json.NewDecoder(request.Body).Decode(&user)
	if err != nil {
		*statusCode = http.StatusBadRequest
		return
	}
	user.Email = strings.ToLower(user.Email)
	pathVars := mux.Vars(request)
	userId := GetIdFromPath(pathVars, "userId")

	if user.Id != userId {
		*statusCode = http.StatusBadRequest
		return
	}
	err = user.Update()
	if err != nil {
		if err.Error() == "incorrect current password" {
			*statusCode = http.StatusBadRequest
			return
		}
	}
}
