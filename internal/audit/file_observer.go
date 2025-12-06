package audit

import (
	"encoding/json"
	"os"
)

type FileObserver struct {
	FilePath string
}

func (f *FileObserver) Notify(event AuditEvent) {
	data, _ := json.Marshal(event)
	os.WriteFile(f.FilePath, append(data, '\n'), 0644)
}
