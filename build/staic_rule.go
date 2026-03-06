package build

import "github.com/google/blueprint"

var CCRule = pkgCtx.StaticRule(
	"cc",
	blueprint.RuleParams{
		Command:     "mkdir -p $$(dirname $out) && $cc -MMD -MF $depfile -c $cflags -o $out $in $incs",
		Description: "CC $out",
	},
	"cc", "cflags", "depfile", "incs",
)

var LinkRule = pkgCtx.StaticRule(
	"link",
	blueprint.RuleParams{
		Command:     "mkdir -p $$(dirname $out) && $cc $ldflags -o $out $in $libs",
		Description: "LINK $out",
	},
	"cc", "ldflags", "libs",
)

var LibRule = pkgCtx.StaticRule(
	"ar",
	blueprint.RuleParams{
		Command:     "$arcmd crs $out $in",
		Description: "LIB $out",
	},
	"arcmd",
)
