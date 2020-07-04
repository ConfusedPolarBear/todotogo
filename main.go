// Copyright 2020 Matt Montgomery
// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"log"
	"io/ioutil"
	"fmt"

	"./pkg/todo"
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

	raw, _ := ioutil.ReadFile("todo.txt")
	contents := string(raw)

	for _, task := range todo.ParseAll(contents) {
		fmt.Printf("%s\n", task)
	}
}

func printTask(task string) {
	log.Printf("Parsed task \"%s\" as \"%s\"", task, todo.ParseTask(task))
}

