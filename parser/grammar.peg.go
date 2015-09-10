package parser

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

type namespace struct {
	scope     string
	namespace string
}

type typeDef struct {
	name string
	typ  *Type
}

type exception *Struct

type union *Struct

type include string

func toIfaceSlice(v interface{}) []interface{} {
	if v == nil {
		return nil
	}
	return v.([]interface{})
}

func ifaceSliceToString(v interface{}) string {
	ifs := toIfaceSlice(v)
	b := make([]byte, len(ifs))
	for i, v := range ifs {
		b[i] = v.([]uint8)[0]
	}
	return string(b)
}

// toStruct converts a union to a struct with all fields optional.
func unionToStruct(u union) *Struct {
	st := (*Struct)(u)
	for _, f := range st.Fields {
		f.Optional = true
	}
	return st
}

var g = &grammar{
	rules: []*rule{
		{
			name: "Grammar",
			pos:  position{line: 52, col: 1, offset: 727},
			expr: &actionExpr{
				pos: position{line: 52, col: 11, offset: 739},
				run: (*parser).callonGrammar1,
				expr: &seqExpr{
					pos: position{line: 52, col: 11, offset: 739},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 52, col: 11, offset: 739},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 52, col: 14, offset: 742},
							label: "statements",
							expr: &zeroOrMoreExpr{
								pos: position{line: 52, col: 25, offset: 753},
								expr: &seqExpr{
									pos: position{line: 52, col: 27, offset: 755},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 52, col: 27, offset: 755},
											name: "Statement",
										},
										&ruleRefExpr{
											pos:  position{line: 52, col: 37, offset: 765},
											name: "__",
										},
									},
								},
							},
						},
						&choiceExpr{
							pos: position{line: 52, col: 44, offset: 772},
							alternatives: []interface{}{
								&ruleRefExpr{
									pos:  position{line: 52, col: 44, offset: 772},
									name: "EOF",
								},
								&ruleRefExpr{
									pos:  position{line: 52, col: 50, offset: 778},
									name: "SyntaxError",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "SyntaxError",
			pos:  position{line: 96, col: 1, offset: 1933},
			expr: &actionExpr{
				pos: position{line: 96, col: 15, offset: 1949},
				run: (*parser).callonSyntaxError1,
				expr: &anyMatcher{
					line: 96, col: 15, offset: 1949,
				},
			},
		},
		{
			name: "Include",
			pos:  position{line: 100, col: 1, offset: 2004},
			expr: &actionExpr{
				pos: position{line: 100, col: 11, offset: 2016},
				run: (*parser).callonInclude1,
				expr: &seqExpr{
					pos: position{line: 100, col: 11, offset: 2016},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 100, col: 11, offset: 2016},
							val:        "include",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 100, col: 21, offset: 2026},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 100, col: 23, offset: 2028},
							label: "file",
							expr: &ruleRefExpr{
								pos:  position{line: 100, col: 28, offset: 2033},
								name: "Literal",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 100, col: 36, offset: 2041},
							name: "EOS",
						},
					},
				},
			},
		},
		{
			name: "Statement",
			pos:  position{line: 104, col: 1, offset: 2086},
			expr: &choiceExpr{
				pos: position{line: 104, col: 13, offset: 2100},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 104, col: 13, offset: 2100},
						name: "Include",
					},
					&ruleRefExpr{
						pos:  position{line: 104, col: 23, offset: 2110},
						name: "Namespace",
					},
					&ruleRefExpr{
						pos:  position{line: 104, col: 35, offset: 2122},
						name: "Const",
					},
					&ruleRefExpr{
						pos:  position{line: 104, col: 43, offset: 2130},
						name: "Enum",
					},
					&ruleRefExpr{
						pos:  position{line: 104, col: 50, offset: 2137},
						name: "TypeDef",
					},
					&ruleRefExpr{
						pos:  position{line: 104, col: 60, offset: 2147},
						name: "Struct",
					},
					&ruleRefExpr{
						pos:  position{line: 104, col: 69, offset: 2156},
						name: "Exception",
					},
					&ruleRefExpr{
						pos:  position{line: 104, col: 81, offset: 2168},
						name: "Union",
					},
					&ruleRefExpr{
						pos:  position{line: 104, col: 89, offset: 2176},
						name: "Service",
					},
				},
			},
		},
		{
			name: "Namespace",
			pos:  position{line: 106, col: 1, offset: 2185},
			expr: &actionExpr{
				pos: position{line: 106, col: 13, offset: 2199},
				run: (*parser).callonNamespace1,
				expr: &seqExpr{
					pos: position{line: 106, col: 13, offset: 2199},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 106, col: 13, offset: 2199},
							val:        "namespace",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 106, col: 25, offset: 2211},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 106, col: 27, offset: 2213},
							label: "scope",
							expr: &oneOrMoreExpr{
								pos: position{line: 106, col: 33, offset: 2219},
								expr: &charClassMatcher{
									pos:        position{line: 106, col: 33, offset: 2219},
									val:        "[a-z.-]",
									chars:      []rune{'.', '-'},
									ranges:     []rune{'a', 'z'},
									ignoreCase: false,
									inverted:   false,
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 106, col: 42, offset: 2228},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 106, col: 44, offset: 2230},
							label: "ns",
							expr: &ruleRefExpr{
								pos:  position{line: 106, col: 47, offset: 2233},
								name: "Identifier",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 106, col: 58, offset: 2244},
							name: "EOS",
						},
					},
				},
			},
		},
		{
			name: "Const",
			pos:  position{line: 113, col: 1, offset: 2355},
			expr: &actionExpr{
				pos: position{line: 113, col: 9, offset: 2365},
				run: (*parser).callonConst1,
				expr: &seqExpr{
					pos: position{line: 113, col: 9, offset: 2365},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 113, col: 9, offset: 2365},
							val:        "const",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 113, col: 17, offset: 2373},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 113, col: 19, offset: 2375},
							label: "typ",
							expr: &ruleRefExpr{
								pos:  position{line: 113, col: 23, offset: 2379},
								name: "FieldType",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 113, col: 33, offset: 2389},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 113, col: 35, offset: 2391},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 113, col: 40, offset: 2396},
								name: "Identifier",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 113, col: 51, offset: 2407},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 113, col: 53, offset: 2409},
							val:        "=",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 113, col: 57, offset: 2413},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 113, col: 59, offset: 2415},
							label: "value",
							expr: &ruleRefExpr{
								pos:  position{line: 113, col: 65, offset: 2421},
								name: "ConstValue",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 113, col: 76, offset: 2432},
							name: "EOS",
						},
					},
				},
			},
		},
		{
			name: "Enum",
			pos:  position{line: 121, col: 1, offset: 2540},
			expr: &actionExpr{
				pos: position{line: 121, col: 8, offset: 2549},
				run: (*parser).callonEnum1,
				expr: &seqExpr{
					pos: position{line: 121, col: 8, offset: 2549},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 121, col: 8, offset: 2549},
							val:        "enum",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 121, col: 15, offset: 2556},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 121, col: 17, offset: 2558},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 121, col: 22, offset: 2563},
								name: "Identifier",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 121, col: 33, offset: 2574},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 121, col: 35, offset: 2576},
							val:        "{",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 121, col: 39, offset: 2580},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 121, col: 42, offset: 2583},
							label: "values",
							expr: &zeroOrMoreExpr{
								pos: position{line: 121, col: 49, offset: 2590},
								expr: &seqExpr{
									pos: position{line: 121, col: 50, offset: 2591},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 121, col: 50, offset: 2591},
											name: "EnumValue",
										},
										&ruleRefExpr{
											pos:  position{line: 121, col: 60, offset: 2601},
											name: "__",
										},
									},
								},
							},
						},
						&litMatcher{
							pos:        position{line: 121, col: 65, offset: 2606},
							val:        "}",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 121, col: 69, offset: 2610},
							name: "EOS",
						},
					},
				},
			},
		},
		{
			name: "EnumValue",
			pos:  position{line: 144, col: 1, offset: 3126},
			expr: &actionExpr{
				pos: position{line: 144, col: 13, offset: 3140},
				run: (*parser).callonEnumValue1,
				expr: &seqExpr{
					pos: position{line: 144, col: 13, offset: 3140},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 144, col: 13, offset: 3140},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 144, col: 18, offset: 3145},
								name: "Identifier",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 144, col: 29, offset: 3156},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 144, col: 31, offset: 3158},
							label: "value",
							expr: &zeroOrOneExpr{
								pos: position{line: 144, col: 37, offset: 3164},
								expr: &seqExpr{
									pos: position{line: 144, col: 38, offset: 3165},
									exprs: []interface{}{
										&litMatcher{
											pos:        position{line: 144, col: 38, offset: 3165},
											val:        "=",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 144, col: 42, offset: 3169},
											name: "_",
										},
										&ruleRefExpr{
											pos:  position{line: 144, col: 44, offset: 3171},
											name: "IntConstant",
										},
									},
								},
							},
						},
						&zeroOrOneExpr{
							pos: position{line: 144, col: 58, offset: 3185},
							expr: &ruleRefExpr{
								pos:  position{line: 144, col: 58, offset: 3185},
								name: "ListSeparator",
							},
						},
					},
				},
			},
		},
		{
			name: "TypeDef",
			pos:  position{line: 155, col: 1, offset: 3364},
			expr: &actionExpr{
				pos: position{line: 155, col: 11, offset: 3376},
				run: (*parser).callonTypeDef1,
				expr: &seqExpr{
					pos: position{line: 155, col: 11, offset: 3376},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 155, col: 11, offset: 3376},
							val:        "typedef",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 155, col: 21, offset: 3386},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 155, col: 23, offset: 3388},
							label: "typ",
							expr: &ruleRefExpr{
								pos:  position{line: 155, col: 27, offset: 3392},
								name: "FieldType",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 155, col: 37, offset: 3402},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 155, col: 39, offset: 3404},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 155, col: 44, offset: 3409},
								name: "Identifier",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 155, col: 55, offset: 3420},
							name: "EOS",
						},
					},
				},
			},
		},
		{
			name: "Struct",
			pos:  position{line: 162, col: 1, offset: 3510},
			expr: &actionExpr{
				pos: position{line: 162, col: 10, offset: 3521},
				run: (*parser).callonStruct1,
				expr: &seqExpr{
					pos: position{line: 162, col: 10, offset: 3521},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 162, col: 10, offset: 3521},
							val:        "struct",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 162, col: 19, offset: 3530},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 162, col: 21, offset: 3532},
							label: "st",
							expr: &ruleRefExpr{
								pos:  position{line: 162, col: 24, offset: 3535},
								name: "StructLike",
							},
						},
					},
				},
			},
		},
		{
			name: "Exception",
			pos:  position{line: 163, col: 1, offset: 3575},
			expr: &actionExpr{
				pos: position{line: 163, col: 13, offset: 3589},
				run: (*parser).callonException1,
				expr: &seqExpr{
					pos: position{line: 163, col: 13, offset: 3589},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 163, col: 13, offset: 3589},
							val:        "exception",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 163, col: 25, offset: 3601},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 163, col: 27, offset: 3603},
							label: "st",
							expr: &ruleRefExpr{
								pos:  position{line: 163, col: 30, offset: 3606},
								name: "StructLike",
							},
						},
					},
				},
			},
		},
		{
			name: "Union",
			pos:  position{line: 164, col: 1, offset: 3657},
			expr: &actionExpr{
				pos: position{line: 164, col: 9, offset: 3667},
				run: (*parser).callonUnion1,
				expr: &seqExpr{
					pos: position{line: 164, col: 9, offset: 3667},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 164, col: 9, offset: 3667},
							val:        "union",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 164, col: 17, offset: 3675},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 164, col: 19, offset: 3677},
							label: "st",
							expr: &ruleRefExpr{
								pos:  position{line: 164, col: 22, offset: 3680},
								name: "StructLike",
							},
						},
					},
				},
			},
		},
		{
			name: "StructLike",
			pos:  position{line: 165, col: 1, offset: 3727},
			expr: &actionExpr{
				pos: position{line: 165, col: 14, offset: 3742},
				run: (*parser).callonStructLike1,
				expr: &seqExpr{
					pos: position{line: 165, col: 14, offset: 3742},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 165, col: 14, offset: 3742},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 165, col: 19, offset: 3747},
								name: "Identifier",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 165, col: 30, offset: 3758},
							name: "__",
						},
						&litMatcher{
							pos:        position{line: 165, col: 33, offset: 3761},
							val:        "{",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 165, col: 37, offset: 3765},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 165, col: 40, offset: 3768},
							label: "fields",
							expr: &ruleRefExpr{
								pos:  position{line: 165, col: 47, offset: 3775},
								name: "FieldList",
							},
						},
						&litMatcher{
							pos:        position{line: 165, col: 57, offset: 3785},
							val:        "}",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 165, col: 61, offset: 3789},
							name: "EOS",
						},
					},
				},
			},
		},
		{
			name: "FieldList",
			pos:  position{line: 175, col: 1, offset: 3923},
			expr: &actionExpr{
				pos: position{line: 175, col: 13, offset: 3937},
				run: (*parser).callonFieldList1,
				expr: &labeledExpr{
					pos:   position{line: 175, col: 13, offset: 3937},
					label: "fields",
					expr: &zeroOrMoreExpr{
						pos: position{line: 175, col: 20, offset: 3944},
						expr: &seqExpr{
							pos: position{line: 175, col: 21, offset: 3945},
							exprs: []interface{}{
								&ruleRefExpr{
									pos:  position{line: 175, col: 21, offset: 3945},
									name: "Field",
								},
								&ruleRefExpr{
									pos:  position{line: 175, col: 27, offset: 3951},
									name: "__",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Field",
			pos:  position{line: 184, col: 1, offset: 4111},
			expr: &actionExpr{
				pos: position{line: 184, col: 9, offset: 4121},
				run: (*parser).callonField1,
				expr: &seqExpr{
					pos: position{line: 184, col: 9, offset: 4121},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 184, col: 9, offset: 4121},
							label: "id",
							expr: &ruleRefExpr{
								pos:  position{line: 184, col: 12, offset: 4124},
								name: "IntConstant",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 184, col: 24, offset: 4136},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 184, col: 26, offset: 4138},
							val:        ":",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 184, col: 30, offset: 4142},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 184, col: 32, offset: 4144},
							label: "req",
							expr: &zeroOrOneExpr{
								pos: position{line: 184, col: 36, offset: 4148},
								expr: &ruleRefExpr{
									pos:  position{line: 184, col: 36, offset: 4148},
									name: "FieldReq",
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 184, col: 46, offset: 4158},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 184, col: 48, offset: 4160},
							label: "typ",
							expr: &ruleRefExpr{
								pos:  position{line: 184, col: 52, offset: 4164},
								name: "FieldType",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 184, col: 62, offset: 4174},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 184, col: 64, offset: 4176},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 184, col: 69, offset: 4181},
								name: "Identifier",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 184, col: 80, offset: 4192},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 184, col: 82, offset: 4194},
							label: "def",
							expr: &zeroOrOneExpr{
								pos: position{line: 184, col: 86, offset: 4198},
								expr: &seqExpr{
									pos: position{line: 184, col: 87, offset: 4199},
									exprs: []interface{}{
										&litMatcher{
											pos:        position{line: 184, col: 87, offset: 4199},
											val:        "=",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 184, col: 91, offset: 4203},
											name: "_",
										},
										&ruleRefExpr{
											pos:  position{line: 184, col: 93, offset: 4205},
											name: "ConstValue",
										},
									},
								},
							},
						},
						&zeroOrOneExpr{
							pos: position{line: 184, col: 106, offset: 4218},
							expr: &ruleRefExpr{
								pos:  position{line: 184, col: 106, offset: 4218},
								name: "ListSeparator",
							},
						},
					},
				},
			},
		},
		{
			name: "FieldReq",
			pos:  position{line: 199, col: 1, offset: 4478},
			expr: &actionExpr{
				pos: position{line: 199, col: 12, offset: 4491},
				run: (*parser).callonFieldReq1,
				expr: &choiceExpr{
					pos: position{line: 199, col: 13, offset: 4492},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 199, col: 13, offset: 4492},
							val:        "required",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 199, col: 26, offset: 4505},
							val:        "optional",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "Service",
			pos:  position{line: 203, col: 1, offset: 4576},
			expr: &actionExpr{
				pos: position{line: 203, col: 11, offset: 4588},
				run: (*parser).callonService1,
				expr: &seqExpr{
					pos: position{line: 203, col: 11, offset: 4588},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 203, col: 11, offset: 4588},
							val:        "service",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 203, col: 21, offset: 4598},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 203, col: 23, offset: 4600},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 203, col: 28, offset: 4605},
								name: "Identifier",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 203, col: 39, offset: 4616},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 203, col: 41, offset: 4618},
							label: "extends",
							expr: &zeroOrOneExpr{
								pos: position{line: 203, col: 49, offset: 4626},
								expr: &seqExpr{
									pos: position{line: 203, col: 50, offset: 4627},
									exprs: []interface{}{
										&litMatcher{
											pos:        position{line: 203, col: 50, offset: 4627},
											val:        "extends",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 203, col: 60, offset: 4637},
											name: "__",
										},
										&ruleRefExpr{
											pos:  position{line: 203, col: 63, offset: 4640},
											name: "Identifier",
										},
										&ruleRefExpr{
											pos:  position{line: 203, col: 74, offset: 4651},
											name: "__",
										},
									},
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 203, col: 79, offset: 4656},
							name: "__",
						},
						&litMatcher{
							pos:        position{line: 203, col: 82, offset: 4659},
							val:        "{",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 203, col: 86, offset: 4663},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 203, col: 89, offset: 4666},
							label: "methods",
							expr: &zeroOrMoreExpr{
								pos: position{line: 203, col: 97, offset: 4674},
								expr: &seqExpr{
									pos: position{line: 203, col: 98, offset: 4675},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 203, col: 98, offset: 4675},
											name: "Function",
										},
										&ruleRefExpr{
											pos:  position{line: 203, col: 107, offset: 4684},
											name: "__",
										},
									},
								},
							},
						},
						&choiceExpr{
							pos: position{line: 203, col: 113, offset: 4690},
							alternatives: []interface{}{
								&litMatcher{
									pos:        position{line: 203, col: 113, offset: 4690},
									val:        "}",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 203, col: 119, offset: 4696},
									name: "EndOfServiceError",
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 203, col: 138, offset: 4715},
							name: "EOS",
						},
					},
				},
			},
		},
		{
			name: "EndOfServiceError",
			pos:  position{line: 218, col: 1, offset: 5056},
			expr: &actionExpr{
				pos: position{line: 218, col: 21, offset: 5078},
				run: (*parser).callonEndOfServiceError1,
				expr: &anyMatcher{
					line: 218, col: 21, offset: 5078,
				},
			},
		},
		{
			name: "Function",
			pos:  position{line: 222, col: 1, offset: 5144},
			expr: &actionExpr{
				pos: position{line: 222, col: 12, offset: 5157},
				run: (*parser).callonFunction1,
				expr: &seqExpr{
					pos: position{line: 222, col: 12, offset: 5157},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 222, col: 12, offset: 5157},
							label: "oneway",
							expr: &zeroOrOneExpr{
								pos: position{line: 222, col: 19, offset: 5164},
								expr: &seqExpr{
									pos: position{line: 222, col: 20, offset: 5165},
									exprs: []interface{}{
										&litMatcher{
											pos:        position{line: 222, col: 20, offset: 5165},
											val:        "oneway",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 222, col: 29, offset: 5174},
											name: "__",
										},
									},
								},
							},
						},
						&labeledExpr{
							pos:   position{line: 222, col: 34, offset: 5179},
							label: "typ",
							expr: &ruleRefExpr{
								pos:  position{line: 222, col: 38, offset: 5183},
								name: "FunctionType",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 222, col: 51, offset: 5196},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 222, col: 54, offset: 5199},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 222, col: 59, offset: 5204},
								name: "Identifier",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 222, col: 70, offset: 5215},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 222, col: 72, offset: 5217},
							val:        "(",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 222, col: 76, offset: 5221},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 222, col: 79, offset: 5224},
							label: "arguments",
							expr: &ruleRefExpr{
								pos:  position{line: 222, col: 89, offset: 5234},
								name: "FieldList",
							},
						},
						&litMatcher{
							pos:        position{line: 222, col: 99, offset: 5244},
							val:        ")",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 222, col: 103, offset: 5248},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 222, col: 106, offset: 5251},
							label: "exceptions",
							expr: &zeroOrOneExpr{
								pos: position{line: 222, col: 117, offset: 5262},
								expr: &ruleRefExpr{
									pos:  position{line: 222, col: 117, offset: 5262},
									name: "Throws",
								},
							},
						},
						&zeroOrOneExpr{
							pos: position{line: 222, col: 125, offset: 5270},
							expr: &ruleRefExpr{
								pos:  position{line: 222, col: 125, offset: 5270},
								name: "ListSeparator",
							},
						},
					},
				},
			},
		},
		{
			name: "FunctionType",
			pos:  position{line: 245, col: 1, offset: 5651},
			expr: &actionExpr{
				pos: position{line: 245, col: 16, offset: 5668},
				run: (*parser).callonFunctionType1,
				expr: &labeledExpr{
					pos:   position{line: 245, col: 16, offset: 5668},
					label: "typ",
					expr: &choiceExpr{
						pos: position{line: 245, col: 21, offset: 5673},
						alternatives: []interface{}{
							&litMatcher{
								pos:        position{line: 245, col: 21, offset: 5673},
								val:        "void",
								ignoreCase: false,
							},
							&ruleRefExpr{
								pos:  position{line: 245, col: 30, offset: 5682},
								name: "FieldType",
							},
						},
					},
				},
			},
		},
		{
			name: "Throws",
			pos:  position{line: 252, col: 1, offset: 5789},
			expr: &actionExpr{
				pos: position{line: 252, col: 10, offset: 5800},
				run: (*parser).callonThrows1,
				expr: &seqExpr{
					pos: position{line: 252, col: 10, offset: 5800},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 252, col: 10, offset: 5800},
							val:        "throws",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 252, col: 19, offset: 5809},
							name: "__",
						},
						&litMatcher{
							pos:        position{line: 252, col: 22, offset: 5812},
							val:        "(",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 252, col: 26, offset: 5816},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 252, col: 29, offset: 5819},
							label: "exceptions",
							expr: &ruleRefExpr{
								pos:  position{line: 252, col: 40, offset: 5830},
								name: "FieldList",
							},
						},
						&litMatcher{
							pos:        position{line: 252, col: 50, offset: 5840},
							val:        ")",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "FieldType",
			pos:  position{line: 256, col: 1, offset: 5873},
			expr: &actionExpr{
				pos: position{line: 256, col: 13, offset: 5887},
				run: (*parser).callonFieldType1,
				expr: &labeledExpr{
					pos:   position{line: 256, col: 13, offset: 5887},
					label: "typ",
					expr: &choiceExpr{
						pos: position{line: 256, col: 18, offset: 5892},
						alternatives: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 256, col: 18, offset: 5892},
								name: "BaseType",
							},
							&ruleRefExpr{
								pos:  position{line: 256, col: 29, offset: 5903},
								name: "ContainerType",
							},
							&ruleRefExpr{
								pos:  position{line: 256, col: 45, offset: 5919},
								name: "Identifier",
							},
						},
					},
				},
			},
		},
		{
			name: "DefinitionType",
			pos:  position{line: 263, col: 1, offset: 6029},
			expr: &actionExpr{
				pos: position{line: 263, col: 18, offset: 6048},
				run: (*parser).callonDefinitionType1,
				expr: &labeledExpr{
					pos:   position{line: 263, col: 18, offset: 6048},
					label: "typ",
					expr: &choiceExpr{
						pos: position{line: 263, col: 23, offset: 6053},
						alternatives: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 263, col: 23, offset: 6053},
								name: "BaseType",
							},
							&ruleRefExpr{
								pos:  position{line: 263, col: 34, offset: 6064},
								name: "ContainerType",
							},
						},
					},
				},
			},
		},
		{
			name: "BaseType",
			pos:  position{line: 267, col: 1, offset: 6101},
			expr: &actionExpr{
				pos: position{line: 267, col: 12, offset: 6114},
				run: (*parser).callonBaseType1,
				expr: &choiceExpr{
					pos: position{line: 267, col: 13, offset: 6115},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 267, col: 13, offset: 6115},
							val:        "bool",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 267, col: 22, offset: 6124},
							val:        "byte",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 267, col: 31, offset: 6133},
							val:        "i16",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 267, col: 39, offset: 6141},
							val:        "i32",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 267, col: 47, offset: 6149},
							val:        "i64",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 267, col: 55, offset: 6157},
							val:        "double",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 267, col: 66, offset: 6168},
							val:        "string",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 267, col: 77, offset: 6179},
							val:        "binary",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "ContainerType",
			pos:  position{line: 271, col: 1, offset: 6236},
			expr: &actionExpr{
				pos: position{line: 271, col: 17, offset: 6254},
				run: (*parser).callonContainerType1,
				expr: &labeledExpr{
					pos:   position{line: 271, col: 17, offset: 6254},
					label: "typ",
					expr: &choiceExpr{
						pos: position{line: 271, col: 22, offset: 6259},
						alternatives: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 271, col: 22, offset: 6259},
								name: "MapType",
							},
							&ruleRefExpr{
								pos:  position{line: 271, col: 32, offset: 6269},
								name: "SetType",
							},
							&ruleRefExpr{
								pos:  position{line: 271, col: 42, offset: 6279},
								name: "ListType",
							},
						},
					},
				},
			},
		},
		{
			name: "MapType",
			pos:  position{line: 275, col: 1, offset: 6311},
			expr: &actionExpr{
				pos: position{line: 275, col: 11, offset: 6323},
				run: (*parser).callonMapType1,
				expr: &seqExpr{
					pos: position{line: 275, col: 11, offset: 6323},
					exprs: []interface{}{
						&zeroOrOneExpr{
							pos: position{line: 275, col: 11, offset: 6323},
							expr: &ruleRefExpr{
								pos:  position{line: 275, col: 11, offset: 6323},
								name: "CppType",
							},
						},
						&litMatcher{
							pos:        position{line: 275, col: 20, offset: 6332},
							val:        "map<",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 275, col: 27, offset: 6339},
							name: "WS",
						},
						&labeledExpr{
							pos:   position{line: 275, col: 30, offset: 6342},
							label: "key",
							expr: &ruleRefExpr{
								pos:  position{line: 275, col: 34, offset: 6346},
								name: "FieldType",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 275, col: 44, offset: 6356},
							name: "WS",
						},
						&litMatcher{
							pos:        position{line: 275, col: 47, offset: 6359},
							val:        ",",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 275, col: 51, offset: 6363},
							name: "WS",
						},
						&labeledExpr{
							pos:   position{line: 275, col: 54, offset: 6366},
							label: "value",
							expr: &ruleRefExpr{
								pos:  position{line: 275, col: 60, offset: 6372},
								name: "FieldType",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 275, col: 70, offset: 6382},
							name: "WS",
						},
						&litMatcher{
							pos:        position{line: 275, col: 73, offset: 6385},
							val:        ">",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "SetType",
			pos:  position{line: 283, col: 1, offset: 6484},
			expr: &actionExpr{
				pos: position{line: 283, col: 11, offset: 6496},
				run: (*parser).callonSetType1,
				expr: &seqExpr{
					pos: position{line: 283, col: 11, offset: 6496},
					exprs: []interface{}{
						&zeroOrOneExpr{
							pos: position{line: 283, col: 11, offset: 6496},
							expr: &ruleRefExpr{
								pos:  position{line: 283, col: 11, offset: 6496},
								name: "CppType",
							},
						},
						&litMatcher{
							pos:        position{line: 283, col: 20, offset: 6505},
							val:        "set<",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 283, col: 27, offset: 6512},
							name: "WS",
						},
						&labeledExpr{
							pos:   position{line: 283, col: 30, offset: 6515},
							label: "typ",
							expr: &ruleRefExpr{
								pos:  position{line: 283, col: 34, offset: 6519},
								name: "FieldType",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 283, col: 44, offset: 6529},
							name: "WS",
						},
						&litMatcher{
							pos:        position{line: 283, col: 47, offset: 6532},
							val:        ">",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "ListType",
			pos:  position{line: 290, col: 1, offset: 6605},
			expr: &actionExpr{
				pos: position{line: 290, col: 12, offset: 6618},
				run: (*parser).callonListType1,
				expr: &seqExpr{
					pos: position{line: 290, col: 12, offset: 6618},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 290, col: 12, offset: 6618},
							val:        "list<",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 290, col: 20, offset: 6626},
							name: "WS",
						},
						&labeledExpr{
							pos:   position{line: 290, col: 23, offset: 6629},
							label: "typ",
							expr: &ruleRefExpr{
								pos:  position{line: 290, col: 27, offset: 6633},
								name: "FieldType",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 290, col: 37, offset: 6643},
							name: "WS",
						},
						&litMatcher{
							pos:        position{line: 290, col: 40, offset: 6646},
							val:        ">",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "CppType",
			pos:  position{line: 297, col: 1, offset: 6720},
			expr: &actionExpr{
				pos: position{line: 297, col: 11, offset: 6732},
				run: (*parser).callonCppType1,
				expr: &seqExpr{
					pos: position{line: 297, col: 11, offset: 6732},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 297, col: 11, offset: 6732},
							val:        "cpp_type",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 297, col: 22, offset: 6743},
							label: "cppType",
							expr: &ruleRefExpr{
								pos:  position{line: 297, col: 30, offset: 6751},
								name: "Literal",
							},
						},
					},
				},
			},
		},
		{
			name: "ConstValue",
			pos:  position{line: 301, col: 1, offset: 6785},
			expr: &choiceExpr{
				pos: position{line: 301, col: 14, offset: 6800},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 301, col: 14, offset: 6800},
						name: "Literal",
					},
					&ruleRefExpr{
						pos:  position{line: 301, col: 24, offset: 6810},
						name: "DoubleConstant",
					},
					&ruleRefExpr{
						pos:  position{line: 301, col: 41, offset: 6827},
						name: "IntConstant",
					},
					&ruleRefExpr{
						pos:  position{line: 301, col: 55, offset: 6841},
						name: "ConstMap",
					},
					&ruleRefExpr{
						pos:  position{line: 301, col: 66, offset: 6852},
						name: "ConstList",
					},
					&ruleRefExpr{
						pos:  position{line: 301, col: 78, offset: 6864},
						name: "Identifier",
					},
				},
			},
		},
		{
			name: "IntConstant",
			pos:  position{line: 303, col: 1, offset: 6876},
			expr: &actionExpr{
				pos: position{line: 303, col: 15, offset: 6892},
				run: (*parser).callonIntConstant1,
				expr: &seqExpr{
					pos: position{line: 303, col: 15, offset: 6892},
					exprs: []interface{}{
						&zeroOrOneExpr{
							pos: position{line: 303, col: 15, offset: 6892},
							expr: &charClassMatcher{
								pos:        position{line: 303, col: 15, offset: 6892},
								val:        "[-+]",
								chars:      []rune{'-', '+'},
								ignoreCase: false,
								inverted:   false,
							},
						},
						&oneOrMoreExpr{
							pos: position{line: 303, col: 21, offset: 6898},
							expr: &ruleRefExpr{
								pos:  position{line: 303, col: 21, offset: 6898},
								name: "Digit",
							},
						},
					},
				},
			},
		},
		{
			name: "DoubleConstant",
			pos:  position{line: 307, col: 1, offset: 6959},
			expr: &actionExpr{
				pos: position{line: 307, col: 18, offset: 6978},
				run: (*parser).callonDoubleConstant1,
				expr: &seqExpr{
					pos: position{line: 307, col: 18, offset: 6978},
					exprs: []interface{}{
						&zeroOrOneExpr{
							pos: position{line: 307, col: 18, offset: 6978},
							expr: &charClassMatcher{
								pos:        position{line: 307, col: 18, offset: 6978},
								val:        "[+-]",
								chars:      []rune{'+', '-'},
								ignoreCase: false,
								inverted:   false,
							},
						},
						&zeroOrMoreExpr{
							pos: position{line: 307, col: 24, offset: 6984},
							expr: &ruleRefExpr{
								pos:  position{line: 307, col: 24, offset: 6984},
								name: "Digit",
							},
						},
						&litMatcher{
							pos:        position{line: 307, col: 31, offset: 6991},
							val:        ".",
							ignoreCase: false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 307, col: 35, offset: 6995},
							expr: &ruleRefExpr{
								pos:  position{line: 307, col: 35, offset: 6995},
								name: "Digit",
							},
						},
						&zeroOrOneExpr{
							pos: position{line: 307, col: 42, offset: 7002},
							expr: &seqExpr{
								pos: position{line: 307, col: 44, offset: 7004},
								exprs: []interface{}{
									&charClassMatcher{
										pos:        position{line: 307, col: 44, offset: 7004},
										val:        "['Ee']",
										chars:      []rune{'\'', 'E', 'e', '\''},
										ignoreCase: false,
										inverted:   false,
									},
									&ruleRefExpr{
										pos:  position{line: 307, col: 51, offset: 7011},
										name: "IntConstant",
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "ConstList",
			pos:  position{line: 311, col: 1, offset: 7078},
			expr: &actionExpr{
				pos: position{line: 311, col: 13, offset: 7092},
				run: (*parser).callonConstList1,
				expr: &seqExpr{
					pos: position{line: 311, col: 13, offset: 7092},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 311, col: 13, offset: 7092},
							val:        "[",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 311, col: 17, offset: 7096},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 311, col: 20, offset: 7099},
							label: "values",
							expr: &zeroOrMoreExpr{
								pos: position{line: 311, col: 27, offset: 7106},
								expr: &seqExpr{
									pos: position{line: 311, col: 28, offset: 7107},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 311, col: 28, offset: 7107},
											name: "ConstValue",
										},
										&ruleRefExpr{
											pos:  position{line: 311, col: 39, offset: 7118},
											name: "__",
										},
										&zeroOrOneExpr{
											pos: position{line: 311, col: 42, offset: 7121},
											expr: &ruleRefExpr{
												pos:  position{line: 311, col: 42, offset: 7121},
												name: "ListSeparator",
											},
										},
										&ruleRefExpr{
											pos:  position{line: 311, col: 57, offset: 7136},
											name: "__",
										},
									},
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 311, col: 62, offset: 7141},
							name: "__",
						},
						&litMatcher{
							pos:        position{line: 311, col: 65, offset: 7144},
							val:        "]",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "ConstMap",
			pos:  position{line: 320, col: 1, offset: 7317},
			expr: &actionExpr{
				pos: position{line: 320, col: 12, offset: 7330},
				run: (*parser).callonConstMap1,
				expr: &seqExpr{
					pos: position{line: 320, col: 12, offset: 7330},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 320, col: 12, offset: 7330},
							val:        "{",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 320, col: 16, offset: 7334},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 320, col: 19, offset: 7337},
							label: "values",
							expr: &zeroOrMoreExpr{
								pos: position{line: 320, col: 26, offset: 7344},
								expr: &seqExpr{
									pos: position{line: 320, col: 27, offset: 7345},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 320, col: 27, offset: 7345},
											name: "ConstValue",
										},
										&ruleRefExpr{
											pos:  position{line: 320, col: 38, offset: 7356},
											name: "__",
										},
										&litMatcher{
											pos:        position{line: 320, col: 41, offset: 7359},
											val:        ":",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 320, col: 45, offset: 7363},
											name: "__",
										},
										&ruleRefExpr{
											pos:  position{line: 320, col: 48, offset: 7366},
											name: "ConstValue",
										},
										&ruleRefExpr{
											pos:  position{line: 320, col: 59, offset: 7377},
											name: "__",
										},
										&choiceExpr{
											pos: position{line: 320, col: 63, offset: 7381},
											alternatives: []interface{}{
												&litMatcher{
													pos:        position{line: 320, col: 63, offset: 7381},
													val:        ",",
													ignoreCase: false,
												},
												&andExpr{
													pos: position{line: 320, col: 69, offset: 7387},
													expr: &litMatcher{
														pos:        position{line: 320, col: 70, offset: 7388},
														val:        "}",
														ignoreCase: false,
													},
												},
											},
										},
										&ruleRefExpr{
											pos:  position{line: 320, col: 75, offset: 7393},
											name: "__",
										},
									},
								},
							},
						},
						&litMatcher{
							pos:        position{line: 320, col: 80, offset: 7398},
							val:        "}",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "Literal",
			pos:  position{line: 336, col: 1, offset: 7644},
			expr: &actionExpr{
				pos: position{line: 336, col: 11, offset: 7656},
				run: (*parser).callonLiteral1,
				expr: &choiceExpr{
					pos: position{line: 336, col: 12, offset: 7657},
					alternatives: []interface{}{
						&seqExpr{
							pos: position{line: 336, col: 13, offset: 7658},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 336, col: 13, offset: 7658},
									val:        "\"",
									ignoreCase: false,
								},
								&zeroOrMoreExpr{
									pos: position{line: 336, col: 17, offset: 7662},
									expr: &choiceExpr{
										pos: position{line: 336, col: 18, offset: 7663},
										alternatives: []interface{}{
											&litMatcher{
												pos:        position{line: 336, col: 18, offset: 7663},
												val:        "\\\"",
												ignoreCase: false,
											},
											&charClassMatcher{
												pos:        position{line: 336, col: 25, offset: 7670},
												val:        "[^\"]",
												chars:      []rune{'"'},
												ignoreCase: false,
												inverted:   true,
											},
										},
									},
								},
								&litMatcher{
									pos:        position{line: 336, col: 32, offset: 7677},
									val:        "\"",
									ignoreCase: false,
								},
							},
						},
						&seqExpr{
							pos: position{line: 336, col: 40, offset: 7685},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 336, col: 40, offset: 7685},
									val:        "'",
									ignoreCase: false,
								},
								&zeroOrMoreExpr{
									pos: position{line: 336, col: 45, offset: 7690},
									expr: &choiceExpr{
										pos: position{line: 336, col: 46, offset: 7691},
										alternatives: []interface{}{
											&litMatcher{
												pos:        position{line: 336, col: 46, offset: 7691},
												val:        "\\'",
												ignoreCase: false,
											},
											&charClassMatcher{
												pos:        position{line: 336, col: 53, offset: 7698},
												val:        "[^']",
												chars:      []rune{'\''},
												ignoreCase: false,
												inverted:   true,
											},
										},
									},
								},
								&litMatcher{
									pos:        position{line: 336, col: 60, offset: 7705},
									val:        "'",
									ignoreCase: false,
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Identifier",
			pos:  position{line: 343, col: 1, offset: 7906},
			expr: &actionExpr{
				pos: position{line: 343, col: 14, offset: 7921},
				run: (*parser).callonIdentifier1,
				expr: &seqExpr{
					pos: position{line: 343, col: 14, offset: 7921},
					exprs: []interface{}{
						&oneOrMoreExpr{
							pos: position{line: 343, col: 14, offset: 7921},
							expr: &choiceExpr{
								pos: position{line: 343, col: 15, offset: 7922},
								alternatives: []interface{}{
									&ruleRefExpr{
										pos:  position{line: 343, col: 15, offset: 7922},
										name: "Letter",
									},
									&litMatcher{
										pos:        position{line: 343, col: 24, offset: 7931},
										val:        "_",
										ignoreCase: false,
									},
								},
							},
						},
						&zeroOrMoreExpr{
							pos: position{line: 343, col: 30, offset: 7937},
							expr: &choiceExpr{
								pos: position{line: 343, col: 31, offset: 7938},
								alternatives: []interface{}{
									&ruleRefExpr{
										pos:  position{line: 343, col: 31, offset: 7938},
										name: "Letter",
									},
									&ruleRefExpr{
										pos:  position{line: 343, col: 40, offset: 7947},
										name: "Digit",
									},
									&charClassMatcher{
										pos:        position{line: 343, col: 48, offset: 7955},
										val:        "[._]",
										chars:      []rune{'.', '_'},
										ignoreCase: false,
										inverted:   false,
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "ListSeparator",
			pos:  position{line: 347, col: 1, offset: 8007},
			expr: &charClassMatcher{
				pos:        position{line: 347, col: 17, offset: 8025},
				val:        "[,;]",
				chars:      []rune{',', ';'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "Letter",
			pos:  position{line: 348, col: 1, offset: 8030},
			expr: &charClassMatcher{
				pos:        position{line: 348, col: 10, offset: 8041},
				val:        "[A-Za-z]",
				ranges:     []rune{'A', 'Z', 'a', 'z'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "Digit",
			pos:  position{line: 349, col: 1, offset: 8050},
			expr: &charClassMatcher{
				pos:        position{line: 349, col: 9, offset: 8060},
				val:        "[0-9]",
				ranges:     []rune{'0', '9'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "SourceChar",
			pos:  position{line: 353, col: 1, offset: 8071},
			expr: &anyMatcher{
				line: 353, col: 14, offset: 8086,
			},
		},
		{
			name: "Comment",
			pos:  position{line: 354, col: 1, offset: 8088},
			expr: &choiceExpr{
				pos: position{line: 354, col: 11, offset: 8100},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 354, col: 11, offset: 8100},
						name: "MultiLineComment",
					},
					&ruleRefExpr{
						pos:  position{line: 354, col: 30, offset: 8119},
						name: "SingleLineComment",
					},
				},
			},
		},
		{
			name: "MultiLineComment",
			pos:  position{line: 355, col: 1, offset: 8137},
			expr: &seqExpr{
				pos: position{line: 355, col: 20, offset: 8158},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 355, col: 20, offset: 8158},
						val:        "/*",
						ignoreCase: false,
					},
					&zeroOrMoreExpr{
						pos: position{line: 355, col: 25, offset: 8163},
						expr: &seqExpr{
							pos: position{line: 355, col: 27, offset: 8165},
							exprs: []interface{}{
								&notExpr{
									pos: position{line: 355, col: 27, offset: 8165},
									expr: &litMatcher{
										pos:        position{line: 355, col: 28, offset: 8166},
										val:        "*/",
										ignoreCase: false,
									},
								},
								&ruleRefExpr{
									pos:  position{line: 355, col: 33, offset: 8171},
									name: "SourceChar",
								},
							},
						},
					},
					&litMatcher{
						pos:        position{line: 355, col: 47, offset: 8185},
						val:        "*/",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "MultiLineCommentNoLineTerminator",
			pos:  position{line: 356, col: 1, offset: 8190},
			expr: &seqExpr{
				pos: position{line: 356, col: 36, offset: 8227},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 356, col: 36, offset: 8227},
						val:        "/*",
						ignoreCase: false,
					},
					&zeroOrMoreExpr{
						pos: position{line: 356, col: 41, offset: 8232},
						expr: &seqExpr{
							pos: position{line: 356, col: 43, offset: 8234},
							exprs: []interface{}{
								&notExpr{
									pos: position{line: 356, col: 43, offset: 8234},
									expr: &choiceExpr{
										pos: position{line: 356, col: 46, offset: 8237},
										alternatives: []interface{}{
											&litMatcher{
												pos:        position{line: 356, col: 46, offset: 8237},
												val:        "*/",
												ignoreCase: false,
											},
											&ruleRefExpr{
												pos:  position{line: 356, col: 53, offset: 8244},
												name: "EOL",
											},
										},
									},
								},
								&ruleRefExpr{
									pos:  position{line: 356, col: 59, offset: 8250},
									name: "SourceChar",
								},
							},
						},
					},
					&litMatcher{
						pos:        position{line: 356, col: 73, offset: 8264},
						val:        "*/",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "SingleLineComment",
			pos:  position{line: 357, col: 1, offset: 8269},
			expr: &choiceExpr{
				pos: position{line: 357, col: 21, offset: 8291},
				alternatives: []interface{}{
					&seqExpr{
						pos: position{line: 357, col: 22, offset: 8292},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 357, col: 22, offset: 8292},
								val:        "//",
								ignoreCase: false,
							},
							&zeroOrMoreExpr{
								pos: position{line: 357, col: 27, offset: 8297},
								expr: &seqExpr{
									pos: position{line: 357, col: 29, offset: 8299},
									exprs: []interface{}{
										&notExpr{
											pos: position{line: 357, col: 29, offset: 8299},
											expr: &ruleRefExpr{
												pos:  position{line: 357, col: 30, offset: 8300},
												name: "EOL",
											},
										},
										&ruleRefExpr{
											pos:  position{line: 357, col: 34, offset: 8304},
											name: "SourceChar",
										},
									},
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 357, col: 52, offset: 8322},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 357, col: 52, offset: 8322},
								val:        "#",
								ignoreCase: false,
							},
							&zeroOrMoreExpr{
								pos: position{line: 357, col: 56, offset: 8326},
								expr: &seqExpr{
									pos: position{line: 357, col: 58, offset: 8328},
									exprs: []interface{}{
										&notExpr{
											pos: position{line: 357, col: 58, offset: 8328},
											expr: &ruleRefExpr{
												pos:  position{line: 357, col: 59, offset: 8329},
												name: "EOL",
											},
										},
										&ruleRefExpr{
											pos:  position{line: 357, col: 63, offset: 8333},
											name: "SourceChar",
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "__",
			pos:  position{line: 359, col: 1, offset: 8349},
			expr: &zeroOrMoreExpr{
				pos: position{line: 359, col: 6, offset: 8356},
				expr: &choiceExpr{
					pos: position{line: 359, col: 8, offset: 8358},
					alternatives: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 359, col: 8, offset: 8358},
							name: "Whitespace",
						},
						&ruleRefExpr{
							pos:  position{line: 359, col: 21, offset: 8371},
							name: "EOL",
						},
						&ruleRefExpr{
							pos:  position{line: 359, col: 27, offset: 8377},
							name: "Comment",
						},
					},
				},
			},
		},
		{
			name: "_",
			pos:  position{line: 360, col: 1, offset: 8388},
			expr: &zeroOrMoreExpr{
				pos: position{line: 360, col: 5, offset: 8394},
				expr: &choiceExpr{
					pos: position{line: 360, col: 7, offset: 8396},
					alternatives: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 360, col: 7, offset: 8396},
							name: "Whitespace",
						},
						&ruleRefExpr{
							pos:  position{line: 360, col: 20, offset: 8409},
							name: "MultiLineCommentNoLineTerminator",
						},
					},
				},
			},
		},
		{
			name: "WS",
			pos:  position{line: 361, col: 1, offset: 8445},
			expr: &zeroOrMoreExpr{
				pos: position{line: 361, col: 6, offset: 8452},
				expr: &ruleRefExpr{
					pos:  position{line: 361, col: 6, offset: 8452},
					name: "Whitespace",
				},
			},
		},
		{
			name: "Whitespace",
			pos:  position{line: 363, col: 1, offset: 8465},
			expr: &charClassMatcher{
				pos:        position{line: 363, col: 14, offset: 8480},
				val:        "[ \\t\\r]",
				chars:      []rune{' ', '\t', '\r'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "EOL",
			pos:  position{line: 364, col: 1, offset: 8488},
			expr: &litMatcher{
				pos:        position{line: 364, col: 7, offset: 8496},
				val:        "\n",
				ignoreCase: false,
			},
		},
		{
			name: "EOS",
			pos:  position{line: 365, col: 1, offset: 8501},
			expr: &choiceExpr{
				pos: position{line: 365, col: 7, offset: 8509},
				alternatives: []interface{}{
					&seqExpr{
						pos: position{line: 365, col: 7, offset: 8509},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 365, col: 7, offset: 8509},
								name: "__",
							},
							&litMatcher{
								pos:        position{line: 365, col: 10, offset: 8512},
								val:        ";",
								ignoreCase: false,
							},
						},
					},
					&seqExpr{
						pos: position{line: 365, col: 16, offset: 8518},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 365, col: 16, offset: 8518},
								name: "_",
							},
							&zeroOrOneExpr{
								pos: position{line: 365, col: 18, offset: 8520},
								expr: &ruleRefExpr{
									pos:  position{line: 365, col: 18, offset: 8520},
									name: "SingleLineComment",
								},
							},
							&ruleRefExpr{
								pos:  position{line: 365, col: 37, offset: 8539},
								name: "EOL",
							},
						},
					},
					&seqExpr{
						pos: position{line: 365, col: 43, offset: 8545},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 365, col: 43, offset: 8545},
								name: "__",
							},
							&ruleRefExpr{
								pos:  position{line: 365, col: 46, offset: 8548},
								name: "EOF",
							},
						},
					},
				},
			},
		},
		{
			name: "EOF",
			pos:  position{line: 367, col: 1, offset: 8553},
			expr: &notExpr{
				pos: position{line: 367, col: 7, offset: 8561},
				expr: &anyMatcher{
					line: 367, col: 8, offset: 8562,
				},
			},
		},
	},
}

