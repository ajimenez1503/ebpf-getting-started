# ebpf-go XDP counter example

Source: [ebpf-go getting started guide](https://ebpf-go.dev/guides/getting-started/#ebpf-c-program).

## Requirements (inside the VM)

- Linux kernel 5.7+ (for `bpf_link` support)
- LLVM/Clang 11+ (`clang`, `llvm-strip`)
- libbpf headers and kernel headers
- Go toolchain (version supported by `github.com/cilium/ebpf`)
- Module dependency:

```bash
go get github.com/cilium/ebpf
```

## Files

`counter.c` (eBPF program)

```c
//go:build ignore
#include <linux/bpf.h>
#include <bpf/bpf_helpers.h>

struct {
    __uint(type, BPF_MAP_TYPE_ARRAY);
    __type(key, __u32);
    __type(value, __u64);
    __uint(max_entries, 1);
} pkt_count SEC(".maps");

SEC("xdp")
int count_packets()
{
    __u32 key = 0;
    __u64 *count = bpf_map_lookup_elem(&pkt_count, &key);
    if (count) {
        __sync_fetch_and_add(count, 1);
    }
    return XDP_PASS;
}

char __license[] SEC("license") = "Dual MIT/GPL";
```

`gen.go` (bpf2go scaffolding)

```go
//go:generate go run github.com/cilium/ebpf/cmd/bpf2go counter ./counter.c -- -nostdinc -I/usr/include -O2 -g
package main
```

`main.go` (loader)

```go
package main

import (
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
        log.Fatal("Removing memlock:", err)
    }

    var objs counterObjects
    if err := loadCounterObjects(&objs, nil); err != nil {
        log.Fatal("Loading eBPF objects:", err)
    }
    defer objs.Close()

    iface, err := net.InterfaceByName("eth0") // set to a real interface
    if err != nil {
        log.Fatalf("Getting interface: %v", err)
    }

    l, err := link.AttachXDP(link.XDPOptions{
        Program:   objs.CountPackets,
        Interface: iface.Index,
    })
    if err != nil {
        log.Fatal("Attaching XDP:", err)
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
                log.Fatal("Map lookup:", err)
            }
            log.Printf("Received %d packets", count)
        case <-stop:
            log.Print("Received signal, exiting..")
            return
        }
    }
}
```

## Build and run (inside the VM)

```bash
cd running-ebpf-in-mac/examples/counter
go mod tidy
go generate ./...
go build ./...
sudo ./counter
```
