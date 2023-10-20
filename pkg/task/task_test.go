package task

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"testing"
)

func populateMockTasksMap() {
	/* - [systemd] tid: 1, pid: 1, rptid: 0
	 *   - [p1t1] tid: 2, pid: 2, rptid: 1
	 *   - [p1t2] tid: 3, pid: 2, rptid: 1
	 *    - [p2t1] tid: 7, pid: 7, rptid: 1
	 *    - [p2t2] tid: 8, pid: 7, rptid: 1
	 *     - [p3t1] tid: 10, pid: 10, rptid: 8
	 */
	tasksList = []*task{
		{
			Info: TaskInfo{
				Tid:              1,
				Pid:              1,
				IsChildSubreaper: 0,
				RealParentTid:    0,
				RealParentPid:    0,
				ParentTid:        0,
				ParentPid:        0,
				Comm:             [commLen]byte{'s', 'y', 's', 't', 'e', 'm', 'd'},
			},
		},
		{
			Info: TaskInfo{
				Tid:              2,
				Pid:              2,
				IsChildSubreaper: 0,
				RealParentTid:    1,
				RealParentPid:    1,
				ParentTid:        1,
				ParentPid:        1,
				Comm:             [commLen]byte{'p', '1', 't', '1'},
			},
		},
		{
			Info: TaskInfo{
				Tid:              3,
				Pid:              2,
				IsChildSubreaper: 0,
				RealParentTid:    1,
				RealParentPid:    1,
				ParentTid:        1,
				ParentPid:        1,
				Comm:             [commLen]byte{'p', '1', 't', '2'},
			},
		},
		{
			Info: TaskInfo{
				Tid:              7,
				Pid:              7,
				IsChildSubreaper: 1,
				RealParentTid:    1,
				RealParentPid:    1,
				ParentTid:        1,
				ParentPid:        1,
				Comm:             [commLen]byte{'p', '2', 't', '1'},
			},
		},
		{
			Info: TaskInfo{
				Tid:              8,
				Pid:              7,
				IsChildSubreaper: 1,
				RealParentTid:    1,
				RealParentPid:    1,
				ParentTid:        1,
				ParentPid:        1,
				Comm:             [commLen]byte{'p', '2', 't', '2'},
			},
		},
		{
			Info: TaskInfo{
				Tid:              10,
				Pid:              10,
				IsChildSubreaper: 0,
				RealParentTid:    8,
				RealParentPid:    7,
				ParentTid:        8,
				ParentPid:        7,
				Comm:             [commLen]byte{'p', '3', 't', '1'},
			},
		},
	}

	/* Order the tasks List */
	sort.SliceStable(tasksList, func(i, j int) bool {
		return tasksList[i].Info.Tid < tasksList[j].Info.Tid
	})
}

func init() {
	populateMockTasksMap()
}

func TestIsMainThread(t *testing.T) {
	var tests = []struct {
		name         string
		task         task
		isMainThread bool
	}{
		{"MainThreadTrue", task{Info: TaskInfo{Tid: 1, Pid: 1}}, true},
		{"MainThreadFalse", task{Info: TaskInfo{Tid: 1, Pid: 2}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.task.isMainThread() != tt.isMainThread {
				t.Errorf("Expected:\n%t\n\nGot:\n%t\n", tt.isMainThread, tt.task.isMainThread())
			}
		})
	}
}

