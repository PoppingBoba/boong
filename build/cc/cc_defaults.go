package cc

import "github.com/google/blueprint"

type CDefaults struct {
	blueprint.SimpleName
	Properties struct {
		Cflags []string
	}

	outCflags []string
}

func (c *CDefaults) GenerateBuildActions(ctx blueprint.ModuleContext) {

}

func (c *CDefaults) String() string {
	return c.Name()
}

type defaultTag struct {
	blueprint.BaseDependencyTag
}

func GetCDefault() (blueprint.Module, []interface{}) {
	m := &CDefaults{}
	return m, []interface{}{&m.SimpleName.Properties, &m.Properties}
}
