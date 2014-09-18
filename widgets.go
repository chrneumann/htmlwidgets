// This file is part of htmlwidgets.
// Copyright 2014 Christian Neumann <cneumann@datenkarussell.de>

// htmlwidgets is free software: you can redistribute it and/or modify it under
// the terms of the GNU Lesser General Public License as published by the Free
// Software Foundation, either version 3 of the License, or (at your option) any
// later version.

// htmlwidgets is distributed in the hope that it will be useful, but WITHOUT
// ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS
// FOR A PARTICULAR PURPOSE. See the GNU Lesser General Public License for more
// details.

// You should have received a copy of the GNU Lesser General Public License
// along with htmlwidgets. If not, see <http://www.gnu.org/licenses/>.

package htmlwidgets

import (
	"net/url"
	"regexp"
	"strconv"
)

// WidgetRenderData contains the data needed for widget rendering.
type WidgetRenderData struct {
	WidgetBase
	// Template is the id of the template to be used to render the widget.
	Template string
	// Data contains any widget dependent data used to render the widget.
	Data interface{}
}

type Widget interface {
	// GetRenderData returns the data needed to render the widget.
	GetRenderData() WidgetRenderData
	// Fill reads the given values to fill into the app struct.
	Fill(url.Values) bool
	Base() *WidgetBase
}

// WidgetBase contains common fields used by widgets.
type WidgetBase struct {
	Id, Label, Description string
	// Errors contains any validation errors.
	Errors []string
	form   *Form
}

func (w WidgetBase) GetRenderData() WidgetRenderData {
	value, _ := w.form.getNestedField(w.Id)
	return WidgetRenderData{
		WidgetBase: w,
		Template:   "text",
		Data:       value.Interface()}
}

func (w *WidgetBase) Base() *WidgetBase {
	return w
}

type TextWidget struct {
	WidgetBase
	MinLength       int
	Regexp          string
	ValidationError string
}

func (w *TextWidget) GetRenderData() WidgetRenderData {
	rd := w.Base().GetRenderData()
	rd.Template = "text"
	return rd
}

func (w *TextWidget) Fill(values url.Values) bool {
	value := ""
	if len(values[w.Id]) == 1 {
		value = values[w.Id][0]
	}
	w.form.findNestedField(w.Id, value)

	validated := true
	if len(value) < w.MinLength {
		validated = false
	}
	if validated && len(w.Regexp) > 0 {
		if matched, _ := regexp.MatchString(w.Regexp, value); !matched {
			validated = false
		}
	}
	if !validated {
		w.Errors = append(w.Errors, w.ValidationError)
		return false
	}
	return true
}

type BoolWidget struct{ WidgetBase }

func (w *BoolWidget) GetRenderData() WidgetRenderData {
	rd := w.Base().GetRenderData()
	rd.Template = "checkbox"
	return rd
}

func (w *BoolWidget) Fill(values url.Values) bool {
	if len(values[w.Id]) != 0 {
		if v, err := strconv.ParseBool(values[w.Id][0]); err == nil {
			w.form.findNestedField(w.Id, v)
		}
	}
	return true
}

type IntegerWidget struct{ WidgetBase }

func (w *IntegerWidget) GetRenderData() WidgetRenderData {
	rd := w.Base().GetRenderData()
	rd.Template = "text"
	return rd
}

func (w *IntegerWidget) Fill(values url.Values) bool {
	v, err := strconv.ParseInt(values[w.Id][0], 0, 0)
	if err != nil {
		return false
	}
	w.form.findNestedField(w.Id, int(v))
	return true
}

// SelectOption is an option to choose from in a SelectWidget
type SelectOption struct {
	Value, Description string
	Selected           bool
}

// SelectWidget allows to choose one from multiple options.
type SelectWidget struct {
	WidgetBase
	Options []SelectOption
}

func (w *SelectWidget) Fill(values url.Values) bool {
	value := w.Options[0].Value
	if len(values[w.Id]) == 1 {
		for i, option := range w.Options {
			if option.Value == values[w.Id][0] {
				value = option.Value
				w.Options[i].Selected = true
			} else {
				w.Options[i].Selected = false
			}
		}
	}
	w.form.findNestedField(w.Id, value)
	return true
}

func (w SelectWidget) GetRenderData() WidgetRenderData {
	return WidgetRenderData{
		WidgetBase: w.WidgetBase,
		Template:   "select",
		Data:       w.Options}
}

type HiddenWidget struct {
	WidgetBase
}

func (w *HiddenWidget) GetRenderData() WidgetRenderData {
	rd := w.Base().GetRenderData()
	rd.Template = "hidden"
	return rd
}

func (w *HiddenWidget) Fill(values url.Values) bool {
	value := ""
	if len(values[w.Id]) == 1 {
		value = values[w.Id][0]
	}
	w.form.findNestedField(w.Id, value)
	return true
}
