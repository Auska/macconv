# MAC地址格式转换工具 (macconv)

## 项目概述

macconv 是一个用 Go 语言编写的命令行工具，用于网络地址格式转换和常见网络操作。该工具基于 Cobra 框架构建，提供了多种网络相关的实用功能，包括 MAC 地址格式转换、CIDR 掩码转换、端口检查、DHCP 选项 43 转换等。项目经过全面重构，具有统一的错误处理、结构化日志记录和全面的测试覆盖。

## 主要技术

- **语言**: Go 1.22.0
- **框架**: Cobra (CLI 框架)
- **依赖库**:
  - github.com/spf13/cobra v1.8.0
  - github.com/spf13/pflag v1.0.5
  - github.com/inconshreveable/mousetrap v1.1.0

## 项目结构

```
macconv/
├── cmd/              # 命令实现目录
│   ├── root.go       # 根命令定义和版本信息处理
│   ├── mac.go        # MAC 地址转换命令
│   ├── ip.go         # CIDR 掩码转换命令（支持反掩码和IPv6）
│   ├── tcp.go        # 端口检查命令（支持域名和持续监控）
│   ├── dhcp.go       # DHCP 选项 43 转换命令
│   ├── version.go    # 版本信息命令
│   ├── mac_test.go   # MAC 地址转换测试
│   └── ip_test.go    # CIDR 转换测试
├── pkg/              # 共享包目录
│   ├── errors/       # 统一错误处理系统
│   │   ├── errors.go
│   │   └── errors_test.go
│   ├── logger/       # 结构化日志记录系统
│   │   └── logger.go
│   └── validator/    # 通用验证工具包
│       ├── validator.go
│       └── validator_test.go
├── .github/          # GitHub Actions 工作流
│   └── workflows/
│       └── test.yml  # 多平台构建和CI/CD配置
├── dist/             # 多平台构建输出目录
├── main.go           # 程序入口点
├── go.mod            # Go 模块定义
├── go.sum            # 依赖校验和
├── Makefile          # 构建脚本（支持多平台构建）
├── .golangci.yml     # 代码质量检查配置
├── .gitignore        # Git 忽略文件
├── LICENSE           # 许可证文件
├── README.md         # 项目说明文档
└── IFLOW.md          # 项目上下文文档
```

## 构建和运行

### 本地构建

```bash
# 构建项目
make build

# 运行项目
make run

# 清理构建产物
make clean

# 运行测试
make test

# 运行测试并生成覆盖率报告
make test-coverage

# 代码格式化
make fmt

# 代码检查
make lint

# 或者直接使用 Go 命令
go build -o macconv .
./macconv --help
```

### 跨平台构建

```bash
# 构建所有支持的平台
make build-all

# 支持的平台：
# - Linux (amd64)
# - Windows (amd64)
# - macOS (amd64, arm64)
```

### 日志级别控制

```bash
# 默认日志级别为 warn
./macconv mac 00:11:22:33:44:55

# 设置日志级别
./macconv -l debug mac 00:11:22:33:44:55
./macconv -l info mac 00:11:22:33:44:55
./macconv -l error mac 00:11:22:33:44:55
```

## 功能命令

### MAC 地址转换

```bash
macconv mac 00:11:22:33:44:55
```

将 MAC 地址转换为多种格式（无分隔符、冒号分隔、点分隔、连字符分隔，以及大写形式）。支持输入验证和错误处理。

**输出示例：**
```
001122334455
00:11:22:33:44:55
0011.2233.4455
0011-2233-4455
001122334455
00:11:22:33:44:55
0011.2233.4455
0011-2233-4455
```

### CIDR 掩码转换

```bash
macconv ip 192.168.1.0/24
macconv ip 2001:db8::/32
```

计算并显示 CIDR 地址范围、子网掩码、反掩码、网络 ID、广播地址和主机数量。支持 IPv4 和 IPv6 地址。

**IPv4 输出示例：**
```
CIDR Address Range: 192.168.1.1 - 192.168.1.254
Subnet Mask: 255.255.255.0
Inverse Mask: 0.0.0.255
Network ID: 192.168.1.0
Broadcast Address: 192.168.1.255
Total Hosts: 254
```

