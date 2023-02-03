package beego

import (
	"github.com/beego/beego/v2/server/web"
)

type UserController struct {
	web.Controller
}

func (c *UserController) GetUser() {
	c.Ctx.WriteString("this is hello from Hud!")
}

func (c *UserController) CreateUser() {
	u := &User{}
	err := c.Ctx.BindJSON(u)
	if err != nil {
		c.Ctx.WriteString(err.Error())
		return
	}

	_ = c.Ctx.JSONResp(u)
}

type User struct {
	Name string
}
