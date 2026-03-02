package build

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/google/blueprint"
)

type Config struct {
	SrcPath       string
	OutPath       string
	RelSrcPath    string
	CCompilerPath string
	CCompiler     string
	Arch          string
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

func (c *Config) AddCompileObjects(pctx blueprint.ModuleContext, srcs []string, cflags []string, compilers Compilers) []string {
	var objs []string

	for _, src := range srcs {
		// RelSrcPath is for build.ninja on output directory
		in := filepath.Join(c.RelSrcPath, pctx.ModuleDir(), src)

		base := filepath.Base(src)
		obj := filepath.Join("obj", pctx.ModuleName(), base+".o")
		dep := obj + ".d"
		objs = append(objs, obj)

		var cc string
		file_ext := filepath.Ext(src)
		if file_ext == ".cpp" || file_ext == ".cc" {
			cc = compilers.CXX
		} else {
			cc = compilers.CC
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
					"cflags":  strings.Join(cflags, " "),
					"depfile": dep,
				},
			},
		)
	}

	return objs
}
