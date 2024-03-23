package iterators

import (
	"errors"
	"fmt"
	"io"

	"github.com/cilium/ebpf"
	"github.com/cilium/ebpf/asm"
	"github.com/cilium/ebpf/link"
	"github.com/cilium/ebpf/rlimit"
	"golang.org/x/sys/unix"
)

//go:generate go run github.com/cilium/ebpf/cmd/bpf2go iter ./bpf/iter.bpf.c

func eBPFIteratorSupport() error {
	// Check if the kernel supports BPF iterators.
	progSpec := &ebpf.ProgramSpec{
		Type:       ebpf.Tracing,
		AttachType: ebpf.AttachTraceIter,
		AttachTo:   "task",
		Instructions: asm.Instructions{
			asm.LoadImm(asm.R0, 0, asm.DWord),
			asm.Return(),
		},
	}

	prog, err := ebpf.NewProgramWithOptions(progSpec, ebpf.ProgramOptions{
		LogDisabled: true,
	})

	if err == nil {
		if closeErr := prog.Close(); closeErr != nil {
			return fmt.Errorf("failed to close test ebpf program: %w", closeErr)
		}
		return nil
	}

	switch {
	// EINVAL occurs when attempting to create a program with an unknown type.
	// E2BIG occurs when ProgLoadAttr contains non-zero bytes past the end
	// of the struct known by the running kernel, meaning the kernel is too old
	// to support the given prog type.
	case errors.Is(err, unix.EINVAL), errors.Is(err, unix.E2BIG):
		return fmt.Errorf("eBPF iterators are not supported by the running kernel")
	case errors.Is(err, unix.EPERM):
		return fmt.Errorf("not enough privileges, retry with sudo or with {CAP_BPF,CAP_PERFMON,CAP_SYS_RESOURCE}")
	}
	return fmt.Errorf("error while probing BPF Iterators in the kernel. %v", err)
}

func checkeBPFSupport() error {
	// Allow the current process to lock memory for eBPF resources.
	if err := rlimit.RemoveMemlock(); err != nil {
		return fmt.Errorf("failed to bump memlock rlimit: %w", err)
	}

	// Detect necessary BPF features
	if err := eBPFIteratorSupport(); err != nil {
		return fmt.Errorf("eBPF iterators are not supported by the running kernel: %w", err)
	}

	return nil
}

// GetTasksReader returns a reader for the task iterator.
func GetTasksReader() (io.ReadCloser, error) {
	if err := checkeBPFSupport(); err != nil {
		return nil, err
	}

	// Load pre-compiled programs and maps into the kernel.
	objs := iterObjects{}
	if err := loadIterObjects(&objs, nil); err != nil {
		return nil, fmt.Errorf("unable to load BPF objects into the kernel: %w", err)
	}
	defer func() {
		_ = objs.Close()
	}()

	// open taskInfo iterator
	opts := link.IterOptions{
		Program: objs.DumpTask,
	}

	iter, err := link.AttachIter(opts)
	if err != nil {
		return nil, fmt.Errorf("unable to attach taskInfo iter: %w", err)
	}

	reader, err := iter.Open()
	if err != nil {
		return nil, fmt.Errorf("unable to open a reader for taskInfo iter: %w", err)
	}
	return reader, nil
}

// GetFilesReader returns a reader for the file iterator.
func GetFilesReader() (io.ReadCloser, error) {
	if err := checkeBPFSupport(); err != nil {
		return nil, err
	}

	// Load pre-compiled programs and maps into the kernel.
	objs := iterObjects{}
	if err := loadIterObjects(&objs, nil); err != nil {
		return nil, fmt.Errorf("unable to load BPF objects into the kernel: %w", err)
	}
	defer func() {
		_ = objs.Close()
	}()

	// Open fileInfo iterator
	opts := link.IterOptions{
		Program: objs.DumpTaskFile,
	}

	iter, err := link.AttachIter(opts)
	if err != nil {
		return nil, fmt.Errorf("unable to attach fileInfo iter: %w", err)
	}

	reader, err := iter.Open()
	if err != nil {
		return nil, fmt.Errorf("unable to open a reader for fileInfo iter: %w", err)
	}
	return reader, nil
}
