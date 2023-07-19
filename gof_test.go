package main

import (
	"io/ioutil"
	"os"
	"testing"

	. "github.com/flocks/gof/parse"
)

func init() {
	os.Setenv("PWD", "$PWD")
	FileExist = fileExist
}

func fileExist(path string) bool {
	files := []string{
		"/home/fteissier/ledger/vault-ts/apps/cli/src/cli.ts",
		"/home/fteissier/ledger/vault-ts/apps/cli/src/index.ts",
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
	input, _ := ioutil.ReadFile("./samples/pnpm-lint.txt")
	result := FindFiles(string(input))
	expected := []Filematch{
		{
			FilePath: "/home/fteissier/ledger/vault-ts/apps/cli/src/cli.ts",
			Line:     23,
			Col:      7,
			Desc:     "  warning  'a' is assigned a value but never used  @typescript-eslint/no-unused-vars",
		},
		{
			FilePath: "/home/fteissier/ledger/vault-ts/apps/cli/src/cli.ts",
			Line:     194,
			Col:      7,
			Desc:     "  warning  'b' is assigned a value but never used  @typescript-eslint/no-unused-vars",
		},
		{
			FilePath: "/home/fteissier/ledger/vault-ts/apps/cli/src/registerTransports.ts",
			Line:     14,
			Col:      8,
			Desc:     "  warning  'b' is assigned a value but never used  @typescript-eslint/no-unused-vars",
		},
	}
	compareFiles(expected, result, t)
}
func TestLinterUnix(t *testing.T) {
	input, _ := ioutil.ReadFile("./samples/pnpm-lint-format-unix.txt")
	result := FindFiles(string(input))
	expected := []Filematch{
		{
			FilePath: "/home/fteissier/ledger/vault-ts/apps/cli/src/cli.ts",
			Line:     12,
			Col:      7,
			Desc:     " 'a' is assigned a value but never used. [Warning/@typescript-eslint/no-unused-vars]",
		},
		{
			FilePath: "/home/fteissier/ledger/vault-ts/apps/cli/src/index.ts",
			Line:     7,
			Col:      7,
			Desc:     " 'a' is assigned a value but never used. [Warning/@typescript-eslint/no-unused-vars]",
		},
		{
			FilePath: "/home/fteissier/ledger/vault-ts/apps/cli/src/registerTransports.ts",
			Line:     14,
			Col:      7,
			Desc:     " 'b' is assigned a value but never used. [Warning/@typescript-eslint/no-unused-vars]",
		},
	}
	compareFiles(expected, result, t)
}
func TestGrep(t *testing.T) {
	input, _ := ioutil.ReadFile("./samples/grep.txt")
	result := FindFiles(string(input))
	expected := []Filematch{
		{
			FilePath: "$PWD/src/components/DeviceInteraction/index.jsx",
			Line:     62,
			Col:      11,
			Desc:     "  appVersion: string,",
		},
		{
			FilePath: "$PWD/src/components/Onboarding/index.jsx",
			Line:     10,
			Col:      23,
			Desc:     "    appVersion: window.config.APP_VERSION,",
		},
	}
	compareFiles(expected, result, t)
}

func TestGitStatus(t *testing.T) {
	input, _ := ioutil.ReadFile("./samples/git-show.txt")
	result := FindFiles(string(input))
	expected := []Filematch{
		{
			FilePath: "$PWD/src/Root.jsx",
		},
		{
			FilePath: "$PWD/src/components/DeviceInteraction/index.jsx",
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
		if val.CompareWith(item) {
			return true
		}
	}
	return false
}

func compareFiles(expected []Filematch, result []Filematch, t *testing.T) bool {
	for _, val := range result {
		if !containsF(expected, val) {
			t.Fatalf(`file %v is missing in expected`, val.FilePath)
		}
	}
	for _, val := range expected {
		if !containsF(result, val) {
			t.Fatalf(`file %v is missing in actual result`, val.FilePath)
		}
	}

	return true
}