func (c *current) onGrammar1(statements interface{}) (interface{}, error) {
	thrift := &Thrift{
		Includes:   make(map[string]string),
		Namespaces: make(map[string]string),
		Typedefs:   make(map[string]*Type),
		Constants:  make(map[string]*Constant),
		Enums:      make(map[string]*Enum),
		Structs:    make(map[string]*Struct),
		Exceptions: make(map[string]*Struct),
		Unions:     make(map[string]*Struct),
		Services:   make(map[string]*Service),
	}
	stmts := toIfaceSlice(statements)
	for _, st := range stmts {
		switch v := st.([]interface{})[0].(type) {
		case *namespace:
			thrift.Namespaces[v.scope] = v.namespace
		case *Constant:
			thrift.Constants[v.Name] = v
		case *Enum:
			thrift.Enums[v.Name] = v
		case *typeDef:
			thrift.Typedefs[v.name] = v.typ
		case *Struct:
			thrift.Structs[v.Name] = v
		case exception:
			thrift.Exceptions[v.Name] = (*Struct)(v)
		case union:
			thrift.Unions[v.Name] = unionToStruct(v)
		case *Service:
			thrift.Services[v.Name] = v
		case include:
			name := string(v)
			if ix := strings.LastIndex(name, "."); ix > 0 {
				name = name[:ix]
			}
			thrift.Includes[name] = string(v)
		default:
			return nil, fmt.Errorf("parser: unknown value %#v", v)
		}
	}
	return thrift, nil
}

