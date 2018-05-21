package awsset

import (
	"io"
	"io/ioutil"
	"net/http"
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

func TestAsFileInfo(t *testing.T) {
	contentStr := "hello world!"
	content := []byte(contentStr)
	fileName := "temp.txt"
	dir, path, err := createTempDirAndFile(fileName, contentStr)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir) // clean up
	fi, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	var f os.FileInfo = NewFile(path, content, fi.Mode(), fi.ModTime())

	if f.Name() != fileName ||
		f.Size() != int64(len(content)) ||
		f.Mode() != fi.Mode() ||
		f.ModTime() != fi.ModTime() ||
		f.IsDir() ||
		f.Sys() != nil {
		t.Errorf("File maybe an unexpected File %#v", f)
	}
}

func TestReaddir(t *testing.T) {
	contentStr := "hello world!"
	content := []byte(contentStr)
	fileName := "temp.txt"
	dir, path, err := createTempDirAndFile(fileName, contentStr)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir) // clean up
	fi, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	var f http.File = NewFile(path, content, fi.Mode(), fi.ModTime())
	files, err := f.Readdir(1)
	if files != nil || err != nil {
		t.Errorf("File maybe an unexpected File %#v", f)
	}
}

func TestReadAndClose(t *testing.T) {
	contentStr := "hello world!"
	content := []byte(contentStr)
	fileName := "temp.txt"
	dir, path, err := createTempDirAndFile(fileName, contentStr)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir) // clean up
	fi, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	var f http.File = NewFile(path, content, fi.Mode(), fi.ModTime())
	readContent := make([]byte, len(content))
	n, err := f.Read(readContent)
	actualFile, ok := f.(*File)
	if n != len(content) || err != nil || string(readContent) != contentStr || !ok || actualFile.buf == nil {
		t.Errorf("File maybe an unexpected File %#v", f)
	}
	err = f.Close()
	if err != nil || actualFile.buf != nil {
		t.Errorf("File maybe an unexpected File %#v", f)
	}
}

func TestStat(t *testing.T) {
	contentStr := "hello world!"
	content := []byte(contentStr)
	fileName := "temp.txt"
	dir, path, err := createTempDirAndFile(fileName, contentStr)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir) // clean up
	fi, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	var f http.File = NewFile(path, content, fi.Mode(), fi.ModTime())

	statFi, err := f.Stat()
	if err != nil || statFi.Name() != fi.Name() {
		t.Errorf("File maybe an unexpected FileInfo %#v, expected %#v", statFi, fi)
	}
}

func TestSeek(t *testing.T) {
	contentStr := "hello world!"
	content := []byte(contentStr)
	fileName := "temp.txt"
	dir, path, err := createTempDirAndFile(fileName, contentStr)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir) // clean up
	fi, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	var f http.File = NewFile(path, content, fi.Mode(), fi.ModTime())

	offset := int64(6)
	expectedLen := 6
	seeked, err := f.Seek(offset, io.SeekStart)
	readContent := make([]byte, expectedLen)
	if seeked != int64(expectedLen) || err != nil {
		t.Errorf("Seek returned unexpected result %#v, %#v", seeked, err)
	}
	n, err := f.Read(readContent)
	if n != expectedLen || err != nil || string(readContent) != contentStr[offset:] {
		t.Errorf("Read returned unexpected result %#v, expected: %#v", string(readContent), contentStr[offset:])
	}
	_ = f.Close()
}
