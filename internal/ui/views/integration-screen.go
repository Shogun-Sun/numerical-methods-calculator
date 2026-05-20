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

	// 1. Создаем массивы для полей ввода (чтобы потом легко прочитать в цикле)
	var xEntries [5]*widget.Entry
	var yEntries [5]*widget.Entry

	// Заголовки для таблицы ввода
	inputGrid := container.NewGridWithColumns(3,
		widget.NewLabelWithStyle("Узел (i)", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Значение X_i", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Значение Y_i", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
	)

	// Дефолтные данные из твоей задачи (чтобы не вводить руками каждый раз при тестах)
	defaultX := []string{"0.725", "0.727", "0.729", "0.731", "0.733"}
	defaultY := []string{"0.66314", "0.66463", "0.66612", "0.66761", "0.66769"}

	// Заполняем сетку вводами в цикле
	for i := 0; i < 5; i++ {
		xEntries[i] = widget.NewEntry()
		xEntries[i].SetText(defaultX[i]) // Устанавливаем дефолтное значение

		yEntries[i] = widget.NewEntry()
		yEntries[i].SetText(defaultY[i]) // Устанавливаем дефолтное значение

		// Добавляем строку в сетку: Номер узла, поле X, поле Y
		inputGrid.Add(widget.NewLabel(fmt.Sprintf("Узел %d", i)))
		inputGrid.Add(xEntries[i])
		inputGrid.Add(yEntries[i])
	}

	// 2. Создаем блок для вывода результатов (изначально пустой)
	resultCard := container.NewVBox()
	errorLabel := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	errorLabel.Hide() // Показываем только при ошибках парсинга

	// Названия методов для красивого отображения в таблице результатов
	methodNames := map[string]string{
		"rect_left":    "Метод левых прямоугольников",
		"rect_right":   "Метод правых прямоугольников",
		"trapezoidal":  "Метод трапеций",
		"simpson":      "Метод Симпсона",
		"newton_cotes": "Формула Ньютона-Котеса (Бооля)",
	}

	// 3. Логика кнопки «Рассчитать»
	calcBtn := widget.NewButton("Рассчитать", func() {
		errorLabel.Hide()
		resultCard.Objects = nil // Очищаем старые результаты перед новым расчетом

		X := make([]float64, 5)
		Y := make([]float64, 5)

		// Парсим строки из полей ввода в float64
		for i := 0; i < 5; i++ {
			var errX, errY error
			X[i], errX = strconv.ParseFloat(xEntries[i].Text, 64)
			Y[i], errY = strconv.ParseFloat(yEntries[i].Text, 64)

			if errX != nil || errY != nil {
				errorLabel.SetText("Ошибка: Введены некорректные числа!")
				errorLabel.Show()
				return
			}
		}

		// Вызываем твою функцию вычислений
		results, err := methods.NumericalIntegration(X, Y)
		if err != nil {
			errorLabel.SetText(fmt.Sprintf("Ошибка расчета: %v", err))
			errorLabel.Show()
			return
		}

		// Строим красивую псевдотаблицу результатов
		resultGrid := container.NewGridWithColumns(2,
			widget.NewLabelWithStyle("Метод интегрирования", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			widget.NewLabelWithStyle("Результат", fyne.TextAlignTrailing, fyne.TextStyle{Bold: true}),
		)

		// Выводим данные. Спецификатор %.7f принудительно сохранит ровную колонку знаков
		for key, val := range results {
			name := methodNames[key]
			resultGrid.Add(widget.NewLabel(name))
			resultGrid.Add(widget.NewLabelWithStyle(fmt.Sprintf("%.7f", val), fyne.TextAlignTrailing, fyne.TextStyle{Monospace: true}))
		}

		// Добавляем сетку результатов в контейнер и обновляем экран
		resultCard.Add(widget.NewSeparator())
		resultCard.Add(widget.NewLabelWithStyle("Результаты вычислений:", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}))
		resultCard.Add(resultGrid)
		resultCard.Refresh()
	})

	// Кнопка возврата
	backBtn := widget.NewButton("Назад в меню", func() {
		ShowMainWindow(w)
	})

	// Собираем весь экран воедино в один вертикальный контейнер со скроллом
	mainLayout := container.NewVScroll(container.NewVBox(
		widget.NewLabelWithStyle("Входные данные таблицы", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		inputGrid,
		calcBtn,
		errorLabel,
		resultCard,
		widget.NewSeparator(),
		backBtn,
	))

	w.SetContent(mainLayout)
}
