// Copyright 2018 Tim Möhlmann. All rights reserved.
// This project is licensed under the BSD 3-Clause
// See the LICENSE file for details.

package view

import (
	"bytes"
	"testing"
)

var templateSets = [][]string{
	{
		"a",
		"b",
		"c",
	},
	{
		"top",
		"bottom",
	},
}

//TestSetTemplates tests the setting and parsing of the common templates.
//The last template set will remain the Common.templates for the remainder of the tests.
func TestSetTemplates(test *testing.T) {
	C.Base = "test/"
	for _, ts := range templateSets {
		if err := C.SetTemplates(ts...); err != nil {
			test.Error(
				"For: ", ts, "; ",
				err.Error(),
			)
			continue
		}
		for _, t := range ts {
			if C.templates.Lookup(t) == nil {
				test.Error(
					"For: ", ts,
					"; Expected: ", t,
					"; Got: nil",
				)
				break
			}
		}
	}
}

var testView *View
var testWriter bytes.Buffer

func TestNew(test *testing.T) {
	if C.templates == nil {
		test.Error("Common.templates not set")
		return
	}

	var err error
	testView, err = New(&testWriter, templateSets[0]...)
	if err != nil {
		test.Error("Error in calling New(): ", err.Error())
		return
	}
	for _, ts := range templateSets {
		for _, t := range ts {
			if testView.t.Lookup(t) == nil {
				test.Error(
					"For: ", templateSets,
					"; Expected: ", t,
					"; Got: nil",
				)
				return
			}
		}
	}
}

type testData struct {
	A, B, C, Top, Bottom string
}

var td = testData{"A", "B", "C", "Top", "Bottom"}
var exp = "TopABCBottom"

func TestRender(test *testing.T) {
	if err := testView.Render("a", td); err != nil {
		test.Error("Render error: ", err.Error())
		return
	}
	res := testWriter.String()
	if res != exp {
		test.Error(
			"Render error. ",
			"Expected: ", exp,
			"Got: ", res,
		)
	}
}
