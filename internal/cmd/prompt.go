package cmd

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
)

func runConfirm(domain string) bool {
	confirm := false
	prompt := &survey.Confirm{
		Message: fmt.Sprintf("Do you want to delete domain %s?", domain),
	}

	if err := survey.AskOne(prompt, &confirm); err != nil {
		return false
	}

	return confirm
}
