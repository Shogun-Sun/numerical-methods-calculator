package views

import (
	"fmt"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"numerical-methods-calculator/cmd/methods"
)

func showEquationsScreen(w fyne.Window) {
	w.SetTitle("Численное решение нелинейных уравнений")

	parentApp := fyne.CurrentApp()

	equationEntry := widget.NewEntry()
	equationEntry.SetText("ln(x) + x - 2")

	x0Entry := widget.NewEntry()
	x0Entry.SetText("1.0")

	x1Entry := widget.NewEntry()
	x1Entry.SetText("2.0")

	epsEntry := widget.NewEntry()
	epsEntry.SetText("0.001")

	maxIterationsEntry := widget.NewEntry()
	maxIterationsEntry.SetText("10000")

	methodSelect := widget.NewSelect([]string{"Метод хорд"}, func(value string) {})
	methodSelect.SetSelected("Метод хорд")

	formGrid := container.NewGridWithColumns(2,
		widget.NewLabel("Введите уравнение f(x) = 0:"), equationEntry,
		widget.NewLabel("Неподвижная точка (x0):"), x0Entry,
		widget.NewLabel("Начальное приближение (x1):"), x1Entry,
		widget.NewLabel("Точность (eps):"), epsEntry,
		widget.NewLabel("Максимальное количество итераций:"), maxIterationsEntry,
		widget.NewLabel("Выберите метод решения:"), methodSelect,
	)

	errorLabel := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	errorLabel.Hide()

	resultCard := container.NewVBox()

	helpBtn := widget.NewButton("Как вводить уравнения? (Инструкция)", func() {
		openHelpWindow()
	})

	plotBtn := widget.NewButton("Показать график уравнения", func() {
		OpenGraphWindow(parentApp, equationEntry.Text, errorLabel)
	})

	calcBtn := widget.NewButton("Рассчитать корень", func() {
		errorLabel.Hide()
		resultCard.Objects = nil

		x0, errX0 := strconv.ParseFloat(x0Entry.Text, 64)
		x1, errX1 := strconv.ParseFloat(x1Entry.Text, 64)
		eps, errEps := strconv.ParseFloat(epsEntry.Text, 64)

		if errX0 != nil || errX1 != nil || errEps != nil {
			errorLabel.SetText("Ошибка: Проверьте корректность ввода чисел (x0, x1, eps)!")
			errorLabel.Show()
			return
		}

		f, err := methods.MakeFunctionFromString(equationEntry.Text)
		if err != nil {
			errorLabel.SetText(fmt.Sprintf("Ошибка парсинга уравнения: %v", err))
			errorLabel.Show()
			return
		}

		maxIterations, errMaxIter := strconv.Atoi(maxIterationsEntry.Text)
		if errMaxIter != nil || maxIterations <= 0 {
			errorLabel.SetText("Ошибка: Максимальное число итераций должно быть целым положительным числом!")
			errorLabel.Show()
			return
		}

		var root float64
		var steps []methods.IterationStep
		var calcErr error

		switch methodSelect.Selected {
		case "Метод хорд":
			root, steps, calcErr = methods.ChordMethod(f, x0, x1, eps, maxIterations)
		default:
			calcErr = fmt.Errorf("выбран неизвестный метод")
		}

		if calcErr != nil {
			errorLabel.SetText(fmt.Sprintf("Ошибка расчета: %v", calcErr))
			errorLabel.Show()
			return
		}

		finalGrid := container.NewGridWithColumns(2,
			widget.NewLabel("Найденный корень уравнения:"),
			widget.NewLabelWithStyle(fmt.Sprintf("%.6f", root), fyne.TextAlignTrailing, fyne.TextStyle{Monospace: true, Bold: true}),
			widget.NewLabel("Всего шагов (итераций):"),
			widget.NewLabelWithStyle(strconv.Itoa(len(steps)), fyne.TextAlignTrailing, fyne.TextStyle{Monospace: true}),
		)

		tableGrid := container.NewGridWithColumns(4,
			widget.NewLabelWithStyle("Итер. (i)", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
			widget.NewLabelWithStyle("x_n", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
			widget.NewLabelWithStyle("f(x_n)", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
			widget.NewLabelWithStyle("x_{n+1}", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		)

		for _, step := range steps {
			tableGrid.Add(widget.NewLabelWithStyle(strconv.Itoa(step.Num), fyne.TextAlignCenter, fyne.TextStyle{Monospace: true}))
			tableGrid.Add(widget.NewLabelWithStyle(fmt.Sprintf("%.6f", step.Xn), fyne.TextAlignCenter, fyne.TextStyle{Monospace: true}))
			tableGrid.Add(widget.NewLabelWithStyle(fmt.Sprintf("%.6f", step.Fxn), fyne.TextAlignCenter, fyne.TextStyle{Monospace: true}))
			tableGrid.Add(widget.NewLabelWithStyle(fmt.Sprintf("%.6f", step.NextX), fyne.TextAlignCenter, fyne.TextStyle{Monospace: true}))
		}

		scrollTable := container.NewScroll(tableGrid)
		scrollTable.SetMinSize(fyne.NewSize(0, 200))

		resultCard.Add(widget.NewSeparator())
		resultCard.Add(widget.NewLabelWithStyle("Результат вычислений:", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}))
		resultCard.Add(finalGrid)
		resultCard.Add(widget.NewSeparator())
		resultCard.Add(widget.NewLabelWithStyle("Таблица промежуточных шагов:", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}))
		resultCard.Add(scrollTable)
		resultCard.Refresh()
	})

	backBtn := widget.NewButton("Назад в меню", func() {
		ShowMainWindow(w)
	})

	buttonsGrid := container.NewGridWithColumns(3, helpBtn, plotBtn, calcBtn)

	mainLayout := container.NewVBox(
		widget.NewLabelWithStyle("Параметры уравнения", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		formGrid,
		layout.NewSpacer(),
		buttonsGrid,
		errorLabel,
		resultCard,
		layout.NewSpacer(),
		backBtn,
	)

	w.SetContent(container.NewVScroll(mainLayout))
}
