# awsset
library and cmd tools set for managment assets like go-bindata or go-assets

## install

To use asset builder, get all packages.
```
go get github.com/deadcheat/awsset/...
```

To only use generated file, get single package
```
go get github.com/deadcheat/awsset
```

## How to use

On command-line, awsset acts as asset builder like as go-assets-builder or go-bindata
```
> awsset -h 
NAME:
   awsset - make a binary contain some assets

USAGE:
   awsset [global options] command [command options] [arguments...]

VERSION:
   0.1.0

COMMANDS:
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --except value, -e value   Regular expressions you want files to ignore
   --out value, -o value      Output file name, result will be displaed to standard-out when it's skipped
   --package value, -p value  Package name for output (default: "main")
   --var value, -t value      Variable name for output assets (default: "Assets")
   --help, -h                 show help
   --version, -v              print the version
```

