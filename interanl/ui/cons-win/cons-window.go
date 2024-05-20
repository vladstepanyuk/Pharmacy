package cons_win

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"pharmacy/interanl/pharmacy"
)

func NewMenu(app fyne.App, m *pharmacy.Model) fyne.CanvasObject {
	c := container.NewVBox(
		widget.NewButton("Добавить клиента", addConsumerCallback(app, m)),
		widget.NewButton("Все клиенты", getListButtonCallback(app, m)),
		widget.NewButton("Обновить клиентf", updateConsumerCallback(app, m)),
		widget.NewButton("Удалить клиента", deleteConsumerCallback(app, m)),
	)

	return c
}
