# csp

Uses various Computational Intelligence algorithms to produce solutions for the Cutting Stock Problem.

## Usage

You'll need to have `go` installed. The following assumes it's all set up and you have `$GOPATH/bin` in your `$PATH`. Then you'll want to put this project directory at `$GOPATH/src/github.com/sprusr/csp`.

Now you can install it and start using it.

```bash
# install the csp command
cd csp
go install .

# show a list of commands
csp --help

# compare AIS, ACO and random search on instance 2, after 10 iterations, over 40 runs
csp compare -I 2 -i 10 -r 40
```

To the marker: sorry if it being in `go` has caused you undue stress!

## Architecture

The `csp` command uses the `cobra` library to handle command line flags and sub-commands. Each command is found in `csp/cmd`, with `root.go` handling calls without any sub-command. All actual logic can be found in the three files in the `csp` directory.
