package configbin

import (
	"time"

	"github.com/deadcheat/awsset"
)

// Assets a generated file system
var Assets = awsset.NewFS(
	map[string][]string{
		"/config": []string{
			"configfile.toml",
		},
	},
	map[string]*awsset.File{
		"/config":                 awsset.NewFile("/config", nil, 0x800001ed, time.Unix(1526447308, 1526447308936443679)),
		"/config/configfile.toml": awsset.NewFile("/config/configfile.toml", []byte(_Assets355e50d6eca05369f3315665c674fa974e2d230a), 0x1a4, time.Unix(1526447443, 1526447443457464453)),
	},
)

// binary data
var (
	_Assets355e50d6eca05369f3315665c674fa974e2d230a = "[server]\nHost = \"127.0.0.1\"\nPort = 3000\n"
)

