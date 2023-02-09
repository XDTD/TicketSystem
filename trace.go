package main

import (
	"fmt"
	"math/rand"
)

type Trace struct {
	threadnum  int
	routenum   int // route is designed from 1 to 3
	coachnum   int // coach is arranged from 1 to 5
	seatnum    int // seat is allocated from 1 to 20
	stationnum int // station is designed from 1 to 5

	testnum int
	retpc   int // return ticket operation is 10% percent
	buypc   int // buy ticket operation is 30% percent
	inqpc   int //inquiry ticket operation is 60% percent
}

func (trace *Trace) initTrace() {
	trace.threadnum = 4
	trace.routenum = 3
	trace.coachnum = 5
	trace.seatnum = 10
	trace.stationnum = 8

	trace.testnum = 1000
	trace.retpc = 30
	trace.buypc = 60
	trace.inqpc = 100

}

func (trace *Trace) passengerName() string {
	uid := rand.Intn(trace.testnum)
	return fmt.Sprintln("passenger", uid)
}
