package producer

// CloseAndWait ожидает отсылки всех сообщений
func (s *AsyncProducer) CloseAndWait() {
	s.wg.Wait()
}
