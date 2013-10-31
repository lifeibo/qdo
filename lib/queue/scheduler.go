package queue

import (
	"bytes"
	"fmt"
	"time"

	"github.com/borgenk/qdo/third_party/github.com/syndtr/goleveldb/leveldb/comparer"
	"github.com/borgenk/qdo/third_party/github.com/syndtr/goleveldb/leveldb/opt"

	"github.com/borgenk/qdo/lib/log"
)

type Scheduler struct {
	// Conveyor.
	Conveyor *Conveyor `json:"-"`

	// Conveyor status signal.
	notifySignal chan convSignal `json:"-"`

	// Conveyor scheduler id name.
	ScheduleId string `json:"schedule_id"`

	// Conveyor scheduler list name.
	ScheduleList string `json:"schedule_list"`

	// How often scheduler checks schedule list in seconds.
	Rate time.Duration `json:"rate"`
}

func NewScheduler(conveyor *Conveyor) *Scheduler {
	scheduler := &Scheduler{
		Conveyor: conveyor,
		Rate:     5 * time.Second,
	}
	return scheduler
}

func (sched *Scheduler) Start() {
	sched.notifySignal = make(chan convSignal)

	for {
		select {
		case sig := <-sched.notifySignal:
			if sig == stop {
				log.Infof("stopping scheduler for conveyor %s", sched.Conveyor.ID)
				return
			}
		default:
		}

		// Fetch all tasks scheduled earlier then right now, starting with the
		// oldest (lowest) possible item.
		stop := append(sched.Conveyor.waitKeyStart, []byte(string(time.Now().Unix()))...)
		iter := db.NewIterator(nil)
		for iter.Seek(sched.Conveyor.waitKeyStart); iter.Valid(); iter.Next() { // start := []byte("w\x00")
			k := iter.Key()
			v := iter.Value()

			if comparer.DefaultComparer.Compare(k, stop) > 0 { // []byte(fmt.Sprintf("w\x00%d", now))
				// This might be a task is sechduled in the future or some other
				// stored value. All scheduled tasks up until right now is read.
				break
			}

			// Parse out task id in order to avoid decode/encode gob.
			i := bytes.LastIndex(k, []byte(startPoint))
			taskId := string(k[i:len(k)])
			log.Infof("placing scheduled task %s into queue", taskId)

			// TODO: Batch add / removal of task.
			sched.Conveyor.add("", v)
			db.Delete(k, nil)
		}
		iter.Release()

		time.Sleep(sched.Rate)
	}
}

func (sched *Scheduler) Stop() {

}

func (sched *Scheduler) Add() {

}

// Reschedule task.
func (sched *Scheduler) Reschedule(task *Task) (int32, error) {
	if task.Delay == 0 {
		task.Delay = 1
	}
	task.Delay = task.Delay * 2
	task.Tries = task.Tries + 1

	t, err := GobEncode(task)
	if err != nil {
		log.Error("", err)
		return 0, err
	}

	retryAt := time.Now().Unix() + int64(task.Delay)

	// Key format: [convid]\x00q\x00[timestamp]\x00[taskid]
	key := fmt.Sprintf("w\x00%d\x00%s", retryAt, task.ID)

	// sched.conv.waitKeyStart

	err = db.Put([]byte(key), t, nil)
	if err != nil {
		log.Error("add task to db failed", err)
		return 0, err
	}

	return task.Delay, nil
}
