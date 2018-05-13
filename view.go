// Copyright 2018 Tim Möhlmann. All rights reserved.
// This project is licensed under the BSD 3-Clause
// See the LICENSE file for details.

// Package view manages template loading and caching. Data inclusion to the template and execution.
// It features common settings of template name prefix and suffix (defaults to templates/ and .html).
// It holds a set of Common templates, which are automatically included in every new View.
package view

import (
	"html/template"
	"io"
)

//parse adds one or more templates to a Template object pointer. It takes one or more template names in its arguments.
//It returns error if len(tn) == 0 or in case Template.ParseFiles fails.
func parse(to *template.Template, tn ...string) (err error) {
	for _, v := range tn {
		if to, err = to.ParseFiles(C.Base + v + C.Ext); err != nil {
			return
		}
	}
	return
}

//Common keeps the common templates and common tempate data
type Common struct {
	Base      string //Basename, ussualy the directory where the templates reside
	Ext       string //Extention for the template file
	templates *template.Template
}

//C declares Common with default values
var C = Common{
	Base: "templates/", //Default directory
	Ext:  ".html",      //Default extention
}

//SetTemplates sets the common templates. The common templates will be cached and available in every new view.
//This method creates a new template object and all previous addes templates are lost.
//It returns an error if one of the templates fail to parse. In this case, if there was a previous template set they will remain un-affected.
func (c *Common) SetTemplates(tn ...string) (err error) {
	t := template.New("")
	if err = parse(t, tn...); err != nil {
		return
	}
	c.templates = t
	return
}

//View holds the templates and the output writer.
type View struct {
	t *template.Template
	w io.Writer
}

//New creates a new View and assoctiates it with the output Writer from the first argument.
//It needs 0 or more template names that need to be loaded for this view. Templates may be loaded from cache.
//This template set will be merged with the common template set.
func New(w io.Writer, tn ...string) (v *View, err error) {
	t, _ := C.templates.Clone() //Clone never returns error
	if len(tn) > 0 {
		if err = parse(t, tn...); err != nil {
			return
		}
	}
	v = &View{
		t: t,
		w: w,
	}
	return
}

//Render the view. tmpl is the name of the main template being rendered. Data will be passed directly to the template.
func (v *View) Render(tmpl string, data interface{}) (err error) {
	if err = v.t.ExecuteTemplate(v.w, tmpl, data); err != nil {
		return
	}
	return
}
