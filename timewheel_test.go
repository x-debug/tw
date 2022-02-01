package tw

import (
	"testing"
	"time"
)

func TestTimeWheel_SetTimer(t *testing.T) {
	wheel := NewTimeWheel(1, 60)
	defer wheel.StopTimer()

	var validate int
	validate = 0
	_ = wheel.SetTimer("timer1", 2*time.Second, func() {
		validate = 1
	})
	time.Sleep(5 * time.Second)
	if validate != 1 {
		t.Fail()
	}
}

func TestTimeWheel_StopTimer(t *testing.T) {
	wheel := NewTimeWheel(1, 60)
	wheel.StopTimer()

	var validate int
	validate = 0
	_ = wheel.SetTimer("timer1", 2*time.Second, func() {
		validate = 1
	})
	time.Sleep(5 * time.Second)
	if validate == 1 {
		t.Fail()
	}
}

func TestTimeWheel_RemoveTimer(t *testing.T) {
	wheel := NewTimeWheel(1, 60)
	defer wheel.StopTimer()

	var validate int
	validate = 0
	_ = wheel.SetTimer("timer1", 2*time.Second, func() {
		validate = 1
	})
	_ = wheel.SetTimer("timer2", 5*time.Second, func() {
		validate = 2
	})
	_ = wheel.RemoveTimer("timer1")
	time.Sleep(8 * time.Second)
	if validate != 2 {
		t.Fail()
	}
}
