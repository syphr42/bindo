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

package github

import (
	"flag"
	"fmt"

	"github.com/syphr42/bindo/pkg/cmd"
)

type GitHubCommand struct {
	cmd.AbstractCommand

	host  string
	owner string
	name  string
}

func NewGitHubCommand() *GitHubCommand {
	flags := flag.NewFlagSet("github", flag.ContinueOnError)

	cmd := &GitHubCommand{
		AbstractCommand: cmd.AbstractCommand{
			Flags: flags,
		},

		host:  "",
		owner: "",
		name:  "",
	}

	cmd.Flags.StringVar(&cmd.host, "host", "github.com", "GitHub instance hostname")
	cmd.Flags.StringVar(&cmd.owner, "owner", "", "owner of the repository")
	cmd.Flags.StringVar(&cmd.name, "name", "", "name of the repository")

	return cmd
}

func (github *GitHubCommand) Run() error {
	fmt.Println("host = ", github.host, "; owner = ", github.owner, "; name = ", github.name)
	return nil
}
