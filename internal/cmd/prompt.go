package cmd

import (
	"errors"
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/edsonmichaque/dnsimple-cli/internal/config"
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

func promptConfig(c *config.Config) (*config.Config, string, error) {
	baseURL := prodBaseURL

	accountID, err := promptAccountID(c.Account)
	if err != nil {
		return nil, "", err
	}

	accessToken, err := promptAccessToken(c.AccessToken)
	if err != nil {
		return nil, "", err
	}

	env, err := promptEnvironment(envProd)
	if err != nil {
		return nil, "", err
	}

	if env == envSandbox {
		baseURL = sandboxBaseURL
	}

	if env == envDev {
		baseURL, err = promptBaseURL(prodBaseURL)
		if err != nil {
			return nil, "", err
		}
	}

	fileFormat, err := promptFileFormat(configFormatJSON)
	if err != nil {
		return nil, "", err
	}

	confirmation, err := promptConfirmation("Do you want to save?", true)
	if err != nil {
		return nil, "", err
	}

	if !confirmation {
		return nil, "", errors.New("did not confirm")
	}

	cfg := config.Config{
		Account:     accountID,
		AccessToken: accessToken,
	}

	if env == envDev {
		cfg.BaseURL = baseURL
	}

	if env == envSandbox {
		cfg.Sandbox = true
	}

	return &cfg, fileFormat, nil
}

func promptAccessToken(value string) (string, error) {
	prompt := &survey.Input{
		Message: "Access Token",
		Default: value,
	}

	var token string
	if err := survey.AskOne(prompt, &token); err != nil {
		return "", err
	}

	return token, nil
}

func promptAccountID(value string) (string, error) {
	prompt := &survey.Input{
		Message: "Account ID",
		Default: value,
	}

	var accountID string
	if err := survey.AskOne(prompt, &accountID); err != nil {
		return "", err
	}

	return accountID, nil
}

func promptEnvironment(value string) (string, error) {
	prompt := &survey.Select{
		Message: "Environment",
		Options: []string{envProd, envSandbox, envDev},
		Default: value,
	}

	var env string
	if err := survey.AskOne(prompt, &env); err != nil {
		return "", err
	}

	return env, nil
}

func promptBaseURL(value string) (string, error) {
	prompt := &survey.Input{
		Message: "Base URL",
		Default: value,
	}

	var baseURL string
	if err := survey.AskOne(prompt, &baseURL); err != nil {
		return "", err
	}

	return baseURL, nil
}

func promptFileFormat(value string) (string, error) {
	prompt := &survey.Select{
		Message: "File format",
		Options: []string{configFormatJSON, configFormatYAML, configFormatTOML},
		Default: value,
	}

	var fileFormat string
	if err := survey.AskOne(prompt, &fileFormat); err != nil {
		return "", err
	}

	return fileFormat, nil
}

func promptConfirmation(msg string, value bool) (bool, error) {
	var confirmation bool

	prompt := &survey.Confirm{
		Message: msg,
		Default: value,
	}

	if err := survey.AskOne(prompt, &confirmation); err != nil {
		return false, err
	}

	return confirmation, nil
}
