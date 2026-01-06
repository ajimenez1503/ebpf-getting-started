package main

import (
	"bytes"
	_ "embed"
	"flag"
	"log"
	"os"
	"os/signal"

	"github.com/cilium/ebpf"
	"github.com/cilium/ebpf/link"
	"github.com/cilium/ebpf/rlimit"
)

//go:generate clang -O2 -g -target bpf -c kprobe.c -o kprobe.o -I/usr/include/aarch64-linux-gnu -I/usr/include/asm-generic

//go:embed kprobe.o
var bpfProgram []byte

func main() {
	if err := rlimit.RemoveMemlock(); err != nil {
		log.Fatalf("removing memlock limit: %v", err)
	}

	symbol := flag.String("symbol", "__arm64_sys_execve", "kprobe target symbol (e.g. __arm64_sys_execve, __x64_sys_execve, do_execveat_common)")
	flag.Parse()

	spec, err := ebpf.LoadCollectionSpecFromReader(bytes.NewReader(bpfProgram))
	if err != nil {
		log.Fatalf("loading collection spec: %v", err)
	}

	var objs struct {
		KprobeExecve *ebpf.Program `ebpf:"kprobe_execve"`
	}
	if err := spec.LoadAndAssign(&objs, nil); err != nil {
		log.Fatalf("loading objects: %v", err)
	}
	defer objs.KprobeExecve.Close()

	kp, err := link.Kprobe(*symbol, objs.KprobeExecve, nil)
	if err != nil {
		log.Fatalf("attaching kprobe: %v", err)
	}
	defer kp.Close()

	log.Printf("kprobe attached to %s; check trace_pipe for logs (sudo cat /sys/kernel/debug/tracing/trace_pipe)", *symbol)

	// Wait for SIGINT/SIGTERM
	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, os.Interrupt)
	<-sigch
	log.Println("signal received, exiting")
}
