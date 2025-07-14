# Contributing to ConvertContent2UTF8

我们欢迎并感谢您对本项目的贡献！请阅读以下指南以了解如何参与贡献。

## 如何贡献

### 报告 Bug

1. 确保该 bug 尚未被报告，请先搜索 [Issues](https://github.com/mirbf/ConvertContent2UTF8/issues)
2. 创建新的 Issue，使用清晰的标题和描述
3. 提供以下信息：
   - Go 版本
   - 操作系统
   - 重现步骤
   - 预期行为 vs 实际行为
   - 相关的错误日志

### 提交功能请求

1. 搜索现有的 Issues 确保功能尚未被请求
2. 创建新的 Issue，描述：
   - 功能的详细描述
   - 使用场景和动机
   - 可能的实现方案

### 提交代码

1. Fork 本仓库
2. 创建功能分支：`git checkout -b feature/amazing-feature`
3. 提交您的更改：`git commit -m 'Add amazing feature'`
4. 推送到分支：`git push origin feature/amazing-feature`
5. 提交 Pull Request

## 开发指南

### 环境设置

```bash
# 克隆仓库
git clone https://github.com/mirbf/ConvertContent2UTF8.git
cd ConvertContent2UTF8

# 安装依赖
go mod download

# 运行测试
go test -v
```

### 编码规范

- 遵循 Go 官方编码规范
- 使用 `gofmt` 格式化代码
- 添加必要的注释，特别是公共函数
- 保持函数简洁，单一职责
- 使用有意义的变量和函数名

### 测试要求

- 新功能必须包含测试
- 测试覆盖率应保持在 80% 以上
- 测试应包括正常情况和边界情况
- 运行 `go test -race` 确保并发安全

### 提交信息格式

使用以下格式：

```
类型(作用域): 简短描述

详细描述（可选）

Fixes #123
```

类型：
- `feat`: 新功能
- `fix`: bug 修复
- `docs`: 文档更新
- `test`: 测试相关
- `refactor`: 重构
- `perf`: 性能优化
- `chore`: 构建过程或辅助工具的变动

### Pull Request 指南

1. 确保所有测试通过
2. 更新相关文档
3. 在 PR 描述中说明：
   - 更改的内容
   - 测试的内容
   - 相关的 Issue

## 代码审查

所有提交都需要代码审查。我们会检查：

- 代码质量和规范
- 测试覆盖率
- 文档完整性
- 向后兼容性
- 性能影响

## 发布流程

1. 更新版本号
2. 更新 CHANGELOG.md
3. 创建 Git 标签
4. 发布到 GitHub Releases

## 问题和支持

如有疑问，请：

1. 查看现有的 Issues 和 Discussions
2. 创建新的 Issue
3. 联系维护者

感谢您的贡献！