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
	ID        uint32
	IncNumber uint32
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

	ID := m.ID

	B = byte(ID & Mask)
	Encoded = append(Encoded, B)

	ID = ID >> 8

	B = byte(ID & Mask)
	Encoded = append(Encoded, B)

	ID = ID >> 8

	B = byte(ID & Mask)
	Encoded = append(Encoded, B)

	ID = ID >> 8

	B = byte(ID & Mask)
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
	if len(BS) != 12 {
		panic(
			fmt.Sprintf(
				"Need 12 bytes to decode message, instead of %v",
				len(BS),
			),
		)
	}

	const Mask = 0xFF

	State := uint32(BS[0])

	for x := 1; x < 4; x++ {
		State += uint32(BS[x]) << uint32(8*x)
	}

	ID := uint32(BS[4])

	for x := 5; x < 8; x++ {
		ID += uint32(BS[x]) << uint32(8*(x-4))
	}

	IncNumber := uint32(BS[8])

	for x := 9; x < 12; x++ {
		IncNumber += uint32(BS[x]) << uint32(8*(x-8))
	}

	return Message{
		State:     State,
		ID:        ID,
		IncNumber: IncNumber,
	}
}
