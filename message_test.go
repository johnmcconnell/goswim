package goswim

import (
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
)

func TestEncodeAndDecode(t *testing.T) {
	assert := assert.New(t)

	IP := uint32(127 << 24)
	IP += 10 << 16
	IP += 10 << 8
	IP += 2

	m := Message{
		State:     Suspected,
		IP:        IP,
		Port:      8888,
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

func TestURL(t *testing.T) {
	assert := assert.New(t)

	IP := uint32(127 << 24)
	IP += 10 << 16
	IP += 10 << 8
	IP += 2

	m := Message{
		State:     Suspected,
		IP:        IP,
		Port:      8888,
		IncNumber: 42,
	}

	assert.Equal(
		"127.10.10.2:8888",
		m.URL(),
		"The urls match",
	)
}

func TestBuild(t *testing.T) {
	assert := assert.New(t)

	IP := uint32(127 << 24)
	IP += 10 << 16
	IP += 10 << 8
	IP += 2

	M := Message{
		State:     Alive,
		IP:        IP,
		Port:      3000,
		IncNumber: 67,
	}

	URL := M.URL()

	assert.Equal(
		"127.10.10.2:3000",
		URL,
		"URL match",
	)

	ExpectedAddr, err := net.ResolveUDPAddr(
		"udp",
		URL,
	)

	assert.Nil(
		err,
		"no error",
	)

	Addr, err := M.UDPAddr()

	assert.Nil(
		err,
		"no error",
	)

	assert.Equal(
		*ExpectedAddr,
		*Addr,
	)

	NewM := BuildMessage(
		M.State,
		Addr,
		M.IncNumber,
	)

	assert.Equal(
		M,
		NewM,
		"they match",
	)
}

func TestEncode(t *testing.T) {
	assert := assert.New(t)

	IP := uint32(127 << 24)
	IP += 10 << 16
	IP += 10 << 8
	IP += 2

	m := Message{
		State:     Alive,
		IP:        IP,
		Port:      2,
		IncNumber: 67,
	}

	Encoded := EncodeMessage(
		m,
	)

	assert.Equal(
		16,
		len(Encoded),
		"Only takes 16 bytes",
	)
}
