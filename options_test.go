package convertcontent2utf8

import (
	"testing"
)

func TestWithProgress(t *testing.T) {
	called := false
	callback := func(p Progress) {
		called = true
	}

	config := getDefaultConfig()
	option := WithProgress(callback)
	option(config)

	if config.ProgressCallback == nil {
		t.Error("ProgressCallback should not be nil")
	}

	// 测试回调函数
	config.ProgressCallback(Progress{})
	if !called {
		t.Error("Callback function was not called")
	}
}

func TestWithTargetEncoding(t *testing.T) {
	testCases := []string{"GBK", "BIG5", "Shift_JIS", "UTF-16"}

	for _, encoding := range testCases {
		config := getDefaultConfig()
		option := WithTargetEncoding(encoding)
		option(config)

		if config.TargetEncoding != encoding {
			t.Errorf("Expected TargetEncoding %s, got %s", encoding, config.TargetEncoding)
		}
	}
}

func TestWithConcurrency(t *testing.T) {
	testCases := []int{1, 2, 8, 16}

	for _, limit := range testCases {
		config := getDefaultConfig()
		option := WithConcurrency(limit)
		option(config)

		if config.ConcurrencyLimit != limit {
			t.Errorf("Expected ConcurrencyLimit %d, got %d", limit, config.ConcurrencyLimit)
		}
	}
}

func TestWithFileFilter(t *testing.T) {
	filter := func(filename string) bool {
		return filename == "test.txt"
	}

	config := getDefaultConfig()
	option := WithFileFilter(filter)
	option(config)

	if config.FileFilter == nil {
		t.Error("FileFilter should not be nil")
	}

	// 测试过滤器函数
	if !config.FileFilter("test.txt") {
		t.Error("Filter should return true for test.txt")
	}
	if config.FileFilter("other.txt") {
		t.Error("Filter should return false for other.txt")
	}
}

func TestWithBackup(t *testing.T) {
	testCases := []bool{true, false}

	for _, backup := range testCases {
		config := getDefaultConfig()
		option := WithBackup(backup)
		option(config)

		if config.CreateBackup != backup {
			t.Errorf("Expected CreateBackup %v, got %v", backup, config.CreateBackup)
		}
	}
}

func TestWithOverwrite(t *testing.T) {
	testCases := []bool{true, false}

	for _, overwrite := range testCases {
		config := getDefaultConfig()
		option := WithOverwrite(overwrite)
		option(config)

		if config.OverwriteExisting != overwrite {
			t.Errorf("Expected OverwriteExisting %v, got %v", overwrite, config.OverwriteExisting)
		}
	}
}

func TestWithMinConfidence(t *testing.T) {
	testCases := []float64{0.5, 0.7, 0.9, 1.0}

	for _, confidence := range testCases {
		config := getDefaultConfig()
		option := WithMinConfidence(confidence)
		option(config)

		if config.MinConfidence != confidence {
			t.Errorf("Expected MinConfidence %f, got %f", confidence, config.MinConfidence)
		}
	}
}

func TestWithDryRun(t *testing.T) {
	testCases := []bool{true, false}

	for _, dryRun := range testCases {
		config := getDefaultConfig()
		option := WithDryRun(dryRun)
		option(config)

		if config.DryRun != dryRun {
			t.Errorf("Expected DryRun %v, got %v", dryRun, config.DryRun)
		}
	}
}

func TestWithSkipHidden(t *testing.T) {
	testCases := []bool{true, false}

	for _, skip := range testCases {
		config := getDefaultConfig()
		option := WithSkipHidden(skip)
		option(config)

		if config.SkipHidden != skip {
			t.Errorf("Expected SkipHidden %v, got %v", skip, config.SkipHidden)
		}
	}
}

func TestWithRecursive(t *testing.T) {
	testCases := []bool{true, false}

	for _, recursive := range testCases {
		config := getDefaultConfig()
		option := WithRecursive(recursive)
		option(config)

		if config.Recursive != recursive {
			t.Errorf("Expected Recursive %v, got %v", recursive, config.Recursive)
		}
	}
}

func TestWithMaxFileSize(t *testing.T) {
	testCases := []int64{1024, 1024 * 1024, 10 * 1024 * 1024}

	for _, size := range testCases {
		config := getDefaultConfig()
		option := WithMaxFileSize(size)
		option(config)

		if config.MaxFileSize != size {
			t.Errorf("Expected MaxFileSize %d, got %d", size, config.MaxFileSize)
		}
	}
}

