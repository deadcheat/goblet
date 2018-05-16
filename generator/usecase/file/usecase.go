package file

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/deadcheat/awsset"
	"github.com/deadcheat/awsset/generator"
)

// UseCase file usecase
type UseCase struct {
	rr         generator.RegexpRepository
	fileMap    map[string]*awsset.File
	dirMap     map[string][]string
	validPaths []string
}

// New return new UseCase
func New(rr generator.RegexpRepository) generator.UseCase {
	return &UseCase{rr: rr}
}

// LoadFiles load files for given paths, except what matches given ignore path regex
func (u *UseCase) LoadFiles(paths []string, ignorePatterns []string) (*generator.Entity, error) {
	if err := u.rr.CompilePatterns(ignorePatterns); err != nil {
		return nil, err
	}

	u.fileMap = make(map[string]*awsset.File)
	u.validPaths = make([]string, 0)
	u.dirMap = make(map[string][]string)
	for i := range paths {
		path := paths[i]
		if err := u.addFile(path); err != nil {
			return nil, err
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
	if u.rr.MatchAny(path) {
		return nil
	}
	fi, err := os.Stat(path)
	if err != nil {
		return err
	}
	u.validPaths = append(u.validPaths, vPath)
	if fi.IsDir() {
		children := u.dirMap[vPath]
		if children == nil {
			children = make([]string, 0)
		}
		var files []os.FileInfo
		files, err = ioutil.ReadDir(path)
		if err != nil {
			log.Fatal(err)
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
		d := awsset.NewFromFileInfo(fi, vPath, nil)
		u.fileMap[vPath] = d
	} else {
		data, err := ioutil.ReadFile(path)
		if err != nil {
			log.Fatal(err)
		}
		file := awsset.NewFromFileInfo(fi, vPath, data)
		u.fileMap[vPath] = file
	}
	return nil
}
