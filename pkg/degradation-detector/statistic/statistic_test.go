package statistic

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShamosCoefficient10(t *testing.T) {
	assert.InDelta(t, 1.000643, shamosBias(10), 0.00001)
}

func TestShamosCoefficient100(t *testing.T) {
	assert.InDelta(t, 1.043987, shamosBias(100), 0.00001)
}

func TestShamosCoefficient200(t *testing.T) {
	assert.InDelta(t, 1.046203, shamosBias(200), 0.001)
}

func TestShamosEstimator(t *testing.T) {
	data := []int{
		5691, 5855, 5720, 6339, 5829, 5496, 5427, 5586, 5859, 5603, 5868, 5761, 5440, 5590, 5870, 5781, 5632, 6092, 5636, 5849, 5730, 5639, 5678, 5857, 5655, 5486, 5877, 5639, 5668,
		5864, 5602, 5855, 6049, 5741, 5794, 5822, 5704, 5707, 6167, 5923, 5765, 5648, 5775, 5578, 5541, 5919, 5498, 5436, 5857, 5508, 5739, 5820, 5662, 5582, 5565, 5708, 5587, 5813,
		5618, 5796, 5682, 5778, 5848, 6034, 5847, 5653, 5783, 6006, 5647, 5509, 5869, 5738, 5709, 5762, 5793, 5607, 5620, 5580, 5710, 5641, 5673, 5794, 5937, 5708, 5705, 5747, 5679,
		5963, 6240, 5958, 5915, 5737, 6000, 5747, 5529, 5562, 5909, 5713, 5680, 5729, 5656, 5820, 5670, 5884, 5686, 5662, 5848, 5710, 5707, 5821, 5564, 6029, 6045, 5765, 5727, 5653,
		5766, 5784, 5893, 5755, 5756, 5836, 5652, 5971, 6000, 5689, 6110, 5953, 6102, 5747, 5872, 5808, 5891, 5839, 5719, 5865, 6114, 5811, 5687, 5834, 5759, 5873, 6114, 6314, 5757,
		5849, 5901, 5993, 5226, 5338, 5356, 5299, 5426, 5479, 5687, 5594, 5497, 5735, 19094, 19217, 19264, 18976, 19040, 19348, 19092, 5777, 5810, 5636, 5600, 5681, 5528, 5573, 5494,
		5613, 5603, 5509, 5455, 5552, 5773, 5903, 5385,
	}
	v, e := shamosEstimator(data)
	require.NoError(t, e)
	assert.InDelta(t, 191.4089, v, 0.1)
}
