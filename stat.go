package gostat

import (
	"github.com/gonum/stat"
	"math"
	"sort"
)

// MAD returns a median absolute deviation (MAD) product.
//
// The MAD algorithm, a recognized statistical methodology, is a robust analog
// to a more commonly used outlier technique, which uses standard deviation
// from the mean. MAD uses deviation from the median, which is less susceptible
// to distortion caused by outlying values.
//
// MAD is determined by calculating the deviation from the median as follows:
//
// 1. subtracting each charge in the distribution from the median charge to determine each respective deviation from the median;
//
// 2. the absolute values of the deviations from the median are arrayed in order from lowest to highest and the median of the absolute deviations is determined;
//
// 3. the median of the absolute deviations (from the median) is multiplied by the constant of 1.4826;
//
// 4. this product is defined as the MAD.
func MAD(x []float64) float64 {
	if len(x) == 0 {
		return -1.0
	}
	median := Median(x)
	series := make([]float64, len(x))
	for i := 0; i < len(x); i++ {
		series[i] = math.Abs(median - x[i])
	}

	return 1.4826 * Median(series)
}

// Median returns the median by arraying the data for a given slice
// from lowest to highest and identifying the value at which half of the data
// are higher and half are lower
func Median(x []float64) float64 {
	series := append([]float64{}, x...)
	sort.Float64s(series)

	k := len(series) / 2

	if len(series)%2 == 1 {
		return series[k]
	}

	return 0.5 * (series[k-1] + series[k])
}

// MovStdDev returns moving standard deviation, a slice of local k-point
// standard deviation values, where each standard deviation is calculated
// over a sliding window of length k across neighboring elements of x.
// Set center to true for center moving standard deviation or to false
// for trailing moving standard deviation.
func MovStdDev(x, weights []float64, k int, omitNaNs, trailing, fullWnd bool) []float64 {
	rolling := RollingWindow(x, k, omitNaNs, trailing, fullWnd)
	stdDevs := make([]float64, len(rolling))
	for i := 0; i < len(rolling); i++ {
		stdDevs[i] = stat.StdDev(rolling[i], weights)
	}

	return stdDevs
}

// Volatility calculates historical volatility as annualized standard
// deviation of logarithmic returns
func Volatility(x []float64, periodicity float64) float64 {
	var rets []float64
	// calculate returns
	for i := 1; i < len(x); i++ {
		rets = append(rets, math.Log(x[i]/x[i-1]))
	}
	stdev := stat.StdDev(rets, nil)
	return stdev * math.Sqrt(periodicity)
}

// Normalize is normalizing a set of scores x using the standard deviation.
// This normalization is known as Z-scores. With elementary algebraic
// manipulations, it can be shown that a set of Z-score has a mean equal of
// zero and a standard deviation of one. Therefore, Z-scores constitute an unit
// free measure which can be used to compare observations measured with
// different units.
func Normalize(x, weights []float64) []float64 {
	zscores := make([]float64, len(x))
	mean := stat.Mean(x, weights)
	stdDev := stat.StdDev(x, weights)
	for i := 0; i < len(x); i++ {
		if stdDev != 0.0 {
			zscores[i] = (x[i] - mean) / stdDev
		} else {
			zscores[i] = x[i] - mean
		}
	}
	return zscores
}

// RollingWindow splits slice x into a sliding window of length k. The window
// size is automatically truncated at the endpoints when there are not enough
// elements to fill the window.
//
// - trailing - if there are more windows than length of x, do not select center windows
//
// - omitNaNs - omit NaN values
//
// - fullWnd  - discard any window that uses fewer elements than k
func RollingWindow(x []float64, k int, omitNaNs, trailing, fullWnd bool) [][]float64 {

	var rets [][]float64
	var v []float64

	if omitNaNs {
		v = filterNaNs(x)
	} else {
		v = x
	}

	if !fullWnd {
		for i := 1; i < k; i++ {
			rets = append(rets, v[:i])
		}
	}

	for i := 0; i+k-1 < len(v); i++ {
		rets = append(rets, v[i:i+k])
	}

	if !fullWnd {
		for i := k - 1; i > 0; i-- {
			rets = append(rets, v[len(v)-i:])
		}
	}

	if len(x) >= len(rets) {
		return rets
	}

	if !trailing {
		diff := len(rets) - len(v)
		var trim int
		if math.Mod(float64(diff), 2.) == 0 {
			trim = diff / 2
		} else {
			trim = (diff - 1) / 2
		}
		return rets[trim : len(x)+trim]
	}

	return rets[:len(x)]
}

func filterNaNs(x []float64) []float64 {
	var v []float64
	for i := 0; i < len(x); i++ {
		if isRealVal(x[i]) {
			v = append(v, x[i])
		}
	}
	return v
}

func isRealVal(x float64) bool {
	return !math.IsNaN(x) && !math.IsInf(x, 0)
}