**IPv6 输出示例：**
```
CIDR Address Range: 2001:db8:: - 2001:db8:ffff:ffff:ffff:ffff:ffff:ffff
Subnet Mask: ffff:ffff::
Inverse Mask: ::ffff:ffff:ffff:ffff:ffff:ffff
Network ID: 2001:db8::
Network Type: IPv6
Total Hosts: Very large number (>2^63)
```

### 端口检查

```bash
macconv tcp 192.168.1.1 22
macconv tcp baidu.com 80
```

持续检查指定主机的端口是否开放。支持 IP 地址和域名。端口关闭时持续检查，端口开放时需要连续 5 次检查成功才确认开放。使用 Ctrl+C 可以停止检查。每次检查间隔 1 秒。

**特性：**
- 支持域名解析和 IP 地址
- 显示解析的 IP 地址（仅限域名）
- 开放端口状态和进度在一行显示
- 智能区分 IPv4 和 IPv6 地址

**输出示例：**
```
# 域名检查
2025-12-05T22:05:05+08:00 [1] Port 80 on baidu.com (111.63.65.247) is OPEN ✓ (1/5 consecutive checks)
2025-12-05T22:05:06+08:00 [2] Port 80 on baidu.com (111.63.65.247) is OPEN ✓ (2/5 consecutive checks)
...
2025-12-05T22:05:10+08:00 [5] Port 80 on baidu.com (111.63.65.247) is OPEN ✓ (5/5 consecutive checks) - CONFIRMED

# IP 地址检查
2025-12-05T22:05:13+08:00 [1] Port 22 on 127.0.0.1 is OPEN ✓ (1/5 consecutive checks)
```

### DHCP 选项 43 转换

```bash
macconv dhcp 192.168.1.1
macconv dhcp 192.168.1.1 192.168.1.2
```

将 IPv4 地址转换为 DHCP 选项 43 的 PXE 和 ACS 格式（包括十六进制字符串和字节表示）。支持单个或多个 IP 地址转换。

**输出示例：**
```
PXE Format: 8005000a8000000201a8c001
ACS Format: 010a8000000201a8c001
PXE Format (Bytes): 0x80 0x05 0x00 0x0a 0x80 0x00 0x00 0x02 0x01 0xa8 0xc0 0x01
ACS Format (Bytes): 0x01 0x0a 0x80 0x00 0x00 0x02 0x01 0xa8 0xc0 0x01
```

### 版本信息

```bash
macconv version
./macconv --version
```

显示工具的版本信息、构建日期、作者和联系方式。

## 架构特性

### 错误处理系统
- 统一的错误类型定义（ValidationError, NetworkError, FileSystemError, ParseError）
- 错误包装和链式处理
- 结构化错误信息输出

### 日志记录系统
- 多级别日志支持（debug, info, warn, error）
- 默认日志级别为 warn，减少噪音
- 格式化日志输出，包含时间戳和级别标识

### 验证工具包
- 通用的输入验证函数
- MAC 地址、IP 地址、端口号、CIDR 格式验证
- 文件路径安全验证

### 测试覆盖
- 核心功能的单元测试
- 错误处理和边界条件测试
- IPv4 和 IPv6 兼容性测试

## 开发约定

1. **代码风格**: 遵循 Go 语言标准代码风格
2. **版权信息**: 所有源文件包含版权声明 `Copyright © 2024-2025 Auska <luodan0709@live.cn>`
3. **提交信息**: 使用英文提交信息
4. **测试**: 使用 `go test ./pkg/...` 运行测试
5. **构建**: 使用 `go build` 并通过 `-ldflags="-s -w"` 减小二进制文件大小
6. **代码质量**: 使用 golangci-lint 进行代码检查
7. **日志级别**: 默认使用 warn 级别，避免过多输出

## 版本信息

当前版本: 0.2.0  
作者: LuoDan  
邮箱: luodan0709@live.cn

## 许可证

请参考 LICENSE 文件了解项目的许可证信息。