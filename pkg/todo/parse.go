// Copyright 2020 Matt Montgomery
// SPDX-License-Identifier: GPL-3.0-or-later

package todo

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

type Task struct {
	Completed      bool		// If this task is completed
	Priority       string		// Priority for this task if set
	CompletionDate time.Time	// Date this task was completed
	CreationDate   time.Time	// Date this task was created
	Description    string		// Task description, including all tags and contexts
	// Tags           []string		// All tags in the task description
	// Contexts       []string		// All contexts in the task description
	DueDate        time.Time	// Key value pair holding the due date for this task
	// Data           map[string]string	// All key value pairs in the task
}

var EmptyDate = time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)

func (t Task) String() string {
	var complete, priority, completion, creation string

	if t.Completed {
		complete = "x "
	}

	if t.Priority != "" {
		priority = fmt.Sprintf("(%s) ", t.Priority)
	}

	// Check completion and creation
	if t.CompletionDate != EmptyDate {
		completion = format(t.CompletionDate) + " "
	}

	if t.CreationDate != EmptyDate {
		creation = format(t.CreationDate) + " "
	}

	return complete + priority + completion + creation + t.Description
}

func format(d time.Time) string {
	return fmt.Sprintf("%d-%02d-%02d", d.Year(), d.Month(), d.Day())
}

func ParseAll(contents string) []Task {
	var tasks []Task

	// Handles newlines on Windows
	contents = strings.ReplaceAll(contents, "\r", "")
	lines := strings.Split(contents, "\n")

	for _, line := range lines {
		task := ParseTask(line)
		if task.Description != "" {
			tasks = append(tasks, ParseTask(line))
		}
	}

	return tasks
}

// Formatting rules can be found at https://github.com/todotxt/todo.txt
func ParseTask(raw string) Task {
	// completion creation description description description+tag @context due:YYYY-MM-DD
	// (A) 2020-07-02 2020-07-01 task description goes here +tag @context due:2020-07-02

	var task Task

	if len(raw) <= 1 {
		return task
	}

	dateRegex := regexp.MustCompile("^[0-9]{4}-[0-9]{2}-[0-9]{2}")	// 0000-00-00
	dueRegex := regexp.MustCompile("due:[0-9]{4}-[0-9]{2}-[0-9]{2}")	// due:0000-00-00
	priorityRegex := regexp.MustCompile("^\\([A-Z]\\)")	// ([A-Z])
	dateLayout := "2006-01-02"

	// Parse completion status
	// If the task is completed, mark it as such and remove the "x " prefix
	if strings.HasPrefix(raw, "x ") {
		task.Completed = true
		raw = raw[2:]
	}

	// Parse priority
	// If the next field in the string looks like a priority, pop and save it
	priority := strings.Fields(raw)[0]
	if priorityRegex.MatchString(priority) {
		task.Priority = strings.ReplaceAll(priority, "(", "")
		task.Priority = strings.ReplaceAll(task.Priority, ")", "")

		// Remove the priority and trailing space
		remove := len(priority) + 1
		raw = raw[remove:]
	}

	// Parse completion and creation dates
	for i := 0; i <= 1; i++ {
		// Returns the date or "" if nothing was found
		date := dateRegex.FindString(raw)
		if date != "" {
			remove := len(date) + 1
			raw = raw[remove:]

			if i == 0 {
				task.CompletionDate, _ = time.Parse(dateLayout, date)
			} else {
				task.CreationDate, _ = time.Parse(dateLayout, date)
			}
		}
	}

	// Handles the case where only one date is given
	if task.CreationDate == EmptyDate && task.CompletionDate != EmptyDate {
		task.CreationDate = task.CompletionDate
		task.CompletionDate = EmptyDate
	}

	// Parse description
	task.Description = raw

	// Check for a due date
	due := dueRegex.FindString(raw)
	due = strings.ReplaceAll(due, "due:", "")
	task.DueDate, _ = time.Parse(dateLayout, due)

	return task
}
