package util

import (
	"math"
)

// MovingWindow is a moving window
type MovingWindow struct {
	index    int
	length   int
	capacity int

	variance float64
	stddev   float64

	sum     float64
	average float64

	pool []float64
}

// NewMovingWindow returns a new moving window.
func NewMovingWindow(size int) MovingWindow {
	return MovingWindow{
		pool:     make([]float64, size),
		capacity: size,
	}
}

func (mw *MovingWindow) calcFinal() (float64, float64) {
	if mw.length > 1 {
		// mw.stddev = math.Sqrt(mw.variance / (mw.length - 1))
		// okay so this came from dpayne/cli-visualizer
		mw.stddev = (mw.variance / float64(mw.length-1)) - math.Pow(mw.average, 2)
		mw.stddev = math.Sqrt(math.Abs(mw.stddev))
	} else {
		mw.stddev = 0
	}

	if mw.length > 0 {
		mw.average = mw.sum / float64(mw.length)
	} else {
		mw.average = 0
	}

	return mw.average, mw.stddev
}

// Update updates the moving window
func (mw *MovingWindow) Update(value float64) (float64, float64) {
	if mw.length < mw.capacity {

		mw.length++

		mw.sum += value
		mw.variance += math.Pow(value, 2)

	} else {
		mw.sum += value - mw.pool[mw.index]
		mw.variance += math.Pow(value, 2) - math.Pow(mw.pool[mw.index], 2)
	}

	mw.pool[mw.index] = value

	if mw.index++; mw.index >= mw.capacity {
		mw.index = 0
	}

	return mw.calcFinal()
}

// Drop removes count items from the window
// TODO(winter): look into a better index calculation
func (mw *MovingWindow) Drop(count int) (float64, float64) {
	if mw.length <= 0 {
		return mw.calcFinal()
	}

	for count > 0 && mw.length > 0 {

		var idx = (mw.index - mw.length)
		if idx < 0 {
			idx = mw.capacity + idx
		}

		mw.sum -= mw.pool[idx]
		mw.variance -= math.Pow(mw.pool[idx], 2)

		mw.length--

		count--
	}

	// If we dont have enough length for standard dev, clear variance
	if mw.length < 2 {
		mw.variance = 0
		if mw.length < 1 {
			mw.length = 0
			// same idea with sum. just clear it so we dont have a rouding issue
			mw.sum = 0
		}
	}

	return mw.calcFinal()
}

// Len returns how many items in the window
func (mw *MovingWindow) Len() int {
	// logical length
	return mw.length
}

// Cap returns max size of window
func (mw *MovingWindow) Cap() int {
	return mw.capacity
}

// Mean is the moving window average
func (mw *MovingWindow) Mean() float64 {
	return mw.average
}

// StdDev is the moving average std
func (mw *MovingWindow) StdDev() float64 {
	return mw.stddev
}

// Stats returns the statistics of this window
func (mw *MovingWindow) Stats() (float64, float64) {
	return mw.average, mw.stddev
}
