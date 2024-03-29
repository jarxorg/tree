package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	"github.com/jarxorg/io2"
	"github.com/jarxorg/tree"
)

func TestRun(t *testing.T) {
	stdinOrg := os.Stdin
	defer func() { os.Stdin = stdinOrg }()

	mustReadFileString := func(file string) string {
		bin, err := ioutil.ReadFile(file)
		if err != nil {
			t.Fatal(err)
		}
		return string(bin)
	}

	tests := []struct {
		stdin  string
		args   []string
		want   string
		errstr string
	}{
		{
			args: []string{},
			want: mustReadFileString("testdata/usage"),
		}, {
			args: []string{"-h"},
			want: mustReadFileString("testdata/usage"),
		}, {
			args: []string{"-v"},
			want: tree.VERSION + "\n",
		}, {
			stdin: "testdata/store.json",
			args:  []string{".store.book[0]"},
			want:  mustReadFileString("testdata/book-0.json"),
		}, {
			stdin: "testdata/store.yaml",
			args:  []string{".store.book[0]"},
			want:  mustReadFileString("testdata/book-0.yaml"),
		}, {
			stdin: "testdata/store.yaml",
			args:  []string{"-i", "yaml", ".store.book[0]"},
			want:  mustReadFileString("testdata/book-0.yaml"),
		}, {
			args: []string{".store.book[0]", "testdata/store.json"},
			want: mustReadFileString("testdata/book-0.json"),
		}, {
			args: []string{"-o", "yaml", ".store.book[0]", "testdata/store.json"},
			want: mustReadFileString("testdata/book-0.yaml"),
		}, {
			args: []string{".store.book[1:3]|", "testdata/store.json"},
			want: mustReadFileString("testdata/book-1-3.json"),
		}, {
			stdin: "testdata/store.json",
			args:  []string{"-x", ".store.book"},
			want:  mustReadFileString("testdata/book-x"),
		}, {
			stdin: "testdata/store.json",
			args:  []string{"-s", ".store.book[]"},
			want:  mustReadFileString("testdata/book-s"),
		}, {
			stdin: "testdata/book-x",
			args:  []string{"-s", "."},
			want:  mustReadFileString("testdata/book-s"),
		}, {
			stdin: "testdata/book-s",
			args:  []string{"-x", "."},
			want:  mustReadFileString("testdata/book-x"),
		}, {
			stdin: "testdata/store.json",
			args: []string{
				"-t", "{{.title}},{{.author}},{{.category}},{{.price}}",
				".store.book[]",
			},
			want: mustReadFileString("testdata/book.csv"),
		}, {
			stdin: "testdata/store.json",
			args:  []string{"-c", "."},
			want:  mustReadFileString("testdata/store-color.json"),
		}, {
			stdin: "testdata/store.yaml",
			args:  []string{"-c", "."},
			want:  mustReadFileString("testdata/store-color.yaml"),
		}, {
			stdin: "testdata/empty-object.json",
			args: []string{
				"-e", `.author = "Nigel Rees"`,
				"-e", `.category = "reference"`,
				"-e", `.price = 8.95`,
				"-e", `.title = "Sayings of the Century"`,
			},
			want: mustReadFileString("testdata/book-0.json"),
		}, {
			stdin: "testdata/null",
			args:  []string{"..walk"},
		}, {
			args:   []string{"-i", "json", ".", "testdata/invalid-json"},
			errstr: `failed to evaluate testdata/invalid-json: invalid character 'i' looking for beginning of value`,
		}, {
			stdin:  "testdata/invalid-json",
			args:   []string{"-i", "json", "."},
			errstr: `failed to evaluate STDIN: invalid character 'i' looking for beginning of value`,
		}, {
			args:   []string{"-i", "yaml", ".", "testdata/invalid-yaml"},
			errstr: `failed to evaluate testdata/invalid-yaml: yaml: found unexpected end of stream`,
		}, {
			args: []string{".", "testdata/book-0.yaml", "testdata/book-0.yaml"},
			want: mustReadFileString("testdata/book-0.yaml") + "---\n" + mustReadFileString("testdata/book-0.yaml"),
		},
	}
	fn := func(i int) {
		test := tests[i]
		if test.stdin != "" {
			in, err := os.Open(test.stdin)
			if err != nil {
				t.Fatalf("tests[%d] %v", i, err)
			}
			defer in.Close()
			os.Stdin = in
		}

		buf := new(bytes.Buffer)
		r := &runner{
			stderr: io2.NopWriteCloser(buf),
			out:    io2.NopWriteCloser(buf),
		}
		defer r.close()

		err := r.run(append([]string{"tq"}, test.args...))
		if test.errstr != "" {
			if err == nil {
				t.Fatalf("tests[%d] no error", i)
			}
			if err.Error() != test.errstr {
				t.Errorf(`tests[%d] error %s; want %s`, i, err.Error(), test.errstr)
			}
			return
		}
		if err != nil {
			t.Fatalf("tests[%d] %v", i, err)
		}
		if got := buf.String(); got != test.want {
			t.Errorf("tests[%d] got %s; want %s", i, got, test.want)
		}
	}
	for i := range tests {
		fn(i)
	}
}
