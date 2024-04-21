package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBinarySerializer(t *testing.T) {
	encoder := NewBinarySerializer()
	encoder.SerializeLengthUint15(0)          // [0]
	encoder.SerializeLengthUint15(127)        // [127]
	encoder.SerializeLengthUint15(128)        // [128, 128]
	encoder.SerializeLengthUint15(32767)      // [255, 255]
	encoder.SerializeLengthUint30(0)          // [0]
	encoder.SerializeLengthUint30(63)         // [63]
	encoder.SerializeLengthUint30(64)         // [64, 64]
	encoder.SerializeLengthUint30(16383)      // [127, 255]
	encoder.SerializeLengthUint30(16384)      // [128, 64, 0]
	encoder.SerializeLengthUint30(4194303)    // [191, 255, 255]
	encoder.SerializeLengthUint30(4194304)    // [192, 64, 0, 0]
	encoder.SerializeLengthUint30(1073741823) // [255, 255, 255, 255]

	assert.Equal(t, []byte{
		0,
		127,
		128,
		128,
		255,
		255,
		0,
		63,
		64,
		64,
		127,
		255,
		128,
		64,
		0,
		191,
		255,
		255,
		192,
		64,
		0,
		0,
		255,
		255,
		255,
		255,
	}, encoder.output)
}
