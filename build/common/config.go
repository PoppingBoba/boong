package common

import (
	"io/fs"
	"os"
	"path/filepath"

	"github.com/google/blueprint"
)

type ConfigBase struct {
	// Main path of source code
	SrcPath string

	// Main path of output directory
	OutPath string

	// Rel path for Src & Inc path
	RelSrcPath string

	// Target Architecture
	// (like x86_64 or arm64)
	Arch string
}

func (c *ConfigBase) CreateOutPath() error {
	return os.MkdirAll(c.OutPath, 0o755)
}

func (c *ConfigBase) SetSrcPath(pctx *blueprint.Context) {
	pctx.SetSrcDir(c.SrcPath)
}

func (c *ConfigBase) CreateBuildNinja() (*os.File, error) {
	ninjaPath := filepath.Join(c.OutPath, "build.ninja")
	return os.Create(ninjaPath)
}

func (c *ConfigBase) SearchBuildFiles() ([]string, error) {
	var bpFiles []string
	var retErr error

	retErr = filepath.WalkDir(c.SrcPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() && filepath.Base(path) == "Build.bp" {
			bpFiles = append(bpFiles, path)
		}

		return nil
	})

	return bpFiles, retErr
}

func (c *ConfigBase) GetRelIncPath(pctx blueprint.ModuleContext, inIncs []string) []string {
	var incs []string

	for _, inc := range inIncs {
		relInc := filepath.Join(c.RelSrcPath, pctx.ModuleDir(), inc)
		incs = append(incs, relInc)
	}

	return incs
}
