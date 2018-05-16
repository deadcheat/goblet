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
		"/config/configfile.toml": awsset.NewFile("/config/configfile.toml", _Assets355e50d6eca05369f3315665c674fa974e2d230a, 0x1a4, time.Unix(1526447443, 1526447443457464453)),
	},
)

// binary data
var (
	_Assets355e50d6eca05369f3315665c674fa974e2d230a = []byte{0x5b, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x5d, 0xa, 0x48, 0x6f, 0x73, 0x74, 0x20, 0x3d, 0x20, 0x22, 0x31, 0x32, 0x37, 0x2e, 0x30, 0x2e, 0x30, 0x2e, 0x31, 0x22, 0xa, 0x50, 0x6f, 0x72, 0x74, 0x20, 0x3d, 0x20, 0x33, 0x30, 0x30, 0x30, 0xa}
)

