package build

import (
	"path/filepath"
	"strings"

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
		Srcs    []string
		Cflags  []string
		Ldflags []string
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

	objs := cfg.AddCompileObjects(ctx, m.Properties.Srcs, m.Properties.Cflags, compilers)

	out := filepath.Join("bin", ctx.ModuleName())
	ctx.Build(pkgCtx, blueprint.BuildParams{
		Rule:    LinkRule,
		Outputs: []string{out},
		Inputs:  objs,
		Default: true,
		Args: map[string]string{
			"cc":      compilers.CXX,
			"ldflags": strings.Join(m.Properties.Ldflags, " "),
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
