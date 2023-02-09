package main

import (
	"fmt"
	"reflect"
	"testing"
)

func TestParseline(t *testing.T) {
	tls := make([]HistoryLine, 0)
	tl := HistoryLine{
		pretime:       1,
		posttime:      2,
		threadid:      3,
		operationName: "4",
		tid:           5,
		passenger:     "6",
		route:         7,
		coach:         8,
		seat:          9,
		departure:     10,
		arrival:       11,
		// res:           "12",
	}

	var line string = fmt.Sprintln(tl.pretime, " ", tl.posttime, " ", tl.threadid, " ", tl.operationName, " ", tl.tid, " ", tl.passenger, " ", tl.route, " ", tl.coach, " ", tl.departure, " ", tl.arrival, " ", tl.seat)

	if ok := parseline(&tls, line); !ok {
		t.Error("parseline wrong")
	}

	if len(tls) == 0 {
		t.Error("parseline null result")
	} else {
		tlParse := tls[0]
		if !reflect.DeepEqual(tl, tlParse) {
			t.Errorf("parseline not equal")
		}
	}

}
