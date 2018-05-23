# awsset

[![Build Status](https://travis-ci.org/deadcheat/awsset.svg?branch=master)](https://travis-ci.org/deadcheat/awsset) [![Coverage Status](https://coveralls.io/repos/github/deadcheat/awsset/badge.svg?branch=master&service=github)](https://coveralls.io/github/deadcheat/awsset?branch=master&service=github) [![GoDoc](https://godoc.org/github.com/deadcheat/awsset?status.svg)](https://godoc.org/github.com/deadcheat/awsset)

library and cmd tools set for managment assets like go-bindata or go-assets

## install

To use asset builder, get all packages.
```
go get -u github.com/deadcheat/awsset/...
```

To only use generated file, get single package
```
go get -u github.com/deadcheat/awsset
```

## How to use

On command-line, awsset acts as asset builder like as go-assets-builder or go-bindata
```
> awsset -h 
NAME:
   awsset - make a binary contain some assets

USAGE:
   main [global options] command [command options] [arguments...]

VERSION:
   0.2.0

COMMANDS:
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --expression value, -e value  Regular expressions you want files to contain
   --out value, -o value         Output file name, result will be displaed to standard-out when it's skipped
   --package value, -p value     Package name for output (default: "main")
   --name value, -n value        Variable name for output assets (default: "Assets")
   --help, -h                    show help
   --version, -v                 print the version
```

## Examples
see example dir for full of codes.

### http static file

Generated asset is generated as implementation of http.FileSystem
```
	http.Handle("/", http.FileServer(assetsbin.Assets))
	log.Println("start server localhost:3000")
	if err := http.ListenAndServe(":3000", nil); err != nil {
		log.Fatal(err)
	}
```

Sometimes we changed root on http request from "/" such as "/statics/",
awsset.FileSystem has `WithPrefix` func.
```
	http.Handle("/static/", http.FileServer(assetsbin.Assets.WithPrefix("/static/")))
	log.Println("start server localhost:3000")
	if err := http.ListenAndServe(":3000", nil); err != nil {
		log.Fatal(err)
	}
```

### reading config file with config library like [github.com/spf13/viper](https://github.com/spf13/viper)
awsset.File has bytes.Reader, so you can use awsset.File directly
```
	viper.SetConfigType("toml")
	f, _ := configbin.Assets.File("/config/configfile.toml")
	viper.ReadConfig(f)
	var s Server
	_ = viper.UnmarshalKey("server", &s)
```
