# Abbey CLI

An interactive way to set up Abbey. Built on Go with [Cobra](https://github.com/spf13/cobra) & [BubbleTea](https://github.com/charmbracelet/bubbletea).

The Abbey CLI will 
* Pull existing examples
* Guide you on how to configure the example - such as setting the time expiry for your access policies
* Commit and deploy your changes to GitHub

## Prerequisites
Before running the CLI, you must have
* An Abbey account
  * You can create one at [app.abbey.io](https://app.abbey.io/)
* A GitHub account

## Installing
You can install the CLI via
1. Cloning this repo
2. Building it via `go build -o abbey main.go`, or if you'd like to make it accessible everywhere `go build -o $GOPATH/bin/abbey`

## Interactive Mode
You can run the CLI in interactive mode simply by running `./abbey init` and it will guide you the rest of the way.

## Flags
There are a number of flags available as well to speed things up. You can view the full list by running `./abbey init --help`

```bash
./abbey init --help
Initializes an Abbey example


Usage:

  abbey init [flags]


Flags:

  -a, --accessOutput string   Location of Access Output
  -h, --help                  help for init
      --repo string           Git repo name in the format github-username/repo-name
  -r, --reviewer string       Abbey email address of the reviewer
  -t, --timeExpiry string     Time expiry of permissions

```

![abbey_cli_short](https://github.com/abbeylabs/cli/assets/8519403/a4df3fc9-a188-4665-959f-18e78bb507da)

## Need help?
If you have any questions along the way, feel free to get in touch - we'd love to chat with you
* Join our [Slack](https://join.slack.com/t/abbey-io/shared_invite/zt-255094l05-wtoSqqAuGJtI5BnpNZl73Q) 
* Email us at [hello@abbey.io](mailto:hello@abbey.io)
* Post a [Github Issue](https://github.com/abbeylabs/cli/issues/new)

