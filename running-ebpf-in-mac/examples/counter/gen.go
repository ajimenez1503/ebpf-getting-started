//go:generate go run github.com/cilium/ebpf/cmd/bpf2go counter ./counter.c -- -nostdinc -I/usr/include -I/usr/include/aarch64-linux-gnu -I/usr/include/x86_64-linux-gnu -O2 -g
package main

