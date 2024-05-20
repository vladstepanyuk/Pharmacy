package orders_win

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

func addConsumerCallback(app fyne.App, m *pharmacy.Model) func() {
	return func() {
		cIdField := widget.NewEntry()
		cIdField.SetPlaceHolder("Consumer ID")
		dIdField := widget.NewEntry()
		dIdField.SetPlaceHolder("Drug ID")
		useField := widget.NewEntry()
		useField.SetPlaceHolder("Use ID")
		diseaseField := widget.NewEntry()
		diseaseField.SetPlaceHolder("Disease")
		countField := widget.NewEntry()
		countField.SetPlaceHolder("Count")

		w := app.NewWindow("")

		statusText := widget.NewLabel("")

		dao, err := m.GetDaoByType(daos.RecipeE)
		if err != nil {
			panic(err)
		}

		w.SetContent(container.NewGridWithRows(3,
			container.NewAdaptiveGrid(5, cIdField, dIdField, useField, diseaseField, countField),
			widget.NewButton("добавить", func() {
				cid, err := strconv.Atoi(cIdField.Text)
				if err != nil {
					statusText.SetText(err.Error())
					return
				}
				did, err := strconv.Atoi(cIdField.Text)
				if err != nil {
					statusText.SetText(err.Error())
					return
				}
				c, err := strconv.Atoi(countField.Text)
				if err != nil {
					statusText.SetText(err.Error())
					return
				}
				useId, err := pharmacy.GetUseId(m, useField.Text)
				if err != nil {
					statusText.SetText(err.Error())
					return
				}

				r := daos.Recipe{
					ConsumerId: cid, DoctorId: 1,
					DrugId:  did,
					UseId:   useId,
					Disease: diseaseField.Text,
					Count:   c,
				}

				ctx, f := context.WithTimeout(context.Background(), time.Second)
				defer f()

				id, err := dao.Insert(ctx, r)
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

func updateConsumerCallback(app fyne.App, m *pharmacy.Model) func() {
	return func() {
		idField := widget.NewEntry()
		idField.SetPlaceHolder("ID")
		cIdField := widget.NewEntry()
		cIdField.SetPlaceHolder("Consumer ID")
		dIdField := widget.NewEntry()
		dIdField.SetPlaceHolder("Drug ID")
		useField := widget.NewEntry()
		useField.SetPlaceHolder("Use ID")
		diseaseField := widget.NewEntry()
		diseaseField.SetPlaceHolder("Disease")
		countField := widget.NewEntry()
		countField.SetPlaceHolder("Count")

		w := app.NewWindow("")

		statusText := widget.NewLabel("")

		dao, err := m.GetDaoByType(daos.RecipeE)
		if err != nil {
			panic(err)
		}

		w.SetContent(container.NewGridWithRows(3,
			container.NewAdaptiveGrid(6, idField, cIdField, dIdField, useField, diseaseField, countField),
			widget.NewButton("добавить", func() {
				id, err := strconv.Atoi(idField.Text)
				if err != nil {
					statusText.SetText(err.Error())
					return
				}
				cid, err := strconv.Atoi(cIdField.Text)
				if err != nil {
					statusText.SetText(err.Error())
					return
				}
				did, err := strconv.Atoi(cIdField.Text)
				if err != nil {
					statusText.SetText(err.Error())
					return
				}
				c, err := strconv.Atoi(countField.Text)
				if err != nil {
					statusText.SetText(err.Error())
					return
				}

				useId, err := pharmacy.GetUseId(m, useField.Text)
				if err != nil {
					statusText.SetText(err.Error())
					return
				}

				r := daos.Recipe{
					ConsumerId: cid, DoctorId: 1,
					DrugId:  did,
					UseId:   useId,
					Disease: diseaseField.Text,
					Count:   c,
				}

				ctx, f := context.WithTimeout(context.Background(), time.Second)
				defer f()

				err = dao.Update(ctx, id, r)
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

func deleteConsumerCallback(app fyne.App, m *pharmacy.Model) func() {
	return func() {
		idField := widget.NewEntry()
		idField.SetPlaceHolder("ID")

		w := app.NewWindow("")

		statusText := widget.NewLabel("")

		dao, err := m.GetDaoByType(daos.RecipeE)
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
