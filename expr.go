package beherit

import (
	"fmt"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/vm"
)

type ExprMatcher interface {
	Match(key string) bool
}

type ExpressionCompiler struct {
	matcher  ExprMatcher
	programs map[string]*vm.Program
}

func NewExpressionCompiler(m ExprMatcher) *ExpressionCompiler {
	return &ExpressionCompiler{
		matcher:  m,
		programs: make(map[string]*vm.Program),
	}
}

func (ec *ExpressionCompiler) Compile(s string) (*vm.Program, error) {
	if p, ok := ec.programs[s]; ok {
		return p, nil
	}
	program, err := expr.Compile(s)
	if err != nil {
		return nil, fmt.Errorf("could not compile expression: %w", err)
	}
	ec.programs[s] = program
	return program, nil
}

func (ec *ExpressionCompiler) CompileMap(params map[string]any) (map[string]any, error) {
	result := make(map[string]any)
	for key, value := range params {
		switch v := value.(type) {
		case []any:
			s := make([]any, len(v))
			for i, item := range v {
				switch item := item.(type) {
				case map[string]any:
					r, err := ec.CompileMap(item)
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
			r, err := ec.CompileMap(v)
			if err != nil {
				return nil, err
			}
			result[key] = r
		case string: // check if string is an expression
			if !ec.matcher.Match(key) {
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
