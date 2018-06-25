package file

import (
	"errors"
	"flag"
	"path/filepath"
	"testing"
	"time"

	"github.com/deadcheat/goblet"
	"github.com/deadcheat/gonch"

	"github.com/deadcheat/goblet/generator"
	"github.com/deadcheat/goblet/generator/mock"
	"github.com/deadcheat/goblet/generator/values"
	"github.com/golang/mock/gomock"

	"github.com/urfave/cli"
)

func TestMount(t *testing.T) {

	// Prepare mock
	c := gomock.NewController(t)
	defer c.Finish()

	m := mock.NewMockUseCase(c)
	p := New(m)
	// Success Pattern
	a := &cli.App{}
	if err := p.Mount(a); err != nil {
		t.Error("Mount should not return any errors", err)
	}

	// error pattern
	if err := p.Mount(struct{}{}); err != ErrIllegalMount {
		t.Error("Mount should return ErrIllegalMount", err)
	}
}

var fs = flag.NewFlagSet("test", flag.ContinueOnError)

func TestActionNoArgs(t *testing.T) {

	// Prepare mock
	c := gomock.NewController(t)
	defer c.Finish()

	m := mock.NewMockUseCase(c)
	p := New(m)
	a := &cli.App{}
	p.Mount(a)
	a.Flags = values.FlagDefs
	ctx := cli.NewContext(a, fs, nil)
	if err := p.action(ctx); err != ErrNoArguments {
		t.Error("Mount should return ErrNoArguments", err)
	}
}

var errTest = errors.New("error test")

func TestActionFailLoadFiles(t *testing.T) {

	// Prepare mock
	c := gomock.NewController(t)
	defer c.Finish()

	m := mock.NewMockUseCase(c)
	m.EXPECT().LoadFiles([]string{"config"}, nil).Return(nil, errTest)

	p := New(m)
	a := &cli.App{}
	p.Mount(a)
	set := flag.NewFlagSet("test", flag.ContinueOnError)
	if err := set.Parse([]string{"config"}); err != nil {
		panic(err)
	}

	a.Flags = values.FlagDefs
	ctx := cli.NewContext(a, set, nil)
	if err := p.action(ctx); err != errTest {
		t.Error("Mount should return errTest", err)
	}
}

func TestActionFailGoformat(t *testing.T) {

	// Prepare mock
	c := gomock.NewController(t)
	defer c.Finish()

	m := mock.NewMockUseCase(c)
	m.EXPECT().LoadFiles([]string{"config"}, nil).Return(&generator.Entity{FileMap: nil}, nil)

	p := New(m)
	a := &cli.App{}
	p.Mount(a)
	set := flag.NewFlagSet("test", flag.ContinueOnError)
	if err := set.Parse([]string{"config"}); err != nil {
		panic(err)
	}

	a.Flags = values.FlagDefs
	ctx := cli.NewContext(a, set, nil)
	if err := p.action(ctx); err == nil {
		t.Error("Mount should return errTest", err)
	}
}

func TestActionSuccessWithStdout(t *testing.T) {

	// Prepare mock
	c := gomock.NewController(t)
	defer c.Finish()

	mockReturn := &generator.Entity{
		DirMap: map[string][]string{
			"/tmpdir/testdir": {"testfile.png"},
		},
		FileMap: map[string]*goblet.File{
			"/tmpdir/testdir/testfile.png": {
				Path:       "/tmpdir/testdir/testfile.png",
				Data:       []byte("test"),
				FileMode:   0x1a4,
				ModifiedAt: time.Now(),
			},
			"/tmpdir/testdir": {
				Path:       "/tmpdir/testdir",
				Data:       []byte("test"),
				FileMode:   0x800001ed,
				ModifiedAt: time.Now(),
			},
		},
		Paths: []string{
			"/tmpdir/testdir",
			"/tmpdir/testdir/testfile.png",
		},
	}

	m := mock.NewMockUseCase(c)
	m.EXPECT().LoadFiles([]string{"config"}, nil).Return(mockReturn, nil)

	p := New(m)
	a := &cli.App{}
	p.Mount(a)
	set := flag.NewFlagSet("test", flag.ContinueOnError)
	if err := set.Parse([]string{"config"}); err != nil {
		panic(err)
	}

	a.Flags = values.FlagDefs
	set.String("name", "Asset", "")
	set.String("package", "assets", "")

	// stdout pattern
	ctx := cli.NewContext(a, set, nil)
	if err := p.action(ctx); err != nil {
		t.Error("Mount should not return any error: ", err)
	}
}

