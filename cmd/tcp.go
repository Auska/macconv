/*
Copyright © 2024 LuoDan<luodan0709@live.cn>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/spf13/cobra"
)

// tcpCmd represents the tcp command
var tcpCmd = &cobra.Command{
	Use:   "tcp",
	Short: "Check host port",
	Long: `Check if the host port is open. For example:

macconv tcp 192.168.1.1 22`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			fmt.Printf("missing arguments")
			return
		}
		checkPort(args[0], args[1])
	},
}

func init() {
	rootCmd.AddCommand(tcpCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// tcpCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// tcpCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
func checkPort(ipStr, portStr string) {

	if net.ParseIP(ipStr) == nil {
		fmt.Printf("invalid IP address")
		return
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		fmt.Printf("invalid port")
		return
	}

	var target string
	if net.ParseIP(ipStr).To4() == nil {
		target = fmt.Sprintf("%s:%d", ipStr, port) // 默认使用IPv4格式
	} else {
		target = fmt.Sprintf("[%s]:%d", ipStr, port) // 如果是IPv6，则调整格式
	}
	count := 0
	for count < 5 {
		// 尝试连接到目标主机的指定端口
		now := time.Now()
		conn, err := net.DialTimeout("tcp", target, 2*time.Second)
		if err != nil {
			fmt.Printf("%v Port %d on %s is close\n", now.Format(time.RFC3339), port, ipStr)
		} else {
			fmt.Printf("%v Port %d on %s is open\n", now.Format(time.RFC3339), port, ipStr)
			count++
			conn.Close()
		}
		time.Sleep(time.Second)

	}
}
