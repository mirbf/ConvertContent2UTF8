package convertcontent2utf8

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	encoding "github.com/mirbf/encoding-processor"
)

// ConvertFile 转换单个文件
func ConvertFile(inputFile, outputFile string, options ...Option) (*ConvertResult, error) {
	config := applyOptions(options)

	// 参数验证
	if inputFile == "" {
		return nil, fmt.Errorf("input file path cannot be empty")
	}
	if outputFile == "" {
		return nil, fmt.Errorf("output file path cannot be empty")
	}

	// 检查输入文件是否存在
	if _, err := os.Stat(inputFile); os.IsNotExist(err) {
		return nil, fmt.Errorf("input file does not exist: %s", inputFile)
	}

	// 进度回调
	if config.ProgressCallback != nil {
		progress := Progress{
			CurrentFile:    inputFile,
			ProcessedFiles: 0,
			TotalFiles:     1,
			Status:         StatusStarting,
			StartTime:      time.Now(),
			ErrorCount:     0,
		}
		config.ProgressCallback(progress)
	}

	start := time.Now()

	// 创建智能编码处理器
	processor := encoding.NewSmartProcessor()

	// 进度更新 - 处理中
	if config.ProgressCallback != nil {
		progress := Progress{
			CurrentFile:    inputFile,
			ProcessedFiles: 0,
			TotalFiles:     1,
			Status:         StatusProcessing,
			StartTime:      time.Now(),
			ElapsedTime:    time.Since(start),
			ErrorCount:     0,
		}
		config.ProgressCallback(progress)
	}

	// 读取文件内容
	data, err := os.ReadFile(inputFile)
	if err != nil {
		// 进度更新 - 失败
		if config.ProgressCallback != nil {
			progress := Progress{
				CurrentFile:    inputFile,
				ProcessedFiles: 0,
				TotalFiles:     1,
				Status:         StatusFailed,
				StartTime:      time.Now(),
				ElapsedTime:    time.Since(start),
				ErrorCount:     1,
			}
			config.ProgressCallback(progress)
		}
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// 使用智能转换（自动检测源编码）
	processorResult, err := processor.SmartConvert(data, config.TargetEncoding)
	if err != nil {
		// 进度更新 - 失败
		if config.ProgressCallback != nil {
			progress := Progress{
				CurrentFile:    inputFile,
				ProcessedFiles: 0,
				TotalFiles:     1,
				Status:         StatusFailed,
				StartTime:      time.Now(),
				ElapsedTime:    time.Since(start),
				ErrorCount:     1,
			}
			config.ProgressCallback(progress)
		}
		return nil, fmt.Errorf("failed to convert file %s: %w", inputFile, err)
	}

	// 写入转换后的数据
	err = os.WriteFile(outputFile, processorResult.Data, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to write file %s: %w", outputFile, err)
	}

	// 构建返回结果
	result := &ConvertResult{
		InputFile:           inputFile,
		OutputFile:          outputFile,
		SourceEncoding:      processorResult.SourceEncoding,
		TargetEncoding:      processorResult.TargetEncoding,
		BytesProcessed:      processorResult.BytesProcessed,
		ProcessingTime:      time.Since(start),
		DetectionConfidence: 1.0, // 智能转换的置信度设为1.0
		BackupFile:          "",
		ProcessorResult:     nil,
	}

	// 进度更新 - 完成
	if config.ProgressCallback != nil {
		progress := Progress{
			CurrentFile:    inputFile,
			ProcessedFiles: 1,
			TotalFiles:     1,
			Status:         StatusCompleted,
			StartTime:      time.Now(),
			ElapsedTime:    time.Since(start),
			ProcessedBytes: result.BytesProcessed,
			ErrorCount:     0,
		}
		config.ProgressCallback(progress)
	}

	return result, nil
}

// ConvertFiles 批量转换文件
func ConvertFiles(files []string, options ...Option) (*BatchResult, error) {
	config := applyOptions(options)

	if len(files) == 0 {
		return nil, fmt.Errorf("file list cannot be empty")
	}

	start := time.Now()
	batchResult := &BatchResult{
		TotalFiles:     len(files),
		ProcessingTime: 0,
		Results:        make([]*ConvertResult, 0, len(files)),
		Errors:         make([]FileError, 0),
	}

	// 进度初始化
	if config.ProgressCallback != nil {
		progress := Progress{
			CurrentFile:    "",
			ProcessedFiles: 0,
			TotalFiles:     len(files),
			Status:         StatusStarting,
			StartTime:      start,
			ErrorCount:     0,
		}
		config.ProgressCallback(progress)
	}

	// 并发控制
	semaphore := make(chan struct{}, config.ConcurrencyLimit)
	var wg sync.WaitGroup
	var mutex sync.Mutex

	// 处理每个文件
	for i, file := range files {
		wg.Add(1)
		go func(index int, filePath string) {
			defer wg.Done()
			semaphore <- struct{}{}        // 获取信号量
			defer func() { <-semaphore }() // 释放信号量

			// 应用文件过滤器
			if config.FileFilter != nil && !config.FileFilter(filePath) {
				mutex.Lock()
				batchResult.SkippedFiles++
				batchResult.ProcessedFiles++
				mutex.Unlock()

				if config.ProgressCallback != nil {
					mutex.Lock()
					progress := Progress{
						CurrentFile:    filePath,
						ProcessedFiles: batchResult.ProcessedFiles,
						TotalFiles:     len(files),
						Status:         StatusSkipped,
						StartTime:      start,
						ElapsedTime:    time.Since(start),
						ErrorCount:     len(batchResult.Errors),
					}
					mutex.Unlock()
					config.ProgressCallback(progress)
				}
				return
			}

			// 生成输出文件名
			outputFile := generateOutputFileName(filePath, config.TargetEncoding)

			// 转换单个文件（不使用ConvertFile以避免重复的进度回调）
			processor := encoding.NewSmartProcessor()

			// 读取文件内容
			data, err := os.ReadFile(filePath)
			if err != nil {
				mutex.Lock()
				batchResult.FailedFiles++
				batchResult.Errors = append(batchResult.Errors, FileError{
					File:      filePath,
					Operation: "read",
					Error:     err.Error(),
					Timestamp: time.Now(),
				})
				mutex.Unlock()
				return
			}

			// 使用智能转换（自动检测源编码）
			processorResult, err := processor.SmartConvert(data, config.TargetEncoding)
			if err != nil {
				mutex.Lock()
				batchResult.FailedFiles++
				batchResult.Errors = append(batchResult.Errors, FileError{
					File:      filePath,
					Operation: "convert",
					Error:     err.Error(),
					Timestamp: time.Now(),
				})
				mutex.Unlock()
				return
			}

			// 写入转换后的数据
			err = os.WriteFile(outputFile, processorResult.Data, 0644)
			if err != nil {
				mutex.Lock()
				batchResult.FailedFiles++
				batchResult.Errors = append(batchResult.Errors, FileError{
					File:      filePath,
					Operation: "write",
					Error:     err.Error(),
					Timestamp: time.Now(),
				})
				mutex.Unlock()
				return
			}

			mutex.Lock()
			batchResult.ProcessedFiles++
			batchResult.SuccessfulFiles++
			batchResult.TotalBytes += processorResult.BytesProcessed

			// 构建我们的ConvertResult
			ourResult := &ConvertResult{
				InputFile:           filePath,
				OutputFile:          outputFile,
				SourceEncoding:      processorResult.SourceEncoding,
				TargetEncoding:      processorResult.TargetEncoding,
				BytesProcessed:      processorResult.BytesProcessed,
				ProcessingTime:      processorResult.ConversionTime,
				DetectionConfidence: 1.0, // 智能转换的置信度设为1.0
				BackupFile:          "",
				ProcessorResult:     nil,
			}
			batchResult.Results = append(batchResult.Results, ourResult)

			// 进度回调
			if config.ProgressCallback != nil {
				status := StatusCompleted

				progress := Progress{
					CurrentFile:    filePath,
					ProcessedFiles: batchResult.ProcessedFiles,
					TotalFiles:     len(files),
					Status:         status,
					StartTime:      start,
					ElapsedTime:    time.Since(start),
					ErrorCount:     len(batchResult.Errors),
				}

				if batchResult.ProcessedFiles > 0 {
					remaining := len(files) - batchResult.ProcessedFiles
					avgTime := time.Since(start) / time.Duration(batchResult.ProcessedFiles)
					progress.EstimatedTime = time.Duration(remaining) * avgTime
				}

				config.ProgressCallback(progress)
			}
			mutex.Unlock()
		}(i, file)
	}

	wg.Wait()
	batchResult.ProcessingTime = time.Since(start)

	return batchResult, nil
}

// ConvertDirectory 递归转换目录中的文件
func ConvertDirectory(dirPath string, options ...Option) (*BatchResult, error) {
	config := applyOptions(options)

	// 检查目录是否存在
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("directory does not exist: %s", dirPath)
	}

	// 收集文件列表
	files, err := collectFiles(dirPath, config)
	if err != nil {
		return nil, fmt.Errorf("failed to collect files from directory %s: %w", dirPath, err)
	}

	if len(files) == 0 {
		return &BatchResult{
			TotalFiles:     0,
			ProcessingTime: 0,
			Results:        make([]*ConvertResult, 0),
			Errors:         make([]FileError, 0),
		}, nil
	}

	// 使用ConvertFiles处理收集到的文件
	return ConvertFiles(files, options...)
}

