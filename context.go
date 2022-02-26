package meshreflow

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var rCmd = regexp.MustCompile(`^[a-z]+ `)
var rArg = regexp.MustCompile(` \{[a-z]+:num\}`)
var rNum = regexp.MustCompile(`\s*[-]?\d[\d,]*[\.]?[\d{2}]*`)

type RExp = *regexp.Regexp

type cmdMatcher struct {
	name        string
	argMatchers []argMatcher
}

type argMatcher struct {
	name    string
	matcher RExp
}

type parsedCommand struct {
	command        string
	commandName    string
	commandMatcher cmdMatcher
	argNames       []string
	argValues      []string
}

type ArgConverter struct{}

type Context struct {
	commandMatchers map[string][]cmdMatcher
	commands        []string
	meshPoints      []float64
}

func NewContext() *Context {
	return &Context{
		commandMatchers: map[string][]cmdMatcher{},
		commands:        make([]string, 0),
		meshPoints:      make([]float64, 0),
	}
}

func (ctx *Context) PushCmd(cmd string) (parsedCommand, error) {
	cmdMatchIdx := rCmd.FindStringIndex(cmd)
	if cmdMatchIdx == nil {
		return parsedCommand{}, fmt.Errorf("faild to parse command name from command '%s'", cmd)
	}
	cmdName := cmd[cmdMatchIdx[0] : cmdMatchIdx[1]-1]
	cmdMatchers := ctx.commandMatchers[cmdName]
	if cmdMatchers == nil {
		return parsedCommand{}, fmt.Errorf("no matching parser found for command '%s'", cmd)
	}

testCmdMatcher:
	for _, cmdMatcher := range cmdMatchers {
		argNames := []string{}
		argValues := []string{}
		cmdFragment := cmd[cmdMatchIdx[1]-1:]
		numOfArgMatchers := len(cmdMatcher.argMatchers)

		for i, argMatcher := range cmdMatcher.argMatchers {
			argValue, rest, wasParsed := parseArg(cmdFragment, argMatcher)
			if wasParsed {
				argNames = append(argNames, argMatcher.name)
				argValues = append(argValues, argValue)
				if strings.TrimSpace(rest) == "" && numOfArgMatchers-1 == i {
					break
				}
				if strings.TrimSpace(rest) != "" && numOfArgMatchers-1 == i {
					continue testCmdMatcher
				}
				cmdFragment = rest
			} else {
				continue testCmdMatcher
			}
		}
		return parsedCommand{command: cmd, commandName: cmdName, commandMatcher: cmdMatcher, argNames: argNames, argValues: argValues}, nil
	}
	return parsedCommand{}, fmt.Errorf("can't find command parser for '%s'", cmd)
}

func parseArg(cmd string, argM argMatcher) (string, string, bool) {
	argMatchIdx := argM.matcher.FindStringIndex(cmd)
	if argMatchIdx == nil {
		return "", "", false
	}
	argValue := strings.TrimSpace(cmd[argMatchIdx[0]:argMatchIdx[1]])
	return argValue, cmd[argMatchIdx[1]:], true
}

func (ctx *Context) AddCmdPattern(cmdPattern string) error {
	cmdMatchIdx := rCmd.FindStringIndex(cmdPattern)
	if cmdMatchIdx == nil {
		return fmt.Errorf("faild to parse command name from command pattern %s", cmdPattern)
	}
	cmdName := cmdPattern[cmdMatchIdx[0] : cmdMatchIdx[1]-1]
	cmd := cmdPattern[cmdMatchIdx[1]-1:]
	argMatchers := parseArgPattern(cmd, make([]argMatcher, 0))
	matcher := cmdMatcher{name: cmdName, argMatchers: argMatchers}

	if cmdMatchers, found := ctx.commandMatchers[cmdName]; found {
		cmdMatchers = append(cmdMatchers, matcher)
		ctx.commandMatchers[cmdName] = cmdMatchers
	} else {
		ctx.commandMatchers[cmdName] = []cmdMatcher{matcher}
	}
	return nil
}

func parseArgPattern(cmd string, argMatchers []argMatcher) []argMatcher {
	matchIdx := rArg.FindStringIndex(cmd)
	if matchIdx == nil {
		return argMatchers
	}
	argPattern := cmd[matchIdx[0]+2 : matchIdx[1]-1] // without '{' and '}'
	argSplit := strings.SplitN(argPattern, ":", 2)
	var matcher RExp = nil

	switch argSplit[1] {
	case "num":
		matcher = rNum
	default:
		panic("failed to parse matcher for " + argPattern)
	}

	argMatchers = append(argMatchers, argMatcher{name: argSplit[0], matcher: matcher})
	cmd = cmd[matchIdx[1]:]
	if len(cmd) > 0 {
		return parseArgPattern(cmd, argMatchers)
	} else {
		return argMatchers
	}
}

func (ctx *Context) PerformCmd(parsedCmd parsedCommand) error {
	funcName := strings.Title(parsedCmd.commandName)
	for _, argName := range parsedCmd.argNames {
		funcName = funcName + strings.Title(argName)
	}
	method := reflect.ValueOf(ctx).MethodByName(funcName)
	if method.IsValid() {
		numIn := method.Type().NumIn()
		in := make([]reflect.Value, numIn)
		argValues := parsedCmd.argValues
		if len(argValues) == numIn {
			fmt.Printf("method %T\n", method)
			for i := 0; i < numIn; i++ {
				t := method.Type().In(i)

				switch t.Kind() {
				case reflect.Float64:
					argVal, err := strconv.ParseFloat(argValues[i], 64)
					if err != nil {
						return fmt.Errorf("failed to parse as Float64 the argument value from '%s'", argValues[i])
					}
					in[i] = reflect.ValueOf(argVal)
				default:
					fmt.Printf("unhandled argument kind %s \n", t.Kind())
				}

				fmt.Printf("%T\n", t)
			}
			method.Call(in)
			return nil
		}
		return fmt.Errorf("num of arguments in parsed command %d doesn't match with the method num of input arguments '%d'", len(argValues), numIn)
	}
	return fmt.Errorf("can't find method by name '%s'", funcName)
}

func (ctx *Context) ExtrudeLength(length float64) error {
	fmt.Printf("-> ExtrudeLength %f\n", length)
	return nil
}

func (ctx *Context) InsetLength(length float64) error {
	fmt.Printf("-> InsetLength %f\n", length)
	return nil
}

func (ctx *Context) OutsetLength(length float64) error {
	fmt.Printf("-> OutsetLength %f\n", length)
	return nil
}

func (ctx *Context) RectLength(length float64) error {
	fmt.Printf("-> RectLength %f\n", length)
	return nil
}

func NotImplementedYet() error {
	return fmt.Errorf("not implemented yet")
}
