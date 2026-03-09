package build

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/blueprint"
)

type Config struct {
	// Main path of source code
	SrcPath string

	// Main path of output directory
	OutPath string

	// Rel path for Src & Inc path
	RelSrcPath string

	// C & C++ compiler path
	// (Most for Embedded Compiler path)
	CCompilerPath string

	// Compiler target
	// (like GCC or Clang)
	CCompiler string

	// Target Architecture
	// (like x86_64 or arm64)
	Arch string
}

type BuildInfo struct {
	Srcs      []string
	Incs      []string
	Cflags    []string
	Libs      []string
	Compilers Compilers
}

func incPathsToOpts(paths []string) string {
	var ret []string

	for _, v := range paths {
		ret = append(ret, "-I"+v)
	}

	return strings.Join(ret, " ")
}

func (c *Config) CreateOutPath() error {
	return os.MkdirAll(c.OutPath, 0o755)
}

func (c *Config) SetSrcPath(pctx *blueprint.Context) {
	pctx.SetSrcDir(c.SrcPath)
}

func (c *Config) CreateBuildNinja() (*os.File, error) {
	ninjaPath := filepath.Join(c.OutPath, "build.ninja")
	return os.Create(ninjaPath)
}

func (c *Config) SearchBuildFiles() ([]string, error) {
	var bpFiles []string
	var retErr error

	retErr = filepath.WalkDir(c.SrcPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() && filepath.Base(path) == "Build.bp" {
			bpFiles = append(bpFiles, path)
		}

		return nil
	})

	return bpFiles, retErr
}

func (c *Config) GetRelIncPath(pctx blueprint.ModuleContext, inIncs []string) []string {
	var incs []string

	for _, inc := range inIncs {
		relInc := filepath.Join(c.RelSrcPath, pctx.ModuleDir(), inc)
		incs = append(incs, relInc)
	}

	return incs
}

func (c *Config) AddCompileObjects(pctx blueprint.ModuleContext, buildInfo BuildInfo) []string {
	var objs []string

	for _, src := range buildInfo.Srcs {
		// RelSrcPath is for build.ninja on output directory
		in := filepath.Join(c.RelSrcPath, pctx.ModuleDir(), src)

		base := filepath.Base(src)
		obj := filepath.Join("obj", pctx.ModuleName(), base+".o")
		dep := obj + ".d"
		objs = append(objs, obj)

		var cc string
		file_ext := filepath.Ext(src)
		if file_ext == ".cpp" || file_ext == ".cc" {
			cc = buildInfo.Compilers.CXX
		} else {
			cc = buildInfo.Compilers.CC
		}

		pctx.Build(
			pkgCtx,
			blueprint.BuildParams{
				Rule:    CCRule,
				Outputs: []string{obj},
				Inputs:  []string{in},
				Depfile: dep,
				Deps:    blueprint.DepsGCC,
				Args: map[string]string{
					"cc":      cc,
					"cflags":  strings.Join(buildInfo.Cflags, " "),
					"depfile": dep,
					"incs":    incPathsToOpts(buildInfo.Incs),
				},
			},
		)
	}

	return objs
}
