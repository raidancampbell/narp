package main

import (
	"context"
	"fmt"
	"github.com/mdlayher/arp"
	"net"
	"os"
	"os/signal"
)

func main() {
	ctx := context.Background()
	iFaces, err := net.Interfaces()
	if err != nil {
		panic(err)
	}
	for _, iFace := range iFaces {
		addrs, err := iFace.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.To4() == nil || ip.IsLoopback() {
				continue
			}

			client, err := arp.Dial(&iFace)
			if err != nil {
				fmt.Printf(err.Error())
				continue
			}
			go watchAndNarp(ctx, client, ip, iFace.HardwareAddr)
		}
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	for range c {
		ctx.Done()
	}
}

func watchAndNarp(ctx context.Context, client *arp.Client, ip net.IP, addr net.HardwareAddr) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		packet, _, err := client.Read()
		if err != nil {
			panic(err)
		}

		if packet.SenderIP.Equal(net.IPv4zero) && !packet.TargetIP.Equal(net.IPv4zero) {
			fmt.Printf("ARP Probe requesting if IP '%s' is available\n", packet.TargetIP.String())
			fmt.Printf("Telling '%s' that we (%s) own it...\n", packet.SenderHardwareAddr, ip.String())
			err := client.Reply(packet, addr, ip)
			if err != nil {
				fmt.Printf("ARP probe reply failed with '%s'.  Retrying...", err.Error())
				err := client.Reply(packet, addr, ip)
				if err != nil {
					fmt.Println("retry failed, giving up")
				}
			}
		}
	}
}
