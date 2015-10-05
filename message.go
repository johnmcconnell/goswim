package goswim

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

// Membership States
const (
	Alive = iota
	Suspected
	Failed
	MessageSize = 16
)

// Message ...
type Message struct {
	State     uint32
	IP        uint32
	Port      uint32
	IncNumber uint32
}

// String ...
func (m Message) String() string {
	S := ""

	switch m.State {
	case Alive:
		S = "Alive"
	case Suspected:
		S = "Suspected"
	case Failed:
		S = "Failed"
	}

	return fmt.Sprintf(
		"URL[%v] Status[%v] Inc[%v]",
		m.URL(),
		S,
		m.IncNumber,
	)
}

// BuildMessageS ...
func BuildMessageS(s, h, p, i string) (Message, error) {
	Zero := Message{}

	var S uint32

	switch s {
	case "Alive":
		S = Alive
	case "Suspected":
		S = Suspected
	case "Failed":
		S = Failed
	}

	I, err := strconv.ParseUint(i, 10, 32)

	if err != nil {
		return Zero, err
	}

	URL := h + ":" + p

	Addr, err := net.ResolveUDPAddr(
		"udp",
		URL,
	)

	if err != nil {
		return Zero, err
	}

	return BuildMessage(S, Addr, uint32(I)), nil
}

// BuildMessage ...
func BuildMessage(s uint32, addr *net.UDPAddr, i uint32) Message {
	IPv4 := addr.IP.To4()

	var IP uint32

	IP += uint32(IPv4[0])
	IP = IP << 8
	IP += uint32(IPv4[1])
	IP = IP << 8
	IP += uint32(IPv4[2])
	IP = IP << 8
	IP += uint32(IPv4[3])

	m := Message{
		State:     s,
		IP:        IP,
		Port:      uint32(addr.Port),
		IncNumber: i,
	}

	return m
}

// BuildMessage2 ...
func BuildMessage2(s uint32, host, port string, i uint32) (Message, error) {
	Zero := Message{}

	Port, err := strconv.ParseUint(port, 10, 32)

	if err != nil {
		return Zero, err
	}

	Subs := strings.Split(host, ".")
	if len(Subs) != 4 {
		return Zero, fmt.Errorf(
			"Invalid length host[%v]",
			host,
		)
	}

	Part, err := strconv.Atoi(Subs[0])
	if err != nil {
		return Zero, err
	}
	IP := uint32(Part) << 24

	Part, err = strconv.Atoi(Subs[1])
	if err != nil {
		return Zero, err
	}
	IP += uint32(Part) << 16

	Part, err = strconv.Atoi(Subs[2])
	if err != nil {
		return Zero, err
	}
	IP += uint32(Part) << 8

	Part, err = strconv.Atoi(Subs[3])
	if err != nil {
		return Zero, err
	}
	IP += uint32(Part)

	m := Message{
		State:     s,
		IP:        IP,
		Port:      uint32(Port),
		IncNumber: i,
	}

	return m, nil
}

// URL ...
func (m Message) URL() string {
	return fmt.Sprintf(
		"%v:%v",
		m.IPv4S(),
		m.Port,
	)
}

// PortS ...
func (m Message) PortS() string {
	return fmt.Sprintf(
		"%v",
		m.Port,
	)
}

// ID ...
func (m Message) ID() uint64 {
	ID := uint64(m.IP) << 32
	ID += uint64(m.Port)

	return ID
}

// IPv4B ....
func (m Message) IPv4B() (byte, byte, byte, byte) {
	B1 := byte(m.IP >> 24)
	B2 := byte((m.IP >> 16) & 0xFF)
	B3 := byte((m.IP >> 8) & 0xFF)
	B4 := byte(m.IP & 0xFF)

	return B1, B2, B3, B4
}

// IPv4S ....
func (m Message) IPv4S() string {
	B1, B2, B3, B4 := m.IPv4B()

	return fmt.Sprintf(
		"%v.%v.%v.%v",
		B1,
		B2,
		B3,
		B4,
	)
}

// UDPAddr ...
func (m Message) UDPAddr() (*net.UDPAddr, error) {
	addr, err := net.ResolveUDPAddr(
		"udp",
		m.URL(),
	)

	return addr, err
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
	if len(BS) < MessageSize {
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
