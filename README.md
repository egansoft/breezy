# Breezy
The simplest possible web nano-framework. No code needed.

## Usage
Usage: breezy [OPTIONS] FILE
Run a Breezy server using the routes defined in FILE.

FILE consists of a newline seperated list of routes, which have the forms:
  shell command:
      /my/url/path $ command to run
    Args can be specified with brackets:
      /my/url/[arg1]/path $ command to [arg1] run
    The request body is piped into stdin.
  filesystem root:
      /url/path : relative/filesystem/path
    Static files from the filesystem path are served on the url path, as if the
    filesystem were mounted on the url path.

If multiple routes match a given url, the first one is used.

Options:
  -p, --port     specify the port to run on, with 8080 as the default
  -s, --shell    specify a shell to use, with sh as the default
  -d, --debug    enable debug mode
  -h, --help     display this message and exit
