## MAC地址格式转换

用法：`macconv mac 00:00:00:00:00:00`

```
Used to convert mac addresses between different devices.
For example:
        macconv mac 00:11:22:33:44:55
        macconv ip 192.168.1.1/24
        macconv tcp 192.168.1.1 22
        macconv dhcp 192.168.1.1

Usage:
  macconv [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  dhcp        DHCP option 43 conversion
  help        Help about any command
  ### CIDR 掩码转换

```bash
macconv ip 192.168.1.1/24
```

计算并显示 CIDR 地址范围、子网掩码、反掩码、网络 ID、广播地址和主机数量。
  mac         Convert mac address
  ### 端口检查

```bash
macconv tcp 192.168.1.1 22
macconv tcp baidu.com 80
```

持续检查指定主机的端口是否开放。支持 IP 地址和域名。端口关闭时持续检查，端口开放时需要连续 5 次检查成功才确认开放。使用 Ctrl+C 可以停止检查。每次检查间隔 1 秒。
  version     Print version.

Flags:
  -h, --help               help for macconv
  -l, --log-level string   Set log level (debug, info, warn, error) (default "warn")
  -t, --toggle             Help message for toggle

Use "macconv [command] --help" for more information about a command.
```
