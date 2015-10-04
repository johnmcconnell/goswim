package goswim

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEncodeAndDecode(t *testing.T) {
	assert := assert.New(t)

	ID := uint32(127 << 24)
	ID += 10 << 16
	ID += 10 << 8
	ID += 2

	m := Message{
		State:     Suspected,
		ID:        ID,
		IncNumber: 42,
	}

	r := DecodeMessage(
		EncodeMessage(
			m,
		),
	)

	assert.Equal(
		m,
		r,
		"The two messages match",
	)
}

func TestEncode(t *testing.T) {
	assert := assert.New(t)

	ID := uint32(127 << 24)
	ID += 10 << 16
	ID += 10 << 8
	ID += 2

	m := Message{
		State:     Alive,
		ID:        ID,
		IncNumber: 67,
	}

	Encoded := EncodeMessage(
		m,
	)

	assert.Equal(
		12,
		len(Encoded),
		"Only takes 12 bytes",
	)
}
