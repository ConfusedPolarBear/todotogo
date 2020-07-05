// Copyright 2020 Matt Montgomery
// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"testing"
	"time"
	"fmt"

	"github.com/ConfusedPolarBear/todotogo/pkg/todo"
)

func getMessage(task string, prop string, expect interface{}, actual interface{}) string {
	return fmt.Sprintf("Task %s failed (%s).\nExpected: %v\nActual:   %v", task, prop, expect, actual)
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
	}, t)
}

func TestFullTaskComplete(t *testing.T) {
	areTasksEqual("x (C) priority C +test due:2020-07-01", &todo.Task{
		Completed: true,
		Priority: "C",
		Description: "priority C +test due:2020-07-01",
		DueDate: time.Date(2020, 7, 1, 0, 0, 0, 0, time.UTC),
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
	}, t)
}