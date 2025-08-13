package http_server

import (
	"encoding/json"
	"net/http"
	"svc-task_master/src/application"
	"svc-task_master/src/ports_adapters/primary/http_server/dto"
)

type Server struct {
	app *application.App
}

func NewServer(app *application.App) *Server {
	return &Server{
		app: app,
	}
}

func response(w http.ResponseWriter, data any, status int, err error) {
	res := dto.Response{}
	if err != nil {
		errorr := err.Error()
		res.Error = &errorr
	}
	if data != nil {
		res.Data = data
	}
	res.Status = status
	body, _ := json.Marshal(res)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	w.Write(body)
}
