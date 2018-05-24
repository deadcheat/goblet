package file

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/deadcheat/goblet/generator/mock"
	"github.com/golang/mock/gomock"
)

func TestAddFile_ForSingleFiles(t *testing.T) {
	// Prepare dir and file
	contentStr := "hello world!"
	dir, path, _ := createTempDirAndFile("temp.txt", contentStr)
	defer os.RemoveAll(dir) // clean up

	// Prepare mock
	c := gomock.NewController(t)
	defer c.Finish()

	m := mock.NewMockRegexpRepository(c)
	// successfull pattern
	m.EXPECT().MatchAny(path).Return(true)
	// should not be contained
	wrongPath := filepath.Join(dir, "doesnotmatch.txt")
	ioutil.WriteFile(wrongPath, []byte("test"), 0666)
	m.EXPECT().MatchAny(wrongPath).Return(false)
	// file can not be opened
	closedPath := filepath.Join(dir, "closed.txt")
	ioutil.WriteFile(closedPath, []byte(""), 0000)
	m.EXPECT().MatchAny(closedPath).Return(true)

	// create usecase
	iu := New(m)

	u := iu.(*UseCase)
	// success pattern
	err := u.addFile(path)
	if err != nil {
		t.Error("addFile should not return any errors", err)
	}
	err = u.addFile(wrongPath)
	if err != ErrFileIsNotMatchExpression {
		t.Error("addFile should return ErrFileIsNotMatchExpression but returned ", err)
	}
	// filename does not exist
	brokenPath := filepath.Join(dir, "/this/is/match/but/broken.txt")
	err = u.addFile(brokenPath)
	if err == nil {
		t.Error("addFile should return some error")
	}
	err = u.addFile(closedPath)
	if err == nil {
		t.Error("addFile should return some error")
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
