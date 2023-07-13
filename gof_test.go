package main

import (
	"io/ioutil"
	"os"
	"testing"
)

func fileExist(path string) bool {
	os.Setenv("PWD", "$PWD")
	files := []string{
		"/home/fteissier/ledger/vault-ts/apps/cli/src/cli.ts",
		"/home/fteissier/ledger/vault-ts/apps/cli/src/registerTransports.ts",
		"$PWD/src/Root.jsx",
		"$PWD/src/components/DeviceInteraction/index.jsx",
		"$PWD/src/components/Onboarding/index.jsx",
		"$PWD/src/components/legacy/DeviceInteraction/index.jsx",
		"$PWD/src/device/interactions/generateWrappingKeys.js",
		"$PWD/src/device/interactions/hsmFlows.js",
		"$PWD/src/network/fetchF.js",
		"$PWD/src/notifications/Notification.jsx",
		"$PWD/src/notifications/useListenEvents.js",
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
func TestGrep(t *testing.T) {
	FileExist = fileExist
	input, _ := ioutil.ReadFile("./samples/grep.txt")
	result := FindFiles(string(input))
	expected := []Filematch{
		{
			filePath: "$PWD/src/components/DeviceInteraction/index.jsx",
			line:     0,
			col:      0,
		},
		{
			filePath: "$PWD/src/components/Onboarding/index.jsx",
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
			filePath: "$PWD/src/Root.jsx",
			line:     0,
			col:      0,
		},
		{
			filePath: "$PWD/src/components/DeviceInteraction/index.jsx",
			line:     0,
			col:      0,
		},
		{
			filePath: "$PWD/src/components/Onboarding/index.jsx",
			line:     0,
			col:      0,
		},
		{
			filePath: "$PWD/src/components/legacy/DeviceInteraction/index.jsx",
			line:     0,
			col:      0,
		},
		{
			filePath: "$PWD/src/device/interactions/generateWrappingKeys.js",
			line:     0,
			col:      0,
		},
		{
			filePath: "$PWD/src/device/interactions/hsmFlows.js",
			line:     0,
			col:      0,
		},
		{
			filePath: "$PWD/src/network/fetchF.js",
			line:     0,
			col:      0,
		},
		{
			filePath: "$PWD/src/notifications/Notification.jsx",
			line:     0,
			col:      0,
		},
		{
			filePath: "$PWD/src/notifications/useListenEvents.js",
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
			t.Fatalf(`file %v is missing in expected`, val.filePath)
		}
	}
	for _, val := range expected {
		if !containsF(result, val) {
			t.Fatalf(`file %v is missing in actual result`, val.filePath)
		}
	}

	return true
}
