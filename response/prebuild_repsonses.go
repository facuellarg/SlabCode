package response

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

//SomthingWrongResponse default response
func SomthingWrongResponse(err error) (int, echo.Map) {
	return http.StatusBadGateway, echo.Map{
		"message": SOMETHING_WRONG,
		"error":   err.Error(),
	}
}

//AlreadyExistResponse default response
func AlreadyExistResponse(table string) (int, echo.Map) {
	return http.StatusConflict, echo.Map{
		"message": fmt.Sprintf(ALREADY_EXISTS, table),
	}
}

//ForbiddenResourceResponse default response
func ForbiddenResourceResponse() (int, echo.Map) {
	return http.StatusConflict, echo.Map{
		"message": fmt.Sprintf(ALREADY_EXISTS, FORBIDDEN_RESOURCE),
	}
}

//BadRequestResponse default response
func BadRequestResponse() (int, echo.Map) {
	return http.StatusBadRequest, echo.Map{
		"message": BAD_REQUEST,
	}
}
