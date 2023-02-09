package main

import (
	"fmt"
	"math/rand"
	"sync/atomic"
	"time"

	"github.com/timandy/routine"
)

var th *ThreadId
var threadLocal = routine.NewThreadLocal()
var ge *GenerateHistory

// GenerateHistory相关
var sLock int32 = 0

type ThreadId struct {
	// Atomic integer containing the next thread ID to be assigned
	nextId int64
}

func initThreadId() {
	threadLocal.Set(atomic.AddInt64(&th.nextId, 1))
}
func (this *ThreadId) get() int64 {
	return threadLocal.Get().(int64)
}

type GenerateHistory struct {
	threadnum    int  //input
	testnum      int  //input
	isSequential bool //input
	msec         int
	nsec         int
	totalPc      int
	fin          []bool

	/****************Manually Set Testing Information **************/

	routenum   int // route is designed from 1 to 3
	coachnum   int // coach is arranged from 1 to 5
	seatnum    int // seat is allocated from 1 to 20
	stationnum int // station is designed from 1 to 5

	tds           *TicketingDs
	methodList    []string
	freqList      []int
	currentTicket []Ticket
	currentRes    []string
	soldTicket    [][]Ticket
	initLock      int32 // use as atomic bool
	r             rand.Rand
	//	final static AtomicInteger tidGen = new AtomicInteger(0);
	// final static Random rand = new Random();
}

func initGenerateHistory() {
	ge.routenum = 3
	ge.coachnum = 3
	ge.seatnum = 3
	ge.stationnum = 3
	ge.initLock = 0
	ge.r = *rand.New(rand.NewSource(time.Now().UnixNano()))
	ge.initialization()
}

func (generateHistory *GenerateHistory) initialization() {
	generateHistory.tds = newTicketingDs(generateHistory.routenum, generateHistory.coachnum, generateHistory.seatnum, generateHistory.stationnum, generateHistory.threadnum)
	for i := 0; i < generateHistory.threadnum; i++ {
		threadTickets := make([]Ticket, maxThreadNum)
		generateHistory.soldTicket = append(generateHistory.soldTicket, threadTickets)
		generateHistory.currentTicket = make([]Ticket, maxThreadNum)
		generateHistory.currentRes = append(generateHistory.currentRes, " ")
	}
	generateHistory.methodList = append(generateHistory.methodList, "refundTicket")
	generateHistory.freqList = append(generateHistory.freqList, 10)
	generateHistory.methodList = append(generateHistory.methodList, "buyTicket")
	generateHistory.freqList = append(generateHistory.freqList, 30)
	generateHistory.methodList = append(generateHistory.methodList, "inquiry")
	generateHistory.freqList = append(generateHistory.freqList, 60)
	generateHistory.totalPc = 100

}

func (this *GenerateHistory) exOthNotFin(tNum, tid int) bool {
	flag := false
	for k := 0; k < tNum; k++ {
		if k == tid {
			continue
		}
		flag = (flag || !(this.fin[k]))
	}
	return flag
}

func (*GenerateHistory) SLOCK_TAKE() {
	for !atomic.CompareAndSwapInt32(&sLock, 0, 1) {

	}
}

func (*GenerateHistory) SLOCK_GIVE() {
	atomic.StoreInt32(&sLock, 0)
}

func (*GenerateHistory) SLOCK_TRY() bool {
	return atomic.LoadInt32(&sLock) == 0
}

func (generateHistory *GenerateHistory) getPassengerName() string {
	uid := generateHistory.r.Int()
	return fmt.Sprintln("passenger", uid)
}

func (generateHistory *GenerateHistory) print(preTime, postTime int64, actionName string) {
	ticket := generateHistory.currentTicket[th.get()]
	fmt.Println(preTime, " ", postTime, " ", th.get(), " ", actionName, " ", ticket.tid, " ", ticket.passenger, " ", ticket.route, " ", ticket.coach, " ", ticket.departure, " ", ticket.arrival, " ", ticket.seat, " ", generateHistory.currentTicket[th.get()])
}

func (generateHistory *GenerateHistory) execute(num int) bool {
	var route, departure, arrival int
	ticket = Ticket{}
	switch(num){
	  case 0:  //refund
		if len(soldTicket[th.get()]) == 0 {
			return false
		}
		n := rand.nextInt(soldTicket.get(ThreadId.get()).size())
		ticket = soldTicket.get(ThreadId.get()).remove(n);
		if(ticket == null){
		  return false
		}
		currentTicket.set(ThreadId.get(), ticket);
	    flag := tds.refundTicket(ticket);
		currentRes.set(ThreadId.get(), "true"); 
		return flag
	  case 1:  //buy
		passenger := getPassengerName();
		route = gegenerateHistory.r.Int(routenum) + 1;
		departure = rand.nextInt(stationnum - 1) + 1;
		arrival = departure + rand.nextInt(stationnum - departure) + 1;
		ticket = tds.buyTicket(passenger, route, departure, arrival);
		if(ticket == null){
		  ticket = new Ticket();
		  ticket.passenger = passenger;
		  ticket.route = route;
		  ticket.departure = departure;
		  ticket.arrival = arrival;
		  ticket.seat = 0;
		  currentTicket.set(ThreadId.get(), ticket);
		  currentRes.set(ThreadId.get(), "false");
		  return true;
		}
		currentTicket.set(ThreadId.get(), ticket);
		currentRes.set(ThreadId.get(), "true");
		soldTicket.get(ThreadId.get()).add(ticket);
		return true;
	  case 2: 
		ticket.passenger = getPassengerName();
		ticket.route = rand.nextInt(routenum) + 1;
		ticket.departure = rand.nextInt(stationnum - 1) + 1;
		ticket.arrival = ticket.departure + rand.nextInt(stationnum - ticket.departure) + 1; // arrival is always greater than departure
		ticket.seat = tds.inquiry(ticket.route, ticket.departure, ticket.arrival);
		currentTicket.set(ThreadId.get(), ticket);
		currentRes.set(ThreadId.get(), "true"); 
		return true
	  default: 
		System.out.println("Error in execution.");
		return false
	  

	}
}