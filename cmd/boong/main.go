package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/PoppingBoba/boong/build"
	"github.com/google/blueprint"
)

func main() {
	var src string
	var out string
	var compiler_c string
	var compiler_path string

	flag.StringVar(&src, "src", ".", "Boong's root path")
	flag.StringVar(&out, "out", "out", "Boong's output directory")
	flag.StringVar(&compiler_c, "compiler", "clang", "Boong's used compiler")
	flag.StringVar(&compiler_path, "path", "", "Boong's compiler path")
	flag.Parse()

	srcDir, err := filepath.Abs(src)
	checkFail(err)

	outDir, err := filepath.Abs(out)
	checkFail(err)

	checkFail(os.MkdirAll(outDir, 0o755))

	relToSrc, err := filepath.Rel(outDir, srcDir)
	checkFail(err)

	cfg := build.Config{
		SrcPath:       srcDir,
		OutPath:       outDir,
		RelToSrcPath:  relToSrc,
		CCompiler:     compiler_c,
		CCompilerPath: compiler_path,
	}

	bpCtx := blueprint.NewContext()
	bpCtx.SetSrcDir(srcDir)

	RegisterBoongModule(bpCtx)

	_, parseErrs := bpCtx.ParseFileList(".", []string{"Build.bp"}, cfg)
	checkFailMany(parseErrs)

	_, depErrs := bpCtx.ResolveDependencies(cfg)
	checkFailMany(depErrs)

	_, actErrs := bpCtx.PrepareBuildActions(cfg)
	checkFailMany(actErrs)

	ninjaPath := filepath.Join(outDir, "build.ninja")
	f, err := os.Create(ninjaPath)
	checkFail(err)
	defer f.Close()

	bw := bufio.NewWriter(f)
	checkFail(bpCtx.WriteBuildFile(bw, false, "build.ninja"))
	checkFail(bw.Flush())

	fmt.Printf("Boooong Run Done : %s", ninjaPath)

}
