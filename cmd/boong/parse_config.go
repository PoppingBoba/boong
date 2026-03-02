package main

import (
	"flag"
	"path/filepath"

	"github.com/PoppingBoba/boong/build"
)

func parseAndGetConfig() build.Config {
	var src string
	var out string
	var compiler_c string
	var compiler_path string
	var arch string

	flag.StringVar(&src, "src", ".", "Boong's root path")
	flag.StringVar(&out, "out", "out", "Boong's output directory")
	flag.StringVar(&compiler_c, "compiler", "clang", "Boong's used compiler")
	flag.StringVar(&compiler_path, "path", "", "Boong's compiler path")
	flag.StringVar(&arch, "arch", "", "Build Architecture")
	flag.Parse()

	srcDir, err := filepath.Abs(src)
	checkFail(err)

	outDir, err := filepath.Abs(out)
	checkFail(err)

	relSrcPath, err := filepath.Rel(outDir, srcDir)
	checkFail(err)

	cfg := build.Config{
		SrcPath:       srcDir,
		OutPath:       outDir,
		RelSrcPath:    relSrcPath,
		CCompiler:     compiler_c,
		CCompilerPath: compiler_path,
		Arch:          arch,
	}

	return cfg
}
