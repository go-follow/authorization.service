package transport

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/go-follow/authorization.service/internal/service"
	"github.com/go-follow/authorization.service/internal/service/model"
	"github.com/go-follow/authorization.service/pkg/server"
)

//HTTP - http transport
type HTTP struct {
	svc service.Service
}

//NewHTTP - инициализация http transport, со всеми методами
func NewHTTP(svc service.Service, r chi.Router) {
	h := HTTP{svc}
	r.Group(func(ropen chi.Router) {
		ropen.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
			render.JSON(w, r, map[string]string{"success": "pong"})
		})
		ropen.Get("/check_access", h.checkAccess)
		ropen.Get("/login", h.login)
		ropen.Get("/logout", h.logout)
		ropen.Post("/add_user", h.addUser)

	})
}
func (h *HTTP) checkAccess(w http.ResponseWriter, r *http.Request) {
	token, err := r.Cookie("token")
	if err != nil {
		server.SendErrorJSON(w, r, 401, fmt.Errorf("токен авторизации не действителен"))
		return
	}
	role, err := h.svc.CheckAccess(token.Value)
	if err != nil {
		server.SendErrorJSON(w, r, 401, fmt.Errorf("токен авторизации не действителен"))
		return
	}
	render.JSON(w, r, map[string]string{"user_role": role})
}

func (h *HTTP) login(w http.ResponseWriter, r *http.Request) {
	username, password, _ := r.BasicAuth()
	status, token, err := h.svc.Login(username, password)
	if err != nil {
		server.SendErrorJSON(w, r, status, err)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		Expires:  time.Now().Add(time.Hour * 3),
		MaxAge:   10000,
		HttpOnly: true,
	})
}

func (h *HTTP) logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   "",
		Path:    "/",
		Expires: time.Unix(0, 0),
		//MaxAge:   10000,
		HttpOnly: true,
	})
}

func (h *HTTP) addUser(w http.ResponseWriter, r *http.Request) {
	token, err := r.Cookie("token")
	if err != nil {
		server.SendErrorJSON(w, r, 400, fmt.Errorf("токен авторизации не действителен"))
		return
	}
	role, err := h.svc.CheckAccess(token.Value)
	if err != nil {
		server.SendErrorJSON(w, r, 401, fmt.Errorf("токен авторизации не действителен"))
		return
	}
	if role != "admin" {
		server.SendErrorJSON(w, r, 403, fmt.Errorf("недостаточно прав для выполнения запроса"))
		return
	}
	enable, err := strconv.ParseBool(r.FormValue("enable"))
	if err != nil {
		enable = false
	}
	u := &model.User{
		Name:     r.FormValue("user_name"),
		Password: r.FormValue("password"),
		Enable:   enable,
		Role:     r.FormValue("role"),
	}

	status, _, err := h.svc.AddUser(u)
	if err != nil {
		server.SendErrorJSON(w, r, status, err)
		return
	}
}