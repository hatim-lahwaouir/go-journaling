package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type Thoughts struct {
	Thoughts []Thought `json:"thoughts"`
}

type Thought struct {
	Content string    `json:"content"`
	Date    time.Time `json:"date"`
}

func help() {

	fmt.Println("this is the only available cmd")
	fmt.Println("write -> for writing a thoughts")
	fmt.Println("delete -> for deleting a thoughts")
	fmt.Println("edit -> for edeting a thoughts")
	fmt.Println("show -> for showing your thoughts")

	os.Exit(0)
}

var cmds = []string{"write", "show"}
var journalFileName string = "jornal.json"

func PrintError(message string, err string) {

	fmt.Println(message)
	fmt.Println("here is the error", err)
	os.Exit(1)
}

func createIfNotExists() *os.File {
	home := os.Getenv("HOME")
	journalfilePath := filepath.Join(home, journalFileName)

	if _, err := os.Stat(journalfilePath); errors.Is(err, os.ErrNotExist) {
		f, err := os.Create(journalfilePath)
		if err != nil {
			PrintError("we can't create the journal file", err.Error())
		}

		data, _ := json.Marshal(&Thoughts{Thoughts: []Thought{}})
		f.Write(data)
		return f
	}

	f, err := os.Open(journalfilePath)
	if err != nil {
		PrintError("we can't open the journal file", err.Error())
	}
	return f
}

func write() {
	fd := createIfNotExists()
	defer fd.Close()
	bytes, err := os.ReadFile(fd.Name())
	if err != nil {
		PrintError("We can't read the journal file", err.Error())
	}

	var data Thoughts
	json.Unmarshal(bytes, &data)
	fmt.Println(data)

	// now let's create a thought and append it to our json file

	var thought Thought
	fmt.Println("input text:")
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		PrintError("we can't read from stdin ", err.Error())
	}
	thought.Content = line
	thought.Date = time.Now()

	data.Thoughts = append(data.Thoughts, thought)
	bytes, _ = json.Marshal(data)
	os.WriteFile(fd.Name(), bytes, 0666)
}

func show() {
	fd := createIfNotExists()

	defer fd.Close()

	bytes, err := os.ReadFile(fd.Name())
	if err != nil {
		PrintError("We can't read the journal file", err.Error())
	}
	var data Thoughts

	json.Unmarshal(bytes, &data)

	fmt.Println("=> Thoughts:")
	for _, thought := range data.Thoughts {
		fmt.Println("*********************")
		fmt.Println(thought.Date.Format("Jan 02, 2006 3:04 PM"))

		fmt.Println(thought.Content)
		fmt.Println("*********************")
	}
}

func main() {

	// checking args
	args := os.Args
	if len(args) < 2 {
		help()
	}
	cmd := args[1]
	var cmdIndex int = -1
	for i, str := range cmds {
		if str == cmd {
			cmdIndex = i
		}
	}

	if cmdIndex == -1 {
		help()
	}

	switch cmds[cmdIndex] {
	case "write":
		write()
	case "show":
		show()
	}
}
