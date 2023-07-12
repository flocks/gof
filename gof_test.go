package main

import (
	"io/ioutil"
	"os"
	"testing"
)

func fileExist(path string) bool {
	os.Setenv("PWD", "/path/to/pwd")
	files := []string{
		"/home/fteissier/ledger/vault-ts/apps/cli/src/cli.ts",
		"/home/fteissier/ledger/vault-ts/apps/cli/src/registerTransports.ts",
		"/path/to/pwd/src/Root.jsx",
		"/path/to/pwd/src/components/DeviceInteraction/index.jsx",
		"/path/to/pwd/src/components/Onboarding/index.jsx",
		"/path/to/pwd/src/components/legacy/DeviceInteraction/index.jsx",
		"/path/to/pwd/src/device/interactions/generateWrappingKeys.js",
		"/path/to/pwd/src/device/interactions/hsmFlows.js",
		"/path/to/pwd/src/network/fetchF.js",
		"/path/to/pwd/src/notifications/Notification.jsx",
		"/path/to/pwd/src/notifications/useListenEvents.js",
	}
	if contains(files, path) {
		return true
	}
	return false
}

func TestLinter(t *testing.T) {
	FileExist = fileExist
	input, _ := ioutil.ReadFile("./samples/pnpm-lint.txt")
	result := FindFiles(string(input))
	expected := []Filematch{
		{
			filePath: "/home/fteissier/ledger/vault-ts/apps/cli/src/cli.ts",
			line:     0,
			col:      0,
		},
		{
			filePath: "/home/fteissier/ledger/vault-ts/apps/cli/src/registerTransports.ts",
			line:     0,
			col:      0,
		},
	}
	compareFiles(expected, result, t)
}
func TestLinterUnix(t *testing.T) {
	FileExist = fileExist
	input, _ := ioutil.ReadFile("./samples/pnpm-lint-format-unix.txt")
	result := FindFiles(string(input))
	expected := []Filematch{
		{
			filePath: "/home/fteissier/ledger/vault-ts/apps/cli/src/cli.ts",
			line:     0,
			col:      0,
		},
		{
			filePath: "/home/fteissier/ledger/vault-ts/apps/cli/src/registerTransports.ts",
			line:     0,
			col:      0,
		},
	}
	compareFiles(expected, result, t)
}

func TestGitStatus(t *testing.T) {
	input, _ := ioutil.ReadFile("./samples/git-show.txt")
	result := FindFiles(string(input))
	expected := []Filematch{
		{
			filePath: "/path/to/pwd/src/Root.jsx",
			line:     0,
			col:      0,
		},
		{
			filePath: "/path/to/pwd/src/components/DeviceInteraction/index.jsx",
			line:     0,
			col:      0,
		},
		{
			filePath: "/path/to/pwd/src/components/Onboarding/index.jsx",
			line:     0,
			col:      0,
		},
		{
			filePath: "/path/to/pwd/src/components/legacy/DeviceInteraction/index.jsx",
			line:     0,
			col:      0,
		},
		{
			filePath: "/path/to/pwd/src/device/interactions/generateWrappingKeys.js",
			line:     0,
			col:      0,
		},
		{
			filePath: "/path/to/pwd/src/device/interactions/hsmFlows.js",
			line:     0,
			col:      0,
		},
		{
			filePath: "/path/to/pwd/src/network/fetchF.js",
			line:     0,
			col:      0,
		},
		{
			filePath: "/path/to/pwd/src/notifications/Notification.jsx",
			line:     0,
			col:      0,
		},
		{
			filePath: "/path/to/pwd/src/notifications/useListenEvents.js",
			line:     0,
			col:      0,
		},
	}
	compareFiles(expected, result, t)
}

func contains(arr []string, item string) bool {
	for _, val := range arr {
		if val == item {
			return true
		}
	}
	return false
}
func containsF(arr []Filematch, item Filematch) bool {
	for _, val := range arr {
		if val == item {
			return true
		}
	}
	return false
}

func compareFiles(expected []Filematch, result []Filematch, t *testing.T) bool {
	for _, val := range result {
		if !containsF(expected, val) {
			t.Fatalf(`file %v is missing in expected`, val)
		}
	}
	for _, val := range expected {
		if !containsF(result, val) {
			t.Fatalf(`file %v is missing in actual result`, val)
		}
	}

	return true
}
