package awsset

import (
	"os"
	"reflect"
	"testing"
	"time"
)

var (
	f = []*File{
		&File{
			Path:       "/tmp/test",
			Data:       nil,
			FileMode:   os.ModeDir,
			ModifiedAt: time.Now(),
		},
		&File{
			Path:       "/tmp/test/hoge.txt",
			Data:       []byte("hogehoge"),
			FileMode:   0x800001ed,
			ModifiedAt: time.Now(),
		},
		&File{
			Path:       "/tmp/test/fuga.txt",
			Data:       []byte("fuga"),
			FileMode:   0x800001ed,
			ModifiedAt: time.Now(),
		},
	}
	files = map[string]*File{
		"/tmp/test":          f[0],
		"/tmp/test/hoge.txt": f[1],
		"/tmp/test/fuga.txt": f[2],
	}
	dirs = map[string][]string{
		"/tmp/test": []string{
			"hoge.png", "fuga.xml",
		},
	}
)

func TestNewFS(t *testing.T) {

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

func TestWithPrefix(t *testing.T) {
	path := "/static/"

	// get actual
	actual := NewFS(dirs, files).WithPrefix(path)

	// expected data
	expected := &FileSystem{
		Dirs:       dirs,
		Files:      files,
		pathPrefix: path,
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("expected and actual are not equal, expected %#+v and actual %#+v", expected, actual)
	}

}

func TestWithPrefix_PanicWhenNil(t *testing.T) {
	var fs *FileSystem
	defer func() {
		if err := recover(); err != nil {
			return
		}
		t.Fail()
	}()
	fs.WithPrefix("")
	t.Fail()
}

func TestNameResolute(t *testing.T) {
	prefix := "/static"
	path := "/static/tmp/test/fuga.txt"
	fs := NewFS(dirs, files).WithPrefix(prefix)

	actual := fs.nameResolute(path)

	expected := "/tmp/test/fuga.txt"

	if actual != expected {
		t.Errorf("path FileSystem.nameResolute returned: %s does not equal expected: %s \n", actual, expected)
	}

	fs = NewFS(dirs, files)

	actual = fs.nameResolute(path)

	expected = path

	if actual != expected {
		t.Errorf("path FileSystem.nameResolute returned: %s does not equal expected: %s \n", actual, expected)
	}
}
