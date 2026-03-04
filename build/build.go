package build

import "github.com/google/blueprint"

// Package Context
var pkgCtx = blueprint.NewPackageContext("github.com/PoppingBoba/boong/build")

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
