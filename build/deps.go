package build

import "github.com/google/blueprint"

func CDepsMutator(ctx blueprint.BottomUpMutatorContext) {
	switch m := ctx.Module().(type) {
	case *CBinary:
		if len(m.Properties.Static_libs) > 0 {
			ctx.AddDependency(
				m,
				staticLibTag{},
				m.Properties.Static_libs...,
			)
		}
	case *CLibraryStatic:
		if len(m.Properties.Static_libs) > 0 {
			ctx.AddDependency(
				m,
				staticLibTag{},
				m.Properties.Static_libs...,
			)
		}
	}
}
