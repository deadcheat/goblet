package awsset

import (
	"os"
	"reflect"
	"testing"
	"time"
)

func TestNewFS(t *testing.T) {
	f := make([]*File, 3)
	f[0] = &File{
		Path:       "/tmp/test",
		Data:       nil,
		FileMode:   os.ModeDir,
		ModifiedAt: time.Now(),
	}
	f[1] = &File{
		Path:       "/tmp/test/hoge.txt",
		Data:       []byte("hogehoge"),
		FileMode:   0x800001ed,
		ModifiedAt: time.Now(),
	}
	f[2] = &File{
		Path:       "/tmp/test/fuga.txt",
		Data:       []byte("fuga"),
		FileMode:   0x800001ed,
		ModifiedAt: time.Now(),
	}
	files := map[string]*File{
		"/tmp/test":          f[0],
		"/tmp/test/hoge.txt": f[1],
		"/tmp/test/fuga.txt": f[2],
	}
	dirs := map[string][]string{
		"/tmp/test": []string{
			"hoge.png", "fuga.xml",
		},
	}

	// expected data
	expected := &FileSystem{
		Dirs:       dirs,
		Files:      files,
		pathPrefix: "",
	}

	// get actual
	actual := NewFS(dirs, files)

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("expected and actual are not equal, expected %#+v and actual %#+v", expected, actual)
	}
}
