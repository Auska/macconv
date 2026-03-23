# AGENTS.md

This file provides guidance to CodeBuddy Code when working with code in this repository.

## 项目概述

macconv 是一个 Go CLI 工具（Go 1.22+，Cobra 框架），用于网络地址格式转换和常见网络操作：MAC 地址转换、CIDR 掩码转换、TCP 端口检查、DHCP 选项 43 转换。

## 常用命令

```bash
# 构建
make build                    # 等同于 go build -ldflags="-s -w ..." -o macconv .

# 测试（运行所有测试，包括 cmd/）
make test                     # 运行 ./pkg/... 和 -tags=unit ./cmd/...
go test ./pkg/...             # 仅 pkg/ 测试
go test ./pkg/errors/...      # 单个包
go test ./pkg/errors/... -run TestNew  # 单个测试函数
go test -tags=unit ./cmd/... -run TestMacAddress  # cmd 测试需要 unit 标签

# 覆盖率
make test-coverage            # 生成 coverage.out 和 coverage.html

# 代码检查
make lint                     # golangci-lint run
make fmt                      # gofmt
make vet                      # go vet

# 跨平台构建
make build-all                # linux/windows/darwin (amd64/arm64)

# 清理
make clean
```

## 测试注意事项

**`cmd/` 和 `pkg/` 的测试方式不同：**
- `pkg/errors/`、`pkg/logger/`、`pkg/validator/` 的测试文件没有构建标签，直接运行 `go test ./pkg/...` 即可
- `cmd/` 下的测试文件（`mac_test.go`、`ip_test.go` 等）全部使用 `//go:build unit` 构建标签，**必须加 `-tags=unit`** 才能运行
- `make test` 会运行所有测试（包括 cmd/ 下的单元测试）

## 项目架构

```
macconv/
├── main.go           # 入口，通过 -ldflags 注入 version 和 buildDate
├── cmd/              # Cobra 子命令实现
│   ├── root.go       # 根命令 + 全局 flag（-v version, -l log-level）
│   ├── mac.go        # mac 命令
│   ├── ip.go         # ip 命令（CIDR/IPv4/IPv6）
│   ├── tcp.go        # tcp 命令（端口持续检查）
│   ├── dhcp.go       # dhcp 命令（选项 43 转换）
│   └── version.go    # version 子命令
└── pkg/
    ├── errors/       # AppError 类型系统（ValidationError, NetworkError, FileSystemError, ParseError）
    ├── logger/       # 结构化日志（debug/info/warn/error 级别，支持实例方法和包级函数两种调用方式）
    └── validator/    # 输入验证（MAC、IP、端口、CIDR、文件路径）
```

**数据流：** `main.go` → `cmd.Execute()` → Cobra 路由到子命令 → 调用 `pkg/validator` 验证输入 → 使用 `pkg/logger` 输出结果/错误

## 关键约束

### depguard 导入限制
`.golangci.yml` 中配置了严格的 depguard 规则，**每个文件只允许导入特定的包**，违反会导致 lint 失败：
- `cmd/tcp.go` 不允许导入 `macconv/pkg/errors` 和 `macconv/pkg/validator`
- `cmd/version.go` 只允许导入 `cobra`
- `pkg/validator/validator.go` 只允许导入 `macconv/pkg/errors`

修改任何文件前，先检查 `.golangci.yml` 中对应的 `depguard.rules` 条目。

### 版本信息注入
`main.go` 中的 `version` 和 `buildDate` 变量通过 `make build` 的 `-ldflags` 在构建时注入。`cmd.SetVersionInfo()` 将它们传递给 `cmd` 包。

### 日志系统
`pkg/logger` 提供两种等价的调用方式，代码中混用：
- 包级函数：`logger.Debug("msg %s", arg)`
- 实例方法：`logger.DefaultLogger.Debug("msg %s", arg)`

默认日志级别为 `warn`，通过全局 `-l` flag 控制。

### 错误处理
所有业务错误通过 `pkg/errors` 创建，使用 `errors.New(type, msg)` 或 `errors.Wrap(type, msg, err)`。错误类型可用于 `errors.Is` / 类型断言判断。

## 开发约定

1. **版权声明**: 所有源文件包含 `Copyright © 2024-2025 Auska <luodan0709@live.cn>`
2. **提交信息**: 使用英文
3. **构建优化**: 使用 `-ldflags="-s -w"` 减小二进制体积
4. **正则表达式**: 使用包级 `var` 预编译（参考 `pkg/validator/validator.go`）
5. **常量**: 避免魔法数字，使用包级 `const`
6. **CI**: GitHub Actions 在 push/PR 到 main/develop 时运行 test + lint + 多平台 build
