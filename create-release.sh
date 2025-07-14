#!/bin/bash

# GitHub Release Script for ConvertContent2UTF8 v1.0.0

echo "🚀 Creating GitHub Release for ConvertContent2UTF8 v1.0.0..."

# Release notes content
RELEASE_NOTES="$(cat <<'EOF'
# 🎉 ConvertContent2UTF8 v1.0.0 - First Stable Release

A powerful, production-ready Go library for intelligent text file encoding conversion to UTF-8.

## ✨ Key Features

### 🧠 Smart Processing
- **Intelligent Encoding Detection**: Powered by encoding-processor v0.3.0+
- **Flexible Conversion**: Single file, batch files, and directory processing
- **Real-time Progress**: Detailed callbacks with status monitoring

### 🚀 High Performance  
- **Concurrent Processing**: Configurable concurrency limits
- **Memory Efficient**: Optimized for large-scale file conversion
- **Production Ready**: 80.5% test coverage, comprehensive error handling

### 🔤 Extensive Format Support
- **Unicode**: UTF-8, UTF-16, UTF-32 variants
- **Chinese**: GBK, GB2312, GB18030, BIG5
- **Japanese**: Shift_JIS, EUC-JP
- **Korean**: EUC-KR  
- **Western**: ISO-8859 series, Windows-125x series
- **Others**: KOI8-R, CP866, Macintosh

## 📦 Installation

\`\`\`bash
go get github.com/mirbf/ConvertContent2UTF8@v1.0.0
\`\`\`

## 🏃 Quick Start

\`\`\`go
import "github.com/mirbf/ConvertContent2UTF8"

// Single file conversion
result, err := convertcontent2utf8.ConvertFile("input.txt", "output.txt")

// Batch conversion with progress
files := []string{"file1.txt", "file2.txt", "file3.txt"}
result, err := convertcontent2utf8.ConvertFiles(files,
    convertcontent2utf8.WithProgress(func(p convertcontent2utf8.Progress) {
        fmt.Printf("Progress: %.1f%% - %s\n", 
            float64(p.ProcessedFiles)/float64(p.TotalFiles)*100, 
            p.CurrentFile)
    }),
    convertcontent2utf8.WithConcurrency(8),
)

// Directory conversion  
result, err := convertcontent2utf8.ConvertDirectory("/path/to/dir",
    convertcontent2utf8.WithRecursive(true),
    convertcontent2utf8.WithFileFilter(func(f string) bool {
        return strings.HasSuffix(f, ".txt")
    }),
)
\`\`\`

## 📋 Requirements
- **Go**: 1.20+
- **Dependencies**: encoding-processor v0.3.0+

## 🔗 Resources
- **Documentation**: https://pkg.go.dev/github.com/mirbf/ConvertContent2UTF8
- **Examples**: [examples/main.go](examples/main.go)
- **Issues**: https://github.com/mirbf/ConvertContent2UTF8/issues
- **Contributing**: [CONTRIBUTING.md](CONTRIBUTING.md)

## 🛡️ Security
See [SECURITY.md](SECURITY.md) for security policies and best practices.

---

**Full Changelog**: https://github.com/mirbf/ConvertContent2UTF8/commits/v1.0.0

🤖 Generated with [Claude Code](https://claude.ai/code)
EOF
)"

echo "📝 Release notes prepared"
echo "🏷️ Tag v1.0.0 already created and pushed"
echo "✅ GitHub Release ready to be created manually at:"
echo "   https://github.com/mirbf/ConvertContent2UTF8/releases/new?tag=v1.0.0"
echo ""
echo "📋 Release Notes Content:"
echo "================================"
echo "$RELEASE_NOTES"