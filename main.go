// Copyright 2020 Matt Montgomery
// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"io/ioutil"
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/ConfusedPolarBear/todotogo/pkg/todo"
)

/* MVP Functions to implement:
 * add/a
 * rm/r
 * do/d
 * archive/ar
 * edit/e
 * list/l

 * Other potential functions:
 * find/f - loads the contents in multiselect fzf OR an interactive prompt that searches for the given substring
 * With no argument, incomplete tasks from 1-6 days ago should be displayed along with tasks for the next 7 days up to X in each direction
 */

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Parse all flags
	filenameFlag := flag.String("f", "todo.txt", "Input filename")
	outputFlag := flag.String("o", "", "Output filename")
	autoBackupFlag := flag.Bool("b", false, "Disables automatic backup. (dangerous!)")
	
	flag.Parse()

	filename := *filenameFlag
	outputFilename := *outputFlag
	backup := !(*autoBackupFlag)
	command := flag.Arg(0)		// optional command (add, rm, etc.)

	// Parse input file
	raw, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("Unable to open %s: %s", filename, err)
	}
	tasks := todo.ParseAll(string(raw))

	if command == "list" || command == "" {
		listTasks(tasks)

	} else if command == "add" || command == "a" {
		backupOriginal(backup, filename, raw)

		description := strings.Join(flag.Args()[1:], " ")
		task := todo.ParseTask(description)
		task.CreationDate = time.Now()

		if outputFilename != "" {
			filename = outputFilename
		}

		appendTask(filename, raw, task)
		log.Printf("Successfully added task %s", task)

	} else {
		log.Printf("Unknown subcommand %s", command)
	}
}

func backupOriginal(enabled bool, filename string, contents []byte) {
	if !enabled {
		return
	}

	backup := filename + ".bak"
	backupErr := ioutil.WriteFile(backup, contents, 0644)
	
	log.Printf("Backing up original file %s as %s", filename, backup)

	if backupErr != nil {
		log.Fatalf("Unable to create backup %s: %s", backup, backupErr)
	}
}

func appendTask(filename string, original []byte, task todo.Task) {
	contents := string(original)
	contents += task.String() + "\n"

	ioutil.WriteFile(filename, []byte(contents), 0644)
}

func listTasks(tasks []todo.Task) {
	todo.SortByDate(tasks)
	for number, task := range tasks {
		fmt.Printf("%03d %s\n", number + 1, task)
	}
}
