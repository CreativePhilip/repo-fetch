package utils

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func StdinInput(message string) (string, error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("%s: ", message)
	value, err := reader.ReadString('\n')

	if err != nil {
		return "", err
	}

	return strings.Replace(value, "\n", "", -1), nil
}

func StdinSelectFrom(message string, options []string) (string, error) {
	for idx, option := range options {
		fmt.Printf("%d. %s\n", idx, option)
	}

	fmt.Printf("%s: ", message)
	reader := bufio.NewReader(os.Stdin)
	for {
		selectedIndexStr, err := reader.ReadString('\n')

		if err != nil {
			return "", err
		}

		selectedIndex, err := strconv.Atoi(strings.Replace(selectedIndexStr, "\n", "", -1))

		if err != nil {
			return "", err
		}

		if !(selectedIndex >= 0 && selectedIndex <= len(options)) {
			fmt.Println("Please select a valid option")
			continue
		}

		return options[selectedIndex], nil
	}

}
