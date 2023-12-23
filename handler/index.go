package handler

import (
	"github.com/baseba/got/view/indexView"
	"github.com/labstack/echo/v4"
)

type IndexHandler struct{}

func (h IndexHandler) HandleIndexShow(c echo.Context) error {
	return render(c, indexView.Show())
}
