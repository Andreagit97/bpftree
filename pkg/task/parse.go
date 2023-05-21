package task

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/cilium/ebpf"
	"github.com/cilium/ebpf/features"
	"github.com/cilium/ebpf/link"
	"github.com/cilium/ebpf/rlimit"
	"golang.org/x/sys/unix"
)

//go:generate go run github.com/cilium/ebpf/cmd/bpf2go iter ./bpf/iter.bpf.c

const (
	// the magic number is on 2 bytes.
	magicLittleEndian uint16 = 0xeb9f
	magicBigEndian    uint16 = 0x9feb
)

/* If we want to read taskInfo from a file. */
var (
	captureFilePath                  = ""
	byteOrder       binary.ByteOrder = binary.LittleEndian
)

// SetCaptureFilePath is used to set the capture file.
func SetCaptureFilePath(filepath string) {
	captureFilePath = filepath
}

func decodeByteEndianness(reader io.ReadCloser, numBytes uint32, data any) error {
	rawBytes := make([]byte, numBytes)
	readBytes, err := reader.Read(rawBytes)
	if errors.Is(err, io.EOF) {
		/* The caller should check for end of file */
		return io.EOF
	}
	if err != nil {
		displayError(fmt.Sprintf("cannot read '%d' bytes: %v", numBytes, err))
		return err
	}
	if uint32(readBytes) != numBytes {
		err = fmt.Errorf("read '%d' bytes instead of '%d'", readBytes, numBytes)
		return err
	}

	byteBuffer := bytes.NewBuffer(rawBytes)
	err = binary.Read(byteBuffer, byteOrder, data)
	if err != nil {
		displayError(fmt.Sprintf("cannot decode '%d' bytes: %v", numBytes, err))
		return err
	}
	return nil
}

func parseMagicNumber(reader io.ReadCloser) error {
	var magic uint16
	if err := decodeByteEndianness(reader, 2, &magic); err != nil {
		return err
	}

	// Swap endianness if necessary.
	if magic != magicLittleEndian {
		if magic == magicBigEndian {
			byteOrder = binary.BigEndian
		} else {
			return fmt.Errorf("the provided file doesn't have the right format")
		}
	}
	return nil
}

func bpfDetectionFeatures() error {
	/* Try here if we support all necessary BPF helpers and functionalities. */
	if err := features.HaveProgramType(ebpf.Tracing); err != nil {
		switch {
		case errors.Is(err, ebpf.ErrNotSupported):
			displayError("BPF iterators are not supported",
				"by the running kernel (requires >= v5.8).")
		case errors.Is(err, unix.EPERM):
			displayError("Not enough privileges, retry with sudo",
				"or with {CAP_BPF,CAP_PERFMON,CAP_SYS_RESOURCE}.")
		default:
			displayError("error while probing Iter BPF progs:", err)
		}
		return err
	}

	/* Allow the current process to lock memory for eBPF resources. */
	if err := rlimit.RemoveMemlock(); err != nil {
		displayError("unable to bump RLIMIT_MEMLOCK:", err)
		return err
	}

	return nil
}

func openBPFIterator() (io.ReadCloser, error) {
	/* Detect necessary BPF features */
	if err := bpfDetectionFeatures(); err != nil {
		return nil, err
	}

	/* Load pre-compiled programs and maps into the kernel. */
	objs := iterObjects{}
	if err := loadIterObjects(&objs, nil); err != nil {
		displayError("unable to load BPF objects into the kernel:", err)
		return nil, err
	}
	defer objs.Close()

	opts := link.IterOptions{
		Program: objs.DumpTask,
	}

	iter, err := link.AttachIter(opts)
	if err != nil {
		displayError("unable to attach Iter prog:", err)
		return nil, err
	}

	reader, err := iter.Open()
	if err != nil {
		displayError("unable to open a reader for the iter:", err)
		return nil, err
	}
	return reader, nil
}

func openCaptureFile() (io.ReadCloser, error) {
	reader, err := os.Open(captureFilePath)
	if err != nil {
		displayError(fmt.Sprintf("unable to open file '%s': %v", captureFilePath, err))
		return nil, err
	}
	return reader, nil
}

// PopulateTaskInfo populate the thread table reading from a capture file
// or using the BPF iterator.
func PopulateTaskInfo() error {
	var err error
	var reader io.ReadCloser
	if captureFilePath == "" {
		/* We need to read from the system with BPF. */
		reader, err = openBPFIterator()
	} else {
		/* We can use the provided file. */
		reader, err = openCaptureFile()
	}

	if err != nil {
		return err
	}

	defer reader.Close()

	/* 1. Parse Magic number */
	if err := parseMagicNumber(reader); err != nil {
		displayError("cannot parse the magic number:", err)
		return err
	}

	/* 2. Parse Header len */
	if err := parseHeaderLen(reader); err != nil {
		displayError("cannot parse the header len:", err)
		return err
	}

	/* 3. Parse Header */
	if err := parseHeader(reader); err != nil {
		displayError("cannot parse the header:", err)
		return err
	}

	/* 4. Parse task Info. */
	if err := parseTaskInfos(reader); err != nil {
		displayError("unable to parse task info:", err)
		return err
	}

	/* Order the list of tasks by tid */
	sort.SliceStable(tasksList, func(i, j int) bool {
		return tasksList[i].Info.Tid < tasksList[j].Info.Tid
	})

	/* Compute children and order them by tid */
	computeChildren()
	return nil
}
