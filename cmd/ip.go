/*
Copyright © 2024 LuoDan(luodan0709@live.cn)

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

	"github.com/seancfoley/ipaddress-go/ipaddr"
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
		cidr(args[0])
	}
}
func cidr(str string) {
	addrString := ipaddr.NewIPAddressString(str)
	addr := addrString.GetAddress()
	if addr == nil {
		fmt.Printf("Invalid IP address: %s\n", str)
		return
	}
	version := addrString.GetIPVersion()
	segments := addr.GetSegments()
	bitLength := addr.GetBitCount()
	size := addr.GetCount()
	prefixLen := addr.GetNetworkPrefixLen()
	mask := addr.GetNetworkMask()

	// three different ways to get the network address
	networkAddr, _ := addr.Mask(mask)
	networkAddr = networkAddr.WithoutPrefixLen()
	networkAddrAnotherWay := addr.GetLower().WithoutPrefixLen()
	zeroHost, _ := addr.ToZeroHost()
	networkAddrOneMoreWay := zeroHost.WithoutPrefixLen()

	fmt.Printf("%v address: %v\nprefix length: %v\nbit length: %v\nsegments: %v\n"+
		"size: %v\nnetwork mask: %v\nnetwork address: %v\n\n",
		version, addr, prefixLen, bitLength, segments, size, mask, networkAddr)

	_, _ = networkAddrAnotherWay, networkAddrOneMoreWay
}
