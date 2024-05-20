package tb_win

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

func showTechsInStorageButtonCallback(app fyne.App, m *pharmacy.Model) func() {

	return func() {

		ctx, f := context.WithTimeout(context.Background(), time.Second)
		defer f()

		dao, err := m.GetDaoByType(daos.TechnologyBookE)
		if err != nil {
			panic(err)
		}

		techs := make(map[daos.ID]daos.TechnologyBook)
		err = dao.List(ctx, &techs)
		if err != nil {
			w := app.NewWindow("error")
			w.SetContent(widget.NewLabel(err.Error()))
			w.Show()
			return
		}

		ids := make([]daos.ID, 0, len(techs))
		for id := range techs {
			ids = append(ids, id)
		}

		table := utils.GetTableForTechs(techs, ids)

		newW := app.NewWindow("list techs")
		newW.Resize(fyne.NewSize(400, 400))
		newW.SetContent(table)
		newW.Show()

	}
}

func showTechsFilteredInStorageButtonCallback(app fyne.App, m *pharmacy.Model) func() {

	return func() {
		newW := app.NewWindow("")
		newW.Resize(fyne.NewSize(400, 200))

		idField := widget.NewEntry()
		idField.SetPlaceHolder("id")

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
		checkboxes = append(checkboxes, idField)
		statusText := widget.NewLabel("")

		newW.SetContent(container.NewGridWithRows(3,
			container.NewAdaptiveGrid(len(tnames)+1, checkboxes...),
			widget.NewButton("загрузить", func() {
				idForReq := make([]daos.ID, 0)
				for i := range tnames {
					if checkboxes[i].(*widget.Check).Checked {
						idForReq = append(idForReq, ids[i]...)
					}
				}

				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()

				id, err := strconv.Atoi(idField.Text)
				if err != nil {
					statusText.SetText(err.Error())
					return
				}

				techs, err := m.GetTechsForDrugs(ctx, []daos.ID{id}, idForReq)
				if err != nil {
					statusText.SetText(err.Error())
					return
				}

				ids := make([]daos.ID, 0, len(techs))
				for id := range techs {
					ids = append(ids, id)
				}
				table := utils.GetTableForTechs(techs, ids)

				w := app.NewWindow("list techs")
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

func getMostPopulardrugs(app fyne.App, m *pharmacy.Model) func() {

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

				drugsC, err := m.GetTenMostPopularDrugs(ctx, idForReq)
				if err != nil {
					statusText.SetText(err.Error())
					return
				}

				ids := make([]daos.ID, 0, len(drugsC))
				for id := range drugsC {
					ids = append(ids, id)
				}

				drugs := make(map[daos.ID]daos.Drug, len(drugsC))
				for id, dc := range drugsC {
					drugs[id] = daos.Drug{Name: dc.Name, TypeId: dc.TypeId, CriticalCount: dc.CriticalCount, Price: dc.Price}
				}
				table := utils.GetTableForDrugs(m, drugs, ids)

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
