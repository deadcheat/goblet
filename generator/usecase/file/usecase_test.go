package file

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/deadcheat/gonch"

	"github.com/deadcheat/goblet/generator"
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

	m := mock.NewMockPathMatcherRepository(c)
	// successful pattern
	filenameSuccess := "success.txt"
	path := filepath.Join(d.Dir(), filenameSuccess)
	d.AddFile("success", filenameSuccess, content, 0666)
	m.EXPECT().Match(path).Return(true)

	// add Expects for Prepare
	emptyPatterns := make([]string, 0)
	op := generator.OptionFlagEntity{
		IncludePatterns: emptyPatterns,
	}
	m.EXPECT().Prepare(op).AnyTimes().Return(nil)

	// create usecase
	iu := New([]generator.PathMatcherRepository{m})
	_, err := iu.LoadFiles([]string{d.Dir()}, op)
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

	m := mock.NewMockPathMatcherRepository(c)

	// add Expects for Prepare
	emptyPatterns := make([]string, 0)
	dummyError := errors.New("dummy")
	op := generator.OptionFlagEntity{
		IncludePatterns: emptyPatterns,
	}
	m.EXPECT().Prepare(op).AnyTimes().Return(dummyError)

	// create usecase
	iu := New([]generator.PathMatcherRepository{m})
	_, err := iu.LoadFiles([]string{d.Dir()}, op)
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

	m := mock.NewMockPathMatcherRepository(c)

	// should not be contained
	fileDoesNotPermitted := "fail.txt"
	permitWrongPath := filepath.Join(d.Dir(), fileDoesNotPermitted)
	d.AddFile("success", fileDoesNotPermitted, content, 0000)
	m.EXPECT().Match(permitWrongPath).AnyTimes().Return(true)

	// add Expects for Prepare
	emptyPatterns := make([]string, 0)
	op := generator.OptionFlagEntity{
		IncludePatterns: emptyPatterns,
	}
	m.EXPECT().Prepare(op).AnyTimes().Return(nil)

	// create usecase
	iu := New([]generator.PathMatcherRepository{m})
	_, err := iu.LoadFiles([]string{d.Dir(), permitWrongPath}, op)
	if err == nil {
		t.Error("addFile should return error")
	}

}

func TestLoadFilesSkipNotMatchedFile(t *testing.T) {
	// Prepare dir and file
	contentStr := "hello world!"
	content := []byte(contentStr)
	d := gonch.New("", "tmpdir")
	defer d.Close() // clean up

	// Prepare mock
	c := gomock.NewController(t)
	defer c.Finish()

	m := mock.NewMockPathMatcherRepository(c)

	// should not be contained
	fileDoesNotMatch := "doesnotmatch.txt"
	notMatchPath := filepath.Join(d.Dir(), fileDoesNotMatch)
	d.AddFile("wrong path file", fileDoesNotMatch, content, 0666)
	m.EXPECT().Match(notMatchPath).AnyTimes().Return(false)

	// add Expects for Prepare
	emptyPatterns := make([]string, 0)
	op := generator.OptionFlagEntity{
		IncludePatterns: emptyPatterns,
	}
	m.EXPECT().Prepare(op).AnyTimes().Return(nil)

	// create usecase
	iu := New([]generator.PathMatcherRepository{m})
	result, err := iu.LoadFiles([]string{d.Dir(), notMatchPath}, op)
	if err != nil {
		t.Error("addFile should not return error")
	}

	if len(result.DirMap[d.Dir()]) != 0 {
		t.Error("some files are included in dirmap")
	}

	_, ok := result.FileMap[notMatchPath]
	if ok {
		t.Error("unmatched file are included in filemap")
	}

	for _, v := range result.Paths {
		if v == notMatchPath {
			t.Error("unmatched file are included in path slice")
		}
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

	m := mock.NewMockPathMatcherRepository(c)
	// successful pattern
	filenameSuccess := "success.txt"
	path := filepath.Join(d.Dir(), filenameSuccess)
	d.AddFile("success", filenameSuccess, content, 0666)
	m.EXPECT().Match(path).Return(true)
	// should not be contained
	fileDoesNotMatch := "doesnotmatch.txt"
	pathDoesNotMatch := filepath.Join(d.Dir(), fileDoesNotMatch)
	d.AddFile("wrong path file", fileDoesNotMatch, content, 0666)
	m.EXPECT().Match(pathDoesNotMatch).Return(false)
	// file can not be opened
	d.AddFile("wrong path file", pathDoesNotMatch, content, 0666)
	closedPath := filepath.Join(d.Dir(), "closed.txt")
	ioutil.WriteFile(closedPath, []byte(""), 0000)
	m.EXPECT().Match(closedPath).Return(true)

	// create usecase
	iu := New([]generator.PathMatcherRepository{m})

	u := iu.(*UseCase)
	// success pattern
	err := u.addFile(path)
	if err != nil {
		t.Error("addFile should not return any errors", err)
	}
	err = u.addFile(pathDoesNotMatch)
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
	m := mock.NewMockPathMatcherRepository(c)

	// successful pattern
	filenameSuccess1 := "/child/success1.txt"
	path1 := filepath.Join(d.Dir(), filenameSuccess1)
	d.AddFile("success", filenameSuccess1, content, 0666)
	m.EXPECT().Match(path1).Return(true)

	filenameSuccess2 := "/child/success2.txt"
	path2 := filepath.Join(d.Dir(), filenameSuccess2)
	d.AddFile("success", filenameSuccess2, content, 0666)
	m.EXPECT().Match(path2).Return(true)

	// permission denied dir
	deniedDir := "/permissiondeny"
	if err := d.AddDir("faildir", deniedDir, 0000); err != nil {
		panic(err)
	}

	filenameSuccess3 := "/fail.txt"
	path3 := filepath.Join(d.Dir(), filenameSuccess3)
	d.AddFile("success", filenameSuccess3, content, 0000)
	m.EXPECT().Match(path3).AnyTimes().Return(true)

	// permission denied file in permitted dir
	permittedDir := "/permitteddir"
	if err := d.AddDir("permitteddir", permittedDir, os.ModePerm); err != nil {
		panic(err)
	}
	deniedFile := filepath.Join(permittedDir, "permitted_denied.txt")
	deniedFilePath := filepath.Join(d.Dir(), deniedFile)
	if err := d.AddFile("permitted_denied", deniedFile, content, 0000); err != nil {
		panic(err)
	}
	m.EXPECT().Match(deniedFilePath).Return(true)

	// create usecase
	iu := New([]generator.PathMatcherRepository{m})
	u := iu.(*UseCase)

	// success pattern
	targetDir := filepath.Join(d.Dir(), "/child")
	if err := u.addFile(targetDir); err != nil {
		t.Error("addFile should not return any errors", err)
	}

	// when directory is not permittedf
	if err := u.addFile(filepath.Join(d.Dir(), deniedDir)); err == nil {
		t.Error("addFile should return any error when dir is denied")
	}

	// when file in dir is permitted
	if err := u.addFile(filepath.Join(d.Dir(), permittedDir)); err == nil {
		t.Error("addFile should return any error when file in dir is denied")
	}
	// when file in dir is not permitted
	if err := u.addFile(path3); err == nil {
		t.Error("addFile should return any error when file in dir is denied")
	}
}
