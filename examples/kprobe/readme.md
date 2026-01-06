# kprobe example (based on cilium/ebpf examples/kprobe)

Reference: <https://github.com/cilium/ebpf/tree/main/examples/kprobe>

This example attaches a kprobe to an exec syscall and logs events from the kernel to user space. On arm64 the symbol is typically `__arm64_sys_execve`; on x86_64 it is `__x64_sys_execve`. Override with `-symbol` if needed.

```bash
cd examples/kprobe
go mod tidy          # once
go generate ./...
go build
sudo ./kprobe -symbol __arm64_sys_execve   # for arm64
# for x86_64 use: sudo ./kprobe -symbol __x64_sys_execve
# in another shell, watch logs:
sudo cat /sys/kernel/tracing/trace_pipe
```

Notes:

- Requires kernel headers, libbpf, clang/llvm, and Go (already provisioned in the Lima setup).

- If you see header include errors for `asm/types.h`, create a compat symlink (arm64): `sudo ln -s /usr/include/aarch64-linux-gnu/asm /usr/include/asm`.
