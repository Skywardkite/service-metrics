package audit

type AuditEvent struct {
	TS        int64    `json:"ts"`         // unix timestamp события
	Metrics   []string `json:"metrics"`    // наименование полученных метрик
	IPAddress string   `json:"ip_address"` // IP адрес входящего запроса
}

type AuditPublisher struct {
	observers []Observer
}

func NewAuditPublisher() *AuditPublisher {
	return &AuditPublisher{
		observers: make([]Observer, 0),
	}
}

// Добавляем нового наблюдателя
func (p *AuditPublisher) Subscribe(o Observer) {
	p.observers = append(p.observers, o)
}

// Публикуем событие всем подписчикам
func (p *AuditPublisher) Publish(event AuditEvent) {
	for _, o := range p.observers {
		o.Notify(event)
	}
}
