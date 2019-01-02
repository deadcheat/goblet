package file

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/deadcheat/goblet"
	"github.com/deadcheat/goblet/generator"
)

var (
	// ErrFileIsNotMatchExpression return when regexp doesn't match specified name
	ErrFileIsNotMatchExpression = errors.New("file path does not match specified pattern")
)

// UseCase file usecase
type UseCase struct {
	rr         generator.RegexpRepository
	fileMap    map[string]*goblet.File
	dirMap     map[string][]string
	validPaths []string
}

// New return new UseCase
func New(rr generator.RegexpRepository) generator.UseCase {

	fileMap := make(map[string]*goblet.File)
	validPaths := make([]string, 0)
	dirMap := make(map[string][]string)
	return &UseCase{
		rr:         rr,
		fileMap:    fileMap,
		validPaths: validPaths,
		dirMap:     dirMap,
	}
}

// LoadFiles load files for given paths, except what matches given ignore path regex
func (u *UseCase) LoadFiles(paths []string, includePatterns []string) (*generator.Entity, error) {
	if err := u.rr.CompilePatterns(includePatterns); err != nil {
		return nil, err
	}

	for i := range paths {
		path := paths[i]
		if err := u.addFile(path); err != nil {
			log.Printf("path %s is not matched pattern given in 'expression(e)' flag")
			continue
		}
	}
	e := &generator.Entity{
		DirMap:  u.dirMap,
		FileMap: u.fileMap,
		Paths:   u.validPaths,
	}

	return e, nil
}

func (u *UseCase) addFile(path string) (err error) {
	vPath := filepath.Join("/", path)
	fi, err := os.Stat(path)
	if err != nil {
		return err
	}
	u.validPaths = append(u.validPaths, vPath)
	if !fi.IsDir() {
		if !u.rr.MatchAny(path) {
			return ErrFileIsNotMatchExpression
		}
		var data []byte
		data, err = ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		file := goblet.NewFromFileInfo(fi, vPath, data)
		u.fileMap[vPath] = file
		return nil
	}
	children := u.dirMap[vPath]
	if children == nil {
		children = make([]string, 0)
	}
	var files []os.FileInfo
	files, err = ioutil.ReadDir(path)
	if err != nil {
		return err
	}
	// add files recursively in dir
	for i := range files {
		f := files[i]
		childPath := filepath.Join(path, f.Name())
		err = u.addFile(childPath)
		if err != nil {
			return err
		}
		children = append(children, filepath.Base(childPath))
	}
	u.dirMap[vPath] = children
	d := goblet.NewFromFileInfo(fi, vPath, nil)
	u.fileMap[vPath] = d
	return nil
}
