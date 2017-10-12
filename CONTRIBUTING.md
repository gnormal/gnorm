# Contributing

Of course, contributions are more than welcome. Please read these guidelines for making the process as painless as possible.

## Discussion

Development discussion should take place on the #gnorm channel of [gopher
slack](https://gophers.slack.com/).

There is a separate #gnorm-dev channel that has the github app to post github activity to the channel, to make it easy to follow.

## Tasks

See the [Nearterm Todos](https://github.com/gnormal/gnorm/projects/1) github
project on the repo to see what needs to be done soon.  If there's something
you'd like to do that's not on the list, feel free to file a new issue, or talk
to the devs on the #gnorm channel in [gopher slack](https://gophers.slack.com/).

## Dependency Management

gnorm uses the official dep tool for managing dependencies.

`go get -u github.com/golang/dep/cmd/dep`

If you add a dependency to the binary, make sure to update the vendor directory by running `dep ensure` and adding the resulting files to the repo.

## Setting up development environment

Imports points to `gnorm.org/gnorm` instead of `github.com/gnormal/gnorm` so a safer way to set this up is to follow these simple steps.

```shell
mkdir $GOPATH/src/gnorm.org

cd $GOPATH/src/gnorm.org

git clone git@github.com:gnormal/gnorm.git

cd gnorm 
```

After this, you need to make sure that the [dep](https://github.com/golang/dep) tool is installed, and `$GOPATH/bin` is in your system `PATH`. Once complete, proceed to run the following commands:

```shell
dep ensure
```

Next, make sure that [mage](https://github.com/magefile/mage) is installed. Mage will help you build the binary and ensure everything is working:

```
mage build
```

Running `gnorm` command should give you something like this, then you are all set up and good to `go`.

```
$ gnorm
A flexible code generator that turns your DB schema into
runnable code.  See full docs at https://gnorm.org

Usage:
  gnorm [command]

Available Commands:
  gen         Generate code from DB schema
  help        Help about any command
  init        Generates the files needed to run GNORM.
  preview     Preview the data that will be sent to your templates
  version     Displays the version of GNORM.

Flags:
  -h, --help   help for gnorm

Use "gnorm [command] --help" for more information about a command.
```

## Formatting

- All code *must* be go-formatted.
- All code *must* pass go-lint and go-vet.
- Comments *should* be wrapped at column 80.

## Documentation

If you change the behavior of gnorm, you must update the documentation.  The
docs are a hugo site under the ./site/ directory.  They are automatically
deployed when your code is committed, so that your updates and the docs are
always in sync.

## Testing

Tests must use the normal go testing package and call t.Parallel() at the start
of the test.  This makes it more likely that we catch race conditions during
tests, and don't rely on the order of tests running.

Tests must pass the race detector (run `go test -race ./...`).

## Git-fu

Please rebase your PR onto the current master before submitting:

`git rebase origin/master`

Please make PRs a single commit to start.  If changes are needed, add these as additional commits to the PR (they'll be squashed when the commit is accepted).

To squash a bunch of commits on a dev branch into one, do the following:

use `git log` and count the number of commits you want to squash.

run `git reset HEAD~N --soft` where N is __ONE LESS THAN__ the number of commits.

rerun `git log` and ensure your first commit is still in the log.  

If it is, add the rest of the changes to it with `git commit -a --amend` (optionally add `--no-edit` to keep the previous commit message, or `-m foo` to change the commit message).

If you accidentally reset all your changes, just commit as normal with `git
commit -a`.

## Commit and PR messages

Please be verbose when explaining changes in commit messages and PRs.  When writing commit messages, the first line is the title.  Make it up to ~100 characters, with a description that explains what the commit does. e.g. "adds a new foo to bar" or "converts X to use Y".

Put two line returns after the title and then add the body of the commit.  Explain both the high level of what has been done, and more importantly, *why* this change is being written.

Both of these apply to commits as well as PRs.

## Documentation

Documentation is of the utmost importance for users of gnorm. To ensure that the documentation is always up to date, many parts of the gnorm documentation are generated from the source directly.  The gnorm.org website code exists in the same repo under the /site directory, so that changes to the code are updated simultaneously with changes to the documentation.  Go generate is used liberally along with [gocog](https://github.com/natefinch/gocog) to generate a lot of the documentation.  After any change, before committing, run `go generate ./...` and note if any changes were made in the files under the /site directory.

