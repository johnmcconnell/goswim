package goswim

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var (
	timeout = 1 * time.Millisecond
)

func TestSWIMSuccess(t *testing.T) {
	assert := assert.New(t)

	MS := []Message{
		Message{
			State:     Alive,
			IP:        1,
			Port:      3000,
			IncNumber: 42,
		},
		Message{
			State:     Alive,
			IP:        2,
			Port:      3000,
			IncNumber: 1,
		},
	}

	List := NewMemberList(
		MS,
		timeout,
	)

	Now := time.Now()

	List.Awaiting(
		MS[0],
		Now,
	)

	assert.Equal(
		MS[0],
		*List.OutstandingAck(),
		"Haven't received Ack",
	)

	List.Awaiting(
		MS[0],
		Now,
	)

	List.Received(
		MS[0],
	)

	assert.Nil(
		List.OutstandingAck(),
		"Nothing outstanding",
	)
}

func TestSelect(t *testing.T) {
	assert := assert.New(t)

	MS := []Message{
		Message{
			State:     Alive,
			IP:        1,
			Port:      3000,
			IncNumber: 42,
		},
		Message{
			State:     Alive,
			IP:        2,
			Port:      3000,
			IncNumber: 1,
		},
	}

	List := NewMemberList(
		MS,
		timeout,
	)

	RandMS := List.Select(2)

	assert.Equal(
		2,
		len(RandMS),
	)

	assert.NotEqual(
		RandMS[0],
		RandMS[1],
		fmt.Sprintf(
			"%v != %v",
			RandMS[0],
			RandMS[1],
		),
	)
}

func TestSWIMSuspicion(t *testing.T) {
	assert := assert.New(t)

	MS := []Message{
		Message{
			State:     Alive,
			IP:        1,
			Port:      3000,
			IncNumber: 42,
		},
		Message{
			State:     Alive,
			IP:        2,
			Port:      3000,
			IncNumber: 1,
		},
	}

	List := NewMemberList(
		MS,
		timeout,
	)

	Now := time.Now()

	List.Awaiting(
		MS[0],
		Now,
	)

	Now = Now.Add(
		10 * time.Millisecond,
	)

	List.Awaiting(
		MS[1],
		Now,
	)

	assert.EqualValues(
		Suspected,
		List.Entries[MS[0].ID()].State,
		"Now suspected",
	)

	List.Received(
		MS[1],
	)

	Now = Now.Add(
		timeout + 1,
	)

	List.Awaiting(
		MS[0],
		Now,
	)

	assert.EqualValues(
		Failed,
		List.Entries[MS[0].ID()].State,
		"Now failed",
	)

	List.Received(
		MS[0],
	)

	assert.EqualValues(
		Failed,
		List.Entries[MS[0].ID()].State,
		"Failed because no update to incarnation number",
	)

	RejoinM := MS[0]
	RejoinM.IncNumber = MS[0].IncNumber + 1

	List.Received(
		RejoinM,
	)

	assert.EqualValues(
		Alive,
		List.Entries[MS[0].ID()].State,
		"Now alive",
	)
}
