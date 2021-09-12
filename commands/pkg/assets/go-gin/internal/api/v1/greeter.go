package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// @Summary Greeter service
// @Id 1
// @version 1.0
// @produce application/json
// @Param name query string true "Input name"
// @Success 200 {object} GreeterResponse
// @Router /v1/greeter [get]
func Greeter(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, &GreeterResponse{
		Message: fmt.Sprintf("Hello %s!", ctx.Query("name")),
	})
}

// GreeterResponse Response of /v1/greeter.
type GreeterResponse struct {
	Message string
}
