package main

import (
	"numerical-methods-calculator/internal/ui/views"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func main() {

	app := app.New()
	window := app.NewWindow("Численные методы - Калькулятор")

	window.Resize(fyne.NewSize(400, 500))

	views.ShowMainWindow(window)

	window.ShowAndRun()
}
