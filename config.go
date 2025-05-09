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
	Name        string           `yaml:"name" expr:"name"`
	Description string           `yaml:"description" expr:"description"`
	Components  []map[string]any `yaml:"components" expr:"components"`
}

type Config struct {
	Blueprints []*Blueprint                 `yaml:"blueprints"`
	Triggers   map[string][]*TriggerCommand `yaml:"triggers"`
}
