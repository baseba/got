package handler

import (
	"github.com/baseba/got/model"
	"github.com/baseba/got/view/userView"
	"github.com/labstack/echo/v4"
)

type UserHandler struct{}

func (h UserHandler) HandleUserShow(c echo.Context) error {
	user := model.User{
		Email: "user@user.com",
	}
	return render(c, userView.Show(user))
}
