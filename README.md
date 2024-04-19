## MAC地址格式转换

用法：`macconv mac 00:00:00:00:00:00`

```
Used to convert mac addresses between different devices.
For example:
        macconv mac 00:11:22:33:44:55
        macconv ip 192.168.1.1/24
        macconv tcp 192.168.1.1 22

Usage:
  macconv [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  ip          CIDR mask conversion
  mac         Convert mac address
  tcp         Check host port
  version     Print version.

Flags:
  -h, --help     help for macconv
  -t, --toggle   Help message for toggle

Use "macconv [command] --help" for more information about a command.
```
