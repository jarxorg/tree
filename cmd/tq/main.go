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

type stringList []string

func (l *stringList) String() string {
	return strings.Join(*l, ",")
}

func (l *stringList) Set(value string) error {
	*l = append(*l, value)
	return nil
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

type inputReader struct {
	io.ReadSeekCloser
}

func newInputReader(fargs []string) (*inputReader, error) {
	if len(fargs) == 0 {
		return newStdinReader()
	}

	var rs []io.ReadSeekCloser
	closeOnDefer := true
	defer func() {
		if closeOnDefer {
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
	closeOnDefer = false
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

type runner struct {
	flagSet      *flag.FlagSet
	isVersion    bool
	isHelp       bool
	isExpand     bool
	isSlurp      bool
	isRaw        bool
	outputFile   string
	tmplText     string
	tmpl         *template.Template
	inputFormat  string
	outputFormat string
	editExprs    stringList

	out              io.WriteCloser
	outputYAMLCalled int
}

func (r *runner) initFlagSet(args []string) error {
	s := flag.NewFlagSet(args[0], flag.ExitOnError)
	r.flagSet = s

	s.BoolVar(&r.isVersion, "v", false, "print version")
	s.BoolVar(&r.isHelp, "h", false, "help for "+cmd)
	s.BoolVar(&r.isExpand, "x", false, "expand results")
	s.BoolVar(&r.isSlurp, "s", false, "slurp all results into an array")
	s.BoolVar(&r.isRaw, "r", false, "output raw strings")
	s.StringVar(&r.outputFile, "O", "", "output file")
	s.StringVar(&r.tmplText, "t", "", "golang text/template string")
	s.StringVar(&r.inputFormat, "i", "", "input format (json or yaml)")
	s.StringVar(&r.outputFormat, "o", "json", "output format (json or yaml)")
	s.Var(&r.editExprs, "e", "edit expression")
	s.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s\n\nUsage:\n  %s\n\n", desc, usage)
		fmt.Fprintln(os.Stderr, "Flags:")
		s.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n%s", examplesText)
	}
	return s.Parse(args[1:])
}

func (r *runner) close() {
	if r.out != nil {
		r.out.Close()
		r.out = nil
	}
}

func (r *runner) run(args []string) error {
	defer r.close()

	if err := r.initFlagSet(args); err != nil {
		return err
	}
	if r.isVersion {
		fmt.Println(tree.VERSION)
		return nil
	}
	if r.isHelp || (r.flagSet.Arg(0) == "" && len(r.editExprs) == 0) {
		r.flagSet.Usage()
		return nil
	}
	if r.tmplText != "" {
		tmpl, err := template.New("").Parse(r.tmplText)
		if err != nil {
			return err
		}
		r.tmpl = tmpl
	}

	var fargs []string
	if args := r.flagSet.Args(); len(args) > 1 {
		fargs = args[1:]
	}
	if len(fargs) == 0 && term.IsTerminal(0) {
		r.flagSet.Usage()
		return nil
	}

	if r.outputFile != "" {
		out, err := os.OpenFile(r.outputFile, os.O_CREATE|os.O_WRONLY, os.ModePerm)
		if err != nil {
			return err
		}
		r.out = out
	} else {
		r.out = os.Stdout
	}

	in, err := newInputReader(fargs)
	if err != nil {
		return err
	}
	defer in.Close()

	return r.evaluate(in)
}

func (r *runner) evaluate(in io.ReadSeeker) error {
	switch r.inputFormat {
	case "yaml":
		return r.evaluateYAML(in)
	case "json":
		return r.evaluateJSON(in)
	}
	fns := []func(io.Reader) error{r.evaluateJSON, r.evaluateYAML}
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

func (r *runner) evaluateJSON(in io.Reader) error {
	dec := json.NewDecoder(in)
	for dec.More() {
		n, err := tree.DecodeJSON(dec)
		if err != nil {
			return &decodeError{err}
		}
		if err := r.evaluateNode(n); err != nil {
			return err
		}
	}
	return nil
}

func (r *runner) evaluateYAML(in io.Reader) error {
	dec := yaml.NewDecoder(in)
	for {
		n, err := tree.DecodeYAML(dec)
		if err != nil {
			if err == io.EOF {
				break
			}
			return &decodeError{err}
		}
		if err := r.evaluateNode(n); err != nil {
			return err
		}
	}
	return nil
}

func (r *runner) evaluateNode(node tree.Node) error {
	for _, expr := range r.editExprs {
		if err := tree.Edit(&node, expr); err != nil {
			return err
		}
	}
	expr := r.flagSet.Arg(0)
	if expr == "" {
		expr = "."
	}
	results, err := tree.Find(node, expr)
	if err != nil {
		return err
	}
	if len(results) == 0 {
		return nil
	}
	if r.isSlurp {
		results = []tree.Node{tree.Array(results)}
	}
	if r.isExpand {
		cb := func(_ interface{}, v tree.Node) error {
			return r.output(v)
		}
		for _, result := range results {
			if err := result.Each(cb); err != nil {
				return err
			}
		}
		return nil
	}
	for _, result := range results {
		if err := r.output(result); err != nil {
			return err
		}
	}
	return nil
}

func (r *runner) output(node tree.Node) error {
	if r.isRaw && node.Type().IsValue() {
		if _, err := fmt.Fprintln(r.out, node.Value().String()); err != nil {
			return err
		}
		return nil
	}
	if r.tmpl != nil {
		if err := r.tmpl.Execute(os.Stdout, node); err != nil {
			return err
		}
		if _, err := fmt.Fprintln(r.out); err != nil {
			return err
		}
		return nil
	}
	switch r.outputFormat {
	case "yaml":
		return r.outputYAML(node)
	}
	return r.outputJSON(node)
}

func (r *runner) outputYAML(node tree.Node) error {
	if r.outputYAMLCalled > 0 {
		if _, err := fmt.Fprintln(r.out, "---"); err != nil {
			return err
		}
	}
	r.outputYAMLCalled++

	b, err := tree.MarshalYAML(node)
	if err != nil {
		return err
	}
	if _, err := r.out.Write(b); err != nil {
		return err
	}
	return nil
}

func (r *runner) outputJSON(node tree.Node) error {
	b, err := json.MarshalIndent(node, "", "  ")
	if err != nil {
		return err
	}
	if _, err := r.out.Write(b); err != nil {
		return err
	}
	return nil
}

func main() {
	r := &runner{}
	if err := r.run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}
