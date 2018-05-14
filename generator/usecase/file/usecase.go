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
	fsMap      map[string]*awsset.FileSystem
	validPaths []string
}

// New return new UseCase
func New(rr generator.RegexpRepository) generator.UseCase {
	return &UseCase{rr: rr}
}

// LoadFiles load files for given paths, except what matches given ignore path regex
func (u *UseCase) LoadFiles(paths []string, ignorePatterns []string) (*generator.Entity, error) {
	u.rr.CompilePatterns(ignorePatterns)

	u.fsMap = make(map[string]*awsset.FileSystem)
	u.validPaths = make([]string, 0)
	for i := range paths {
		path := paths[i]
		u.addFile(path)
	}
	e := &generator.Entity{
		FsMap: u.fsMap,
		Paths: u.validPaths,
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
	if fi.IsDir() {
		var files []os.FileInfo
		files, err = ioutil.ReadDir(path)
		if err != nil {
			log.Fatal(err)
		}
		for i := range files {
			f := files[i]
			u.addFile(filepath.Join(path, f.Name()))
		}
		return
	}
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	fs := &awsset.FileSystem{
		FileInfo: fi,
		Path:     path,
		Data:     data,
	}
	u.fsMap[path] = fs
	u.validPaths = append(u.validPaths, path)
}
