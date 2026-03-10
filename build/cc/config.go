package cc

import (
	"path/filepath"
	"strings"

	"github.com/PoppingBoba/boong/build/common"
	"github.com/google/blueprint"
)

type Config struct {
	common.ConfigBase

	// C & C++ compiler path
	// (Most for Embedded Compiler path)
	CCompilerPath string

	// Compiler target
	// (like GCC or Clang)
	CCompiler string
}

type BuildInfo struct {
	Srcs      []string
	Incs      []string
	Cflags    []string
	Cppflags  []string
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

func (c *Config) AddCompileObjects(pctx blueprint.ModuleContext, buildInfo BuildInfo) []string {
	var objs []string
	var flags []string

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
			flags = append(buildInfo.Cflags, buildInfo.Cppflags...)
			cc = buildInfo.Compilers.CXX
		} else {
			flags = buildInfo.Cflags
			cc = buildInfo.Compilers.CC
		}

		pctx.Build(
			common.PkgCtx,
			blueprint.BuildParams{
				Rule:    common.CCRule,
				Outputs: []string{obj},
				Inputs:  []string{in},
				Depfile: dep,
				Deps:    blueprint.DepsGCC,
				Args: map[string]string{
					"cc":      cc,
					"cflags":  strings.Join(flags, " "),
					"depfile": dep,
					"incs":    incPathsToOpts(buildInfo.Incs),
				},
			},
		)
	}

	return objs
}
