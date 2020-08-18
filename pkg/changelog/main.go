package changelog

import (
	"fmt"
	"log"
	"strings"

	"semtag/pkg/terminal"
)

const (
	EnvVarGitCommitUrl    = "GIT_COMMIT_URL"
	EnvVarGitTagUrl       = "GIT_TAG_URL"
	ChangelogDefaultRegex = `^[0-9]+\.[0-9]+\.[0-9]+$`
)

func GenerateChangeLog(regex string) {
	commitUrl, err := terminal.GetEnv(EnvVarGitCommitUrl)
	if err != nil {
		log.Fatal(err)
	}
	tagUrl, err := terminal.GetEnv(EnvVarGitTagUrl)
	if err != nil {
		log.Fatal(err)
	}

	script := fmt.Sprintf("export REGEX=%q ; ", regex) + `cat <<"EOF" | bash > CHANGELOG.md
#!/usr/bin/env bash
set -uo pipefail

previous_tag=0
for current_tag in $(git tag --sort=-creatordate | grep -E "$REGEX" ); do

  if [ "$previous_tag" != 0 ]; then
    tag_date=$(git log -1 --pretty=format:'%ad' --date=iso8601 ${previous_tag})
    printf "## [${previous_tag}](https://bitbucket.org/projects/test/repos/my-project/tags/${previous_tag})\n${tag_date}\n\n"
    git log ${current_tag}...${previous_tag} --pretty=format:'*  %s ([%h](https://bitbucket.org/projects/test/repos/my-project/commits/%H))' --reverse | grep -v Merge
    printf "\n\n"
  fi

  previous_tag=${current_tag}
done
EOF`
	script = strings.Replace(script, "https://bitbucket.org/projects/test/repos/my-project/commits", commitUrl, 1)
	script = strings.Replace(script, "https://bitbucket.org/projects/test/repos/my-project/tags", tagUrl, 1)
	_, err = terminal.Shell(script)
	if err != nil {
		log.Panic(err)
	}
}
