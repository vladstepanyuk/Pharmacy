package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"pharmacy/interanl/pharmacy"
	co_win "pharmacy/interanl/ui/co-win"
	cons_win "pharmacy/interanl/ui/cons-win"
	"pharmacy/interanl/ui/drug-win"
	orders_win "pharmacy/interanl/ui/orders-win"
	storage_win "pharmacy/interanl/ui/storage-win"
	tb_win "pharmacy/interanl/ui/tb-win"
)

func NewMainWindow(app fyne.App, m *pharmacy.Model) fyne.Window {
	myWindow := app.NewWindow("Table Widget")

	tab := container.NewAppTabs(
		container.NewTabItem("лекарства", drug_win.NewMenu(app, m)),
		container.NewTabItem("склад", storage_win.NewMenu(app, m)),
		container.NewTabItem("справочник технологий", tb_win.NewMenu(app, m)),
		container.NewTabItem("справочник клиентов", cons_win.NewMenu(app, m)),
		container.NewTabItem("рецепты", orders_win.NewMenu(app, m)),
		container.NewTabItem("заказы", co_win.NewMenu(app, m)),
	)

	myWindow.Resize(fyne.NewSize(800, 800))
	myWindow.SetContent(tab)

	return myWindow
}