func (p *parser) callonGrammar1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onGrammar1(stack["statements"])
}

func (c *current) onSyntaxError1() (interface{}, error) {
	return nil, errors.New("parser: syntax error")
}

func (p *parser) callonSyntaxError1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onSyntaxError1()
}

func (c *current) onInclude1(file interface{}) (interface{}, error) {
	return include(file.(string)), nil
}

func (p *parser) callonInclude1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onInclude1(stack["file"])
}

func (c *current) onNamespace1(scope, ns interface{}) (interface{}, error) {
	return &namespace{
		scope:     ifaceSliceToString(scope),
		namespace: string(ns.(Identifier)),
	}, nil
}

func (p *parser) callonNamespace1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onNamespace1(stack["scope"], stack["ns"])
}

func (c *current) onConst1(typ, name, value interface{}) (interface{}, error) {
	return &Constant{
		Name:  string(name.(Identifier)),
		Type:  typ.(*Type),
		Value: value,
	}, nil
}

func (p *parser) callonConst1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onConst1(stack["typ"], stack["name"], stack["value"])
}

func (c *current) onEnum1(name, values interface{}) (interface{}, error) {
	vs := toIfaceSlice(values)
	en := &Enum{
		Name:   string(name.(Identifier)),
		Values: make(map[string]*EnumValue, len(vs)),
	}
	// Assigns numbers in order. This will behave badly if some values are
	// defined and other are not, but I think that's ok since that's a silly
	// thing to do.
	next := 0
	for _, v := range vs {
		ev := v.([]interface{})[0].(*EnumValue)
		if ev.Value < 0 {
			ev.Value = next
		}
		if ev.Value >= next {
			next = ev.Value + 1
		}
		en.Values[ev.Name] = ev
	}
	return en, nil
}

