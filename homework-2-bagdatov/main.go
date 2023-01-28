package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"v0/office"
)

const (
	errMsg     = "\n\033[31mERROR\033[0m"
	successMsg = "\033[32mSUCCESS\033[0m\n"
)

func main() {
	file, err := os.Open("employeeActions")
	if err != nil {
		log.Println(err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var content [][]string

	for scanner.Scan() {
		text := scanner.Text()
		if !strings.ContainsAny(text, ",") {
			log.Println("incorrect formatting")
			return
		}

		line := strings.Split(text, ",")
		content = append(content, line)
	}

	if err := scanner.Err(); err != nil {
		log.Println(err)
	}

	for i := range content {
		processEmployee(content[i])
		fmt.Println()
	}

}

func processEmployee(content []string) {
	var employee office.Employee
	var err error

	for i := range content {
		if i == 0 {
			// first element in slice must be employee type
			employee, err = office.NewEmployeeFactory(content[i])
			if err != nil {
				fmt.Println(err, errMsg)
				return
			}
			fmt.Printf("new employee: %T\n", employee)
			continue
		}

		if i == 1 && content[1] != "office" {
			// second element in slice must be office
			fmt.Printf("first location of %T is not an office %s", employee, errMsg)
			return
		}

		if err := moveEmployee(content[i], employee); err != nil {
			fmt.Println(err, errMsg)
			return
		}

		fmt.Printf("%T entered %s\n", employee, content[i])
	}
	fmt.Printf(successMsg)
}

func moveEmployee(newLocation string, employee office.Employee) error {
	location, err := office.NewLocationFactory(newLocation)
	if err != nil {
		return fmt.Errorf("%w: %s", err, newLocation)
	}

	err = employee.MoveToLocation(location)
	if err != nil {
		return err
	}
	return nil
}
