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
		Srcs        []string
		Cflags      []string
		Ldflags     []string
		Static_libs []string
	}
}

func (m *CBinary) setRules(ctx blueprint.ModuleContext, compilers Compilers) {
	cfg := ctx.Config().(Config)

	objs := cfg.AddCompileObjects(ctx, m.Properties.Srcs, m.Properties.Cflags, compilers)

	var libs []string

	ctx.VisitDepsDepthFirst(func(m blueprint.Module) {
		if l, ok := m.(*CLibraryStatic); ok && l.outLib != "" {
			libs = append(libs, l.outLib)
		}
	})

	out := filepath.Join("bin", ctx.ModuleName())
	ctx.Build(pkgCtx, blueprint.BuildParams{
		Rule:      LinkRule,
		Outputs:   []string{out},
		Inputs:    objs,
		Default:   true,
		Implicits: libs,
		Args: map[string]string{
			"cc":      compilers.CXX,
			"ldflags": strings.Join(m.Properties.Ldflags, " "),
			"libs":    strings.Join(libs, " "),
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