func TestGetFieldMatrix(t *testing.T) {
	var tests = []struct {
		name          string
		task          task
		expectedField string
		fieldGetter   func(t task) string
	}{
		{"GetTid",
			task{Info: TaskInfo{
				Tid: 1,
			}},
			"1",
			func(t task) string { return fmt.Sprint(t.getTid()) }},
		{"GetPid",
			task{Info: TaskInfo{
				Pid: 2,
			}},
			"2",
			func(t task) string { return fmt.Sprint(t.getPid()) }},
		{"GetParentTid",
			task{Info: TaskInfo{
				ParentTid: 3,
			}},
			"3",
			func(t task) string { return fmt.Sprint(t.getParentTid()) }},
		{"GetParentPid",
			task{Info: TaskInfo{
				ParentPid: 4,
			}},
			"4",
			func(t task) string { return fmt.Sprint(t.getParentPid()) }},
		{"GetRealParentTid",
			task{Info: TaskInfo{
				RealParentTid: 5,
			}},
			"5",
			func(t task) string { return fmt.Sprint(t.getRealParentTid()) }},
		{"GetRealParentPid",
			task{Info: TaskInfo{
				RealParentPid: 6,
			}},
			"6",
			func(t task) string { return fmt.Sprint(t.getRealParentPid()) }},
		{"GetComm",
			task{Info: TaskInfo{
				Comm: [commLen]byte{'0'},
			}},
			"0",
			func(t task) string { return t.getComm() }},
		{"GetCommWithoutTerminator",
			task{Info: TaskInfo{
				Comm: [commLen]byte{
					'p', 'p', 'p', 'p', 'p', 'p', 'p', 'p', 'p', 'p', 'p', 'p', 'p', 'p', 'p', 'p',
				},
			}},
			"pppppppppppppppp",
			func(t task) string { return t.getComm() }},
		{"GetCommBytes",
			task{Info: TaskInfo{
				Comm: [commLen]byte{99, 111, 110, 116, 97, 105, 110, 101, 114, 100},
			}},
			"containerd",
			func(t task) string { return t.getComm() }},
		{"GetCommContainerd",
			task{Info: TaskInfo{
				Comm: [commLen]byte{'c', 'o', 'n', 't', 'a', 'i', 'n', 'e', 'r', 'd'},
			}},
			"containerd",
			func(t task) string { return t.getComm() }},
		{"GetReaperTrue",
			task{Info: TaskInfo{
				IsChildSubreaper: 1,
			}},
			"true",
			func(t task) string { return fmt.Sprint(t.isChildSubReaper()) }},
		{"GetReaperFalse",
			task{Info: TaskInfo{
				IsChildSubreaper: 0,
			}},
			"false",
			func(t task) string { return fmt.Sprint(t.isChildSubReaper()) }},
		{"GetNsLevel",
			task{Info: TaskInfo{
				NsLevel: 7,
			}},
			"7",
			func(t task) string { return fmt.Sprint(t.getNsLevel()) }},
		{"GetVTid",
			task{Info: TaskInfo{
				VTid: 8,
			}},
			"8",
			func(t task) string { return fmt.Sprint(t.getVTid()) }},
		{"GetVPid",
			task{Info: TaskInfo{
				VPid: 9,
			}},
			"9",
			func(t task) string { return fmt.Sprint(t.getVPid()) }},
		{"GetPgid",
			task{Info: TaskInfo{
				Pgid: 10,
			}},
			"10",
			func(t task) string { return fmt.Sprint(t.getPgid()) }},
		{"GetVPgid",
			task{Info: TaskInfo{
				VPgid: 11,
			}},
			"11",
			func(t task) string { return fmt.Sprint(t.getVPgid()) }},
		{"GetSid",
			task{Info: TaskInfo{
				Sid: 12,
			}},
			"12",
			func(t task) string { return fmt.Sprint(t.getSid()) }},
		{"GetVSid",
			task{Info: TaskInfo{
				VSid: 13,
			}},
			"13",
			func(t task) string { return fmt.Sprint(t.getVSid()) }},
		{"GetExePath",
			task{Info: TaskInfo{
				ExePath: [exePathLen]byte{'/', 'u', 's', 'r', '/'},
			}},
			"/usr/",
			func(t task) string { return t.getExePath() }},
		{"GetExePathBytes",
			task{Info: TaskInfo{
				ExePath: [exePathLen]byte{99, 111, 110, 116, 97, 105, 110, 101, 114, 100},
			}},
			"containerd",
			func(t task) string { return t.getExePath() }},
		{"GetExePathContainerd",
			task{Info: TaskInfo{
				ExePath: [exePathLen]byte{'c', 'o', 'n', 't', 'a', 'i', 'n', 'e', 'r', 'd'},
			}},
			"containerd",
			func(t task) string { return t.getExePath() }},
		{"GetLoginUid",
			task{Info: TaskInfo{
				LoginUID: -2,
			}},
			"-2",
			func(t task) string { return fmt.Sprint(t.getLoginUID()) }},
		{"GetEUid",
			task{Info: TaskInfo{
				EUID: -1,
			}},
			"-1",
			func(t task) string { return fmt.Sprint(t.getEUID()) }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.fieldGetter(tt.task) != tt.expectedField {
				t.Errorf("Expected:\n%s\n\nGot:\n%s\n", tt.expectedField, tt.fieldGetter(tt.task))
			}
		})
	}
}

func TestGetTaskFromTidCrash(t *testing.T) {
	if os.Getenv("GETTASKFROMTID") == "1" {
		getTaskFromTid(1000000)
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestGetTaskFromTidCrash")
	cmd.Env = append(os.Environ(), "GETTASKFROMTID=1")
	err := cmd.Run()
	var exitError *exec.ExitError
	if errors.As(err, &exitError) {
		return
	}
	t.Fatalf("process ran with err %v, want exit status 1", err)
}

func TestGetTaskFromTidPresence(t *testing.T) {
	/* in the tasksList populate by `init` we have a task with tid 2. */
	task := getTaskFromTid(2)
	if task.getTid() != 2 {
		t.Errorf("Expected task with tid:\n%d\n\nGot task with tid:\n%d\n", 2, task.getTid())
	}
}

func TestGetSelectedTasksFromField(t *testing.T) {
	if os.Getenv("GETSELECTEDTASKS") == "1" {
		getSelectedTasksFromField("not-existent", "33huh")
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestGetSelectedTasksFromField")
	cmd.Env = append(os.Environ(), "GETSELECTEDTASKS=1")
	err := cmd.Run()
	var exitError *exec.ExitError
	if errors.As(err, &exitError) {
		return
	}
	t.Fatalf("process ran with err %v, want exit status 1", err)
}

func TestGetTaskFromPrefixUnique(t *testing.T) {
	taskSlice := getSelectedTasksFromField("comm", "p2t2")
	expectedSliceLen := 1
	if len(taskSlice) != expectedSliceLen {
		t.Errorf("Expected :\n%v\n\nGot:\n%v\n", expectedSliceLen, taskSlice)
	}
}

func TestComputeChildren(t *testing.T) {
	computeChildren()
	if len(getTaskFromTid(1).getChildren()) != 4 {
		t.Errorf("Expected :\n%d\n\nGot:\n%d\n", 4, len(getTaskFromTid(1).getChildren()))
	}

	if len(getTaskFromTid(8).getChildren()) != 1 {
		t.Errorf("Expected :\n%d\n\nGot:\n%d\n", 1, len(getTaskFromTid(8).getChildren()))
	}

	if len(getTaskFromTid(2).getChildren()) != 0 {
		t.Errorf("Expected :\n%d\n\nGot:\n%d\n", 0, len(getTaskFromTid(2).getChildren()))
	}
}
