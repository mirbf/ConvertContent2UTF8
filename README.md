# ConvertContent2UTF8

[![Go Report Card](https://goreportcard.com/badge/github.com/mirbf/ConvertContent2UTF8)](https://goreportcard.com/report/github.com/mirbf/ConvertContent2UTF8)
[![GoDoc](https://godoc.org/github.com/mirbf/ConvertContent2UTF8?status.svg)](https://godoc.org/github.com/mirbf/ConvertContent2UTF8)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![GitHub release](https://img.shields.io/github/release/mirbf/ConvertContent2UTF8.svg)](https://github.com/mirbf/ConvertContent2UTF8/releases)

一个专注于文本文档UTF8编码转换的Go语言库，基于[encoding-processor](https://github.com/mirbf/encoding-processor)构建，提供批量处理和进度监控功能。

[English](README_EN.md) | 中文

## ✨ 特性

- 🔄 **单文件转换**: 支持单个文件的编码转换
- 📁 **批量处理**: 支持多文件批量转换，内置并发控制  
- 🌊 **目录递归**: 支持目录递归遍历和转换
- 📊 **进度监控**: 实时进度回调，支持自定义进度显示
- ⚙️ **高度可配置**: 丰富的配置选项，满足不同场景需求
- 🛡️ **完善错误处理**: 结构化错误信息和恢复机制
- 🔍 **智能编码检测**: 基于encoding-processor的智能编码检测
- 💾 **安全备份**: 支持自动备份和文件恢复
- 🚀 **高性能**: 并发处理，支持大规模文件转换

## 📦 安装

```bash
go get github.com/mirbf/ConvertContent2UTF8
```

## 🚀 快速开始

### 基础使用

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/mirbf/ConvertContent2UTF8"
)

func main() {
    // 单文件转换
    result, err := ConvertContent2UTF8.ConvertFile("input.txt", "output.txt")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("转换完成: %s -> %s\n", result.SourceEncoding, result.TargetEncoding)
}
```

### 批量文件转换

```go
files := []string{"file1.txt", "file2.txt", "file3.txt"}

result, err := ConvertContent2UTF8.ConvertFiles(files,
    ConvertContent2UTF8.WithOverwrite(true),
    ConvertContent2UTF8.WithConcurrency(8),
)

if err != nil {
    log.Fatal(err)
}

fmt.Printf("处理完成: 成功 %d, 失败 %d\n", 
    result.SuccessfulFiles, result.FailedFiles)
```

### 目录递归转换

```go
result, err := ConvertContent2UTF8.ConvertDirectory("/path/to/directory",
    ConvertContent2UTF8.WithRecursive(true),
    ConvertContent2UTF8.WithFileFilter(func(filename string) bool {
        return strings.HasSuffix(filename, ".txt") || 
               strings.HasSuffix(filename, ".md")
    }),
)
```

### 带进度回调的转换

```go
result, err := ConvertContent2UTF8.ConvertFiles(files,
    ConvertContent2UTF8.WithProgress(func(p ConvertContent2UTF8.Progress) {
        percentage := float64(p.ProcessedFiles) / float64(p.TotalFiles) * 100
        fmt.Printf("进度: %.1f%% - %s [%s]\n", 
            percentage, p.CurrentFile, p.Status)
        
        if p.EstimatedTime > 0 {
            fmt.Printf("预计剩余时间: %v\n", p.EstimatedTime)
        }
    }),
)
```

## 📚 API 接口

### 核心函数

```go
// 单文件转换
func ConvertFile(inputFile, outputFile string, options ...Option) (*ConvertResult, error)

// 批量文件转换  
func ConvertFiles(files []string, options ...Option) (*BatchResult, error)

// 目录递归转换
func ConvertDirectory(dirPath string, options ...Option) (*BatchResult, error)
```

### 配置选项

| 选项 | 说明 | 默认值 |
|------|------|--------|
| `WithProgress(callback)` | 进度回调函数 | 无 |
| `WithTargetEncoding(encoding)` | 目标编码 | UTF-8 |
| `WithConcurrency(limit)` | 并发限制 | 4 |
| `WithFileFilter(filter)` | 文件过滤器 | .txt文件 |
| `WithBackup(create)` | 创建备份 | true |
| `WithOverwrite(overwrite)` | 覆盖已存在文件 | false |
| `WithMinConfidence(confidence)` | 最小检测置信度 | 0.8 |
| `WithDryRun(dryRun)` | 试运行模式 | false |
| `WithSkipHidden(skip)` | 跳过隐藏文件 | true |
| `WithRecursive(recursive)` | 递归处理目录 | false |
| `WithMaxFileSize(size)` | 最大文件大小限制 | 100MB |

## 📊 数据结构

### Progress 进度信息

```go
type Progress struct {
    CurrentFile    string         // 当前处理的文件
    ProcessedFiles int            // 已处理文件数
    TotalFiles     int            // 总文件数
    Status         ProgressStatus // 当前状态
    StartTime      time.Time      // 开始时间
    ElapsedTime    time.Duration  // 已耗时
    EstimatedTime  time.Duration  // 预计剩余时间
    ProcessedBytes int64          // 已处理字节数
    ErrorCount     int            // 错误数量
}
```

### ConvertResult 转换结果

```go
type ConvertResult struct {
    InputFile           string        // 输入文件
    OutputFile          string        // 输出文件
    SourceEncoding      string        // 源编码
    TargetEncoding      string        // 目标编码
    BytesProcessed      int64         // 处理字节数
    ProcessingTime      time.Duration // 处理时间
    DetectionConfidence float64       // 检测置信度
    BackupFile          string        // 备份文件
}
```

### BatchResult 批量处理结果

```go
type BatchResult struct {
    TotalFiles      int              // 总文件数
    ProcessedFiles  int              // 已处理文件数
    SuccessfulFiles int              // 成功文件数
    FailedFiles     int              // 失败文件数
    SkippedFiles    int              // 跳过文件数
    TotalBytes      int64            // 总字节数
    ProcessingTime  time.Duration    // 处理时间
    Results         []*ConvertResult // 详细结果
    Errors          []FileError      // 错误列表
}
```

## 🔤 支持的编码

基于 [encoding-processor](https://github.com/mirbf/encoding-processor)，支持以下编码格式：

- **Unicode**: UTF-8, UTF-16, UTF-16LE, UTF-16BE, UTF-32*, UTF-32LE*, UTF-32BE*
- **中文**: GBK, GB2312, GB18030, BIG5  
- **日文**: Shift_JIS, EUC-JP
- **韩文**: EUC-KR
- **西欧**: ISO-8859-1, ISO-8859-2, ISO-8859-5, ISO-8859-15
- **Windows**: Windows-1250, Windows-1251, Windows-1252, Windows-1254
- **其他**: KOI8-R, CP866, Macintosh

## 💡 使用示例

查看 [examples/main.go](./examples/main.go) 了解完整的使用示例。

运行示例：

```bash
cd examples
go run main.go
```

## 🧪 测试

```bash
# 运行所有测试
go test -v

# 运行性能测试
go test -bench=.

# 查看测试覆盖率
go test -cover

# 并发安全测试
go test -race
```

## 🎯 设计原则

- **专注职责**: 专注于批量处理和用户便利性，编码检测和转换依赖成熟的encoding-processor库
- **并发安全**: 所有公共接口都支持并发调用
- **错误透明**: 完善的错误处理和恢复机制
- **进度可视**: 详细的进度信息和回调支持
- **配置灵活**: 丰富的配置选项满足不同场景需求

## 📋 依赖

- **Go**: 1.20+
- **github.com/mirbf/encoding-processor**: v0.3.0+ - 编码检测和转换核心库

### 版本兼容性

本库要求 `encoding-processor` 最低版本为 v0.3.0，这确保了：
- 智能编码检测功能的可用性
- API 稳定性和向后兼容
- 使用者可以安全升级到 v0.3.x 的任意版本

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！请阅读 [CONTRIBUTING.md](CONTRIBUTING.md) 了解贡献指南。

## 📜 许可证

[MIT License](LICENSE)

## 🆕 更新日志

查看 [CHANGELOG.md](CHANGELOG.md) 了解版本更新记录。

## 🔒 安全

查看 [SECURITY.md](SECURITY.md) 了解安全政策和最佳实践。

---

<p align="center">
  ⭐ 如果这个项目对你有帮助，请给我们一个 star！
</p>