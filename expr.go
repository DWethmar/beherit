package beherit

import (
	"fmt"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/vm"
)

type ExpressionCompiler struct {
	programs map[string]*vm.Program
}

func NewExpressionCompiler() *ExpressionCompiler {
	return &ExpressionCompiler{
		programs: make(map[string]*vm.Program),
	}
}

func (ec *ExpressionCompiler) Compile(params map[string]any) (map[string]any, error) {
	result := make(map[string]any)
	for key, value := range params {
		switch v := value.(type) {
		case []any:
			s := make([]any, len(v))
			for i, item := range v {
				switch item := item.(type) {
				case map[string]any:
					r, err := ec.Compile(item)
					if err != nil {
						return nil, err
					}
					s[i] = r
				default:
					s[i] = item
				}
			}
			result[key] = s
		case map[string]any:
			r, err := ec.Compile(v)
			if err != nil {
				return nil, err
			}
			result[key] = r
		case string: // check if string is an expression
			if key[len(key)-1] != '$' {
				result[key] = v
				continue
			}
			var program *vm.Program
			var err error
			if p, ok := ec.programs[v]; ok { // check if expression is already compiled
				program = p
			} else {
				program, err = expr.Compile(v)
				if err != nil {
					return nil, fmt.Errorf("could not compile expression: %w", err)
				}
				ec.programs[v] = program
			}
			result[key[:len(key)-1]] = program
		default:
			result[key] = v
		}
	}
	return result, nil
}

func (ec *ExpressionCompiler) Run(params map[string]any, env any) (map[string]any, error) {
	result := make(map[string]any)
	for key, value := range params {
		switch v := value.(type) {
		case []any:
			s := make([]any, len(v))
			for i, item := range v {
				switch item := item.(type) {
				case map[string]any:
					r, err := ec.Run(item, env)
					if err != nil {
						return nil, err
					}
					s[i] = r
				default:
					s[i] = item
				}
			}
			result[key] = s
		case map[string]any:
			r, err := ec.Run(v, env)
			if err != nil {
				return nil, err
			}
			result[key] = r
		case *vm.Program:
			r, err := expr.Run(v, env)
			if err != nil {
				return nil, err
			}
			result[key] = r
		default:
			result[key] = v
		}
	}
	return result, nil
}
