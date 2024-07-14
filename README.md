# 🌴 bpftree

## 🗞️ Why bpftree?

During my working routine, I often use tools like `pstree`, `pidof` combined with some `/proc` queries. Consider a concrete case in which you have just spawned a `tail` process inside a container and you want to know its full lineage. With tools like `pstree` you should probably `grep` the tail process and then try to trace the lineage by hand, nothing impossible but of course quite noisy... Moreover imagine you want to know some details about one of the ancestors, in the best case you have to scan `/proc` for this info, but in the worst one `/proc` can't help you either because the information you are searching for is not available... So how `bpftree` come into play? Using the power of eBPF iterators we can scan all tasks running in the kernel and aggregate them generating graphs like trees and process lineages. Moreover, since we are using eBPF we can scrape almost all data we want from the kernel without strict constraints due to the `/proc` interface!
Let's see `bpftree` in action.

We said we want to obtain the `tail` lineage, well, let's ask `bpftree` to do it for us:

```bash
sudo bpftree lineage comm tail
# or short form
sudo bpftree l c tail
```

Here we are asking `bpftree` to print the lineage for all processes with `comm==tail`. No need for `pidof tail` or some "fancy greps" you can just use the process name.

In this example, we have only one `tail` process in our system:

```plain
📜 Task Lineage for 'comm=tail'
⬇️ [tail] tid: 24996, pid: 24996, rptid: 24776, rppid: 24776
⬇️ [zsh] tid: 24776, pid: 24776, rptid: 24755, rppid: 24755
⬇️ [gnome-terminal-] tid: 24755, pid: 24755, rptid: 1297, rppid: 1297
⬇️ [systemd]💀 tid: 1297, pid: 1297, rptid: 1, rppid: 1
⬇️ [systemd] tid: 1, pid: 1, rptid: 0, rppid: 0
```

