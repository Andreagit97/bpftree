// go:build ignore
#include "vmlinux.h"
#include <bpf/bpf_helpers.h>
#include <bpf/bpf_core_read.h>
#include <bpf/bpf_tracing.h>

#define FILE_PATH_MAX_LEN 1024

char _license[] SEC("license") = "GPL";

struct exported_file_info
{
	char file_path[FILE_PATH_MAX_LEN];
	uint64_t ino;
	uint32_t fd;
	uint32_t pad;
} typedef exported_file_info;

exported_file_info data;

volatile const int target_thread_id;

SEC("iter/task_file")
int dump_task_file(struct bpf_iter__task_file *ctx)
{
	struct seq_file *seq = ctx->meta->seq;
	struct task_struct *task = ctx->task;
	struct file *file = ctx->file;
	if(task == NULL || file == NULL)
	{
		return 0;
	}

	if(task->pid != target_thread_id)
	{
		return 0;
	}

	data.fd = ctx->fd;
	data.ino = file->f_inode->i_ino;

	if(bpf_core_enum_value_exists(enum bpf_func_id, BPF_FUNC_d_path) &&
	   (bpf_core_enum_value(enum bpf_func_id, BPF_FUNC_d_path) == BPF_FUNC_d_path))
	{
		// According to the manual if the path is too long the helper
		// doesn't populate the path so we return a too-long.
		if(bpf_d_path(&(file->f_path), data.file_path, FILE_PATH_MAX_LEN) < 0)
		{
			// The path was too long `TL`
			data.file_path[0] = 'T';
			data.file_path[1] = 'L';
			data.file_path[2] = '\0';
		}
	}
	else
	{
		// The path is not available `NA`, we cannot recover it!
		data.file_path[0] = 'N';
		data.file_path[1] = 'A';
		data.file_path[2] = '\0';
	}

	bpf_seq_write(seq, &data, sizeof(data));
	return 0;
}
