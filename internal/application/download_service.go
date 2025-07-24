package application

import (
	"fmt"

	"github.com/sant0x00/downloader-music/internal/domain"
)

type DownloadService struct {
	scraper    domain.WebScraper
	downloader domain.DownloadService
	repository domain.ClipeRepository
	logger     domain.Logger
}

func NewDownloadService(
	scraper domain.WebScraper,
	downloader domain.DownloadService,
	repository domain.ClipeRepository,
	logger domain.Logger,
) *DownloadService {
	return &DownloadService{
		scraper:    scraper,
		downloader: downloader,
		repository: repository,
		logger:     logger,
	}
}

func (s *DownloadService) DownloadAllClipes(baseURL string) error {
	s.logger.Info("Iniciando processo de download de todos os clipes")

	s.logger.Info("Fazendo scraping da lista de clipes", "url", baseURL)
	clipes, err := s.scraper.ScrapClipesList(baseURL)
	if err != nil {
		s.logger.Error("Erro ao fazer scraping da lista", err)
		return fmt.Errorf("erro ao obter lista de clipes: %w", err)
	}

	if len(clipes) == 0 {
		s.logger.Warn("Nenhum clipe encontrado na página")
		return fmt.Errorf("nenhum clipe encontrado")
	}

	s.logger.Info("Lista de clipes obtida", "total", len(clipes))

	s.logger.Info("Obtendo detalhes dos clipes")
	var clipesValidos []domain.ClipeMusical

	for i, clipe := range clipes {
		s.logger.Debug("Processando clipe", "index", i+1, "total", len(clipes), "titulo", clipe.Titulo)

		clipeDetalhado, err := s.scraper.ScrapClipeDetails(clipe)
		if err != nil {
			s.logger.Error("Erro ao obter detalhes do clipe", err, "titulo", clipe.Titulo)
			continue
		}

		if !clipeDetalhado.IsValid() {
			s.logger.Warn("Clipe inválido, pulando", "titulo", clipe.Titulo, "url_download", clipeDetalhado.URLDownload)
			continue
		}

		clipesValidos = append(clipesValidos, clipeDetalhado)
	}

	if len(clipesValidos) == 0 {
		s.logger.Error("Nenhum clipe válido encontrado", fmt.Errorf("sem clipes para download"))
		return fmt.Errorf("nenhum clipe válido encontrado")
	}

	s.logger.Info("Clipes válidos encontrados", "total", len(clipesValidos))

	var clipesParaDownload []domain.ClipeMusical
	for _, clipe := range clipesValidos {
		filename := clipe.GetSanitizedFilename()
		if !s.repository.Exists(filename) {
			clipesParaDownload = append(clipesParaDownload, clipe)
		} else {
			s.logger.Info("Clipe já existe, pulando", "titulo", clipe.Titulo, "arquivo", filename)
		}
	}

	if len(clipesParaDownload) == 0 {
		s.logger.Info("Todos os clipes já foram baixados")
		return nil
	}

	s.logger.Info("Clipes para download", "novos", len(clipesParaDownload), "existentes", len(clipesValidos)-len(clipesParaDownload))

	outputDir := s.repository.GetOutputDirectory()
	err = s.downloader.DownloadBatch(clipesParaDownload, outputDir)
	if err != nil {
		s.logger.Error("Erro durante download em lote", err)
		return fmt.Errorf("erro durante download: %w", err)
	}

	s.logger.Info("Processo de download concluído com sucesso", "total_baixados", len(clipesParaDownload))
	return nil
}

func (s *DownloadService) CheckForNewClipes(baseURL string) ([]domain.ClipeMusical, error) {
	s.logger.Info("Verificando novos clipes disponíveis")

	clipes, err := s.scraper.ScrapClipesList(baseURL)
	if err != nil {
		return nil, fmt.Errorf("erro ao obter lista de clipes: %w", err)
	}

	var novosClipes []domain.ClipeMusical
	for _, clipe := range clipes {
		filename := clipe.GetSanitizedFilename()
		if !s.repository.Exists(filename) {
			novosClipes = append(novosClipes, clipe)
		}
	}

	s.logger.Info("Verificação concluída", "total_clipes", len(clipes), "novos_clipes", len(novosClipes))
	return novosClipes, nil
}

func (s *DownloadService) DownloadSpecificClipe(baseURL, titulo string) error {
	s.logger.Info("Procurando clipe específico", "titulo", titulo)

	clipes, err := s.scraper.ScrapClipesList(baseURL)
	if err != nil {
		return fmt.Errorf("erro ao obter lista de clipes: %w", err)
	}

	var clipeEncontrado *domain.ClipeMusical
	for _, clipe := range clipes {
		if clipe.Titulo == titulo {
			clipeEncontrado = &clipe
			break
		}
	}

	if clipeEncontrado == nil {
		return fmt.Errorf("clipe não encontrado: %s", titulo)
	}

	clipeDetalhado, err := s.scraper.ScrapClipeDetails(*clipeEncontrado)
	if err != nil {
		return fmt.Errorf("erro ao obter detalhes do clipe: %w", err)
	}

	if !clipeDetalhado.IsValid() {
		return fmt.Errorf("clipe inválido: %s", titulo)
	}

	filename := clipeDetalhado.GetSanitizedFilename()
	if s.repository.Exists(filename) {
		s.logger.Info("Clipe já existe", "titulo", titulo, "arquivo", filename)
		return nil
	}

	outputDir := s.repository.GetOutputDirectory()
	err = s.downloader.Download(clipeDetalhado, outputDir)
	if err != nil {
		return fmt.Errorf("erro no download: %w", err)
	}

	s.logger.Info("Download do clipe específico concluído", "titulo", titulo)
	return nil
}
