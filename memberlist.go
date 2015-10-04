package goswim

// MemberList ...
type MemberList struct {
	Entries map[uint32]Message
}

// NewMemberList ...
func NewMemberList(MS []Message) *MemberList {
	Entries := make(map[uint32]Message)

	for _, M := range MS {
		Entries[M.ID] = M
	}

	M := MemberList{
		Entries: Entries,
	}

	return &M
}

// Update ...
func (m *MemberList) Update(M Message) bool {
	e, ok := m.Entries[M.ID]

	if !ok {
		m.Entries[M.ID] = M
		return true
	}

	if M.IncNumber > e.IncNumber {
		if M.State > e.State {
			m.Entries[M.ID] = M

			return true
		}
	}

	return false
}
