package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"strings"
)

var (
	unchecked = "\u2610" // ☐
	checked   = "\u2611" // ☑
)

func write() {

	clear()

	file, err := os.Create("output.csv")

	if err != nil {
		panic(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	header := []string{"Status", "Task"}

	if err := writer.Write(header); err != nil {
		panic(err)
	}

	records := [][]string{}

	scanner := bufio.NewScanner(os.Stdin)
	for {

		fmt.Print("Enter your task (or 'quit' to stop): ")

		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				fmt.Println("input error:", err)
			}
			break // EOF
		}

		task := strings.TrimSpace(scanner.Text())

		if strings.EqualFold(task, "quit") {
			for _, r := range records {
				if err := writer.Write(r); err != nil {
					panic(err)
				}
			}
			fmt.Println("Tasks saved.")
			return
		}

		if task == "" {
			fmt.Println("No task entered!")
			continue
		}

		status := unchecked
		records = append(records, []string{status, task})
		fmt.Println("Added:", status, task)

	}
}
func flipStatusByNumber() {

	clear()

	file, err := os.Open("output.csv")
	if err != nil {
		panic(err)
	}
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	file.Close()
	if err != nil {
		panic(err)
	}

	if len(records) <= 1 {
		fmt.Println("No tasks found.")
		return
	}

	fmt.Println("\nYour TODO List:")
	for i, record := range records {
		if i == 0 {
			continue // skip header
		}
		fmt.Printf("%d. %s %s\n", i, record[0], record[1])
	}

	fmt.Print("\nEnter the number of the task to flip: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	input := strings.TrimSpace(scanner.Text())

	var index int
	_, err = fmt.Sscanf(input, "%d", &index)
	if err != nil || index <= 0 || index >= len(records) {
		fmt.Println("Invalid number.")
		return
	}

	if records[index][0] == unchecked {
		records[index][0] = checked
	} else {
		records[index][0] = unchecked
	}

	file, err = os.Create("output.csv")
	if err != nil {
		panic(err)
	}
	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, record := range records {
		if err := writer.Write(record); err != nil {
			panic(err)
		}
	}

}

func read() {
	clear()

	file, err := os.Open("output.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	/*
		headers, err := reader.Read()
		if err != nil {
			panic(err)
		}

		fmt.Println("Headers: ", headers)
	*/

	for {

		record, err := reader.Read()
		if err != nil {
			break
		}

		fmt.Println(record)
	}

	fmt.Print("\nPress Enter to return to the menu...")
	bufio.NewScanner(os.Stdin).Scan()
}

func clear() {
	fmt.Print("\033[H\033[2J")
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		clear()
		fmt.Print("TODO Time! 'write', 'read', 'flip', or 'done': ")

		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				fmt.Println("input error:", err)
			}
			return
		}

		resp := strings.ToLower(strings.TrimSpace(scanner.Text()))

		switch resp {
		case "write":
			write()
		case "read":
			read()
		case "flip":
			flipStatusByNumber()
		case "done", "quit", "exit", "":
			fmt.Println("Bye!")
			return
		default:
			fmt.Println("Please type 'write', 'read', 'flip', or 'quit'.")
		}
	}
}