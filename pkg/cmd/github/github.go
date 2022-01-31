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
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/syphr42/bindo/pkg/cmd"
)

type GitHubCommand struct {
	cmd.AbstractCommand

	host       string
	owner      string
	name       string
	preRelease bool
}

type release struct {
	Name       string `json:"name"`
	Tag        string `json:"tag_name"`
	PreRelease bool   `json:"prerelease"`
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

	return cmd
}

func (github *GitHubCommand) Run() error {
	releases, err := getReleases(github)

	if err != nil {
		return err
	}

	for _, release := range releases {
		if github.preRelease || !release.PreRelease {
			handleRelease(release)
			break
		}
	}

	return nil
}

func handleRelease(release release) {
	fmt.Println("name =", release.Name, "tag =", release.Tag, "pre-release =", release.PreRelease)
}

func buildUrl(github *GitHubCommand) string {
	return "https://api." + github.host + "/repos/" + github.owner + "/" + github.name + "/releases"
}

func getReleases(github *GitHubCommand) ([]release, error) {
	client := http.Client{
		Timeout: time.Second * 2,
	}

	req, err := http.NewRequest(http.MethodGet, buildUrl(github), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "bindo")

	res, getErr := client.Do(req)
	if getErr != nil {
		return nil, getErr
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		return nil, readErr
	}

	return parseReleases(body)
}

func parseReleases(jsonData []byte) ([]release, error) {
	var releases []release

	err := json.Unmarshal(jsonData, &releases)
	if err != nil {
		return nil, err
	}

	return releases, nil
}
