package storage_win

import (
	"context"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"pharmacy/interanl/pharmacy"
	"pharmacy/interanl/pharmacy/daos"
	"strconv"
	"time"
)

func addDrugsToStorageButtonCallback(app fyne.App, m *pharmacy.Model) func() {
	return func() {
		w := app.NewWindow("добавить на склад")

		idField := widget.NewEntry()
		idField.SetPlaceHolder("ID")
		countField := widget.NewEntry()
		countField.SetPlaceHolder("Count")

		statusText := widget.NewLabel("")

		w.SetContent(container.NewGridWithRows(3,
			container.NewAdaptiveGrid(2, idField, countField),
			widget.NewButton("применить", func() {
				id, err := strconv.ParseUint(idField.Text, 10, 64)
				if err != nil {
					statusText.SetText("ОШИБКА: id не число")
					return
				}
				c, err := strconv.ParseUint(countField.Text, 10, 64)
				if err != nil {
					statusText.SetText("ОШИБКА: Count не число")
					return
				}

				ctx, f := context.WithTimeout(context.Background(), time.Second)
				defer f()

				err = m.AddDrugsToStorage(ctx, daos.ID(id), int(c))
				if err != nil {
					statusText.SetText("ОШИБКА: " + err.Error())
					return
				}
				statusText.SetText("ok")
			}),
			statusText,
		))
		w.Show()
	}
}

func takeDrugsOfStorageButtonCallback(app fyne.App, m *pharmacy.Model) func() {
	return func() {
		w := app.NewWindow("списать со склада")

		idField := widget.NewEntry()
		idField.SetPlaceHolder("ID")
		countField := widget.NewEntry()
		countField.SetPlaceHolder("Count")

		statusText := widget.NewLabel("")

		w.SetContent(container.NewGridWithRows(3,
			container.NewAdaptiveGrid(2, idField, countField),
			widget.NewButton("применить", func() {
				id, err := strconv.ParseUint(idField.Text, 10, 64)
				if err != nil {
					statusText.SetText("ОШИБКА: id не число")
					return
				}
				c, err := strconv.ParseUint(countField.Text, 10, 64)
				if err != nil {
					statusText.SetText("ОШИБКА: Count не число")
					return
				}

				ctx, f := context.WithTimeout(context.Background(), time.Second)
				defer f()

				err = m.TakeOutOfStorage(ctx, daos.ID(id), int(c))
				if err != nil {
					statusText.SetText("ОШИБКА: " + err.Error())
					return
				}
				statusText.SetText("ok")
			}),
			statusText,
		))
		w.Show()
	}
}
