package orders_win

import (
	"context"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"pharmacy/interanl/pharmacy"
	"pharmacy/interanl/pharmacy/daos"
	"strconv"
	"time"
)

func getTableForDrugs(m *pharmacy.Model, drugs map[daos.ID]daos.Recipe, ids []int) fyne.CanvasObject {
	return widget.NewTable(
		func() (int, int) {
			return len(drugs), 7
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("wideeeeeeee text")
		},
		func(id widget.TableCellID, object fyne.CanvasObject) {
			label := object.(*widget.Label)
			r := drugs[ids[id.Row]]
			switch id.Col {
			case 0:
				label.SetText(strconv.Itoa(ids[id.Row]))
			case 1:
				label.SetText(strconv.Itoa(r.ConsumerId))
			case 2:
				label.SetText(strconv.Itoa(r.DoctorId))
			case 3:
				label.SetText(strconv.Itoa(r.DrugId))
			case 4:
				label.SetText(strconv.Itoa(r.UseId))
			case 5:
				label.SetText(r.Disease)
			case 6:
				label.SetText(strconv.Itoa(r.Count))
			}
		},
	)

}

func getListButtonCallback(app fyne.App, m *pharmacy.Model) func() {
	return func() {
		ctx, f := context.WithTimeout(context.Background(), time.Second)
		defer f()
		dao, err := m.GetDaoByType(daos.RecipeE)
		if err != nil {
			fyne.LogError(err.Error(), err)
			return
		}

		cs := make(map[daos.ID]daos.Recipe)

		err = dao.List(ctx, &cs)
		if err != nil {
			fyne.LogError(err.Error(), err)
			return
		}

		ids := make([]daos.ID, 0, len(cs))
		for id := range cs {
			ids = append(ids, id)
		}
		table := getTableForDrugs(m, cs, ids)

		newW := app.NewWindow("list")
		newW.Resize(fyne.NewSize(400, 400))
		newW.SetContent(table)
		newW.Show()
	}
}
