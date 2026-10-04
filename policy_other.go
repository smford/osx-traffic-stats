//go:build !darwin

package main

import "fmt"

func setAccessoryPolicy() {}

func configureFixedStatusItem() {}

func showAboutBox() {
	fmt.Printf("%s\n%s\n%s\n", appName, aboutMessage(), githubURL)
}
