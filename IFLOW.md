# MAC地址格式转换工具 (macconv)

## 项目概述

macconv 是一个用 Go 语言编写的命令行工具，用于网络地址格式转换和常见网络操作。该工具基于 Cobra 框架构建，提供了多种网络相关的实用功能，包括 MAC 地址格式转换、CIDR 掩码转换、端口检查、DHCP 选项 43 转换等。

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
│   ├── root.go       # 根命令定义
│   ├── mac.go        # MAC 地址转换命令
│   ├── ip.go         # CIDR 掩码转换命令
│   ├── tcp.go        # 端口检查命令
│   ├── dhcp.go       # DHCP 选项 43 转换命令
│   └── version.go    # 版本信息命令
├── pkg/              # 共享包目录
│   ├── errors/       # 错误处理包
│   ├── logger/       # 日志记录包
│   └── validator/    # 验证工具包
├── .github/          # GitHub Actions 工作流
│   └── workflows/
│       └── test.yml  # 自动构建和测试配置
├── main.go           # 程序入口点
├── go.mod            # Go 模块定义
├── go.sum            # 依赖校验和
├── Makefile          # 构建脚本
├── LICENSE           # 许可证文件
└── README.md         # 项目说明文档
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

# 或者直接使用 Go 命令
go build -o macconv .
./macconv --help
```

### 跨平台构建

项目支持通过 GitHub Actions 进行跨平台构建，自动生成 Windows 和 Linux 平台的可执行文件。

## 功能命令

### MAC 地址转换

```bash
macconv mac 00:11:22:33:44:55
```

将 MAC 地址转换为多种格式（无分隔符、冒号分隔、点分隔、连字符分隔，以及大写形式）。

### CIDR 掩码转换

```bash
macconv ip 192.168.1.1/24
```

计算并显示 CIDR 地址范围、子网掩码、反掩码、网络 ID、广播地址和主机数量。

### 端口检查

```bash
macconv tcp 192.168.1.1 22
macconv tcp baidu.com 80
```

持续检查指定主机的端口是否开放。支持 IP 地址和域名。端口关闭时持续检查，端口开放时需要连续 5 次检查成功才确认开放。使用 Ctrl+C 可以停止检查。每次检查间隔 1 秒。

### DHCP 选项 43 转换

```bash
macconv dhcp 192.168.1.1
```

将 IPv4 地址转换为 DHCP 选项 43 的 PXE 和 ACS 格式（包括十六进制字符串和字节表示）。

### 版本信息

```bash
macconv version
```

显示工具的版本信息、作者和联系方式。

## 开发约定

1. **代码风格**: 遵循 Go 语言标准代码风格
2. **版权信息**: 所有源文件包含版权声明 `Copyright © 2024-2025 Auska <luodan0709@live.cn>`
3. **提交信息**: 使用英文提交信息
4. **测试**: 使用 `go test ./...` 运行测试
5. **构建**: 使用 `go build` 并通过 `-ldflags="-s -w"` 减小二进制文件大小

## 版本信息

当前版本: 0.1.3  
作者: LuoDan  
邮箱: luodan0709@live.cn

## 许可证

请参考 LICENSE 文件了解项目的许可证信息。