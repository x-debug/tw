package tw

import (
	"container/list"
	"errors"
	"log"
	"time"
)

var (
	// ErrTimerKeyFound returns when the key of the timer is found
	ErrTimerKeyFound    = errors.New("timer's key found")

	// ErrTimerKeyNotFound returns when the key of the timer is not found
	ErrTimerKeyNotFound = errors.New("timer's key is not exist")
)

type task struct {
	key    string
	circle int
	pos    int
	fn     func()
	slot   *list.List
}

type setingTasker struct {
	key      string
	duration time.Duration
	callback func()
}

type removingTasker struct {
	key string
}

//TimeWheel is an implementation of Simple Timing Wheels.
type TimeWheel struct {
	tickerPos  int
	interval   int
	sizeOfSlot int
	wheel      []*list.List
	wMap       map[string]*list.Element
	setChan    chan *setingTasker
	removeChan chan *removingTasker
	ticker     *time.Ticker
}

//NewTimeWheel construct time wheel,interval is the minimum span for the ticker to travel,
//sizeOfSlot is the total num of wheel slott
func NewTimeWheel(interval int, sizeOfSlot int) *TimeWheel {
	wheel := &TimeWheel{
		interval:   interval,
		tickerPos:  0,
		sizeOfSlot: sizeOfSlot,
		ticker:     time.NewTicker(time.Duration(interval) * time.Second),
		setChan:    make(chan *setingTasker),
		removeChan: make(chan *removingTasker),
		wMap:       make(map[string]*list.Element)}
	wheel.initWheel(wheel.sizeOfSlot)
	go wheel.run()
	return wheel
}

func (tw *TimeWheel) initWheel(slotNumber int) {
	tw.wheel = make([]*list.List, tw.sizeOfSlot)
	for i := 0; i < tw.sizeOfSlot; i++ {
		tw.wheel[i] = list.New()
	}
}

func (tw *TimeWheel) run() {
	for {
		select {
		case tasker := <-tw.setChan:
			tw.setTime(tasker)
		case tasker := <-tw.removeChan:
			tw.removeTime(tasker)
		case <-tw.ticker.C:
			tw.runTicker()
		}
	}
}

//SetTimer set a timer, key is the name of the timer,
//d will start to execute after how much time has passed, f is a callback of executing
func (tw *TimeWheel) SetTimer(key string, d time.Duration, f func()) error {
	if _, ok := tw.wMap[key]; ok {
		return ErrTimerKeyFound
	}

	tasker := &setingTasker{key: key, duration: d, callback: f}
	tw.setChan <- tasker
	return nil
}

//RemoveTimer remove a timer, key is the name of the timer
func (tw *TimeWheel) RemoveTimer(key string) error {
	if _, ok := tw.wMap[key]; !ok {
		return ErrTimerKeyNotFound
	}

	tasker := &removingTasker{key: key}
	tw.removeChan <- tasker
	return nil
}

//StopTimer stop all-timer, time wheel will stop schedule
func (tw *TimeWheel) StopTimer() {
	tw.ticker.Stop()
}

func (tw *TimeWheel) getPosition(d time.Duration) (circle int, pos int) {
	circle = (int(d) / int(time.Second) / tw.interval) / tw.sizeOfSlot
	pos = (tw.tickerPos + int(d)/int(time.Second)/tw.interval) % tw.sizeOfSlot
	return
}

func (tw *TimeWheel) setTime(tasker *setingTasker) {
	key, d, fn := tasker.key, tasker.duration, tasker.callback
	if int(d) < tw.interval {
		d = time.Duration(tw.interval)
	}
	circle, pos := tw.getPosition(d)
	slot := tw.wheel[pos]
	tw.insertSlot(slot, key, &task{circle: circle, pos: pos, fn: fn, slot: slot, key: key})
}

func (tw *TimeWheel) insertSlot(slot *list.List, key string, tasker *task) *list.Element {
	element := slot.PushBack(tasker)
	tw.wMap[key] = element
	return element
}

func (tw *TimeWheel) removeTime(rmTasker *removingTasker) {
	if element, ok := tw.wMap[rmTasker.key]; ok {
		tasker := element.Value.(*task)
		tasker.slot.Remove(element)
		delete(tw.wMap, rmTasker.key)
	}
}

func (tw *TimeWheel) runTicker() {
	tw.tickerPos = (tw.tickerPos + 1) % tw.sizeOfSlot
	slot := tw.wheel[tw.tickerPos]
	for e := slot.Front(); e != nil; e = e.Next() {
		tasker := e.Value.(*task)
		if tasker.circle > 0 {
			tasker.circle -= 1
			continue
		} else {
			tw.removeTime(&removingTasker{key: tasker.key})
			go func() {
				defer func() {
					if err := recover(); err != nil {
						log.Printf("tasker error: %v\n", err)
					}
				}()

				tasker.fn()
			}()
		}
	}
}
