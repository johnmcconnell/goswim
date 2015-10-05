package goswim

import (
	"bytes"
	"fmt"
	"math/rand"
	"time"
	"sync"
)

// MemberList ...
type MemberList struct {
	Lock       *sync.Mutex
	Entries    map[uint64]*Message
	Expecting  *Message
	Suspicions map[uint64]time.Time
	Timeout    time.Duration
}

// String ...
func (m MemberList) String() string {
	var buffer bytes.Buffer

	buffer.WriteString("\n===== MEMBERSHIP =====\n")

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
	Entries := make(map[uint64]*Message)

	for _, M := range MS {
		Copy := M

		Entries[M.ID()] = &Copy
	}

	Suspicions := make(map[uint64]time.Time)

	M := MemberList{
		Entries:    Entries,
		Suspicions: Suspicions,
		Timeout:    T,
		Lock:       &sync.Mutex{},
	}

	return &M
}

// CheckSuspicionTimeouts ...
func (m *MemberList) CheckSuspicionTimeouts(Now time.Time) {
	m.Lock.Lock()
	defer m.Lock.Unlock()

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

	m.Lock.Lock()
	defer m.Lock.Unlock()

	// Missed the previous expecting
	if m.Expecting != nil {
		E := m.Entries[m.Expecting.ID()]
		if E.State == Alive {
			E.State = Suspected

			m.Suspicions[E.ID()] = Now
		}
	}

	m.Expecting = &M
}

// Received ...
func (m *MemberList) Received(M Message) {
	m.Lock.Lock()

	if m.Expecting != nil {
		if m.Expecting.IP == M.IP {
			m.Expecting = nil
		}
	}

	m.Lock.Unlock()

	m.Update(M, true)
}

// OutstandingAck ...
func (m *MemberList) OutstandingAck() *Message {
	Message := m.Expecting

	return Message
}

// Select ...
func (m *MemberList) Select(L int) []Message {
	m.Lock.Lock()
	defer m.Lock.Unlock()

	Entries := m.Entries

	LE := len(Entries)

	if LE < L {
		L = LE
	}

	if L == 0 {
		return nil
	}

	Selection := make([]Message, L)

	Perm := rand.Perm(
		LE,
	)

	ListEntries := make([]Message, LE)

	I := 0
	for _, E := range Entries {
		if I == LE {
			break
		}

		ListEntries[I] = *E
		I++
	}

	for I := range Selection {
		E := ListEntries[Perm[I]]
		Selection[I] = E
	}

	return Selection
}

// Updates ...
func (m *MemberList) Updates(MS ...Message) {
	for _, M := range MS {
		m.Update(M, false)
	}
}

// Update ...
func (m *MemberList) Update(M Message, base bool) bool {
	m.Lock.Lock()
	defer m.Lock.Unlock()

	e, ok := m.Entries[M.ID()]

	if !ok {
		m.Entries[M.ID()] = &M
		return true
	}

	if M.IncNumber == e.IncNumber {
		if M.State > e.State {
			m.Entries[M.ID()] = &M

			return true
		}
	}

	if M.IncNumber > e.IncNumber {
		m.Entries[M.ID()] = &M

		delete(m.Suspicions, M.ID())

		return true
	}

	if base {
		e.State = Alive
		delete(m.Suspicions, M.ID())
	}

	return false
}
