package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mirbf/ConvertContent2UTF8"
)

func main() {
	fmt.Println("ConvertContent2UTF8 库使用示例")
	fmt.Println("================================")

	// 创建示例目录和文件
	if err := setupExampleFiles(); err != nil {
		log.Fatalf("创建示例文件失败: %v", err)
	}
	defer cleanupExampleFiles()

	// 示例1: 单文件转换
	fmt.Println("\n1. 单文件转换示例")
	singleFileExample()

	// 示例2: 批量文件转换
	fmt.Println("\n2. 批量文件转换示例")
	batchFilesExample()

	// 示例3: 目录递归转换
	fmt.Println("\n3. 目录递归转换示例")
	directoryExample()

	// 示例4: 带进度回调的转换
	fmt.Println("\n4. 带进度回调的转换示例")
	progressExample()

	// 示例5: 自定义选项转换
	fmt.Println("\n5. 自定义选项转换示例")
	customOptionsExample()
}

// 单文件转换示例
func singleFileExample() {
	inputFile := filepath.Join("temp_examples", "single_test.txt")
	outputFile := filepath.Join("temp_examples", "single_output.txt")

	fmt.Printf("转换文件: %s -> %s\n", inputFile, outputFile)

	result, err := ConvertContent2UTF8.ConvertFile(inputFile, outputFile,
		ConvertContent2UTF8.WithOverwrite(true),
		ConvertContent2UTF8.WithBackup(true),
	)

	if err != nil {
		fmt.Printf("转换失败: %v\n", err)
		return
	}

	fmt.Printf("转换成功!\n")
	fmt.Printf("  源编码: %s\n", result.SourceEncoding)
	fmt.Printf("  目标编码: %s\n", result.TargetEncoding)
	fmt.Printf("  处理字节数: %d\n", result.BytesProcessed)
	fmt.Printf("  处理时间: %v\n", result.ProcessingTime)
	fmt.Printf("  检测置信度: %.2f\n", result.DetectionConfidence)
	if result.BackupFile != "" {
		fmt.Printf("  备份文件: %s\n", result.BackupFile)
	}
}

// 批量文件转换示例
func batchFilesExample() {
	files := []string{
		filepath.Join("temp_examples", "batch1.txt"),
		filepath.Join("temp_examples", "batch2.txt"),
		filepath.Join("temp_examples", "batch3.txt"),
	}

	fmt.Printf("批量转换 %d 个文件\n", len(files))

	result, err := ConvertContent2UTF8.ConvertFiles(files,
		ConvertContent2UTF8.WithOverwrite(true),
		ConvertContent2UTF8.WithConcurrency(2),
		ConvertContent2UTF8.WithBackup(false),
	)

	if err != nil {
		fmt.Printf("批量转换失败: %v\n", err)
		return
	}

	fmt.Printf("批量转换完成!\n")
	fmt.Printf("  总文件数: %d\n", result.TotalFiles)
	fmt.Printf("  处理文件数: %d\n", result.ProcessedFiles)
	fmt.Printf("  成功文件数: %d\n", result.SuccessfulFiles)
	fmt.Printf("  失败文件数: %d\n", result.FailedFiles)
	fmt.Printf("  跳过文件数: %d\n", result.SkippedFiles)
	fmt.Printf("  总处理字节数: %d\n", result.TotalBytes)
	fmt.Printf("  总处理时间: %v\n", result.ProcessingTime)

	if len(result.Errors) > 0 {
		fmt.Printf("  错误信息:\n")
		for _, fileErr := range result.Errors {
			fmt.Printf("    %s: %s\n", fileErr.File, fileErr.Error)
		}
	}
}

// 目录递归转换示例
func directoryExample() {
	dirPath := filepath.Join("temp_examples", "subdir")

	fmt.Printf("递归转换目录: %s\n", dirPath)

	result, err := ConvertContent2UTF8.ConvertDirectory(dirPath,
		ConvertContent2UTF8.WithRecursive(true),
		ConvertContent2UTF8.WithOverwrite(true),
		ConvertContent2UTF8.WithFileFilter(func(filename string) bool {
			// 只处理.txt文件
			return strings.HasSuffix(strings.ToLower(filename), ".txt")
		}),
	)

	if err != nil {
		fmt.Printf("目录转换失败: %v\n", err)
		return
	}

	fmt.Printf("目录转换完成!\n")
	fmt.Printf("  找到文件数: %d\n", result.TotalFiles)
	fmt.Printf("  成功转换: %d\n", result.SuccessfulFiles)
	fmt.Printf("  处理时间: %v\n", result.ProcessingTime)
}

