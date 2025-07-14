package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	encoding "github.com/mirbf/encoding-processor"
	"github.com/mirbf/ConvertContent2UTF8"
)

func main() {
	fmt.Println("批量转换 /Users/apple/Desktop/test/ok 目录下的文件")
	fmt.Println("=================================================")

	sourceDir := "/Users/apple/Desktop/test/ok"
	utf8Dir := "/Users/apple/Desktop/test/UTF8"
	alreadyUTF8Dir := "/Users/apple/Desktop/test/has"
	errorDir := "/Users/apple/Desktop/test/error"

	// 检查源目录是否存在
	if _, err := os.Stat(sourceDir); os.IsNotExist(err) {
		log.Fatalf("源目录不存在: %s", sourceDir)
	}

	// 创建目标目录
	dirs := []string{utf8Dir, alreadyUTF8Dir, errorDir}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Fatalf("创建目录失败 %s: %v", dir, err)
		}
	}

	fmt.Printf("源目录: %s\n", sourceDir)
	fmt.Printf("转换成功目录: %s\n", utf8Dir)
	fmt.Printf("已是UTF8目录: %s\n", alreadyUTF8Dir)
	fmt.Printf("转换失败目录: %s\n", errorDir)
	fmt.Println()

	// 进度回调函数
	progressCallback := func(p ConvertContent2UTF8.Progress) {
		percentage := float64(p.ProcessedFiles) / float64(p.TotalFiles) * 100

		switch p.Status {
		case ConvertContent2UTF8.StatusStarting:
			fmt.Printf("🚀 开始处理 %d 个文件...\n", p.TotalFiles)
		case ConvertContent2UTF8.StatusProcessing:
			fmt.Printf("⏳ [%.1f%%] 正在处理: %s\n", percentage, filepath.Base(p.CurrentFile))
		case ConvertContent2UTF8.StatusCompleted:
			fmt.Printf("✅ [%.1f%%] 完成: %s\n", percentage, filepath.Base(p.CurrentFile))
		case ConvertContent2UTF8.StatusFailed:
			fmt.Printf("❌ [%.1f%%] 失败: %s\n", percentage, filepath.Base(p.CurrentFile))
		case ConvertContent2UTF8.StatusSkipped:
			fmt.Printf("⏭️  [%.1f%%] 跳过: %s\n", percentage, filepath.Base(p.CurrentFile))
		}

		if p.ProcessedFiles == p.TotalFiles {
			fmt.Printf("\n📊 处理完成统计:\n")
			fmt.Printf("   总计耗时: %v\n", p.ElapsedTime)
			fmt.Printf("   错误数量: %d\n", p.ErrorCount)
		}
	}

	// 收集所有文件并检查编码
	start := time.Now()
	files, err := collectAllFiles(sourceDir)
	if err != nil {
		log.Fatalf("收集文件失败: %v", err)
	}

	if len(files) == 0 {
		fmt.Println("源目录中没有找到任何文件")
		return
	}

	fmt.Printf("找到 %d 个文件，开始处理...\n\n", len(files))

	// 分类处理文件
	alreadyUTF8Files := []*FileInfo{}
	needConvertFiles := []*FileInfo{}
	errorFiles := []*FileInfo{}

	// 检查每个文件的编码
	for i, fileInfo := range files {
		percentage := float64(i+1) / float64(len(files)) * 100
		fmt.Printf("[%.1f%%] 检查文件编码: %s\n", percentage, filepath.Base(fileInfo.Path))

		// 检测文件编码（降低置信度要求）
		config := encoding.GetDefaultDetectorConfig()
		config.MinConfidence = 0.3 // 降低置信度要求到0.3
		detector := encoding.NewDetector(config)
		detectResult, err := detector.DetectFileEncoding(fileInfo.Path)
		if err != nil {
			fmt.Printf("❌ 编码检测失败: %s\n", err)
			fileInfo.Error = err.Error()
			errorFiles = append(errorFiles, fileInfo)
			continue
		}

		fileInfo.DetectedEncoding = detectResult.Encoding
		fileInfo.Confidence = detectResult.Confidence

		// 如果已经是UTF-8，直接移动到has目录
		if strings.ToUpper(detectResult.Encoding) == "UTF-8" {
			fmt.Printf("ℹ️  已是UTF-8编码，移动到has目录\n")
			alreadyUTF8Files = append(alreadyUTF8Files, fileInfo)
		} else {
			fmt.Printf("🔄 需要转换编码: %s -> UTF-8\n", detectResult.Encoding)
			needConvertFiles = append(needConvertFiles, fileInfo)
		}
	}

	fmt.Printf("\n编码检查完成:\n")
	fmt.Printf("  已是UTF-8: %d 个文件\n", len(alreadyUTF8Files))
	fmt.Printf("  需要转换: %d 个文件\n", len(needConvertFiles))
	fmt.Printf("  检测失败: %d 个文件\n\n", len(errorFiles))

	// 移动已经是UTF-8的文件
	if len(alreadyUTF8Files) > 0 {
		fmt.Println("移动已是UTF-8的文件到has目录...")
		moveFiles(alreadyUTF8Files, alreadyUTF8Dir, "已是UTF-8")
	}

	// 移动检测失败的文件
	if len(errorFiles) > 0 {
		fmt.Println("移动编码检测失败的文件到error目录...")
		moveFiles(errorFiles, errorDir, "检测失败")
	}

	// 转换需要转换的文件
	if len(needConvertFiles) > 0 {
		fmt.Printf("开始转换 %d 个文件...\n", len(needConvertFiles))
		successFiles, failFiles := convertFiles(needConvertFiles, progressCallback)

		// 移动转换成功的文件
		if len(successFiles) > 0 {
			fmt.Println("移动转换成功的文件到UTF8目录...")
			moveFiles(successFiles, utf8Dir, "转换成功")
		}

		// 移动转换失败的文件
		if len(failFiles) > 0 {
			fmt.Println("移动转换失败的文件到error目录...")
			moveFiles(failFiles, errorDir, "转换失败")
		}
	}

	fmt.Printf("\n📈 最终结果:\n")
	fmt.Printf("   总文件数: %d\n", len(files))
	fmt.Printf("   已是UTF-8: %d (移动到has目录)\n", len(alreadyUTF8Files))
	fmt.Printf("   转换成功: %d (移动到UTF8目录)\n", countFilesByStatus(files, "success"))
	fmt.Printf("   转换失败: %d (移动到error目录)\n", countFilesByStatus(files, "error"))
	fmt.Printf("   总耗时: %v\n", time.Since(start))

	fmt.Printf("\n🎉 所有操作完成！\n")
}

