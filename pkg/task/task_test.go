package task

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/Andreagit97/bpftree/pkg/render"
	"github.com/enescakir/emoji"
)

func TestGetReaperImage(t *testing.T) {
	var tests = []struct {
		name          string
		task          task
		expectedImage emoji.Emoji
	}{
		{"IsReaper", task{Info: &taskInfo{IsChildSubreaper: 1}}, render.GetImageReaper()},
		{"IsNotReaper", task{Info: &taskInfo{IsChildSubreaper: 0}}, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.task.getReaperImage(); got != tt.expectedImage {
				t.Errorf("Expected:\n%v\n\nGot:\n%v\n", tt.expectedImage, got)
			}
		})
	}
}

func TestUsers(t *testing.T) {
	var tests = []struct {
		name             string
		task             task
		expectedName     string
		expectedUserList map[string]string
		setup            func()
	}{
		{
			name:         "InvalidUser",
			task:         task{Info: &taskInfo{EUID: -1, LoginUID: -1}},
			expectedName: "",
			expectedUserList: map[string]string{
				"1000": "testuser",
			},
			setup: func() {
				userList = make(map[string]string)
				userList["1000"] = "testuser"
			},
		},
		{
			name:         "ExistingUserInTheTable",
			task:         task{Info: &taskInfo{EUID: 1000, LoginUID: 1000}},
			expectedName: "testuser",
			expectedUserList: map[string]string{
				"1000": "testuser",
			},
			setup: func() {
				userList = make(map[string]string)
				userList["1000"] = "testuser"
			},
		},
		{
			name: "ExistingUserInTheSystem",
			// Hopefully the root user is always present
			task:         task{Info: &taskInfo{EUID: 0, LoginUID: 0}},
			expectedName: "root",
			expectedUserList: map[string]string{
				"0": "root",
			},
			setup: func() {
				userList = make(map[string]string)
			},
		},
		{
			name:             "NotExistingUser",
			task:             task{Info: &taskInfo{EUID: 999999, LoginUID: 999999}},
			expectedName:     "",
			expectedUserList: map[string]string{},
			setup: func() {
				userList = make(map[string]string)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			if got := tt.task.getEUIDName(); got != tt.expectedName {
				t.Errorf("Expected:\n%s\n\nGot:\n%s\n", tt.expectedName, got)
			}
			if got := tt.task.getLoginUIDName(); got != tt.expectedName {
				t.Errorf("Expected:\n%s\n\nGot:\n%s\n", tt.expectedName, got)
			}
			if !reflect.DeepEqual(userList, tt.expectedUserList) {
				t.Errorf("Expected:\n%v\n\nGot:\n%v\n", tt.expectedUserList, userList)
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
		{"GetComm",
			task{Info: &taskInfo{
				Comm: [commLen]byte{'0'},
			}},
			"0",
			func(t task) string { return t.getComm() }},
		{"GetCommWithoutTerminator",
			task{Info: &taskInfo{
				Comm: [commLen]byte{
					'p', 'p', 'p', 'p', 'p', 'p', 'p', 'p', 'p', 'p', 'p', 'p', 'p', 'p', 'p', 'p',
				},
			}},
			"pppppppppppppppp",
			func(t task) string { return t.getComm() }},
		{"GetCommBytes",
			task{Info: &taskInfo{
				Comm: [commLen]byte{99, 111, 110, 116, 97, 105, 110, 101, 114, 100},
			}},
			"containerd",
			func(t task) string { return t.getComm() }},
		{"GetCommContainerd",
			task{Info: &taskInfo{
				Comm: [commLen]byte{'c', 'o', 'n', 't', 'a', 'i', 'n', 'e', 'r', 'd'},
			}},
			"containerd",
			func(t task) string { return t.getComm() }},
		{"GetCommEmpty",
			task{Info: &taskInfo{
				Comm: [commLen]byte{},
			}},
			"",
			func(t task) string { return t.getComm() }},
		{"GetFormattedCommMainThread", task{Info: &taskInfo{Tid: 1, Pid: 1, Comm: [commLen]byte{'c', 'o', 'm', 'm'}}}, "[comm]", func(t task) string { return t.getFormattedComm() }},
		{"GetFormattedCommSecondaryThread", task{Info: &taskInfo{Tid: 2, Pid: 1, Comm: [commLen]byte{'c', 'o', 'm', 'm'}}}, "{comm}", func(t task) string { return t.getFormattedComm() }},
		{"GetTid",
			task{Info: &taskInfo{
				Tid: 1,
			}},
			"1",
			func(t task) string { return fmt.Sprint(t.GetTid()) }},
		{"GetPid",
			task{Info: &taskInfo{
				Pid: 2,
			}},
			"2",
			func(t task) string { return fmt.Sprint(t.getPid()) }},
		{"GetParentTid",
			task{Info: &taskInfo{
				ParentTid: 3,
			}},
			"3",
			func(t task) string { return fmt.Sprint(t.getParentTid()) }},
		{"GetParentPid",
			task{Info: &taskInfo{
				ParentPid: 4,
			}},
			"4",
			func(t task) string { return fmt.Sprint(t.getParentPid()) }},
		{"GetRealParentTid",
			task{Info: &taskInfo{
				RealParentTid: 5,
			}},
			"5",
			func(t task) string { return fmt.Sprint(t.GetRealParentTid()) }},
		{"GetRealParentPid",
			task{Info: &taskInfo{
				RealParentPid: 6,
			}},
			"6",
			func(t task) string { return fmt.Sprint(t.getRealParentPid()) }},
		{"GetReaperTrue",
			task{Info: &taskInfo{
				IsChildSubreaper: 1,
			}},
			"true",
			func(t task) string { return fmt.Sprint(t.isChildSubReaper()) }},
		{"GetReaperFalse",
			task{Info: &taskInfo{
				IsChildSubreaper: 0,
			}},
			"false",
			func(t task) string { return fmt.Sprint(t.isChildSubReaper()) }},
		{"isMainThreadTrue",
			task{Info: &taskInfo{
				Tid: 1,
				Pid: 1,
			}},
			"true",
			func(t task) string { return fmt.Sprint(t.isMainThread()) }},
		{"isMainThreadFalse",
			task{Info: &taskInfo{
				Tid: 1,
				Pid: 2,
			}},
			"false",
			func(t task) string { return fmt.Sprint(t.isMainThread()) }},
		{"GetNsLevel",
			task{Info: &taskInfo{
				NsLevel: 7,
			}},
			"7",
			func(t task) string { return fmt.Sprint(t.getNsLevel()) }},
		{"GetVTid",
			task{Info: &taskInfo{
				VTid: 8,
			}},
			"8",
			func(t task) string { return fmt.Sprint(t.getVTid()) }},
		{"GetPgid",
			task{Info: &taskInfo{
				Pgid: 10,
			}},
			"10",
			func(t task) string { return fmt.Sprint(t.getPgid()) }},
		{"GetSid",
			task{Info: &taskInfo{
				Sid: 12,
			}},
			"12",
			func(t task) string { return fmt.Sprint(t.getSid()) }},
		{"GetExePath",
			task{Info: &taskInfo{
				ExePath: [exePathLen]byte{'/', 'u', 's', 'r', '/'},
			}},
			"/usr/",
			func(t task) string { return t.getExePath() }},
		{"GetExePathBytes",
			task{Info: &taskInfo{
				ExePath: [exePathLen]byte{99, 111, 110, 116, 97, 105, 110, 101, 114, 100},
			}},
			"containerd",
			func(t task) string { return t.getExePath() }},
		{"GetExePathContainerd",
			task{Info: &taskInfo{
				ExePath: [exePathLen]byte{'c', 'o', 'n', 't', 'a', 'i', 'n', 'e', 'r', 'd'},
			}},
			"containerd",
			func(t task) string { return t.getExePath() }},
		{"GetExePathEmpty",
			task{Info: &taskInfo{
				ExePath: [exePathLen]byte{},
			}},
			"",
			func(t task) string { return t.getExePath() }},
		{"GetLoginUid",
			task{Info: &taskInfo{
				LoginUID: -2,
			}},
			"-2",
			func(t task) string { return fmt.Sprint(t.getLoginUID()) }},
		{"GetEUid",
			task{Info: &taskInfo{
				EUID: -1,
			}},
			"-1",
			func(t task) string { return fmt.Sprint(t.getEUID()) }},
		{"GetCmdLine",
			task{Info: &taskInfo{
				CmdLine: [exePathLen]byte{'e', 'x', 'e', 0, 'a', 'r', 'g', '1', 0, 'a', 'r', 'g', '2', 0},
			}},
			"exe,arg1,arg2",
			func(t task) string { return t.getCmdLine() }},
		{"GetCmdLineExtra0",
			task{Info: &taskInfo{
				CmdLine: [exePathLen]byte{'e', 'x', 'e', 0, 'a', 'r', 'g', '1', 0, 'a', 'r', 'g', '2', 0, 0, 0},
			}},
			"exe,arg1,arg2",
			func(t task) string { return t.getCmdLine() }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.fieldGetter(tt.task) != tt.expectedField {
				t.Errorf("Expected:\n%s\n\nGot:\n%s\n", tt.expectedField, tt.fieldGetter(tt.task))
			}
		})
	}
}
