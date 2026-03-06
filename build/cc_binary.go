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
		Srcs               []string
		Cflags             []string
		Ldflags            []string
		Static_libs        []string
		Local_include_dirs []string
	}
}

func incArrToIncOpt(pctx blueprint.ModuleContext, cfg Config, incs []string) []string {
	var ret []string

	for _, v := range incs {
		path := filepath.Join(cfg.RelSrcPath, pctx.ModuleDir(), v)
		ret = append(ret, "-I"+path)
	}

	return ret
}

func (m *CBinary) setRules(ctx blueprint.ModuleContext, compilers Compilers) {
	cfg := ctx.Config().(Config)

	var buildInfo BuildInfo

	ctx.VisitDepsDepthFirst(func(m blueprint.Module) {
		if l, ok := m.(*CLibraryStatic); ok && l.outLib != "" {
			buildInfo.Libs = append(buildInfo.Libs, l.outLib)
			buildInfo.Incs = append(buildInfo.Incs, l.exportIncludeDirs...)
		}
	})

	objs := cfg.AddCompileObjects(ctx, buildInfo)

	out := filepath.Join("bin", ctx.ModuleName())
	ctx.Build(pkgCtx, blueprint.BuildParams{
		Rule:      LinkRule,
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
