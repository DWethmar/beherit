package beherit_test

import (
	"regexp"
	"testing"

	"github.com/dwethmar/beherit"
	"github.com/expr-lang/expr/vm"
	"github.com/google/go-cmp/cmp"
)

var MatchByEndingDollarRegex = regexp.MustCompile(`\$$`)

func TestCompileMapExpressions(t *testing.T) {
	t.Run("compile param expressions", func(t *testing.T) {
		ec := beherit.NewExpressionCompiler(&beherit.RegexExprMatcher{
			KeyRegex: MatchByEndingDollarRegex,
		})
		p, err := ec.Compile("1 + 1")
		if err != nil {
			t.Error(err)
		}
		if p == nil {
			t.Error("expected a program, got nil")
		}

		p2, err := ec.Compile("1 + 1")
		if err != nil {
			t.Error(err)
		}
		if p != p2 {
			t.Error("expected the same program")
		}
	})
	t.Run("compile param expressions", func(t *testing.T) {
		ec := beherit.NewExpressionCompiler(&beherit.RegexExprMatcher{
			KeyRegex: MatchByEndingDollarRegex,
		})
		programs := map[string]any{
			"Position": map[string]any{
				"X$": `let x = 1; 1 + 1 + x`,
				"Y":  23,
			},
			"Test": map[string]any{
				"Name$": "1 + 1",
			},
			"Test2": map[string]any{
				"Name$": "1 + 1",
			},
		}
		got, err := ec.CompileMap(programs)
		if err != nil {
			t.Error(err)
		}
		// check if component.Position.X is of type *vm.Program
		v1, ok := got["Position"].(map[string]any)["X"].(*vm.Program)
		if !ok {
			t.Errorf("expected type *vm.Program, got %T", v1)
		}
		// check if component.Position.Y is of type int
		v2, ok := got["Position"].(map[string]any)["Y"].(int)
		if !ok {
			t.Errorf("expected type int, got %T", v2)
		}
		// check if component.Test.Name is of type *vm.Program
		v3, ok := got["Test"].(map[string]any)["Name"].(*vm.Program)
		if !ok {
			t.Errorf("expected type *vm.Program, got %T", v3)
		}
		// check if component.Test2.Name is of type *vm.Program
		v4, ok := got["Test2"].(map[string]any)["Name"].(*vm.Program)
		if !ok {
			t.Errorf("expected type *vm.Program, got %T", v4)
		}
		if v3 != v4 {
			t.Error("expected component.Test.Name to be equal to component.Test2.Name")
		}
	})
}

func TestRunParamExpressions(t *testing.T) {
	t.Run("run param expressions", func(t *testing.T) {
		ec := beherit.NewExpressionCompiler(&beherit.RegexExprMatcher{
			KeyRegex: MatchByEndingDollarRegex,
		})
		params := map[string]any{
			"Component": map[string]any{
				"Position": map[string]any{
					"X$": `let x = 1; a + b + x`,
					"Y":  23,
				},
				"Test": map[string]any{
					"Name$": "string(a * b)+'name'",
				},
			},
		}
		compiled, err := ec.CompileMap(params)
		if err != nil {
			t.Error(err)
		}
		env := map[string]any{
			"a": 2,
			"b": 3,
		}
		result, err := ec.Run(compiled, env)
		if err != nil {
			t.Error(err)
		}
		expected := map[string]any{
			"Component": map[string]any{
				"Position": map[string]any{
					"X": 6,
					"Y": 23,
				},
				"Test": map[string]any{
					"Name": "6name",
				},
			},
		}
		if diff := cmp.Diff(expected, result); diff != "" {
			t.Errorf("unexpected result (-want +got):\n%s", diff)
		}
	})

	t.Run("run param expressions with nested expressions", func(t *testing.T) {
		ec := beherit.NewExpressionCompiler(&beherit.RegexExprMatcher{
			KeyRegex: MatchByEndingDollarRegex,
		})
		params := map[string]any{
			"Component": map[string]any{
				"Position": map[string]any{
					"X$": `let x = 1; a + b + x`,
					"Y":  23,
				},
				"Test": map[string]any{
					"Name$": "string(a * b)+'name'",
				},
			},
		}
		compiled, err := ec.CompileMap(params)
		if err != nil {
			t.Error(err)
		}
		env := map[string]any{
			"a": 2,
			"b": 3,
		}
		result, err := ec.Run(compiled, env)
		if err != nil {
			t.Error(err)
		}
		expected := map[string]any{
			"Component": map[string]any{
				"Position": map[string]any{
					"X": 6,
					"Y": 23,
				},
				"Test": map[string]any{
					"Name": "6name",
				},
			},
		}
		if diff := cmp.Diff(expected, result); diff != "" {
			t.Errorf("unexpected result (-want +got):\n%s", diff)
		}
	})
}
