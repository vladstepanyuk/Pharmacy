package tb_win

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

func addTechButtonCallback(app fyne.App, m *pharmacy.Model) func() {
	return func() {
		w := app.NewWindow("добавить технологию")
		w.Resize(fyne.NewSize(400, 200))

		descrField := widget.NewEntry()
		descrField.SetPlaceHolder("описание")

		statusText := widget.NewLabel("")

		dao, err := m.GetDaoByType(daos.TechnologyBookE)
		if err != nil {
			panic(err)
		}

		w.SetContent(container.NewGridWithRows(3,
			descrField,
			widget.NewButton("добавить", func() {

				dtxt := descrField.Text

				ctx, f := context.WithTimeout(context.Background(), time.Second)
				defer f()

				tb := daos.TechnologyBook{
					Description: dtxt,
				}
				id, err := dao.Insert(ctx, tb)
				if err != nil {
					statusText.SetText(err.Error())
					return
				}
				statusText.SetText("ok, id нового элемента: " + strconv.Itoa(id))
			}),
			statusText,
		))

		w.Show()
	}
}

func updateTechButtonCallback(app fyne.App, m *pharmacy.Model) func() {
	return func() {
		w := app.NewWindow("обновить технологию")
		w.Resize(fyne.NewSize(400, 200))

		idField := widget.NewEntry()
		idField.SetPlaceHolder("id")
		descrField := widget.NewEntry()
		descrField.SetPlaceHolder("описание")

		statusText := widget.NewLabel("")

		dao, err := m.GetDaoByType(daos.TechnologyBookE)
		if err != nil {
			panic(err)
		}

		w.SetContent(container.NewGridWithRows(3,
			container.NewAdaptiveGrid(2,
				idField,
				descrField,
			),
			widget.NewButton("добавить", func() {

				id, err := strconv.ParseUint(idField.Text, 10, 64)
				if err != nil {
					statusText.SetText("ОШИБКА: id не число")
					return
				}

				dtxt := descrField.Text

				ctx, f := context.WithTimeout(context.Background(), time.Second)
				defer f()

				tb := daos.TechnologyBook{
					Description: dtxt,
				}

				err = dao.Update(ctx, daos.ID(id), tb)
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

func deleteTechButtonCallback(app fyne.App, m *pharmacy.Model) func() {
	return func() {
		w := app.NewWindow("обновить технологию")
		w.Resize(fyne.NewSize(400, 200))

		idField := widget.NewEntry()
		idField.SetPlaceHolder("id")

		statusText := widget.NewLabel("")

		dao, err := m.GetDaoByType(daos.TechnologyBookE)
		if err != nil {
			panic(err)
		}

		w.SetContent(container.NewGridWithRows(3,
			container.NewAdaptiveGrid(2,
				idField,
			),
			widget.NewButton("добавить", func() {

				id, err := strconv.ParseUint(idField.Text, 10, 64)
				if err != nil {
					statusText.SetText("ОШИБКА: id не число")
					return
				}

				ctx, f := context.WithTimeout(context.Background(), time.Second)
				defer f()

				err = dao.Delete(ctx, daos.ID(id))
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

func setTechForDrugButtonCallback(app fyne.App, m *pharmacy.Model) func() {
	return func() {
		w := app.NewWindow("")
		w.Resize(fyne.NewSize(400, 200))

		idField := widget.NewEntry()
		idField.SetPlaceHolder("id лекарства")
		techField := widget.NewEntry()
		techField.SetPlaceHolder("id технологии")

		statusText := widget.NewLabel("")

		w.SetContent(container.NewGridWithRows(3,
			container.NewAdaptiveGrid(2,
				idField,
				techField,
			),
			widget.NewButton("добавить", func() {

				idd, err := strconv.ParseUint(idField.Text, 10, 64)
				if err != nil {
					statusText.SetText("ОШИБКА: id лекарства не число")
					return
				}
				idt, err := strconv.ParseUint(techField.Text, 10, 64)
				if err != nil {
					statusText.SetText("ОШИБКА: id технологии не число")
					return
				}

				ctx, f := context.WithTimeout(context.Background(), time.Second)
				defer f()

				err = m.AddDrugIdForTech(ctx, daos.ID(idt), daos.ID(idd))
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

func setCompForTechButtonCallback(app fyne.App, m *pharmacy.Model) func() {
	return func() {
		w := app.NewWindow("")
		w.Resize(fyne.NewSize(400, 200))

		idField := widget.NewEntry()
		idField.SetPlaceHolder("id компонента")
		techField := widget.NewEntry()
		techField.SetPlaceHolder("id технологии")
		countField := widget.NewEntry()
		countField.SetPlaceHolder("кол-во")

		statusText := widget.NewLabel("")

		w.SetContent(container.NewGridWithRows(3,
			container.NewAdaptiveGrid(3,
				idField,
				techField,
				countField,
			),
			widget.NewButton("добавить", func() {

				idd, err := strconv.ParseUint(idField.Text, 10, 64)
				if err != nil {
					statusText.SetText("ОШИБКА: id лекарства не число")
					return
				}
				idt, err := strconv.ParseUint(techField.Text, 10, 64)
				if err != nil {
					statusText.SetText("ОШИБКА: id технологии не число")
					return
				}
				count, err := strconv.Atoi(countField.Text)
				if err != nil {
					statusText.SetText("ОШИБКА: кол-во не число")
					return
				}

				ctx, f := context.WithTimeout(context.Background(), time.Second)
				defer f()

				err = m.DiffCompIdForTech(ctx, daos.ID(idt), daos.ID(idd), count)

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
