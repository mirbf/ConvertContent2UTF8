package convertcontent2utf8

import (
	"time"

	encoding "github.com/mirbf/encoding-processor"
)

// ProgressStatus 进度状态枚举
type ProgressStatus string

const (
	StatusStarting   ProgressStatus = "starting"
	StatusProcessing ProgressStatus = "processing"
	StatusCompleted  ProgressStatus = "completed"
	StatusFailed     ProgressStatus = "failed"
	StatusSkipped    ProgressStatus = "skipped"
)

// Progress 进度信息结构
type Progress struct {
	// 核心进度信息
	CurrentFile    string         `json:"current_file"`
	ProcessedFiles int            `json:"processed_files"`
	TotalFiles     int            `json:"total_files"`
	Status         ProgressStatus `json:"status"`

	// 时间信息
	StartTime     time.Time     `json:"start_time"`
	ElapsedTime   time.Duration `json:"elapsed_time"`
	EstimatedTime time.Duration `json:"estimated_time,omitempty"`

	// 可选详细信息
	FileSize       int64 `json:"file_size,omitempty"`
	ProcessedBytes int64 `json:"processed_bytes,omitempty"`
	ErrorCount     int   `json:"error_count"`
}

// ConvertResult 单文件转换结果
type ConvertResult struct {
	InputFile           string                      `json:"input_file"`
	OutputFile          string                      `json:"output_file"`
	SourceEncoding      string                      `json:"source_encoding"`
	TargetEncoding      string                      `json:"target_encoding"`
	BytesProcessed      int64                       `json:"bytes_processed"`
	ProcessingTime      time.Duration               `json:"processing_time"`
	DetectionConfidence float64                     `json:"detection_confidence"`
	BackupFile          string                      `json:"backup_file,omitempty"`
	ProcessorResult     *encoding.FileProcessResult `json:"-"` // 底层库结果
}

// BatchResult 批量转换结果
type BatchResult struct {
	TotalFiles      int              `json:"total_files"`
	ProcessedFiles  int              `json:"processed_files"`
	SuccessfulFiles int              `json:"successful_files"`
	FailedFiles     int              `json:"failed_files"`
	SkippedFiles    int              `json:"skipped_files"`
	TotalBytes      int64            `json:"total_bytes"`
	ProcessingTime  time.Duration    `json:"processing_time"`
	Results         []*ConvertResult `json:"results"`
	Errors          []FileError      `json:"errors,omitempty"`
}

// FileError 文件处理错误
type FileError struct {
	File      string    `json:"file"`
	Operation string    `json:"operation"`
	Error     string    `json:"error"`
	Timestamp time.Time `json:"timestamp"`
}

// Config 配置结构
type Config struct {
	// 目标编码，默认UTF-8
	TargetEncoding string

	// 进度回调函数
	ProgressCallback func(Progress)

	// 并发限制，默认为4
	ConcurrencyLimit int

	// 文件过滤器
	FileFilter func(string) bool

	// encoding-processor选项
	CreateBackup      bool
	OverwriteExisting bool
	MinConfidence     float64
	DryRun            bool

	// 其他选项
	SkipHidden  bool  // 跳过隐藏文件
	Recursive   bool  // 目录递归处理
	MaxFileSize int64 // 最大文件大小限制
}

// Option 配置选项函数类型
type Option func(*Config)
