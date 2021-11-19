
type ChatUsecase interface {
	GetHistory(idFrom, idTo) []models.Message, error
}
