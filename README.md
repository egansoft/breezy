# Breezy
Breezy is the simplest possible web nano-framework. No code needed.  

Sometimes all you need from a web server is just
* Retrieving errors from your logs with a `cat /logs/serverlog | grep ERROR`
* Restarting a job with a `supervisorctl restart mybigjob`
* Downloading some static files hosted at output/images

Breezy makes these tasks a breeze. Simply throw this into a text file (like routes.txt):

    /find/errors $ cat /logs/serverlog | grep ERROR
    /restart/big/job $ supervisorctl restart mybigjob
    /download : output/images

And run `breezy routes.txt`. You can specify arguments with brackets:

    /find/[needle]/limit/[n] $ cat /logs/serverlog | grep [needle] | head [n]
    /restart/[name]/job $ supervisorctl restart [name]

You can even pipe POST data into STDIN:

    /upload/[filename] $ cat > uploads/[filename]

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
