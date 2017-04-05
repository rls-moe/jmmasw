# Just Make Me A Static Website

JMMASW is a simple tool for a simple problem: Static Websites.

## Operation

1. Load all files with the extension ".html" or ".tmpl"
2. Parse every file into a Go Template
3. Load the data file ("./data.json" per default)
4. Run all templates with a name that ends in ".html"
5. Write the output to files unless it's empty

## Usage

JMMASW accepts the following commandline parameters:

* `data` - Location of the data file
* `dir` - Project location, default is current directory
* `static-out` - Output directory, default is `./jasw-out` (For "Just a static Website Out")
* `ignore-no-data` - Defaults to true. If set to false then a missing data file results in termination of the program

## Building

JMMASW is built using Go1.8, so you should at minimum install Go1.8 for this to work but older
versions might work too (but no promises)

To install simply run the following commands:

```
go get go.rls.moe/jmmasw
go install go.rls.moe/jmmasw
```

## Templates

For Information on the Template format, refer to the godoc entries on `html/template`
and `text/template`.

Templates are defined by their filename but if you manually defined templates this does not apply.

Only templates with a name that ends in ".html" will be rendered, ".tmpl" files will only be parsed
and are intended for common functionality.

### Template Functions

* `file` accepts a single string as parameter and attempts to read the specified file form the file system
* `json` parses a given string into json format, the top-level structure must be a map, simple arrays are not permitted
* `markdown` renders a given string into HTML using Markdown Processors without any sanitization (don't use on user inputs!)
* `dict` accepts any number of key-value pairs with a string-typed key and returns the result. This allows to combine several variables into a single pipeline

## Roadmap

* ~~Include Markdown Parser and Markdown File Loader~~
* ~~Function to parse additional .json files from within the templates~~
* Method to ignore specific directories
* Function to output files that aren't present as templates
    * Example: index.html generates index-en.html and index-de.html instead of index.html

## Wait what?

If you're wondering why this tool exists: I made it so I can generate my website for
multiple languages without having to copy-paste half the website all the time.

I considered Hugo, which is an excellent static page generator but it is also way too
complicated for this application. JMMASW works with "raw" HTML templates and makes
no assumption about your website.

You can take your website as it is now and provided it doesn't contain invalid
go-template code, it will come out the other end of JMMASW without change.
