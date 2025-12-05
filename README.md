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
  tcp         Check host port
  version     Print version.

Flags:
  -h, --help               help for macconv
  -l, --log-level string   Set log level (debug, info, warn, error) (default "warn")
  -t, --toggle             Help message for toggle

Use "macconv [command] --help" for more information about a command.
```
