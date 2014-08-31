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
	"log"
	"net/url"
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
	Id, Label, Description string
	// Errors contains any validation errors.
	Errors []string
	Form   *Form
}

func (w WidgetBase) GetRenderData() WidgetRenderData {
	value, _ := w.Form.getNestedField(w.Id)
	return WidgetRenderData{WidgetBase: w, Data: value.Interface()}
}

func (w *WidgetBase) Fill(values url.Values) bool {
	log.Println(values[w.Id])
	w.Form.setNestedField(w.Id, values[w.Id][0])
	return true
}

func (w *WidgetBase) Base() *WidgetBase {
	return w
}

type TextWidget struct{ WidgetBase }
type IntegerWidget struct{ WidgetBase }
type CheckboxWidget struct{ WidgetBase }