func (p *parser) callonEnum1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onEnum1(stack["name"], stack["values"])
}

func (c *current) onEnumValue1(name, value interface{}) (interface{}, error) {
	ev := &EnumValue{
		Name:  string(name.(Identifier)),
		Value: -1,
	}
	if value != nil {
		ev.Value = int(value.([]interface{})[2].(int64))
	}
	return ev, nil
}

func (p *parser) callonEnumValue1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onEnumValue1(stack["name"], stack["value"])
}

func (c *current) onTypeDef1(typ, name interface{}) (interface{}, error) {
	return &typeDef{
		name: string(name.(Identifier)),
		typ:  typ.(*Type),
	}, nil
}

func (p *parser) callonTypeDef1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onTypeDef1(stack["typ"], stack["name"])
}

func (c *current) onStruct1(st interface{}) (interface{}, error) {
	return st.(*Struct), nil
}

func (p *parser) callonStruct1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onStruct1(stack["st"])
}

func (c *current) onException1(st interface{}) (interface{}, error) {
	return exception(st.(*Struct)), nil
}

func (p *parser) callonException1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onException1(stack["st"])
}

func (c *current) onUnion1(st interface{}) (interface{}, error) {
	return union(st.(*Struct)), nil
}

func (p *parser) callonUnion1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onUnion1(stack["st"])
}

