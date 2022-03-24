package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"text/template"

	"github.com/jarxorg/tree"
	"golang.org/x/term"
	"gopkg.in/yaml.v2"
)

const (
	cmd          = "tq"
	desc         = cmd + " is a command-line JSON/YAML processor."
	usage        = cmd + " [flags] [query] ([file...])"
	examplesText = `Examples:
  % echo '{"colors": ["red", "green", "blue"]}' | tq '.colors[0]'
  "red"

  % echo '{"users":[{"id":1,"name":"one"},{"id":2,"name":"two"}]}' | tq -x -t '{{.id}}: {{.name}}' '.users'
  1: one
  2: two

  % echo '{}' | tq -e '.colors = ["red", "green"]' -e '.colors += "blue"' .
  {
    "colors": [
      "red",
      "green",
      "blue"
    ]
  }
`
)

type format string

func (f *format) String() string {
	return string(*f)
}

func (f *format) Set(value string) error {
	switch value {
	case "json":
		*f = "json"
		return nil
	case "yaml":
		*f = "yaml"
		return nil
	}
	return fmt.Errorf("unknown format")
}

type stringList []string

func (l *stringList) String() string {
	return strings.Join(*l, ",")
}

func (l *stringList) Set(value string) error {
	*l = append(*l, value)
	return nil
}

var (
	isHelp       bool
	isExpand     bool
	isSlurp      bool
	isRaw        bool
	tmplText     string
	tmpl         *template.Template
	inputFormat  = format("yaml")
	outputFormat = format("json")
	editExprs    stringList
)

func init() {
	flag.BoolVar(&isHelp, "h", false, "help for "+cmd)
	flag.BoolVar(&isExpand, "x", false, "expand results")
	flag.BoolVar(&isSlurp, "s", false, "slurp all results into an array")
	flag.BoolVar(&isRaw, "r", false, "output raw strings")
	flag.StringVar(&tmplText, "t", "", "golang text/template string (ignore -o flag)")
	flag.Var(&inputFormat, "i", "input format (json or yaml)")
	flag.Var(&outputFormat, "o", "output format (json or yaml)")
	flag.Var(&editExprs, "e", "edit expression")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s\n\nUsage:\n  %s\n\n", desc, usage)
		fmt.Fprintln(os.Stderr, "Flags:")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n%s", examplesText)
	}
}

func main() {
	flag.Parse()
	if isHelp || (flag.Arg(0) == "" && len(editExprs) == 0) {
		flag.Usage()
		return
	}
	handleError(run())
}

func handleError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}

func run() error {
	if tmplText != "" {
		var err error
		tmpl, err = template.New("").Parse(tmplText)
		if err != nil {
			return err
		}
	}

	var fargs []string
	if args := flag.Args(); len(args) > 1 {
		fargs = args[1:]
	}
	if len(fargs) == 0 && term.IsTerminal(0) {
		flag.Usage()
		return nil
	}

	in, err := newInputReader(fargs)
	if err != nil {
		return err
	}
	defer in.Close()

	switch inputFormat {
	case "yaml":
		return evaluateYAML(in)
	}
	return evaluateJSON(in)
}

type inputReader struct {
	io.Reader
	cs []io.Closer
}

func newInputReader(fargs []string) (*inputReader, error) {
	if len(fargs) == 0 {
		return &inputReader{Reader: os.Stdin}, nil
	}
	rs := make([]io.Reader, len(fargs))
	cs := make([]io.Closer, len(fargs))
	ok := false
	defer func() {
		if !ok {
			for _, c := range cs {
				if c != nil {
					c.Close()
				}
			}
		}
	}()
	for i, farg := range fargs {
		var err error
		f, err := os.Open(farg)
		if err != nil {
			return nil, err
		}
		rs[i] = f
		cs[i] = f
	}
	ok = true
	return &inputReader{Reader: io.MultiReader(rs...), cs: cs}, nil
}

func (r *inputReader) Close() error {
	var errs []error
	for _, c := range r.cs {
		if err := c.Close(); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return errs[0]
	}
	return nil
}

func evaluateJSON(in io.Reader) error {
	dec := json.NewDecoder(in)
	for dec.More() {
		n, err := tree.DecodeJSON(dec)
		if err != nil {
			return err
		}
		if err := evaluate(n); err != nil {
			return err
		}
	}
	return nil
}

func evaluateYAML(in io.Reader) error {
	dec := yaml.NewDecoder(in)
	for {
		n, err := tree.DecodeYAML(dec)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		if err := evaluate(n); err != nil {
			return err
		}
	}
	return nil
}

func evaluate(node tree.Node) error {
	for _, expr := range editExprs {
		if err := tree.Edit(&node, expr); err != nil {
			return err
		}
	}
	expr := flag.Arg(0)
	if expr == "" {
		expr = "."
	}
	rs, err := tree.Find(node, expr)
	if err != nil {
		return err
	}
	if len(rs) == 0 {
		return nil
	}
	if isSlurp {
		rs = []tree.Node{tree.Array(rs)}
	}
	if isExpand {
		cb := func(_ interface{}, v tree.Node) error {
			return output(v)
		}
		for _, r := range rs {
			if err := r.Each(cb); err != nil {
				return err
			}
		}
		return nil
	}
	for _, r := range rs {
		if err := output(r); err != nil {
			return err
		}
	}
	return nil
}

func output(node tree.Node) error {
	if isRaw && node.Type().IsValue() {
		fmt.Println(node.Value().String())
		return nil
	}
	if tmpl != nil {
		if err := tmpl.Execute(os.Stdout, node); err != nil {
			return err
		}
		fmt.Println()
		return nil
	}
	switch outputFormat {
	case "yaml":
		return outputYAML(node)
	}
	return outputJSON(node)
}

var outputYAMLCalled = 0

func outputYAML(node tree.Node) error {
	if outputYAMLCalled > 0 {
		fmt.Println("---")
	}
	out, err := tree.MarshalYAML(node)
	if err != nil {
		return err
	}
	fmt.Print(string(out))
	outputYAMLCalled++
	return nil
}

func outputJSON(node tree.Node) error {
	out, err := json.MarshalIndent(node, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(out))
	return nil
}
