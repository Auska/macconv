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

	"github.com/spf13/cobra"
)

// ipCmd represents the ip command
var ipCmd = &cobra.Command{
	Use:   "ip",
	Short: "CIDR mask conversion",
	Long: `CIDR mask conversion. For example:

macconv ip 192.168.1.1/24`,
	Run: convertIPAddress(),
}

func init() {
	rootCmd.AddCommand(ipCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// ipCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// ipCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func convertIPAddress() func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Printf("missing arguments")
			return
		}
		first, last, mask, err := calculateCIDRInfo(args[0])
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		fmt.Println("CIDR Address Range:", first, "-", last)
		fmt.Println("Subnet Mask:", mask)
		fmt.Println("Network ID:", first)
	}
}
func calculateCIDRInfo(cidr string) (string, string, string, error) {
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return "", "", "", err
	}

	// 计算网络号
	network := ip.Mask(ipnet.Mask)

	// 计算地址段
	first := network
	last := net.IP(make([]byte, len(network)))
	copy(last, network)
	for i := range last {
		last[i] |= ^ipnet.Mask[i]
	}

	// 转换掩码为字符串
	mask := net.IP(ipnet.Mask).String()

	return first.String(), last.String(), mask, nil
}
