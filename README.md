# Breezy
The simplest possible web nano-framework. No code needed.

## Usage
    Usage: breezy [OPTIONS] PORT FILE
    Run a Breezy server at the given PORT using the routes defined in FILE.

    FILE consists of a newline seperated list of routes, which have the forms:
      shell command:
        /url/[arg1]/with/[arg2]/args $ cmd_name --option arg1 arg2
        which runs the command with the arguments specified in the url, with the
        request body piped in to stdin
      filesystem root:
        /url/path : relative/filesystem/path
        which serves static files from the system path on that url path, as if the
        filesystem were mounted on the url path

    If multiple routes match a given url, the first one is used.

    Options:
      -p, --port     specify the port to run on, with 8080 as the default
      -s, --shell    specify a shell to use, with sh as the default
      -d, --debug    enable debug mode
      -h, --help     display this message and exit
