package main

import (
	"github.com/PoppingBoba/boong/build"
	"github.com/google/blueprint"
)

func RegisterBoongModule(ctx *blueprint.Context) {
	ctx.RegisterModuleType("cc_binary", build.GetCBinary)
	ctx.RegisterModuleType("cc_library_static", build.GetCLibraryStatic)
	ctx.RegisterBottomUpMutator("cc_deps", build.CDepsMutator)
}
