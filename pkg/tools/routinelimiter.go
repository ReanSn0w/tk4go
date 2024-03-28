package tools

// NewRoutineLimiter Создает новую структуру для ограничения
// количества одновременно запущенных горутин
func NewRoutineLimiter(max int) *RoutineLimiter {
	return &RoutineLimiter{
		ch: make(chan struct{}, max),
	}
}

type RoutineLimiter struct {
	ch chan struct{}
}

// Run запускает задачу в горутине
// в случае если доступных слотов нет,
// то ожидает освобождения слота
// перед запуском задачи
func (r *RoutineLimiter) Run(task func()) {
	r.ch <- struct{}{}

	go func() {
		task()
		<-r.ch
	}()
}
