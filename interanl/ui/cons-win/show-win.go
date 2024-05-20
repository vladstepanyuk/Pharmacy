package cons_win

import (
	"context"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"pharmacy/interanl/pharmacy"
	"pharmacy/interanl/pharmacy/daos"
	"strconv"
	"time"
)

func getTableForConsumers(drugs map[daos.ID]daos.Consumer, ids []int) fyne.CanvasObject {
	return widget.NewTable(
		func() (int, int) {
			return len(drugs), 5
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("wideeeeeeee text")
		},
		func(id widget.TableCellID, object fyne.CanvasObject) {
			label := object.(*widget.Label)
			c := drugs[ids[id.Row]]
			switch id.Col {
			case 0:
				label.SetText(strconv.Itoa(ids[id.Row]))
			case 1:
				label.SetText(c.Name)
			case 2:
				label.SetText(c.Phone)
			case 3:
				label.SetText(c.Address)
			case 4:
				label.SetText(c.BirthDate.String())
			}
		},
	)

}

func getListButtonCallback(app fyne.App, m *pharmacy.Model) func() {
	return func() {
		ctx, f := context.WithTimeout(context.Background(), time.Second)
		defer f()
		dao, err := m.GetDaoByType(daos.ConsumerE)
		if err != nil {
			fyne.LogError(err.Error(), err)
			return
		}

		cs := make(map[daos.ID]daos.Consumer)

		err = dao.List(ctx, &cs)
		if err != nil {
			fyne.LogError(err.Error(), err)
			return
		}

		ids := make([]daos.ID, 0, len(cs))
		for id := range cs {
			ids = append(ids, id)
		}
		table := getTableForConsumers(cs, ids)

		newW := app.NewWindow("list drugs")
		newW.Resize(fyne.NewSize(400, 400))
		newW.SetContent(table)
		newW.Show()
	}
}
