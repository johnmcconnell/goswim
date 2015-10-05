package goswim

import (
	"bytes"
	"fmt"
	"sort"
	"time"
)

// InfQueueEntries ...
type InfQueueEntries []InfQueueEntry

// InfQueue ...
type InfQueue struct {
	Entries map[uint64]*InfQueueEntry
	Size    int
	Decay   int
}

// InfQueueEntry ...
type InfQueueEntry struct {
	Message
	Updated time.Time
	Sent    int
}

// String ...
func (e InfQueueEntry) String() string {
	return fmt.Sprintf(
		"%v --> %v:%v",
		e.Message,
		e.Updated,
		e.Sent,
	)
}

// String ...
func (m InfQueue) String() string {
	var buffer bytes.Buffer

	buffer.WriteString("\n===== Messages =====\n")

	for _, e := range m.Entries {
		buffer.WriteString(
			e.String() + "\n",
		)
	}

	return buffer.String()
}

// List ...
func (m InfQueue) List() InfQueueEntries {
	L := len(m.Entries)

	Entries := make(InfQueueEntries, L)

	I := 0
	for _, E := range m.Entries {
		Entries[I] = *E
		I++
	}

	Entries.Sort()

	return Entries
}

// NewInfQueue ...
func NewInfQueue(MS []Message, size, decay int, now time.Time) *InfQueue {
	Entries := make(map[uint64]*InfQueueEntry)

	for _, M := range MS {

		Entries[M.ID()] = &InfQueueEntry{
			Message: M,
			Updated: now,
			Sent:    0,
		}
	}

	M := InfQueue{
		Entries: Entries,
		Size:    size,
		Decay:   decay,
	}

	return &M
}

// Messages ...
func (m *InfQueue) Messages() []Message {
	L := len(m.Entries)
	S := L - m.Size

	if S < 0 {
		S = 0
	}

	Entries := m.List()

	var Messages []Message

	for i := L; i > S; i-- {
		E := Entries[i-1]

		E.Sent++

		if E.Sent < m.Decay {
			Messages = append(
				Messages,
				E.Message,
			)
		} else {
			delete(m.Entries, E.Message.ID())
		}
	}

	return Messages
}

// Update ...
func (m *InfQueue) Update(MS []Message, Now time.Time) {
	for _, M := range MS {
		Orig, ok := m.Entries[M.ID()]

		if !ok {
			E := InfQueueEntry{
				Message: M,
				Updated: Now,
				Sent:    0,
			}

			m.Entries[M.ID()] = &E

			continue
		}

		AllowUpdate := Orig.Message.IncNumber > M.IncNumber

		if Orig.Message.IncNumber == M.IncNumber {
			AllowUpdate = M.State > M.State
		}

		if AllowUpdate {
			m.Entries[M.ID()].Updated = Now
			m.Entries[M.ID()].Sent = 0
			m.Entries[M.ID()].Message = M
		}
	}
}

// Less ...
func (E InfQueueEntries) Less(i, j int) bool {
	I := E[i]
	J := E[j]
	return I.Updated.Before(J.Updated)
}

// Swap ...
func (E *InfQueueEntries) Swap(i, j int) {
	(*E)[i], (*E)[j] = (*E)[j], (*E)[i]
}

// Len ...
func (E InfQueueEntries) Len() int {
	return len(E)
}

// Del ...
func (E *InfQueueEntries) Del(i int) {
	(*E) = append((*E)[:i], (*E)[i+1:]...)
}

// Sort ...
func (E *InfQueueEntries) Sort() {
	sort.Sort(E)
}
