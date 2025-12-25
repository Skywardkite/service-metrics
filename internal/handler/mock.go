package handler

import (
	"context"

	"github.com/Skywardkite/service-metrics/internal/audit"
)

type MockService struct{}

func (m *MockService) UpdateMetric(ctx context.Context, metricType, metricName, metricValue string) error {
	return nil
}

type MockAudit struct{}

func (m *MockAudit) Publish(event audit.AuditEvent) {}
