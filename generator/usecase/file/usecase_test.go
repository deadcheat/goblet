package file

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/deadcheat/gonch"

	"github.com/deadcheat/goblet/generator/mock"
	"github.com/golang/mock/gomock"
)

func TestAddFile_ForSingleFiles(t *testing.T) {
	// Prepare dir and file
	contentStr := "hello world!"
	content := []byte(contentStr)
	d := gonch.New("", "tmpdir")
	defer d.Close() // clean up

	// Prepare mock
	c := gomock.NewController(t)
	defer c.Finish()

	m := mock.NewMockRegexpRepository(c)
	// successfull pattern
	filenameSuccess := "success.txt"
	path := filepath.Join(d.Dir(), filenameSuccess)
	d.AddFile("success", filenameSuccess, content, 0666)
	m.EXPECT().MatchAny(path).Return(true)
	// should not be contained
	fileWrongPath := "doesnotmatch.txt"
	wrongPath := filepath.Join(d.Dir(), fileWrongPath)
	d.AddFile("wrong path file", fileWrongPath, content, 0666)
	m.EXPECT().MatchAny(wrongPath).Return(false)
	// file can not be opened
	d.AddFile("wrong path file", wrongPath, content, 0666)
	closedPath := filepath.Join(d.Dir(), "closed.txt")
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
	brokenPath := filepath.Join(d.Dir(), "/this/is/match/but/broken.txt")
	err = u.addFile(brokenPath)
	if err == nil {
		t.Error("addFile should return some error")
	}
	err = u.addFile(closedPath)
	if err == nil {
		t.Error("addFile should return some error")
	}
}
