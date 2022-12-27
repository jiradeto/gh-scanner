package workerhandler

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/AlekSi/pointer"
	"github.com/jiradeto/gh-scanner/app/domain/entities"
	"github.com/jiradeto/gh-scanner/app/infrastructure/interfaces/github"
	"github.com/jiradeto/gh-scanner/app/infrastructure/interfaces/messagequeue"
	repositoryusecase "github.com/jiradeto/gh-scanner/app/usecases/repository"
	"github.com/jiradeto/gh-scanner/app/utils/gerrors"
)

func listFiles(root string) ([]string, error) {
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})

	return files, err
}

func search(fileName string, outPath string) ([]entities.ScanFinding, error) {
	var findings []entities.ScanFinding
	f, err := os.Open(fileName)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer f.Close()

	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(f)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)
	line := int64(1)
	for scanner.Scan() {
		// Check if the line contains the words "banana" or "apple"
		for _, rule := range ScannerRules {
			if strings.Contains(scanner.Text(), rule.Word) {
				fmt.Printf("%s:%d\n", fileName, line)
				filePath := strings.Replace(fileName, outPath+"/", "", -1)
				scanFinding := entities.ScanFinding{
					Type:   rule.Type,
					RuleID: rule.RuleId,
					Location: entities.Location{
						Path: filePath,
						Positions: entities.Positions{
							Begin: entities.Begin{
								Line: line,
							},
						},
					},
					Metadata: entities.Metadata{
						Description: rule.Description,
						Severity:    rule.Severity,
					},
				}
				findings = append(findings, scanFinding)
			}
		}
		line++
	}

	// Check for errors while scanning the file
	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}
	return findings, nil
}

func (w *Worker) PerformScan(params *messagequeue.StartScannerMessage) error {
	fmt.Println("PerformScan", params)
	ctx := context.Background()
	if params == nil {
		return gerrors.NewInternalError(nil)
	}
	queuedAt := time.Now()
	scanResult, err := w.RepositoryUsecase.UpdateOneScanResult(ctx, repositoryusecase.UpdateOneScanResultInput{
		ID:       pointer.ToString(params.ResultId),
		Status:   pointer.ToString(entities.ScanResultStatusInProgress.String()),
		QueuedAt: &queuedAt,
	})
	if err != nil {
		return err
	}

	github := github.New("scanning", params.URL)
	defer github.ClearnDirectory()
	err = github.Clone()
	if err != nil {
		return err
	}
	outputDirectory := *github.GetOutputPath()
	files, err := listFiles(outputDirectory)
	if err != nil {
		return err
	}
	fmt.Println("Total files", len(files))
	var allFindings []entities.ScanFinding
	for _, file := range files {
		findings, err := search(file, outputDirectory)
		if err != nil {
			continue
		}
		allFindings = append(allFindings, findings...)
	}

	fmt.Println("total findings ", len(allFindings))
	finishedAt := time.Now()
	_, err = w.RepositoryUsecase.UpdateOneScanResult(ctx, repositoryusecase.UpdateOneScanResultInput{
		ID:         scanResult.ID,
		Status:     pointer.ToString(entities.ScanResultStatusSusccess.String()),
		Findings:   allFindings,
		FinishedAt: &finishedAt,
	})
	if err != nil {
		return err
	}

	return nil
}
