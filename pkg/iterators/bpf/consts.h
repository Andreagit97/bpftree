#pragma once

// Extracted from the kernel
#define MAX_PID_NS_LEVEL 32
#define MS_NOUSER 1 << 31
#define PPM_OVERLAYFS_SUPER_MAGIC 0x794c7630

// Limits
// todo!: some of these definitions should be shared with the userspace
#define UINT32_MAX (4294967295U)
#define EXE_PATH_MAX_LEN 512
#define CMDLINE_MAX_LEN 512
#define FILE_PATH_MAX_LEN 512
#define SAFE_1024_ACCESS(x) x &(512 - 1)

// "TL" stands for "Too Long"
#define TOO_LONG "TL"
#define TOO_LONG_SIZE 3

// "NA" stands for "not available"
#define NOT_AVAILABLE "NA"
#define NOT_AVAILABLE_SIZE 3

// File flags
#define MOUNTABLE 1 << 0
#define UPPER_LAYER 1 << 1
#define LOWER_LAYER 1 << 2
