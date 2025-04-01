package post_root

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/spanwalla/url-shortener/internal/controller/http/api"
	"github.com/spanwalla/url-shortener/internal/service"
)

type handler struct {
	service service.Shortener
}

func New(service service.Shortener) api.Handler {
	return &handler{service: service}
}

type createAliasInput struct {
	URI string `json:"uri" validate:"required,uri" example:"https://google.com"`
}

type createAliasResponse struct {
	Alias string `json:"alias" example:"46g1B3ZgAy"`
}

// Handle handles HTTP-request
//
//	@Summary		Получить сокращённую ссылку
//	@Description	Получает сокращённую ссылку из полной.
//	@Param			uri	body	createAliasInput	true	"Информация об URI"
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	createAliasResponse
//	@Success		201	{object}	createAliasResponse
//	@Failure		400	{object}	echo.HTTPError
//	@Failure		500	{object}	echo.HTTPError
//	@Router			/ [post]
func (h *handler) Handle(c echo.Context) error {
	var input createAliasInput

	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}

	if err := c.Validate(input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	alias, created, err := h.service.Shorten(c.Request().Context(), input.URI)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}

	responseCode := http.StatusOK
	if created {
		responseCode = http.StatusCreated
	}

	return c.JSON(responseCode, createAliasResponse{alias})
}
