package domain

type ClipeRepository interface {
	FindAll() ([]ClipeMusical, error)
	Save(clipe ClipeMusical) error
	Exists(filename string) bool
	GetOutputDirectory() string
	CreateDirectoryStructure(clipe ClipeMusical) error
}

type WebScraper interface {
	ScrapClipesList(url string) ([]ClipeMusical, error)
	ScrapClipeDetails(clipe ClipeMusical) (ClipeMusical, error)
}

type DownloadService interface {
	Download(clipe ClipeMusical, destPath string) error
	DownloadBatch(clipes []ClipeMusical, destPath string) error
	SetProgressCallback(callback func(current, total int64, filename string))
}

type Logger interface {
	Info(msg string, fields ...interface{})
	Error(msg string, err error, fields ...interface{})
	Debug(msg string, fields ...interface{})
	Warn(msg string, fields ...interface{})
}
