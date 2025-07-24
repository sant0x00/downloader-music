package domain

import (
	"fmt"
	"time"
)

type ClipeMusical struct {
	ID             string
	Titulo         string
	Descricao      string
	URL            string
	URLDownload    string
	TamanhoArquivo int64
	DataPublicacao time.Time
	NomeArquivo    string
	Ano            int
}

func (c *ClipeMusical) IsValid() bool {
	return c.Titulo != "" && c.URL != "" && c.URLDownload != ""
}

func (c *ClipeMusical) GetSanitizedFilename() string {
	if c.NomeArquivo != "" {
		return c.NomeArquivo
	}

	filename := c.Titulo
	sanitized := ""
	for _, char := range filename {
		switch {
		case char >= 'a' && char <= 'z',
			char >= 'A' && char <= 'Z',
			char >= '0' && char <= '9':
			sanitized += string(char)
		case char == ' ':
			sanitized += "_"
		case char == 'ã', char == 'á', char == 'à', char == 'â':
			sanitized += "a"
		case char == 'é', char == 'ê':
			sanitized += "e"
		case char == 'í':
			sanitized += "i"
		case char == 'õ', char == 'ó', char == 'ô':
			sanitized += "o"
		case char == 'ú', char == 'û':
			sanitized += "u"
		case char == 'ç':
			sanitized += "c"
		default:
			// Ignore any other special characters
		}
	}

	c.NomeArquivo = sanitized + ".mp3"
	return c.NomeArquivo
}

func (c *ClipeMusical) GetDirectoryPath() string {
	if c.Ano > 0 {
		return fmt.Sprintf("%d", c.Ano)
	}
	return "outros"
}
