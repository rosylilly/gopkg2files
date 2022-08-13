package gopkg2files

import (
	"fmt"
	"go/build"
	"io"
	"log"
	"path/filepath"
	"sync"
)

type Resolver struct {
	option *Option
	logger *log.Logger
	build  build.Context
	mode   build.ImportMode

	packagesMu sync.Mutex
	packages   map[string]*build.Package

	filesMu      sync.Mutex
	GoFiles      []string
	CgoFiles     []string
	TestGoFiles  []string
	XTestGoFiles []string
}

func NewResolver(option *Option) *Resolver {
	return &Resolver{
		option:       option,
		logger:       option.Logger(),
		build:        build.Default,
		mode:         build.ImportComment,
		packagesMu:   sync.Mutex{},
		packages:     map[string]*build.Package{},
		filesMu:      sync.Mutex{},
		GoFiles:      []string{},
		CgoFiles:     []string{},
		TestGoFiles:  []string{},
		XTestGoFiles: []string{},
	}
}

func (r *Resolver) PrintFiles(w io.Writer) error {
	for _, file := range r.GoFiles {
		if _, err := fmt.Fprintln(w, file); err != nil {
			return err
		}
	}
	if r.option.Cgo {
		for _, file := range r.CgoFiles {
			if _, err := fmt.Fprintln(w, file); err != nil {
				return err
			}
		}
	}
	if r.option.Test {
		for _, file := range r.TestGoFiles {
			if _, err := fmt.Fprintln(w, file); err != nil {
				return err
			}
		}
	}
	if r.option.XTest {
		for _, file := range r.XTestGoFiles {
			if _, err := fmt.Fprintln(w, file); err != nil {
				return err
			}
		}
	}
	return nil
}

func (r *Resolver) ResolveAll() error {
	for _, pkg := range r.option.Packages {
		if err := r.Resolve(pkg); err != nil {
			return err
		}
	}
	return nil
}

func (r *Resolver) Resolve(nameOrDir string) error {
	if err := r.ResolvePackage(nameOrDir); err != nil {
		return r.ResolveDir(nameOrDir)
	}
	return nil
}

func (r *Resolver) ResolveDir(dir string) error {
	pkg, err := r.build.ImportDir(dir, r.mode)
	if err != nil {
		return err
	}

	return r.resolve(pkg)
}

func (r *Resolver) ResolvePackage(name string) error {
	pkg, err := r.build.Import(name, r.option.WorkingDirectory, r.mode)
	if err != nil {
		return err
	}

	return r.resolve(pkg)
}

func (r *Resolver) resolve(pkg *build.Package) error {
	r.packagesMu.Lock()

	if _, ok := r.packages[pkg.ImportPath]; ok {
		r.packagesMu.Unlock()
		return nil
	}

	r.packages[pkg.ImportPath] = pkg
	r.packagesMu.Unlock()

	r.logger.Printf("package: %s", pkg.ImportPath)
	if pkg.Goroot && !r.option.Goroot {
		r.logger.Println("skip: $GOROOT")
		return nil
	}

	r.filesMu.Lock()
	goFiles := make([]string, len(pkg.GoFiles))
	for i, file := range pkg.GoFiles {
		goFiles[i] = filepath.Join(pkg.Dir, file)
	}
	r.GoFiles = append(r.GoFiles, goFiles...)
	cgoFiles := make([]string, len(pkg.CgoFiles))
	for i, file := range pkg.CgoFiles {
		cgoFiles[i] = filepath.Join(pkg.Dir, file)
	}
	r.CgoFiles = append(r.CgoFiles, cgoFiles...)
	testGoFiles := make([]string, len(pkg.TestGoFiles))
	for i, file := range pkg.TestGoFiles {
		testGoFiles[i] = filepath.Join(pkg.Dir, file)
	}
	r.TestGoFiles = append(r.TestGoFiles, testGoFiles...)
	xTestGoFiles := make([]string, len(pkg.XTestGoFiles))
	for i, file := range pkg.XTestGoFiles {
		xTestGoFiles[i] = filepath.Join(pkg.Dir, file)
	}
	r.XTestGoFiles = append(r.XTestGoFiles, xTestGoFiles...)
	r.filesMu.Unlock()

	for _, imp := range pkg.Imports {
		if err := r.Resolve(imp); err != nil {
			return err
		}
	}
	return nil
}
