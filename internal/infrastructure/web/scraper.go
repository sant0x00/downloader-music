package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/sant0x00/downloader-music/internal/domain"
)

type JWAPIResponse struct {
	PubName string                     `json:"pubName"`
	Files   map[string]JWLanguageFiles `json:"files"`
}

type JWLanguageFiles struct {
	MP3 []JWAudioFile `json:"MP3"`
}

type JWAudioFile struct {
	Title    string `json:"title"`
	File     JWFile `json:"file"`
	FileSize int    `json:"filesize"`
}

type JWFile struct {
	URL              string `json:"url"`
	ModifiedDatetime string `json:"modifiedDatetime"`
	Checksum         string `json:"checksum"`
}

type JWScraper struct {
	client        *http.Client
	userAgent     string
	delay         time.Duration
	logger        domain.Logger
	downloadURL   string
	downloadCache map[string]string
}

func NewJWScraper(userAgent string, delay time.Duration, logger domain.Logger) *JWScraper {
	return &JWScraper{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		userAgent:     userAgent,
		delay:         delay,
		logger:        logger,
		downloadURL:   "https://b.jw-cdn.org/apis/pub-media/GETPUBMEDIALINKS?output=json&pub=osg&fileformat=MP3%2CAAC&alllangs=0&langwritten=T&txtCMSLang=T",
		downloadCache: make(map[string]string),
	}
}

func (s *JWScraper) ScrapClipesList(url string) ([]domain.ClipeMusical, error) {
	s.logger.Info("Iniciando scraping da lista de clipes", "url", url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar requisição: %w", err)
	}

	req.Header.Set("User-Agent", s.userAgent)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "pt-BR,pt;q=0.9,en;q=0.8")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("erro ao fazer requisição: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code inválido: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("erro ao parsear HTML: %w", err)
	}

	var clipes []domain.ClipeMusical

	doc.Find("h2 a[href*='/biblioteca/musica-canticos/clipes-musicais/']").Each(func(i int, sel *goquery.Selection) {
		href, exists := sel.Attr("href")
		if !exists {
			return
		}

		titulo := strings.TrimSpace(sel.Text())
		if titulo == "" {
			return
		}

		fullURL := href
		if strings.HasPrefix(href, "/") {
			fullURL = "https://www.jw.org" + href
		}

		id := s.extractClipeID(href)

		ano := s.extractYearFromTitle(titulo)

		clipe := domain.ClipeMusical{
			ID:     id,
			Titulo: titulo,
			URL:    fullURL,
			Ano:    ano,
		}

		clipes = append(clipes, clipe)
	})

	s.logger.Info("Scraping da lista concluído", "total_clipes", len(clipes))
	return clipes, nil
}

func (s *JWScraper) ScrapClipeDetails(clipe domain.ClipeMusical) (domain.ClipeMusical, error) {
	s.logger.Debug("Obtendo detalhes do clipe", "titulo", clipe.Titulo, "url", clipe.URL)

	downloadURL, err := s.findDownloadURLForClipe(clipe.Titulo)
	if err != nil {
		s.logger.Error("Erro ao buscar URL de download", err, "titulo", clipe.Titulo)
		return clipe, err
	}

	if downloadURL != "" {
		clipe.URLDownload = downloadURL
		s.logger.Debug("URL de download encontrada", "titulo", clipe.Titulo, "url", downloadURL)
	}

	if clipe.Ano == 0 {
		clipe.Ano = s.extractYearFromTitle(clipe.Titulo)
	}

	time.Sleep(s.delay)

	clipe.Descricao = fmt.Sprintf("Clipe musical: %s", clipe.Titulo)

	s.logger.Debug("Detalhes obtidos", "titulo", clipe.Titulo, "download_url", clipe.URLDownload, "ano", clipe.Ano)
	return clipe, nil
}

func (s *JWScraper) extractClipeID(url string) string {
	parts := strings.Split(url, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return ""
}

func (s *JWScraper) extractYearFromTitle(titulo string) int {
	re := regexp.MustCompile(`\b(202[0-9])\b`)
	matches := re.FindStringSubmatch(titulo)
	if len(matches) > 1 {
		if year, err := strconv.Atoi(matches[1]); err == nil {
			return year
		}
	}

	re2 := regexp.MustCompile(`congresso de (\d{4})`)
	matches2 := re2.FindStringSubmatch(titulo)
	if len(matches2) > 1 {
		if year, err := strconv.Atoi(matches2[1]); err == nil {
			return year
		}
	}

	re3 := regexp.MustCompile(`de (\d{4})`)
	matches3 := re3.FindStringSubmatch(titulo)
	if len(matches3) > 1 {
		if year, err := strconv.Atoi(matches3[1]); err == nil {
			return year
		}
	}

	return 0
}

func (s *JWScraper) findDownloadURLForClipe(titulo string) (string, error) {
	if url, exists := s.downloadCache[titulo]; exists {
		return url, nil
	}

	if len(s.downloadCache) == 0 {
		err := s.loadDownloadCache()
		if err != nil {
			return "", err
		}
	}

	if url, exists := s.downloadCache[titulo]; exists {
		return url, nil
	}

	s.logger.Warn("URL de download não encontrada para clipe", "titulo", titulo)
	return "", nil
}

func (s *JWScraper) loadDownloadCache() error {
	s.logger.Info("Carregando cache de downloads via API JSON", "url", s.downloadURL)

	req, err := http.NewRequest("GET", s.downloadURL, nil)
	if err != nil {
		return fmt.Errorf("erro ao criar requisição: %w", err)
	}

	req.Header.Set("User-Agent", s.userAgent)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-Language", "pt-BR,pt;q=0.9,en;q=0.8")

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("erro ao fazer requisição: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status code inválido: %d", resp.StatusCode)
	}

	var apiResponse JWAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return fmt.Errorf("erro ao decodificar JSON: %w", err)
	}

	foundLinks := 0

	for langCode, langFiles := range apiResponse.Files {
		s.logger.Debug("Processando idioma", "lang", langCode, "mp3_count", len(langFiles.MP3))

		for _, audioFile := range langFiles.MP3 {
			if audioFile.File.URL != "" && audioFile.Title != "" {
				titulo := strings.TrimSpace(audioFile.Title)

				titulo = strings.TrimPrefix(titulo, "Reproduzir")
				titulo = strings.TrimSpace(titulo)

				if titulo != "" && len(titulo) > 2 {
					s.downloadCache[titulo] = audioFile.File.URL
					s.logger.Debug("Link de download adicionado ao cache (API)",
						"titulo", titulo,
						"url", audioFile.File.URL,
						"filesize", audioFile.FileSize)
					foundLinks++
				}
			}
		}
	}

	s.logger.Info("Cache de downloads carregado via API", "total_links", len(s.downloadCache))
	return nil
}
