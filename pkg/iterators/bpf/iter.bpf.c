// go:build ignore
#include "vmlinux.h"
#include <bpf/bpf_helpers.h>
#include <bpf/bpf_core_read.h>
#include <bpf/bpf_tracing.h>
#include "consts.h"

char _license[] SEC("license") = "GPL";

extern void bpf_rcu_read_lock(void) __ksym;
extern void bpf_rcu_read_unlock(void) __ksym;

/////////////////////////
// Task iterator
/////////////////////////

typedef struct
{
	int32_t pid;
	int32_t tgid;
	int32_t pgid;
	int32_t sid;
	int32_t vpid;
	int32_t parent_pid;
	int32_t parent_tgid;
	int32_t real_parent_pid;
	int32_t real_parent_tgid;
	uint8_t is_child_subreaper; /* Please note that a sub-reaper is
				       different from a reaper */
	uint32_t ns_level;
	int64_t login_uid;
	int64_t e_uid;
	char comm[TASK_COMM_LEN];
	// todo!: we can use a new field exe_path_len and try to avoid the huge buffers. In this way
	// we can catch very long paths on demand.
	char exe_path[EXE_PATH_MAX_LEN];
	char cmdline[CMDLINE_MAX_LEN];
} __attribute__((packed, aligned(1))) exported_task_info;

/* Global variable to read task_info
 * We can do that because we iterate sequentially over
 * all tasks, not in parallel.
 * In this way we save space in the BPF stack
 */
exported_task_info task_info;

SEC("iter.s/task")
int dump_task(struct bpf_iter__task *ctx)
{
	struct seq_file *seq = ctx->meta->seq;
	struct task_struct *task = ctx->task;
	if(task == NULL)
	{
		// todo!: we can save some metrics and print them to the user in case we miss some
		// tasks/files.
		return 0;
	}

	/* pid */
	task_info.pid = task->pid;

	/* tgid */
	task_info.tgid = task->tgid;

	/* parent_pid */
	task_info.parent_pid = task->parent->pid;

	/* parent_tgid */
	task_info.parent_tgid = task->parent->tgid;

	/* real_parent_pid */
	task_info.real_parent_pid = task->real_parent->pid;

	/* real_parent_tgid */
	task_info.real_parent_tgid = task->real_parent->tgid;

	/* comm */
	__builtin_memcpy(task_info.comm, task->comm, TASK_COMM_LEN);

	/* is_child_subreaper */
	struct signal_struct *signal = task->signal;
	task_info.is_child_subreaper = BPF_CORE_READ_BITFIELD(signal, is_child_subreaper);

	/* ns_level, current pid_ns */
	uint32_t ns_level = task->thread_pid->level;
	task_info.ns_level = ns_level;

	struct upid upid = {0};
	struct pid *pid_struct = NULL;

	/* vpid, `pid` referred to the current pid_ns */
	pid_struct = task->thread_pid;
	BPF_CORE_READ_INTO(&upid, pid_struct, numbers[ns_level & (MAX_PID_NS_LEVEL - 1)]);
	task_info.vpid = upid.nr;

	/* pgid */
	pid_struct = task->signal->pids[PIDTYPE_PGID];
	BPF_CORE_READ_INTO(&upid, pid_struct, numbers[0]);
	task_info.pgid = upid.nr;

	/* sid */
	pid_struct = task->signal->pids[PIDTYPE_SID];
	BPF_CORE_READ_INTO(&upid, pid_struct, numbers[0]);
	task_info.sid = upid.nr;

	/* exe_path */
	if(bpf_core_enum_value_exists(enum bpf_func_id, BPF_FUNC_d_path) &&
	   (bpf_core_enum_value(enum bpf_func_id, BPF_FUNC_d_path) == BPF_FUNC_d_path))
	{
		// We need to enter the RCU critical section otherwise when we dereference the
		// `task` pointer we will obtain an `untrusted` pointer instead of a
		// `PTR_TO_BTF_ID`. See this
		// https://lore.kernel.org/bpf/CAGQdkDtgqwJkHyx+txp6hQD83qUaRRWuo7nYMVGZq79xw+kgTA@mail.gmail.com/T/#me219e6226349c757071894bbd7673df08e02745c
		bpf_rcu_read_lock();
		if(task->mm != NULL)
		{
			struct file *exe_file = task->mm->exe_file;
			/* We need this extra check otherwise we risk a kernel panic.
			 * `exe_file` could be null.
			 */
			if(exe_file != NULL)
			{
				// According to the manual if the path is too long the helper
				// doesn't populate the path so we return a too-long.
				if(bpf_d_path(&(exe_file->f_path), task_info.exe_path,
					      EXE_PATH_MAX_LEN) < 0)
				{
					// The path was too long
					__builtin_memcpy(task_info.exe_path, TOO_LONG,
							 TOO_LONG_SIZE);
				}
			}
		}
		else
		{
			__builtin_memcpy(task_info.exe_path, NOT_AVAILABLE, NOT_AVAILABLE_SIZE);
		}
		bpf_rcu_read_unlock();
	}
	else
	{
		__builtin_memcpy(task_info.exe_path, NOT_AVAILABLE, NOT_AVAILABLE_SIZE);
	}

	/* `loginuid` is an uint32_t but we use 64 bit in this way we can provide the user with a
	 * user-friendly info:
	 * - we return `-1` when loginuid is `UINT32_MAX`.
	 * - we return `-2` when for some reason we are not able to extract this info from the
	 * 	kernel. for example on COS the extraction path is different, or maybe the kernel is
	 * 	compiled without the `CONFIG_AUDIT` config.
	 */
	task_info.login_uid = -2;
	if(bpf_core_field_exists(task->loginuid))
	{
		task_info.login_uid = (s64)task->loginuid.val;
		if(task_info.login_uid == UINT32_MAX)
		{
			task_info.login_uid = -1;
		}
	}

	task_info.e_uid = task->cred->euid.val;
	if(task_info.e_uid == UINT32_MAX)
	{
		/* Like loginuid `-1` here is user friendly */
		task_info.e_uid = -1;
	}

	// todo!: we need to manage the `set_proctitle` issue:
	// https://github.com/falcosecurity/libs/issues/988
	unsigned long arg_start_pointer = task->mm->arg_start;
	unsigned long arg_end_pointer = task->mm->arg_end;
	const uint16_t cmdline_len = arg_end_pointer - arg_start_pointer >= CMDLINE_MAX_LEN
					     ? CMDLINE_MAX_LEN - 1
					     : arg_end_pointer - arg_start_pointer;

	int ret = bpf_copy_from_user_task(&task_info.cmdline[0], SAFE_1024_ACCESS(cmdline_len),
					  (void *)arg_start_pointer, task, 0);
	if(ret != 0)
	{
		__builtin_memcpy(task_info.cmdline, NOT_AVAILABLE, NOT_AVAILABLE_SIZE);
	}

	bpf_seq_write(seq, &task_info, sizeof(task_info));
	return 0;
}

