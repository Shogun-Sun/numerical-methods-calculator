package widgets

import (
	"fyne.io/fyne/v2/widget"
)

func ButtonTheme(text string, onClick func()) *widget.Button {
	btn := widget.NewButton(text, onClick)
	return btn
}
