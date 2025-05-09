package beherit

import "github.com/expr-lang/expr/vm"

type TriggerCommand struct {
	Command      string         `yaml:"command"`
	Vars         map[string]any `yaml:"vars"`
	ConditionStr string         `yaml:"condition"`
	condition    *vm.Program    `yaml:"-"`
	Params       map[string]any `yaml:"params"`
}

type Blueprint struct {
	Name        string           `yaml:"name"`
	Description string           `yaml:"description"`
	Components  []map[string]any `yaml:"components"`
}

type Config struct {
	Blueprints []*Blueprint                 `yaml:"blueprints"`
	Triggers   map[string][]*TriggerCommand `yaml:"triggers"`
}
