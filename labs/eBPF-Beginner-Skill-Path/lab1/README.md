# Lab 1 — Hello World eBPF (tracepoint `sys_enter_execve`)

This folder mirrors the iximiuz “From Zero to Your First eBPF Program” lab, wired to run locally or in the provided playground.

## Files
- `main.go` — user-space loader that generates/loads eBPF bytecode and attaches it to the `sys_enter_execve` tracepoint.
- `hello.c` — kernel-space eBPF program printing “Hello world” via `bpf_printk`.
- `go.mod` — minimal module definition so `go generate` / `go build` work in isolation.
- `vmlinux.h` — needs to be generated once (see below). Not committed due to size/host specificity.

## Prereqs (playground already satisfies these)
- Go toolchain (>=1.22)
- `clang`/`llc`
- `bpftool`
- Linux kernel headers with BTF (for `vmlinux.h`)

## One-time setup (generate `vmlinux.h`)
```bash
bpftool btf dump file /sys/kernel/btf/vmlinux format c > vmlinux.h
```

## Build steps
```bash
# Inside labs/eBPF-Beginner-Skill-Path/lab1
go generate ./...   # runs bpf2go to compile hello.c -> hello_bpf.* and hello_bpf.o
go build ./...      # produces ./lab1 binary
```

## Run
```bash
sudo ./lab1
```

In another terminal, view printk output:
```bash
sudo cat /sys/kernel/debug/tracing/trace_pipe
```

Trigger an event so the tracepoint fires (e.g., run `uname -a`), then look for “Hello world from eBPF!” lines in `trace_pipe`.

## Cleanup
```bash
sudo pkill lab1 || true
```

