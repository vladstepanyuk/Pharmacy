package co_win

import (
	"context"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"pharmacy/interanl/pharmacy"
	"pharmacy/interanl/pharmacy/daos"
	"pharmacy/interanl/ui/utils"
	"strconv"
	"time"
)

func getLateConsumers(app fyne.App, m *pharmacy.Model) func() {
	return func() {
		w := app.NewWindow("")
		w.Resize(fyne.NewSize(400, 200))

		ctx, f := context.WithTimeout(context.Background(), time.Second)
		defer f()

		cs, err := m.GetLateConsumers(ctx)
		if err != nil {
			w.SetContent(widget.NewLabel(err.Error()))
			w.Show()
			return
		}
		ids := make([]daos.ID, 0, len(cs))
		for id := range cs {
			ids = append(ids, id)
		}

		table := utils.GetTableForConsumers(cs, ids)

		w.SetContent(table)
		w.Show()
	}
}

func getConsumersWhoWait(app fyne.App, m *pharmacy.Model) func() {
	return func() {
		newW := app.NewWindow("")
		newW.Resize(fyne.NewSize(400, 200))

		tnames := pharmacy.StdTypesNames()
		checkboxes := make([]fyne.CanvasObject, 0, len(tnames))
		ids := make([][]daos.ID, len(tnames))
		for i, tn := range tnames {
			checkboxes = append(checkboxes, widget.NewCheck(tn, func(b bool) {}))
			idForType, err := pharmacy.GetTypesIDByName(m, tn)
			if err != nil {
				panic(err)
			}

			ids[i] = idForType
		}

		statusText := widget.NewLabel("")

		newW.SetContent(container.NewGridWithRows(3,
			container.NewAdaptiveGrid(len(tnames), checkboxes...),
			widget.NewButton("загрузить", func() {
				idForReq := make([]daos.ID, 0)
				for i := range tnames {
					if checkboxes[i].(*widget.Check).Checked {
						idForReq = append(idForReq, ids[i]...)
					}
				}

				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()

				cs, err := m.GetWaitingConsumers(ctx, idForReq)
				if err != nil {
					statusText.SetText(err.Error())
					return
				}

				ids := make([]daos.ID, 0, len(cs))
				for id := range cs {
					ids = append(ids, id)
				}
				table := utils.GetTableForConsumers(cs, ids)

				w := app.NewWindow("")
				w.Resize(fyne.NewSize(400, 400))
				w.SetContent(table)
				w.Show()

			}),
			statusText,
		),
		)

		newW.Show()
	}
}

func getBuyersFromPeriod(app fyne.App, m *pharmacy.Model) func() {
	return func() {
		typeField := widget.NewEntry()
		typeField.SetPlaceHolder("Type Field")
		dIdField := widget.NewEntry()
		dIdField.SetPlaceHolder("Drug Id Field")
		beginField := widget.NewEntry()
		beginField.SetPlaceHolder("Begin: YYYY-MM-DD")
		endField := widget.NewEntry()
		endField.SetPlaceHolder("End: YYYY-MM-DD")

		w := app.NewWindow("")
		w.Resize(fyne.NewSize(400, 200))
		statusText := widget.NewLabel("")
		w.SetContent(container.NewGridWithRows(3,
			container.NewAdaptiveGrid(4, typeField, dIdField, beginField, endField),
			widget.NewButton("загрузить", func() {
				dId, err := strconv.Atoi(dIdField.Text)
				if err != nil {
					statusText.SetText(err.Error())
					return
				}
				begin, err := time.Parse("2006-01-02", beginField.Text)
				if err != nil {
					statusText.SetText(err.Error())
					return
				}
				end, err := time.Parse("2006-01-02", endField.Text)
				if err != nil {
					statusText.SetText(err.Error())
					return
				}

				tIds, err := pharmacy.GetTypesIDByName(m, typeField.Text)
				if err != nil {
					statusText.SetText(err.Error())
					return
				}
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()

				cs, err := m.GetConsumersOfDrugs(ctx, []daos.ID{dId}, tIds, begin, end)
				if err != nil {
					statusText.SetText(err.Error())
					return
				}

				ids := make([]daos.ID, 0, len(cs))
				for id := range cs {
					ids = append(ids, id)
				}
				table := utils.GetTableForConsumers(cs, ids)

				w := app.NewWindow("")
				w.Resize(fyne.NewSize(400, 400))
				w.SetContent(table)
				w.Show()

			}),
			statusText,
		),
		)

		w.Show()

	}
}

func getOrdersInProgress(app fyne.App, m *pharmacy.Model) func() {
	return func() {
		w := app.NewWindow("")
		w.Resize(fyne.NewSize(400, 400))

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		orders, err := m.GetInProgressOrders(ctx)
		if err != nil {
			w.SetContent(widget.NewLabel(err.Error()))
			w.Show()
			return
		}
		ids := make([]daos.ID, 0, len(orders))
		for id := range orders {
			ids = append(ids, id)
		}

		table := utils.GetTableForOrders(orders, ids)

		w.SetContent(table)
		w.Show()

	}
}

func getCompsForOrdersInProgress(app fyne.App, m *pharmacy.Model) func() {
	return func() {
		w := app.NewWindow("")
		w.Resize(fyne.NewSize(400, 400))

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		drugs, err := m.GetComponentsForInProgressOrders(ctx)
		if err != nil {
			w.SetContent(widget.NewLabel(err.Error()))
			w.Show()
			return
		}
		ids := make([]daos.ID, 0, len(drugs))
		for id := range drugs {
			ids = append(ids, id)
		}

		table := utils.GetTableForDrugs(m, drugs, ids)

		w.SetContent(table)
		w.Show()

	}
}

func getBestConsumers(app fyne.App, m *pharmacy.Model) func() {
	return func() {
		typeField := widget.NewEntry()
		typeField.SetPlaceHolder("Type Field")
		dIdField := widget.NewEntry()
		dIdField.SetPlaceHolder("Drug Id Field")

		w := app.NewWindow("")
		w.Resize(fyne.NewSize(400, 200))
		statusText := widget.NewLabel("")
		w.SetContent(container.NewGridWithRows(3,
			container.NewAdaptiveGrid(2, typeField, dIdField),
			widget.NewButton("загрузить", func() {
				dId, err := strconv.Atoi(dIdField.Text)
				if err != nil {
					statusText.SetText(err.Error())
					return
				}

				tIds, err := pharmacy.GetTypesIDByName(m, typeField.Text)
				if err != nil {
					statusText.SetText(err.Error())
					return
				}
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()

				csc, err := m.GetMostOftenConsumers(ctx, []daos.ID{dId}, tIds)
				if err != nil {
					statusText.SetText(err.Error())
					return
				}

				ids := make([]daos.ID, 0, len(csc))
				for id := range csc {
					ids = append(ids, id)
				}

				cs := make(map[daos.ID]daos.Consumer)

				for id, cc := range csc {
					cs[id] = daos.Consumer{Name: cc.Name, Phone: cc.Phone, Address: cc.Address, BirthDate: cc.BirthDate}
				}

				table := utils.GetTableForConsumers(cs, ids)

				w := app.NewWindow("")
				w.Resize(fyne.NewSize(400, 400))
				w.SetContent(table)
				w.Show()

			}),
			statusText,
		),
		)

		w.Show()

	}

}
