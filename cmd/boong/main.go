package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"path/filepath"

	"github.com/google/blueprint"
)

func main() {
	cfg := parseAndGetConfig()

	checkFail(cfg.CreateOutPath())

	bpCtx := blueprint.NewContext()
	cfg.SetSrcPath(bpCtx)

	RegisterBoongModule(bpCtx)

	var bpFiles []string
	_ = filepath.WalkDir(cfg.SrcPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() && filepath.Base(path) == "Build.bp" {
			bpFiles = append(bpFiles, path)
		}

		return nil
	})

	for index, bpFile := range bpFiles {
		fmt.Printf("[%d] List of bp file : %s\n", index, bpFile)
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
