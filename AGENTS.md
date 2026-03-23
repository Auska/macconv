# AGENTS.md

This file provides guidance to CodeBuddy Code when working with code in this repository.

## 项目概述

macconv 是一个 Go CLI 工具（Go 1.22+，Cobra 框架），用于网络地址格式转换和常见网络操作：MAC 地址转换、CIDR 掩码转换、TCP 端口检查、DHCP 选项 43 转换。

## 常用命令

```bash
# 构建
make build                    # 等同于 go build -ldflags="-s -w ..." -o macconv .
make run                      # go run . （快速运行不生成二进制）
make build-all                # linux/windows/darwin (amd64/arm64)

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

# 组合命令
make all                      # 默认目标：fmt → vet → test → build
make dev                      # 开发模式：fmt → vet → test → run
make release                  # 完整发布：clean → fmt → vet → test → build-all

# 其他
make deps                     # go mod download && go mod tidy
make docs                     # godoc -http=:6060

# 清理
make clean
```

## 测试注意事项

**`cmd/` 和 `pkg/` 的测试方式不同：**
- `pkg/errors/`、`pkg/logger/`、`pkg/validator/` 的测试文件没有构建标签，直接运行 `go test ./pkg/...` 即可
- `cmd/` 下的测试文件（`mac_test.go`、`ip_test.go` 等）全部使用 `//go:build unit` 构建标签，**必须加 `-tags=unit`** 才能运行
- `cmd/test_helpers.go`（同样带 `//go:build unit` 标签）提供 `captureOutput()` / `restoreOutput()` 辅助函数，所有 cmd 测试依赖它
- `make test` 会运行所有测试（包括 cmd/ 下的单元测试）

## 项目架构

```
macconv/
├── main.go           # 入口，通过 -ldflags 注入 version 和 buildDate
├── cmd/              # Cobra 子命令实现（所有子命令通过 init() 调用 rootCmd.AddCommand() 注册）
│   ├── root.go       # 根命令 + 全局 flag（-v version, -l log-level）
│   ├── mac.go        # mac 命令（唯一允许导入 pkg/validator 的 cmd 文件）
│   ├── ip.go         # ip 命令（CIDR/IPv4/IPv6）
│   ├── tcp.go        # tcp 命令（端口持续检查）
│   ├── dhcp.go       # dhcp 命令（选项 43 转换，仅 IPv4）
│   ├── version.go    # version 子命令
│   └── test_helpers.go  # [//go:build unit] 测试辅助函数（stdout 捕获）
└── pkg/
    ├── errors/       # AppError 类型系统（ValidationError, NetworkError, FileSystemError, ParseError）
    ├── logger/       # 结构化日志（debug/info/warn/error 级别）
    └── validator/    # 输入验证（MAC、IP、端口、CIDR、文件路径）
```

**数据流：** `main.go` → `cmd.Execute()` → Cobra 路由到子命令 → 调用 `pkg/validator` 验证输入 → 使用 `pkg/logger` 输出结果/错误

**注意：** 由于 depguard 限制，`cmd/tcp.go` 和 `cmd/dhcp.go` 无法导入 `pkg/validator`，因此自行实现了端口解析和 IP 验证逻辑。

## 关键约束

### depguard 导入限制

`.golangci.yml` 中配置了严格的 depguard 规则，**每个文件只允许导入特定的包**，违反会导致 lint 失败。完整规则如下：

| 文件 | 允许导入的包 |
|------|-------------|
| `cmd/` (默认) | `cobra`, `pkg/errors`, `pkg/logger` |
| `cmd/mac.go` | `cobra`, `pkg/errors`, `pkg/logger`, **`pkg/validator`** |
| `cmd/tcp.go` | `cobra`, `pkg/logger` |
| `cmd/version.go` | `cobra` |
| `pkg/validator/validator.go` | `pkg/errors` |

