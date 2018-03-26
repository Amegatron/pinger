package main

import (
	"fmt"
	"github.com/Amegatron/pinger/config"
	"github.com/tatsushid/go-fastping"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

type RemoteInfo struct {
	domain      string
	addr        *net.IPAddr
	lastSuccess time.Time
	failure     bool
}

var remotes map[string]*RemoteInfo

func init() {
	remotes = make(map[string]*RemoteInfo)
	for _, remote := range config.Remotes() {
		ra, err := resolveRemote(remote)

		if err != nil {
			fmt.Println("Could not resolve", remote, ":", err)
		} else {
			remotes[ra.String()] = &RemoteInfo{remote, ra, time.Now(), false}
		}
	}
}

func main() {
	pinger := fastping.NewPinger()
	pinger.OnRecv = func(addr *net.IPAddr, duration time.Duration) {
		remote := remotes[addr.String()]
		remote.lastSuccess = time.Now()
		if remote.failure == true {
			remote.failure = false
			// Success event
		}
		fmt.Println(addr.IP, "OK", duration)
	}
	pinger.OnIdle = func() {
		for _, remote := range remotes {
			if remote.failure {
				continue
			}

			if time.Since(remote.lastSuccess) > time.Duration(time.Second*5) {
				fmt.Println(remote.domain, "is unreachable")
				remote.failure = true
				// Failure event
			}
		}
	}
	pinger.MaxRTT = time.Second * 2
	for _, remoteInfo := range remotes {
		pinger.AddIPAddr(remoteInfo.addr)
	}
	pinger.RunLoop()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)

loop:
	for {
		select {
		case <-c:
			fmt.Println("Got interrupted")
			break loop
		case <-time.After(time.Second):
			break
		}
	}
	signal.Stop(c)
	fmt.Println("Finished")
}

func resolveRemote(remote string) (*net.IPAddr, error) {
	netProto := "ip4:icmp"
	if strings.Index(remote, ":") != -1 {
		netProto = "ip6:ipv6-icmp"
	}
	ra, err := net.ResolveIPAddr(netProto, remote)
	return ra, err
}
