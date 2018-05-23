package file

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/deadcheat/awsset/generator/mock"
	"github.com/golang/mock/gomock"
)

func TestAddFile(t *testing.T) {
	// Prepare dir and file
	contentStr := "hello world!"
	dir, path, err := createTempDirAndFile("temp.txt", contentStr)
	defer os.RemoveAll(dir) // clean up

	// Prepare mock
	c := gomock.NewController(t)
	defer c.Finish()

	m := mock.NewMockRegexpRepository(c)
	m.EXPECT().MatchAny(path).Return(true)
	// create usecase
	iu := New(m)

	u := iu.(*UseCase)
	// success pattern
	err = u.addFile(path)
	if err != nil {
		t.Error("addFile should not return any errors", err)
	}
}

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