> __Please note__: The symbol 💀 means that the process is a [`child_subreaper`](https://elixir.bootlin.com/linux/latest/source/include/linux/sched/signal.h#L132).

Let's consider we want to know some extra info, something like the session id (`sid`). Let's configure `bpftree` to print tasks according to a specific format:

```bash
sudo bpftree lineage comm tail --format 'tid,pid,sid'
# or short form
sudo bpftree l c tail -f 't,p,s'
```

Here we are asking the same thing as before but we are telling `bpftree` to print tasks according to a specific format. Please note that when you specify a custom format all the displayed fields are abbreviated (`tid` becomes `t` and so on) this is done to support long lists of fields.

```plain
📜 Task Lineage for comm: tail
⬇️ [tail] t: 24996, p: 24996, s: 24776
⬇️ [zsh] t: 24776, p: 24776, s: 24776
⬇️ [gnome-terminal-] t: 24755, p: 24755, s: 24755
⬇️ [systemd]💀 t: 1297, p: 1297, s: 1297
⬇️ [systemd] t: 1, p: 1, s: 1
```

You can find the full list of fields with all the associated abbreviations [here](#✔️-supported-fields) or you can directly use the `fields` command:

```bash
sudo bpftree fields
# or short form
sudo bpftree f
```

> __Please note__: if you run bpftree inside a container you will see only the processes inside the current pid namespace.

## ⛓️ Requirements

The "known supported architectures" are:

- x86_64
- aarch64
- s390x

"known supported architectures" means that we tried `bpftree` on these architectures but others could be supported as well.

Right now the unique BPF requirement is the presence of BPF iterators in the running system. They were introduced in kernel version [5.8](https://github.com/torvalds/linux/commit/eaaacd23910f2d7c4b22d43f591002cc217d294b) so all newer kernel versions should be supported

## ✔️ Supported fields

| Fields      | Description                                                         |
| ----------- | ------------------------------------------------------------------- |
| tid,t       | thread id (tid) of the current task (init namespace)                |
| vtid,vt     | thread id (tid) of the current task (task namespace)                |
| pid,p       | process id (pid) of the current task (init namespace)               |
| vpid,vp     | process id (pid) of the current task (task namespace)               |
| pgid,pg     | process group id (pgid) of the current task (init namespace)        |
| vpgid,vpg   | process group id (pgid) of the current task (task namespace)        |
| sid,s       | session id (sid) of the current task (init namespace)               |
| vsid,vs     | session id (sid) of the current task (task namespace)               |
| ptid,pt     | parent thread id (ptid) of the current task (init namespace)        |
| ppid,pp     | parent process id (ppid) of the current task (init namespace)       |
| rptid,rpt   | real parent thread id (rptid) of the current task (init namespace)  |
| rppid,rpp   | real parent process id (rppid) of the current task (init namespace) |
| comm,c      | human readeable process name ('task->comm')                         |
| reaper,r    | true if the current process is a child_sub_reaper                   |
| ns_level,ns | pid namespace level of the actual thread                            |
| exepath,e   | full executable path of the current task                            |

In the `Fields` column, we have the full field name as a first string and its short name as a second one (e.g. "`tid,t`", `tid` is the full name while `t` is the short one).

To specify the format you can use either the full name or the short one:

```bash
## Full
sudo ./bpftree info tid 1 --format "tid,pid" 
## Short
sudo ./bpftree i t 1 -f "t,p"
## They will produce the same result
🗞️ [systemd] t: 1, p: 1
```

The output always uses the short format to improve readability with long field lists

## 🎮 Available commands

### ℹ️ info

This command shows all tasks that match a certain condition. The condition is expressed through a field, the field name is provided as the first argument while the value is provided as the second one.

```bash
sudo bpftree info pid 8056
# or in short form
sudo bpftree i p 8056
```

```plain
ℹ️ Task Info for pid: 8056
🗞️ [plugin_host-3.3] tid: 8056, pid: 8056, rptid: 7972, rppid: 7972
🗞️ {thread_queue} tid: 8061, pid: 8056, rptid: 7972, rppid: 7972
🗞️ {process_status} tid: 8064, pid: 8056, rptid: 7972, rppid: 7972
```

Here the condition is: *"print all tasks with pid==8056"*. So print all tasks that belong to the thread group with id `8056`.

To write a condition you can use all fields shown by the `fields` command. These same fields can be also used to customize the command output but by default, all tasks are rendered in the default "info format":

```plain
[comm] tid: ..., pid: ..., rptid: ..., rppid: ...
```

> __Please note__: When the rendered task is a leader thread its comm is displayed like this `[comm]` while for secondary threads we use this notation `{comm}`.

If you want to obtain also the `vtid` for all the selected tasks, you have just to customize the format. All fields must be separated by ',' (e.g. 'tid,pid,reaper') and DON'T have white spaces between them.

```bash
sudo bpftree info pid 3 --format 'tid,pid,rptid,rppid,vtid'
# or in short form
sudo bpftree i p 3 --format 't,p,rpt,rpp,vt'
```

```plain
ℹ️ Task Info for pid: 8056
🗞️ [plugin_host-3.3] t: 8056, p: 8056, rpt: 7972, rpp: 7972, vt: 8056
🗞️ {thread_queue} t: 8061, p: 8056, rpt: 7972, rpp: 7972, vt: 8061
🗞️ {process_status} t: 8064, p: 8056, rpt: 7972, rpp: 7972, vt: 8064
```

### 📜 lineage

This command shows the lineage for all tasks that match a certain condition. The command features are the same as the "info" ones.

```bash
sudo bpftree lineage pid 8056
# or in short form
sudo bpftree l p 8056
```

```plain
📜 Task Lineage for pid: 8056
⬇️ [plugin_host-3.3] tid: 8056, pid: 8056, rptid: 7972, rppid: 7972
⬇️ [sublime_text] tid: 7972, pid: 7972, rptid: 1486, rppid: 1486
⬇️ [systemd]💀 tid: 1486, pid: 1486, rptid: 1, rppid: 1
⬇️ [systemd] tid: 1, pid: 1, rptid: 0, rppid: 0

⬇️ {thread_queue} tid: 8061, pid: 8056, rptid: 7972, rppid: 7972
⬇️ [sublime_text] tid: 7972, pid: 7972, rptid: 1486, rppid: 1486
⬇️ [systemd]💀 tid: 1486, pid: 1486, rptid: 1, rppid: 1
⬇️ [systemd] tid: 1, pid: 1, rptid: 0, rppid: 0

⬇️ {process_status} tid: 8064, pid: 8056, rptid: 7972, rppid: 7972
⬇️ [sublime_text] tid: 7972, pid: 7972, rptid: 1486, rppid: 1486
⬇️ [systemd]💀 tid: 1486, pid: 1486, rptid: 1, rppid: 1
⬇️ [systemd] tid: 1, pid: 1, rptid: 0, rppid: 0
```

`bpftree` shows the lineage for all the tasks with `pid==8056`.

Also here the default "lineage format" is:

```plain
[comm] tid: ..., pid: ..., rptid: ..., rppid: ...
```

### 🌴 tree

This command shows the tree for all tasks that match a certain condition. The command features are the same as the "info" ones.

```bash
sudo bpftree tree pid 8056
# or in short form
sudo bpftree t p 8056
```

```plain
🌴 Task Tree for pid: 8056
[plugin_host-3.3] tid: 8056, pid: 8056

{thread_queue} tid: 8061, pid: 8056

{process_status} tid: 8064, pid: 8056
```

`bpftree` shows the tree for all the tasks with `pid==8056`. Here is quite useless since all the involved tasks are tree leaves.

The default "tree format" is:

```plain
[comm] tid: ..., pid: ...
```

but you can always customize it with the `--format` flag.

### 🗺️ fields

This command shows all fields that can be used both to write conditions and to customize the output.

```bash
sudo bpftree fields
# or short form
sudo bpftree f
```

### 🔻 dump

This command allows you to save read system tasks into a file. Let's imagine you are using an enterprise tool and you find a bug with a specific system configuration. You could dump a capture file and send it to the support team to analyze it. Thanks to the `--capture` flag `bpftree` can replay previously dumped files coming also from systems with different architectures.

```bash
# Dump the file
sudo bpftree dump output-file.tree
# Read it with the capture flag
sudo bpftree info tid 1 --capture output-file.tree
```

On the capture file, you can use all commands listed above like in a live run.

The capture files follow the versioning of the bpftree tool. This means that a particular version of bpftree can read only capture files with the same Major and Minor and with a lower or equal Patch. If needed we could create a wrapper tool that calls the right bpftree version for the provided capture file.

## For developers

### 🏗️ Build from source

As with many go project you have just to type a bunch of commands in the root folder of the project:

```bash
go generate ./...
go build .
```

### 🏎️ Run

In the root folder of the project:

```bash
sudo ./bpftree --help
```

As an example try to log the full process tree of your system 🕹️

```bash
sudo ./bpftree tree tid 1
# or short form
sudo ./bpftree t t 1
```

### 🧪 Run tests

In the root folder of the project run:

```bash
# Unit tests
sudo go test ./pkg/...

# All tests
sudo -E env "PATH=$PATH" go test ./... -count=1
```

### Run linters

To run golangci-lint

```bash
golangci-lint run  --timeout=900s
```

> __Note__: be sure that the used golangci-lint version is compatible with your go version. Tested with `golangci-lint 1.59.1 built with go1.22.3` and `go1.22.1`

### ➕ Add a new field

- BPF side
   1. add the new field to `exported_task_info` or `exported_file_info` struct
   2. instrument the code to collect the new field
- Userspace side
   1. add the new field to `taskInfo` or `fileInfo` struct
   3. add a getter method for the new field in `task.go`
   4. add an enum `allowedField` for the new field
   5. add a new entry for the field into `allowedFieldsSlice`
- Tests
   1. add an entry in the `TestGetFieldMatrix` for the new getter method
   2. update `testdata/script/cmd_fields/cmd_fields.txtar` file with the new output of `sudo ./bpftree f`
   3. update `testdata/script/fields/systemd_fields.txtar` tests with the new field
