/*
Copyright 2022 Gregory P. Moyer.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"fmt"
	"os"

	"github.com/syphr42/bindo/internal/cmd"
	"github.com/syphr42/bindo/internal/cmd/github"
	"github.com/syphr42/bindo/internal/cmd/help"
)

var helpCommand = help.NewHelpCommand()

var commands = []cmd.Command{
	github.NewGitHubCommand(),
	helpCommand,
}

func root(args []string) error {
	if len(args) < 1 {
		return run(helpCommand, args)
	}

	subcommand := args[0]

	for _, cmd := range commands {
		if cmd.Name() == subcommand {
			return run(cmd, args[1:])
		}
	}

	return fmt.Errorf("unknown subcommand: %s", subcommand)
}

func main() {
	if err := root(os.Args[1:]); err != nil {
		fmt.Println(err)
		run(helpCommand, []string{})
		os.Exit(1)
	}
}

func run(cmd cmd.Command, args []string) error {
	cmd.Init(args)
	return cmd.Run()
}
