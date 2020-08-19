# Semantic Tagger ![GitHub Actions](https://img.shields.io/github/workflow/status/mpdred/semantic-tagger/Pipeline/master) ![tag](https://img.shields.io/github/v/tag/mpdred/semantic-tagger) ![last commit](https://img.shields.io/github/last-commit/mpdred/semantic-tagger)

Increment a version number as per [Semantic Versioning 2.0.0 specifications](https://semver.org/)


If you don't provide the current version, Git tags will be checked to see if they contain a version string as defined at https://semver.org/ 

To determine the type of change, the [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/) rules are used (you can override this with a cli argument):
> The commit message should be structured as follows:
>
>  <type>[optional scope]: <description>
>
>  [optional body]
>
>  [optional footer(s)]
>
> The commit contains the following structural elements, to communicate intent to the consumers of your library:
> - fix: a commit of the type fix patches a bug in your codebase (this correlates with PATCH in semantic versioning).
> - feat: a commit of the type feat introduces a new feature to the codebase (this correlates with MINOR in semantic versioning).
> - BREAKING CHANGE: a commit that has a footer BREAKING CHANGE:, or appends a ! after the type/scope, introduces a breaking API change (correlating with MAJOR in semantic versioning). A BREAKING CHANGE can be part of commits of any type.
> - types other than fix: and feat: are allowed, for example @commitlint/config-conventional (based on the the Angular convention) recommends build:, chore:, ci:, docs:, style:, refactor:, perf:, test:, and others.
> - footers other than BREAKING CHANGE: <description> may be provided and follow a convention similar to git trailer format.

A `BREAKING CHANGE` increments the major number, and resets the feature and patch number to zero (e.g. 4.0.7 -> 5.0.0)

A `feat`ure increments the minor number, and resets the patch number to zero (e.g. 4.0.7 -> 4.1.0)

All other types increment the patch number (e.g. 4.0.7 -> 4.0.8)



## Docs
- [how to test/build](docs/build.md) the project
- what [command line arguments](docs/usage.md) are available and  how to use the compiled binary for creating Git tags or update version numbers in files
- see the shell script for [Git configuration](docs/git.sh) for various hack configurations when running _Semantic Tagger_ in a CI executor environment (e.g. GitLab, Bitbucket, etc.)