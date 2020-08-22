// Copyright 2020 Matt Montgomery
// SPDX-License-Identifier: GPL-3.0-or-later

package todo

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"
)

type Task struct {
	Completed      bool			// If this task is completed
	Priority       string		// Priority for this task if set
	CompletionDate time.Time	// Date this task was completed
	CreationDate   time.Time	// Date this task was created
	Description    string		// Task description, including all tags and contexts
	// Tags           []string		// All tags in the task description
	// Contexts       []string		// All contexts in the task description
	DueDate        time.Time	// Key value pair holding the due date for this task
	// Data           map[string]string	// All key value pairs in the task
	Deleted        bool			// If the task was deleted (exclude the task from the list)
	Hash           string		// Unique identifier for this task
}

var EmptyDate = time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)
const seperator = "+=+=+=+=+="

type ByDate []Task
func (a ByDate) Less(i, j int) bool {
	lhs, rhs := a[i], a[j]

	// If the dates are equal, check if one side has the separator since it always goes at the end of that days' tasks
	// Otherwise, just sort normally
	if lhs.DueDate == rhs.DueDate {
		return !strings.HasPrefix(lhs.Description, seperator)
	} else {
		return lhs.DueDate.Before(rhs.DueDate)
	}
}
func (a ByDate) Len() int  { return len(a) }
func (a ByDate) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

func (t Task) String() string {
	var complete, priority, completion, creation string

	if t.Completed {
		complete = "x "
	}

	if t.Priority != "" {
		priority = fmt.Sprintf("(%s) ", t.Priority)
	}

	// Check completion and creation
	if !time.Time.IsZero(t.CompletionDate) {
		completion = format(t.CompletionDate) + " "
	}

	if !time.Time.IsZero(t.CreationDate) {
		creation = format(t.CreationDate) + " "
	}

	// Remove the due date from seperators
	if strings.HasPrefix(t.Description, seperator) {
		t.Description = strings.Fields(t.Description)[0]
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

	if len(raw) == 0 {
		return task
	}

	dateFormat := "[0-9]{4}-[0-9]{2}-[0-9]{2}"
	dateRegex := regexp.MustCompile("^" + dateFormat)	// 0000-00-00
	dueRegex := regexp.MustCompile("due:" + dateFormat)	// due:0000-00-00
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
	if time.Time.IsZero(task.CreationDate) && !time.Time.IsZero(task.CompletionDate) {
		task.CreationDate = task.CompletionDate
		task.CompletionDate = EmptyDate
	}

	// Parse description
	task.Description = raw

	// Check for a due date
	due := dueRegex.FindString(raw)
	due = strings.ReplaceAll(due, "due:", "")
	task.DueDate, _ = time.Parse(dateLayout, due)

	task.Deleted = false

	// Calculate hash
	hash := sha256.Sum256([]byte(task.String()))
	task.Hash = hex.EncodeToString(hash[:])

	return task
}

func SortByDate(raw []Task) []Task {
	sort.Sort(ByDate(raw))
	return raw
}
