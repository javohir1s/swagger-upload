package handler

import (
	"net/http"
	"ret/api/models"
	"ret/pkg/helpers"

	"github.com/gin-gonic/gin"
)

// CreateCity godoc
// @Summary Create City
// @Description Create City
// @Tags City
// @Accept json
// @Produce json
// @Param object body models.CreateCity true "CreateCityRequestBody"
// @Success 201 {object} Response{data=models.City} "CityBody"
// @Failure 400 {object} Response{data=string} "Invalid Argument"
// @Failure 500 {object} Response{data=string} "Server Error"
// @Router /city [post]
func (h *Handler) CreateCity(c *gin.Context) {
	var city = models.CreateCity{}
	err := c.ShouldBindJSON(&city)
	if err != nil {
		c.JSON(400, "ShouldBindJSON err:"+err.Error())
		return
	}

	resp, err := h.strg.City().Create(city)
	if err != nil {
		handleResponse(c, http.StatusBadRequest, "Does not create"+err.Error())
		return
	}
	handleResponse(c, http.StatusCreated, resp)
}

// City GetById godoc
// @Summary Get City  by ID
// @Description Get City  by ID
// @Tags City
// @Accept json
// @Produce json
// @Param id path string true "City  ID"
// @Success 200 {object} Response{data=models.City} "City Body"
// @Failure 400 {object} Response{data=string} "Invalid Argument"
// @Failure 500 {object} Response{data=string} "Server Error"
// @Router /city/{id} [get]
func (h *Handler) CityGetById(c *gin.Context) {
	var id = c.Param("id")
	if !helpers.IsValidUUID(id) {
		handleResponse(c, http.StatusBadRequest, "id is not uuid")
		return
	}

	resp, err := h.strg.City().GetById(models.CityPrimaryKey{Id: id})
	if err != nil {
		handleResponse(c, 500, "City does not exist: "+err.Error())
		return
	}

	handleResponse(c, http.StatusOK, resp)
}

// CityGetList godoc
// @Summary Get List of cities
// @Description Get List of cities
// @Tags City
// @Accept json
// @Produce json
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {object} Response{data=models.GetListCityResponse} "GetListCityResponseBody"
// @Router /city [get]
func (h *Handler) CityGetList(c *gin.Context) {
	var city models.GetListCityRequest
	err := c.ShouldBindQuery(&city)
	if err != nil {
		handleResponse(c, http.StatusBadRequest, "Error while binding data: "+err.Error())
		return
	}
	resp, err := h.strg.City().GetList(city)
	if err != nil {
		handleResponse(c, 500, "city does not exist: "+err.Error())
		return
	}

	handleResponse(c, http.StatusOK, resp)
}

// CityUpdate godoc
// @Router /city/{id} [put]
// @Summary Update City
// @Description Update City
// @Tags City
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Param object body models.UpdateCity true "UpdateCityRequestBody"
// @Success 202 {string} string "Updated"
// @Failure 400 {object} Response{data=string} "Invalid Argument"
// @Failure 500 {object} Response{data=string} "Server Error"
func (h *Handler) CityUpdate(c *gin.Context) {
	var city = models.UpdateCity{}

	id := c.Param("id")
	if !helpers.IsValidUUID(id) {
		handleResponse(c, http.StatusBadRequest, "id not valid uuid")
		return
	}

	err := c.ShouldBindJSON(&city)
	if err != nil {
		handleResponse(c, http.StatusBadRequest, err)
		return
	}
	city.Guid = id

	resp, err := h.strg.City().Update(city)
	if err != nil {
		handleResponse(c, 500, "city does not update: "+err.Error())
		return
	}

	handleResponse(c, http.StatusAccepted, resp)
}

// CityDelete godoc
// @Router /city/{id} [delete]
// @Summary Delete City
// @Description Delete City
// @Tags City
// @Accept json
// @Produce json
// @Param id path string true "City ID"
// @Success 204 {string} models.NoContent ""
// @Failure 400 {object} Response{data=string} "Invalid Argument"
// @Failure 500 {object} Response{data=string} "Server Error"
func (h *Handler) CityDelete(c *gin.Context) {
	id := c.Param("id")
	if !helpers.IsValidUUID(id) {
		handleResponse(c, http.StatusBadRequest, "id not valid uuid")
		return
	}

	err := h.strg.City().Delete(models.CityPrimaryKey{Id: id})
	if err != nil {
		handleResponse(c, 500, "City does not delete: "+err.Error())
		return
	}

	handleResponse(c, http.StatusNoContent, nil)
}

// UploadCities godoc
// @Summary Загрузка городов
// @Description Загрузка городов из файла
// @Tags City
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Файл JSON с городами"
// @Success 200 {string} string "Файл успешно загружен"
// @Failure 400 {object} Response{data=string} "Неверный аргумент"
// @Failure 500 {object} Response{data=string} "Ошибка сервера"
// @Router /upload [post]
func (h *Handler) UploadCities(c *gin.Context) {

	file, err := c.FormFile("file")
	if err != nil {
		handleResponse(c, http.StatusBadRequest, "Ошибка при получении файла: "+err.Error())
		return
	}

	if file.Header.Get("Content-Type") != "application/json" {
		handleResponse(c, http.StatusBadRequest, "Неверный формат файла.  загрузите файл JSON.")
		return
	}

	uploadPath := "uploads/"

	err = c.SaveUploadedFile(file, uploadPath+file.Filename)
	if err != nil {
		handleResponse(c, http.StatusInternalServerError, "Ошибка при сохранении файла: "+err.Error())
		return
	}

	filePath := uploadPath + file.Filename

	err = h.strg.City().ImportFromFile(filePath)
	if err != nil {
		handleResponse(c, http.StatusInternalServerError, "Ошибка при импорте данных: "+err.Error())
		return
	}

	handleResponse(c, http.StatusOK, "Файл успешно загружен ")
}