func (c *current) onStructLike1(name, fields interface{}) (interface{}, error) {
	st := &Struct{
		Name: string(name.(Identifier)),
	}
	if fields != nil {
		st.Fields = fields.([]*Field)
	}
	return st, nil
}

func (p *parser) callonStructLike1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onStructLike1(stack["name"], stack["fields"])
}

func (c *current) onFieldList1(fields interface{}) (interface{}, error) {
	fs := fields.([]interface{})
	flds := make([]*Field, len(fs))
	for i, f := range fs {
		flds[i] = f.([]interface{})[0].(*Field)
	}
	return flds, nil
}

func (p *parser) callonFieldList1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onFieldList1(stack["fields"])
}

func (c *current) onField1(id, req, typ, name, def interface{}) (interface{}, error) {
	f := &Field{
		ID:   int(id.(int64)),
		Name: string(name.(Identifier)),
		Type: typ.(*Type),
	}
	if req != nil && !req.(bool) {
		f.Optional = true
	}
	if def != nil {
		f.Default = def.([]interface{})[2]
	}
	return f, nil
}

func (p *parser) callonField1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onField1(stack["id"], stack["req"], stack["typ"], stack["name"], stack["def"])
}

func (c *current) onFieldReq1() (interface{}, error) {
	return !bytes.Equal(c.text, []byte("optional")), nil
}

