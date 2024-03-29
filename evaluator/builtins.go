package evaluator

import (
    "fmt"
    "monkey/object"
)

var builtins = map[string]*object.Builtin{
	"len": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			switch arg := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			case *object.Array:
                return &object.Integer{Value: int64(len(arg.Values))}
			default:
				return newError("argument to `len` not supported, got %s", args[0].Type())
			}
		},
	},
    "print": {
        Fn: func(args ...object.Object) object.Object {
            if len(args) != 1 {
                return newError("wrong number of arguments. got=%d, want=1", len(args))
            }

            fmt.Printf("%s\n", args[0].Inspect())

            return NULL
        },
    },
}
