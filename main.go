// Copyright 2020 Matt Montgomery
// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"io/ioutil"
	"flag"
	"fmt"
	"log"
	"strings"
	"strconv"
	"time"

	"github.com/ConfusedPolarBear/todotogo/pkg/todo"
)

/* MVP Functions to implement:
 * Date formats to support:
	due:today
	due:tomorrow
	due:saturday
	due:sat
	
 * rm/r			delete task X
 * do/d			mark task X as done
 * archive/ar	move all completed tasks to filename-archive.txt
 * edit/e		save the description to a temp file, exec editor and save

 Sort completed tasks at the bottom

 * Implemented
 * add/a
 * list/l

 * Other potential functions:
 * find/f - loads the contents in multiselect fzf OR an interactive prompt that searches for the given substring
 * With no argument, incomplete tasks from 1-6 days ago should be displayed along with tasks for the next 7 days up to X in each direction
 */

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Parse all flags
	filenameFlag := flag.String("f", "todo.txt", "Input filename")
	// outputFlag := flag.String("o", "", "Output filename")
	autoBackupFlag := flag.Bool("b", false, "Disables automatic backup. (dangerous!)")
	
	flag.Parse()

	filename := *filenameFlag
	// outputFilename := *outputFlag
	backup := !(*autoBackupFlag)
	command := flag.Arg(0)		// optional command (add, rm, etc.)

	// Parse input file
	raw, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("Unable to open %s: %s", filename, err)
	}
	tasks := todo.ParseAll(string(raw))

	if command == "help" || command == "h" {
		printHelp()

	} else if command == "list" || command == "l" || command == "" {
		listTasks(tasks)

	} else if command == "add" || command == "a" {
		description := strings.Join(flag.Args()[1:], " ")
		task := todo.ParseTask(description)
		task.CreationDate = time.Now()

		if description == "" {
			log.Fatalf("Error: you must specify a task")
		}

		backupOriginal(backup, filename, raw)

		/*if outputFilename != "" {
			filename = outputFilename
		}*/

		appendTask(filename, raw, task)
		log.Printf("Successfully added task %s", task)

	} else if command == "do" || command == "d" {
		provided, numbers := numbersToTasks(flag.Args()[1:], tasks)
		
		log.Printf("Marking the following tasks as completed:")
		listNumberedTasks(provided, numbers)
		
		backupOriginal(backup, filename, raw)

		for _, task := range numbers {
			log.Printf("marking number %d as complete", task)
		}

	} else {
		log.Printf("Unknown subcommand %s", command)
		printHelp()
	}
}

func printHelp() {
	log.Printf("Available commands:")
	log.Printf("[l]ist: List all tasks")
	log.Printf("[a]dd:  Add new task")
	log.Printf("[d]o:   Marks the task(s) as completed")
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

func listNumberedTasks(tasks []todo.Task, numbers []int) {
	todo.SortByDate(tasks)
	for i, task := range tasks {
		fmt.Printf("%03d %s\n", numbers[i] + 1, task)
	}
}

// This function accepts a slice of strings { "1", "2", "6" } and returns the tasks that correspond to those numbers.
// TODO: unit test this
func numbersToTasks(numbers []string, tasks []todo.Task) ([]todo.Task, []int) {
	var ret []todo.Task
	var parsed []int

	for _, i := range numbers {
		// Subtract 1 from the task number since listTasks adds 1
		longIndex, err := strconv.ParseInt(i, 10, 32)
		index := int(longIndex) - 1

		log.Printf("Parsed task number %d", index)

		if index < 0 || index > len(tasks) || err != nil {
			log.Fatalf("Error: cannot find task with index %d", index)
		}

		ret = append(ret, tasks[index])
		parsed = append(parsed, index)
	}

	return ret, parsed
}