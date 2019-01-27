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
	rr         []generator.PathMatcherRepository
	fileMap    map[string]*goblet.File
	dirMap     map[string][]string
	validPaths []string
}

// New return new UseCase
func New(rr []generator.PathMatcherRepository) generator.UseCase {

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

// LoadFiles load files for given paths, filtered by what matches given include path regex
func (u *UseCase) LoadFiles(paths []string, option generator.OptionFlagEntity) (*generator.Entity, error) {
	for i := range u.rr {
		repo := u.rr[i]
		if err := repo.Prepare(option); err != nil {
			return nil, err
		}
	}
	for i := range paths {
		path := paths[i]
		if err := u.addFile(path, option); err != nil {
			if err == ErrFileIsNotMatchExpression {
				log.Printf("%s is excluded because it is not matched with specified condition", path)
				continue
			}
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

func (u *UseCase) addFile(path string, option generator.OptionFlagEntity) (err error) {
	vPath := filepath.Join("/", path)
	fi, err := os.Stat(path)
	if err != nil {
		return err
	}
	if !fi.IsDir() {
		for i := range u.rr {
			repo := u.rr[i]
			if !repo.Match(path) {
				return ErrFileIsNotMatchExpression
			}
		}
		u.validPaths = append(u.validPaths, vPath)
		var data []byte
		data, err = ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		file := goblet.NewFromFileInfo(fi, vPath, data)
		u.fileMap[vPath] = file
		return nil
	}
	children, found := u.dirMap[vPath]
	if !found {
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
		err = u.addFile(childPath, option)
		if err != nil {
			if err == ErrFileIsNotMatchExpression {
				log.Printf("%s is excluded because it is not matched with specified condition", childPath)
				continue
			}
			return err
		}
		children = append(children, filepath.Base(childPath))
	}
	if option.ExcludeEmptyDir && len(children) == 0 {
		return nil
	}
	u.dirMap[vPath] = children
	u.validPaths = append(u.validPaths, vPath)
	u.fileMap[vPath] = goblet.NewFromFileInfo(fi, vPath, nil)
	return nil
}
