package views

import (
	"fmt"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"numerical-methods-calculator/cmd/methods"
)

func showIntegrationScreen(w fyne.Window) {
	w.SetTitle("Численное интегрирование")

	var xEntries []*widget.Entry
	var yEntries []*widget.Entry

	rowsContainer := container.NewVBox()

	appendRow := func(defaultX, defaultY string) {
		index := len(xEntries)

		xEntry := widget.NewEntry()
		xEntry.SetText(defaultX)
		xEntries = append(xEntries, xEntry)

		yEntry := widget.NewEntry()
		yEntry.SetText(defaultY)
		yEntries = append(yEntries, yEntry)

		rowGrid := container.NewGridWithColumns(3,
			widget.NewLabel(fmt.Sprintf("Узел %d", index)),
			xEntry,
			yEntry,
		)
		rowsContainer.Add(rowGrid)
	}

	initX := []string{"0.725", "0.727", "0.729", "0.731", "0.733"}
	initY := []string{"0.66314", "0.66463", "0.66612", "0.66761", "0.66769"}
	for i := 0; i < 5; i++ {
		appendRow(initX[i], initY[i])
	}

	headerGrid := container.NewGridWithColumns(3,
		widget.NewLabelWithStyle("Узел (i)", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Значение X_i", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Значение Y_i", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
	)

	var counterLabel *widget.Label

	addBtn := widget.NewButton("+ Добавить узел", func() {
		if len(xEntries) >= 6 {
			return
		}
		nextXStr := ""
		if len(xEntries) >= 2 {
			xPrev, _ := strconv.ParseFloat(xEntries[len(xEntries)-1].Text, 64)
			xPrev2, _ := strconv.ParseFloat(xEntries[len(xEntries)-2].Text, 64)
			step := xPrev - xPrev2
			nextXStr = fmt.Sprintf("%.3f", xPrev+step)
		}
		appendRow(nextXStr, "")
		counterLabel.SetText(fmt.Sprintf("Всего узлов: %d", len(xEntries)))
		rowsContainer.Refresh()
	})

	removeBtn := widget.NewButton("- Удалить узел", func() {
		if len(xEntries) <= 2 {
			return
		}
		lastIdx := len(rowsContainer.Objects) - 1
		rowsContainer.Objects = rowsContainer.Objects[:lastIdx]

		xEntries = xEntries[:len(xEntries)-1]
		yEntries = yEntries[:len(yEntries)-1]

		counterLabel.SetText(fmt.Sprintf("Всего узлов: %d", len(xEntries)))
		rowsContainer.Refresh()
	})

	counterLabel = widget.NewLabel(fmt.Sprintf("Всего узлов: %d", len(xEntries)))

	controlsGrid := container.NewGridWithColumns(3, removeBtn, counterLabel, addBtn)

	resultCard := container.NewVBox()
	errorLabel := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	errorLabel.Hide()

	methodNames := map[string]string{
		"rect_left":    "Метод левых прямоугольников",
		"rect_right":   "Метод правых прямоугольников",
		"trapezoidal":  "Метод трапеций",
		"simpson":      "Метод Симпсона",
		"newton_cotes": "Формула Ньютона-Котеса (высший доступный порядок)",
	}

	calcBtn := widget.NewButton("Рассчитать интегралы", func() {
		errorLabel.Hide()
		resultCard.Objects = nil

		totalNodes := len(xEntries)
		X := make([]float64, totalNodes)
		Y := make([]float64, totalNodes)

		for i := 0; i < totalNodes; i++ {
			var errX, errY error
			X[i], errX = strconv.ParseFloat(xEntries[i].Text, 64)
			Y[i], errY = strconv.ParseFloat(yEntries[i].Text, 64)

			if errX != nil || errY != nil {
				errorLabel.SetText(fmt.Sprintf("Ошибка: Некорректные числа в узле %d!", i))
				errorLabel.Show()
				return
			}
		}

		results, err := methods.NumericalIntegration(X, Y)
		if err != nil {
			errorLabel.SetText(fmt.Sprintf("Ошибка расчета: %v", err))
			errorLabel.Show()
			return
		}

		resultGrid := container.NewGridWithColumns(2,
			widget.NewLabelWithStyle("Метод интегрирования", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			widget.NewLabelWithStyle("Результат", fyne.TextAlignTrailing, fyne.TextStyle{Bold: true}),
		)

		for key, val := range results {
			name := methodNames[key]
			resultGrid.Add(widget.NewLabel(name))
			resultGrid.Add(widget.NewLabelWithStyle(fmt.Sprintf("%.7f", val), fyne.TextAlignTrailing, fyne.TextStyle{Monospace: true}))
		}

		if _, exists := results["simpson"]; !exists {
			resultCard.Add(widget.NewLabelWithStyle("* Метод Симпсона скрыт, так как он требует нечетное число узлов.", fyne.TextAlignLeading, fyne.TextStyle{Italic: true}))
		}

		resultCard.Add(widget.NewSeparator())
		resultCard.Add(widget.NewLabelWithStyle("Результаты вычислений:", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}))
		resultCard.Add(resultGrid)
		resultCard.Refresh()
	})

	backBtn := widget.NewButton("Назад в меню", func() {
		ShowMainWindow(w)
	})

	mainLayout := container.NewVScroll(container.NewVBox(
		widget.NewLabelWithStyle("Входные данные таблицы", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		headerGrid,
		rowsContainer,
		controlsGrid,
		widget.NewSeparator(),
		calcBtn,
		errorLabel,
		resultCard,
		widget.NewSeparator(),
		backBtn,
	))

	w.SetContent(mainLayout)
}
