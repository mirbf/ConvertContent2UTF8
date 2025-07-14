package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	converter "github.com/mirbf/ConvertContent2UTF8"
)

func main() {
	var (
		inputPath    = flag.String("input", "", "输入文件或目录路径")
		outputPath   = flag.String("output", "", "输出路径（可选，默认覆盖原文件）")
		recursive    = flag.Bool("recursive", true, "递归处理子目录")
		createBackup = flag.Bool("backup", true, "创建备份文件")
		dryRun       = flag.Bool("dry-run", false, "仅检测编码，不实际转换")
		verbose      = flag.Bool("verbose", false, "显示详细信息")
		concurrent   = flag.Int("concurrent", 10, "并发处理数量")
		confidence   = flag.Float64("confidence", 0.8, "编码检测置信度阈值")
	)

	flag.Parse()

	if *inputPath == "" {
		fmt.Println("用法: go run main.go -input <文件或目录路径>")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// 检查输入路径是否存在
	stat, err := os.Stat(*inputPath)
	if err != nil {
		log.Fatalf("错误: 无法访问路径 %s: %v", *inputPath, err)
	}

	// 配置选项
	options := []converter.Option{
		converter.WithTargetEncoding("UTF-8"),
		converter.WithMinConfidence(*confidence),
		converter.WithRecursive(*recursive),
		converter.WithBackup(*createBackup),
		converter.WithDryRun(*dryRun),
		converter.WithConcurrency(*concurrent),
		converter.WithOverwrite(true),
	}

	// 添加进度回调
	if *verbose {
		options = append(options, converter.WithProgress(func(progress converter.Progress) {
			switch progress.Status {
			case converter.StatusStarting:
				fmt.Printf("开始处理... 总文件数: %d\n", progress.TotalFiles)
			case converter.StatusProcessing:
				fmt.Printf("处理中: %s\n", progress.CurrentFile)
			case converter.StatusCompleted:
				fmt.Printf("✓ 完成: %s\n", progress.CurrentFile)
			case converter.StatusFailed:
				fmt.Printf("✗ 失败: %s\n", progress.CurrentFile)
			case converter.StatusSkipped:
				fmt.Printf("- 跳过: %s\n", progress.CurrentFile)
			}

			if progress.ProcessedFiles > 0 {
				percentage := float64(progress.ProcessedFiles) / float64(progress.TotalFiles) * 100
				fmt.Printf("进度: %d/%d (%.1f%%) - 错误: %d\n",
					progress.ProcessedFiles, progress.TotalFiles, percentage, progress.ErrorCount)
			}
		}))
	}

	// 添加文件过滤器，只处理文本文件
	options = append(options, converter.WithFileFilter(func(filename string) bool {
		// 常见文本文件扩展名
		textExts := []string{".txt", ".md", ".csv", ".log", ".json", ".xml", ".html", ".htm", ".js", ".css", ".py", ".go", ".java", ".c", ".cpp", ".h"}
		dotIndex := strings.LastIndex(filename, ".")
		if dotIndex == -1 {
			return false // 没有扩展名的文件不处理
		}
		ext := strings.ToLower(filename[dotIndex:])
		for _, textExt := range textExts {
			if ext == textExt {
				return true
			}
		}
		return false
	}))

	var result *converter.BatchResult

	if stat.IsDir() {
		// 处理目录
		fmt.Printf("处理目录: %s\n", *inputPath)
		result, err = converter.ConvertDirectory(*inputPath, options...)
	} else {
		// 处理单个文件
		fmt.Printf("处理文件: %s\n", *inputPath)
		outputFile := *outputPath
		if outputFile == "" {
			outputFile = *inputPath // 覆盖原文件
		}

		convertResult, err := converter.ConvertFile(*inputPath, outputFile, options...)
		if err != nil {
			log.Fatalf("转换失败: %v", err)
		}

		// 创建一个BatchResult包装单个文件结果
		result = &converter.BatchResult{
			TotalFiles:      1,
			ProcessedFiles:  1,
			SuccessfulFiles: 1,
			Results:         []*converter.ConvertResult{convertResult},
			ProcessingTime:  convertResult.ProcessingTime,
			TotalBytes:      convertResult.BytesProcessed,
		}
	}

	if err != nil {
		log.Fatalf("处理失败: %v", err)
	}

	// 打印结果摘要
	fmt.Println("\n========== 处理完成 ==========")
	fmt.Printf("总文件数: %d\n", result.TotalFiles)
	fmt.Printf("成功转换: %d\n", result.SuccessfulFiles)
	fmt.Printf("失败: %d\n", result.FailedFiles)
	fmt.Printf("跳过: %d\n", result.SkippedFiles)
	fmt.Printf("处理时间: %v\n", result.ProcessingTime)
	fmt.Printf("处理字节数: %d\n", result.TotalBytes)

	if len(result.Errors) > 0 {
		fmt.Println("\n错误详情:")
		for _, fileErr := range result.Errors {
			fmt.Printf("- %s: %s\n", fileErr.File, fileErr.Error)
		}
	}

	if *verbose && len(result.Results) > 0 {
		fmt.Println("\n详细结果:")
		for _, res := range result.Results {
			fmt.Printf("文件: %s\n", res.InputFile)
			fmt.Printf("  源编码: %s\n", res.SourceEncoding)
			fmt.Printf("  目标编码: %s\n", res.TargetEncoding)
			fmt.Printf("  置信度: %.2f\n", res.DetectionConfidence)
			fmt.Printf("  处理字节: %d\n", res.BytesProcessed)
			fmt.Printf("  处理时间: %v\n", res.ProcessingTime)
			if res.BackupFile != "" {
				fmt.Printf("  备份文件: %s\n", res.BackupFile)
			}
			fmt.Println()
		}
	}
}
