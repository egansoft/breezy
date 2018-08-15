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

You can even pipe request data into STDIN:

    /upload/[filename] $ cat > uploads/[filename]

## Setup
    go install github.com/egansoft/breezy

## Usage
Using breezy is as simple as creating a text file with a route on every line. A route is either a shell command or a filesystem. A shell command route is of the form

    /my/url/path/[with]/my/[args] $ my shell command --[with] [args]

In this simple example, whenever someone hits your server on a URL like `/my/url/path/arg1/my/arg2`, the command `my shell command --arg1 arg2` is run. The HTTP message body is piped into STDIN, and STDOUT is outputed as the response. If your command has a non-zero exit code, no data is returned from the request and the HTTP status code is 500 (internal server error). The shell command doesn't have to be in bash: by specifying the shell argument when running breezy, you can write it in bash, zsh, python, ruby, or javascript, with /bin/sh as the default choice.

A filesystem is of the form

    /my/url/path : my/static/files

In this example, whenever someone hits your server on a URL like `/my/url/path/images/picture.png`, breezy looks for a file called `my/static/files/images/picture.png`, relative to where its being run. If no such file exists, breezy returns a 404 HTTP status code (page not found). 

If multiple routes match a given URL, then the one that is first in the file is used to serve requests on that URL. After saving this routes file, you can run breezy with `breezy [OPTIONS] FILE`, where `[OPTIONS]` are the following:

* -p, --port: specify the port to run on, with 8080 as the default
* -s, --shell: specify the shell to use, with /bin/sh as the default
* -d, --debug: enable debug mode, returning information about how the request was handled for every request
* -h, --help: display the help message

## Security
Breezy is a simple program, making it really easy to make insecure webservers. Please don't write routes that allow people to execute arbritrary commands on your server, like `/[command] $ [command]`. 
