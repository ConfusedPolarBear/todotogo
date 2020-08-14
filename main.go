// Copyright 2020 Matt Montgomery
// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/ConfusedPolarBear/todotogo/pkg/todo"
)

/* MVP Functions to implement:
 * Date formats to support:
	due:today
	due:tomorrow
	due:saturday
	due:sat

 * Output cleanup:
	Sort completed tasks at the bottom
 
 * Other potential functions:
 	find/f - loads the contents in multiselect fzf OR an interactive prompt that searches for the given substring
 	With no argument, incomplete tasks from 1-6 days ago should be displayed along with tasks for the next 7 days up to X in each direction

 * ============ Implemented ============
	add/a		add new task
	list/l		list current tasks
 	do/d		mark task X as done
 	rm/r		delete task X
 	archive/ar	move all completed tasks to filename-done.txt
 	edit/e		save the description to a temp file, exec editor and save
 */

 type Tasks = []todo.Task

 var backup bool
 var filename string

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Parse all flags
	filenameFlag := flag.String("f", "todo.txt", "Input filename")
	autoBackupFlag := flag.Bool("b", false, "Disables automatic backup. (dangerous!)")

	flag.Parse()

	filename = *filenameFlag
	backup = !(*autoBackupFlag)
	command := flag.Arg(0)		// optional command (add, rm, etc.)

	// Parse any extra arguments
	args := flag.Args()
	extra := ""
	if len(args) > 1 {
		extra = strings.Join(args[1:], " ")
	}

	// Parse initial task list
	tasks := loadTasks(filename, true)

	if command == "help" || command == "h" {
		printHelp()

	} else if command == "list" || command == "l" || command == "" {
		listTasks(tasks)

	} else if command == "add" || command == "a" {
		task := todo.ParseTask(extra)
		task.CreationDate = time.Now()

		if extra == "" {
			log.Fatalf("Error: you must specify a task")
		}

		backupOriginal(backup, filename)

		tasks = append(tasks, task)
		writeTasks(filename, tasks)

		log.Printf("Successfully added task %s", task)

	} else if (command == "do" || command == "d") {
		markTasks(extra, tasks, true)

	} else if (command == "undo" || command == "u") {
		markTasks(extra, tasks, false)

	} else if command == "rm" || command == "r" {
		_, numbers := numbersToTasks(extra, tasks, "Removed the following tasks:")

		for _, task := range numbers {
			tasks[task].Deleted = true
		}

		writeTasks(filename, tasks)

	} else if command == "archive" || command == "ar" {
		archiveName := strings.ReplaceAll(filename, ".txt", "-done.txt")
		archived := loadTasks(archiveName, false)
		var remaining Tasks

		backupOriginal(backup, filename)

		log.Printf("Archived the following tasks:")
		for _, task := range tasks {
			if !task.Completed {
				remaining = append(remaining, task)
				continue
			}

			archived = append(archived, task)
			log.Printf("%s", task)
		}

		writeTasks(archiveName, archived)
		writeTasks(filename, remaining)

	} else if command == "edit" || command == "e" {
		provided := strings.Fields(extra)
		if len(provided) == 0 {
			log.Fatalf("You must provide at least one task number")
		}

		backupOriginal(backup, filename)

		for index, raw := range provided {
			i, _ := strconv.ParseInt(raw, 10, 32)
			i -= 1

			log.Printf("Editing task %d/%d (%d): %s", index + 1, len(provided), i + 1, tasks[i])
			new := editTask(tasks[i].String())
			log.Printf("New contents of task %d: %s", i + 1, new)

			tasks[i] = todo.ParseTask(new)
		}

		writeTasks(filename, tasks)

	} else {
		log.Printf("Unknown subcommand %s", command)
		printHelp()
	}
}

