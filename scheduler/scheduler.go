package scheduler

var sched chan func() // scheduler buffer

// Scheduler
func RunScheduler() {
	sched = make(chan func(), 16)
	for {
		fn := <- sched
		go fn()
	}
}
func Schedule(fn func()) {
	sched <- fn
}
