package methods

import (
	"errors"
	"math"
	"strconv"
	"strings"
)

func NumericalIntegration(X, Y []float64) (map[string]float64, error) {
	n := len(X)
	if n < 2 || len(Y) < 2 {
		return nil, errors.New("недостаточно точек для вычисления (требуется как минимум 2)")
	}
	if n != len(Y) {
		return nil, errors.New("количества координат X и Y не совпадают")
	}

	h := X[1] - X[0]

	precision := detectPrecision(Y) + 2

	m := n - 1

	if m > 6 {
		return nil, errors.New("количество узлов больше 7 (порядок m > 6) пока не поддерживается формулой Ньютона-Котеса")
	}

	ncResult, err := newtonCotes(X[0], X[m], m, Y)
	if err != nil {
		return nil, err
	}

	results := map[string]float64{
		"rect_left":    roundTo(rectangleLeftMethod(h, Y), precision),
		"rect_right":   roundTo(rectangleRightMethod(h, Y), precision),
		"trapezoidal":  roundTo(trapezoidalMethod(h, Y), precision),
		"newton_cotes": roundTo(ncResult, precision),
	}

	if n >= 3 && m%2 == 0 {
		results["simpson"] = roundTo(simpsonMethod(h, Y), precision)
	}

	return results, nil
}

func detectPrecision(values []float64) int {
	maxDigits := 0
	for _, val := range values {
		str := strconv.FormatFloat(val, 'f', -1, 64)
		if idx := strings.Index(str, "."); idx != -1 {
			digits := len(str) - idx - 1
			if digits > maxDigits {
				maxDigits = digits
			}
		}
	}
	return maxDigits
}

func roundTo(val float64, precision int) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}

func newtonCotes(a, b float64, m int, Y []float64) (float64, error) {
	coeffs, err := getCotesCoefficients(m, a, b)
	if err != nil {
		return 0, err
	}

	var integral float64
	for i := 0; i <= m; i++ {
		integral += coeffs[i] * Y[i]
	}
	return integral, nil
}

func getCotesCoefficients(m int, a, b float64) ([]float64, error) {
	width := b - a
	coeffs := make([]float64, m+1)

	switch m {
	case 1:
		coeffs[0] = width / 2.0
		coeffs[1] = width / 2.0
	case 2:
		coeffs[0] = width / 6.0
		coeffs[2] = width / 6.0
		coeffs[1] = (4.0 * width) / 6.0
	case 3:
		coeffs[0] = width / 8.0
		coeffs[3] = width / 8.0
		coeffs[1] = (3.0 * width) / 8.0
		coeffs[2] = (3.0 * width) / 8.0
	case 4:
		coeffs[0] = (7.0 * width) / 90.0
		coeffs[4] = (7.0 * width) / 90.0
		coeffs[1] = (16.0 * width) / 45.0
		coeffs[3] = (16.0 * width) / 45.0
		coeffs[2] = (2.0 * width) / 15.0
	case 5:
		coeffs[0] = (19.0 * width) / 288.0
		coeffs[5] = (19.0 * width) / 288.0
		coeffs[1] = (25.0 * width) / 96.0
		coeffs[4] = (25.0 * width) / 96.0
		coeffs[2] = (25.0 * width) / 144.0
		coeffs[3] = (25.0 * width) / 144.0
	case 6:
		coeffs[0] = (41.0 * width) / 840.0
		coeffs[6] = (41.0 * width) / 840.0
		coeffs[1] = (9.0 * width) / 35.0
		coeffs[5] = (9.0 * width) / 35.0
		coeffs[2] = (9.0 * width) / 280.0
		coeffs[4] = (9.0 * width) / 280.0
		coeffs[3] = (34.0 * width) / 105.0
	default:
		return nil, errors.New("высокий порядок m не поддерживается данной таблицей")
	}
	return coeffs, nil
}

func rectangleLeftMethod(h float64, Y []float64) float64 {
	var sum float64
	for i := 0; i < len(Y)-1; i++ {
		sum += Y[i]
	}
	return h * sum
}

func rectangleRightMethod(h float64, Y []float64) float64 {
	var sum float64
	for i := 1; i < len(Y); i++ {
		sum += Y[i]
	}
	return h * sum
}

func trapezoidalMethod(h float64, Y []float64) float64 {
	n := len(Y)
	var sum float64
	for i := 1; i < n-1; i++ {
		sum += Y[i]
	}
	totalSum := ((Y[0] + Y[n-1]) / 2.0) + sum
	return h * totalSum
}

func simpsonMethod(h float64, Y []float64) float64 {
	n := len(Y)
	sum := Y[0] + Y[n-1]

	for i := 1; i < n-1; i++ {
		if i%2 == 1 {
			sum += 4.0 * Y[i]
		} else {
			sum += 2.0 * Y[i]
		}
	}
	return (h / 3.0) * sum
}
