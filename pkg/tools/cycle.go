package tools

import (
	"time"
)

// Функция создает новую структуру для решения
// повторяющихся задач
func NewCycleTask(task func()) *CycleTask {

	return &CycleTask{
		timer: time.NewTimer(0),
		done:  make(chan bool),
		task:  task,
	}
}

// Структура для работы с циклической задачей
type CycleTask struct {
	timer *time.Timer
	done  chan bool
	task  func()
}

// Запуск задачи единоразово здесь и сейчас
func (r *CycleTask) Once() {
	r.task()
}

// Останавливает цикл выполнения задачи
func (r *CycleTask) Stop() {
	if r.timer.Stop() {
		r.done <- true
	}
}

// Запускает цикл выполнения задачи
// в случае если задача уже запущена
func (r *CycleTask) Run(d time.Duration) {
	r.timer.Reset(d)

	go func() {
		loop := true

		for loop {
			select {
			case <-r.timer.C:
				r.task()
				r.timer.Reset(d)
			case <-r.done:
				r.timer.Stop()
				loop = false
			}
		}
	}()
}
