package audit

// AuditEvent - событие аудита, которое будем рассылать всем подписчикам.
type AuditEvent struct {
	// TS - unix timestamp события.
	TS int64 `json:"ts"`
	// Metrics - наименование полученных метрик.
	Metrics []string `json:"metrics"`
	// IPAddress - IP адрес входящего запроса.
	IPAddress string `json:"ip_address"`
}

type AuditPublisher struct {
	observers []Observer
}

func NewAuditPublisher() *AuditPublisher {
	return &AuditPublisher{
		observers: make([]Observer, 0),
	}
}

// Subscribe добавляет нового наблюдателя.
func (p *AuditPublisher) Subscribe(o Observer) {
	p.observers = append(p.observers, o)
}

// Publish рассылает событие всем подписчикам.
func (p *AuditPublisher) Publish(event AuditEvent) {
	for _, o := range p.observers {
		o.Notify(event)
	}
}
