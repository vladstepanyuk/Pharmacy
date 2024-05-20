package storage_win

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"pharmacy/interanl/pharmacy"
)

func NewMenu(app fyne.App, m *pharmacy.Model) fyne.CanvasObject {
	c := container.NewVBox(
		widget.NewButton("Добавить товар на склад", addDrugsToStorageButtonCallback(app, m)),
		widget.NewButton("Списать товар со склада", takeDrugsOfStorageButtonCallback(app, m)),
		widget.NewButton("Показать товары на складе", showDrugsInStorageButtonCallback(app, m)),
		widget.NewButton("Показать расходы за период", showDrugsAmountButtonCallback(app, m)),
	)

	return c
}
