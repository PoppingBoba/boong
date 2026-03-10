package common

import "github.com/google/blueprint"

var CCRule = PkgCtx.StaticRule(
	"cc",
	blueprint.RuleParams{
		Command:     "mkdir -p $$(dirname $out) && $cc -MMD -MF $depfile -c $cflags -o $out $in $incs",
		Description: "CC $out",
	},
	"cc", "cflags", "depfile", "incs",
)

var LinkRule = PkgCtx.StaticRule(
	"link",
	blueprint.RuleParams{
		Command:     "mkdir -p $$(dirname $out) && $cc $ldflags -o $out $in $libs",
		Description: "LINK $out",
	},
	"cc", "ldflags", "libs",
)

var LibRule = PkgCtx.StaticRule(
	"ar",
	blueprint.RuleParams{
		Command:     "$arcmd crs $out $in",
		Description: "LIB $out",
	},
	"arcmd",
)
