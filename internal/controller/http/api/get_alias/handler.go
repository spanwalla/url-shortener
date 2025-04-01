package get_alias

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/spanwalla/url-shortener/internal/controller/http/api"
	"github.com/spanwalla/url-shortener/internal/service"
)

type handler struct {
	service service.Expander
}

func New(service service.Expander) api.Handler {
	return &handler{service: service}
}

type getUriInput struct {
	Alias string `param:"alias" validate:"required,len=10"`
}

type getUriResponse struct {
	URI string `json:"uri" example:"https://github.com/spanwalla/url-shortener"`
}

// Handle handles HTTP-request
//
//	@Summary		Получить полную ссылку
//	@Description	Получает полную ссылку по сокращённому идентификатору.
//	@Param			alias	path	string	true	"Alias"	example(46g1B3ZgAy)	len(10)
//	@Produce		json
//	@Success		200	{object}	getUriResponse
//	@Failure		400	{object}	echo.HTTPError
//	@Failure		404	{object}	echo.HTTPError
//	@Failure		500	{object}	echo.HTTPError
//	@Router			/:alias [get]
func (h *handler) Handle(c echo.Context) error {
	var input getUriInput

	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}

	if err := c.Validate(input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	uri, err := h.service.Expand(c.Request().Context(), input.Alias)
	if err != nil {
		if errors.Is(err, service.ErrURINotFound) {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}

	return c.JSON(http.StatusOK, getUriResponse{uri})
}
