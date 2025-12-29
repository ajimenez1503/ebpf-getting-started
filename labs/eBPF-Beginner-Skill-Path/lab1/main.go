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

	log.Println("tracepoint attached: sys_enter_execve")
	log.Println("run `sudo cat /sys/kernel/debug/tracing/trace_pipe` in another terminal to see printk output")

	// Wait for interrupt so the program stays alive.
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	log.Println("shutting down")
}

