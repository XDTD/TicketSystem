package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

/**********Manually Modified ***********/
var isPosttime = true
var detail = false

const (
	routenum     int = 3
	coachnum     int = 3
	seatnum      int = 3
	stationnum   int = 3
	debugMode    int = 1
	maxThreadNum int = 64
)

type HistoryLine struct {
	pretime       int64
	posttime      int64
	threadid      int
	operationName string
	tid           int64
	passenger     string
	route         int
	coach         int
	seat          int
	departure     int
	arrival       int
	res           string
}

type Replay struct {
	threadNum  int
	methodList []string
	history    []HistoryLine
	object     *TicketingDs
}

func parseline(historyList *[]HistoryLine, line string) bool {
	tl := HistoryLine{}
	fmt.Sscanf(line, "%v %v %v %v %v %v %v %v %v %v %v", &tl.pretime, &tl.posttime, &tl.threadid, &tl.operationName, &tl.tid, &tl.passenger, &tl.route, &tl.coach, &tl.departure, &tl.arrival, &tl.seat)
	*historyList = append(*historyList, tl)
	return true
}

func (this *Replay) initialization() {
	this.object = newTicketingDs(routenum, coachnum, seatnum, stationnum, this.threadNum)
	this.methodList = append(this.methodList, "refundTicket")
	this.methodList = append(this.methodList, "buyTicket")
	this.methodList = append(this.methodList, "inquiry")
}

func (this *Replay) execute(methodName string, line HistoryLine, line_num int) bool {
	ticket := Ticket{}
	flag := false
	ticket.tid = line.tid
	ticket.passenger = line.passenger
	ticket.route = line.route
	ticket.coach = line.coach
	ticket.departure = line.departure
	ticket.arrival = line.arrival
	ticket.seat = line.seat
	if methodName == "buyTicket" {
		if line.res == "false" {
			num := this.object.Inquiry(ticket.route, ticket.departure, ticket.arrival)
			if num == 0 {
				return true
			} else {
				fmt.Println("Error: TicketSoldOut", " ", line.pretime, " ", line.posttime, " ", line.threadid, " ", line.route, " ", line.departure, " ", line.arrival)
				fmt.Println("RemainTicket", " ", num, " ", line.route, " ", line.departure, " ", line.arrival)
				return false
			}
		}
		ticket1 := this.object.BuyTicket(ticket.passenger, ticket.route, ticket.departure, ticket.arrival)
		if ticket1 != nil && line.res == "true" &&
			ticket.passenger == ticket1.passenger && ticket.route == ticket1.route &&
			ticket.coach == ticket1.coach && ticket.departure == ticket1.departure &&
			ticket.arrival == ticket1.arrival && ticket.seat == ticket1.seat {
			return true
		} else {
			fmt.Println("Error: Ticket is bought", " ", line.pretime, " ", line.posttime, " ", line.threadid, " ", ticket.tid, " ", ticket.passenger, " ", ticket.route, " ", ticket.coach, " ", ticket.departure, " ", ticket.arrival, " ", ticket.seat)
			return false
		}
	} else if methodName == "refundTicket" {
		flag = this.object.RefundTicket(ticket)
		if (flag && line.res == "true") || (!flag && line.res == "false") {
			return true
		} else {
			fmt.Println("Error: Ticket is refunded", " ", line.pretime, " ", line.posttime, " ", line.threadid, " ", ticket.tid, " ", ticket.passenger, " ", ticket.route, " ", ticket.coach, " ", ticket.departure, " ", ticket.arrival, " ", ticket.seat)
			return false
		}
	} else if methodName == "inquiry" {
		num := this.object.Inquiry(line.route, line.departure, line.arrival)
		if num == line.seat {
			return true
		} else {
			fmt.Println("Error: RemainTicket", " ", line.pretime, " ", line.posttime, " ", line.threadid, " ", line.route, " ", line.departure, " ", line.arrival, " ", line.seat)
			fmt.Println("Real RemainTicket is", " ", line.seat, " ", ", Expect RemainTicket is", " ", num, ", ", line.route, " ", line.departure, " ", line.arrival)
			return false
		}
	}
	fmt.Println("No match method name")
	return false
}

/***********************VeriLin*************** */

func writeHistoryToFile(historyList []HistoryLine, filename string) {
	f, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	writeHistory := func(historyList []HistoryLine) {
		for i := range historyList {
			tl := historyList[i]
			line := fmt.Sprintln(tl.pretime, " ", tl.posttime, " ", tl.threadid, " ", tl.operationName, " ", tl.tid, " ", tl.passenger, " ", tl.route, " ", tl.coach, " ", tl.departure, " ", tl.arrival, " ", tl.seat)
			_, err := f.WriteString(line + "\n")
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	writeHistory(historyList)

}

func readHistory(historyList []HistoryLine, filename string) bool {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	i := 0
	for scanner.Scan() {
		if !parseline(&historyList, scanner.Text()) {
			log.Fatal("Error in parsing line ", i)
			return false
		}
		i++
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
		return false
	}

	return true
}

func (this *Replay) checkline(historyList []HistoryLine, index int) bool {
	line := historyList[index]

	if debugMode == 1 {
		if index == 158 {
			fmt.Println("Debugging line ", index, " ")
		}
	}

	for i := range this.methodList {
		if line.operationName == this.methodList[i] {
			flag := this.execute(this.methodList[i], line, index)
			fmt.Println("Line ", index, " executing ", this.methodList[i], " res: ", flag, " tid = ", line.tid)
			return flag
		}
	}

	return false

}