func printHelp() {
	log.Printf("Available commands:")
	log.Printf("[a]dd:     Adds new task")
	log.Printf("[ar]chive: Moves all completed tasks to FILENAME-done.txt")
	log.Printf("[d]o:      Marks the task(s) as complete")
	log.Printf("[e]dit:    Interactively edit the provided tasks")
	log.Printf("[l]ist:    Lists all tasks (default if no action is specified)")
	log.Printf("[u]ndo:    Marks the task(s) as incomplete")
}

func editTask(original string) string {
	// Create a temporary file to hold the task
	file, tmpErr := ioutil.TempFile("/tmp", "task.")
	if tmpErr != nil {
			log.Fatalf("Unable to create temp file: %s", tmpErr)
	}
	tmp := file.Name()
	defer os.Remove(tmp)

	// Write out the contents of the task
	if err := ioutil.WriteFile(tmp, []byte(original), 0600); err != nil {
		log.Fatalf("Unable to write to temp file: %s", err)
	}

	// Execute the default editor
	cmd := exec.Command("editor", tmp)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatalf("Unable to execute editor: %s", err)
	}

	// Read back the contents of the file and return the processed output
	raw, err := ioutil.ReadFile(tmp)
	if err != nil {
		log.Fatalf("Unable to open %s: %s", tmp, err)
	}
	contents := strings.ReplaceAll(string(raw), "\n", " ")

	return contents
}

func backupOriginal(enabled bool, filename string) {
	if !enabled {
		return
	}

	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("Unable to open %s: %s", filename, err)
	}

	backup := filename + ".bak"
	backupErr := ioutil.WriteFile(backup, contents, 0644)

	if backupErr != nil {
		log.Fatalf("Unable to create backup %s: %s", backup, backupErr)
	}

	log.Printf("Backed up original file %s as %s", filename, backup)
}

func markTasks(input string, tasks Tasks, complete bool) {
	msg := "Marked the following tasks as complete:"
	if !complete {
		msg = strings.ReplaceAll(msg, "complete", "incomplete")
	}

	_, numbers := numbersToTasks(input, tasks, msg)

	for _, task := range numbers {
		tasks[task].Completed = complete
	}

	writeTasks(filename, tasks)
}

func loadTasks(filename string, fatal bool) Tasks {
	raw, err := ioutil.ReadFile(filename)
	if err != nil {
		if fatal {
			log.Fatalf("Unable to open %s: %s", filename, err)
		} else {
			raw = []byte{}
		}
	}

	tasks := todo.ParseAll(string(raw))
	return todo.SortByDate(tasks)
}

func writeTasks(filename string, tasks Tasks) {
	contents := ""

	for _, task := range tasks {
		if task.Deleted {
			continue
		}

		contents += fmt.Sprintf("%s\n", task)
	}

	ioutil.WriteFile(filename, []byte(contents), 0644)
}

func listTasks(tasks Tasks) {
	for number, task := range tasks {
		fmt.Printf("%03d %s\n", number + 1, task)
	}
}

func listNumberedTasks(tasks Tasks, numbers []int) {
	for i, task := range tasks {
		fmt.Printf("%03d %s\n", numbers[i] + 1, task)
	}
}

// rawNumbers is a string of space seperated numbers ("1 2 6") and returns the tasks that correspond to those numbers
func numbersToTasks(rawNumbers string, tasks Tasks, msg string) (Tasks, []int) {
	var ret Tasks
	var parsed []int

	numbers := strings.Fields(rawNumbers)

	for _, i := range numbers {
		// Subtract 1 from the task number since listTasks adds 1
		longIndex, err := strconv.ParseInt(i, 10, 32)
		index := int(longIndex) - 1

		if index < 0 || index > len(tasks) || err != nil {
			log.Fatalf("Error: cannot find task with index %d", index)
		}

		ret = append(ret, tasks[index])
		parsed = append(parsed, index)
	}

	if msg != "" {
		backupOriginal(backup, filename)

		log.Printf(msg)
		listNumberedTasks(ret, parsed)
	}

	return ret, parsed
}
