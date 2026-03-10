package cc

import (
	"path/filepath"
	"strings"

	"github.com/PoppingBoba/boong/build/common"
	"github.com/google/blueprint"
)

type Compilers struct {
	CC  string
	CXX string
}

// For C Binary
type CBinary struct {
	blueprint.SimpleName
	Properties struct {
		// C/C++ Sources
		Srcs []string
		// C/C++ common compiler flags
		Cflags []string
		// Only for C++ extensive compiler flags
		Cppflags []string
		// Global C/C++ flags
		Defaults []string
		// Static Librareies
		Ldflags []string
		// Export Include Directories for dependency
		Static_libs []string
		// Local include directories
		Local_include_dirs []string
	}
}

func setCompiler(config Config) *Compilers {
	compilers := &Compilers{}

	if config.CCompiler == "" || config.CCompiler == "clang" {
		compilers.CC = "clang"
		compilers.CXX = "clang++"
	}

	if config.CCompiler == "gcc" {
		compilers.CC = "gcc"
		compilers.CXX = "g++"
	}

	if config.CCompilerPath != "" {
		compiler_path := config.CCompilerPath + "/"
		compilers.CC = compiler_path + compilers.CC
		compilers.CXX = compiler_path + compilers.CXX
	}

	return compilers
}

func (m *CBinary) setRules(ctx blueprint.ModuleContext, compilers Compilers) {
	cfg := ctx.Config().(Config)

	var buildInfo BuildInfo

	buildInfo.Srcs = m.Properties.Srcs
	buildInfo.Incs = cfg.GetRelIncPath(ctx, m.Properties.Local_include_dirs)
	buildInfo.Cflags = m.Properties.Cflags
	buildInfo.Cppflags = m.Properties.Cppflags
	buildInfo.Compilers = compilers

	ctx.VisitDepsDepthFirst(func(m blueprint.Module) {
		if d, ok := m.(*CDefaults); ok {
			if len(d.outCflags) > 0 {
				buildInfo.Cflags = append(buildInfo.Cflags, d.outCflags...)
			}

			if len(d.outCppflags) > 0 {
				buildInfo.Cppflags = append(buildInfo.Cppflags, d.outCppflags...)
			}
		}

		if l, ok := m.(*CLibraryStatic); ok {
			if l.outLib != "" {
				buildInfo.Libs = append(buildInfo.Libs, l.outLib)
			}
			if len(l.exportIncludeDirs) > 0 {
				buildInfo.Incs = append(buildInfo.Incs, l.exportIncludeDirs...)
			}
		}
	})

	objs := cfg.AddCompileObjects(ctx, buildInfo)

	out := filepath.Join("bin", ctx.ModuleName())
	ctx.Build(common.PkgCtx, blueprint.BuildParams{
		Rule:      common.LinkRule,
		Outputs:   []string{out},
		Inputs:    objs,
		Default:   true,
		Implicits: buildInfo.Libs,
		Args: map[string]string{
			"cc":      compilers.CXX,
			"ldflags": strings.Join(m.Properties.Ldflags, " "),
			"libs":    strings.Join(buildInfo.Libs, " "),
		},
	})
}

func (c *CBinary) GenerateBuildActions(ctx blueprint.ModuleContext) {
	cfg := ctx.Config().(Config)

	compilers := setCompiler(cfg)
	c.setRules(ctx, *compilers)
}

func (c *CBinary) String() string {
	return c.Name()
}

func GetCBinary() (blueprint.Module, []interface{}) {
	m := &CBinary{}
	return m, []interface{}{&m.SimpleName.Properties, &m.Properties}
}
