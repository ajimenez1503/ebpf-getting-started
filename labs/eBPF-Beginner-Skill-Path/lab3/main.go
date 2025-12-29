package main

//go:generate go run github.com/cilium/ebpf/cmd/bpf2go -target bpf hello hello.c -- -I./

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/cilium/ebpf/link"
	"github.com/cilium/ebpf/rlimit"
)

func main() {
	if err := rlimit.RemoveMemlock(); err != nil {
		log.Fatalf("failed to adjust rlimit: %v", err)
	}

	objs := helloObjects{}
	if err := loadHelloObjects(&objs, nil); err != nil {
		log.Fatalf("loading BPF objects: %v", err)
	}
	defer objs.Close()

	tp, err := link.Tracepoint("syscalls", "sys_enter_execve", objs.HandleExecveTp, nil)
	if err != nil {
		log.Fatalf("attaching tracepoint: %v", err)
	}
	defer tp.Close()

	log.Println("lab3 running: tracepoint attached (sys_enter_execve)")
	log.Println("Inspect with: sudo bpftool prog list && sudo bpftool map list")
	log.Println("Dump exec_count: sudo bpftool map dump id <ID>")
	log.Println("Monitor runtime: sudo bpftop (Ctrl+C to exit bpftop)")
	log.Println("press Ctrl+C here to exit the program")

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	log.Println("shutting down")
}

