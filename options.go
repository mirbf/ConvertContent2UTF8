package convertcontent2utf8

import (
	"path/filepath"
	"strings"
)

// WithProgress 设置进度回调
func WithProgress(callback func(Progress)) Option {
	return func(c *Config) {
		c.ProgressCallback = callback
	}
}

// WithTargetEncoding 设置目标编码
func WithTargetEncoding(encoding string) Option {
	return func(c *Config) {
		c.TargetEncoding = encoding
	}
}

// WithConcurrency 设置并发限制
func WithConcurrency(limit int) Option {
	return func(c *Config) {
		c.ConcurrencyLimit = limit
	}
}

// WithFileFilter 设置文件过滤器
func WithFileFilter(filter func(string) bool) Option {
	return func(c *Config) {
		c.FileFilter = filter
	}
}

// WithBackup 设置是否创建备份
func WithBackup(create bool) Option {
	return func(c *Config) {
		c.CreateBackup = create
	}
}

// WithOverwrite 设置是否覆盖已存在文件
func WithOverwrite(overwrite bool) Option {
	return func(c *Config) {
		c.OverwriteExisting = overwrite
	}
}

// WithMinConfidence 设置最小检测置信度
func WithMinConfidence(confidence float64) Option {
	return func(c *Config) {
		c.MinConfidence = confidence
	}
}

// WithDryRun 设置试运行模式
func WithDryRun(dryRun bool) Option {
	return func(c *Config) {
		c.DryRun = dryRun
	}
}

// WithSkipHidden 设置是否跳过隐藏文件
func WithSkipHidden(skip bool) Option {
	return func(c *Config) {
		c.SkipHidden = skip
	}
}

// WithRecursive 设置是否递归处理目录
func WithRecursive(recursive bool) Option {
	return func(c *Config) {
		c.Recursive = recursive
	}
}

// WithMaxFileSize 设置最大文件大小限制
func WithMaxFileSize(size int64) Option {
	return func(c *Config) {
		c.MaxFileSize = size
	}
}

// getDefaultConfig 获取默认配置
func getDefaultConfig() *Config {
	return &Config{
		TargetEncoding:    "UTF-8",
		ConcurrencyLimit:  4,
		CreateBackup:      true,
		OverwriteExisting: false,
		MinConfidence:     0.8,
		DryRun:            false,
		SkipHidden:        true,
		Recursive:         false,
		MaxFileSize:       100 * 1024 * 1024, // 100MB
		FileFilter: func(filename string) bool {
			// 默认只处理.txt文件
			return filepath.Ext(strings.ToLower(filename)) == ".txt"
		},
	}
}

// applyOptions 应用配置选项
func applyOptions(options []Option) *Config {
	config := getDefaultConfig()
	for _, option := range options {
		option(config)
	}
	return config
}
