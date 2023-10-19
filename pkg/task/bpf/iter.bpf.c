// go:build ignore
#include "vmlinux.h"
#include <bpf/bpf_helpers.h>
#include <bpf/bpf_core_read.h>
#include <bpf/bpf_tracing.h>

/* Extracted from the kernel */
#define MAX_PID_NS_LEVEL 32

#define EXE_PATH_MAX_LEN 1024

char _license[] SEC("license") = "GPL";

struct exported_task_info
{
	int32_t pid;
	int32_t tgid;
	int32_t parent_pid;
	int32_t parent_tgid;
	int32_t real_parent_pid;
	int32_t real_parent_tgid;
	char comm[TASK_COMM_LEN];
	uint32_t is_child_subreaper; /* Please note that a sub-reaper is
				       different from a reaper */
	uint32_t ns_level;
	int32_t vpid;
	int32_t vtgid;
	int32_t pgid;
	int32_t vpgid;
	int32_t sid;
	int32_t vsid;
	char exe_path[EXE_PATH_MAX_LEN];
	int64_t loginuid;
	int64_t euid;
} typedef exported_task_info;

/* If necessary we could expand this header */
struct header
{
	uint32_t version;
	uint32_t exported_task_info_size;
} typedef header;

/* Global variable to read data
 * We can do that because we iterate sequentially over
 * all tasks, not in parallel.
 * In this way we save space in the BPF stack
 */
exported_task_info data;
header h;

/* Keep the number of task struct visited */
uint64_t counter = 0;

#define UINT32_MAX (4294967295U)

/* used to check the endianness */
const uint16_t magic = 0xeB9F;
/* Initial version equals to 0 */
const uint32_t header_version = 0;

SEC("iter/task")
int dump_task(struct bpf_iter__task *ctx)
{
	struct seq_file *seq = ctx->meta->seq;
	struct task_struct *task = ctx->task;
	if(task == NULL)
	{
		return 0;
	}

	/* At the first iteration we send:
	 * - 1. the magic number to test endianness
	 * - 2. the header len
	 * - 3. the header
	 *   - the header version
	 *   - the size of exported_task_info
	 * - 4. List of task info
	 */
	counter++;
	if(counter == 1)
	{
		/* 1. Send the magic number */
		bpf_seq_write(seq, &magic, sizeof(magic));

		/* 2. Send the header size */
		u32 header_size = sizeof(h);
		bpf_seq_write(seq, &header_size, sizeof(header_size));

		/* 3. Send the header */
		h.version = header_version;
		h.exported_task_info_size = sizeof(data);
		bpf_seq_write(seq, &h, sizeof(h));
	}

	/* pid */
	data.pid = task->pid;

	/* tgid */
	data.tgid = task->tgid;

	/* parent_pid */
	data.parent_pid = task->parent->pid;

	/* parent_tgid */
	data.parent_tgid = task->parent->tgid;

	/* real_parent_pid */
	data.real_parent_pid = task->real_parent->pid;

	/* real_parent_tgid */
	data.real_parent_tgid = task->real_parent->tgid;

	/* comm */
	__builtin_memcpy(data.comm, task->comm, TASK_COMM_LEN);

	/* is_child_subreaper */
	struct signal_struct *signal = task->signal;
	data.is_child_subreaper = BPF_CORE_READ_BITFIELD(signal, is_child_subreaper);

	/* ns_level, current pid_ns */
	uint32_t ns_level = task->thread_pid->level;
	data.ns_level = ns_level;

	struct upid upid = {0};
	struct pid *pid_struct = NULL;

	/* vpid, `pid` referred to the current pid_ns */
	pid_struct = task->thread_pid;
	BPF_CORE_READ_INTO(&upid, pid_struct, numbers[ns_level & (MAX_PID_NS_LEVEL - 1)]);
	data.vpid = upid.nr;

	/* vtgid, `tgid` referred to the current pid_ns */
	pid_struct = task->signal->pids[PIDTYPE_TGID];
	BPF_CORE_READ_INTO(&upid, pid_struct, numbers[ns_level & (MAX_PID_NS_LEVEL - 1)]);
	data.vtgid = upid.nr;

	/* pgid */
	pid_struct = task->signal->pids[PIDTYPE_PGID];
	BPF_CORE_READ_INTO(&upid, pid_struct, numbers[0]);
	data.pgid = upid.nr;

	/* vpgid */
	BPF_CORE_READ_INTO(&upid, pid_struct, numbers[ns_level & (MAX_PID_NS_LEVEL - 1)]);
	data.vpgid = upid.nr;

	/* sid */
	pid_struct = task->signal->pids[PIDTYPE_SID];
	BPF_CORE_READ_INTO(&upid, pid_struct, numbers[0]);
	data.sid = upid.nr;

	/* vsid */
	BPF_CORE_READ_INTO(&upid, pid_struct, numbers[ns_level & (MAX_PID_NS_LEVEL - 1)]);
	data.vsid = upid.nr;

	/* exe_path */
	if(bpf_core_enum_value_exists(enum bpf_func_id, BPF_FUNC_d_path) &&
	   (bpf_core_enum_value(enum bpf_func_id, BPF_FUNC_d_path) == BPF_FUNC_d_path))
	{
		/* We need this extra check otherwise we risk a kernel panic.
		 * `exe_file` could be null.
		 */
		if(task->mm != NULL)
		{
			struct file *exe_file = task->mm->exe_file;
			if(exe_file != NULL)
			{
				/* Right now we don't care about the return value */
				bpf_d_path(&(exe_file->f_path), data.exe_path, EXE_PATH_MAX_LEN);
			}
		}
		else
		{
			data.exe_path[0] = '\0';
		}
	}
	else
	{
		data.exe_path[0] = '\0';
	}

	/* `loginuid` is an uint32_t but we use 64 bit in this way we can provide the user with a
	 * user-friendly info:
	 * - we return `-1` when loginuid is `UINT32_MAX`.
	 * - we return `-2` when for some reason we are not able to extract this info from the
	 * 	kernel. for example on COS the extraction path is different, or maybe the kernel is
	 * 	compiled without the `CONFIG_AUDIT` config.
	 */
	data.loginuid = -2;
	if(bpf_core_field_exists(task->loginuid))
	{
		data.loginuid = (s64)task->loginuid.val;
		if(data.loginuid == UINT32_MAX)
		{
			data.loginuid = -1;
		}
	}

	data.euid = task->cred->euid.val;
	if(data.euid == UINT32_MAX)
	{
		/* Like loginuid `-1` here is user friendly */
		data.euid = -1;
	}

	bpf_seq_write(seq, &data, sizeof(data));
	return 0;
}
