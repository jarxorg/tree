package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"text/template"

	"github.com/jarxorg/io2"
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
	case "json", "yaml":
		*f = format(value)
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
	isVersion    bool
	isHelp       bool
	isExpand     bool
	isSlurp      bool
	isRaw        bool
	tmplText     string
	tmpl         *template.Template
	inputFormat  = format("")
	outputFormat = format("json")
	editExprs    stringList
)

func init() {
	flag.BoolVar(&isVersion, "v", false, "print version")
	flag.BoolVar(&isHelp, "h", false, "help for "+cmd)
	flag.BoolVar(&isExpand, "x", false, "expand results")
	flag.BoolVar(&isSlurp, "s", false, "slurp all results into an array")
	flag.BoolVar(&isRaw, "r", false, "output raw strings")
	flag.StringVar(&tmplText, "t", "", "golang text/template string (ignore -o flag)")
	flag.Var(&inputFormat, "i", "input format (json or yaml)")
	flag.Var(&outputFormat, "o", "output format (json or yaml, default json)")
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
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}

func run() error {
	if isVersion {
		fmt.Println(tree.VERSION)
		return nil
	}
	if isHelp || (flag.Arg(0) == "" && len(editExprs) == 0) {
		flag.Usage()
		return nil
	}
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

	return evaluate(in)
}

type inputReader struct {
	io.ReadSeekCloser
}

func newInputReader(fargs []string) (*inputReader, error) {
	if len(fargs) == 0 {
		return newStdinReader()
	}

	var rs []io.ReadSeekCloser
	ok := false
	defer func() {
		if !ok {
			for _, r := range rs {
				r.Close()
			}
		}
	}()
	isYaml := func(f string) bool {
		return strings.HasSuffix(f, ".yaml") || strings.HasSuffix(f, ".yml")
	}
	for _, farg := range fargs {
		var err error
		f, err := os.Open(farg)
		if err != nil {
			return nil, err
		}
		if len(rs) > 0 && isYaml(farg) {
			rs = append(rs, io2.NopReadSeekCloser(strings.NewReader("\n---\n")))
		}
		rs = append(rs, f)
	}
	mr, err := io2.MultiReadSeekCloser(rs...)
	if err != nil {
		return nil, err
	}
	ok = true
	return &inputReader{ReadSeekCloser: mr}, nil
}

func newStdinReader() (*inputReader, error) {
	tmp, err := os.CreateTemp("", "*.stdin")
	if err != nil {
		return nil, err
	}
	r := io2.DelegateReadSeekCloser(tmp)
	r.CloseFunc = func() error {
		_ = tmp.Close()
		return os.Remove(tmp.Name())
	}
	if _, err := io.Copy(tmp, os.Stdin); err != nil {
		r.Close()
		return nil, err
	}
	if _, err := r.Seek(0, io.SeekStart); err != nil {
		r.Close()
		return nil, err
	}
	return &inputReader{ReadSeekCloser: r}, nil
}

func evaluate(in io.ReadSeeker) error {
	switch inputFormat {
	case "json":
		return evaluateJSON(in)
	case "yaml":
		return evaluateYAML(in)
	}
	fns := []func(io.Reader) error{evaluateJSON, evaluateYAML}
	var errs []string
	for i, fn := range fns {
		if i > 0 {
			if _, err := in.Seek(0, io.SeekStart); err != nil {
				return err
			}
		}
		if err := fn(in); err != nil {
			errs = append(errs, err.Error())
			if !isDecodeError(err) {
				break
			}
			continue
		}
		return nil
	}
	return errors.New(strings.Join(errs, "; "))
}

type decodeError struct {
	err error
}

func (e *decodeError) Error() string {
	return e.err.Error()
}

func isDecodeError(err error) bool {
	_, ok := err.(*decodeError)
	return ok
}

func evaluateJSON(in io.Reader) error {
	dec := json.NewDecoder(in)
	for dec.More() {
		n, err := tree.DecodeJSON(dec)
		if err != nil {
			return &decodeError{err}
		}
		if err := evaluateNode(n); err != nil {
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
			return &decodeError{err}
		}
		if err := evaluateNode(n); err != nil {
			return err
		}
	}
	return nil
}

func evaluateNode(node tree.Node) error {
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