// FileInfo 文件信息结构
type FileInfo struct {
	Path             string
	DetectedEncoding string
	Confidence       float64
	Error            string
	Status           string // "success", "error", "utf8"
}

// collectAllFiles 收集目录中的所有文件
func collectAllFiles(dirPath string) ([]*FileInfo, error) {
	var files []*FileInfo

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 跳过目录
		if info.IsDir() {
			return nil
		}

		// 跳过隐藏文件
		if strings.HasPrefix(info.Name(), ".") {
			return nil
		}

		// 文件类型过滤
		ext := strings.ToLower(filepath.Ext(path))
		// 去除时间戳后缀，如 .txt.20250713231607 -> .txt
		if ext != "" && strings.Contains(ext, ".") {
			// 如果扩展名包含数字，可能是时间戳，尝试提取真正的扩展名
			parts := strings.Split(filepath.Base(path), ".")
			if len(parts) >= 2 {
				// 找到最后一个非数字的扩展名
				for i := len(parts) - 2; i >= 0; i-- {
					if !isAllDigits(parts[i]) {
						ext = "." + parts[i]
						break
					}
				}
			}
		}

		textExts := []string{".txt", ".log", ".md", ".csv", ".json", ".xml", ".html", ".css", ".js", ".py", ".go", ".java", ".c", ".cpp", ".h"}
		isTextFile := false
		for _, textExt := range textExts {
			if ext == textExt {
				isTextFile = true
				break
			}
		}
		// 无扩展名的文件也可能是文本文件
		if ext == "" {
			isTextFile = true
		}

		if !isTextFile {
			return nil
		}

		files = append(files, &FileInfo{
			Path: path,
		})
		return nil
	})

	return files, err
}

