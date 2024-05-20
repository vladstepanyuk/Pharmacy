package orders_win

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"pharmacy/interanl/pharmacy"
)

func NewMenu(app fyne.App, m *pharmacy.Model) fyne.CanvasObject {
	c := container.NewVBox(
		widget.NewButton("Добавить рецепт", addConsumerCallback(app, m)),
		widget.NewButton("Все рецепт", getListButtonCallback(app, m)),
		widget.NewButton("Обновить рецепт", updateConsumerCallback(app, m)),
		widget.NewButton("Удалить рецепт", deleteConsumerCallback(app, m)),
	)

	return c
}
