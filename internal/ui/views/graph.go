package views

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"numerical-methods-calculator/cmd/methods"
)

func OpenGraphWindow(parentApp fyne.App, equationText string, errorLabel *widget.Label) {
	errorLabel.Hide()

	f, err := methods.MakeFunctionFromString(equationText)
	if err != nil {
		errorLabel.SetText(fmt.Sprintf("Ошибка парсинга для графика: %v", err))
		errorLabel.Show()
		return
	}

	graphWin := parentApp.NewWindow("Визуализация корней: " + equationText)

	chartWidget := methods.NewInteractiveChart(f, 600, 400)

	closeBtn := widget.NewButton("Закрыть график", func() {
		graphWin.Close()
	})

	windowContent := container.NewBorder(
		closeBtn,
		nil,
		nil,
		chartWidget,
	)

	graphWin.SetContent(windowContent)
	graphWin.Resize(fyne.NewSize(650, 480))
	graphWin.Show()
}
