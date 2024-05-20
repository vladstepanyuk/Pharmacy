package drug_win

import (
	"context"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"pharmacy/interanl/pharmacy"
	"pharmacy/interanl/pharmacy/daos"
	"pharmacy/interanl/ui/utils"
	"strconv"
	"strings"
	"time"
)

func getListButtonCallback(app fyne.App, m *pharmacy.Model) func() {
	return func() {
		ctx, f := context.WithTimeout(context.Background(), time.Second)
		defer f()
		dao, err := m.GetDaoByType(daos.DrugE)
		if err != nil {
			fyne.LogError(err.Error(), err)
			return
		}

		drugs := make(map[daos.ID]daos.Drug)

		err = dao.List(ctx, &drugs)
		if err != nil {
			fyne.LogError(err.Error(), err)
			return
		}

		ids := make([]daos.ID, 0, len(drugs))
		for id := range drugs {
			ids = append(ids, id)
		}
		table := utils.GetTableForDrugs(m, drugs, ids)

		newW := app.NewWindow("list drugs")
		newW.Resize(fyne.NewSize(400, 400))
		newW.SetContent(table)
		newW.Show()
	}
}

func renderDetailedInfo(app fyne.App, m *pharmacy.Model, ddi pharmacy.DrugDetail) {
	newW := app.NewWindow(ddi.Name)
	newW.Resize(fyne.NewSize(400, 200))

	var b strings.Builder

	b.WriteString("Название лекарства:\n\t")
	b.WriteString(ddi.Name)
	b.WriteString("\nТип лекарства:\n\t")
	b.WriteString(pharmacy.GetTypeNameByID(m, ddi.TypeId))
	b.WriteString("\nЦена:\n\t")
	b.WriteString(strconv.Itoa(ddi.Price))
	b.WriteString("\nКритическое количество:\n\t")
	b.WriteString(strconv.Itoa(ddi.CriticalCount))
	b.WriteString("\nСпособ приготовления:\n\t")
	b.WriteString(ddi.TechDescription)
	b.WriteString("\nКомпоненты:")

	ids := make([]daos.ID, 0, len(ddi.Components))
	//ddi.Components[0].
	for id := range ddi.Components {
		ids = append(ids, id)
	}

	components := widget.NewTable(
		func() (int, int) {
			return len(ids), 4
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("wideeeeeeee text")
		},
		func(id widget.TableCellID, object fyne.CanvasObject) {
			label := object.(*widget.Label)
			comp := ddi.Components[ids[id.Row]]
			switch id.Col {
			case 0:
				label.SetText(strconv.Itoa(ids[id.Row]))
			case 1:
				label.SetText(comp.Name)
			case 2:
				label.SetText(strconv.Itoa(comp.Price))
			case 3:
				label.SetText(strconv.Itoa(comp.Count))
			case 4:
			}
		},
	)

	newW.SetContent(container.NewGridWithRows(2,
		widget.NewLabel(b.String()),
		components,
	))
	newW.Show()
}

func getDetailedInfoWindow(app fyne.App, m *pharmacy.Model) fyne.CanvasObject {

	idField := widget.NewEntry()
	idField.SetPlaceHolder("введите айди")

	statusText := widget.NewLabel("")

	return container.NewGridWithRows(
		3,
		idField,
		widget.NewButton("найти", func() {

			id, err := strconv.ParseUint(idField.Text, 10, 64)
			if err != nil {
				statusText.SetText("error: id не число")
				return
			}

			ctx, f := context.WithTimeout(context.Background(), time.Second)
			defer f()

			detailedInf, err := m.GetDrugDetailInfo(ctx, daos.ID(id))

			if err != nil {
				fyne.LogError("произошла ошибка", err)
				statusText.SetText("error: " + err.Error())
				return
			}
			renderDetailedInfo(app, m, detailedInf)
			statusText.SetText("ok")

		}),
		statusText,
	)

}

func getDetailedInfoButtonCallback(app fyne.App, m *pharmacy.Model) func() {
	return func() {

		newW := app.NewWindow("search drug")
		newW.Resize(fyne.NewSize(400, 200))

		newW.SetContent(
			getDetailedInfoWindow(app, m),
		)

		newW.Show()
	}
}

func getRunOutDrugsWindow(m *pharmacy.Model) fyne.CanvasObject {
	ctx, f := context.WithTimeout(context.Background(), time.Second)
	defer f()
	drugs, err := m.GetCriticalDrugs(ctx)
	if err != nil {
		label := widget.NewLabel("ОШИБКА: " + err.Error())
		return label
	}

	ids := make([]daos.ID, 0, len(drugs))
	for id := range drugs {
		ids = append(ids, id)
	}
	table := utils.GetTableForDrugs(m, drugs, ids)

	return table
}

func getRunOutDrugsButtonCallback(app fyne.App, m *pharmacy.Model) func() {
	return func() {

		newW := app.NewWindow("run out drugs")
		newW.Resize(fyne.NewSize(400, 200))

		newW.SetContent(
			getRunOutDrugsWindow(m),
		)

		newW.Show()

	}
}

func getMinimalDrugsButtonCallback(app fyne.App, m *pharmacy.Model) func() {
	return func() {

		newW := app.NewWindow("выберите категории")
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

		newW.SetContent(container.NewGridWithRows(
			3,
			container.NewAdaptiveGrid(len(checkboxes), checkboxes...),
			widget.NewButton("загрузить", func() {
				idForReq := make([]daos.ID, 0)
				for i, cb := range checkboxes {
					if cb.(*widget.Check).Checked {
						idForReq = append(idForReq, ids[i]...)
					}
				}

				ctx, f := context.WithTimeout(context.Background(), time.Second)
				defer f()

				fmt.Println(idForReq)
				drugs, err := m.GetDrugsWithMinStorageCount(ctx, idForReq)
				if err != nil {
					statusText.SetText("ОШИБКА: " + err.Error())
					return
				}

				dids := make([]daos.ID, 0, len(drugs))
				for id := range drugs {
					dids = append(dids, id)
				}
				table := utils.GetTableForDrugs(m, drugs, dids)

				ww := app.NewWindow("table")
				ww.SetContent(table)
				ww.Show()
				statusText.SetText("ok")
			}),
			statusText,
		))

		newW.Show()

	}
}
