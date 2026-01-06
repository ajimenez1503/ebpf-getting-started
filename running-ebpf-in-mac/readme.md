# Run eBPF on macOS (Apple Silicon or Intel)

macOS cannot run eBPF programs directly because eBPF is a Linux kernel feature. The simplest workflow is to develop inside a lightweight Linux VM. Below is a minimal, repeatable setup using Lima (works well on Apple Silicon and Intel) with the required toolchain installed automatically.

## Prerequisites

- Homebrew installed on macOS.
- ~10 GB free disk space for the VM image.

## 1) Install Lima

```bash
brew install lima
```

## 2) Create a Lima config in your project

Save this as `lima.yaml` in the project root:

```yaml
images:
  - location: "https://cloud-images.ubuntu.com/releases/24.04/release-20240423/ubuntu-24.04-server-cloudimg-arm64.img"
    arch: aarch64

mounts:
  - location: "~"
    writable: true  # optional; set false if you prefer read-only

provision:
  - mode: system
    script: |
      apt-get update
      apt-get install -y apt-transport-https ca-certificates curl clang llvm jq
      apt-get install -y libelf-dev libpcap-dev libbfd-dev binutils-dev build-essential make
      apt-get install -y linux-headers-$(uname -r)
      apt-get install -y linux-tools-common linux-tools-$(uname -r)
      apt-get install -y bpfcc-tools libbpf-dev
      apt-get install -y golang-go git pkg-config python3-pip
```

## 3) Start the VM

```bash
limactl start --name ebpf lima.yaml
```

This downloads the Ubuntu image, creates the VM, mounts your home directory, and installs the toolchain via the provision script.

## 4) Enter the VM shell

```bash
limactl shell ebpf
```

You are now inside Linux with `clang`, `llvm`, `libbpf-dev`, `bpftool`, and `bpftrace` available. Your macOS home directory is mounted, so you can edit code on macOS and build/run inside the VM.

## 5) Quick sanity checks inside the VM

```bash
uname -a          # confirm Linux kernel
bpftool version   # from linux-tools
bpftrace -V       # from bpfcc-tools
go version        # Go toolchain for cilium/ebpf programs
```

## 6) Build and run your eBPF programs

- Compile with `clang -target bpf` (or use `libbpf`/`bpftool gen skeleton` as usual).
- Run loaders/test harnesses from inside the VM so they can attach to the Linux kernel.
- For Go-based loaders or tooling using `cilium/ebpf`, the provisioned `golang-go`, `git`, `pkg-config`, and kernel headers support building and running Go eBPF programs. See [ebpf-go docs](https://ebpf-go.dev/).

Minimum requirements for Go eBPF examples with `cilium/ebpf` ([ebpf-go docs](https://ebpf-go.dev/)):

- Linux kernel 5.7+ (for `bpf_link` support).
- LLVM/Clang 11+ (`clang`, `llvm-strip`).
- `libbpf` headers and `bpftrace` (installed via `bpftrace` package).
- Linux kernel headers.
- A Go toolchain version supported by the `cilium/ebpf` module (match the module’s `go` directive).

## Notes and alternatives

- If you want a separate source directory instead of mounting your whole home folder, adjust the `mounts` section accordingly.
- The Medium guide covers a similar VM-based approach on macOS; using Lima keeps things lightweight and automatable, but any Linux VM with the same packages (clang/llvm, libelf-dev, libbpf-dev, bpfcc-tools, linux-tools) will work. See [Medium: Setting up eBPF in macOS – a beginner’s guide](https://medium.com/@kalaiarasanbalaraman/setting-up-ebpf-in-macos-a-beginners-guide-42c59182b41f).

With this setup, you can iterate on eBPF programs from macOS while compiling and running against a real Linux kernel.
