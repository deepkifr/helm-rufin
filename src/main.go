package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"log"
)

var HelmArgs []string
var secretPrefix = "@secretsmanager@"
var secretsmanagerPattern = regexp.MustCompile(`.+` + secretPrefix + `(arn:aws:secretsmanager:[a-z]{2}-[a-z]+-[0-9]:[0-9]+:[a-z]+.+)$`)
var cleartextFile string
var cleartextFiles []string
var keepCleartext bool

func main() {

	for _, arg := range os.Args[1:] {
		if strings.HasSuffix(arg, ".yaml") {
			if containsSecrets(arg) {

				// helm command will be run with the new file containing secrets
				cleartextFile = replaceSecrets(arg, getSecretsmanagerSecret)
				HelmArgs = append(HelmArgs, cleartextFile)
				cleartextFiles = append(cleartextFiles, cleartextFile)
			} else {
				HelmArgs = append(HelmArgs, arg)
			}

		} else if arg == "--keep" {
			keepCleartext = true
		} else {
				HelmArgs = append(HelmArgs, arg)
		}
	}

	fmt.Println(strings.Join(HelmArgs, " "))

	if !keepCleartext {
		for _, file := range cleartextFiles {
	  		fmt.Fprintln(os.Stderr, "Deleting cleartext file : ", file)
	  		e := os.Remove(file)
	  		if e != nil {
		  		log.Fatal(e)
	  		}
		}
	}
}
