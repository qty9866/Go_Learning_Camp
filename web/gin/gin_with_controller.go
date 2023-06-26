package gin

import (
	"github.com/gin-gonic/gin"
)

type UserController struct{}

func (c *UserController) GetUser(ctx *gin.Context) {
	//panic("something wrong!")
	ctx.String(200, "this is Hello form Hud")
}
