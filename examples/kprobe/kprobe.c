//go:build ignore
// +build ignore

#include <linux/bpf.h>
#include <bpf/bpf_helpers.h>

SEC("kprobe/execve")
int kprobe_execve(struct pt_regs *ctx)
{
    bpf_printk("execve called\n");
    return 0;
}

char LICENSE[] SEC("license") = "Dual MIT/GPL";

