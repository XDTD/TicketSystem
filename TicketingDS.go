package main

import (
	"sync"
	"sync/atomic"
)

type TicketingDs struct {
	routenum   int //车次
	coachnum   int //车厢数
	seatnum    int //座位数
	stationnum int //车站数
	threadnum  int //线程数

	tickets    [][][]int64
	maxTid     int64
	ticketCopy sync.Map //保证tid不重复
}

func newTicketingDs(routenum, coachnum, seatnum, stationnum, threadnum int) *TicketingDs {
	tickets := make([][][]int64, routenum)
	for i := range tickets {
		tickets[i] = make([][]int64, coachnum)
		for j := range tickets[i] {
			tickets[i][j] = make([]int64, seatnum)
		}
	}
	//用于生成tid
	var ticketCopy sync.Map
	return &TicketingDs{
		routenum:   routenum,
		coachnum:   coachnum,
		seatnum:    seatnum,
		stationnum: stationnum,
		threadnum:  threadnum,
		tickets:    tickets,
		maxTid:     1,
		ticketCopy: ticketCopy,
	}
}

func newDefualtTicketingDs() *TicketingDs {
	var routenum int = 5
	var coachnum int = 8
	var seatnum int = 100
	var stationnum int = 10
	var threadnum int = 16
	return newTicketingDs(routenum, coachnum, seatnum, stationnum, threadnum)
}

// 功能：买票，买指定车次从departure到arrival的票，如果有就返回买的票，如果没有就返回Null
// 输入： route 指定车次 ； departure 出发站; arrival到达站
// 输出： t  买的票
func (this *TicketingDs) BuyTicket(passenger string, route, departure, arrival int) *Ticket {
	//车票从1开始计数，需要减一
	route--
	departure--
	arrival--
	//得到座位的位图
	var seatBit int64 = (0x01 << (arrival - departure)) - 1
	seatBit = seatBit << departure
	for i := 0; i < this.coachnum; i++ {

		for j := 0; j < this.seatnum; j++ {
			//检测座位是否有票
			var oldTickets int64 = this.tickets[route][i][j]
			var result int64 = seatBit & oldTickets
			if result == 0 {
				var newTickets = oldTickets | seatBit
				for ok := atomic.CompareAndSwapInt64(&this.tickets[route][i][j], oldTickets, newTickets); !ok; {
					oldTickets = this.tickets[route][i][j]
					result = seatBit & oldTickets
					if result != 0 {
						break
					}
					newTickets = oldTickets | seatBit
				}
			}

			//买到票了
			if result == 0 {
				// 车票生成
				t := Ticket{}
				t.tid = atomic.AddInt64(&this.maxTid, 1)
				t.passenger = passenger
				t.route = route + 1
				t.coach = i + 1
				t.seat = j + 1
				t.departure = departure + 1
				t.arrival = arrival + 1
				// 存储票
				this.ticketCopy.Store(t.tid, t)
				return &t
			}

		}

	}
	//遍历完没买到票就是null
	return nil
}

// 功能：查票，返回指定车次从departure到arrival的余票
// 输入： route 指定车次 ； departure 出发站; arrival到达站
// 输出： remainTickets 余票数目
func (this *TicketingDs) Inquiry(route, departure, arrival int) int {
	//从1开始计数的
	route--
	departure--
	arrival--

	var remainTickets int = 0
	//得到座位位图
	var seatBit int64 = (0x01 << (arrival - departure)) - 1
	seatBit = seatBit << departure
	for i := 0; i < this.coachnum; i++ {
		for j := 0; j < this.seatnum; j++ {
			//检测座位是否有票
			var result int = int(seatBit & this.tickets[route][i][j])
			//有票了
			if result == 0 {
				remainTickets++
			}
		}
	}
	return remainTickets
}

// 功能：退票，返回退票结果，真票是true，假票是false
// 输入： ticket  Ticket 要退的票
// 输出： result boolean 退票的结果
func (this *TicketingDs) RefundTicket(ticket Ticket) bool {
	//车票是从1开始计数的
	var route int = ticket.route - 1
	var coach int = ticket.coach - 1
	var seat int = ticket.seat - 1
	var arrival int = ticket.arrival - 1
	var departure int = ticket.departure - 1

	//检测是不是假票
	var realTicket Ticket
	realTicketInterface, ok := this.ticketCopy.Load(ticket.tid)
	if ok {
		realTicket, ok = (realTicketInterface).(Ticket)
		if ok && realTicket.arrival == ticket.arrival &&
			realTicket.departure == ticket.departure &&
			realTicket.route == ticket.route &&
			realTicket.coach == ticket.coach &&
			realTicket.seat == ticket.seat {
			//尝试删除票，删除成功开始退票
			this.ticketCopy.Delete(ticket.tid)
			//得到座位的位图
			var seatBit int64 = (0x01 << (arrival - departure)) - 1
			seatBit = seatBit << departure
			seatBit = ^seatBit

			// 修改座位信息
			var oldTickets int64 = this.tickets[route][coach][seat]
			var newTickets int64 = oldTickets & seatBit
			for ok := atomic.CompareAndSwapInt64(&this.tickets[route][coach][seat], oldTickets, newTickets); !ok; {
				oldTickets = this.tickets[route][coach][seat]
				newTickets = oldTickets & seatBit
			}
			return true

		}
	}

	return false
}
