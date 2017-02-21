package gostat

import (
	"github.com/gonum/floats"
	"github.com/gonum/stat"
	"math"
	"sort"
	"testing"
)

func TestMedianValue_Odd(t *testing.T) {
	series := []float64{2., 1., 3.}
	median := Median(series)

	if got, want := median, 2.0; got != want {
		t.Errorf("Expected median=%f, got=%f", want, got)
	}

	if sort.Float64sAreSorted(series) {
		t.Errorf("slice should not got sorted %v", series)
	}
}

func TestMedianValue_Even(t *testing.T) {
	series := []float64{2., 1., 4., 3.}
	median := Median(series)

	if got, want := median, 2.5; got != want {
		t.Errorf("Expected median=%f, got=%f", want, got)
	}

	if sort.Float64sAreSorted(series) {
		t.Errorf("slice should not got sorted %v", series)
	}
}

func TestMAD(t *testing.T) {
	x := []float64{2., 6., 6., 12., 17., 25., 32.}
	mad := MAD(x)

	if got, want := mad, 8.8956; got != want {
		t.Errorf("Expected MAD=%f, got=%f", want, got)
	}
}

func TestMAD_Empty(t *testing.T) {
	x := []float64{}
	mad := MAD(x)

	if got, want := mad, -1.; got != want {
		t.Errorf("Expected MAD=%f, got=%f", want, got)
	}
}

func TestRollingWindow(t *testing.T) {
	x := []float64{1., 2., 3., 4., 5.}
	rolling := RollingWindow(x, 3, false, false, false)
	if got, want := len(rolling), len(x); got != want {
		t.Errorf("Expected number of elements=%d, got=%d", want, got)
	}
	compareArrays([]float64{1., 2.}, rolling[0], t)
	compareArrays([]float64{1., 2., 3.}, rolling[1], t)
	compareArrays([]float64{2., 3., 4.}, rolling[2], t)
	compareArrays([]float64{3., 4., 5.}, rolling[3], t)
	compareArrays([]float64{4., 5.}, rolling[4], t)
}

func TestRollingWindow_EvenWindow(t *testing.T) {
	x := []float64{1., 2., 3., 4., 5.}
	rolling := RollingWindow(x, 2, false, false, false)
	if got, want := len(rolling), len(x); got != want {
		t.Errorf("Expected number of elements=%d, got=%d", want, got)
	}
	compareArrays([]float64{1.}, rolling[0], t)
	compareArrays([]float64{1., 2.}, rolling[1], t)
	compareArrays([]float64{2., 3.}, rolling[2], t)
	compareArrays([]float64{3., 4.}, rolling[3], t)
	compareArrays([]float64{4., 5.}, rolling[4], t)
}

func TestRollingWindow_Singleton(t *testing.T) {
	x := []float64{1., 2., 3., 4., 5.}
	rolling := RollingWindow(x, 1, false, false, false)
	if got, want := len(rolling), len(x); got != want {
		t.Errorf("Expected number of elements=%d, got=%d", want, got)
	}
	compareArrays([]float64{1.}, rolling[0], t)
	compareArrays([]float64{2.}, rolling[1], t)
	compareArrays([]float64{3.}, rolling[2], t)
	compareArrays([]float64{4.}, rolling[3], t)
	compareArrays([]float64{5.}, rolling[4], t)
}

func TestRollingWindow_OmitNaN(t *testing.T) {
	x := []float64{1., math.NaN(), 3., math.Inf(-1), 5.}
	rolling := RollingWindow(x, 2, true, false, false)
	if got, want := len(rolling), 4; got != want {
		t.Errorf("Expected number of elements=%d, got=%d", want, got)
	}
	compareArrays([]float64{1.}, rolling[0], t)
	compareArrays([]float64{1., 3.}, rolling[1], t)
	compareArrays([]float64{3., 5.}, rolling[2], t)
	compareArrays([]float64{5.}, rolling[3], t)
}

func TestRollingWindow_FullWindow(t *testing.T) {
	x := []float64{1., 2., 3., 4., 5.}
	rolling := RollingWindow(x, 3, false, false, true)
	if got, want := len(rolling), 3; got != want {
		t.Errorf("Expected number of elements=%d, got=%d", want, got)
	}
	compareArrays([]float64{1., 2., 3.}, rolling[0], t)
	compareArrays([]float64{2., 3., 4.}, rolling[1], t)
	compareArrays([]float64{3., 4., 5.}, rolling[2], t)
}

