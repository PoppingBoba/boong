package build

import (
	"path/filepath"

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

	objs := cfg.AddCompileObjects(ctx, l.Properties.Srcs, l.Properties.Cflags, compilers)

	out := filepath.Join("lib", ctx.ModuleName()+".a")
	ctx.Build(
		pkgCtx,
		blueprint.BuildParams{
			Rule:    LibRule,
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
