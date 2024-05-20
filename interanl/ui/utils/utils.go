package utils

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"pharmacy/interanl/pharmacy"
	"pharmacy/interanl/pharmacy/daos"
	"strconv"
)

func GetTableForConsumers(cons map[daos.ID]daos.Consumer, ids []int) fyne.CanvasObject {
	return widget.NewTable(
		func() (int, int) {
			return len(cons), 5
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("wideeeeeeee text")
		},
		func(id widget.TableCellID, object fyne.CanvasObject) {
			label := object.(*widget.Label)
			c := cons[ids[id.Row]]
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

func GetTableForDrugs(m *pharmacy.Model, drugs map[daos.ID]daos.Drug, ids []int) fyne.CanvasObject {
	return widget.NewTable(
		func() (int, int) {
			return len(drugs), 5
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("wideee text")
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
			}
		},
	)

}

type Pair[T1, T2 any] struct {
	a T1
	b T2
}

func GetTableForTechs(drugs map[daos.ID]daos.TechnologyBook, ids []int) fyne.CanvasObject {
	return widget.NewTable(
		func() (int, int) {
			return len(drugs), 2
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("wideeeeee text")
		},
		func(id widget.TableCellID, object fyne.CanvasObject) {
			label := object.(*widget.Label)
			tech := drugs[ids[id.Row]]
			switch id.Col {
			case 0:
				label.SetText(strconv.Itoa(ids[id.Row]))
			case 1:
				label.SetText(tech.Description)
			}
		},
	)

}

func GetTableForOrders(orders map[daos.ID]daos.ConsumerOrder, ids []int) fyne.CanvasObject {
	return widget.NewTable(
		func() (int, int) {
			return len(orders), 5
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("wideeeeee text")
		},
		func(id widget.TableCellID, object fyne.CanvasObject) {
			label := object.(*widget.Label)
			order := orders[ids[id.Row]]
			switch id.Col {
			case 0:
				label.SetText(strconv.Itoa(ids[id.Row]))
			case 1:
				label.SetText(order.OrderTime)
			case 2:
				label.SetText(order.CompleteTime)
			case 3:
				label.SetText(order.Status)
			case 4:
				label.SetText(strconv.Itoa(order.ResipeId))
			}
		},
	)

}
