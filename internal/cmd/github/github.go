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
	"context"
	"flag"
	"fmt"
	"regexp"

	"github.com/google/go-github/v47/github"

	"github.com/syphr42/bindo/internal/cmd"
)

type GitHubCommand struct {
	cmd.AbstractCommand

	host         string
	owner        string
	name         string
	preRelease   bool
	assetPattern string
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
	cmd.Flags.BoolVar(&cmd.preRelease, "prerelease", false, "include pre-releases")
	cmd.Flags.StringVar(&cmd.assetPattern, "pattern", ".*", "pattern to use for identifing an asset to download")

	return cmd
}

func (cmd *GitHubCommand) Run() error {
	releases, err := getReleases(cmd)

	if err != nil {
		return err
	}

	for _, release := range releases {
		if cmd.preRelease || !is(release.Prerelease) {
			handleRelease(cmd, release)
			break
		}
	}

	return nil
}

func getReleases(cmd *GitHubCommand) ([]*github.RepositoryRelease, error) {
	client := github.NewClient(nil)

	ctx := context.Background()
	releases, _, err := client.Repositories.ListReleases(ctx, cmd.owner, cmd.name, nil)
	if err != nil {
		return nil, err
	}

	return releases, nil
}

func handleRelease(cmd *GitHubCommand, release *github.RepositoryRelease) {
	fmt.Println("Release found: ", *release.Name)

	asset := findAsset(cmd, release)
	if asset == nil {
		fmt.Println("No matching asset found.")
	} else {
		fmt.Println("Asset found: ", *asset.Name)
	}
}

func findAsset(cmd *GitHubCommand, release *github.RepositoryRelease) *github.ReleaseAsset {
	r, _ := regexp.Compile(cmd.assetPattern)

	for _, asset := range release.Assets {
		if asset.Name != nil && r.MatchString(*asset.Name) {
			return asset
		}
	}

	return nil
}

func is(value *bool) bool {
	if value == nil {
		return false
	}

	return *value
}
