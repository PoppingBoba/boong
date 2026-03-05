package main

import (
	"bufio"
	"fmt"

	"github.com/google/blueprint"
)

func main() {
	cfg := parseAndGetConfig()

	checkFail(cfg.CreateOutPath())

	bpCtx := blueprint.NewContext()
	cfg.SetSrcPath(bpCtx)

	RegisterBoongModule(bpCtx)

	var bpFiles []string
	bpFiles, searchErr := cfg.SearchBuildFiles()
	checkFail(searchErr)

	if len(bpFiles) == 0 {
		panic("[X] Failed to Get BP Files")
	}

	_, parseErrs := bpCtx.ParseFileList(cfg.SrcPath, bpFiles, cfg)
	checkFailMany(parseErrs)

	_, actErrs := bpCtx.PrepareBuildActions(cfg)
	checkFailMany(actErrs)

	f, err := cfg.CreateBuildNinja()
	checkFail(err)
	defer f.Close()

	bw := bufio.NewWriter(f)
	checkFail(bpCtx.WriteBuildFile(bw, false, "build.ninja"))
	checkFail(bw.Flush())

	fmt.Printf("Boooong Run Done : %s\n", cfg.OutPath)
}