func TestActionSuccessWithFile(t *testing.T) {

	// Prepare mock
	c := gomock.NewController(t)
	defer c.Finish()

	mockReturn := &generator.Entity{
		DirMap: map[string][]string{
			"/tmpdir/testdir": {"testfile.png"},
		},
		FileMap: map[string]*goblet.File{
			"/tmpdir/testdir/testfile.png": {
				Path:       "/tmpdir/testdir/testfile.png",
				Data:       []byte("test"),
				FileMode:   0x1a4,
				ModifiedAt: time.Now(),
			},
			"/tmpdir/testdir": {
				Path:       "/tmpdir/testdir",
				Data:       []byte("test"),
				FileMode:   0x800001ed,
				ModifiedAt: time.Now(),
			},
		},
		Paths: []string{
			"/tmpdir/testdir",
			"/tmpdir/testdir/testfile.png",
		},
	}

	m := mock.NewMockUseCase(c)
	m.EXPECT().LoadFiles([]string{"config"}, nil).Return(mockReturn, nil)

	p := New(m)
	a := &cli.App{}
	p.Mount(a)
	set := flag.NewFlagSet("test", flag.ContinueOnError)
	if err := set.Parse([]string{"config"}); err != nil {
		panic(err)
	}

	a.Flags = values.FlagDefs
	set.String("name", "Asset", "")
	set.String("package", "assets", "")

	// out to file pattern
	d := gonch.New("", "tmpdir")
	defer d.Close()
	contentStr := "hello world!"
	content := []byte(contentStr)
	d.AddFile("testfile", "temp.txt", content, 0777)
	f, _ := d.File("testfile")
	set.String("out", f.Path, "")
	ctx := cli.NewContext(a, set, nil)
	if err := p.action(ctx); err != nil {
		t.Error("Mount should not return any error: ", err)
	}
}

func TestActionFailWhenCouldNotOpenFile(t *testing.T) {

	// Prepare mock
	c := gomock.NewController(t)
	defer c.Finish()

	mockReturn := &generator.Entity{
		DirMap: map[string][]string{
			"/tmpdir/testdir": {"testfile.png"},
		},
		FileMap: map[string]*goblet.File{
			"/tmpdir/testdir/testfile.png": {
				Path:       "/tmpdir/testdir/testfile.png",
				Data:       []byte("test"),
				FileMode:   0x1a4,
				ModifiedAt: time.Now(),
			},
			"/tmpdir/testdir": {
				Path:       "/tmpdir/testdir",
				Data:       []byte("test"),
				FileMode:   0x800001ed,
				ModifiedAt: time.Now(),
			},
		},
		Paths: []string{
			"/tmpdir/testdir",
			"/tmpdir/testdir/testfile.png",
		},
	}

	m := mock.NewMockUseCase(c)
	m.EXPECT().LoadFiles([]string{"config"}, nil).Return(mockReturn, nil)

	p := New(m)
	a := &cli.App{}
	p.Mount(a)
	set := flag.NewFlagSet("test", flag.ContinueOnError)
	if err := set.Parse([]string{"config"}); err != nil {
		panic(err)
	}

	a.Flags = values.FlagDefs
	set.String("name", "Asset", "")
	set.String("package", "assets", "")

	// out to file pattern
	d := gonch.New("", "tmpdir")
	defer d.Close()

	d.AddDir("testdir", "testdir", 0111)
	f, _ := d.File("testdir")
	set.String("out", filepath.Join(f.Path, "testfile.txt"), "")
	ctx := cli.NewContext(a, set, nil)
	if err := p.action(ctx); err == nil {
		t.Error("Mount should return any error: ", err)
	}
}
