package csv

import (
	"testing"

	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"github.com/kylelemons/godebug/pretty"
)

type outer struct {
	Col1       string
	Col2       string
	S          []string      `csv:"slice"`
	M1         map[int]int   `csv:"map1"`
	M2         map[int]inner `csv:"map2"`
	XXX_ignore string
	Inner      inner
	InnerSlice []inner
}

type inner struct {
	Col3  string
	Col4  int
	Col5  float64
	Outer *outer
}

func TestMarshaler_Marshal(t *testing.T) {

	v := []outer{
		{
			Col1: "a",
			Col2: "b",
			S:    []string{"a", "b"},
			M1:   map[int]int{1: 2},
			M2: map[int]inner{
				1: {
					Col3: "x",
					Col4: 7,
					Col5: 7.22,
				},
			},
			Inner: inner{
				Col3: "c",
				Col4: 6,
				Col5: 6.22,
			},
			InnerSlice: []inner{
				{
					Col3: "u",
					Col4: 8,
					Col5: 8.22,
				},
				{
					Col3: "v",
					Col4: 9,
					Col5: 9.22,
				},
			},
		},
		{
			Col1: "e",
			Col2: "f",
			S:    []string{"a", "b", "c"},
			Inner: inner{
				Col3: "g",
				Col4: 6,
				Col5: 6.22,
			},
		},
	}
	v[1].Inner.Outer = &v[1] // introduce cycle

	tests := []struct {
		name    string
		v       interface{}
		want    string
		wantErr bool
	}{
		{
			name: "multiple empty slices",
			v: struct {
				a []string
				b []string
			}{
				a: nil,
				b: nil,
			},
			want: "",
		},
		{
			name: "multiple slices w/o struct",
			v: struct {
				a []string
				b []string
			}{
				a: []string{"a", "b"},
				b: []string{"c", "d"},
			},
			wantErr: true,
		},
		{
			name: "multiple slices w/ struct",
			v: struct {
				a []inner
				b []outer
			}{
				a: []inner{
					{
						Col3: "a",
						Col4: 1,
					},
					{
						Col3: "b",
						Col4: 2,
					},
				},
				b: []outer{
					{
						Col1: "c",
						Col2: "d",
					},
					{
						Col1: "e",
						Col2: "f",
					},
				},
			},
			want: "Col3;Col4;Col5\na;1;0\nb;2;0\n---\nCol1;Col2;slice;map1;map2;Col3;Col4;Col5;InnerSlice\nc;d;;;;;0;0;\ne;f;;;;;0;0;\n",
		},
		{
			name: "deep structure",
			v:    v,
			want: "Col1;Col2;slice;map1;map2;Col3;Col4;Col5;InnerSlice\na;b;a|b;1:2;1:x|7|7.22;c;6;6.22;u|8|8.22|v|9|9.22\ne;f;a|b|c;;;g;6;6.22;\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Marshaler{}
			g, err := m.Marshal(tt.v)
			if (err != nil) != tt.wantErr {
				t.Errorf("CSVMarshaler.Marshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			got := string(g)
			if diff := pretty.Compare(got, tt.want); diff != "" {
				t.Errorf("CSVMarshaler.Marshal() generate unexpected results:\n%s", diff)
			}
		})
	}
}

func TestMarshaler_ContentType(t *testing.T) {
	want := "text/csv"
	m := &Marshaler{}
	if got := m.ContentType(nil); got != want {
		t.Errorf("Marshaler.ContentType() = %v, want %v", got, want)
	}
}

func TestMarshaler_MarshalLanguage(t *testing.T) {

	v := []outer{
		{
			Col1: "a",
			Col2: "b",
			S:    []string{"a", "b"},
			M1:   map[int]int{1: 2},
			M2: map[int]inner{
				1: {
					Col3: "x",
					Col4: 7,
					Col5: 7.22,
				},
			},
			Inner: inner{
				Col3: "c",
				Col4: 6,
				Col5: 6.22,
			},
			InnerSlice: []inner{
				{
					Col3: "u",
					Col4: 8,
					Col5: 8.22,
				},
				{
					Col3: "v",
					Col4: 9,
					Col5: 9.22,
				},
			},
		},
		{
			Col1: "e",
			Col2: "f",
			S:    []string{"a", "b", "c"},
			Inner: inner{
				Col3: "g",
				Col4: 6,
				Col5: 6.22,
			},
		},
	}
	v[1].Inner.Outer = &v[1] // introduce cycle

	tests := []struct {
		name    string
		v       interface{}
		want    string
		wantErr bool
	}{
		{
			name: "multiple empty slices",
			v: struct {
				a []string
				b []string
			}{
				a: nil,
				b: nil,
			},
			want: "",
		},
		{
			name: "multiple slices w/o struct",
			v: struct {
				a []string
				b []string
			}{
				a: []string{"a", "b"},
				b: []string{"c", "d"},
			},
			wantErr: true,
		},
		{
			name: "multiple slices w/ struct",
			v: struct {
				a []inner
				b []outer
			}{
				a: []inner{
					{
						Col3: "a",
						Col4: 1,
					},
					{
						Col3: "b",
						Col4: 2,
					},
				},
				b: []outer{
					{
						Col1: "c",
						Col2: "d",
					},
					{
						Col1: "e",
						Col2: "f",
					},
				},
			},
			want: "Col3;Col4;Col5\na;1;0\nb;2;0\n---\nCol1;Col2;slice;map1;map2;Col3;Col4;Col5;InnerSlice\nc;d;;;;;0;0;\ne;f;;;;;0;0;\n",
		},
		{
			name: "deep structure",
			v:    v,
			want: "Col1;Col2;slice;map1;map2;Col3;Col4;Col5;InnerSlice\na;b;a|b;1:2;1:x|7|7,22;c;6;6,22;u|8|8,22|v|9|9,22\ne;f;a|b|c;;;g;6;6,22;\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := message.NewPrinter(language.German)
			m := &Marshaler{Printf: func(format string, a ...any) string { return p.Sprintf(format, a...) }}
			g, err := m.Marshal(tt.v)

			if (err != nil) != tt.wantErr {
				t.Errorf("CSVMarshaler.Marshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			got := string(g)
			if diff := pretty.Compare(got, tt.want); diff != "" {
				t.Errorf("CSVMarshaler.Marshal() generate unexpected results:\n%s", diff)
			}
		})
	}
}
