package handler

import (
	"net/http"
	"ret/api/models"
	"ret/pkg/helpers"

	"github.com/gin-gonic/gin"
)

// CreateCountry godoc
// @Summary Create Country
// @Description Create Country
// @Tags Country
// @Accept json
// @Produce json
// @Param object body models.CreateCountry true "CreateCountryRequestBody"
// @Success 201 {object} Response{data=models.Country} "CountryBody"
// @Failure 400 {object} Response{data=string} "Invalid Argument"
// @Failure 500 {object} Response{data=string} "Server Error"
// @Router /country [post]
func (h *Handler) CreateCountry(c *gin.Context) {
	var country = models.CreateCountry{}
	err := c.ShouldBindJSON(&country)
	if err != nil {
		c.JSON(400, "ShouldBindJSON err:"+err.Error())
		return
	}
	resp, err := h.strg.Country().Create(country)
	if err != nil {
		handleResponse(c, http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, http.StatusCreated, resp)
}

// CountryGetById godoc
// @Summary Get Country by ID
// @Description Get Country by ID
// @Tags Country
// @Accept json
// @Produce json
// @Param id path string true "Country ID"
// @Success 200 {object} Response{data=models.Country} "CountryBody"
// @Failure 400 {object} Response{data=string} "Invalid Argument"
// @Failure 500 {object} Response{data=string} "Server Error"
// @Router /country/{id} [get]
func (h *Handler) CountryGetById(c *gin.Context) {
	var id = c.Param("id")
	if !helpers.IsValidUUID(id) {
		handleResponse(c, http.StatusBadRequest, "id is not uuid")
		return
	}

	resp, err := h.strg.Country().GetById(models.CountryPrimaryKey{Id: id})
	if err != nil {
		handleResponse(c, 500, "Country does not exist: "+err.Error())
		return
	}

	handleResponse(c, http.StatusOK, resp)
}

// CountryGetList godoc
// @Summary Get List of Countries
// @Description Get List of Countries
// @Tags Country
// @Accept json
// @Produce json
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {object} Response{data=models.GetListCountryResponse} "GetListCountryResponseBody"
// @Router /country [get]
func (h *Handler) CountryGetList(c *gin.Context) {
	var country models.GetListCountryRequest
	err := c.ShouldBindQuery(&country)
	if err != nil {
		handleResponse(c, http.StatusBadRequest, "Error while binding data: "+err.Error())
		return
	}
	resp, err := h.strg.Country().GetList(country)
	if err != nil {
		handleResponse(c, 500, "Country does not exist: "+err.Error())
		return
	}

	handleResponse(c, http.StatusOK, resp)
}

// CountryUpdate godoc
// @Router /country/{id} [put]
// @Summary Update Country
// @Description Update Country
// @Tags Country
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Param object body models.UpdateCountry true "UpdateCountryRequestBody"
// @Success 202 {string} string "Updated"
// @Failure 400 {object} Response{data=string} "Invalid Argument"
// @Failure 500 {object} Response{data=string} "Server Error"
func (h *Handler) CountryUpdate(c *gin.Context) {
	var country = models.UpdateCountry{}

	id := c.Param("id")
	if !helpers.IsValidUUID(id) {
		handleResponse(c, http.StatusBadRequest, "id not valid uuid")
		return
	}

	err := c.ShouldBindJSON(&country)
	if err != nil {
		handleResponse(c, http.StatusBadRequest, err)
		return
	}
	country.Guid = id

	resp, err := h.strg.Country().Update(country)
	if err != nil {
		handleResponse(c, 500, "Country does not update: "+err.Error())
		return
	}

	handleResponse(c, http.StatusAccepted, resp)
}

// CountryDelete godoc
// @Router /country/{id} [delete]
// @Summary Delete Country
// @Description Delete Country
// @Tags Country
// @Accept json
// @Produce json
// @Param id path string true "Country ID"
// @Success 204 {string} models.NoContent ""
// @Failure 400 {object} Response{data=string} "Invalid Argument"
// @Failure 500 {object} Response{data=string} "Server Error"
func (h *Handler) CountryDelete(c *gin.Context) {

	id := c.Param("id")
	if !helpers.IsValidUUID(id) {
		handleResponse(c, http.StatusBadRequest, "id not valid uuid")
		return
	}

	err := h.strg.Country().Delete(models.CountryPrimaryKey{Id: id})
	if err != nil {
		handleResponse(c, 500, "Country does not delete: "+err.Error())
		return
	}

	handleResponse(c, http.StatusNoContent, nil)
}

// UploadCountries godoc
// @Summary Загрузка стран
// @Description Загрузка стран из файла
// @Tags Country
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Файл JSON с городами"
// @Success 200 {string} string "Файл успешно загружен"
// @Failure 400 {object} Response{data=string} "Неверный аргумент"
// @Failure 500 {object} Response{data=string} "Ошибка сервера"
// @Router /upload/{table_slug} [post]
func (h *Handler) UploadCountry(c *gin.Context) {

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
	err = h.strg.Country().ImportFromFileCountry(filePath)
	if err != nil {
		handleResponse(c, http.StatusInternalServerError, "Ошибка при импорте данных: "+err.Error())
		return
	}

	handleResponse(c, http.StatusOK, "Файл успешно загружен")
}
