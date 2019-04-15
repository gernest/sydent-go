package matrixid

import (
	"net/http"

	"github.com/labstack/echo"
)

func Version(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, map[string]interface{}{})
}
