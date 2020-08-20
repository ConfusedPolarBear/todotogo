// Copyright 2020 Matt Montgomery
// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/ConfusedPolarBear/todotogo/pkg/todo"
)

func getMessage(task string, prop string, expect interface{}, actual interface{}) string {
	return fmt.Sprintf("Task \"%s\" failed (%s).\nExpected: %v\nActual:   %v", task, prop, expect, actual)
}

func areTasksEqual(rawTask string, expected *todo.Task, t *testing.T) {
	actual := todo.ParseTask(rawTask)

	if expected.Completed != actual.Completed {
		t.Errorf(getMessage(rawTask, "completed", expected.Completed, actual.Completed))
	}

	if expected.Priority != actual.Priority {
		t.Errorf(getMessage(rawTask, "priority", expected.Priority, actual.Priority))
	}

	if expected.CompletionDate != actual.CompletionDate {
		t.Errorf(getMessage(rawTask, "completion date", expected.CompletionDate, actual.CompletionDate))
	}

	if expected.CreationDate != actual.CreationDate {
		t.Errorf(getMessage(rawTask, "creation date", expected.CreationDate, actual.CreationDate))
	}

	if expected.Description != actual.Description {
		t.Errorf(getMessage(rawTask, "description", expected.Description, actual.Description))
	}

	if expected.DueDate != actual.DueDate {
		t.Errorf(getMessage(rawTask, "due date", expected.DueDate, actual.DueDate))
	}

	if expected.Hash != actual.Hash {
		t.Errorf(getMessage(rawTask, "hash", expected.Hash, actual.Hash))
	}

	if rawTask != expected.String() {
		t.Errorf(getMessage(rawTask, "serialization", rawTask, expected))
	}
}

func TestFullTask(t *testing.T) {
	areTasksEqual("(C) priority C +test due:2020-07-01", &todo.Task{
		Completed: false,
		Priority: "C",
		Description: "priority C +test due:2020-07-01",
		DueDate: time.Date(2020, 7, 1, 0, 0, 0, 0, time.UTC),
		Deleted: false,
		Hash: "d7b12180ca7f611e0da354e6b2cf8eac03d14337dae32815dacab6ef962556cf",
	}, t)
}

func TestFullTaskComplete(t *testing.T) {
	areTasksEqual("x (C) priority C +test due:2020-07-01", &todo.Task{
		Completed: true,
		Priority: "C",
		Description: "priority C +test due:2020-07-01",
		DueDate: time.Date(2020, 7, 1, 0, 0, 0, 0, time.UTC),
		Deleted: false,
		Hash: "8585f8b77bc59ae94ca874872c0bd964ce2ddb7f9916412680fdbd9ddd5bac7d",
	}, t)
}

func TestCompleteWithDates(t *testing.T) {
	areTasksEqual("x (A) 2016-05-20 2016-04-30 measure space for +chapelShelving @chapel due:2016-05-30", &todo.Task{
		Completed: true,
		Priority: "A",
		Description: "measure space for +chapelShelving @chapel due:2016-05-30",
		DueDate: time.Date(2016, 5, 30, 0, 0, 0, 0, time.UTC),
		CompletionDate: time.Date(2016, 5, 20, 0, 0, 0, 0, time.UTC),
		CreationDate: time.Date(2016, 4, 30, 0, 0, 0, 0, time.UTC),
		Deleted: false,
		Hash: "6b44b2a3a47f9cb66b9fee123a42052e7d4f3ebedd3ef3f96dab0b88d0ff2aed",
	}, t)
}

func TestOnlyCreationDate(t *testing.T) {
	areTasksEqual("2020-03-20 Create a centralized dotfiles repo due:2020-03-26", &todo.Task{
		Completed: false,
		Priority: "",
		Description: "Create a centralized dotfiles repo due:2020-03-26",
		DueDate: time.Date(2020, 3, 26, 0, 0, 0, 0, time.UTC),
		CompletionDate: time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC),
		CreationDate: time.Date(2020, 3, 20, 0, 0, 0, 0, time.UTC),
		Deleted: false,
		Hash: "22923561a0c012f499528c84668a5550459dfd94c625f44be4c5c2f4c1d541af",
	}, t)
}

func TestNumberedTasks(t *testing.T) {
	var valid, tasks []todo.Task

	valid = append(valid, todo.ParseTask("2020-08-11 first valid task"))
	valid = append(valid, todo.ParseTask("2020-08-11 second +valid task"))
	valid = append(valid, todo.ParseTask("2020-08-11 final +valid task due:1970-01-02"))

	tasks = append(tasks, todo.ParseTask("2020-08-11 junk task 1 tag:asdf"))
	tasks = append(tasks, valid[0])
	tasks = append(tasks, todo.ParseTask("2020-08-11 junk task 2 tag:asdf2"))
	tasks = append(tasks, todo.ParseTask("2020-08-11 (A) junk task 3 tag:asdf"))
	tasks = append(tasks, valid[1])
	tasks = append(tasks, todo.ParseTask("2020-08-11 (C) junk task 4 +asdfasdfasdf tag:asdf"))
	tasks = append(tasks, valid[2])

	selected, indexes := numbersToTasks("2 5 7", tasks, "")
	if len(selected) != len(indexes) || len(selected) != 3 {
		t.Errorf("Error testing numbered tasks: Expected length of (3, 3) but got (%d, %d)", len(selected), len(indexes))
	}

	for i, task := range selected {
		lhs := task.String()
		rhs := valid[i].String()
		if(lhs != rhs) {
			t.Errorf(getMessage(rhs, "numbered tasks", rhs, lhs))
		}
	}
}
