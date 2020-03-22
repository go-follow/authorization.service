package server

import (
	"net/http"

	"github.com/go-chi/render"
)

//SendErrorJSON - отправка JSON c ошибкой и статусом отличным от 200
func SendErrorJSON(w http.ResponseWriter, r *http.Request, httpCode int, err error) {
	render.Status(r, httpCode)
	render.JSON(w, r, map[string]string{"error": messageFromError(err)})
}

func messageFromError(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}