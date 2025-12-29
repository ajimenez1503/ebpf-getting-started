//go:build ignore

#include "vmlinux.h"
#include <bpf/bpf_helpers.h>

SEC("tracepoint/syscalls/sys_enter_execve")
int handle_execve_tp(struct trace_event_raw_sys_enter *ctx) {
	/* This prints to /sys/kernel/debug/tracing/trace_pipe */
	bpf_printk("Hello world from eBPF!");
	return 0;
}

/* Map program sections to names expected by Go side. */
char _license[] SEC("license") = "GPL";