**关键约束：**
- `cmd/mac.go` 是**唯一**允许导入 `pkg/validator` 的 cmd 文件
- `cmd/tcp.go` 不允许导入 `pkg/errors` 和 `pkg/validator`，因此自行实现 `parsePort()` 并直接 `fmt.Println` 输出
- 全局禁止直接导入 `pflag` 和 `mousetrap`（虽为 cobra 间接依赖）

修改任何文件前，先检查 `.golangci.yml` 中对应的 `depguard.rules` 条目。

### 版本信息注入

`main.go` 中的 `version` 和 `buildDate` 变量通过 `make build` 的 `-ldflags` 在构建时注入。`cmd.SetVersionInfo()` 将它们传递给 `cmd` 包。root 命令的 `--version` / `-v` flag 和独立的 `version` 子命令都调用同一个 `printVersion()` 函数。

### 日志系统

`pkg/logger` 提供两种等价的调用方式，代码中混用：
- 包级函数：`logger.Debug("msg %s", arg)`
- 实例方法：`logger.DefaultLogger.Debug("msg %s", arg)`

默认日志级别为 `warn`，通过全局 `-l` flag 控制。

cmd 代码中常用的便捷辅助函数：
- `logger.PrintError(err)` — 打印错误
- `logger.PrintErrorWithMessage(msg, err)` — 打印带上下文的错误
- `logger.PrintValidationError(msg)` — 打印验证错误（前缀 "Validation error: "）

### 错误处理

所有业务错误通过 `pkg/errors` 创建，使用 `errors.New(type, msg)` 或 `errors.Wrap(type, msg, err)`。错误类型可用于 `errors.Is` / 类型断言判断。

## Lint 配置

`.golangci.yml` 启用了 28 个 linter，以下为影响日常开发的关键配置：

| 规则 | 配置 |
|------|------|
| 行长度 (`lll`) | **140 字符** |
| 圈复杂度 (`gocyclo`) | 阈值 **15** |
| 代码重复 (`dupl`) | 阈值 **100 tokens** |
| 变量遮蔽 (`govet`) | 启用 `check-shadowing` |
| 拼写检查 (`misspell`) | **美式英语** |
| 导入排序 (`goimports`) | `local-prefixes: macconv` |
| `gochecknoinits` | 对 `_test.go` 和 `cmd/` **排除**（允许 cmd 中使用 `init()` 注册子命令） |
| `_test.go` 排除 | `gomnd`, `funlen`, `goconst` |

## CI 配置

GitHub Actions 在 push 到 main/develop、PR 到 main、release 发布时触发。4 个 job 按依赖顺序执行：

1. **test** — 运行测试 + 上传覆盖率到 Codecov
2. **lint** — golangci-lint（`--timeout=5m`）
3. **build** — 依赖 test + lint 通过；多平台构建（**排除** `windows/arm64`）；使用 `CGO_ENABLED=0`、`-trimpath`、`-buildid=` 构建纯静态二进制；生成多算法校验和
4. **release** — 依赖 build；仅在 release 事件时运行，自动上传构建产物到 GitHub Release

**注意：** CI 额外使用了 `-trimpath` 和 `-buildid=` flag，`make build-all` 中没有这些 flag。CI 排除 `windows/arm64`，而 `make build-all` 包含所有平台。

## 开发约定

1. **版权声明**: 所有源文件包含 `Copyright © 2024-2025 Auska <luodan0709@live.cn>`
2. **提交信息**: 使用英文
3. **构建优化**: 使用 `-ldflags="-s -w"` 减小二进制体积
4. **正则表达式**: 使用包级 `var` 预编译（参考 `pkg/validator/validator.go`）
5. **常量**: 避免魔法数字，使用包级 `const`
6. **子命令注册**: cmd 文件通过 `init()` 函数 + `rootCmd.AddCommand()` 注册（`gochecknoinits` 已对 cmd/ 排除）
7. **测试模式**: 所有测试使用 table-driven 模式（`[]struct{...}` + `t.Run`）