func (p *parser) callonFieldReq1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onFieldReq1()
}

func (c *current) onService1(name, extends, methods interface{}) (interface{}, error) {
	ms := methods.([]interface{})
	svc := &Service{
		Name:    string(name.(Identifier)),
		Methods: make(map[string]*Method, len(ms)),
	}
	if extends != nil {
		svc.Extends = string(extends.([]interface{})[2].(Identifier))
	}
	for _, m := range ms {
		mt := m.([]interface{})[0].(*Method)
		svc.Methods[mt.Name] = mt
	}
	return svc, nil
}

func (p *parser) callonService1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onService1(stack["name"], stack["extends"], stack["methods"])
}

func (c *current) onEndOfServiceError1() (interface{}, error) {
	return nil, errors.New("parser: expected end of service")
}

func (p *parser) callonEndOfServiceError1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onEndOfServiceError1()
}

func (c *current) onFunction1(oneway, typ, name, arguments, exceptions interface{}) (interface{}, error) {
	m := &Method{
		Name: string(name.(Identifier)),
	}
	t := typ.(*Type)
	if t.Name != "void" {
		m.ReturnType = t
	}
	if oneway != nil {
		m.Oneway = true
	}
	if arguments != nil {
		m.Arguments = arguments.([]*Field)
	}
	if exceptions != nil {
		m.Exceptions = exceptions.([]*Field)
		for _, e := range m.Exceptions {
			e.Optional = true
		}
	}
	return m, nil
}

func (p *parser) callonFunction1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onFunction1(stack["oneway"], stack["typ"], stack["name"], stack["arguments"], stack["exceptions"])
}

func (c *current) onFunctionType1(typ interface{}) (interface{}, error) {
	if t, ok := typ.(*Type); ok {
		return t, nil
	}
	return &Type{Name: string(c.text)}, nil
}

func (p *parser) callonFunctionType1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onFunctionType1(stack["typ"])
}

func (c *current) onThrows1(exceptions interface{}) (interface{}, error) {
	return exceptions, nil
}

func (p *parser) callonThrows1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onThrows1(stack["exceptions"])
}

func (c *current) onFieldType1(typ interface{}) (interface{}, error) {
	if t, ok := typ.(Identifier); ok {
		return &Type{Name: string(t)}, nil
	}
	return typ, nil
}

func (p *parser) callonFieldType1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onFieldType1(stack["typ"])
}

func (c *current) onDefinitionType1(typ interface{}) (interface{}, error) {
	return typ, nil
}

func (p *parser) callonDefinitionType1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onDefinitionType1(stack["typ"])
}

func (c *current) onBaseType1() (interface{}, error) {
	return &Type{Name: string(c.text)}, nil
}

func (p *parser) callonBaseType1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onBaseType1()
}

func (c *current) onContainerType1(typ interface{}) (interface{}, error) {
	return typ, nil
}

func (p *parser) callonContainerType1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onContainerType1(stack["typ"])
}

func (c *current) onMapType1(key, value interface{}) (interface{}, error) {
	return &Type{
		Name:      "map",
		KeyType:   key.(*Type),
		ValueType: value.(*Type),
	}, nil
}

func (p *parser) callonMapType1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onMapType1(stack["key"], stack["value"])
}

func (c *current) onSetType1(typ interface{}) (interface{}, error) {
	return &Type{
		Name:      "set",
		ValueType: typ.(*Type),
	}, nil
}

func (p *parser) callonSetType1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onSetType1(stack["typ"])
}

func (c *current) onListType1(typ interface{}) (interface{}, error) {
	return &Type{
		Name:      "list",
		ValueType: typ.(*Type),
	}, nil
}

func (p *parser) callonListType1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onListType1(stack["typ"])
}

func (c *current) onCppType1(cppType interface{}) (interface{}, error) {
	return cppType, nil
}

func (p *parser) callonCppType1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onCppType1(stack["cppType"])
}

func (c *current) onIntConstant1() (interface{}, error) {
	return strconv.ParseInt(string(c.text), 10, 64)
}

func (p *parser) callonIntConstant1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onIntConstant1()
}

func (c *current) onDoubleConstant1() (interface{}, error) {
	return strconv.ParseFloat(string(c.text), 64)
}

func (p *parser) callonDoubleConstant1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onDoubleConstant1()
}

func (c *current) onConstList1(values interface{}) (interface{}, error) {
	valueSlice := values.([]interface{})
	vs := make([]interface{}, len(valueSlice))
	for i, v := range valueSlice {
		vs[i] = v.([]interface{})[0]
	}
	return vs, nil
}

func (p *parser) callonConstList1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onConstList1(stack["values"])
}

func (c *current) onConstMap1(values interface{}) (interface{}, error) {
	if values == nil {
		return nil, nil
	}
	vals := values.([]interface{})
	kvs := make([]KeyValue, len(vals))
	for i, kv := range vals {
		v := kv.([]interface{})
		kvs[i] = KeyValue{
			Key:   v[0],
			Value: v[4],
		}
	}
	return kvs, nil
}

