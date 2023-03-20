package restserver

import (
	"net/http"

	"github.com/benchkram/bobc/restserver/generated"
	"github.com/labstack/echo/v4"
)

func (s *S) GetHealth(ctx echo.Context) error {
	res := &generated.Success{
		Message: "Server successfully created",
	}
	return ctx.JSON(http.StatusOK, res)
}
