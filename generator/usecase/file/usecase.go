package file

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"

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
		u.addFile(path)
	}
	sort.Strings(u.validPaths)
	e := &generator.Entity{
		DirMap:  u.dirMap,
		FileMap: u.fileMap,
		Paths:   u.validPaths,
	}

	return e, nil
}

func (u *UseCase) addFile(path string) {
	if u.rr.MatchAny(path) {
		return
	}
	fi, err := os.Stat(path)
	if err != nil {
		return
	}
	u.validPaths = append(u.validPaths, path)
	if fi.IsDir() {
		children := u.dirMap[path]
		if children == nil {
			children = make([]string, 0)
		}
		var files []os.FileInfo
		files, err = ioutil.ReadDir(path)
		if err != nil {
			log.Fatal(err)
		}
		for i := range files {
			f := files[i]
			childPath := filepath.Join(path, f.Name())
			children = append(children, childPath)
			u.addFile(childPath)
		}
		u.dirMap[path] = children
		d := awsset.NewFromFileInfo(fi, path, nil)
		u.fileMap[path] = d
		return
	}
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	fs := awsset.NewFromFileInfo(fi, path, data)
	u.fileMap[path] = fs
}