func (p *parser) callonConstMap1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onConstMap1(stack["values"])
}

func (c *current) onLiteral1() (interface{}, error) {
	if len(c.text) != 0 && c.text[0] == '\'' {
		return strconv.Unquote(`"` + strings.Replace(string(c.text[1:len(c.text)-1]), `\'`, `'`, -1) + `"`)
	}
	return strconv.Unquote(string(c.text))
}

func (p *parser) callonLiteral1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onLiteral1()
}

func (c *current) onIdentifier1() (interface{}, error) {
	return Identifier(string(c.text)), nil
}

func (p *parser) callonIdentifier1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onIdentifier1()
}

var (
	// errNoRule is returned when the grammar to parse has no rule.
	errNoRule = errors.New("grammar has no rule")

	// errInvalidEncoding is returned when the source is not properly
	// utf8-encoded.
	errInvalidEncoding = errors.New("invalid encoding")

	// errNoMatch is returned if no match could be found.
	errNoMatch = errors.New("no match found")
)

// Option is a function that can set an option on the parser. It returns
// the previous setting as an Option.
type Option func(*parser) Option

// Debug creates an Option to set the debug flag to b. When set to true,
// debugging information is printed to stdout while parsing.
//
// The default is false.
func Debug(b bool) Option {
	return func(p *parser) Option {
		old := p.debug
		p.debug = b
		return Debug(old)
	}
}

// Memoize creates an Option to set the memoize flag to b. When set to true,
// the parser will cache all results so each expression is evaluated only
// once. This guarantees linear parsing time even for pathological cases,
// at the expense of more memory and slower times for typical cases.
//
// The default is false.
func Memoize(b bool) Option {
	return func(p *parser) Option {
		old := p.memoize
		p.memoize = b
		return Memoize(old)
	}
}

// Recover creates an Option to set the recover flag to b. When set to
// true, this causes the parser to recover from panics and convert it
// to an error. Setting it to false can be useful while debugging to
// access the full stack trace.
//
// The default is true.
func Recover(b bool) Option {
	return func(p *parser) Option {
		old := p.recover
		p.recover = b
		return Recover(old)
	}
}

// ParseFile parses the file identified by filename.
func ParseFile(filename string, opts ...Option) (interface{}, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ParseReader(filename, f, opts...)
}

// ParseReader parses the data from r using filename as information in the
// error messages.
func ParseReader(filename string, r io.Reader, opts ...Option) (interface{}, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return Parse(filename, b, opts...)
}

// Parse parses the data from b using filename as information in the
// error messages.
func Parse(filename string, b []byte, opts ...Option) (interface{}, error) {
	return newParser(filename, b, opts...).parse(g)
}

// position records a position in the text.
type position struct {
	line, col, offset int
}

func (p position) String() string {
	return fmt.Sprintf("%d:%d [%d]", p.line, p.col, p.offset)
}

// savepoint stores all state required to go back to this point in the
// parser.
type savepoint struct {
	position
	rn rune
	w  int
}

type current struct {
	pos  position // start position of the match
	text []byte   // raw text of the match
}

// the AST types...

type grammar struct {
	pos   position
	rules []*rule
}

type rule struct {
	pos         position
	name        string
	displayName string
	expr        interface{}
}

type choiceExpr struct {
	pos          position
	alternatives []interface{}
}

type actionExpr struct {
	pos  position
	expr interface{}
	run  func(*parser) (interface{}, error)
}

type seqExpr struct {
	pos   position
	exprs []interface{}
}

type labeledExpr struct {
	pos   position
	label string
	expr  interface{}
}

type expr struct {
	pos  position
	expr interface{}
}

type andExpr expr
type notExpr expr
type zeroOrOneExpr expr
type zeroOrMoreExpr expr
type oneOrMoreExpr expr

type ruleRefExpr struct {
	pos  position
	name string
}

type andCodeExpr struct {
	pos position
	run func(*parser) (bool, error)
}

type notCodeExpr struct {
	pos position
	run func(*parser) (bool, error)
}

type litMatcher struct {
	pos        position
	val        string
	ignoreCase bool
}

type charClassMatcher struct {
	pos        position
	val        string
	chars      []rune
	ranges     []rune
	classes    []*unicode.RangeTable
	ignoreCase bool
	inverted   bool
}

type anyMatcher position

// errList cumulates the errors found by the parser.
type errList []error

func (e *errList) add(err error) {
	*e = append(*e, err)
}

func (e errList) err() error {
	if len(e) == 0 {
		return nil
	}
	e.dedupe()
	return e
}

func (e *errList) dedupe() {
	var cleaned []error
	set := make(map[string]bool)
	for _, err := range *e {
		if msg := err.Error(); !set[msg] {
			set[msg] = true
			cleaned = append(cleaned, err)
		}
	}
	*e = cleaned
}

func (e errList) Error() string {
	switch len(e) {
	case 0:
		return ""
	case 1:
		return e[0].Error()
	default:
		var buf bytes.Buffer

		for i, err := range e {
			if i > 0 {
				buf.WriteRune('\n')
			}
			buf.WriteString(err.Error())
		}
		return buf.String()
	}
}

// parserError wraps an error with a prefix indicating the rule in which
// the error occurred. The original error is stored in the Inner field.
type parserError struct {
	Inner  error
	pos    position
	prefix string
}

// Error returns the error message.
func (p *parserError) Error() string {
	return p.prefix + ": " + p.Inner.Error()
}

// newParser creates a parser with the specified input source and options.
func newParser(filename string, b []byte, opts ...Option) *parser {
	p := &parser{
		filename: filename,
		errs:     new(errList),
		data:     b,
		pt:       savepoint{position: position{line: 1}},
		recover:  true,
	}
	p.setOptions(opts)
	return p
}

// setOptions applies the options to the parser.
func (p *parser) setOptions(opts []Option) {
	for _, opt := range opts {
		opt(p)
	}
}

type resultTuple struct {
	v   interface{}
	b   bool
	end savepoint
}

type parser struct {
	filename string
	pt       savepoint
	cur      current

	data []byte
	errs *errList

	recover bool
	debug   bool
	depth   int

	memoize bool
	// memoization table for the packrat algorithm:
	// map[offset in source] map[expression or rule] {value, match}
	memo map[int]map[interface{}]resultTuple

	// rules table, maps the rule identifier to the rule node
	rules map[string]*rule
	// variables stack, map of label to value
	vstack []map[string]interface{}
	// rule stack, allows identification of the current rule in errors
	rstack []*rule

	// stats
	exprCnt int
}

// push a variable set on the vstack.
func (p *parser) pushV() {
	if cap(p.vstack) == len(p.vstack) {
		// create new empty slot in the stack
		p.vstack = append(p.vstack, nil)
	} else {
		// slice to 1 more
		p.vstack = p.vstack[:len(p.vstack)+1]
	}

	// get the last args set
	m := p.vstack[len(p.vstack)-1]
	if m != nil && len(m) == 0 {
		// empty map, all good
		return
	}

	m = make(map[string]interface{})
	p.vstack[len(p.vstack)-1] = m
}

// pop a variable set from the vstack.
func (p *parser) popV() {
	// if the map is not empty, clear it
	m := p.vstack[len(p.vstack)-1]
	if len(m) > 0 {
		// GC that map
		p.vstack[len(p.vstack)-1] = nil
	}
	p.vstack = p.vstack[:len(p.vstack)-1]
}

func (p *parser) print(prefix, s string) string {
	if !p.debug {
		return s
	}

	fmt.Printf("%s %d:%d:%d: %s [%#U]\n",
		prefix, p.pt.line, p.pt.col, p.pt.offset, s, p.pt.rn)
	return s
}

func (p *parser) in(s string) string {
	p.depth++
	return p.print(strings.Repeat(" ", p.depth)+">", s)
}

func (p *parser) out(s string) string {
	p.depth--
	return p.print(strings.Repeat(" ", p.depth)+"<", s)
}

func (p *parser) addErr(err error) {
	p.addErrAt(err, p.pt.position)
}

func (p *parser) addErrAt(err error, pos position) {
	var buf bytes.Buffer
	if p.filename != "" {
		buf.WriteString(p.filename)
	}
	if buf.Len() > 0 {
		buf.WriteString(":")
	}
	buf.WriteString(fmt.Sprintf("%d:%d (%d)", pos.line, pos.col, pos.offset))
	if len(p.rstack) > 0 {
		if buf.Len() > 0 {
			buf.WriteString(": ")
		}
		rule := p.rstack[len(p.rstack)-1]
		if rule.displayName != "" {
			buf.WriteString("rule " + rule.displayName)
		} else {
			buf.WriteString("rule " + rule.name)
		}
	}
	pe := &parserError{Inner: err, prefix: buf.String()}
	p.errs.add(pe)
}

// read advances the parser to the next rune.
func (p *parser) read() {
	p.pt.offset += p.pt.w
	rn, n := utf8.DecodeRune(p.data[p.pt.offset:])
	p.pt.rn = rn
	p.pt.w = n
	p.pt.col++
	if rn == '\n' {
		p.pt.line++
		p.pt.col = 0
	}

	if rn == utf8.RuneError {
		if n > 0 {
			p.addErr(errInvalidEncoding)
		}
	}
}

// restore parser position to the savepoint pt.
func (p *parser) restore(pt savepoint) {
	if p.debug {
		defer p.out(p.in("restore"))
	}
	if pt.offset == p.pt.offset {
		return
	}
	p.pt = pt
}

// get the slice of bytes from the savepoint start to the current position.
func (p *parser) sliceFrom(start savepoint) []byte {
	return p.data[start.position.offset:p.pt.position.offset]
}

func (p *parser) getMemoized(node interface{}) (resultTuple, bool) {
	if len(p.memo) == 0 {
		return resultTuple{}, false
	}
	m := p.memo[p.pt.offset]
	if len(m) == 0 {
		return resultTuple{}, false
	}
	res, ok := m[node]
	return res, ok
}

