package utils

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

const REQUIRED_MESSAGE = "%s is required"

func InputStringPrompt(label string, validate func(label, value, message string) error) string {
	input := bufio.NewScanner(os.Stdin)

	fmt.Printf("%s: ", label)
	for input.Scan() {
		inputValue := input.Text()

		err := validate(label, inputValue, REQUIRED_MESSAGE)
		if err == nil {
			return inputValue
		}

		fmt.Printf("%s\n", err.Error())
		fmt.Printf("%s: ", label)
	}

	return ""
}

func StringPrompt(label string) string {
	return InputStringPrompt(label, RequiredInputValidator)
}

func ListPrompt(label string) []string {
	csvlist := InputStringPrompt(label, NotRequiredInputValidator)
	reg, err := regexp.Compile("[^a-zA-Z0-9,]+")
	if err != nil {
		fmt.Printf("Invalid list")
		return ListPrompt(label)
	}

	scrubbed := reg.ReplaceAllString(csvlist, "")
	if scrubbed == "" {
		return nil
	}

	list := toList(scrubbed)

	return list
}

func RequiredInputValidator(label, input, message string) error {
	if len(strings.TrimSpace(input)) < 1 {
		return fmt.Errorf(message, label)
	}

	return nil
}

func NotRequiredInputValidator(label, input, message string) error {
	return nil
}

func toList(csv string) []string {
	var list []string

	strlist := strings.Split(csv, ",")
	for _, n := range strlist {
		if ok := contains(list, n); !ok {
			list = append(list, n)
		}
	}

	return list
}

func contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}

	return false
}
