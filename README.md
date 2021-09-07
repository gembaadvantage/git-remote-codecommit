# git-remote-codecommit

A [git-remote-helper](https://git-scm.com/docs/gitremote-helpers) that supports basic push and pull functionality when working with CodeCommit repositories using the AWS `codecommit` protocol. When installed the helper acts as a transparent proxy, converting the `codecommit` protocol into AWS V4 authenticated HTTPS requests. Removing the need for dedicated AWS HTTPS credentials.

## Install

Binary downloads can be found on the [Releases](https://github.com/gembaadvantage/git-remote-codecommit/releases) page. Unpack the `git-remote-codecommit` binary and add it to your PATH.

### Homebrew

To use [Homebrew](https://brew.sh/):

```sh
brew tap gembaadvantage/tap
brew install gembaadvantage/tap/git-remote-codecommit
```

### Scoop

To use [Scoop](https://scoop.sh/):

```sh
scoop install git-remote-codecommit
```

### Script

To install using a shell script:

```sh
curl https://raw.githubusercontent.com/gembaadvantage/git-remote-codecommit/main/scripts/install | sh
```

## Quick Start

Clone the repository using your standard git syntax, but provide the clone URL using the `codecommit` protocol format of:

```text
codecommit::<REGION>://<PROFILE>@<REPOSITORY>
```

```sh
$ git clone codecommit::eu-west-1://repository

Cloning into 'repository'...
remote: Counting objects: 167, done.
Receiving objects: 100% (167/167), 96.07 KiB | 634.00 KiB/s, done.
Resolving deltas: 100% (31/31), done.
```

Both `git pull` and `git push` operations will behave as normal.

### AWS Named Profile

Depending on your chosen authentication mechansim, you may need to provide an AWS named profile to authenticate with CodeCommit. To do this, provide the named profile, suffixed with an `@`, before the repository name.

```sh
git clone codecommit::eu-west-1://profile@repository
```
