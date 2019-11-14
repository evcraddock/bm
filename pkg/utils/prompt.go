package utils

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

const requiredMessage = "%s is required"

// InputStringPrompt prompts for a string
func InputStringPrompt(label, defaultValue string, validate func(label, value, message string) error) string {
	input := bufio.NewScanner(os.Stdin)

	fmt.Printf("%s (%s): ", label, defaultValue)
	for input.Scan() {
		inputValue := input.Text()

		if inputValue == "" {
			inputValue = defaultValue
		}

		err := validate(label, inputValue, requiredMessage)
		if err == nil {
			return inputValue
		}

		fmt.Printf("%s\n", err.Error())
		fmt.Printf("%s: ", label)
	}

	return ""
}

// StringPrompt prompts for a string - value is required
func StringPrompt(label, value string) string {
	return InputStringPrompt(label, value, RequiredInputValidator)
}

// ListPrompt prompts for a comma delimited string returns a string array
func ListPrompt(label, value string) []string {
	csvlist := InputStringPrompt(label, value, NotRequiredInputValidator)
	reg, err := regexp.Compile("[^a-zA-Z0-9,]+")
	if err != nil {
		fmt.Printf("Invalid list")
		return ListPrompt(label, value)
	}

	scrubbed := reg.ReplaceAllString(csvlist, "")
	if scrubbed == "" {
		return nil
	}

	list := toList(scrubbed)

	return list
}

// RequiredInputValidator validates an input value to ensure it has a value
func RequiredInputValidator(label, input, message string) error {
	if len(strings.TrimSpace(input)) < 1 {
		return fmt.Errorf(message, label)
	}

	return nil
}

// NotRequiredInputValidator a validator which does not require a value
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
