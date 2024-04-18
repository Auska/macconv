/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
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
		cidr(args[0])
	}
}
func cidr(str string) {
	addrString := ipaddr.NewIPAddressString(str)
	addr := addrString.GetAddress()
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
