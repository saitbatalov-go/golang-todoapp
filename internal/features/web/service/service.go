package web_service

type WebService struct {
	webRepository WebRepository
}

type WebRepository interface {
	GetFile(path string) ([]byte, error)
}

func NewWebService(webRepository WebRepository) *WebService {
	return &WebService{
		webRepository: webRepository,
	}
}
