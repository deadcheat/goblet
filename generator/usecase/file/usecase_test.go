package file

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/deadcheat/gonch"

	"github.com/deadcheat/goblet/generator/mock"
	"github.com/golang/mock/gomock"
)

func TestLoadFilesSuccess(t *testing.T) {
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

	// add Expects for CompilePatterns
	emptyPatterns := make([]string, 0)
	m.EXPECT().CompilePatterns(emptyPatterns).Return(nil)

	// create usecase
	iu := New(m)

	_, err := iu.LoadFiles([]string{d.Dir()}, emptyPatterns)
	if err != nil {
		t.Error("addFile should not return any errors", err)
	}
}

func TestLoadFilesErrorCompile(t *testing.T) {
	// Prepare dir and file
	d := gonch.New("", "tmpdir")
	defer d.Close() // clean up

	// Prepare mock
	c := gomock.NewController(t)
	defer c.Finish()

	m := mock.NewMockRegexpRepository(c)

	// add Expects for CompilePatterns
	emptyPatterns := make([]string, 0)
	dummyError := errors.New("dummy")
	m.EXPECT().CompilePatterns(emptyPatterns).Return(dummyError)

	// create usecase
	iu := New(m)

	_, err := iu.LoadFiles([]string{d.Dir()}, emptyPatterns)
	if err == nil {
		t.Error("addFile should return error")
	}
}

func TestLoadFilesErrorAddFile(t *testing.T) {
	// Prepare dir and file
	contentStr := "hello world!"
	content := []byte(contentStr)
	d := gonch.New("", "tmpdir")
	defer d.Close() // clean up

	// Prepare mock
	c := gomock.NewController(t)
	defer c.Finish()

	m := mock.NewMockRegexpRepository(c)

	// should not be contained
	fileWrongPath := "doesnotmatch.txt"
	wrongPath := filepath.Join(d.Dir(), fileWrongPath)
	d.AddFile("wrong path file", fileWrongPath, content, 0666)
	m.EXPECT().MatchAny(wrongPath).Return(false)

	// add Expects for CompilePatterns
	emptyPatterns := make([]string, 0)
	m.EXPECT().CompilePatterns(emptyPatterns).Return(nil)

	// create usecase
	iu := New(m)

	_, err := iu.LoadFiles([]string{d.Dir()}, emptyPatterns)
	if err == nil {
		t.Error("addFile should return error")
	}
}

func TestAddFileForSingleFiles(t *testing.T) {
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

func TestAddFileForDirectory(t *testing.T) {
	contentStr := "hello world!"
	content := []byte(contentStr)
	d := gonch.New("", "tmpdir")
	defer d.Close() // clean up

	// Prepare mock
	c := gomock.NewController(t)
	defer c.Finish()
	m := mock.NewMockRegexpRepository(c)

	// successfull pattern
	filenameSuccess1 := "/child/success1.txt"
	path1 := filepath.Join(d.Dir(), filenameSuccess1)
	d.AddFile("success", filenameSuccess1, content, 0666)
	m.EXPECT().MatchAny(path1).Return(true)

	filenameSuccess2 := "/child/success2.txt"
	path2 := filepath.Join(d.Dir(), filenameSuccess2)
	d.AddFile("success", filenameSuccess2, content, 0666)
	m.EXPECT().MatchAny(path2).Return(true)

	// permission denied dir
	deniedDir := "/permissiondeny"
	if err := d.AddDir("faildir", deniedDir, 0000); err != nil {
		panic(err)
	}

	// permission denied file in permitted dir
	permittedDir := "/permitteddir"
	if err := d.AddDir("permitteddir", permittedDir, os.ModePerm); err != nil {
		panic(err)
	}
	deniedFile := filepath.Join(permittedDir, "permitted_denied.txt")
	deniedFilePath := filepath.Join(d.Dir(), deniedFile)
	err := d.AddFile("permitted_denied", deniedFile, content, 0000)
	m.EXPECT().MatchAny(deniedFilePath).Return(true)

	// create usecase
	iu := New(m)
	u := iu.(*UseCase)

	// success pattern
	targetDir := filepath.Join(d.Dir(), "/child")
	err = u.addFile(targetDir)
	if err != nil {
		t.Error("addFile should not return any errors", err)
	}

	// when directory is not permittedf
	err = u.addFile(filepath.Join(d.Dir(), deniedDir))
	if err == nil {
		t.Error("addFile should return any error when dir is denied")
	}

	// when file in dir is permitted
	err = u.addFile(filepath.Join(d.Dir(), permittedDir))
	if err == nil {
		t.Error("addFile should return any error when file in dir is denied")
	}
}
