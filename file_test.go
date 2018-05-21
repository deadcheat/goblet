package awsset

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func createTempDirAndFile(fileName, content string) (dir, path string, err error) {
	dir, err = ioutil.TempDir("", "tmpdir")
	if err != nil {
		return
	}

	path = filepath.Join(dir, fileName)
	if err = ioutil.WriteFile(path, []byte(content), 0666); err != nil {
		return
	}
	return
}

func TestNewFromFileInfo(t *testing.T) {
	contentStr := "hello world!"
	dir, path, err := createTempDirAndFile("temp.txt", contentStr)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir) // clean up
	fi, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}

	f := NewFromFileInfo(fi, path, []byte(contentStr))
	if f.Path != path || string(f.Data) != contentStr {
		t.Errorf("NewFromFileInfo returned unexpected File %#v", f)
	}
}

func TestNewFile(t *testing.T) {
	contentStr := "hello world!"
	dir, path, err := createTempDirAndFile("temp.txt", contentStr)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir) // clean up
	fi, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}

	f := NewFile(path, []byte(contentStr), fi.Mode(), fi.ModTime())
	if f == nil || f.Path != path || string(f.Data) != contentStr {
		t.Errorf("NewFromFileInfo returned unexpected File %#v", f)
	}

}
