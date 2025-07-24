package download

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/cheggaaa/pb/v3"
	"github.com/sant0x00/downloader-music/internal/domain"
	"github.com/sant0x00/downloader-music/internal/infrastructure/storage"
)

type HTTPDownloader struct {
	client            *http.Client
	repository        *storage.FileSystemRepository
	logger            domain.Logger
	concurrentWorkers int
	retryAttempts     int
	progressCallback  func(current, total int64, filename string)
}

func NewHTTPDownloader(repository *storage.FileSystemRepository, logger domain.Logger, concurrentWorkers, retryAttempts int, timeoutSeconds int) *HTTPDownloader {
	return &HTTPDownloader{
		client: &http.Client{
			Timeout: time.Duration(timeoutSeconds) * time.Second,
		},
		repository:        repository,
		logger:            logger,
		concurrentWorkers: concurrentWorkers,
		retryAttempts:     retryAttempts,
	}
}

func (d *HTTPDownloader) SetProgressCallback(callback func(current, total int64, filename string)) {
	d.progressCallback = callback
}

func (d *HTTPDownloader) Download(clipe domain.ClipeMusical, destPath string) error {
	if clipe.URLDownload == "" {
		return fmt.Errorf("URL de download não encontrada para o clipe: %s", clipe.Titulo)
	}

	filename := clipe.GetSanitizedFilename()
	if d.repository.Exists(filename) {
		d.logger.Info("Arquivo já existe, pulando", "arquivo", filename)
		return nil
	}

	err := d.repository.CreateDirectoryStructure(clipe)
	if err != nil {
		return err
	}

	filePath := d.repository.GetClipeFilePath(clipe)

	d.logger.Info("Iniciando download", "titulo", clipe.Titulo, "url", clipe.URLDownload, "destino", filePath)

	for attempt := 1; attempt <= d.retryAttempts; attempt++ {
		err = d.downloadFile(clipe.URLDownload, filePath, clipe.Titulo)
		if err == nil {
			d.logger.Info("Download concluído", "titulo", clipe.Titulo, "arquivo", filePath)
			return nil
		}

		if attempt < d.retryAttempts {
			d.logger.Warn("Falha no download, tentando novamente", "titulo", clipe.Titulo, "tentativa", attempt, "erro", err.Error())
			time.Sleep(time.Duration(attempt) * time.Second) // Progressive backoff
		}
	}

	return fmt.Errorf("falha no download após %d tentativas: %w", d.retryAttempts, err)
}

func (d *HTTPDownloader) DownloadBatch(clipes []domain.ClipeMusical, destPath string) error {
	d.logger.Info("Iniciando download em lote", "total_clipes", len(clipes), "workers", d.concurrentWorkers)

	jobs := make(chan domain.ClipeMusical, len(clipes))
	results := make(chan error, len(clipes))

	var wg sync.WaitGroup
	for i := 0; i < d.concurrentWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for clipe := range jobs {
				err := d.Download(clipe, destPath)
				results <- err
			}
		}()
	}

	for _, clipe := range clipes {
		jobs <- clipe
	}
	close(jobs)

	wg.Wait()
	close(results)

	var errors []error
	successCount := 0
	for err := range results {
		if err != nil {
			errors = append(errors, err)
		} else {
			successCount++
		}
	}

	d.logger.Info("Download em lote concluído", "sucessos", successCount, "erros", len(errors), "total", len(clipes))

	if len(errors) > 0 {
		d.logger.Error("Alguns downloads falharam", fmt.Errorf("%d erros encontrados", len(errors)))
		return errors[0]
	}

	return nil
}

func (d *HTTPDownloader) downloadFile(url, filePath, titulo string) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("erro ao criar requisição: %w", err)
	}

	req.Header.Set("User-Agent", "ClipesJW-Downloader/1.0")

	resp, err := d.client.Do(req)
	if err != nil {
		return fmt.Errorf("erro ao fazer requisição: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status code inválido: %d", resp.StatusCode)
	}

	tempFile := filePath + ".tmp"
	out, err := os.Create(tempFile)
	if err != nil {
		return fmt.Errorf("erro ao criar arquivo: %w", err)
	}
	defer out.Close()

	contentLength := resp.ContentLength

	var progressReader io.Reader = resp.Body
	if contentLength > 0 {
		bar := pb.Full.Start64(contentLength)
		bar.Set(pb.Bytes, true)
		progressReader = bar.NewProxyReader(resp.Body)
		defer bar.Finish()
	}

	_, err = io.Copy(out, progressReader)
	if err != nil {
		os.Remove(tempFile) // Clean up temporary file
		return fmt.Errorf("erro ao baixar arquivo: %w", err)
	}

	err = os.Rename(tempFile, filePath)
	if err != nil {
		os.Remove(tempFile) // Clean up temporary file
		return fmt.Errorf("erro ao finalizar arquivo: %w", err)
	}

	if d.progressCallback != nil && contentLength > 0 {
		d.progressCallback(contentLength, contentLength, titulo)
	}

	return nil
}
