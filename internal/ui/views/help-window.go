package views

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func openHelpWindow() {
	app := fyne.CurrentApp()
	helpWin := app.NewWindow("Справка: Ввод уравнений")

	helpText := widget.NewRichTextFromMarkdown(`
# Инструкция по вводу уравнений

Для корректной работы парсера придерживайтесь следующих правил:

* **Знак умножения (*):** Компьютер не понимает пропущенных знаков. Вместо 2x всегда пишите: 2 * x
* **Степени (pow):** Используйте функцию pow(база, степень). 
    * *Пример:* x^2 пишется как pow(x, 2)
* **Функции со скобками:** Всегда оборачивайте аргументы функций в круглые скобки.
    * ln(x) вместо lnx
    * sin(x) вместо sinx

### Примеры правильного ввода:
1. **ln(x) + 2x² - 6 = 0** ➔  ln(x) + 2 * pow(x, 2) - 6
2. **3sin(x) + x² - 1 = 0** ➔  3 * sin(x) + pow(x, 2) - 1
3. **2ln²(x) + 2x² - 3 = 0** ➔  2 * pow(ln(x), 2) + 2 * pow(x, 2) - 3
	`)

	closeBtn := widget.NewButton("Понятно", func() {
		helpWin.Close()
	})

	content := container.NewVBox(
		helpText,
		widget.NewSeparator(),
		closeBtn,
	)

	helpWin.SetContent(content)
	helpWin.Resize(fyne.NewSize(480, 430))

	helpWin.Show()
}
