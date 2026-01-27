package app

import "fmt"

func PrintBuildInfo(buildVersion, buildDate, buildCommit string) {
	fmt.Printf("Build version: %s\n", valueOrNA(buildVersion))
	fmt.Printf("Build date: %s\n", valueOrNA(buildDate))
	fmt.Printf("Build commit: %s\n", valueOrNA(buildCommit))
}

func valueOrNA(v string) string {
	if v == "" {
		return "N/A"
	}

	return v
}
