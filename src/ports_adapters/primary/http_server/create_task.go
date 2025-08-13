package http_server

import (
	"encoding/json"
	"net/http"
	"svc-task_master/src/ports_adapters/primary/http_server/dto"
)

// CreateTask создает новую задачу
// @Summary Создание новой задачи
// @Description Создает новую задачу в системе
// @Tags tasks
// @Accept json
// @Produce json
// @Param task body dto.TaskRequest true "Данные для создания задачи"
// @Success 200 {object} dto.Response{data=domain.Task} "Задача успешно создана"
// @Failure 400 {object} dto.Response "Некорректные данные запроса"
// @Failure 500 {object} dto.Response "Внутренняя ошибка сервера"
// @Router /task [post]
func (s Server) CreateTask(w http.ResponseWriter, r *http.Request) {
	var req dto.TaskRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		response(w, nil, http.StatusBadRequest, err)
		return
	}
	err = req.Validate()
	if err != nil {
		response(w, nil, http.StatusBadRequest, err)
		return
	}
	res, err := s.app.Command.CreateTask.Handle(r.Context(), req)
	if err != nil {
		response(w, nil, http.StatusInternalServerError, err)
		return
	}
	response(w, res, http.StatusOK, nil)

}