// collectFiles 收集目录中的文件
func collectFiles(dirPath string, config *Config) ([]string, error) {
	var files []string

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 跳过目录
		if info.IsDir() {
			// 如果不是递归模式且当前目录不是根目录，跳过
			if !config.Recursive && path != dirPath {
				return filepath.SkipDir
			}
			return nil
		}

		// 跳过隐藏文件
		if config.SkipHidden && strings.HasPrefix(info.Name(), ".") {
			return nil
		}

		// 检查文件大小限制
		if config.MaxFileSize > 0 && info.Size() > config.MaxFileSize {
			return nil
		}

		// 应用文件过滤器
		if config.FileFilter != nil && !config.FileFilter(path) {
			return nil
		}

		files = append(files, path)
		return nil
	})

	return files, err
}

// generateOutputFileName 生成输出文件名
func generateOutputFileName(inputFile, targetEncoding string) string {
	ext := filepath.Ext(inputFile)
	base := strings.TrimSuffix(inputFile, ext)

	// 如果目标编码是UTF-8，保持原文件名
	if strings.ToUpper(targetEncoding) == "UTF-8" {
		return inputFile
	}

	// 其他编码添加后缀
	return fmt.Sprintf("%s_%s%s", base, strings.ToLower(targetEncoding), ext)
}
