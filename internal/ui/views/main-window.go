package views

import (
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
)

func MainWindow() {
	a := app.New()
	w := a.NewWindow("Численне методы")

	w.SetContent(widget.NewLabel("Добро поаловать!"))

	w.Show()

	a.Run()
}
