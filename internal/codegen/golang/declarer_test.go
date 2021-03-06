package golang

import (
	"github.com/jschaf/pggen/internal/casing"
	"github.com/jschaf/pggen/internal/codegen/golang/gotype"
	"github.com/jschaf/pggen/internal/pg"
	"github.com/jschaf/pggen/internal/texts"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFindDeclarer_Declare(t *testing.T) {
	caser := casing.NewCaser()
	caser.AddAcronym("ios", "IOS")
	emptyPkgPath := ""
	tests := []struct {
		name    string
		typ     gotype.Type
		pkgPath string
		want    []string
	}{
		{
			name: "enum - simple",
			typ: gotype.NewEnumType(
				emptyPkgPath,
				pg.EnumType{Name: "device_type", Labels: []string{"ios", "mobile"}},
				caser,
			),
			want: []string{
				texts.Dedent(`
				// DeviceType represents the Postgres enum "device_type".
				type DeviceType string
		
				const (
					DeviceTypeIOS    DeviceType = "ios"
					DeviceTypeMobile DeviceType = "mobile"
				)
		
				func (d DeviceType) String() string { return string(d) }
			`),
			},
		},
		{
			name: "enum - escaping",
			typ: gotype.NewEnumType(
				emptyPkgPath,
				pg.EnumType{Name: "quoting", Labels: []string{"\"\n\t", "`\"`"}},
				casing.NewCaser(),
			),
			want: []string{
				texts.Dedent(`
				// Quoting represents the Postgres enum "quoting".
				type Quoting string
		
				const (
					QuotingUnnamedLabel0 Quoting = "\"\n\t"
					QuotingUnnamedLabel1 Quoting = "` + "`" + `\"` + "`" + `"
				)
		
				func (q Quoting) String() string { return string(q) }
			`),
			},
		},
		{
			name: "composite",
			typ: gotype.CompositeType{
				PgComposite: pg.CompositeType{
					Name:        "some_table",
					ColumnNames: []string{"foo", "bar_baz"},
				},
				PkgPath:    "example.com/foo",
				Pkg:        "foo",
				Name:       "SomeTable",
				FieldNames: []string{"Foo", "BarBaz"},
				FieldTypes: []gotype.Type{gotype.Int16, gotype.PgText},
			},
			pkgPath: "example.com/foo",
			want: []string{
				texts.Dedent(`
				// SomeTable represents the Postgres composite type "some_table".
				type SomeTable struct {
					Foo    int16       ` + "`json:\"foo\"`" + `
					BarBaz pgtype.Text ` + "`json:\"bar_baz\"`" + `
				}
			`),
			},
		},
		{
			name: "nested composite",
			typ: gotype.CompositeType{
				PgComposite: pg.CompositeType{
					Name:        "some_table",
					ColumnNames: []string{"foo", "bar_baz"},
				},
				PkgPath:    "example.com/foo",
				Pkg:        "foo",
				Name:       "SomeTable",
				FieldNames: []string{"Foo", "BarBaz"},
				FieldTypes: []gotype.Type{
					gotype.CompositeType{
						PgComposite: pg.CompositeType{
							Name:        "foo_type",
							ColumnNames: []string{"alpha"},
						},
						PkgPath:    "example.com/foo",
						Pkg:        "foo",
						Name:       "FooType",
						FieldNames: []string{"Alpha"},
						FieldTypes: []gotype.Type{gotype.PgText},
					},
					gotype.PgText,
				},
			},
			pkgPath: "example.com/foo",
			want: []string{
				texts.Dedent(`
				// SomeTable represents the Postgres composite type "some_table".
				type SomeTable struct {
					Foo    FooType     ` + "`json:\"foo\"`" + `
					BarBaz pgtype.Text ` + "`json:\"bar_baz\"`" + `
				}
			`),
				texts.Dedent(`
				// FooType represents the Postgres composite type "foo_type".
				type FooType struct {
					Alpha pgtype.Text ` + "`json:\"alpha\"`" + `
				}
			`),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			decls := FindDeclarers(tt.typ)
			gotStrings := make([]string, len(decls))
			for i, decl := range decls {
				s, err := decl.Declare(tt.pkgPath)
				if err != nil {
					t.Fatal(err)
				}
				gotStrings[i] = s
			}
			assert.Equal(t, tt.want, gotStrings)
		})
	}
}
