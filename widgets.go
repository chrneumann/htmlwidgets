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
	"strconv"
)

// WidgetRenderData contains the data needed for widget rendering.
type WidgetRenderData struct {
	WidgetBase
	// Data contains any widget dependent data used to render the widget.
	Data interface{}
}

type Widget interface {
	// GetRenderData returns the data needed to render the widget
	GetRenderData() WidgetRenderData
	// Fill reads the given values to fill into the app struct.
	Fill(url.Values) bool
	Base() *WidgetBase
}

// WidgetBase contains common fields used by widgets.
type WidgetBase struct {
	Id, Label, Description, Required string
	// Errors contains any validation errors.
	Errors []string
	Form   *Form
}

func (w WidgetBase) GetRenderData() WidgetRenderData {
	value, _ := w.Form.getNestedField(w.Id)
	return WidgetRenderData{WidgetBase: w, Data: value.Interface()}
}

func (w *WidgetBase) Base() *WidgetBase {
	return w
}

type TextWidget struct{ WidgetBase }

func (w *TextWidget) Fill(values url.Values) bool {
	value := ""
	if len(values[w.Id]) == 1 {
		value = values[w.Id][0]
	}
	w.Form.findNestedField(w.Id, value)
	if len(w.Required) > 0 && len(value) == 0 {
		w.Errors = append(w.Errors, w.Required)
		return false
	}
	return true
}

type BoolWidget struct{ WidgetBase }

func (w *BoolWidget) Fill(values url.Values) bool {
	v, err := strconv.ParseBool(values[w.Id][0])
	if err != nil {
		return false
	}
	w.Form.findNestedField(w.Id, v)
	return true
}

type IntegerWidget struct{ WidgetBase }

func (w *IntegerWidget) Fill(values url.Values) bool {
	v, err := strconv.ParseInt(values[w.Id][0], 0, 0)
	if err != nil {
		return false
	}
	w.Form.findNestedField(w.Id, int(v))
	return true
}
