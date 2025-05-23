package beherit

import "github.com/expr-lang/expr/vm"

type Set struct {
	Path  string `yaml:"path"`
	Value string `yaml:"value"`
}

type TriggerCommand struct {
	Command      string         `yaml:"command"`
	Vars         map[string]any `yaml:"vars"`      // 1
	Set          []Set          `yaml:"set"`       // 2
	ConditionStr string         `yaml:"condition"` // 3
	condition    *vm.Program    `yaml:"-"`         //
	Mapping      map[string]any `yaml:"mapping"`   // 4
}

type Config struct {
	Env      map[string]any               `yaml:"env"`
	Triggers map[string][]*TriggerCommand `yaml:"triggers"`
}
