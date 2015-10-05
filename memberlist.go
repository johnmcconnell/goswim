package goswim

import (
	"bytes"
	"fmt"
	"math/rand"
	"time"
)

// MemberList ...
type MemberList struct {
	Entries    map[uint32]*Message
	Expecting  *Message
	Suspicions map[uint32]time.Time
	Timeout    time.Duration
}

// String ...
func (m MemberList) String() string {
	var buffer bytes.Buffer

	buffer.WriteString("===== MEMBERSHIP =====\n")

	for _, e := range m.Entries {
		buffer.WriteString(
			e.String() + "\n",
		)
	}

	buffer.WriteString("===== SUSPICIONS =====\n")

	for id := range m.Suspicions {
		buffer.WriteString(
			fmt.Sprintf(
				"%v\n",
				id,
			),
		)
	}

	return buffer.String()
}

// NewMemberList ...
func NewMemberList(MS []Message, T time.Duration) *MemberList {
	Entries := make(map[uint32]*Message)

	for _, M := range MS {
		Copy := M

		Entries[M.IP] = &Copy
	}

	Suspicions := make(map[uint32]time.Time)

	M := MemberList{
		Entries:    Entries,
		Suspicions: Suspicions,
		Timeout:    T,
	}

	return &M
}

// CheckSuspicionTimeouts ...
func (m *MemberList) CheckSuspicionTimeouts(Now time.Time) {
	for ID, Time := range m.Suspicions {
		D := Now.Sub(Time)

		if D > m.Timeout {
			E := m.Entries[ID]
			E.State = Failed

			delete(m.Suspicions, ID)
		}
	}
}

// Awaiting ...
func (m *MemberList) Awaiting(M Message, Now time.Time) {
	m.CheckSuspicionTimeouts(Now)

	// Missed the previous expecting
	if m.Expecting != nil {
		E := m.Entries[m.Expecting.IP]
		if E.State == Alive {
			E.State = Suspected

			m.Suspicions[E.IP] = Now
		}
	}

	m.Expecting = &M
}

// Received ...
func (m *MemberList) Received(M Message) {
	fmt.Println("Received:", M)

	if m.Expecting != nil {
		if m.Expecting.IP == M.IP {
			m.Expecting = nil
		}
	}

	m.Update(M)
}

// OutstandingAck ...
func (m *MemberList) OutstandingAck() *Message {
	Message := m.Expecting

	return Message
}

// Select ...
func (m *MemberList) Select(L int) []Message {
	Selection := make([]Message, L)

	LE := len(m.Entries)

	Perm := rand.Perm(
		LE,
	)

	ListEntries := make([]Message, LE)

	I := 0
	for _, E := range m.Entries {
		ListEntries[I] = *E

		I++
	}

	for I, P := range Perm {
		Selection[I] = ListEntries[P]
	}

	return Selection
}

// Update ...
func (m *MemberList) Update(M Message) bool {
	e, ok := m.Entries[M.IP]

	if !ok {
		m.Entries[M.IP] = &M
		return true
	}

	if M.IncNumber == e.IncNumber {
		if M.State > e.State {
			m.Entries[M.IP] = &M

			return true
		}
	}

	if M.IncNumber > e.IncNumber {
		m.Entries[M.IP] = &M

		delete(m.Suspicions, M.IP)

		return true
	}

	if M.IncNumber == e.IncNumber {
		if M.State == Alive {
			delete(m.Suspicions, M.IP)
		}
	}

	return false
}
