package http_server

import (
	"net/http"
	"svc-task_master/src/ports_adapters/primary/http_server/dto"
)

// GetTasksSortStatus получает список задач с фильтрацией по статусу
// @Summary Получение списка задач
// @Description Возвращает список задач с возможностью фильтрации по статусу
// @Tags tasks
// @Accept json
// @Produce json
// @Param status query string false "Статус для фильтрации (pending, processing, completed, failed, retrying)"
// @Success 200 {object} dto.Response{data=[]domain.Task} "Список задач получен"
// @Failure 400 {object} dto.Response "Некорректные параметры запроса"
// @Failure 500 {object} dto.Response "Внутренняя ошибка сервера"
// @Router /task [get]
func (s Server) GetTasksSortStatus(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	req := dto.GetTaskWhithFiltersRequest{
		Status: status,
	}
	err := req.Validate()
	if err != nil {
		response(w, nil, http.StatusBadRequest, err)
		return
	}
	res, err := s.app.Query.GetTasks.Handle(r.Context(), req)
	if err != nil {
		response(w, nil, http.StatusInternalServerError, err)
		return
	}
	response(w, res, http.StatusOK, nil)

}
