package main

import (
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/cilium/ebpf/link"
	"github.com/cilium/ebpf/rlimit"
)

func main() {
	if err := rlimit.RemoveMemlock(); err != nil {
		log.Fatal("removing memlock:", err)
	}

	ifaceName := flag.String("iface", "eth0", "network interface to attach XDP")
	flag.Parse()

	iface, err := net.InterfaceByName(*ifaceName)
	if err != nil {
		log.Fatalf("getting interface %s: %v", *ifaceName, err)
	}

	var objs counterObjects
	if err := loadCounterObjects(&objs, nil); err != nil {
		log.Fatal("loading eBPF objects:", err)
	}
	defer objs.Close()

	l, err := link.AttachXDP(link.XDPOptions{
		Program:   objs.CountPackets,
		Interface: iface.Index,
	})
	if err != nil {
		log.Fatal("attaching XDP:", err)
	}
	defer l.Close()

	log.Printf("Counting incoming packets on %s..", iface.Name)
	tick := time.Tick(time.Second)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	for {
		select {
		case <-tick:
			var count uint64
			if err := objs.PktCount.Lookup(uint32(0), &count); err != nil {
				log.Fatal("map lookup:", err)
			}
			log.Printf("Received %d packets", count)
		case <-stop:
			log.Print("Received signal, exiting..")
			return
		}
	}
}

