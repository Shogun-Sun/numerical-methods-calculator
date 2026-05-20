package methods

import (
	"errors"
	"math"
	"strconv"
	"strings"
)

func NumericalIntegration(X, Y []float64) (map[string]float64, error) {
	if len(X) < 5 || len(Y) < 5 {
		return nil, errors.New("Недостаточно точек для вычисления (требуется как минимум 5)")
	}

	//Шаг сетки
	h := X[1] - X[0]

	//Определение точности
	precision := detectPrecision(Y) + 2

	// Вычисление по Ньютону-Котессу:
	ncResult, err := newtonCotes(X[0], X[4], 4, Y)
	if err != nil {
		return nil, err
	}

	results := map[string]float64{
		"rect_left":    roundTo(rectangleLeftMethod(h, Y), precision),
		"rect_right":   roundTo(rectangleRightMethod(h, Y), precision),
		"trapezoidal":  roundTo(trapezoidalMethod(h, Y), precision),
		"simpson":      roundTo(simpsonMethod(h, Y), precision),
		"newton_cotes": roundTo(ncResult, precision),
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
		// c_1^0 = c_1^1 = (b-a)/2
		coeffs[0] = width / 2.0
		coeffs[1] = width / 2.0

	case 2:
		// c_2^0 = c_2^2 = (b-a)/6; c_2^1 = 4(b-a)/6
		coeffs[0] = width / 6.0
		coeffs[2] = width / 6.0
		coeffs[1] = (4.0 * width) / 6.0

	case 3:
		// c_3^0 = c_3^3 = (b-a)/8; c_3^1 = c_3^2 = 3(b-a)/8
		coeffs[0] = width / 8.0
		coeffs[3] = width / 8.0
		coeffs[1] = (3.0 * width) / 8.0
		coeffs[2] = (3.0 * width) / 8.0

	case 4:
		// c_4^0 = c_4^4 = 7(b-a)/90; c_4^1 = c_4^3 = 16(b-a)/45; c_4^2 = 2(b-a)/15
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

// Метод левых прямоугольников (использует y0, y1, y2, y3)
func rectangleLeftMethod(h float64, Y []float64) float64 {
	sum := Y[0] + Y[1] + Y[2] + Y[3]
	return h * sum
}

// Метод правых прямоугольников (использует y1, y2, y3, y4)
func rectangleRightMethod(h float64, Y []float64) float64 {
	sum := Y[1] + Y[2] + Y[3] + Y[4]
	return h * sum
}

// Метод трапеций (составная формула)
func trapezoidalMethod(h float64, Y []float64) float64 {
	sum := ((Y[0] + Y[4]) / 2.0) + Y[1] + Y[2] + Y[3]
	return h * sum
}

// Метод Симпсона (парабол)
func simpsonMethod(h float64, Y []float64) float64 {
	sum := Y[0] + 4.0*(Y[1]+Y[3]) + 2.0*Y[2] + Y[4]
	return (h / 3.0) * sum
}
