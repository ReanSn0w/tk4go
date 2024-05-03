package tools

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// NewShutdownStack создает сруктуру для плавного отключения приложения
func NewShutdownStack(logger Logger) *ShutdownStack {
	return &ShutdownStack{
		fn:  make([]ShutdownFunc, 0),
		log: logger,
	}
}

// ShutdownStack структура для плавного отключения приложения
//
// Содержит стек функций которые необходимо выполнить перед завершением приложения
// и методы для их добавления и выполнения
type ShutdownStack struct {
	fn  []ShutdownFunc
	log Logger
}

type ShutdownFunc func(context.Context)

// Add добавляет функции которые необходимо выполнить
// перед завершением приложения в стек
func (st *ShutdownStack) Add(fn ...ShutdownFunc) {
	st.fn = append(st.fn, fn...)
}

// Wait подписывается на системные уведомления об отклбчении
// и производит запуск функции Shutdown после получения сигналов
// syscall.SIGTERM, syscall.SIGINT
func (st *ShutdownStack) Wait(ctx context.Context, timeout time.Duration) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	select {
	case <-ctx.Done():
		st.log.Logf("[INFO] Получен сигнал отмены из контекста")
	case registredSignal := <-quit:
		st.log.Logf("[INFO] Зарегистрирован системный сигнал: %s", registredSignal.String())
	}

	ctx, done := context.WithTimeout(context.Background(), timeout)
	st.shutdown(ctx, done)
}

// Shutdown производит немедленное отключение приложения
// путем последовательного запуска всех функций отключения
func (st *ShutdownStack) shutdown(ctx context.Context, done func()) {
	go func() {
		for _, fn := range st.fn {
			if fn != nil {
				fn(ctx)
			}
		}

		done()
	}()

	<-ctx.Done()
}
