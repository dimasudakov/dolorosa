package feature

import "time"

type FeatureRequest struct {
	EntityID  func( /* some state struct */ ) string
	Namespace string
	Name      string
	// время, которое добавив к now при формировании запроса получим "ОТ"
	TimeFrom time.Duration
	// время, которое добавив к now при формировании запроса получим "ДО"
	TimeTo time.Duration
	// Если нужно получить фичу не по времени, а только текущий бакет
	SingleBucket bool
	// Если фича холодная - бакета у неё в запросе быть не должно
	ColdFeature bool
}