// convertFiles 转换文件
func convertFiles(files []*FileInfo, progressCallback func(ConvertContent2UTF8.Progress)) ([]*FileInfo, []*FileInfo) {
	var successFiles []*FileInfo
	var failFiles []*FileInfo

	for i, fileInfo := range files {
		percentage := float64(i+1) / float64(len(files)) * 100
		fmt.Printf("[%.1f%%] 转换文件: %s\n", percentage, filepath.Base(fileInfo.Path))

		// 创建临时输出文件
		tempOutputFile := fileInfo.Path + ".utf8.tmp"

		// 使用encoding-processor转换文件
		fileProcessor := encoding.NewDefaultFile()
		processOptions := &encoding.FileProcessOptions{
			TargetEncoding:    "UTF-8",
			MinConfidence:     0.8,
			CreateBackup:      false,
			OverwriteExisting: true,
			DryRun:            false,
			PreserveMode:      true,
			PreserveTime:      true,
		}

		result, err := fileProcessor.ProcessFile(fileInfo.Path, tempOutputFile, processOptions)
		if err != nil {
			fmt.Printf("❌ 转换失败: %s\n", err)
			fileInfo.Error = err.Error()
			fileInfo.Status = "error"
			failFiles = append(failFiles, fileInfo)
			// 清理临时文件
			os.Remove(tempOutputFile)
		} else {
			fmt.Printf("✅ 转换成功: %s -> UTF-8\n", result.SourceEncoding)
			// 用转换后的文件替换原文件
			err = os.Rename(tempOutputFile, fileInfo.Path)
			if err != nil {
				fmt.Printf("❌ 替换文件失败: %s\n", err)
				fileInfo.Error = err.Error()
				fileInfo.Status = "error"
				failFiles = append(failFiles, fileInfo)
				os.Remove(tempOutputFile)
			} else {
				fileInfo.Status = "success"
				successFiles = append(successFiles, fileInfo)
			}
		}
	}

	return successFiles, failFiles
}

// moveFiles 移动文件到指定目录
func moveFiles(files []*FileInfo, targetDir, operation string) {
	successCount := 0
	failCount := 0

	for _, fileInfo := range files {
		// 构建目标文件路径
		fileName := filepath.Base(fileInfo.Path)
		targetPath := filepath.Join(targetDir, fileName)

		// 如果目标文件已存在，添加时间戳后缀避免冲突
		if _, err := os.Stat(targetPath); err == nil {
			ext := filepath.Ext(fileName)
			nameWithoutExt := strings.TrimSuffix(fileName, ext)
			timestamp := time.Now().Format("20060102_150405")
			targetPath = filepath.Join(targetDir, fmt.Sprintf("%s_%s%s", nameWithoutExt, timestamp, ext))
		}

		// 移动文件
		err := os.Rename(fileInfo.Path, targetPath)
		if err != nil {
			fmt.Printf("❌ %s移动失败: %s -> %s (%v)\n", operation, fileName, filepath.Base(targetPath), err)
			failCount++
		} else {
			fmt.Printf("✅ %s移动成功: %s -> %s\n", operation, fileName, filepath.Base(targetPath))
			successCount++
		}
	}

	fmt.Printf("%s移动结果: 成功 %d, 失败 %d\n\n", operation, successCount, failCount)
}

// isAllDigits 检查字符串是否全为数字
func isAllDigits(s string) bool {
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return len(s) > 0
}

// countFilesByStatus 统计指定状态的文件数量
func countFilesByStatus(files []*FileInfo, status string) int {
	count := 0
	for _, file := range files {
		if file.Status == status {
			count++
		}
	}
	return count
}