func TestApplyOptions(t *testing.T) {
	options := []Option{
		WithTargetEncoding("GBK"),
		WithConcurrency(8),
		WithBackup(false),
		WithOverwrite(true),
		WithMinConfidence(0.9),
		WithDryRun(true),
		WithSkipHidden(false),
		WithRecursive(true),
		WithMaxFileSize(50 * 1024 * 1024),
	}

	config := applyOptions(options)

	expected := &Config{
		TargetEncoding:    "GBK",
		ConcurrencyLimit:  8,
		CreateBackup:      false,
		OverwriteExisting: true,
		MinConfidence:     0.9,
		DryRun:            true,
		SkipHidden:        false,
		Recursive:         true,
		MaxFileSize:       50 * 1024 * 1024,
		FileFilter:        config.FileFilter, // 函数不能直接比较
	}

	// 比较除了FileFilter之外的所有字段
	if config.TargetEncoding != expected.TargetEncoding {
		t.Errorf("Expected TargetEncoding %s, got %s", expected.TargetEncoding, config.TargetEncoding)
	}
	if config.ConcurrencyLimit != expected.ConcurrencyLimit {
		t.Errorf("Expected ConcurrencyLimit %d, got %d", expected.ConcurrencyLimit, config.ConcurrencyLimit)
	}
	if config.CreateBackup != expected.CreateBackup {
		t.Errorf("Expected CreateBackup %v, got %v", expected.CreateBackup, config.CreateBackup)
	}
	if config.OverwriteExisting != expected.OverwriteExisting {
		t.Errorf("Expected OverwriteExisting %v, got %v", expected.OverwriteExisting, config.OverwriteExisting)
	}
	if config.MinConfidence != expected.MinConfidence {
		t.Errorf("Expected MinConfidence %f, got %f", expected.MinConfidence, config.MinConfidence)
	}
	if config.DryRun != expected.DryRun {
		t.Errorf("Expected DryRun %v, got %v", expected.DryRun, config.DryRun)
	}
	if config.SkipHidden != expected.SkipHidden {
		t.Errorf("Expected SkipHidden %v, got %v", expected.SkipHidden, config.SkipHidden)
	}
	if config.Recursive != expected.Recursive {
		t.Errorf("Expected Recursive %v, got %v", expected.Recursive, config.Recursive)
	}
	if config.MaxFileSize != expected.MaxFileSize {
		t.Errorf("Expected MaxFileSize %d, got %d", expected.MaxFileSize, config.MaxFileSize)
	}
}

func TestDefaultFileFilter(t *testing.T) {
	config := getDefaultConfig()

	testCases := []struct {
		filename string
		expected bool
	}{
		{"test.txt", true},
		{"TEST.TXT", true},
		{"file.TXT", true},
		{"document.doc", false},
		{"image.jpg", false},
		{"script.sh", false},
		{"readme.md", false},
		{"config.json", false},
		{"/path/to/file.txt", true},
		{"/path/to/file.log", false},
	}

	for _, tc := range testCases {
		result := config.FileFilter(tc.filename)
		if result != tc.expected {
			t.Errorf("FileFilter(%s) = %v, expected %v", tc.filename, result, tc.expected)
		}
	}
}

func TestApplyOptionsOrder(t *testing.T) {
	// 测试选项应用的顺序，后面的选项应该覆盖前面的
	options := []Option{
		WithTargetEncoding("GBK"),
		WithTargetEncoding("BIG5"), // 这个应该覆盖上面的
		WithConcurrency(2),
		WithConcurrency(8), // 这个应该覆盖上面的
	}

	config := applyOptions(options)

	if config.TargetEncoding != "BIG5" {
		t.Errorf("Expected TargetEncoding BIG5, got %s", config.TargetEncoding)
	}
	if config.ConcurrencyLimit != 8 {
		t.Errorf("Expected ConcurrencyLimit 8, got %d", config.ConcurrencyLimit)
	}
}

func TestEmptyOptions(t *testing.T) {
	config := applyOptions([]Option{})
	defaultConfig := getDefaultConfig()

	// 检查几个关键字段是否与默认配置相同
	if config.TargetEncoding != defaultConfig.TargetEncoding {
		t.Errorf("Expected TargetEncoding %s, got %s", defaultConfig.TargetEncoding, config.TargetEncoding)
	}
	if config.ConcurrencyLimit != defaultConfig.ConcurrencyLimit {
		t.Errorf("Expected ConcurrencyLimit %d, got %d", defaultConfig.ConcurrencyLimit, config.ConcurrencyLimit)
	}
	if config.CreateBackup != defaultConfig.CreateBackup {
		t.Errorf("Expected CreateBackup %v, got %v", defaultConfig.CreateBackup, config.CreateBackup)
	}
}

// 测试选项函数的类型
func TestOptionType(t *testing.T) {
	var option Option = WithTargetEncoding("UTF-8")

	// 检查Option能否正确应用
	config := getDefaultConfig()
	option(config)

	if config.TargetEncoding != "UTF-8" {
		t.Errorf("Option function failed to apply, expected UTF-8, got %s", config.TargetEncoding)
	}
}
