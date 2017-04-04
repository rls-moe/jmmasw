# Just Make Me A Static Website

JMMASW is a simple tool for a simple problem: Static Websites.

## Operation

1. Load all files with the extension ".html", ".js", ".css" or ".tmpl"
2. Parse every file into a Go Template
3. Load the data file ("./data.json" per default)
4. Run every template with the given data file as parameter
5. Write the output to files unless it's empty

## Usage

JMMASW accepts the following commandline parameters:

* `data` - Location of the data file
* `dir` - Project location, default is current directory
* `static-out` - Output directory, default is `./jasw-out` (For "Just a static Website Out")
* `ignore-no-data` - Defaults to true. If set to false then a missing data file results in termination of the program

## Templates

For Information on the Template format, refer to the godoc entries on `html/template`
and `text/template`.

## Roadmap

* Include Markdown Parser and Markdown File Loader
* Function to parse additional .json files from within the templates
* Method to ignore specific directories

## Wait what?

If you're wondering why this tool exists: I made it so I can generate my website for
multiple languages without having to copy-paste half the website all the time.

I considered Hugo, which is an excellent static page generator but it is also way to
complicated for this application. JMMASW works with "raw" HTML templates and makes
no assumption about your website.

You can take your website as it is now and provided it doesn't contain invalid
go-template code, it will come out the other end of JMMASW without change.