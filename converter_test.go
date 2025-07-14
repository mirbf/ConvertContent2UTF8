package convertcontent2utf8

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestConvertFile(t *testing.T) {
	// 创建临时测试目录
	testDir := t.TempDir()

	// 创建测试文件
	inputFile := filepath.Join(testDir, "test_input.txt")
	outputFile := filepath.Join(testDir, "test_output.txt")

	// 写入测试内容（UTF-8编码）
	testContent := "这是一个测试文件\nHello World\n测试中文内容"
	err := os.WriteFile(inputFile, []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	t.Run("基本文件转换", func(t *testing.T) {
		result, err := ConvertFile(inputFile, outputFile)
		if err != nil {
			t.Fatalf("ConvertFile failed: %v", err)
		}

		// 验证结果
		if result.InputFile != inputFile {
			t.Errorf("Expected InputFile %s, got %s", inputFile, result.InputFile)
		}
		if result.OutputFile != outputFile {
			t.Errorf("Expected OutputFile %s, got %s", outputFile, result.OutputFile)
		}
		if result.TargetEncoding != "UTF-8" {
			t.Errorf("Expected TargetEncoding UTF-8, got %s", result.TargetEncoding)
		}
		if result.BytesProcessed <= 0 {
			t.Errorf("Expected BytesProcessed > 0, got %d", result.BytesProcessed)
		}

		// 验证输出文件存在
		if _, err := os.Stat(outputFile); os.IsNotExist(err) {
			t.Errorf("Output file was not created: %s", outputFile)
		}
	})

	t.Run("带进度回调的转换", func(t *testing.T) {
		outputFile2 := filepath.Join(testDir, "test_output2.txt")
		progressCalled := false

		result, err := ConvertFile(inputFile, outputFile2, WithProgress(func(p Progress) {
			progressCalled = true
			if p.TotalFiles != 1 {
				t.Errorf("Expected TotalFiles 1, got %d", p.TotalFiles)
			}
		}))

		if err != nil {
			t.Fatalf("ConvertFile with progress failed: %v", err)
		}
		if !progressCalled {
			t.Error("Progress callback was not called")
		}
		if result == nil {
			t.Error("Result should not be nil")
		}
	})

	t.Run("自定义选项", func(t *testing.T) {
		outputFile3 := filepath.Join(testDir, "test_output3.txt")

		result, err := ConvertFile(inputFile, outputFile3,
			WithTargetEncoding("UTF-8"),
			WithBackup(false),
			WithOverwrite(true),
		)

		if err != nil {
			t.Fatalf("ConvertFile with options failed: %v", err)
		}
		if result.BackupFile != "" {
			t.Error("Backup file should be empty when backup is disabled")
		}
	})

	t.Run("错误情况", func(t *testing.T) {
		// 空输入文件路径
		_, err := ConvertFile("", outputFile)
		if err == nil {
			t.Error("Expected error for empty input file")
		}

		// 空输出文件路径
		_, err = ConvertFile(inputFile, "")
		if err == nil {
			t.Error("Expected error for empty output file")
		}

		// 不存在的输入文件
		_, err = ConvertFile("/nonexistent/file.txt", outputFile)
		if err == nil {
			t.Error("Expected error for nonexistent input file")
		}
	})
}

func TestConvertFiles(t *testing.T) {
	testDir := t.TempDir()

	// 创建多个测试文件
	files := []string{
		filepath.Join(testDir, "test1.txt"),
		filepath.Join(testDir, "test2.txt"),
		filepath.Join(testDir, "test3.txt"),
	}

	for i, file := range files {
		content := "测试文件内容 " + string(rune('1'+i))
		err := os.WriteFile(file, []byte(content), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file %s: %v", file, err)
		}
	}

	t.Run("批量文件转换", func(t *testing.T) {
		result, err := ConvertFiles(files, WithOverwrite(true)) // 允许覆盖，因为输入输出是同一个文件
		if err != nil {
			t.Fatalf("ConvertFiles failed: %v", err)
		}

		// 打印错误信息用于调试
		if len(result.Errors) > 0 {
			t.Logf("Found %d errors:", len(result.Errors))
			for _, fileErr := range result.Errors {
				t.Logf("File: %s, Error: %s", fileErr.File, fileErr.Error)
			}
		}

		if result.TotalFiles != len(files) {
			t.Errorf("Expected TotalFiles %d, got %d", len(files), result.TotalFiles)
		}
		if result.ProcessedFiles != len(files) {
			t.Errorf("Expected ProcessedFiles %d, got %d", len(files), result.ProcessedFiles)
		}
		if result.SuccessfulFiles != len(files) {
			t.Errorf("Expected SuccessfulFiles %d, got %d", len(files), result.SuccessfulFiles)
		}
		if result.FailedFiles != 0 {
			t.Errorf("Expected FailedFiles 0, got %d", result.FailedFiles)
		}
	})

	t.Run("带进度回调的批量转换", func(t *testing.T) {
		progressUpdates := 0

		result, err := ConvertFiles(files,
			WithOverwrite(true), // 允许覆盖
			WithProgress(func(p Progress) {
				progressUpdates++
				if p.TotalFiles != len(files) {
					t.Errorf("Expected TotalFiles %d, got %d", len(files), p.TotalFiles)
				}
			}))

		if err != nil {
			t.Fatalf("ConvertFiles with progress failed: %v", err)
		}
		if progressUpdates == 0 {
			t.Error("Progress callback was not called")
		}
		if result.ProcessedFiles != len(files) {
			t.Errorf("Expected ProcessedFiles %d, got %d", len(files), result.ProcessedFiles)
		}
	})

	t.Run("并发限制", func(t *testing.T) {
		result, err := ConvertFiles(files, WithConcurrency(2), WithOverwrite(true))
		if err != nil {
			t.Fatalf("ConvertFiles with concurrency failed: %v", err)
		}
		if result.SuccessfulFiles != len(files) {
			t.Errorf("Expected SuccessfulFiles %d, got %d", len(files), result.SuccessfulFiles)
		}
	})

	t.Run("文件过滤器", func(t *testing.T) {
		// 只处理test1.txt文件
		result, err := ConvertFiles(files,
			WithOverwrite(true),
			WithFileFilter(func(filename string) bool {
				return strings.Contains(filename, "test1.txt")
			}))

		if err != nil {
			t.Fatalf("ConvertFiles with filter failed: %v", err)
		}

		t.Logf("Result: TotalFiles=%d, ProcessedFiles=%d, SuccessfulFiles=%d, SkippedFiles=%d",
			result.TotalFiles, result.ProcessedFiles, result.SuccessfulFiles, result.SkippedFiles)

		// 由于文件过滤是在goroutine中进行的，不会改变TotalFiles
		// 只有SkippedFiles和SuccessfulFiles会变化
		if result.ProcessedFiles != len(files) {
			t.Errorf("Expected ProcessedFiles %d, got %d", len(files), result.ProcessedFiles)
		}
		if result.SkippedFiles != 2 {
			t.Errorf("Expected SkippedFiles 2, got %d", result.SkippedFiles)
		}
		if result.SuccessfulFiles != 1 {
			t.Errorf("Expected SuccessfulFiles 1, got %d", result.SuccessfulFiles)
		}
	})

	t.Run("空文件列表", func(t *testing.T) {
		_, err := ConvertFiles([]string{})
		if err == nil {
			t.Error("Expected error for empty file list")
		}
	})
}

func TestConvertDirectory(t *testing.T) {
	testDir := t.TempDir()

	// 创建子目录结构
	subDir := filepath.Join(testDir, "subdir")
	err := os.MkdirAll(subDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create subdirectory: %v", err)
	}

	// 创建测试文件
	files := []string{
		filepath.Join(testDir, "root1.txt"),
		filepath.Join(testDir, "root2.txt"),
		filepath.Join(subDir, "sub1.txt"),
		filepath.Join(subDir, "sub2.txt"),
	}

	for i, file := range files {
		content := "目录测试文件 " + string(rune('1'+i))
		err := os.WriteFile(file, []byte(content), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file %s: %v", file, err)
		}
	}

	// 创建非txt文件（应该被过滤掉）
	nonTxtFile := filepath.Join(testDir, "test.log")
	err = os.WriteFile(nonTxtFile, []byte("log content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create non-txt file: %v", err)
	}

	t.Run("非递归目录转换", func(t *testing.T) {
		result, err := ConvertDirectory(testDir, WithOverwrite(true))
		if err != nil {
			t.Fatalf("ConvertDirectory failed: %v", err)
		}

		// 非递归模式应该只处理根目录的txt文件
		if result.SuccessfulFiles != 2 {
			t.Errorf("Expected SuccessfulFiles 2, got %d", result.SuccessfulFiles)
		}
	})

	t.Run("递归目录转换", func(t *testing.T) {
		result, err := ConvertDirectory(testDir, WithRecursive(true), WithOverwrite(true))
		if err != nil {
			t.Fatalf("ConvertDirectory recursive failed: %v", err)
		}

		// 递归模式应该处理所有txt文件
		if result.SuccessfulFiles != 4 {
			t.Errorf("Expected SuccessfulFiles 4, got %d", result.SuccessfulFiles)
		}
	})

	t.Run("自定义文件过滤器", func(t *testing.T) {
		// 处理所有文件（包括.log），但排除备份文件
		result, err := ConvertDirectory(testDir,
			WithRecursive(true),
			WithOverwrite(true),
			WithFileFilter(func(filename string) bool {
				// 排除备份文件
				if strings.HasSuffix(filename, ".bak") {
					return false
				}
				return true // 接受所有其他文件
			}),
		)

		if err != nil {
			t.Fatalf("ConvertDirectory with custom filter failed: %v", err)
		}

		t.Logf("Custom filter result: TotalFiles=%d", result.TotalFiles)
		// 应该包括原始的5个文件（4个txt + 1个log）
		// 可能还有临时生成的其他文件，所以我们只检查至少有5个
		if result.TotalFiles < 5 {
			t.Errorf("Expected TotalFiles >= 5, got %d", result.TotalFiles)
		}
	})

	t.Run("不存在的目录", func(t *testing.T) {
		_, err := ConvertDirectory("/nonexistent/directory")
		if err == nil {
			t.Error("Expected error for nonexistent directory")
		}
	})

	t.Run("空目录", func(t *testing.T) {
		emptyDir := filepath.Join(testDir, "empty")
		err := os.MkdirAll(emptyDir, 0755)
		if err != nil {
			t.Fatalf("Failed to create empty directory: %v", err)
		}

		result, err := ConvertDirectory(emptyDir)
		if err != nil {
			t.Fatalf("ConvertDirectory on empty dir failed: %v", err)
		}
		if result.TotalFiles != 0 {
			t.Errorf("Expected TotalFiles 0, got %d", result.TotalFiles)
		}
	})
}

func TestProgressStatus(t *testing.T) {
	testCases := []struct {
		status   ProgressStatus
		expected string
	}{
		{StatusStarting, "starting"},
		{StatusProcessing, "processing"},
		{StatusCompleted, "completed"},
		{StatusFailed, "failed"},
		{StatusSkipped, "skipped"},
	}

	for _, tc := range testCases {
		if string(tc.status) != tc.expected {
			t.Errorf("Expected status %s, got %s", tc.expected, string(tc.status))
		}
	}
}

func TestGenerateOutputFileName(t *testing.T) {
	testCases := []struct {
		input    string
		encoding string
		expected string
	}{
		{"test.txt", "UTF-8", "test.txt"},
		{"test.txt", "utf-8", "test.txt"},
		{"test.txt", "GBK", "test_gbk.txt"},
		{"/path/to/file.txt", "BIG5", "/path/to/file_big5.txt"},
		{"file_without_ext", "UTF-8", "file_without_ext"},
	}

	for _, tc := range testCases {
		result := generateOutputFileName(tc.input, tc.encoding)
		if result != tc.expected {
			t.Errorf("generateOutputFileName(%s, %s) = %s, expected %s",
				tc.input, tc.encoding, result, tc.expected)
		}
	}
}

func TestConfigDefaults(t *testing.T) {
	config := getDefaultConfig()

	if config.TargetEncoding != "UTF-8" {
		t.Errorf("Expected default TargetEncoding UTF-8, got %s", config.TargetEncoding)
	}
	if config.ConcurrencyLimit != 4 {
		t.Errorf("Expected default ConcurrencyLimit 4, got %d", config.ConcurrencyLimit)
	}
	if !config.CreateBackup {
		t.Error("Expected default CreateBackup true")
	}
	if config.OverwriteExisting {
		t.Error("Expected default OverwriteExisting false")
	}
	if config.MinConfidence != 0.8 {
		t.Errorf("Expected default MinConfidence 0.8, got %f", config.MinConfidence)
	}
}

// 性能测试
func BenchmarkConvertFile(b *testing.B) {
	testDir := b.TempDir()
	inputFile := filepath.Join(testDir, "bench_test.txt")

	// 创建较大的测试文件
	testContent := strings.Repeat("这是性能测试内容，包含中文和English mixed content.\n", 1000)
	err := os.WriteFile(inputFile, []byte(testContent), 0644)
	if err != nil {
		b.Fatalf("Failed to create test file: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		outputFile := filepath.Join(testDir, "bench_output_"+string(rune(i))+".txt")
		_, err := ConvertFile(inputFile, outputFile)
		if err != nil {
			b.Fatalf("ConvertFile failed: %v", err)
		}
	}
}
