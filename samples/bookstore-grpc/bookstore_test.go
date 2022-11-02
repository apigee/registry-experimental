package main

import (
	"os/exec"
	"strings"
	"testing"
)

type step struct {
	command string
	fails   bool
	expect  string
}

var steps = []step{
	step{command: "bookstore bookstore reset --json"},
	step{command: "bookstore bookstore list-shelves --json"},
	step{command: "bookstore bookstore create-shelf --json --shelf.id 1 --shelf.theme foo"},
	step{command: "bookstore bookstore create-shelf --json --shelf.id 1 --shelf.theme foo", fails: true},
	step{command: "bookstore bookstore create-shelf --json --shelf.id 2 --shelf.theme foo"},
	step{command: "bookstore bookstore create-shelf --json --shelf.id 3 --shelf.theme foo"},
	step{command: "bookstore bookstore create-shelf --json --shelf.id 4 --shelf.theme foo"},
	step{command: "bookstore bookstore delete-shelf --json --shelf 4"},
	step{command: "bookstore bookstore delete-shelf --json --shelf 4", fails: true},
	step{command: "bookstore bookstore list-shelves --json"},
	step{command: "bookstore bookstore get-shelf --json --shelf 3"},
	step{command: "bookstore bookstore create-book --json --shelf 3 --book.id 1 --book.title foo --book.author bar"},
	step{command: "bookstore bookstore create-book --json --shelf 3 --book.id 1 --book.title foo --book.author bar", fails: true},
	step{command: "bookstore bookstore create-book --json --shelf 3 --book.id 2 --book.title foo --book.author bar"},
	step{command: "bookstore bookstore create-book --json --shelf 3 --book.id 3 --book.title foo --book.author bar"},
	step{command: "bookstore bookstore create-book --json --shelf 3 --book.id 4 --book.title foo --book.author bar"},
	step{command: "bookstore bookstore delete-book --json --shelf 3 --book 4"},
	step{command: "bookstore bookstore delete-book --json --shelf 3 --book 4", fails: true},
	step{command: "bookstore bookstore list-books --json --shelf 3"},
	step{command: "bookstore bookstore get-book --json --shelf 3 --book 1",
		expect: `{
  "id": "1",
  "author": "bar",
  "title": "foo"
}
`},
}

func TestBookstore(t *testing.T) {
	for _, s := range steps {
		parts := strings.Split(s.command, " ")
		out, err := exec.Command(parts[0], parts[1:len(parts)]...).Output()
		if err != nil {
			if !s.fails {
				t.Errorf("failed '%s': %s", s.command, err)
			}
		} else {
			if s.fails {
				t.Errorf("'%s' was expected to fail", s.command)
			}
			if s.expect != "" {
				if string(out) != s.expect {
					t.Errorf("expected '%s' got '%s'", s.expect, string(out))
				}
			}
		}
	}
}
