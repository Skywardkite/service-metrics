package audit

type Observer interface {
	Notify(event AuditEvent)
}
