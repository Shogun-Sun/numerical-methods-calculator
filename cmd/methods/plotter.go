package methods

import (
	"fmt"
	"math"

	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type InteractiveChart struct {
	widget.BaseWidget
	f      func(float64) float64
	minX   float64
	maxX   float64
	minY   float64
	maxY   float64
	width  float32
	height float32
}

func NewInteractiveChart(f func(float64) float64, width, height float32) *InteractiveChart {
	chart := &InteractiveChart{
		f:      f,
		minX:   -4.0,
		maxX:   4.0,
		minY:   -8.0,
		maxY:   8.0,
		width:  width,
		height: height,
	}
	if math.IsNaN(f(0.001)) || math.IsInf(f(0.001), 0) {
		chart.minX = 0.1
		chart.maxX = 8.1
	}
	chart.ExtendBaseWidget(chart)
	return chart
}

func (c *InteractiveChart) CreateRenderer() fyne.WidgetRenderer {
	container := fyne.NewContainerWithoutLayout()
	r := &chartRenderer{chart: c, container: container}
	r.updateGraph()
	return r
}

func (c *InteractiveChart) Scrolled(ev *fyne.ScrollEvent) {
	zoomFactor := 0.1
	if ev.Scrolled.DY < 0 {
		zoomFactor = -0.1
	}

	dx := c.maxX - c.minX
	dy := c.maxY - c.minY

	c.minX += dx * zoomFactor * 0.5
	c.maxX -= dx * zoomFactor * 0.5
	c.minY += dy * zoomFactor * 0.5
	c.maxY -= dy * zoomFactor * 0.5

	if math.IsNaN(c.f(0.001)) || math.IsInf(c.f(0.001), 0) {
		if c.minX < 0.01 {
			c.minX = 0.01
		}
	}

	c.Refresh()
}

func (c *InteractiveChart) Dragged(ev *fyne.DragEvent) {
	dx := c.maxX - c.minX
	dy := c.maxY - c.minY

	moveX := (float64(ev.Dragged.DX) / float64(c.width)) * dx
	moveY := (float64(ev.Dragged.DY) / float64(c.height)) * dy

	c.minX -= moveX
	c.maxX -= moveX
	c.minY += moveY
	c.maxY += moveY

	c.Refresh()
}

func (c *InteractiveChart) DragEnd() {}

type chartRenderer struct {
	chart     *InteractiveChart
	container *fyne.Container
}

func (r *chartRenderer) Destroy() {}

func (r *chartRenderer) Layout(size fyne.Size) {
	r.chart.width = size.Width
	r.chart.height = size.Height
	r.container.Resize(size)
	r.updateGraph()
}

func (r *chartRenderer) MinSize() fyne.Size {
	return fyne.NewSize(r.chart.width, r.chart.height)
}

func (r *chartRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.container}
}

func (r *chartRenderer) Refresh() {
	r.updateGraph()
}

