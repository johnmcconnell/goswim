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
	var List []Message

	for _, M := range m.Entries {
		if M.State != Failed {
			List = append(
				List,
				*M,
			)
		}
	}

	L := len(List)
	Perm := rand.Perm(L)

	p.List = List
	p.Perm = Perm
	p.Index = 0
}

// NextMessage ...
func (p *PingQueue) NextMessage(m *MemberList) *Message {
	if len(p.List) == 0 {
		return nil
	}

	Message := p.List[p.Perm[p.Index]]

	p.Index++

	if p.Index == len(p.List) {
		p.Reset(*m)
	}

	return &Message
}
