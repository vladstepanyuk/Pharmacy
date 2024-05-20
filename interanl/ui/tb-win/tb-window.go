package tb_win

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"pharmacy/interanl/pharmacy"
)

func NewMenu(app fyne.App, m *pharmacy.Model) fyne.CanvasObject {
	c := container.NewVBox(
		widget.NewButton("Добавить технологию приготовления", addTechButtonCallback(app, m)),
		widget.NewButton("Обновить технологию", updateTechButtonCallback(app, m)),
		widget.NewButton("Удалить технологию", deleteTechButtonCallback(app, m)),
		widget.NewButton("Вывести все технологии", showTechsInStorageButtonCallback(app, m)),
		widget.NewButton("Задать технологию для лекарства", setTechForDrugButtonCallback(app, m)),
		widget.NewButton("Добавить компонент для технологии", setCompForTechButtonCallback(app, m)),
		widget.NewButton("Вывести все технологии фильтр", showTechsFilteredInStorageButtonCallback(app, m)),
		widget.NewButton("Вывести самые популярные компоненты", getMostPopulardrugs(app, m)),
	)

	return c
}
