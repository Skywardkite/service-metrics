package handler

import (
	"github.com/Skywardkite/service-metrics/internal/audit"
)

type MockAudit struct{}

func (m *MockAudit) Publish(event audit.AuditEvent) {}
