package co_win

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

func addOrder(app fyne.App, m *pharmacy.Model) func() {
	return func() {
		rIdField := widget.NewEntry()
		rIdField.SetPlaceHolder("Recipe ID")
		statussel := widget.NewSelect([]string{
			"Готов", "Ожидает компонент",
			"Ожидает получения", "В процессе",
		}, func(s string) {})
		completeTime := widget.NewEntry()
		completeTime.SetPlaceHolder("Complete time")
		orderTime := widget.NewEntry()
		orderTime.SetPlaceHolder("Order time")

		w := app.NewWindow("")

		statusText := widget.NewLabel("")

		dao, err := m.GetDaoByType(daos.ConsumerOrderE)
		if err != nil {
			panic(err)
		}

		w.SetContent(container.NewGridWithRows(3,
			container.NewAdaptiveGrid(4, rIdField, statussel, completeTime, orderTime),
			widget.NewButton("добавить", func() {
				co := daos.ConsumerOrder{}

				var oTyp int
				switch statussel.Selected {
				case "Готов":
					oTyp = daos.SReady
				case "Ожидает компонент":
					oTyp = daos.SWaitComps
				case "Ожидает получения":
					oTyp = daos.SWait
				case "В процессе":
					oTyp = daos.SInProgress
				default:
					panic("aflwasef")
				}
				co.Status = pharmacy.GetStatusString(oTyp)

				if completeTime.Text != "" {
					_, err = time.Parse("2006-01-02", completeTime.Text)
					if err != nil {
						statusText.SetText(err.Error())
						return
					}
					co.CompleteTime = completeTime.Text
				} else {
					if oTyp == daos.SReady || oTyp == daos.SWait {
						statusText.SetText("error: complete time must have smth")
						return
					}
					co.CompleteTime = time.Now().Format("2006-01-02")
				}

				_, err := time.Parse("2006-01-02", orderTime.Text)
				if err != nil {
					statusText.SetText(err.Error())
					return
				}
				co.OrderTime = orderTime.Text

				rId, err := strconv.Atoi(rIdField.Text)
				if err != nil {
					statusText.SetText(err.Error())
					return
				}
				co.ResipeId = rId

				ctx, f := context.WithTimeout(context.Background(), time.Second)
				defer f()

				id, err := dao.Insert(ctx, co)
				if err != nil {
					statusText.SetText(err.Error())
					return
				}

				statusText.SetText("ok, new id:" + strconv.Itoa(id))

			}),
			statusText,
		))

		w.Show()

	}
}

func updateOrder(app fyne.App, m *pharmacy.Model) func() {
	return func() {
		idField := widget.NewEntry()
		idField.SetPlaceHolder("ID")
		rIdField := widget.NewEntry()
		rIdField.SetPlaceHolder("Recipe ID")
		statussel := widget.NewSelect([]string{
			"Готов", "Ожидает компонент",
			"Ожидает получения", "В процессе",
		}, func(s string) {})
		completeTime := widget.NewEntry()
		completeTime.SetPlaceHolder("Complete time")
		orderTime := widget.NewEntry()
		orderTime.SetPlaceHolder("Order time")

		w := app.NewWindow("")

		statusText := widget.NewLabel("")

		dao, err := m.GetDaoByType(daos.ConsumerOrderE)
		if err != nil {
			panic(err)
		}

		w.SetContent(container.NewGridWithRows(3,
			container.NewAdaptiveGrid(4, idField, rIdField, statussel, completeTime, orderTime),
			widget.NewButton("добавить", func() {
				id, err := strconv.Atoi(idField.Text)
				if err != nil {
					statusText.SetText(err.Error())
					return
				}
				co := daos.ConsumerOrder{}

				var oTyp int
				switch statussel.Selected {
				case "Готов":
					oTyp = daos.SReady
				case "Ожидает компонент":
					oTyp = daos.SWaitComps
				case "Ожидает получения":
					oTyp = daos.SWait
				case "В процессе":
					oTyp = daos.SInProgress
				default:
					panic("aflwasef")
				}
				co.Status = pharmacy.GetStatusString(oTyp)

				if completeTime.Text != "" {
					_, err = time.Parse("2006-01-02", completeTime.Text)
					if err != nil {
						statusText.SetText(err.Error())
						return
					}
					co.CompleteTime = completeTime.Text
				} else {
					if oTyp == daos.SReady || oTyp == daos.SWait {
						statusText.SetText("error: complete time must have smth")
						return
					}
					co.CompleteTime = time.Now().Format("2006-01-02")
				}

				_, err = time.Parse("2006-01-02", orderTime.Text)
				if err != nil {
					statusText.SetText(err.Error())
					return
				}
				co.OrderTime = orderTime.Text

				rId, err := strconv.Atoi(rIdField.Text)
				if err != nil {
					statusText.SetText(err.Error())
					return
				}
				co.ResipeId = rId

				ctx, f := context.WithTimeout(context.Background(), time.Second)
				defer f()

				err = dao.Update(ctx, id, co)
				if err != nil {
					statusText.SetText(err.Error())
					return
				}

				statusText.SetText("ok")

			}),
			statusText,
		))

		w.Show()

	}
}

func deleteOrder(app fyne.App, m *pharmacy.Model) func() {
	return func() {
		idField := widget.NewEntry()
		idField.SetPlaceHolder("ID")

		w := app.NewWindow("")

		statusText := widget.NewLabel("")

		dao, err := m.GetDaoByType(daos.ConsumerOrderE)
		if err != nil {
			panic(err)
		}

		w.SetContent(container.NewGridWithRows(3,
			container.NewAdaptiveGrid(1, idField),
			widget.NewButton("добавить", func() {

				id, err := strconv.Atoi(idField.Text)
				if err != nil {
					statusText.SetText(err.Error())
					return
				}

				ctx, f := context.WithTimeout(context.Background(), time.Second)
				defer f()

				err = dao.Delete(ctx, id)
				if err != nil {
					statusText.SetText(err.Error())
					return
				}

				statusText.SetText("ok")

			}),
			statusText,
		))

		w.Show()

	}
}
