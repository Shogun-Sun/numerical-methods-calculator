package views

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func ShowMainWindow(w fyne.Window) {
	mainWindow(w)
}

func mainWindow(w fyne.Window) {
	w.SetTitle("Численные методы - Главное меню")

	title := widget.NewLabelWithStyle("Калькулятор численных методов", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	numericalIntegrationBtn := widget.NewButton("Численное интегрирование", func() {
		showIntegrationScreen(w)
	})

	menuContainer := container.NewVBox(
		title,
		widget.NewSeparator(),
		numericalIntegrationBtn,
	)

	w.SetContent(menuContainer)
}
