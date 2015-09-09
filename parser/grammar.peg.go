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

var g = &grammar{
	rules: []*rule{
		{
			name: "Grammar",
			pos:  position{line: 43, col: 1, offset: 534},
			expr: &actionExpr{
				pos: position{line: 43, col: 11, offset: 546},
				run: (*parser).callonGrammar1,
				expr: &seqExpr{
					pos: position{line: 43, col: 11, offset: 546},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 43, col: 11, offset: 546},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 43, col: 14, offset: 549},
							label: "statements",
							expr: &zeroOrMoreExpr{
								pos: position{line: 43, col: 25, offset: 560},
								expr: &seqExpr{
									pos: position{line: 43, col: 27, offset: 562},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 43, col: 27, offset: 562},
											name: "Statement",
										},
										&ruleRefExpr{
											pos:  position{line: 43, col: 37, offset: 572},
											name: "__",
										},
									},
								},
							},
						},
						&choiceExpr{
							pos: position{line: 43, col: 44, offset: 579},
							alternatives: []interface{}{
								&ruleRefExpr{
									pos:  position{line: 43, col: 44, offset: 579},
									name: "EOF",
								},
								&ruleRefExpr{
									pos:  position{line: 43, col: 50, offset: 585},
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
			pos:  position{line: 87, col: 1, offset: 1736},
			expr: &actionExpr{
				pos: position{line: 87, col: 15, offset: 1752},
				run: (*parser).callonSyntaxError1,
				expr: &anyMatcher{
					line: 87, col: 15, offset: 1752,
				},
			},
		},
		{
			name: "Include",
			pos:  position{line: 91, col: 1, offset: 1807},
			expr: &actionExpr{
				pos: position{line: 91, col: 11, offset: 1819},
				run: (*parser).callonInclude1,
				expr: &seqExpr{
					pos: position{line: 91, col: 11, offset: 1819},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 91, col: 11, offset: 1819},
							val:        "include",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 91, col: 21, offset: 1829},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 91, col: 23, offset: 1831},
							label: "file",
							expr: &ruleRefExpr{
								pos:  position{line: 91, col: 28, offset: 1836},
								name: "Literal",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 91, col: 36, offset: 1844},
							name: "EOS",
						},
					},
				},
			},
		},
		{
			name: "Statement",
			pos:  position{line: 95, col: 1, offset: 1889},
			expr: &choiceExpr{
				pos: position{line: 95, col: 13, offset: 1903},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 95, col: 13, offset: 1903},
						name: "Include",
					},
					&ruleRefExpr{
						pos:  position{line: 95, col: 23, offset: 1913},
						name: "Namespace",
					},
					&ruleRefExpr{
						pos:  position{line: 95, col: 35, offset: 1925},
						name: "Const",
					},
					&ruleRefExpr{
						pos:  position{line: 95, col: 43, offset: 1933},
						name: "Enum",
					},
					&ruleRefExpr{
						pos:  position{line: 95, col: 50, offset: 1940},
						name: "TypeDef",
					},
					&ruleRefExpr{
						pos:  position{line: 95, col: 60, offset: 1950},
						name: "Struct",
					},
					&ruleRefExpr{
						pos:  position{line: 95, col: 69, offset: 1959},
						name: "Exception",
					},
					&ruleRefExpr{
						pos:  position{line: 95, col: 81, offset: 1971},
						name: "Union",
					},
					&ruleRefExpr{
						pos:  position{line: 95, col: 89, offset: 1979},
						name: "Service",
					},
				},
			},
		},
		{
			name: "Namespace",
			pos:  position{line: 97, col: 1, offset: 1988},
			expr: &actionExpr{
				pos: position{line: 97, col: 13, offset: 2002},
				run: (*parser).callonNamespace1,
				expr: &seqExpr{
					pos: position{line: 97, col: 13, offset: 2002},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 97, col: 13, offset: 2002},
							val:        "namespace",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 97, col: 25, offset: 2014},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 97, col: 27, offset: 2016},
							label: "scope",
							expr: &oneOrMoreExpr{
								pos: position{line: 97, col: 33, offset: 2022},
								expr: &charClassMatcher{
									pos:        position{line: 97, col: 33, offset: 2022},
									val:        "[a-z.-]",
									chars:      []rune{'.', '-'},
									ranges:     []rune{'a', 'z'},
									ignoreCase: false,
									inverted:   false,
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 97, col: 42, offset: 2031},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 97, col: 44, offset: 2033},
							label: "ns",
							expr: &ruleRefExpr{
								pos:  position{line: 97, col: 47, offset: 2036},
								name: "Identifier",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 97, col: 58, offset: 2047},
							name: "EOS",
						},
					},
				},
			},
		},
		{
			name: "Const",
			pos:  position{line: 104, col: 1, offset: 2158},
			expr: &actionExpr{
				pos: position{line: 104, col: 9, offset: 2168},
				run: (*parser).callonConst1,
				expr: &seqExpr{
					pos: position{line: 104, col: 9, offset: 2168},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 104, col: 9, offset: 2168},
							val:        "const",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 104, col: 17, offset: 2176},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 104, col: 19, offset: 2178},
							label: "typ",
							expr: &ruleRefExpr{
								pos:  position{line: 104, col: 23, offset: 2182},
								name: "FieldType",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 104, col: 33, offset: 2192},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 104, col: 35, offset: 2194},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 104, col: 40, offset: 2199},
								name: "Identifier",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 104, col: 51, offset: 2210},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 104, col: 53, offset: 2212},
							val:        "=",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 104, col: 57, offset: 2216},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 104, col: 59, offset: 2218},
							label: "value",
							expr: &ruleRefExpr{
								pos:  position{line: 104, col: 65, offset: 2224},
								name: "ConstValue",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 104, col: 76, offset: 2235},
							name: "EOS",
						},
					},
				},
			},
		},
		{
			name: "Enum",
			pos:  position{line: 112, col: 1, offset: 2343},
			expr: &actionExpr{
				pos: position{line: 112, col: 8, offset: 2352},
				run: (*parser).callonEnum1,
				expr: &seqExpr{
					pos: position{line: 112, col: 8, offset: 2352},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 112, col: 8, offset: 2352},
							val:        "enum",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 112, col: 15, offset: 2359},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 112, col: 17, offset: 2361},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 112, col: 22, offset: 2366},
								name: "Identifier",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 112, col: 33, offset: 2377},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 112, col: 35, offset: 2379},
							val:        "{",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 112, col: 39, offset: 2383},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 112, col: 42, offset: 2386},
							label: "values",
							expr: &zeroOrMoreExpr{
								pos: position{line: 112, col: 49, offset: 2393},
								expr: &seqExpr{
									pos: position{line: 112, col: 50, offset: 2394},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 112, col: 50, offset: 2394},
											name: "EnumValue",
										},
										&ruleRefExpr{
											pos:  position{line: 112, col: 60, offset: 2404},
											name: "__",
										},
									},
								},
							},
						},
						&litMatcher{
							pos:        position{line: 112, col: 65, offset: 2409},
							val:        "}",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 112, col: 69, offset: 2413},
							name: "EOS",
						},
					},
				},
			},
		},
		{
			name: "EnumValue",
			pos:  position{line: 135, col: 1, offset: 2929},
			expr: &actionExpr{
				pos: position{line: 135, col: 13, offset: 2943},
				run: (*parser).callonEnumValue1,
				expr: &seqExpr{
					pos: position{line: 135, col: 13, offset: 2943},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 135, col: 13, offset: 2943},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 135, col: 18, offset: 2948},
								name: "Identifier",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 135, col: 29, offset: 2959},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 135, col: 31, offset: 2961},
							label: "value",
							expr: &zeroOrOneExpr{
								pos: position{line: 135, col: 37, offset: 2967},
								expr: &seqExpr{
									pos: position{line: 135, col: 38, offset: 2968},
									exprs: []interface{}{
										&litMatcher{
											pos:        position{line: 135, col: 38, offset: 2968},
											val:        "=",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 135, col: 42, offset: 2972},
											name: "_",
										},
										&ruleRefExpr{
											pos:  position{line: 135, col: 44, offset: 2974},
											name: "IntConstant",
										},
									},
								},
							},
						},
						&zeroOrOneExpr{
							pos: position{line: 135, col: 58, offset: 2988},
							expr: &ruleRefExpr{
								pos:  position{line: 135, col: 58, offset: 2988},
								name: "ListSeparator",
							},
						},
					},
				},
			},
		},
		{
			name: "TypeDef",
			pos:  position{line: 146, col: 1, offset: 3167},
			expr: &actionExpr{
				pos: position{line: 146, col: 11, offset: 3179},
				run: (*parser).callonTypeDef1,
				expr: &seqExpr{
					pos: position{line: 146, col: 11, offset: 3179},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 146, col: 11, offset: 3179},
							val:        "typedef",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 146, col: 21, offset: 3189},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 146, col: 23, offset: 3191},
							label: "typ",
							expr: &ruleRefExpr{
								pos:  position{line: 146, col: 27, offset: 3195},
								name: "FieldType",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 146, col: 37, offset: 3205},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 146, col: 39, offset: 3207},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 146, col: 44, offset: 3212},
								name: "Identifier",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 146, col: 55, offset: 3223},
							name: "EOS",
						},
					},
				},
			},
		},
		{
			name: "Struct",
			pos:  position{line: 153, col: 1, offset: 3313},
			expr: &actionExpr{
				pos: position{line: 153, col: 10, offset: 3324},
				run: (*parser).callonStruct1,
				expr: &seqExpr{
					pos: position{line: 153, col: 10, offset: 3324},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 153, col: 10, offset: 3324},
							val:        "struct",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 153, col: 19, offset: 3333},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 153, col: 21, offset: 3335},
							label: "st",
							expr: &ruleRefExpr{
								pos:  position{line: 153, col: 24, offset: 3338},
								name: "StructLike",
							},
						},
					},
				},
			},
		},
		{
			name: "Exception",
			pos:  position{line: 154, col: 1, offset: 3378},
			expr: &actionExpr{
				pos: position{line: 154, col: 13, offset: 3392},
				run: (*parser).callonException1,
				expr: &seqExpr{
					pos: position{line: 154, col: 13, offset: 3392},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 154, col: 13, offset: 3392},
							val:        "exception",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 154, col: 25, offset: 3404},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 154, col: 27, offset: 3406},
							label: "st",
							expr: &ruleRefExpr{
								pos:  position{line: 154, col: 30, offset: 3409},
								name: "StructLike",
							},
						},
					},
				},
			},
		},
		{
			name: "Union",
			pos:  position{line: 155, col: 1, offset: 3460},
			expr: &actionExpr{
				pos: position{line: 155, col: 9, offset: 3470},
				run: (*parser).callonUnion1,
				expr: &seqExpr{
					pos: position{line: 155, col: 9, offset: 3470},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 155, col: 9, offset: 3470},
							val:        "union",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 155, col: 17, offset: 3478},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 155, col: 19, offset: 3480},
							label: "st",
							expr: &ruleRefExpr{
								pos:  position{line: 155, col: 22, offset: 3483},
								name: "StructLike",
							},
						},
					},
				},
			},
		},
		{
			name: "StructLike",
			pos:  position{line: 156, col: 1, offset: 3530},
			expr: &actionExpr{
				pos: position{line: 156, col: 14, offset: 3545},
				run: (*parser).callonStructLike1,
				expr: &seqExpr{
					pos: position{line: 156, col: 14, offset: 3545},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 156, col: 14, offset: 3545},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 156, col: 19, offset: 3550},
								name: "Identifier",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 156, col: 30, offset: 3561},
							name: "__",
						},
						&litMatcher{
							pos:        position{line: 156, col: 33, offset: 3564},
							val:        "{",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 156, col: 37, offset: 3568},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 156, col: 40, offset: 3571},
							label: "fields",
							expr: &ruleRefExpr{
								pos:  position{line: 156, col: 47, offset: 3578},
								name: "FieldList",
							},
						},
						&litMatcher{
							pos:        position{line: 156, col: 57, offset: 3588},
							val:        "}",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 156, col: 61, offset: 3592},
							name: "EOS",
						},
					},
				},
			},
		},
		{
			name: "FieldList",
			pos:  position{line: 166, col: 1, offset: 3726},
			expr: &actionExpr{
				pos: position{line: 166, col: 13, offset: 3740},
				run: (*parser).callonFieldList1,
				expr: &labeledExpr{
					pos:   position{line: 166, col: 13, offset: 3740},
					label: "fields",
					expr: &zeroOrMoreExpr{
						pos: position{line: 166, col: 20, offset: 3747},
						expr: &seqExpr{
							pos: position{line: 166, col: 21, offset: 3748},
							exprs: []interface{}{
								&ruleRefExpr{
									pos:  position{line: 166, col: 21, offset: 3748},
									name: "Field",
								},
								&ruleRefExpr{
									pos:  position{line: 166, col: 27, offset: 3754},
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
			pos:  position{line: 175, col: 1, offset: 3914},
			expr: &actionExpr{
				pos: position{line: 175, col: 9, offset: 3924},
				run: (*parser).callonField1,
				expr: &seqExpr{
					pos: position{line: 175, col: 9, offset: 3924},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 175, col: 9, offset: 3924},
							label: "id",
							expr: &ruleRefExpr{
								pos:  position{line: 175, col: 12, offset: 3927},
								name: "IntConstant",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 175, col: 24, offset: 3939},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 175, col: 26, offset: 3941},
							val:        ":",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 175, col: 30, offset: 3945},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 175, col: 32, offset: 3947},
							label: "req",
							expr: &zeroOrOneExpr{
								pos: position{line: 175, col: 36, offset: 3951},
								expr: &ruleRefExpr{
									pos:  position{line: 175, col: 36, offset: 3951},
									name: "FieldReq",
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 175, col: 46, offset: 3961},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 175, col: 48, offset: 3963},
							label: "typ",
							expr: &ruleRefExpr{
								pos:  position{line: 175, col: 52, offset: 3967},
								name: "FieldType",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 175, col: 62, offset: 3977},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 175, col: 64, offset: 3979},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 175, col: 69, offset: 3984},
								name: "Identifier",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 175, col: 80, offset: 3995},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 175, col: 82, offset: 3997},
							label: "def",
							expr: &zeroOrOneExpr{
								pos: position{line: 175, col: 86, offset: 4001},
								expr: &seqExpr{
									pos: position{line: 175, col: 87, offset: 4002},
									exprs: []interface{}{
										&litMatcher{
											pos:        position{line: 175, col: 87, offset: 4002},
											val:        "=",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 175, col: 91, offset: 4006},
											name: "_",
										},
										&ruleRefExpr{
											pos:  position{line: 175, col: 93, offset: 4008},
											name: "ConstValue",
										},
									},
								},
							},
						},
						&zeroOrOneExpr{
							pos: position{line: 175, col: 106, offset: 4021},
							expr: &ruleRefExpr{
								pos:  position{line: 175, col: 106, offset: 4021},
								name: "ListSeparator",
							},
						},
					},
				},
			},
		},
		{
			name: "FieldReq",
			pos:  position{line: 190, col: 1, offset: 4281},
			expr: &actionExpr{
				pos: position{line: 190, col: 12, offset: 4294},
				run: (*parser).callonFieldReq1,
				expr: &choiceExpr{
					pos: position{line: 190, col: 13, offset: 4295},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 190, col: 13, offset: 4295},
							val:        "required",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 190, col: 26, offset: 4308},
							val:        "optional",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "Service",
			pos:  position{line: 194, col: 1, offset: 4379},
			expr: &actionExpr{
				pos: position{line: 194, col: 11, offset: 4391},
				run: (*parser).callonService1,
				expr: &seqExpr{
					pos: position{line: 194, col: 11, offset: 4391},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 194, col: 11, offset: 4391},
							val:        "service",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 194, col: 21, offset: 4401},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 194, col: 23, offset: 4403},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 194, col: 28, offset: 4408},
								name: "Identifier",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 194, col: 39, offset: 4419},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 194, col: 41, offset: 4421},
							label: "extends",
							expr: &zeroOrOneExpr{
								pos: position{line: 194, col: 49, offset: 4429},
								expr: &seqExpr{
									pos: position{line: 194, col: 50, offset: 4430},
									exprs: []interface{}{
										&litMatcher{
											pos:        position{line: 194, col: 50, offset: 4430},
											val:        "extends",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 194, col: 60, offset: 4440},
											name: "__",
										},
										&ruleRefExpr{
											pos:  position{line: 194, col: 63, offset: 4443},
											name: "Identifier",
										},
										&ruleRefExpr{
											pos:  position{line: 194, col: 74, offset: 4454},
											name: "__",
										},
									},
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 194, col: 79, offset: 4459},
							name: "__",
						},
						&litMatcher{
							pos:        position{line: 194, col: 82, offset: 4462},
							val:        "{",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 194, col: 86, offset: 4466},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 194, col: 89, offset: 4469},
							label: "methods",
							expr: &zeroOrMoreExpr{
								pos: position{line: 194, col: 97, offset: 4477},
								expr: &seqExpr{
									pos: position{line: 194, col: 98, offset: 4478},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 194, col: 98, offset: 4478},
											name: "Function",
										},
										&ruleRefExpr{
											pos:  position{line: 194, col: 107, offset: 4487},
											name: "__",
										},
									},
								},
							},
						},
						&choiceExpr{
							pos: position{line: 194, col: 113, offset: 4493},
							alternatives: []interface{}{
								&litMatcher{
									pos:        position{line: 194, col: 113, offset: 4493},
									val:        "}",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 194, col: 119, offset: 4499},
									name: "EndOfServiceError",
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 194, col: 138, offset: 4518},
							name: "EOS",
						},
					},
				},
			},
		},
		{
			name: "EndOfServiceError",
			pos:  position{line: 209, col: 1, offset: 4859},
			expr: &actionExpr{
				pos: position{line: 209, col: 21, offset: 4881},
				run: (*parser).callonEndOfServiceError1,
				expr: &anyMatcher{
					line: 209, col: 21, offset: 4881,
				},
			},
		},
		{
			name: "Function",
			pos:  position{line: 213, col: 1, offset: 4947},
			expr: &actionExpr{
				pos: position{line: 213, col: 12, offset: 4960},
				run: (*parser).callonFunction1,
				expr: &seqExpr{
					pos: position{line: 213, col: 12, offset: 4960},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 213, col: 12, offset: 4960},
							label: "oneway",
							expr: &zeroOrOneExpr{
								pos: position{line: 213, col: 19, offset: 4967},
								expr: &seqExpr{
									pos: position{line: 213, col: 20, offset: 4968},
									exprs: []interface{}{
										&litMatcher{
											pos:        position{line: 213, col: 20, offset: 4968},
											val:        "oneway",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 213, col: 29, offset: 4977},
											name: "__",
										},
									},
								},
							},
						},
						&labeledExpr{
							pos:   position{line: 213, col: 34, offset: 4982},
							label: "typ",
							expr: &ruleRefExpr{
								pos:  position{line: 213, col: 38, offset: 4986},
								name: "FunctionType",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 213, col: 51, offset: 4999},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 213, col: 54, offset: 5002},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 213, col: 59, offset: 5007},
								name: "Identifier",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 213, col: 70, offset: 5018},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 213, col: 72, offset: 5020},
							val:        "(",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 213, col: 76, offset: 5024},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 213, col: 79, offset: 5027},
							label: "arguments",
							expr: &ruleRefExpr{
								pos:  position{line: 213, col: 89, offset: 5037},
								name: "FieldList",
							},
						},
						&litMatcher{
							pos:        position{line: 213, col: 99, offset: 5047},
							val:        ")",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 213, col: 103, offset: 5051},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 213, col: 106, offset: 5054},
							label: "exceptions",
							expr: &zeroOrOneExpr{
								pos: position{line: 213, col: 117, offset: 5065},
								expr: &ruleRefExpr{
									pos:  position{line: 213, col: 117, offset: 5065},
									name: "Throws",
								},
							},
						},
						&zeroOrOneExpr{
							pos: position{line: 213, col: 125, offset: 5073},
							expr: &ruleRefExpr{
								pos:  position{line: 213, col: 125, offset: 5073},
								name: "ListSeparator",
							},
						},
					},
				},
			},
		},
		{
			name: "FunctionType",
			pos:  position{line: 236, col: 1, offset: 5454},
			expr: &actionExpr{
				pos: position{line: 236, col: 16, offset: 5471},
				run: (*parser).callonFunctionType1,
				expr: &labeledExpr{
					pos:   position{line: 236, col: 16, offset: 5471},
					label: "typ",
					expr: &choiceExpr{
						pos: position{line: 236, col: 21, offset: 5476},
						alternatives: []interface{}{
							&litMatcher{
								pos:        position{line: 236, col: 21, offset: 5476},
								val:        "void",
								ignoreCase: false,
							},
							&ruleRefExpr{
								pos:  position{line: 236, col: 30, offset: 5485},
								name: "FieldType",
							},
						},
					},
				},
			},
		},
		{
			name: "Throws",
			pos:  position{line: 243, col: 1, offset: 5592},
			expr: &actionExpr{
				pos: position{line: 243, col: 10, offset: 5603},
				run: (*parser).callonThrows1,
				expr: &seqExpr{
					pos: position{line: 243, col: 10, offset: 5603},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 243, col: 10, offset: 5603},
							val:        "throws",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 243, col: 19, offset: 5612},
							name: "__",
						},
						&litMatcher{
							pos:        position{line: 243, col: 22, offset: 5615},
							val:        "(",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 243, col: 26, offset: 5619},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 243, col: 29, offset: 5622},
							label: "exceptions",
							expr: &ruleRefExpr{
								pos:  position{line: 243, col: 40, offset: 5633},
								name: "FieldList",
							},
						},
						&litMatcher{
							pos:        position{line: 243, col: 50, offset: 5643},
							val:        ")",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "FieldType",
			pos:  position{line: 247, col: 1, offset: 5676},
			expr: &actionExpr{
				pos: position{line: 247, col: 13, offset: 5690},
				run: (*parser).callonFieldType1,
				expr: &labeledExpr{
					pos:   position{line: 247, col: 13, offset: 5690},
					label: "typ",
					expr: &choiceExpr{
						pos: position{line: 247, col: 18, offset: 5695},
						alternatives: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 247, col: 18, offset: 5695},
								name: "BaseType",
							},
							&ruleRefExpr{
								pos:  position{line: 247, col: 29, offset: 5706},
								name: "ContainerType",
							},
							&ruleRefExpr{
								pos:  position{line: 247, col: 45, offset: 5722},
								name: "Identifier",
							},
						},
					},
				},
			},
		},
		{
			name: "DefinitionType",
			pos:  position{line: 254, col: 1, offset: 5832},
			expr: &actionExpr{
				pos: position{line: 254, col: 18, offset: 5851},
				run: (*parser).callonDefinitionType1,
				expr: &labeledExpr{
					pos:   position{line: 254, col: 18, offset: 5851},
					label: "typ",
					expr: &choiceExpr{
						pos: position{line: 254, col: 23, offset: 5856},
						alternatives: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 254, col: 23, offset: 5856},
								name: "BaseType",
							},
							&ruleRefExpr{
								pos:  position{line: 254, col: 34, offset: 5867},
								name: "ContainerType",
							},
						},
					},
				},
			},
		},
		{
			name: "BaseType",
			pos:  position{line: 258, col: 1, offset: 5904},
			expr: &actionExpr{
				pos: position{line: 258, col: 12, offset: 5917},
				run: (*parser).callonBaseType1,
				expr: &choiceExpr{
					pos: position{line: 258, col: 13, offset: 5918},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 258, col: 13, offset: 5918},
							val:        "bool",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 258, col: 22, offset: 5927},
							val:        "byte",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 258, col: 31, offset: 5936},
							val:        "i16",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 258, col: 39, offset: 5944},
							val:        "i32",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 258, col: 47, offset: 5952},
							val:        "i64",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 258, col: 55, offset: 5960},
							val:        "double",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 258, col: 66, offset: 5971},
							val:        "string",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 258, col: 77, offset: 5982},
							val:        "binary",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "ContainerType",
			pos:  position{line: 262, col: 1, offset: 6039},
			expr: &actionExpr{
				pos: position{line: 262, col: 17, offset: 6057},
				run: (*parser).callonContainerType1,
				expr: &labeledExpr{
					pos:   position{line: 262, col: 17, offset: 6057},
					label: "typ",
					expr: &choiceExpr{
						pos: position{line: 262, col: 22, offset: 6062},
						alternatives: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 262, col: 22, offset: 6062},
								name: "MapType",
							},
							&ruleRefExpr{
								pos:  position{line: 262, col: 32, offset: 6072},
								name: "SetType",
							},
							&ruleRefExpr{
								pos:  position{line: 262, col: 42, offset: 6082},
								name: "ListType",
							},
						},
					},
				},
			},
		},
		{
			name: "MapType",
			pos:  position{line: 266, col: 1, offset: 6114},
			expr: &actionExpr{
				pos: position{line: 266, col: 11, offset: 6126},
				run: (*parser).callonMapType1,
				expr: &seqExpr{
					pos: position{line: 266, col: 11, offset: 6126},
					exprs: []interface{}{
						&zeroOrOneExpr{
							pos: position{line: 266, col: 11, offset: 6126},
							expr: &ruleRefExpr{
								pos:  position{line: 266, col: 11, offset: 6126},
								name: "CppType",
							},
						},
						&litMatcher{
							pos:        position{line: 266, col: 20, offset: 6135},
							val:        "map<",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 266, col: 27, offset: 6142},
							name: "WS",
						},
						&labeledExpr{
							pos:   position{line: 266, col: 30, offset: 6145},
							label: "key",
							expr: &ruleRefExpr{
								pos:  position{line: 266, col: 34, offset: 6149},
								name: "FieldType",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 266, col: 44, offset: 6159},
							name: "WS",
						},
						&litMatcher{
							pos:        position{line: 266, col: 47, offset: 6162},
							val:        ",",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 266, col: 51, offset: 6166},
							name: "WS",
						},
						&labeledExpr{
							pos:   position{line: 266, col: 54, offset: 6169},
							label: "value",
							expr: &ruleRefExpr{
								pos:  position{line: 266, col: 60, offset: 6175},
								name: "FieldType",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 266, col: 70, offset: 6185},
							name: "WS",
						},
						&litMatcher{
							pos:        position{line: 266, col: 73, offset: 6188},
							val:        ">",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "SetType",
			pos:  position{line: 274, col: 1, offset: 6287},
			expr: &actionExpr{
				pos: position{line: 274, col: 11, offset: 6299},
				run: (*parser).callonSetType1,
				expr: &seqExpr{
					pos: position{line: 274, col: 11, offset: 6299},
					exprs: []interface{}{
						&zeroOrOneExpr{
							pos: position{line: 274, col: 11, offset: 6299},
							expr: &ruleRefExpr{
								pos:  position{line: 274, col: 11, offset: 6299},
								name: "CppType",
							},
						},
						&litMatcher{
							pos:        position{line: 274, col: 20, offset: 6308},
							val:        "set<",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 274, col: 27, offset: 6315},
							name: "WS",
						},
						&labeledExpr{
							pos:   position{line: 274, col: 30, offset: 6318},
							label: "typ",
							expr: &ruleRefExpr{
								pos:  position{line: 274, col: 34, offset: 6322},
								name: "FieldType",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 274, col: 44, offset: 6332},
							name: "WS",
						},
						&litMatcher{
							pos:        position{line: 274, col: 47, offset: 6335},
							val:        ">",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "ListType",
			pos:  position{line: 281, col: 1, offset: 6408},
			expr: &actionExpr{
				pos: position{line: 281, col: 12, offset: 6421},
				run: (*parser).callonListType1,
				expr: &seqExpr{
					pos: position{line: 281, col: 12, offset: 6421},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 281, col: 12, offset: 6421},
							val:        "list<",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 281, col: 20, offset: 6429},
							name: "WS",
						},
						&labeledExpr{
							pos:   position{line: 281, col: 23, offset: 6432},
							label: "typ",
							expr: &ruleRefExpr{
								pos:  position{line: 281, col: 27, offset: 6436},
								name: "FieldType",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 281, col: 37, offset: 6446},
							name: "WS",
						},
						&litMatcher{
							pos:        position{line: 281, col: 40, offset: 6449},
							val:        ">",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "CppType",
			pos:  position{line: 288, col: 1, offset: 6523},
			expr: &actionExpr{
				pos: position{line: 288, col: 11, offset: 6535},
				run: (*parser).callonCppType1,
				expr: &seqExpr{
					pos: position{line: 288, col: 11, offset: 6535},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 288, col: 11, offset: 6535},
							val:        "cpp_type",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 288, col: 22, offset: 6546},
							label: "cppType",
							expr: &ruleRefExpr{
								pos:  position{line: 288, col: 30, offset: 6554},
								name: "Literal",
							},
						},
					},
				},
			},
		},
		{
			name: "ConstValue",
			pos:  position{line: 292, col: 1, offset: 6588},
			expr: &choiceExpr{
				pos: position{line: 292, col: 14, offset: 6603},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 292, col: 14, offset: 6603},
						name: "Literal",
					},
					&ruleRefExpr{
						pos:  position{line: 292, col: 24, offset: 6613},
						name: "DoubleConstant",
					},
					&ruleRefExpr{
						pos:  position{line: 292, col: 41, offset: 6630},
						name: "IntConstant",
					},
					&ruleRefExpr{
						pos:  position{line: 292, col: 55, offset: 6644},
						name: "ConstMap",
					},
					&ruleRefExpr{
						pos:  position{line: 292, col: 66, offset: 6655},
						name: "ConstList",
					},
					&ruleRefExpr{
						pos:  position{line: 292, col: 78, offset: 6667},
						name: "Identifier",
					},
				},
			},
		},
		{
			name: "IntConstant",
			pos:  position{line: 294, col: 1, offset: 6679},
			expr: &actionExpr{
				pos: position{line: 294, col: 15, offset: 6695},
				run: (*parser).callonIntConstant1,
				expr: &seqExpr{
					pos: position{line: 294, col: 15, offset: 6695},
					exprs: []interface{}{
						&zeroOrOneExpr{
							pos: position{line: 294, col: 15, offset: 6695},
							expr: &charClassMatcher{
								pos:        position{line: 294, col: 15, offset: 6695},
								val:        "[-+]",
								chars:      []rune{'-', '+'},
								ignoreCase: false,
								inverted:   false,
							},
						},
						&oneOrMoreExpr{
							pos: position{line: 294, col: 21, offset: 6701},
							expr: &ruleRefExpr{
								pos:  position{line: 294, col: 21, offset: 6701},
								name: "Digit",
							},
						},
					},
				},
			},
		},
		{
			name: "DoubleConstant",
			pos:  position{line: 298, col: 1, offset: 6762},
			expr: &actionExpr{
				pos: position{line: 298, col: 18, offset: 6781},
				run: (*parser).callonDoubleConstant1,
				expr: &seqExpr{
					pos: position{line: 298, col: 18, offset: 6781},
					exprs: []interface{}{
						&zeroOrOneExpr{
							pos: position{line: 298, col: 18, offset: 6781},
							expr: &charClassMatcher{
								pos:        position{line: 298, col: 18, offset: 6781},
								val:        "[+-]",
								chars:      []rune{'+', '-'},
								ignoreCase: false,
								inverted:   false,
							},
						},
						&zeroOrMoreExpr{
							pos: position{line: 298, col: 24, offset: 6787},
							expr: &ruleRefExpr{
								pos:  position{line: 298, col: 24, offset: 6787},
								name: "Digit",
							},
						},
						&litMatcher{
							pos:        position{line: 298, col: 31, offset: 6794},
							val:        ".",
							ignoreCase: false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 298, col: 35, offset: 6798},
							expr: &ruleRefExpr{
								pos:  position{line: 298, col: 35, offset: 6798},
								name: "Digit",
							},
						},
						&zeroOrOneExpr{
							pos: position{line: 298, col: 42, offset: 6805},
							expr: &seqExpr{
								pos: position{line: 298, col: 44, offset: 6807},
								exprs: []interface{}{
									&charClassMatcher{
										pos:        position{line: 298, col: 44, offset: 6807},
										val:        "['Ee']",
										chars:      []rune{'\'', 'E', 'e', '\''},
										ignoreCase: false,
										inverted:   false,
									},
									&ruleRefExpr{
										pos:  position{line: 298, col: 51, offset: 6814},
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
			pos:  position{line: 302, col: 1, offset: 6881},
			expr: &actionExpr{
				pos: position{line: 302, col: 13, offset: 6895},
				run: (*parser).callonConstList1,
				expr: &seqExpr{
					pos: position{line: 302, col: 13, offset: 6895},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 302, col: 13, offset: 6895},
							val:        "[",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 302, col: 17, offset: 6899},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 302, col: 20, offset: 6902},
							label: "values",
							expr: &zeroOrMoreExpr{
								pos: position{line: 302, col: 27, offset: 6909},
								expr: &seqExpr{
									pos: position{line: 302, col: 28, offset: 6910},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 302, col: 28, offset: 6910},
											name: "ConstValue",
										},
										&ruleRefExpr{
											pos:  position{line: 302, col: 39, offset: 6921},
											name: "__",
										},
										&zeroOrOneExpr{
											pos: position{line: 302, col: 42, offset: 6924},
											expr: &ruleRefExpr{
												pos:  position{line: 302, col: 42, offset: 6924},
												name: "ListSeparator",
											},
										},
										&ruleRefExpr{
											pos:  position{line: 302, col: 57, offset: 6939},
											name: "__",
										},
									},
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 302, col: 62, offset: 6944},
							name: "__",
						},
						&litMatcher{
							pos:        position{line: 302, col: 65, offset: 6947},
							val:        "]",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "ConstMap",
			pos:  position{line: 311, col: 1, offset: 7120},
			expr: &actionExpr{
				pos: position{line: 311, col: 12, offset: 7133},
				run: (*parser).callonConstMap1,
				expr: &seqExpr{
					pos: position{line: 311, col: 12, offset: 7133},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 311, col: 12, offset: 7133},
							val:        "{",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 311, col: 16, offset: 7137},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 311, col: 19, offset: 7140},
							label: "values",
							expr: &zeroOrMoreExpr{
								pos: position{line: 311, col: 26, offset: 7147},
								expr: &seqExpr{
									pos: position{line: 311, col: 27, offset: 7148},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 311, col: 27, offset: 7148},
											name: "ConstValue",
										},
										&ruleRefExpr{
											pos:  position{line: 311, col: 38, offset: 7159},
											name: "__",
										},
										&litMatcher{
											pos:        position{line: 311, col: 41, offset: 7162},
											val:        ":",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 311, col: 45, offset: 7166},
											name: "__",
										},
										&ruleRefExpr{
											pos:  position{line: 311, col: 48, offset: 7169},
											name: "ConstValue",
										},
										&ruleRefExpr{
											pos:  position{line: 311, col: 59, offset: 7180},
											name: "__",
										},
										&choiceExpr{
											pos: position{line: 311, col: 63, offset: 7184},
											alternatives: []interface{}{
												&litMatcher{
													pos:        position{line: 311, col: 63, offset: 7184},
													val:        ",",
													ignoreCase: false,
												},
												&andExpr{
													pos: position{line: 311, col: 69, offset: 7190},
													expr: &litMatcher{
														pos:        position{line: 311, col: 70, offset: 7191},
														val:        "}",
														ignoreCase: false,
													},
												},
											},
										},
										&ruleRefExpr{
											pos:  position{line: 311, col: 75, offset: 7196},
											name: "__",
										},
									},
								},
							},
						},
						&litMatcher{
							pos:        position{line: 311, col: 80, offset: 7201},
							val:        "}",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "Literal",
			pos:  position{line: 327, col: 1, offset: 7447},
			expr: &actionExpr{
				pos: position{line: 327, col: 11, offset: 7459},
				run: (*parser).callonLiteral1,
				expr: &choiceExpr{
					pos: position{line: 327, col: 12, offset: 7460},
					alternatives: []interface{}{
						&seqExpr{
							pos: position{line: 327, col: 13, offset: 7461},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 327, col: 13, offset: 7461},
									val:        "\"",
									ignoreCase: false,
								},
								&zeroOrMoreExpr{
									pos: position{line: 327, col: 17, offset: 7465},
									expr: &choiceExpr{
										pos: position{line: 327, col: 18, offset: 7466},
										alternatives: []interface{}{
											&litMatcher{
												pos:        position{line: 327, col: 18, offset: 7466},
												val:        "\\\"",
												ignoreCase: false,
											},
											&charClassMatcher{
												pos:        position{line: 327, col: 25, offset: 7473},
												val:        "[^\"]",
												chars:      []rune{'"'},
												ignoreCase: false,
												inverted:   true,
											},
										},
									},
								},
								&litMatcher{
									pos:        position{line: 327, col: 32, offset: 7480},
									val:        "\"",
									ignoreCase: false,
								},
							},
						},
						&seqExpr{
							pos: position{line: 327, col: 40, offset: 7488},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 327, col: 40, offset: 7488},
									val:        "'",
									ignoreCase: false,
								},
								&zeroOrMoreExpr{
									pos: position{line: 327, col: 45, offset: 7493},
									expr: &choiceExpr{
										pos: position{line: 327, col: 46, offset: 7494},
										alternatives: []interface{}{
											&litMatcher{
												pos:        position{line: 327, col: 46, offset: 7494},
												val:        "\\'",
												ignoreCase: false,
											},
											&charClassMatcher{
												pos:        position{line: 327, col: 53, offset: 7501},
												val:        "[^']",
												chars:      []rune{'\''},
												ignoreCase: false,
												inverted:   true,
											},
										},
									},
								},
								&litMatcher{
									pos:        position{line: 327, col: 60, offset: 7508},
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
			pos:  position{line: 334, col: 1, offset: 7709},
			expr: &actionExpr{
				pos: position{line: 334, col: 14, offset: 7724},
				run: (*parser).callonIdentifier1,
				expr: &seqExpr{
					pos: position{line: 334, col: 14, offset: 7724},
					exprs: []interface{}{
						&oneOrMoreExpr{
							pos: position{line: 334, col: 14, offset: 7724},
							expr: &choiceExpr{
								pos: position{line: 334, col: 15, offset: 7725},
								alternatives: []interface{}{
									&ruleRefExpr{
										pos:  position{line: 334, col: 15, offset: 7725},
										name: "Letter",
									},
									&litMatcher{
										pos:        position{line: 334, col: 24, offset: 7734},
										val:        "_",
										ignoreCase: false,
									},
								},
							},
						},
						&zeroOrMoreExpr{
							pos: position{line: 334, col: 30, offset: 7740},
							expr: &choiceExpr{
								pos: position{line: 334, col: 31, offset: 7741},
								alternatives: []interface{}{
									&ruleRefExpr{
										pos:  position{line: 334, col: 31, offset: 7741},
										name: "Letter",
									},
									&ruleRefExpr{
										pos:  position{line: 334, col: 40, offset: 7750},
										name: "Digit",
									},
									&charClassMatcher{
										pos:        position{line: 334, col: 48, offset: 7758},
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
			pos:  position{line: 338, col: 1, offset: 7810},
			expr: &charClassMatcher{
				pos:        position{line: 338, col: 17, offset: 7828},
				val:        "[,;]",
				chars:      []rune{',', ';'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "Letter",
			pos:  position{line: 339, col: 1, offset: 7833},
			expr: &charClassMatcher{
				pos:        position{line: 339, col: 10, offset: 7844},
				val:        "[A-Za-z]",
				ranges:     []rune{'A', 'Z', 'a', 'z'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "Digit",
			pos:  position{line: 340, col: 1, offset: 7853},
			expr: &charClassMatcher{
				pos:        position{line: 340, col: 9, offset: 7863},
				val:        "[0-9]",
				ranges:     []rune{'0', '9'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "SourceChar",
			pos:  position{line: 344, col: 1, offset: 7874},
			expr: &anyMatcher{
				line: 344, col: 14, offset: 7889,
			},
		},
		{
			name: "Comment",
			pos:  position{line: 345, col: 1, offset: 7891},
			expr: &choiceExpr{
				pos: position{line: 345, col: 11, offset: 7903},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 345, col: 11, offset: 7903},
						name: "MultiLineComment",
					},
					&ruleRefExpr{
						pos:  position{line: 345, col: 30, offset: 7922},
						name: "SingleLineComment",
					},
				},
			},
		},
		{
			name: "MultiLineComment",
			pos:  position{line: 346, col: 1, offset: 7940},
			expr: &seqExpr{
				pos: position{line: 346, col: 20, offset: 7961},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 346, col: 20, offset: 7961},
						val:        "/*",
						ignoreCase: false,
					},
					&zeroOrMoreExpr{
						pos: position{line: 346, col: 25, offset: 7966},
						expr: &seqExpr{
							pos: position{line: 346, col: 27, offset: 7968},
							exprs: []interface{}{
								&notExpr{
									pos: position{line: 346, col: 27, offset: 7968},
									expr: &litMatcher{
										pos:        position{line: 346, col: 28, offset: 7969},
										val:        "*/",
										ignoreCase: false,
									},
								},
								&ruleRefExpr{
									pos:  position{line: 346, col: 33, offset: 7974},
									name: "SourceChar",
								},
							},
						},
					},
					&litMatcher{
						pos:        position{line: 346, col: 47, offset: 7988},
						val:        "*/",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "MultiLineCommentNoLineTerminator",
			pos:  position{line: 347, col: 1, offset: 7993},
			expr: &seqExpr{
				pos: position{line: 347, col: 36, offset: 8030},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 347, col: 36, offset: 8030},
						val:        "/*",
						ignoreCase: false,
					},
					&zeroOrMoreExpr{
						pos: position{line: 347, col: 41, offset: 8035},
						expr: &seqExpr{
							pos: position{line: 347, col: 43, offset: 8037},
							exprs: []interface{}{
								&notExpr{
									pos: position{line: 347, col: 43, offset: 8037},
									expr: &choiceExpr{
										pos: position{line: 347, col: 46, offset: 8040},
										alternatives: []interface{}{
											&litMatcher{
												pos:        position{line: 347, col: 46, offset: 8040},
												val:        "*/",
												ignoreCase: false,
											},
											&ruleRefExpr{
												pos:  position{line: 347, col: 53, offset: 8047},
												name: "EOL",
											},
										},
									},
								},
								&ruleRefExpr{
									pos:  position{line: 347, col: 59, offset: 8053},
									name: "SourceChar",
								},
							},
						},
					},
					&litMatcher{
						pos:        position{line: 347, col: 73, offset: 8067},
						val:        "*/",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "SingleLineComment",
			pos:  position{line: 348, col: 1, offset: 8072},
			expr: &choiceExpr{
				pos: position{line: 348, col: 21, offset: 8094},
				alternatives: []interface{}{
					&seqExpr{
						pos: position{line: 348, col: 22, offset: 8095},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 348, col: 22, offset: 8095},
								val:        "//",
								ignoreCase: false,
							},
							&zeroOrMoreExpr{
								pos: position{line: 348, col: 27, offset: 8100},
								expr: &seqExpr{
									pos: position{line: 348, col: 29, offset: 8102},
									exprs: []interface{}{
										&notExpr{
											pos: position{line: 348, col: 29, offset: 8102},
											expr: &ruleRefExpr{
												pos:  position{line: 348, col: 30, offset: 8103},
												name: "EOL",
											},
										},
										&ruleRefExpr{
											pos:  position{line: 348, col: 34, offset: 8107},
											name: "SourceChar",
										},
									},
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 348, col: 52, offset: 8125},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 348, col: 52, offset: 8125},
								val:        "#",
								ignoreCase: false,
							},
							&zeroOrMoreExpr{
								pos: position{line: 348, col: 56, offset: 8129},
								expr: &seqExpr{
									pos: position{line: 348, col: 58, offset: 8131},
									exprs: []interface{}{
										&notExpr{
											pos: position{line: 348, col: 58, offset: 8131},
											expr: &ruleRefExpr{
												pos:  position{line: 348, col: 59, offset: 8132},
												name: "EOL",
											},
										},
										&ruleRefExpr{
											pos:  position{line: 348, col: 63, offset: 8136},
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
			pos:  position{line: 350, col: 1, offset: 8152},
			expr: &zeroOrMoreExpr{
				pos: position{line: 350, col: 6, offset: 8159},
				expr: &choiceExpr{
					pos: position{line: 350, col: 8, offset: 8161},
					alternatives: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 350, col: 8, offset: 8161},
							name: "Whitespace",
						},
						&ruleRefExpr{
							pos:  position{line: 350, col: 21, offset: 8174},
							name: "EOL",
						},
						&ruleRefExpr{
							pos:  position{line: 350, col: 27, offset: 8180},
							name: "Comment",
						},
					},
				},
			},
		},
		{
			name: "_",
			pos:  position{line: 351, col: 1, offset: 8191},
			expr: &zeroOrMoreExpr{
				pos: position{line: 351, col: 5, offset: 8197},
				expr: &choiceExpr{
					pos: position{line: 351, col: 7, offset: 8199},
					alternatives: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 351, col: 7, offset: 8199},
							name: "Whitespace",
						},
						&ruleRefExpr{
							pos:  position{line: 351, col: 20, offset: 8212},
							name: "MultiLineCommentNoLineTerminator",
						},
					},
				},
			},
		},
		{
			name: "WS",
			pos:  position{line: 352, col: 1, offset: 8248},
			expr: &zeroOrMoreExpr{
				pos: position{line: 352, col: 6, offset: 8255},
				expr: &ruleRefExpr{
					pos:  position{line: 352, col: 6, offset: 8255},
					name: "Whitespace",
				},
			},
		},
		{
			name: "Whitespace",
			pos:  position{line: 354, col: 1, offset: 8268},
			expr: &charClassMatcher{
				pos:        position{line: 354, col: 14, offset: 8283},
				val:        "[ \\t\\r]",
				chars:      []rune{' ', '\t', '\r'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "EOL",
			pos:  position{line: 355, col: 1, offset: 8291},
			expr: &litMatcher{
				pos:        position{line: 355, col: 7, offset: 8299},
				val:        "\n",
				ignoreCase: false,
			},
		},
		{
			name: "EOS",
			pos:  position{line: 356, col: 1, offset: 8304},
			expr: &choiceExpr{
				pos: position{line: 356, col: 7, offset: 8312},
				alternatives: []interface{}{
					&seqExpr{
						pos: position{line: 356, col: 7, offset: 8312},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 356, col: 7, offset: 8312},
								name: "__",
							},
							&litMatcher{
								pos:        position{line: 356, col: 10, offset: 8315},
								val:        ";",
								ignoreCase: false,
							},
						},
					},
					&seqExpr{
						pos: position{line: 356, col: 16, offset: 8321},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 356, col: 16, offset: 8321},
								name: "_",
							},
							&zeroOrOneExpr{
								pos: position{line: 356, col: 18, offset: 8323},
								expr: &ruleRefExpr{
									pos:  position{line: 356, col: 18, offset: 8323},
									name: "SingleLineComment",
								},
							},
							&ruleRefExpr{
								pos:  position{line: 356, col: 37, offset: 8342},
								name: "EOL",
							},
						},
					},
					&seqExpr{
						pos: position{line: 356, col: 43, offset: 8348},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 356, col: 43, offset: 8348},
								name: "__",
							},
							&ruleRefExpr{
								pos:  position{line: 356, col: 46, offset: 8351},
								name: "EOF",
							},
						},
					},
				},
			},
		},
		{
			name: "EOF",
			pos:  position{line: 358, col: 1, offset: 8356},
			expr: &notExpr{
				pos: position{line: 358, col: 7, offset: 8364},
				expr: &anyMatcher{
					line: 358, col: 8, offset: 8365,
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
			thrift.Unions[v.Name] = (*Struct)(v)
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