/////////////////////////
// File iterator
/////////////////////////

typedef struct
{
	int32_t task_pid; // todo!: we can do better than this, we use this to correlate the fd to
			  // the task in userspace
	uint32_t fd;
	uint64_t ino;
	uint32_t mount_id;
	uint32_t flags;
	char file_path[FILE_PATH_MAX_LEN];
} __attribute__((packed, aligned(1))) exported_file_info;

exported_file_info file_info;

// We consider the `inode` as a pointer to `vfs_inode` in the struct `struct ovl_inode`
//
// struct ovl_inode {
//	union {
//		struct ovl_dir_cache *cache;	/* directory */
//		const char *lowerdata_redirect;	/* regular file */
//	};
//	const char *redirect;
//	u64 version;
//	unsigned long flags;
//	struct inode vfs_inode;
//	struct dentry *__upperdentry;
//  struct ovl_entry *oe;
//
bool __attribute__((always_inline)) is_upper_layer(struct dentry *dentry)
{
	struct inode *vfs_inode = BPF_CORE_READ(dentry, d_inode);
	unsigned long inode_size = bpf_core_type_size(struct inode);
	if(!inode_size)
	{
		return false;
	}

	struct dentry *upper_dentry = NULL;
	bpf_probe_read_kernel(&upper_dentry, sizeof(upper_dentry), (char *)vfs_inode + inode_size);
	if(upper_dentry == NULL)
	{
		return false;
	}

	// If we are in the upper layer the upper inode should be different from 0
	return BPF_CORE_READ(upper_dentry, d_inode, i_ino) != 0;
}

SEC("iter/task_file")
int dump_task_file(struct bpf_iter__task_file *ctx)
{
	struct seq_file *seq = ctx->meta->seq;
	struct task_struct *task = ctx->task;
	struct file *file = ctx->file;
	if(task == NULL || file == NULL || file->f_inode == NULL)
	{
		return 0;
	}

	file_info.task_pid = task->pid;
	file_info.fd = ctx->fd;
	file_info.ino = file->f_inode->i_ino;

	if(bpf_core_enum_value_exists(enum bpf_func_id, BPF_FUNC_d_path) &&
	   (bpf_core_enum_value(enum bpf_func_id, BPF_FUNC_d_path) == BPF_FUNC_d_path))
	{
		// According to the manual if the path is too long the helper
		// doesn't populate the path so we return a too-long.
		if(bpf_d_path(&(file->f_path), file_info.file_path, FILE_PATH_MAX_LEN) < 0)
		{
			__builtin_memcpy(file_info.file_path, TOO_LONG, TOO_LONG_SIZE);
		}
	}
	else
	{
		__builtin_memcpy(file_info.file_path, NOT_AVAILABLE, NOT_AVAILABLE_SIZE);
	}

	struct dentry *dentry = BPF_CORE_READ(file, f_path.dentry);
	struct super_block *sb = BPF_CORE_READ(dentry, d_sb);
	int flags = BPF_CORE_READ(sb, s_flags);

	file_info.flags = 0;
	// If there is no flag the file is mountable
	if((flags & MS_NOUSER) == 0)
	{
		file_info.flags |= MOUNTABLE;
	}

	if(BPF_CORE_READ(sb, s_magic) == PPM_OVERLAYFS_SUPER_MAGIC)
	{
		if(is_upper_layer(dentry))
		{
			file_info.flags |= UPPER_LAYER;
		}
		else
		{
			file_info.flags |= LOWER_LAYER;
		}
	}

	struct vfsmount *vfsmnt = file->f_path.mnt;
	struct mount *mnt_p = container_of(vfsmnt, struct mount, mnt);
	file_info.mount_id = BPF_CORE_READ(mnt_p, mnt_id);

	bpf_seq_write(seq, &file_info, sizeof(file_info));
	return 0;
}
