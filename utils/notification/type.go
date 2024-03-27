package notification

type INotificationService interface {
	SendNotification(payload EmailPayload) error
}
