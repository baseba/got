package handler

import (
	"github.com/baseba/got/view/slotView"
	"github.com/labstack/echo/v4"
)

type SlotsHandler struct{}

func (h SlotsHandler) HandleSlotsShow(c echo.Context) error {
	if c.Param("room") != "" {

		return render(c, slotView.Show(c.Param("room")))
	}
	return render(c, slotView.Show("0"))
}
