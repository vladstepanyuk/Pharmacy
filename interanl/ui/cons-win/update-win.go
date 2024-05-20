package cons_win

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
		nameField := widget.NewEntry()
		nameField.SetPlaceHolder("Name")
		phoneField := widget.NewEntry()
		phoneField.SetPlaceHolder("Phone")
		addressField := widget.NewEntry()
		addressField.SetPlaceHolder("Address")
		birthdayField := widget.NewEntry()
		birthdayField.SetPlaceHolder("Birthday: YYYY-MM-DD")

		w := app.NewWindow("")

		statusText := widget.NewLabel("")

		dao, err := m.GetDaoByType(daos.ConsumerE)
		if err != nil {
			panic(err)
		}

		w.SetContent(container.NewGridWithRows(3,
			container.NewAdaptiveGrid(4, nameField, phoneField, addressField, birthdayField),
			widget.NewButton("добавить", func() {
				if _, err := strconv.ParseUint(phoneField.Text, 10, 64); len(phoneField.Text) != 11 || err != nil {
					statusText.SetText("невалидный формат телефона")
					return
				}
				t, err := time.Parse("2006-01-02", birthdayField.Text)
				if err != nil {
					statusText.SetText(err.Error())
					return
				}

				c := daos.Consumer{Name: nameField.Text, Phone: phoneField.Text, Address: addressField.Text, BirthDate: t}

				ctx, f := context.WithTimeout(context.Background(), time.Second)
				defer f()

				id, err := dao.Insert(ctx, c)
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
		nameField := widget.NewEntry()
		nameField.SetPlaceHolder("Name")
		phoneField := widget.NewEntry()
		phoneField.SetPlaceHolder("Phone")
		addressField := widget.NewEntry()
		addressField.SetPlaceHolder("Address")
		birthdayField := widget.NewEntry()
		birthdayField.SetPlaceHolder("Birthday: YYYY-MM-DD")

		w := app.NewWindow("")

		statusText := widget.NewLabel("")

		dao, err := m.GetDaoByType(daos.ConsumerE)
		if err != nil {
			panic(err)
		}

		w.SetContent(container.NewGridWithRows(3,
			container.NewAdaptiveGrid(5, idField, nameField, phoneField, addressField, birthdayField),
			widget.NewButton("добавить", func() {
				id, err := strconv.Atoi(idField.Text)
				if err != nil {
					statusText.SetText(err.Error())
					return
				}

				if _, err := strconv.ParseUint(phoneField.Text, 10, 64); len(phoneField.Text) != 11 || err != nil {
					statusText.SetText("невалидный формат телефона")
					return
				}
				t, err := time.Parse("2006-01-02", birthdayField.Text)
				if err != nil {
					statusText.SetText(err.Error())
					return
				}

				c := daos.Consumer{Name: nameField.Text, Phone: phoneField.Text, Address: addressField.Text, BirthDate: t}

				ctx, f := context.WithTimeout(context.Background(), time.Second)
				defer f()

				err = dao.Update(ctx, id, c)
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

		dao, err := m.GetDaoByType(daos.ConsumerE)
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