func TestRollingWindow_Trailing(t *testing.T) {
	x := []float64{1., 2., 3., 4., 5.}
	rolling := RollingWindow(x, 3, false, true, false)
	if got, want := len(rolling), len(x); got != want {
		t.Errorf("Expected number of elements=%d, got=%d", want, got)
	}
	compareArrays([]float64{1.}, rolling[0], t)
	compareArrays([]float64{1., 2.}, rolling[1], t)
	compareArrays([]float64{1., 2., 3.}, rolling[2], t)
	compareArrays([]float64{2., 3., 4.}, rolling[3], t)
	compareArrays([]float64{3., 4., 5.}, rolling[4], t)
}

func TestMovStdDev(t *testing.T) {
	x := []float64{4., 8., 6., -1., -2., -3., -1., 3., 4., 5.}
	m := MovStdDev(x, nil, 3, false, false, false)
	compareArrays([]float64{2.8284, 2., 4.7258, 4.3589, 1., 1., 3.0551, 2.6458, 1., 0.7071}, m, t)
}

func TestMovStdDev_WithNaNs(t *testing.T) {
	x := []float64{4., 8., math.NaN(), -1., -2., -3., math.NaN(), 3., 4., 5.}
	m := MovStdDev(x, nil, 3, false, false, false)
	compareArrays([]float64{2.8284, math.NaN(), math.NaN(), math.NaN(), 1., math.NaN(), math.NaN(), math.NaN(), 1., 0.7071}, m, t)
}

func TestMovStdDev_Trailing(t *testing.T) {
	x := []float64{4., 8., 6., -1., -2., -3., -1., 3., 4., 5.}
	m := MovStdDev(x, nil, 3, false, true, false)
	compareArrays([]float64{math.NaN(), 2.8284, 2., 4.7258, 4.3589, 1., 1., 3.0551, 2.6458, 1.}, m, t)
}

func TestMovStdDev_FullWindow(t *testing.T) {
	x := []float64{4., 8., 6., -1., -2., -3., -1., 3., 4., 5.}
	m := MovStdDev(x, nil, 3, false, false, true)
	compareArrays([]float64{2., 4.7258, 4.3589, 1., 1., 3.0551, 2.6458, 1.}, m, t)
}

func TestMovStdDev_EqualSize(t *testing.T) {
	x := []float64{4., 8., 6.}
	m := MovStdDev(x, nil, 3, false, false, false)
	compareArrays([]float64{2.8284, 2.0000, 1.4142}, m, t)
}

func compareArrays(x, y []float64, t *testing.T) {
	if len(x) != len(y) {
		t.Fatalf("Expected number of elements=%d, got=%d", len(x), len(y))
	}
	for i := 0; i < len(x); i++ {
		if got, want := y[i], x[i]; !floatEquals(got, want) {
			t.Errorf("Expected value at index %d=%f, got=%f", i, want, got)
		}
	}
}

func TestVolatility(t *testing.T) {
	// 22 prices to calculate 21 returns - roughly corresponds to 1 calendar month
	prices := []float64{42.35834, 40.703716, 42.202611, 42.338873, 41.47263,
		42.718463, 41.920351, 42.13448, 42.319407, 41.891153,
		42.80606, 43.117518, 43.068854, 42.319407, 42.932591,
		42.728198, 42.698996, 42.737929, 42.767127, 42.13448,
		42.280473, 43.078585}
	// long-term average for US markets is 252 trading days per year
	vol := Volatility(prices, 252.)
	if got, want := vol, 0.2827; !floatEquals(got, want) {
		t.Errorf("Expected volatility=%f, got=%f", want, got)
	}
}

func TestNormalize(t *testing.T) {
	scores := []float64{35., 36., 46., 68., 70.}
	zscores := Normalize(scores, nil)
	compareArrays([]float64{-0.9412, -0.8824, -0.2941, 1.0000, 1.1176}, zscores, t)
	if got, want := stat.Mean(zscores, nil), 0.; !floatEquals(got, want) {
		t.Errorf("Expected mean=%f, got=%f", want, got)
	}
	if got, want := stat.StdDev(zscores, nil), 1.; !floatEquals(got, want) {
		t.Errorf("Expected standard deviation=%f, got=%f", want, got)
	}
}

const epsilon float64 = 0.0001

func floatEquals(a, b float64) bool {
	if math.IsNaN(a) && math.IsNaN(b) {
		return true
	}
	return floats.EqualWithinAbs(a, b, epsilon)
}
