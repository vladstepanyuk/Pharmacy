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

func getTableForDrugs(m *pharmacy.Model, drugs map[daos.ID]pharmacy.DrugsWithCount, ids []int) fyne.CanvasObject {
	return widget.NewTable(
		func() (int, int) {
			return len(drugs), 6
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("wideeeeee text")
		},
		func(id widget.TableCellID, object fyne.CanvasObject) {
			label := object.(*widget.Label)
			drug := drugs[ids[id.Row]]
			switch id.Col {
			case 0:
				label.SetText(strconv.Itoa(ids[id.Row]))
			case 1:
				label.SetText(drug.Name)
			case 2:
				label.SetText(pharmacy.GetTypeNameByID(m, drug.TypeId))
			case 3:
				label.SetText(strconv.Itoa(drug.Price))
			case 4:
				label.SetText(strconv.Itoa(drug.CriticalCount))
			case 5:
				label.SetText(strconv.Itoa(drug.Count))
			}
		},
	)

}

func showDrugsInStorageButtonCallback(app fyne.App, m *pharmacy.Model) func() {

	return func() {

		ctx, f := context.WithTimeout(context.Background(), time.Second)
		defer f()

		drugsC, err := m.ShowStorage(ctx)
		if err != nil {
			w := app.NewWindow("error")
			w.SetContent(widget.NewLabel(err.Error()))
			w.Show()
			return
		}

		ids := make([]daos.ID, 0, len(drugsC))
		for id := range drugsC {
			ids = append(ids, id)
		}

		table := getTableForDrugs(m, drugsC, ids)

		newW := app.NewWindow("list drugs")
		newW.Resize(fyne.NewSize(400, 400))
		newW.SetContent(table)
		newW.Show()

	}
}

func showDrugsAmountButtonCallback(app fyne.App, m *pharmacy.Model) func() {
	return func() {

		idField := widget.NewEntry()
		idField.SetPlaceHolder("id")
		beginField := widget.NewEntry()
		beginField.SetPlaceHolder("Начало: YYYY-MM-DD")
		endField := widget.NewEntry()
		endField.SetPlaceHolder("Конец: YYYY-MM-DD")

		w := app.NewWindow("drugs amount")
		statusText := widget.NewLabel("")

		b := widget.NewButton("загрузить", func() {
			begin, err := time.Parse("2006-01-02", beginField.Text)
			if err != nil {
				statusText.SetText("ОЩИБКА: у начала периода невалидный формат")
				return
			}

			end, err := time.Parse("2006-01-02", endField.Text)
			if err != nil {
				statusText.SetText("ОЩИБКА: у конца периода невалидный формат")
				return
			}
			id, err := strconv.Atoi(idField.Text)
			if err != nil {
				statusText.SetText("ОЩИБКА: id не число")
				return
			}

			ctx, f := context.WithTimeout(context.Background(), time.Second)
			defer f()

			used, err := m.GetUsedAmmount(ctx, []daos.ID{id}, begin, end)
			if err != nil {
				statusText.SetText("ОЩИБКА: " + err.Error())
				return
			}

			statusText.SetText("Использовано: " + strconv.Itoa(used[id]))
		})

		w.SetContent(container.NewGridWithRows(3,
			container.NewAdaptiveGrid(3, idField, beginField, endField),
			b,
			statusText,
		))

		w.Resize(fyne.NewSize(400, 400))
		w.Show()
	}

}
