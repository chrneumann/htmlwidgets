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
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"time"
)

const (
	RFC3339      = "2006-01-02T15:04:05"
	RFC3339Nano  = "2006-01-02T15:04:05.999999999"
	RFC3339Short = "2006-01-02T15:04"
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
	// HTML classes to assign.
	Classes []string
	form    *Form
}

// Widget returns the corresponding widget.
func (w WidgetBase) Widget() Widget {
	return w.form.WidgetById(w.Id)
}

func (w WidgetBase) GetRenderData() WidgetRenderData {
	value, err := w.form.getNestedField(w.Id)
	if err != nil {
		panic(fmt.Sprintf("form: Could not find field %q in data: %v", w.Id, err))
	}
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
	w.Errors = nil
	value := values.Get(w.Id)
	w.form.findNestedField(w.Id, value, false)
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

type PasswordWidget struct {
	TextWidget
	// If the user has to repeat the password to verify it, specify at
	// least the label and possible error message to be displayed.
	VerifyLabel, VerifyDescription, VerifyError string
}

func (w PasswordWidget) GetRenderData() WidgetRenderData {
	rd := w.TextWidget.GetRenderData()
	rd.Template = "password"
	return rd
}

type TextAreaWidget struct {
	WidgetBase
	MinLength       int
	ValidationError string
}

func (w *TextAreaWidget) GetRenderData() WidgetRenderData {
	rd := w.Base().GetRenderData()
	rd.Template = "textarea"
	return rd
}

func (w *TextAreaWidget) Fill(values url.Values) bool {
	w.Errors = nil
	value := values.Get(w.Id)
	w.form.findNestedField(w.Id, value, false)
	validated := true
	if len(value) < w.MinLength {
		validated = false
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
			w.form.findNestedField(w.Id, v, false)
		}
	} else {
		w.form.findNestedField(w.Id, false, false)
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
	w.form.findNestedField(w.Id, int(v), false)
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
	if len(values[w.Id]) >= 1 {
		for i, option := range w.Options {
			if option.Value == values.Get(w.Id) {
				value = option.Value
				w.Options[i].Selected = true
			} else {
				w.Options[i].Selected = false
			}
		}
	}
	w.form.findNestedField(w.Id, value, false)
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
	value := values.Get(w.Id)
	w.form.findNestedField(w.Id, value, false)
	return true
}

// FileWidget is a file upload widget that can be used to render a
// HTML file input. It will ignore any uploaded file, you have to
// process it b yourself.
//
// If you add this widget to a Form, the EncTypeAttr ob the RenderData
// will be set on rendering.
type FileWidget struct {
	WidgetBase
}

func (w *FileWidget) GetRenderData() WidgetRenderData {
	rd := w.Base().GetRenderData()
	rd.Template = "file"
	return rd
}

func (w *FileWidget) Fill(values url.Values) bool {
	return true
}

type ListWidget struct {
	WidgetBase
	InnerWidget           Widget
	AddLabel, RemoveLabel string
}

func (w *ListWidget) GetRenderData() WidgetRenderData {
	innerValues, err := w.form.getNestedField(w.Id)
	if err != nil {
		panic(fmt.Sprintf("Could not find field %q in data", w.Id))
	}
	var innerRenderData []WidgetRenderData
	for i := 0; i < innerValues.Len(); i++ {
		*(w.InnerWidget.Base()) = WidgetBase{
			Id:   fmt.Sprintf("%v.%d", w.Id, i),
			form: w.form,
		}
		renderData := w.InnerWidget.GetRenderData()
		innerRenderData = append(innerRenderData,
			renderData)
	}
	return WidgetRenderData{
		WidgetBase: w.WidgetBase,
		Template:   "list",
		Data: map[string]interface{}{
			"Fields":      innerRenderData,
			"AddLabel":    w.AddLabel,
			"RemoveLabel": w.RemoveLabel,
		},
	}
}

func (w *ListWidget) Fill(values url.Values) bool {
	valid := true
	addTo := values.Get("htmlwidgets-action--add-to-list") == w.Id
	var remove []string
	var maxIndex int

	// Find highest index
	re := regexp.MustCompile("^" + w.Id + `\.(\d)$`)
	for key, _ := range values {
		matches := re.FindStringSubmatch(key)
		if len(matches) == 2 {
			if idx, err := strconv.Atoi(matches[1]); err == nil &&
				idx > maxIndex {
				maxIndex = idx
			}
		}
	}
	if addTo {
		maxIndex += 1
	}

	// Fill values into inner fields
	for i := 0; i <= maxIndex; i++ {
		id := fmt.Sprintf("%v.%d", w.Id, i)
		if values.Get("htmlwidgets-action--remove-from-list") == id {
			remove = append(remove, id)
			valid = false
			continue
		}
		if _, ok := values[id]; !ok {
			if !addTo {
				remove = append(remove, id)
			} else {
				addTo = false
				valid = false
			}
		}
		*(w.InnerWidget.Base()) = WidgetBase{
			Id:   id,
			form: w.form,
		}
		if !w.InnerWidget.Fill(values) {
			valid = false
		}
	}

	// Remove fields as requested by the remove action
	for _, id := range remove {
		w.form.findNestedField(id, nil, true)
	}

	// Remove fields after the maximum index
	field, err := w.form.getNestedField(w.Id)
	if err != nil {
		panic(err)
	}
	for i := maxIndex + 1; i < field.Len(); i++ {
		id := fmt.Sprintf("%v.%d", w.Id, i)
		w.form.findNestedField(id, nil, true)
	}
	return valid
}

// TimeWidget is a widget that allows to set a date and time in the
// local timezone.
//
// It tries to parse values as defined in the constants RFC3339,
// RFC3339Nano and RFC3339Short and renders the time as RFC3339Short.
type TimeWidget struct {
	WidgetBase
	Location *time.Location
}

func (w *TimeWidget) GetRenderData() WidgetRenderData {
	if w.Location == nil {
		w.Location = time.UTC
	}
	value, err := w.form.getNestedField(w.Id)
	if err != nil {
		panic(fmt.Sprintf("Could not find field %q in data", w.Id))
	}
	timeValue := value.Interface().(time.Time).In(w.Location).Format(RFC3339Short)
	return WidgetRenderData{
		WidgetBase: w.WidgetBase,
		Template:   "time",
		Data:       timeValue}
}

func (w *TimeWidget) Fill(values url.Values) bool {
	if w.Location == nil {
		w.Location = time.UTC
	}
	value := values.Get(w.Id)
	v, err := time.ParseInLocation(RFC3339Nano, value, w.Location)
	if err != nil {
		v, err = time.ParseInLocation(RFC3339, value, w.Location)
	}
	if err != nil {
		v, err = time.ParseInLocation(RFC3339Short, value, w.Location)
	}
	if err != nil {
		v = time.Time{}
	}
	w.form.findNestedField(w.Id, v, false)
	return true
}
