package handler

import (
	"net/http"
	"ret/api/models"
	"ret/pkg/helpers"

	"github.com/gin-gonic/gin"
)

// CreateAirport godoc
// @Summary Create Airport
// @Description Create Airport
// @Tags Airport
// @Accept json
// @Produce json
// @Param object body models.CreateAirport true "CreateAirportRequestBody"
// @Success 201 {object} Response{data=models.Airport} "AirportBody"
// @Failure 400 {object} Response{data=string} "Invalid Argument"
// @Failure 500 {object} Response{data=string} "Server Error"
// @Router /airport [post]
func (h *Handler) CreateAirport(c *gin.Context) {
	var airport = models.CreateAirport{}
	err := c.ShouldBindJSON(&airport)
	if err != nil {
		c.JSON(400, "ShouldBindJSON err:"+err.Error())
		return
	}
	resp, err := h.strg.Airport().Create(airport)
	if err != nil {
		handleResponse(c, http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, http.StatusCreated, resp)
}

// AirportGetById godoc
// @Summary Get Airport by ID
// @Description Get Airportby ID
// @Tags Airport
// @Accept json
// @Produce json
// @Param id path string true "Airport ID"
// @Success 200 {object} Response{data=models.Airport} "AirportBody"
// @Failure 400 {object} Response{data=string} "Invalid Argument"
// @Failure 500 {object} Response{data=string} "Server Error"
// @Router /airport/{id} [get]
func (h *Handler) AirportGetById(c *gin.Context) {
	var id = c.Param("id")
	if !helpers.IsValidUUID(id) {
		handleResponse(c, http.StatusBadRequest, "id is not uuid")
		return
	}

	resp, err := h.strg.Airport().GetById(models.AirportPrimaryKey{Id: id})
	if err != nil {
		handleResponse(c, 500, "Airport does not exist: "+err.Error())
		return
	}

	handleResponse(c, http.StatusOK, resp)
}

// AirportGetList godoc
// @Summary Get List of Airports
// @Description Get List of Airports
// @Tags Airport
// @Accept json
// @Produce json
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {object} Response{data=models.GetListAirportResponse} "GetListAirportResponseBody"
// @Router /airport [get]
func (h *Handler) AirportGetList(c *gin.Context) {
	var airport models.GetListAirportRequest
	err := c.ShouldBindQuery(&airport)
	if err != nil {
		handleResponse(c, http.StatusBadRequest, "Error while binding data: "+err.Error())
		return
	}
	resp, err := h.strg.Airport().GetList(airport)
	if err != nil {
		handleResponse(c, 500, "Airport does not exist: "+err.Error())
		return
	}

	handleResponse(c, http.StatusOK, resp)
}

// AirportUpdate godoc
// @Router /airport/{id} [put]
// @Summary Update Airport
// @Description Update Airport
// @Tags Airport
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Param object body models.UpdateAirport true "UpdateAirportRequestBody"
// @Success 202 {string} string "Updated"
// @Failure 400 {object} Response{data=string} "Invalid Argument"
// @Failure 500 {object} Response{data=string} "Server Error"
func (h *Handler) AirportUpdate(c *gin.Context) {
	var airport = models.UpdateAirport{}

	id := c.Param("id")
	if !helpers.IsValidUUID(id) {
		handleResponse(c, http.StatusBadRequest, "id not valid uuid")
		return
	}

	err := c.ShouldBindJSON(&airport)
	if err != nil {
		handleResponse(c, http.StatusBadRequest, err)
		return
	}
	airport.Id = id

	resp, err := h.strg.Airport().Update(airport)
	if err != nil {
		handleResponse(c, 500, "Airport does not update: "+err.Error())
		return
	}

	handleResponse(c, http.StatusAccepted, resp)
}

// AirportDelete godoc
// @Router /airport/{id} [delete]
// @Summary Delete Airport
// @Description Delete Airport
// @Tags Airport
// @Accept json
// @Produce json
// @Param id path string true "Airport ID"
// @Success 204 {string} models.NoContent ""
// @Failure 400 {object} Response{data=string} "Invalid Argument"
// @Failure 500 {object} Response{data=string} "Server Error"
func (h *Handler) AirportDelete(c *gin.Context) {

	id := c.Param("id")
	if !helpers.IsValidUUID(id) {
		handleResponse(c, http.StatusBadRequest, "id not valid uuid")
		return
	}

	err := h.strg.Airport().Delete(models.AirportPrimaryKey{Id: id})
	if err != nil {
		handleResponse(c, 500, "airport does not delete: "+err.Error())
		return
	}

	handleResponse(c, http.StatusNoContent, nil)
}






// UploadAirports godoc
// @Summary Загрузка аэропортов
// @Description Загрузка аэропортов из файла
// @Tags Airport
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Файл JSON с аэропортами"
// @Success 200 {string} string "Файл успешно загружен"
// @Failure 400 {object} Response{data=string} "Неверный аргумент"
// @Failure 500 {object} Response{data=string} "Ошибка сервера"
// @Router /upload/airport [post]
func (h *Handler) UploadAirport(c *gin.Context) {

	file, err := c.FormFile("file")
	if err != nil {
		handleResponse(c, http.StatusBadRequest, "Ошибка при получении файла: "+err.Error())
		return
	}

	if file.Header.Get("Content-Type") != "application/json" {
		handleResponse(c, http.StatusBadRequest, "Неверный формат файла. загрузите файл JSON.")
		return
	}

	uploadPath := "uploads/"
	err = c.SaveUploadedFile(file, uploadPath+file.Filename)
	if err != nil {
		handleResponse(c, http.StatusInternalServerError, "Ошибка при сохранении файла: "+err.Error())
		return
	}

	filePath := uploadPath + file.Filename
	err = h.strg.Airport().ImportFromFileAirport(filePath)
	if err != nil {
		handleResponse(c, http.StatusInternalServerError, "Ошибка при импорте данных: "+err.Error())
		return
	}

	handleResponse(c, http.StatusOK, "Файл успешно загружен")
}
