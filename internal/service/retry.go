package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
)

// withRetry запускает повторение попыток успешно выполнить передаваемую функцию.
// Выполним операцию сразу, чреез 1, 3 и 5 секунд. Итого 4 раза.
// В случае неудачи или другой ошибки - вернем ее.
// При успехе - возврат nil.
func withRetry(fn func() error) error {
	delays := []time.Duration{time.Second, 3 * time.Second, 5 * time.Second}

	var lastErr error
	for attempt := 0; attempt <= len(delays); attempt++ {
		err := fn()
		if err == nil {
			return nil
		}

		if !isRetriable(err) {
			return err
		}

		lastErr = err
		if attempt < len(delays) {
			time.Sleep(delays[attempt])
		}
	}
	return fmt.Errorf("операция не удалась после %d попыток: %w", len(delays)+1, lastErr)
}

func isRetriable(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case pgerrcode.ConnectionException, pgerrcode.AdminShutdown:
			return true
		}
	}
	return false
}
