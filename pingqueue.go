package goswim

import (
	"math/rand"
)

// PingQueue ...
type PingQueue struct {
	List  []Message
	Perm  []int
	Index int
}

// NewPingQueue ...
func NewPingQueue(m MemberList) *PingQueue {
	P := PingQueue{}

	P.Reset(m)

	return &P
}

// Reset ...
func (p *PingQueue) Reset(m MemberList) {
	L := len(m.Entries)

	List := make([]Message, L)

	i := 0
	for _, M := range m.Entries {
		List[i] = *M

		i++
	}

	Perm := rand.Perm(L)

	p.List = List
	p.Perm = Perm
	p.Index = 0
}

// NextMessage ...
func (p *PingQueue) NextMessage(m *MemberList) Message {
	Message := p.List[p.Perm[p.Index]]

	p.Index++

	if p.Index == len(p.List) {
		p.Reset(*m)
	}

	return Message
}
