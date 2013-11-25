# `gommap`: Golang tools for working with MetaMap

`gommap` is a set of tools for interacting with the National Library of Medicine's [MetaMap](http://metamap.nlm.nih.gov) concept identification and extraction tool. While a lovely piece of software, MetaMap as shipped operates in two different modes, neither of which work for all situations:

1. Batch-processing, e.g. of a directory full of files
2. An "interactive" mode via a REPL-style shell.

While both have their uses, I wanted to use MetaMap as a component in a larger application in a more dynamic way. Really, what I wanted was a REST-ish API for MetaMap; the NLM does provide a web service, but it's not really intended for one-off and realtime queries. So I did what any self-respecting hacker would do- I rolled my own! 

Go happens to be a lovely language for this sort of thing. The `mm_server` executable manages a pool of MetaMap instances as well as I/O channels to each, and multiplexes HTTP requests to an instance. It also includes a trivial example of a web-based UI.

## TODO

- `gommap` relies on a couple of hard-coded paths; that obviously needs to change
- More configurability
- More documentation
- Currently, the `MMOs` struct in `mmo_struct.go` doesn't capture all of the data in MetaMap's output. 
- Refactoring- `mm_server` should only handle network-related stuff, the actual MetaMap I/O should be handled elsewhere.