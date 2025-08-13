package http_server

import (
	"net/http"
	"svc-task_master/src/ports_adapters/primary/http_server/dto"
)

// GetTaskForId получает задачу по ID
// @Summary Получение задачи по ID
// @Description Возвращает задачу по указанному идентификатору
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path string true "ID задачи"
// @Success 200 {object} dto.Response{data=domain.Task} "Задача найдена"
// @Failure 400 {object} dto.Response "Некорректный ID задачи"
// @Failure 500 {object} dto.Response "Внутренняя ошибка сервера"
// @Router /task/{id} [get]
func (s Server) GetTaskForId(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(string)
	req := dto.GetTaskRequest{
		ID: id,
	}
	err := req.Validate()
	if err != nil {
		response(w, nil, http.StatusBadRequest, err)
		return
	}
	res, err := s.app.Query.GetTask.Handle(r.Context(), req)
	if err != nil {
		response(w, nil, http.StatusInternalServerError, err)
		return
	}
	response(w, res, http.StatusOK, nil)

}
