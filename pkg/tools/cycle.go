package tools

import (
	"context"
	"time"
)

// Функция создает новую структуру для решения
// повторяющихся задач
func NewCycleTask(task func()) *CycleTask {
	return &CycleTask{
		timer: time.NewTimer(0),
		task:  task,
	}
}

// Структура для работы с циклической задачей
type CycleTask struct {
	timer  *time.Timer
	task   func()
	cancel func()
}

// Запуск задачи единоразово здесь и сейчас
func (r *CycleTask) Once() {
	r.task()
}

// Останавливает цикл выполнения задачи
func (r *CycleTask) Stop() {
	if r.cancel != nil {
		r.cancel()
	}
}

// Запускает цикл выполнения задачи
// в случае если задача уже запущена
func (r *CycleTask) Run(d time.Duration) {
	ctx, cancel := context.WithCancel(context.Background())

	r.cancel = cancel
	r.timer.Reset(d)

	go func() {
		for {
			select {
			case <-r.timer.C:
				r.task()
				r.timer.Reset(d)
			case <-ctx.Done():
				r.cancel = nil
				break
			}
		}
	}()
}
