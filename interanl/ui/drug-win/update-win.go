package drug_win

import (
	"context"
	"errors"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"pharmacy/interanl/pharmacy"
	"pharmacy/interanl/pharmacy/daos"
	"strconv"
	"time"
)

func validateEntry(name, typeName, pricetxt, crittxt string, prod bool, m *pharmacy.Model) (daos.Drug, error) {
	tId, err := pharmacy.GetTypeByNameAndProduced(m, typeName, prod)
	if err != nil {
		return daos.Drug{}, errors.New("error: неправильный тип лекарства")
	}

	price, err := strconv.ParseUint(pricetxt, 10, 64)
	if err != nil {
		return daos.Drug{}, errors.New("error: цена не является числом")
	}

	crit, err := strconv.ParseUint(crittxt, 10, 64)
	if err != nil {
		return daos.Drug{}, errors.New("error: крит. число не является числом")
	}

	drug := daos.Drug{Name: name, TypeId: tId, CriticalCount: int(crit), Price: int(price)}
	return drug, nil
}

func getInsertWindow(dao *daos.GeneralDao, m *pharmacy.Model) fyne.CanvasObject {
	nameField := widget.NewEntry()
	nameField.SetPlaceHolder("введите имя")
	typeNameField := widget.NewEntry()
	typeNameField.SetPlaceHolder("введите тип")
	priceField := widget.NewEntry()
	priceField.SetPlaceHolder("введите цену")
	critField := widget.NewEntry()
	critField.SetPlaceHolder("введите крит. количество")
	isProduced := widget.NewCheck("производится?", func(checked bool) {})

	hc := container.NewAdaptiveGrid(5, nameField, typeNameField, priceField, critField, isProduced)
	hc.Resize(fyne.NewSize(400, 100))

	statusText := widget.NewLabel("")

	return container.NewGridWithRows(
		3,
		hc,
		widget.NewButton("отправить", func() {

			drug, err := validateEntry(nameField.Text, typeNameField.Text, priceField.Text, critField.Text, isProduced.Checked, m)
			if err != nil {
				statusText.SetText(err.Error())
				return
			}
			fmt.Println(drug)

			ctx, f := context.WithTimeout(context.Background(), time.Second)
			defer f()
			_, err = dao.Insert(ctx, drug)
			if err != nil {
				fyne.LogError("произошла ошибка", err)
				statusText.SetText("error: " + err.Error())
				return
			}
			statusText.SetText("ok")

		}),
		statusText,
	)

}

func getInsertButtonCallback(app fyne.App, m *pharmacy.Model) func() {
	return func() {

		newW := app.NewWindow("insert drugs")
		newW.Resize(fyne.NewSize(400, 200))

		dao, err := m.GetDaoByType(daos.DrugE)
		if err != nil {
			panic(err)
		}

		iWin := getInsertWindow(dao, m)

		newW.SetContent(iWin)

		newW.Show()
	}
}

func getUpdateWindow(dao *daos.GeneralDao, m *pharmacy.Model) fyne.CanvasObject {

	idField := widget.NewEntry()
	idField.SetPlaceHolder("введите айди")
	nameField := widget.NewEntry()
	nameField.SetPlaceHolder("введите имя")
	typeNameField := widget.NewEntry()
	typeNameField.SetPlaceHolder("введите тип")
	priceField := widget.NewEntry()
	priceField.SetPlaceHolder("введите цену")
	critField := widget.NewEntry()
	critField.SetPlaceHolder("введите крит. количество")
	isProduced := widget.NewCheck("производится?", func(checked bool) {})

	hc := container.NewAdaptiveGrid(6, idField, nameField, typeNameField, priceField, critField, isProduced)
	hc.Resize(fyne.NewSize(400, 100))

	statusText := widget.NewLabel("")

	return container.NewGridWithRows(
		3,
		hc,
		widget.NewButton("обновить", func() {

			id, err := strconv.ParseUint(idField.Text, 10, 64)
			if err != nil {
				statusText.SetText("error: id не число")
			}
			drug, err := validateEntry(nameField.Text, typeNameField.Text, priceField.Text, critField.Text, isProduced.Checked, m)
			if err != nil {
				statusText.SetText(err.Error())
			}

			ctx, f := context.WithTimeout(context.Background(), time.Second)
			defer f()

			err = dao.Update(ctx, daos.ID(id), drug)
			if err != nil {
				fyne.LogError("произошла ошибка", err)
				statusText.SetText("error: " + err.Error())
				return
			}
			statusText.SetText("ok")

		}),
		statusText,
	)

}

func getUpdateButtonCallback(app fyne.App, m *pharmacy.Model) func() {

	return func() {

		newW := app.NewWindow("update drugs")
		newW.Resize(fyne.NewSize(400, 200))

		dao, err := m.GetDaoByType(daos.DrugE)
		if err != nil {
			panic(err)
		}

		iWin := getUpdateWindow(dao, m)

		newW.SetContent(iWin)

		newW.Show()
	}

}

func getDeleteWindow(dao *daos.GeneralDao, m *pharmacy.Model) fyne.CanvasObject {

	idField := widget.NewEntry()
	idField.SetText("введите айди")

	statusText := widget.NewLabel("")

	return container.NewGridWithRows(
		3,
		idField,
		widget.NewButton("удалить", func() {

			id, err := strconv.ParseUint(idField.Text, 10, 64)
			if err != nil {
				statusText.SetText("error: id не число")
				return
			}

			ctx, f := context.WithTimeout(context.Background(), time.Second)
			defer f()

			err = dao.Delete(ctx, daos.ID(id))

			if err != nil {
				fyne.LogError("произошла ошибка", err)
				statusText.SetText("error: " + err.Error())
				return
			}
			statusText.SetText("ok")

		}),
		statusText,
	)

}

func getDeleteButtonCallback(app fyne.App, m *pharmacy.Model) func() {

	return func() {

		newW := app.NewWindow("update drugs")
		newW.Resize(fyne.NewSize(400, 200))

		dao, err := m.GetDaoByType(daos.DrugE)
		if err != nil {
			panic(err)
		}

		iWin := getDeleteWindow(dao, m)

		newW.SetContent(iWin)

		newW.Show()
	}

}
