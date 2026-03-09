package main

import (
	"github.com/PoppingBoba/boong/build/cc"
	"github.com/google/blueprint"
)

func RegisterBoongModule(ctx *blueprint.Context) {
	ctx.RegisterModuleType("cc_binary", cc.GetCBinary)
	ctx.RegisterModuleType("cc_library_static", cc.GetCLibraryStatic)
	ctx.RegisterBottomUpMutator("cc_deps", cc.CDepsMutator)
}