func (r *chartRenderer) updateGraph() {
	r.container.Objects = nil

	c := r.chart
	wStr, hStr := float64(c.width), float64(c.height)

	axisColor := theme.ForegroundColor()

	var gridColor color.Color
	if rgba, ok := axisColor.(color.RGBA); ok {
		gridColor = color.RGBA{R: rgba.R, G: rgba.G, B: rgba.B, A: 40}
	} else if nrgba, ok := axisColor.(color.NRGBA); ok {
		gridColor = color.NRGBA{R: nrgba.R, G: nrgba.G, B: nrgba.B, A: 40}
	} else {
		gridColor = color.RGBA{R: 128, G: 128, B: 128, A: 40}
	}

	lineColor := theme.PrimaryColor()

	// 1. Рисуем сетку и подписи по оси X (динамический шаг)
	stepX := math.Pow(10, math.Floor(math.Log10(c.maxX-c.minX))-1)
	if (c.maxX-c.minX)/stepX > 15 {
		stepX *= 2
	}
	startX := math.Ceil(c.minX/stepX) * stepX
	for xVal := startX; xVal <= c.maxX; xVal += stepX {
		screenX := ((xVal - c.minX) / (c.maxX - c.minX)) * wStr
		if math.Abs(xVal) > 1e-9 {
			gridLine := canvas.NewLine(gridColor)
			gridLine.Position1 = fyne.NewPos(float32(screenX), 0)
			gridLine.Position2 = fyne.NewPos(float32(screenX), c.height)
			r.container.Add(gridLine)

			labelX := canvas.NewText(fmt.Sprintf("%.2g", xVal), axisColor)
			labelX.TextSize = 10
			zeroYReal := hStr - ((0.0 - c.minY) / (c.maxY - c.minY) * hStr)
			if zeroYReal < 10 {
				zeroYReal = 10
			}
			if zeroYReal > hStr-20 {
				zeroYReal = hStr - 20
			}
			labelX.Move(fyne.NewPos(float32(screenX)-10, float32(zeroYReal)+5))
			r.container.Add(labelX)
		}
	}

	stepY := math.Pow(10, math.Floor(math.Log10(c.maxY-c.minY))-1)
	if (c.maxY-c.minY)/stepY > 15 {
		stepY *= 2
	}
	startY := math.Ceil(c.minY/stepY) * stepY
	for yVal := startY; yVal <= c.maxY; yVal += stepY {
		screenY := hStr - ((yVal - c.minY) / (c.maxY - c.minY) * hStr)
		if math.Abs(yVal) > 1e-9 {
			gridLine := canvas.NewLine(gridColor)
			gridLine.Position1 = fyne.NewPos(0, float32(screenY))
			gridLine.Position2 = fyne.NewPos(c.width, float32(screenY))
			r.container.Add(gridLine)

			labelY := canvas.NewText(fmt.Sprintf("%.2g", yVal), axisColor)
			labelY.TextSize = 10
			zeroXReal := ((0.0 - c.minX) / (c.maxX - c.minX) * wStr)
			if zeroXReal < 5 {
				zeroXReal = 5
			}
			if zeroXReal > wStr-30 {
				zeroXReal = wStr - 30
			}
			labelY.Move(fyne.NewPos(float32(zeroXReal)-25, float32(screenY)-7))
			r.container.Add(labelY)
		}
	}

	// 3. Главные оси координат (X и Y)
	if c.minY <= 0 && c.maxY >= 0 {
		zeroYReal := hStr - ((0.0 - c.minY) / (c.maxY - c.minY) * hStr)
		xAxis := canvas.NewLine(axisColor)
		xAxis.Position1 = fyne.NewPos(0, float32(zeroYReal))
		xAxis.Position2 = fyne.NewPos(c.width, float32(zeroYReal))
		xAxis.StrokeWidth = 2
		r.container.Add(xAxis)
	}
	if c.minX <= 0 && c.maxX >= 0 {
		zeroXReal := ((0.0 - c.minX) / (c.maxX - c.minX) * wStr)
		yAxis := canvas.NewLine(axisColor)
		yAxis.Position1 = fyne.NewPos(float32(zeroXReal), 0)
		yAxis.Position2 = fyne.NewPos(float32(zeroXReal), c.height)
		yAxis.StrokeWidth = 2
		r.container.Add(yAxis)
	}

	fixedLineColor := color.RGBA{R: 244, G: 67, B: 54, A: 130}
	screenY1 := hStr - ((c.minX - c.minY) / (c.maxY - c.minY) * hStr)
	screenY2 := hStr - ((c.maxX - c.minY) / (c.maxY - c.minY) * hStr)

	yEqualsXLine := canvas.NewLine(fixedLineColor)
	yEqualsXLine.Position1 = fyne.NewPos(0, float32(screenY1))
	yEqualsXLine.Position2 = fyne.NewPos(c.width, float32(screenY2))
	yEqualsXLine.StrokeWidth = 1.5
	r.container.Add(yEqualsXLine)

	labelYX := canvas.NewText("y = x", fixedLineColor)
	labelYX.TextSize = 11
	labelYX.Move(fyne.NewPos(c.width-45, float32(screenY2)-15))
	r.container.Add(labelYX)

	type rootPoint struct {
		mathX   float64
		screenX float64
		screenY float64
	}
	var axisRoots []rootPoint
	var fixedRoots []rootPoint

	var prevX, prevY float64
	firstPoint := true

	for screenX := 0.0; screenX <= wStr; screenX += 1.0 {
		mathX := c.minX + (screenX/wStr)*(c.maxX-c.minX)
		mathY := c.f(mathX)

		if math.IsNaN(mathY) || math.IsInf(mathY, 0) {
			firstPoint = true
			continue
		}

		screenY := hStr - ((mathY - c.minY) / (c.maxY - c.minY) * hStr)
		if screenY < -50 || screenY > hStr+50 {
			firstPoint = true
			continue
		}

		if !firstPoint {
			line := canvas.NewLine(lineColor)
			line.Position1 = fyne.NewPos(float32(prevX), float32(prevY))
			line.Position2 = fyne.NewPos(float32(screenX), float32(screenY))
			line.StrokeWidth = 3
			r.container.Add(line)

			prevMathX := c.minX + ((screenX-1.0)/wStr)*(c.maxX-c.minX)
			prevMathY := c.f(prevMathX)

			if !math.IsNaN(prevMathY) && !math.IsInf(prevMathY, 0) {
				if (prevMathY <= 0 && mathY > 0) || (prevMathY >= 0 && mathY < 0) {
					t := -prevMathY / (mathY - prevMathY)
					exactMathX := prevMathX + t*(mathX-prevMathX)
					exactScreenX := prevX + t*(screenX-prevX)
					exactScreenY := hStr - ((0.0 - c.minY) / (c.maxY - c.minY) * hStr)

					axisRoots = append(axisRoots, rootPoint{
						mathX:   exactMathX,
						screenX: exactScreenX,
						screenY: exactScreenY,
					})
				}

				currentDiff := mathY - mathX
				prevDiff := prevMathY - prevMathX

				if (prevDiff <= 0 && currentDiff > 0) || (prevDiff >= 0 && currentDiff < 0) {
					t := -prevDiff / (currentDiff - prevDiff)
					exactMathX := prevMathX + t*(mathX-prevMathX)
					exactScreenX := prevX + t*(screenX-prevX)
					exactScreenY := hStr - ((exactMathX - c.minY) / (c.maxY - c.minY) * hStr)

					fixedRoots = append(fixedRoots, rootPoint{
						mathX:   exactMathX,
						screenX: exactScreenX,
						screenY: exactScreenY,
					})
				}
			}
		}

		prevX = screenX
		prevY = screenY
		firstPoint = false
	}

	for _, root := range axisRoots {
		rootMarker := canvas.NewCircle(theme.ErrorColor())
		rootMarker.Resize(fyne.NewSize(8, 8))
		rootMarker.Move(fyne.NewPos(float32(root.screenX)-4, float32(root.screenY)-4))
		r.container.Add(rootMarker)

		rootText := canvas.NewText(fmt.Sprintf("x₀ ≈ %.2f", root.mathX), axisColor)
		rootText.TextSize = 10
		rootText.TextStyle = fyne.TextStyle{Bold: true}
		rootText.Move(fyne.NewPos(float32(root.screenX)-18, float32(root.screenY)-22))
		r.container.Add(rootText)
	}

	for _, root := range fixedRoots {
		rootMarker := canvas.NewCircle(theme.WarningColor())
		rootMarker.Resize(fyne.NewSize(8, 8))
		rootMarker.Move(fyne.NewPos(float32(root.screenX)-4, float32(root.screenY)-4))
		r.container.Add(rootMarker)

		rootText := canvas.NewText(fmt.Sprintf("x* ≈ %.2f", root.mathX), axisColor)
		rootText.TextSize = 10
		rootText.TextStyle = fyne.TextStyle{Bold: true}
		rootText.Move(fyne.NewPos(float32(root.screenX)-18, float32(root.screenY)+10))
		r.container.Add(rootText)
	}

	canvas.Refresh(r.container)
}
