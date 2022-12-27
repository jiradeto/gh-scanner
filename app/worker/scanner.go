package scannerworker

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/AlekSi/pointer"
	"github.com/Shopify/sarama"
	"github.com/jiradeto/gh-scanner/app/entities"
	"github.com/jiradeto/gh-scanner/app/infrastructure/interfaces/github"
	"github.com/jiradeto/gh-scanner/app/infrastructure/interfaces/messagequeue"
	repositoryusecase "github.com/jiradeto/gh-scanner/app/usecases/repository"
	"github.com/jiradeto/gh-scanner/app/utils/fileutils"
	"github.com/jiradeto/gh-scanner/app/utils/loggers"
	"github.com/rs/xid"
)

func searchByLine(fileName string, outPath string) ([]entities.ScanFinding, error) {
	var findings []entities.ScanFinding
	f, err := os.Open(fileName)
	if err != nil {
		loggers.Text.Info("unable to open file", err.Error())
		return nil, err
	}
	defer f.Close()

	const maxBufferSize int = 10000000
	scanner := bufio.NewScanner(f)
	buf := make([]byte, maxBufferSize)
	scanner.Buffer(buf, maxBufferSize)
	line := int64(1)
	scannerRules := GetScannerRules()
	for scanner.Scan() {
		for _, rule := range scannerRules {
			if strings.Contains(scanner.Text(), rule.Word) {
				basePath := fileutils.GetBasePath(fileName, outPath)
				newFinding := entities.NewScanFunding(basePath, line, rule)
				findings = append(findings, newFinding)
			}
		}
		line++
	}

	if err := scanner.Err(); err != nil {
		loggers.Text.Info("an error when scanning files")
	}
	return findings, nil
}

func (s *ScannerWorker) RunScanner() error {
	requestID := xid.New().String()
	github := github.New(s.configs.URL)
	defer github.CleanDirectory()
	err := github.Clone()
	if err != nil {
		return err
	}
	outputDirectory := *github.GetOutputPath()
	files, err := fileutils.ListFiles(outputDirectory)
	if err != nil {
		return err
	}
	start := time.Now()
	logInfo := map[string]interface{}{
		"url":        s.configs.URL,
		"totalFiles": len(files),
	}
	loggers.Text.Info(fmt.Sprintf("Scanner %v - started", requestID))
	loggers.JSON.Info(logInfo)
	var allFindings []entities.ScanFinding
	for _, file := range files {
		findings, err := searchByLine(file, outputDirectory)
		if err != nil {
			continue
		}
		allFindings = append(allFindings, findings...)
	}
	s.findings = allFindings
	logInfo = map[string]interface{}{
		"url":           s.configs.URL,
		"totalFiles":    len(files),
		"totalFindings": len(allFindings),
		"totalTime":     time.Since(start),
	}
	loggers.Text.Info(fmt.Sprintf("Scanner #%v - finished", requestID))
	loggers.JSON.Info(logInfo)
	return nil
}

func (s *ScannerWorker) dequeueScanResult() (*entities.ScanResult, error) {
	scanResult, err := s.RepositoryUsecase.UpdateOneScanResult(context.Background(),
		repositoryusecase.UpdateOneScanResultInput{
			ID:       pointer.ToString(s.configs.ResultId),
			Status:   pointer.ToString(entities.ScanResultStatusInProgress.String()),
			QueuedAt: pointer.ToTime(time.Now()),
		})
	if err != nil {
		return nil, err
	}
	return scanResult, nil
}

func (s *ScannerWorker) saveScanResult(status entities.ScanResultStatus) (*entities.ScanResult, error) {
	scanResult, err := s.RepositoryUsecase.UpdateOneScanResult(context.Background(),
		repositoryusecase.UpdateOneScanResultInput{
			ID:         &s.configs.ResultId,
			Status:     pointer.ToString(entities.ScanResultStatusSusccess.String()),
			Findings:   s.findings,
			FinishedAt: pointer.ToTime(time.Now()),
		})
	if err != nil {
		return nil, err
	}
	return scanResult, nil
}

func (s *ScannerWorker) parseConsumerMessage(msg *sarama.ConsumerMessage) (*messagequeue.StartScannerMessage, error) {
	message := string(msg.Value)
	var params messagequeue.StartScannerMessage
	err := json.Unmarshal([]byte(message), &params)
	if err != nil {
		return nil, err
	}
	return &params, nil
}

func (s *ScannerWorker) Start(msg *sarama.ConsumerMessage) {
	configs, err := s.parseConsumerMessage(msg)
	if err != nil || configs == nil {
		loggers.Text.Info("Scanner failed to parse consumer message: ", err.Error())
		return
	}
	s.configs = configs
	if _, err := s.dequeueScanResult(); err != nil {
		loggers.Text.Info("Scanner failed update scan result", err.Error())
		return
	}
	err = s.RunScanner()
	scanStatus := entities.ScanResultStatusSusccess
	if err != nil {
		scanStatus = entities.ScanResultStatusFailure
	}
	_, err = s.saveScanResult(scanStatus)
	if err != nil {
		loggers.Text.Info("Scanner failed update scan result", err.Error())
	}
}
