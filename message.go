package goswim

import (
	"fmt"
)

// Membership States
const (
	Alive = iota
	Suspected
	Failed
)

// Message ...
type Message struct {
	State     uint32
	IP        uint32
	Port      uint32
	IncNumber uint32
}

// URL ...
func (m Message) URL() string {
	return fmt.Sprintf(
		"%v.%v.%v.%v:%v",
		m.IP>>24,
		(m.IP>>16)&0xFF,
		(m.IP>>8)&0xFF,
		m.IP&0xFF,
		m.Port,
	)
}

// Encoded ...
func (m Message) Encoded() []byte {
	return EncodeMessage(m)
}

// EncodeMessage ...
func EncodeMessage(m Message) []byte {
	var Encoded []byte

	const Mask = 0xFF

	State := m.State

	B := byte(State & Mask)
	Encoded = append(Encoded, B)

	State = State >> 8

	B = byte(State & Mask)
	Encoded = append(Encoded, B)

	State = State >> 8

	B = byte(State & Mask)
	Encoded = append(Encoded, B)

	State = State >> 8

	B = byte(State & Mask)
	Encoded = append(Encoded, B)

	IP := m.IP

	B = byte(IP & Mask)
	Encoded = append(Encoded, B)

	IP = IP >> 8

	B = byte(IP & Mask)
	Encoded = append(Encoded, B)

	IP = IP >> 8

	B = byte(IP & Mask)
	Encoded = append(Encoded, B)

	IP = IP >> 8

	B = byte(IP & Mask)
	Encoded = append(Encoded, B)

	Port := m.Port

	B = byte(Port & Mask)
	Encoded = append(Encoded, B)

	Port = Port >> 8

	B = byte(Port & Mask)
	Encoded = append(Encoded, B)

	Port = Port >> 8

	B = byte(Port & Mask)
	Encoded = append(Encoded, B)

	Port = Port >> 8

	B = byte(Port & Mask)
	Encoded = append(Encoded, B)

	IncNumber := m.IncNumber

	B = byte(IncNumber & Mask)
	Encoded = append(Encoded, B)

	IncNumber = IncNumber >> 8

	B = byte(IncNumber & Mask)
	Encoded = append(Encoded, B)

	IncNumber = IncNumber >> 8

	B = byte(IncNumber & Mask)
	Encoded = append(Encoded, B)

	IncNumber = IncNumber >> 8

	B = byte(IncNumber & Mask)
	Encoded = append(Encoded, B)

	return Encoded
}

// DecodeMessage ...
func DecodeMessage(BS []byte) Message {
	if len(BS) != 16 {
		panic(
			fmt.Sprintf(
				"Need 16 bytes to decode message, instead of %v",
				len(BS),
			),
		)
	}

	const Mask = 0xFF

	State := uint32(BS[0])

	for x := 1; x < 4; x++ {
		State += uint32(BS[x]) << uint32(8*x)
	}

	IP := uint32(BS[4])

	for x := 5; x < 8; x++ {
		IP += uint32(BS[x]) << uint32(8*(x-4))
	}

	Port := uint32(BS[8])

	for x := 9; x < 12; x++ {
		Port += uint32(BS[x]) << uint32(8*(x-8))
	}

	IncNumber := uint32(BS[12])

	for x := 13; x < 16; x++ {
		Port += uint32(BS[x]) << uint32(12*(x-12))
	}

	return Message{
		State:     State,
		IP:        IP,
		Port:      Port,
		IncNumber: IncNumber,
	}
}
