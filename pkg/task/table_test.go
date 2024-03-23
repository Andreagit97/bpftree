package task

import (
	"sort"
	"testing"
)

func populateMockTasksMap() {

	/* Considering these 2 lineages:
	 * 📜 Task lineage for 'tid=273408'
	 * ⬇️ [containerd] t: 273408, p: 273408, pg: 273408, s: 273408, vt: 105, vp: 105, vpg: 105, vs: 105, rpt: 273246
	 * ⬇️ [systemd] t: 273246, p: 273246, pg: 273246, s: 273246, vt: 1, vp: 1, vpg: 1, vs: 1, rpt: 273226
	 * ⬇️ [containerd-shim]💀 t: 273226, p: 273226, pg: 273226, s: 1130, vt: 273226, vp: 273226, vpg: 273226, vs: 1130, rpt: 1
	 * ⬇️ [systemd] t: 1, p: 1, pg: 1, s: 1, vt: 1, vp: 1, vpg: 1, vs: 1, rpt: 0
	 *
	 * 📜 Task lineage for 'tid=1130'
	 * ⬇️ [containerd] t: 1130, p: 1130, pg: 1130, s: 1130, vt: 1130, vp: 1130, vpg: 1130, vs: 1130, rpt: 1
	 * ⬇️ [systemd] t: 1, p: 1, pg: 1, s: 1, vt: 1, vp: 1, vpg: 1, vs: 1, rpt: 0
	 */
	tasksList = []*task{
		{
			Info: &taskInfo{
				Tid:              273408,
				Pid:              273408,
				Pgid:             273408,
				Sid:              273408,
				VTid:             105,
				IsChildSubreaper: 0,
				RealParentTid:    273246,
				Comm:             [commLen]byte{'c', 'o', 'n', 't', 'a', 'i', 'n', 'e', 'r', 'd'},
			},
		},
		{
			Info: &taskInfo{
				Tid:              273246,
				Pid:              273246,
				Pgid:             273246,
				Sid:              273246,
				VTid:             1,
				IsChildSubreaper: 0,
				RealParentTid:    273226,
				Comm:             [commLen]byte{'s', 'y', 's', 't', 'e', 'm', 'd'},
			},
		},
		{
			Info: &taskInfo{
				Tid:              273226,
				Pid:              273226,
				Pgid:             273226,
				Sid:              1130,
				VTid:             273226,
				IsChildSubreaper: 1,
				RealParentTid:    1,
				Comm:             [commLen]byte{'c', 'o', 'n', 't', 'a', 'i', 'n', 'e', 'r', 'd', '-', 's', 'h', 'i', 'm'},
			},
		},
		{
			Info: &taskInfo{
				Tid:              1,
				Pid:              1,
				Pgid:             1,
				Sid:              1,
				VTid:             1,
				IsChildSubreaper: 0,
				RealParentTid:    0,
				Comm:             [commLen]byte{'s', 'y', 's', 't', 'e', 'm', 'd'},
			},
		},
		{
			Info: &taskInfo{
				Tid:              1130,
				Pid:              1130,
				Pgid:             1130,
				Sid:              1130,
				VTid:             1130,
				IsChildSubreaper: 0,
				RealParentTid:    1,
				Comm:             [commLen]byte{'c', 'o', 'n', 't', 'a', 'i', 'n', 'e', 'r', 'd'},
			},
		},
	}

	// We sort the slice because we need to have the tasks in order to compute the children.
	sort.SliceStable(tasksList, func(i, j int) bool {
		return tasksList[i].Info.Tid < tasksList[j].Info.Tid
	})
}

func TestGetTaskFromTid(t *testing.T) {
	populateMockTasksMap()

	if getTaskFromTid(999999) != nil {
		t.Errorf("Expected nil, got %v", getTaskFromTid(999999))
	}

	initThread := getTaskFromTid(1)
	if initThread == nil {
		t.Errorf("Expected valid thread, got %v", initThread)
	}

	if initThread.GetTid() != 1 {
		t.Errorf("Expected tid 1, got %d", initThread.GetTid())
	}
}

func TestGetSelectedTasksFromField(t *testing.T) {
	var tests = []struct {
		name        string
		field       allowedField
		value       string
		expectedLen int
		setup       func()
	}{
		{
			name:        "[Comm]NotExisting",
			field:       commField,
			value:       "NotExisting",
			expectedLen: 0,
			setup: func() {
				populateMockTasksMap()
			},
		},
		{
			name:        "[Comm]containerd-shim",
			field:       commField,
			value:       "containerd-shim",
			expectedLen: 1,
			setup: func() {
				populateMockTasksMap()
			},
		},
		{
			name:        "[Comm]containerd",
			field:       commField,
			value:       "containerd",
			expectedLen: 2,
			setup: func() {
				populateMockTasksMap()
			},
		},
		{
			name:        "[rptid]1",
			field:       realParentTidField,
			value:       "1",
			expectedLen: 2,
			setup: func() {
				populateMockTasksMap()
			},
		},
		{
			name:        "[tid]273246",
			field:       tidField,
			value:       "273246",
			expectedLen: 1,
			setup: func() {
				populateMockTasksMap()
			},
		},
		{
			name:        "[vsid]1",
			field:       vSidField,
			value:       "1",
			expectedLen: 2,
			setup: func() {
				populateMockTasksMap()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			taskSliceLen := len(getSelectedTasksFromField(tt.field, tt.value))
			if taskSliceLen != tt.expectedLen {
				t.Errorf("Expected :\n%v\n\nGot:\n%v\n", tt.expectedLen, taskSliceLen)
			}
		})
	}
}

func TestComputeChildren(t *testing.T) {
	var tests = []struct {
		name             string
		tid              int
		expectedChildren int
		setup            func()
	}{
		{
			name:             "Systemd",
			tid:              1,
			expectedChildren: 2,
			setup: func() {
				populateMockTasksMap()
				computeChildren()
			},
		},
		{
			name:             "Containerd",
			tid:              273408,
			expectedChildren: 0,
			setup: func() {
				populateMockTasksMap()
				computeChildren()
			},
		},
		{
			name:             "Containerd-shim",
			tid:              273226,
			expectedChildren: 1,
			setup: func() {
				populateMockTasksMap()
				computeChildren()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			childrenLen := len(getTaskFromTid(tt.tid).TaskChildren)
			if childrenLen != tt.expectedChildren {
				t.Errorf("Expected :\n%d\n\nGot:\n%d\n", tt.expectedChildren, childrenLen)
			}

		})
	}
}

func TestVirtualIDs(t *testing.T) {
	var tests = []struct {
		name        string
		tid         int
		expectedvid int
		call        func(t *task) int
		setup       func()
	}{
		{
			name:        "GetVirtualSid",
			tid:         273408,
			expectedvid: 105,
			call:        func(t *task) int { return t.getVSid() },
			setup: func() {
				populateMockTasksMap()
			},
		},

		{
			name:        "GetVirtualPgid",
			tid:         273246,
			expectedvid: 1,
			call:        func(t *task) int { return t.getVPgid() },
			setup: func() {
				populateMockTasksMap()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			task := getTaskFromTid(tt.tid)
			if task == nil {
				t.Errorf("Expected valid thread, got nil")
			}

			if tt.call(task) != tt.expectedvid {
				t.Errorf("Expected id %d, got %d", tt.expectedvid, tt.call(task))
			}
		})
	}
}
