package changelog

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"semtag/pkg/output"
	"semtag/pkg/terminal"
)

const (
	DefaultChangelogFile = "CHANGELOG.md"
)

type file struct {
	name string
}

func (f *file) String(prefix, suffix string) string {
	var out string
	if prefix != "" {
		out = prefix + out
	}

	if f.name == "" {
		f.name = DefaultChangelogFile
	}
	out += f.name

	if suffix != "" {
		out += suffix
	}

	out += ".md"
	return out
}

// Generate the contents of the changelog and write them to a file
func (f *file) Generate(fileName, regex, commitUrl, tagUrl string) error {
	script := fmt.Sprintf(`
# write the changelog to the specified file name
cat <<"EOF" | bash > %q
#!/usr/bin/env bash
set -euo pipefail

previous_tag=0

# for each of the tags that match the provided regex create a changelog entry
for current_tag in $(git tag --sort=-v:refname | grep -E %s ); do

  if [ "$previous_tag" != 0 ]; then
	# get the commit datetime
    tag_date=$(TZ=UTC git log -1 --pretty=format:'%%ad' --date=iso-local ${previous_tag})

	# set the changelog entry header
	# use the tag url to generate the hyperlinks for tags
    printf "## [${previous_tag}](%s/${previous_tag})\n${tag_date}\n\n"

	# set the changelog entry body
	# use the commit url to generate the hyperlinks for commits
	# ignore merge commits
    TZ=UTC git log ${current_tag}...${previous_tag} --pretty=format:'*  %%s by [%%aN](mailto:%%aE) ([%%h](%s/%%H))' --reverse | grep -v Merge

    printf "\n\n"
  fi

  previous_tag=${current_tag}
done
EOF`,
		fileName, regex, tagUrl, commitUrl)
	_, err := terminal.Shell(script)
	if err != nil {
		return err
	}

	output.Logger().WithFields(logrus.Fields{
		"changelogFile":        fileName,
		"changelogGitTagRegex": regex,
		"changelogBaseUrls":    []string{tagUrl, commitUrl},
	}).Info("changelog generated")
	return nil
}