func (p *parser) setMemoized(pt savepoint, node interface{}, tuple resultTuple) {
	if p.memo == nil {
		p.memo = make(map[int]map[interface{}]resultTuple)
	}
	m := p.memo[pt.offset]
	if m == nil {
		m = make(map[interface{}]resultTuple)
		p.memo[pt.offset] = m
	}
	m[node] = tuple
}

func (p *parser) buildRulesTable(g *grammar) {
	p.rules = make(map[string]*rule, len(g.rules))
	for _, r := range g.rules {
		p.rules[r.name] = r
	}
}

func (p *parser) parse(g *grammar) (val interface{}, err error) {
	if len(g.rules) == 0 {
		p.addErr(errNoRule)
		return nil, p.errs.err()
	}

	// TODO : not super critical but this could be generated
	p.buildRulesTable(g)

	if p.recover {
		// panic can be used in action code to stop parsing immediately
		// and return the panic as an error.
		defer func() {
			if e := recover(); e != nil {
				if p.debug {
					defer p.out(p.in("panic handler"))
				}
				val = nil
				switch e := e.(type) {
				case error:
					p.addErr(e)
				default:
					p.addErr(fmt.Errorf("%v", e))
				}
				err = p.errs.err()
			}
		}()
	}

	// start rule is rule [0]
	p.read() // advance to first rune
	val, ok := p.parseRule(g.rules[0])
	if !ok {
		if len(*p.errs) == 0 {
			// make sure this doesn't go out silently
			p.addErr(errNoMatch)
		}
		return nil, p.errs.err()
	}
	return val, p.errs.err()
}

func (p *parser) parseRule(rule *rule) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseRule " + rule.name))
	}

	if p.memoize {
		res, ok := p.getMemoized(rule)
		if ok {
			p.restore(res.end)
			return res.v, res.b
		}
	}

	start := p.pt
	p.rstack = append(p.rstack, rule)
	p.pushV()
	val, ok := p.parseExpr(rule.expr)
	p.popV()
	p.rstack = p.rstack[:len(p.rstack)-1]
	if ok && p.debug {
		p.print(strings.Repeat(" ", p.depth)+"MATCH", string(p.sliceFrom(start)))
	}

	if p.memoize {
		p.setMemoized(start, rule, resultTuple{val, ok, p.pt})
	}
	return val, ok
}

func (p *parser) parseExpr(expr interface{}) (interface{}, bool) {
	var pt savepoint
	var ok bool

	if p.memoize {
		res, ok := p.getMemoized(expr)
		if ok {
			p.restore(res.end)
			return res.v, res.b
		}
		pt = p.pt
	}

	p.exprCnt++
	var val interface{}
	switch expr := expr.(type) {
	case *actionExpr:
		val, ok = p.parseActionExpr(expr)
	case *andCodeExpr:
		val, ok = p.parseAndCodeExpr(expr)
	case *andExpr:
		val, ok = p.parseAndExpr(expr)
	case *anyMatcher:
		val, ok = p.parseAnyMatcher(expr)
	case *charClassMatcher:
		val, ok = p.parseCharClassMatcher(expr)
	case *choiceExpr:
		val, ok = p.parseChoiceExpr(expr)
	case *labeledExpr:
		val, ok = p.parseLabeledExpr(expr)
	case *litMatcher:
		val, ok = p.parseLitMatcher(expr)
	case *notCodeExpr:
		val, ok = p.parseNotCodeExpr(expr)
	case *notExpr:
		val, ok = p.parseNotExpr(expr)
	case *oneOrMoreExpr:
		val, ok = p.parseOneOrMoreExpr(expr)
	case *ruleRefExpr:
		val, ok = p.parseRuleRefExpr(expr)
	case *seqExpr:
		val, ok = p.parseSeqExpr(expr)
	case *zeroOrMoreExpr:
		val, ok = p.parseZeroOrMoreExpr(expr)
	case *zeroOrOneExpr:
		val, ok = p.parseZeroOrOneExpr(expr)
	default:
		panic(fmt.Sprintf("unknown expression type %T", expr))
	}
	if p.memoize {
		p.setMemoized(pt, expr, resultTuple{val, ok, p.pt})
	}
	return val, ok
}

func (p *parser) parseActionExpr(act *actionExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseActionExpr"))
	}

	start := p.pt
	val, ok := p.parseExpr(act.expr)
	if ok {
		p.cur.pos = start.position
		p.cur.text = p.sliceFrom(start)
		actVal, err := act.run(p)
		if err != nil {
			p.addErrAt(err, start.position)
		}
		val = actVal
	}
	if ok && p.debug {
		p.print(strings.Repeat(" ", p.depth)+"MATCH", string(p.sliceFrom(start)))
	}
	return val, ok
}

func (p *parser) parseAndCodeExpr(and *andCodeExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseAndCodeExpr"))
	}

	ok, err := and.run(p)
	if err != nil {
		p.addErr(err)
	}
	return nil, ok
}

func (p *parser) parseAndExpr(and *andExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseAndExpr"))
	}

	pt := p.pt
	p.pushV()
	_, ok := p.parseExpr(and.expr)
	p.popV()
	p.restore(pt)
	return nil, ok
}

func (p *parser) parseAnyMatcher(any *anyMatcher) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseAnyMatcher"))
	}

	if p.pt.rn != utf8.RuneError {
		start := p.pt
		p.read()
		return p.sliceFrom(start), true
	}
	return nil, false
}

func (p *parser) parseCharClassMatcher(chr *charClassMatcher) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseCharClassMatcher"))
	}

	cur := p.pt.rn
	// can't match EOF
	if cur == utf8.RuneError {
		return nil, false
	}
	start := p.pt
	if chr.ignoreCase {
		cur = unicode.ToLower(cur)
	}

	// try to match in the list of available chars
	for _, rn := range chr.chars {
		if rn == cur {
			if chr.inverted {
				return nil, false
			}
			p.read()
			return p.sliceFrom(start), true
		}
	}

	// try to match in the list of ranges
	for i := 0; i < len(chr.ranges); i += 2 {
		if cur >= chr.ranges[i] && cur <= chr.ranges[i+1] {
			if chr.inverted {
				return nil, false
			}
			p.read()
			return p.sliceFrom(start), true
		}
	}

	// try to match in the list of Unicode classes
	for _, cl := range chr.classes {
		if unicode.Is(cl, cur) {
			if chr.inverted {
				return nil, false
			}
			p.read()
			return p.sliceFrom(start), true
		}
	}

	if chr.inverted {
		p.read()
		return p.sliceFrom(start), true
	}
	return nil, false
}

func (p *parser) parseChoiceExpr(ch *choiceExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseChoiceExpr"))
	}

	for _, alt := range ch.alternatives {
		p.pushV()
		val, ok := p.parseExpr(alt)
		p.popV()
		if ok {
			return val, ok
		}
	}
	return nil, false
}

func (p *parser) parseLabeledExpr(lab *labeledExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseLabeledExpr"))
	}

	p.pushV()
	val, ok := p.parseExpr(lab.expr)
	p.popV()
	if ok && lab.label != "" {
		m := p.vstack[len(p.vstack)-1]
		m[lab.label] = val
	}
	return val, ok
}

func (p *parser) parseLitMatcher(lit *litMatcher) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseLitMatcher"))
	}

	start := p.pt
	for _, want := range lit.val {
		cur := p.pt.rn
		if lit.ignoreCase {
			cur = unicode.ToLower(cur)
		}
		if cur != want {
			p.restore(start)
			return nil, false
		}
		p.read()
	}
	return p.sliceFrom(start), true
}

func (p *parser) parseNotCodeExpr(not *notCodeExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseNotCodeExpr"))
	}

	ok, err := not.run(p)
	if err != nil {
		p.addErr(err)
	}
	return nil, !ok
}

func (p *parser) parseNotExpr(not *notExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseNotExpr"))
	}

	pt := p.pt
	p.pushV()
	_, ok := p.parseExpr(not.expr)
	p.popV()
	p.restore(pt)
	return nil, !ok
}

func (p *parser) parseOneOrMoreExpr(expr *oneOrMoreExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseOneOrMoreExpr"))
	}

	var vals []interface{}

	for {
		p.pushV()
		val, ok := p.parseExpr(expr.expr)
		p.popV()
		if !ok {
			if len(vals) == 0 {
				// did not match once, no match
				return nil, false
			}
			return vals, true
		}
		vals = append(vals, val)
	}
}

func (p *parser) parseRuleRefExpr(ref *ruleRefExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseRuleRefExpr " + ref.name))
	}

	if ref.name == "" {
		panic(fmt.Sprintf("%s: invalid rule: missing name", ref.pos))
	}

	rule := p.rules[ref.name]
	if rule == nil {
		p.addErr(fmt.Errorf("undefined rule: %s", ref.name))
		return nil, false
	}
	return p.parseRule(rule)
}

func (p *parser) parseSeqExpr(seq *seqExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseSeqExpr"))
	}

	var vals []interface{}

	pt := p.pt
	for _, expr := range seq.exprs {
		val, ok := p.parseExpr(expr)
		if !ok {
			p.restore(pt)
			return nil, false
		}
		vals = append(vals, val)
	}
	return vals, true
}

func (p *parser) parseZeroOrMoreExpr(expr *zeroOrMoreExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseZeroOrMoreExpr"))
	}

	var vals []interface{}

	for {
		p.pushV()
		val, ok := p.parseExpr(expr.expr)
		p.popV()
		if !ok {
			return vals, true
		}
		vals = append(vals, val)
	}
}

func (p *parser) parseZeroOrOneExpr(expr *zeroOrOneExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseZeroOrOneExpr"))
	}

	p.pushV()
	val, _ := p.parseExpr(expr.expr)
	p.popV()
	// whether it matched or not, consider it a match
	return val, true
}

func rangeTable(class string) *unicode.RangeTable {
	if rt, ok := unicode.Categories[class]; ok {
		return rt
	}
	if rt, ok := unicode.Properties[class]; ok {
		return rt
	}
	if rt, ok := unicode.Scripts[class]; ok {
		return rt
	}

	// cannot happen
	panic(fmt.Sprintf("invalid Unicode class: %s", class))
}