// 带进度回调的转换示例
func progressExample() {
	files := []string{
		filepath.Join("temp_examples", "progress1.txt"),
		filepath.Join("temp_examples", "progress2.txt"),
		filepath.Join("temp_examples", "progress3.txt"),
		filepath.Join("temp_examples", "progress4.txt"),
	}

	fmt.Printf("带进度显示的批量转换 %d 个文件\n", len(files))

	// 创建进度回调函数
	progressCallback := func(p ConvertContent2UTF8.Progress) {
		percentage := float64(p.ProcessedFiles) / float64(p.TotalFiles) * 100

		switch p.Status {
		case ConvertContent2UTF8.StatusStarting:
			fmt.Printf("开始处理...\n")
		case ConvertContent2UTF8.StatusProcessing:
			fmt.Printf("处理中 [%.1f%%] %s\n", percentage, filepath.Base(p.CurrentFile))
		case ConvertContent2UTF8.StatusCompleted:
			fmt.Printf("完成 [%.1f%%] %s\n", percentage, filepath.Base(p.CurrentFile))
		case ConvertContent2UTF8.StatusFailed:
			fmt.Printf("失败 [%.1f%%] %s\n", percentage, filepath.Base(p.CurrentFile))
		case ConvertContent2UTF8.StatusSkipped:
			fmt.Printf("跳过 [%.1f%%] %s\n", percentage, filepath.Base(p.CurrentFile))
		}

		if p.ProcessedFiles == p.TotalFiles {
			fmt.Printf("总计耗时: %v, 错误数: %d\n", p.ElapsedTime, p.ErrorCount)
		}
	}

	result, err := ConvertContent2UTF8.ConvertFiles(files,
		ConvertContent2UTF8.WithProgress(progressCallback),
		ConvertContent2UTF8.WithOverwrite(true),
		ConvertContent2UTF8.WithConcurrency(1), // 串行处理以便观察进度
	)

	if err != nil {
		fmt.Printf("进度转换失败: %v\n", err)
		return
	}

	fmt.Printf("进度转换完成! 成功: %d, 失败: %d\n", result.SuccessfulFiles, result.FailedFiles)
}

// 自定义选项转换示例
func customOptionsExample() {
	inputFile := filepath.Join("temp_examples", "custom_test.txt")
	outputFile := filepath.Join("temp_examples", "custom_output.txt")

	fmt.Printf("自定义选项转换: %s\n", inputFile)

	result, err := ConvertContent2UTF8.ConvertFile(inputFile, outputFile,
		ConvertContent2UTF8.WithTargetEncoding("UTF-8"),
		ConvertContent2UTF8.WithMinConfidence(0.7),
		ConvertContent2UTF8.WithBackup(true),
		ConvertContent2UTF8.WithOverwrite(true),
		ConvertContent2UTF8.WithDryRun(false), // 如果设置为true，只检测不转换
	)

	if err != nil {
		fmt.Printf("自定义转换失败: %v\n", err)
		return
	}

	fmt.Printf("自定义转换完成!\n")
	fmt.Printf("  处理时间: %v\n", result.ProcessingTime)
	fmt.Printf("  检测置信度: %.2f (阈值: 0.7)\n", result.DetectionConfidence)
}

// 创建示例文件
func setupExampleFiles() error {
	baseDir := "temp_examples"
	subDir := filepath.Join(baseDir, "subdir")

	// 创建目录
	if err := os.MkdirAll(subDir, 0755); err != nil {
		return err
	}

	// 示例文件内容
	files := map[string]string{
		filepath.Join(baseDir, "single_test.txt"): "这是单文件转换测试\nHello World\n测试中文内容",
		filepath.Join(baseDir, "batch1.txt"):      "批量转换测试文件1\n包含中文和English",
		filepath.Join(baseDir, "batch2.txt"):      "批量转换测试文件2\n更多测试内容",
		filepath.Join(baseDir, "batch3.txt"):      "批量转换测试文件3\n最后一个测试文件",
		filepath.Join(subDir, "sub1.txt"):         "子目录文件1\n测试递归转换",
		filepath.Join(subDir, "sub2.txt"):         "子目录文件2\n递归处理测试",
		filepath.Join(baseDir, "progress1.txt"):   "进度测试文件1",
		filepath.Join(baseDir, "progress2.txt"):   "进度测试文件2",
		filepath.Join(baseDir, "progress3.txt"):   "进度测试文件3",
		filepath.Join(baseDir, "progress4.txt"):   "进度测试文件4",
		filepath.Join(baseDir, "custom_test.txt"): "自定义选项测试\n包含各种配置",
	}

	// 创建文件
	for filepath, content := range files {
		if err := os.WriteFile(filepath, []byte(content), 0644); err != nil {
			return err
		}
	}

	return nil
}

// 清理示例文件
func cleanupExampleFiles() {
	time.Sleep(100 * time.Millisecond) // 确保文件操作完成
	os.RemoveAll("temp_examples")
	fmt.Println("\n清理完成: 已删除临时示例文件")
}
