package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

var HelmArgs []string
var secretPrefix = "@secretsmanager@"
var secretsmanagerPattern = regexp.MustCompile(`.+` + secretPrefix + `(arn:aws:secretsmanager:[a-z]{2}-[a-z]+-[0-9]:[0-9]+:[a-z]+.+)$`)
var cleartextFile string
var cleartextFiles []string
var splicedArg []string
var valueName = ""
var fileName string

func main() {

	for _, arg := range os.Args[1:] {
		if strings.HasSuffix(arg, ".yaml") || strings.HasSuffix(arg, ".yml") {
			splicedArg = strings.SplitAfter(arg, "=")
			if len(splicedArg) > 1 {
				valueName = splicedArg[0]
			}
			fileName = splicedArg[len(splicedArg) - 1]
			if containsSecrets(fileName) {

				// helm command will be run with the new file containing secrets
				cleartextFile = replaceSecrets(fileName, getSecretsmanagerSecret)
				HelmArgs = append(HelmArgs, valueName + cleartextFile)
				cleartextFiles = append(cleartextFiles, cleartextFile)
			} else {
				HelmArgs = append(HelmArgs, arg)
			}

		} else if arg == "--keep" {
			fmt.Fprintln(os.Stderr, "Warning: The flag `--keep` will be deprecated in the next version.")
		} else {
			HelmArgs = append(HelmArgs, arg)
		}
	}

	fmt.Println(strings.Join(HelmArgs, " "))
}
