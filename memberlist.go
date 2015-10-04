package goswim

import (
	"bytes"
	"fmt"
)

// MemberList ...
type MemberList struct {
	Entries   map[uint32]Message
	Expecting *Message
}

// String ...
func (m MemberList) String() string {
	var buffer bytes.Buffer

	for _, e := range m.Entries {
		buffer.WriteString(
			fmt.Sprintf(
				"URL[%v] Status[%v] Inc[%v]\n",
				e.URL(),
				e.State,
				e.IncNumber,
			),
		)
	}

	return buffer.String()
}

// NewMemberList ...
func NewMemberList(MS []Message) *MemberList {
	Entries := make(map[uint32]Message)

	for _, M := range MS {
		Entries[M.IP] = M
	}

	M := MemberList{
		Entries: Entries,
	}

	return &M
}

// Awaiting ...
func (m *MemberList) Awaiting(M Message) {
	m.Expecting = &M
}

// Received ...
func (m *MemberList) Received(M Message) {
	if m.Expecting != nil {
		if m.Expecting.IP == M.IP {
			m.Expecting = nil
		}
	}
}

// OutstandingAck ...
func (m *MemberList) OutstandingAck() *Message {
	Message := m.Expecting

	m.Expecting = nil

	return Message
}

// Select ...
func (m *MemberList) Select(L int) []Message {
	return nil
}

// Update ...
func (m *MemberList) Update(M Message) bool {
	e, ok := m.Entries[M.IP]

	if !ok {
		m.Entries[M.IP] = M
		return true
	}

	if M.IncNumber > e.IncNumber {
		if M.State > e.State {
			m.Entries[M.IP] = M

			return true
		}
	}

	return false
}
