package build

import (
	"path/filepath"
	"strings"

	"github.com/google/blueprint"
)

type Config struct {
	SrcPath       string
	OutPath       string
	RelToSrcPath  string
	CCompilerPath string
	CCompiler     string
}

type Compilers struct {
	CC  string
	CXX string
}

// Package Context
var pkgCtx = blueprint.NewPackageContext("github.com/PoppingBoba/boong/build")

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

	// For C/C++
	ccRule := ctx.Rule(
		pkgCtx,
		"cc",
		blueprint.RuleParams{
			Command:     "mkdir -p $$(dirname $out) && $cc -MMD -MF $depfile -c $cflags -o $out $in",
			Description: "CC $out",
		},
		"cc", "cflags", "depfile",
	)

	// For Linking
	linkRule := ctx.Rule(
		pkgCtx,
		"link",
		blueprint.RuleParams{
			Command:     "mkdir -p $$(dirname $out) && $cc $ldflags -o $out $in",
			Description: "LINK $out",
		},
		"cc", "ldflags",
	)

	var objs []string
	for _, src := range m.Properties.Srcs {
		in := filepath.Join(cfg.RelToSrcPath, ctx.ModuleDir(), src)

		base := filepath.Base(src)
		obj := filepath.Join("obj", ctx.ModuleName(), base+".o")
		dep := obj + ".d"
		objs = append(objs, obj)

		var cc string
		file_ext := filepath.Ext(src)
		if file_ext == ".cpp" || file_ext == ".cc" {
			cc = compilers.CXX
		} else {
			cc = compilers.CC
		}

		ctx.Build(
			pkgCtx,
			blueprint.BuildParams{
				Rule:    ccRule,
				Outputs: []string{obj},
				Inputs:  []string{in},
				Depfile: dep,
				Deps:    blueprint.DepsGCC,
				Args: map[string]string{
					"cc":      cc,
					"cflags":  strings.Join(m.Properties.Cflags, " "),
					"depfile": dep,
				},
			},
		)
	}

	out := filepath.Join("bin", ctx.ModuleName())
	ctx.Build(pkgCtx, blueprint.BuildParams{
		Rule:    linkRule,
		Outputs: []string{out},
		Inputs:  objs,
		Default: true,
		Args: map[string]string{
			"cc":      compilers.CC,
			"ldflags": strings.Join(m.Properties.Ldflags, " "),
		},
	})
}

func (c *CBinary) GenerateBuildActions(ctx blueprint.ModuleContext) {
	cfg := ctx.Config().(Config)

	compilers := setCompiler(cfg)
	c.setRules(ctx, *compilers)
}

func (c *CBinary) Name() string {
	return "cc_binary"
}

func (c *CBinary) String() string {
	return c.Name()
}

func GetCBinary() (blueprint.Module, []interface{}) {
	m := &CBinary{}
	return m, []interface{}{&m.SimpleName.Properties, &m.Properties}
}
