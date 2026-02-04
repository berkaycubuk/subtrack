package services

type MockTelegramService struct {
	SendNotificationFunc func(name string, price float64, currency string, days int, cycle string, paymentDate string) error
	SendMessageFunc      func(text string) error
	HealthCheckFunc      func() error
}

func (m *MockTelegramService) SendNotification(name string, price float64, currency string, days int, cycle string, paymentDate string) error {
	if m.SendNotificationFunc != nil {
		return m.SendNotificationFunc(name, price, currency, days, cycle, paymentDate)
	}
	return nil
}

func (m *MockTelegramService) SendMessage(text string) error {
	if m.SendMessageFunc != nil {
		return m.SendMessageFunc(text)
	}
	return nil
}

func (m *MockTelegramService) HealthCheck() error {
	if m.HealthCheckFunc != nil {
		return m.HealthCheckFunc()
	}
	return nil
}
