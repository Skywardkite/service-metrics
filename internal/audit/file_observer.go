package audit

import (
	"encoding/json"
	"os"
	"sync"

	"github.com/Skywardkite/service-metrics/internal/logger"
)

type FileObserver struct {
	FilePath string
	mu       sync.Mutex
}

func (f *FileObserver) Notify(event AuditEvent) {
	data, err := json.Marshal(event)
	if err != nil {
		logger.Sugar.Errorw("Can't marshal audit event", "error", err)
		return
	}

	f.mu.Lock()
	defer f.mu.Unlock()

	file, err := os.OpenFile(f.FilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		logger.Sugar.Errorw("Can't open audit file", "error", err)
		return
	}
	defer file.Close()

	if _, err := file.Write(append(data, '\n')); err != nil {
		logger.Sugar.Errorw("Can't write audit event", "error", err)
	}
}
