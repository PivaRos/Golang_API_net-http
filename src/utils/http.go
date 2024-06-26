package utils

import (
	"log"
	"net/http"

	"github.com/go-playground/validator"
	"go.mongodb.org/mongo-driver/mongo"
)

func HandleError(w http.ResponseWriter, err error) {
	var errorCode int
	switch err.(type) {
	case mongo.CommandError:
		if err == mongo.ErrNoDocuments {
			errorCode = http.StatusNotFound
		} else {
			errorCode = http.StatusInternalServerError
		}
	case validator.ValidationErrors:
		errorCode = http.StatusBadRequest
	default:
		errorCode = http.StatusInternalServerError
	}
	log.Println("error is :", err.Error())
	http.Error(w, err.Error(), errorCode)
}
