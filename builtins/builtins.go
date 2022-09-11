package builtins

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
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

	"input": {
		Fn: func(args ...objects.Object) objects.Object {
			if len(args) == 0 || len(args) == 1 {
				if len(args) == 1 {
					if arg, ok := args[0].(*objects.String); ok {
						println(arg.Value)
					}
				}
				inputReader := bufio.NewReader(os.Stdin)
				input, _ := inputReader.ReadString('\n')
				// check if string is int, float, string
				if i, err := strconv.Atoi(input); err == nil {
					return &objects.Integer{Value: int64(i)}
				} else if f, err := strconv.ParseFloat(input, 64); err == nil {
					return &objects.Float{Value: float64(f)}
				} else {
					return &objects.String{Value: input}
				}
			} else {
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
