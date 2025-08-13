package http_server

import (
	"encoding/json"
	"net/http"
	"svc-task_master/src/ports_adapters/primary/http_server/dto"
)

// UpdateStatusTask обновляет статус задачи
// @Summary Обновление статуса задачи
// @Description Обновляет статус задачи по указанному идентификатору
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path string true "ID задачи"
// @Param task body dto.UpdateTaskStatusRequest true "Данные для обновления статуса"
// @Success 200 {object} dto.Response{data=domain.Task} "Статус задачи обновлен"
// @Failure 400 {object} dto.Response "Некорректные данные запроса"
// @Failure 500 {object} dto.Response "Внутренняя ошибка сервера"
// @Router /task/{id} [put]
func (s Server) UpdateStatusTask(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(string)
	var req dto.UpdateTaskStatusRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		response(w, nil, http.StatusBadRequest, err)
	}
	req.Id = id

	err = req.Validate()
	if err != nil {
		response(w, nil, http.StatusBadRequest, err)
		return
	}
	res, err := s.app.Command.UpdateTask.Handle(r.Context(), req)
	if err != nil {
		response(w, nil, http.StatusInternalServerError, err)
		return
	}
	response(w, res, http.StatusOK, nil)

}
