# Instafeed

`instafeed` is a very simple utility that allows you to generate a RSS feed of an [Instagram](https://www.instagram.com/) user.

## Requirements

* The Go language compiler (version >= 1.9)
* A valid Instagram account

## Installation

### Pre-compiled binaries

Pre-compiled binaries are available for [stable releases](https://github.com/falzm/instafeed/releases).

### Using `go get`

```console
go get github.com/falzm/instafeed
```

### From source

At the top of the sources directory, just type `make`. If everything went well, you should end up with binary named `instafeed` in your current directory.

## Usage

Run `instafeed -h` for usage help.

`instafeed` expects the `IG_LOGIN` and `IG_PASSWORD` environment variables set to your Instagram login and password respectively, and the username of the Instagram user provided as argument. On successful execution, it prints the resulting RSS feed on the standard output.

If multiple Instagram users arguments are provided, `instafeed` will bulk all retrieved posts into a single feed sorted in reverse chronological order (i.e. from newest to oldest).

Example:

```console
$ export IG_LOGIN="your_instagram_login" IG_PASSWORD="********"

# Fetching a single user feed
$ instafeed marutaro > marutaro.xml

# Fetching multiple user feeds, limiting to 10 posts per user
$ instafeed -n 10 marutaro nala_cat realgrumpycat > feeds.xml

# Fetching multiple user feeds from a list
$ cat > ~/.instagram_users <<EOF
marutaro
nala_cat
realgrumpycat
EOF
$ instafeed -n 10 -l ~/.instagram_users > feeds.xml
```

In order to avoid re-logging into Instagram upon each execution, `instafeed` stores the user profile into a local file that is read at runtime (your password is **not** stored). By default this file is located at `$HOME/.instafeed`, but a different path can be specified using the `-f` option.
