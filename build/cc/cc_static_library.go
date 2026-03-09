package cc

import (
	"fmt"
	"path/filepath"

	"github.com/PoppingBoba/boong/build/common"
	"github.com/google/blueprint"
)

type CLibraryStatic struct {
	blueprint.SimpleName
	Properties struct {
		Srcs                []string
		Cflags              []string
		Static_libs         []string
		Export_include_dirs []string
		Local_include_dirs  []string
	}

	outLib            string
	objs              []string
	exportIncludeDirs []string
}

type staticLibTag struct {
	blueprint.BaseDependencyTag
}

func (l *CLibraryStatic) setRules(ctx blueprint.ModuleContext, compilers Compilers) {
	cfg := ctx.Config().(Config)

	var buildInfo BuildInfo

	buildInfo.Srcs = l.Properties.Srcs
	buildInfo.Incs = cfg.GetRelIncPath(ctx, l.Properties.Export_include_dirs)
	buildInfo.Cflags = l.Properties.Cflags
	buildInfo.Compilers = compilers

	ctx.VisitDepsDepthFirst(func(m blueprint.Module) {
		fmt.Println("FUCK")
		if l, ok := m.(*CLibraryStatic); ok {
			if l.outLib != "" {
				buildInfo.Libs = append(buildInfo.Libs, l.outLib)
			}
			if len(l.exportIncludeDirs) > 0 {
				buildInfo.Incs = append(buildInfo.Incs, l.exportIncludeDirs...)
			}
		}
	})

	// Pass the
	if len(buildInfo.Srcs) > 0 {
		objs := cfg.AddCompileObjects(ctx, buildInfo)

		out := filepath.Join("lib", ctx.ModuleName()+".a")
		ctx.Build(
			common.PkgCtx,
			blueprint.BuildParams{
				Rule:    common.LibRule,
				Outputs: []string{out},
				Inputs:  objs,
				Default: false,
				Args: map[string]string{
					"arcmd": "ar",
				},
			},
		)
		l.outLib = out
	}

	l.exportIncludeDirs = buildInfo.Incs
}

func (l *CLibraryStatic) GenerateBuildActions(ctx blueprint.ModuleContext) {
	cfg := ctx.Config().(Config)

	compilers := setCompiler(cfg)
	l.setRules(ctx, *compilers)
}

func (l *CLibraryStatic) String() string {
	return l.Name()
}

func (l *CLibraryStatic) LibraryFileName() string {
	return l.outLib
}

func GetCLibraryStatic() (blueprint.Module, []interface{}) {
	m := &CLibraryStatic{}
	return m, []interface{}{&m.SimpleName.Properties, &m.Properties}
}
