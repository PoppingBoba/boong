package cc

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
		if len(m.Properties.Defaults) > 0 {
			ctx.AddDependency(
				m,
				defaultTag{},
				m.Properties.Defaults...,
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
		if len(m.Properties.Defaults) > 0 {
			ctx.AddDependency(
				m,
				defaultTag{},
				m.Properties.Defaults...,
			)
		}
	}
}
