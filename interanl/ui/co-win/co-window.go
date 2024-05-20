package co_win

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"pharmacy/interanl/pharmacy"
)

func NewMenu(app fyne.App, m *pharmacy.Model) fyne.CanvasObject {
	c := container.NewVBox(
		widget.NewButton("Добавить заказ", addOrder(app, m)),
		widget.NewButton("Удалить заказ", updateOrder(app, m)),
		widget.NewButton("Обновить заказ", deleteOrder(app, m)),
		widget.NewButton("Опоздавшие клиенты", getLateConsumers(app, m)),
		widget.NewButton("Ждущие клиенты", getConsumersWhoWait(app, m)),
		widget.NewButton("Покупатели лекарств за период", getBuyersFromPeriod(app, m)),
		widget.NewButton("Получить заказы находящиеся в производстве", getOrdersInProgress(app, m)),
		widget.NewButton("Получить препараты нужные для заказов в производстве", getCompsForOrdersInProgress(app, m)),
		widget.NewButton("Получить самых преданых покупателей", getBestConsumers(app, m)),
	)

	return c
}
