package builtins

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/EVFUBS/AlphaLang/objects"
)

var BuiltIns = map[string]objects.Builtin{
	"len": {
		Fn: func(args ...objects.Object) objects.Object {
			if len(args) != 1 {
				return objects.NewError("wrong number of arguments. got=%d, want=1", len(args))
			}

			switch arg := args[0].(type) {
			case *objects.String:
				return &objects.Integer{Value: int64(len(arg.Value))}
			case *objects.Array:
				return &objects.Integer{Value: int64(len(arg.Elements))}
			default:
				return objects.NewError("argument to `len` not supported, got %s", args[0].Type())
			}
		},
	},
	"append": {
		Fn: func(args ...objects.Object) objects.Object {
			if len(args) == 2 {
				//append array
				if arg, ok := args[0].(*objects.Array); ok {
					arg.Elements = append(arg.Elements, args[1])
				}
				return nil
			} else if len(args) == 3 {
				//append hash
				if arg, ok := args[0].(*objects.Hash); ok {
					key, ok := args[1].(objects.Hashable)
					if !ok {
						return objects.NewError("unusable as hash key: %s", args[1].Type())
					}
					hashKey := key.HashKey()
					arg.Pairs[hashKey] = objects.HashPair{Key: args[1], Value: args[2]}
				}
				return nil
			} else {
				return objects.NewError("arguement to `append` not supported, got %s, %s", args[0].Type(), args[1].Type())
			}
		},
	},
	"pop": {
		Fn: func(args ...objects.Object) objects.Object {
			if len(args) != 1 {
				return objects.NewError("wrong number of arguments. got=%d, want=1", len(args))
			}
			if arg, ok := args[0].(*objects.Array); ok {
				last := arg.Elements[len(arg.Elements)-1]
				arg.Elements = arg.Elements[:len(arg.Elements)-1]
				return last
			} else {
				return objects.NewError("argument to `pop` not supported, got %s", args[0].Type())
			}
		},
	},

	"println": {
		Fn: func(args ...objects.Object) objects.Object {
			for _, arg := range args {
				fmt.Println(arg.Inspect())
			}
			return nil
		},
	},
	"print": {
		Fn: func(args ...objects.Object) objects.Object {
			for _, arg := range args {
				fmt.Print(arg.Inspect())
			}
			return nil
		},
	},

	"int": {
		Fn: func(args ...objects.Object) objects.Object {
			if len(args) != 1 {
				return objects.NewError("wrong number of arguments. got=%d, want=1", len(args))
			}
			switch arg := args[0].(type) {
			case *objects.Integer:
				return arg
			case *objects.Boolean:
				if arg.Value {
					return &objects.Integer{Value: 1}
				}
				return &objects.Integer{Value: 0}
			case *objects.String:
				value, err := strconv.ParseInt(arg.Value, 0, 64)
				if err != nil {
					return objects.NewError("could not convert %q to integer", arg.Value)
				}
				return &objects.Integer{Value: value}
			default:
				return objects.NewError("argument to `int` not supported, got %s", args[0].Type())
			}
		},
	},

	"input": {
		Fn: func(args ...objects.Object) objects.Object {
			if len(args) != 1 {
				return objects.NewError("wrong number of arguments. got=%d, want=1", len(args))
			}
			switch arg := args[0].(type) {
			case *objects.String:
				fmt.Println(arg.Value)
				reader := bufio.NewReader(os.Stdin)
				input, err := reader.ReadString('\n')
				if err != nil {
					fmt.Println(err)
				}
				input = strings.TrimSpace(input)
				return &objects.String{Value: input}
			default:
				return objects.NewError("argument to `input` not supported, got %s", args[0].Type())
			}
		},
	},

	"rand": {
		Fn: func(args ...objects.Object) objects.Object {
			//rand rumber between two integers
			source := rand.NewSource(time.Now().UnixNano())
			rand := rand.New(source)
			if len(args) == 2 {
				if arg1, ok := args[0].(*objects.Integer); ok {
					if arg2, ok := args[1].(*objects.Integer); ok {
						return &objects.Integer{Value: int64(rand.Intn(int(arg2.Value-arg1.Value)) + int(arg1.Value))}
					}
				}
			} else if len(args) == 1 {
				if arg1, ok := args[0].(*objects.Integer); ok {
					return &objects.Integer{Value: int64(rand.Intn(int(arg1.Value)))}
				}
			}
			return objects.NewError("argument to `rand` not supported, got %s, %s", args[0].Type(), args[1].Type())
		},
	},
}
