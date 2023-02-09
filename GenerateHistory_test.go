package main

import (
	"math/rand"
	"testing"
	"time"
)

func TestRand(t *testing.T) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	a := r.Intn(1000)
	b := r.Intn(1000)
	if a == b {
		t.Errorf("no rand")
	}
}

func TestIsTicketNil(t *testing.T) {

	var x Ticket
	if !isTicketNil(x) {
		t.Errorf("should be nil when not initializtion")
	}
	x = Ticket{}
	if !isTicketNil(x) {
		t.Errorf("should be nil when not initializtion")
	}
	ticketList := make([]Ticket, 10)
	x = ticketList[0]
	if !isTicketNil(x) {
		t.Errorf("should be nil when not initializtion")
	}
	x.arrival = 1
	if isTicketNil(x) {
		t.Errorf("should not be nil when not initializtion")
	}
}
