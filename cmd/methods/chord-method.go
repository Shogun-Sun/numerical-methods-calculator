package methods

import (
	"errors"
	"math"
)

type IterationStep struct {
	Num   int
	Xn    float64
	Fxn   float64
	NextX float64
}

func ChordMethod(f func(float64) float64, x0, x1, eps float64, maxIterations int) (float64, []IterationStep, error) {
	fx0 := f(x0)
	xn := x1
	var steps []IterationStep

	for i := 1; i <= maxIterations; i++ {
		fxn := f(xn)

		if math.Abs(fxn-fx0) < 1e-15 {
			return 0, steps, errors.New("деление на ноль: значения функции в точках совпали")
		}

		nextX := xn - ((xn-x0)/(fxn-fx0))*fxn

		steps = append(steps, IterationStep{
			Num:   i,
			Xn:    xn,
			Fxn:   fxn,
			NextX: nextX,
		})

		if math.Abs(nextX-xn) < eps {
			return nextX, steps, nil
		}

		xn = nextX
	}

	return xn, steps, errors.New("превышено максимальное количество итераций")
}
