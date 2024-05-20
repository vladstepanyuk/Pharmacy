package drug_win

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"pharmacy/interanl/pharmacy"
)

func NewMenu(app fyne.App, m *pharmacy.Model) fyne.CanvasObject {
	c := container.NewVBox(
		widget.NewButton("Все лекарства", getListButtonCallback(app, m)),
		widget.NewButton("Добавить лекарство", getInsertButtonCallback(app, m)),
		widget.NewButton("Обновить информацию о лекарстве", getUpdateButtonCallback(app, m)),
		widget.NewButton("Удалить информацию о лекарстве", getDeleteButtonCallback(app, m)),
		widget.NewButton("Детальная информация о лекарстве", getDetailedInfoButtonCallback(app, m)),
		widget.NewButton("Получить заканчивающиеся лекарства", getRunOutDrugsButtonCallback(app, m)),
		widget.NewButton("Получить лекарствв с минимальным запасом", getMinimalDrugsButtonCallback(app, m)),
	)

	return c
}
