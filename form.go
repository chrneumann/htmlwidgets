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
	"html/template"
	"net/url"
	"reflect"
	"strings"
)

// RenderData contains the data needed for form rendering.
type RenderData struct {
	Widgets []WidgetRenderData
	Errors  []string
	// EncTypeAttr is set to 'enctype="multipart/form-data"' if the Form
	// contains a File widget. Should be used as optional attribute for the form
	// element if the form may contain file input elements.
	EncTypeAttr template.HTMLAttr
	Action      string
}

// Form represents an html form.
type Form struct {
	Widgets []Widget
	data    interface{}
	errors  map[string][]string
	// Action defines the action parameter of the HTML form
	Action string
}

func (f *Form) AddWidget(widget Widget, id, label, description,
	required string) {
	*(widget.Base()) = WidgetBase{id, label, description, required, nil, f}
	f.Widgets = append(f.Widgets, widget)
}

// NewForm creates a new Form with data stored in the
// given pointer to a structure.
//
// In panics if data is not a pointer to a struct.
func NewForm(data interface{}) *Form {
	if dataType := reflect.TypeOf(data); (dataType.Kind() != reflect.Ptr ||
		dataType.Elem().Kind() != reflect.Struct) &&
		dataType.Kind() != reflect.Map {
		panic("NewForm(data, widgets) expects data to" +
			" be a map or a pointer to a struct.")
	}
	form := Form{data: data, Widgets: make([]Widget, 0),
		errors: make(map[string][]string, 0)}
	return &form
}

// RenderData returns a RenderData struct for the form.
//
// It panics if a registered widget is not present in the data struct.
func (f Form) RenderData() (renderData *RenderData) {
	renderData = new(RenderData)
	renderData.Action = f.Action
	renderData.Widgets = make([]WidgetRenderData, 0)
	for _, widget := range f.Widgets {
		/*
			} else if _, ok := widget.(*FileWidget); ok {
				renderData.EncTypeAttr = `enctype="multipart/form-data"`
			}
			value, err := f.getNestedField(widget.Id)
			if err != nil {
				value = reflect.ValueOf("")
			}
		*/
		renderData.Widgets = append(renderData.Widgets,
			widget.GetRenderData())

		/*
			WidgetRenderData{
			Label:       widget.Label,
			Input:       widget.HTML(widget.Id, value.Interface()),
			Description: widget.Description,
			Errors:      f.errors[widget.Id]})
		*/
	}
	renderData.Errors = f.errors[""]
	return
}

/*

// AddError adds an error to a widget's error list.
//
// To add global form errors, use an empty string as the widget's name.
func (f *Form) AddError(widget string, error string) {
	if f.errors[widget] == nil {
		f.errors[widget] = make([]string, 0, 1)
	}
	f.errors[widget] = append(f.errors[widget], error)
}

const (
	rawField = iota
	boolField
	stringField
)

type widgetType struct {
	IsArray   bool
	ValueType int
}

*/

// getNestedField searches for the given nested field in the given data
func (f Form) getNestedField(field string) (reflect.Value, error) {
	return f.findNestedField(field, nil)
}

// findNestedField searches for the given field in the form data.
//
// If setValue is given, it will be set to the field.
func (f *Form) findNestedField(field string, setValue interface{}) (reflect.Value, error) {
	parts := strings.Split(field, ".")
	value := reflect.ValueOf(f.data)
	for len(parts) != 0 {
		setIt := len(parts) == 1 && setValue != nil
		part := parts[0]
		switch value.Type().Kind() {
		case reflect.Ptr, reflect.Interface:
			value = value.Elem()
			continue
		case reflect.Struct:
			value = value.FieldByName(part)
		case reflect.Map:
			if setIt {
				value.SetMapIndex(reflect.ValueOf(part), reflect.ValueOf(setValue))
				return reflect.Value{}, nil
			}
			value = value.MapIndex(reflect.ValueOf(part))
		default:
			return reflect.Value{},
				fmt.Errorf("form: Can't find field %q in data", field)
		}
		if !value.IsValid() {
			return reflect.Value{},
				fmt.Errorf("form: Invalid field %q in data", field)
		}
		parts = parts[1:]
	}
	if setValue != nil {
		if value.Type().Kind() == reflect.Ptr {
			v := reflect.New(value.Type().Elem())
			v.Elem().Set(reflect.ValueOf(setValue))
			value.Set(v)
		} else {
			value.Set(reflect.ValueOf(setValue))
		}
	}
	if value.Type().Kind() == reflect.Interface {
		value = value.Elem()
	}
	return value, nil
}

// Fill fills the form data with the given values and validates the form.
//
// It panics if a widget has been set up which is not present in the
// app data struct.
//
// Values that don't match a widget will be ignored.
//
// Returns true iff the form validates.
func (f *Form) Fill(values url.Values) bool {
	ret := true
	for _, widget := range f.Widgets {
		if ok := widget.Fill(values); !ok {
			ret = false
		}
		/*
			widgetValue, err := f.getNestedField(widget.Id)
			if err != nil {
				continue
			}
			widgetType := widgetValue.Type()
			if fieldType.Kind() == reflect.Slice {
				fieldType = fieldType.Elem()
			}
			if paramValue, ok := values[field.Id]; ok {
				for _, value := range paramValue {
					f.setNestedField(field.Id, value)
				}
			} else {
				f.setNestedField(field.Id, "")
			}
		*/
	}
	//return f.validate()
	return ret
}

/*

// validate validates the currently present data.
//
// Resets any previous errors.
// Returns true iff the data validates.
func (f *Form) validate() bool {
	anyError := false
	for _, field := range f.Fields {
		value, err := f.getNestedField(field.Id)
		if err != nil {
			return false
		}
		if field.Validator != nil {
			if errors := field.Validator(value.Interface()); errors != nil {
				f.errors[field.Id] = errors
				anyError = true
			}
		}
	}
	return !anyError
}

// Validator is a function which validates the given data and returns error
// messages if the data does not validate.
type Validator func(interface{}) []string

// And is a Validator that collects errors of all given validators.
func And(vs ...Validator) Validator {
	return func(value interface{}) []string {
		errors := []string{}
		for _, v := range vs {
			errors = append(errors, v(value)...)
		}
		if len(errors) == 0 {
			return nil
		}
		return errors
	}
}

// Required creates a Validator to check for non empty values.
//
// msg is set as validation error.
func Required(msg string) Validator {
	return func(value interface{}) []string {
		if value == reflect.Zero(reflect.TypeOf(value)).Interface() {
			return []string{msg}
		}
		return nil
	}
}

// Regex creates a Validator to check a string for a matching regexp.
//
// If the expression does not match the string to be validated,
// the given error msg is returned.
func Regex(exp, msg string) Validator {
	return func(value interface{}) []string {
		if matched, _ := regexp.MatchString(exp, value.(string)); !matched {
			return []string{msg}
		}
		return nil
	}
}




















// timeConverter converts a string to a time.Time
func timeConverter(in string) reflect.Value {
	out, err := time.Parse(time.RFC3339, in)
	if err != nil {
		out, err = time.Parse("2006-01-02", in)
	}
	if err != nil {
		out, _ = time.Parse("15:04:05", in)
	}
	return reflect.ValueOf(out)
}

type DateTimeWidget int

func (t DateTimeWidget) HTML(field string, value interface{}) template.HTML {
	var out string
	if obj, ok := value.(time.Time); ok {
		out = obj.Format(time.RFC3339)
	} else if obj, ok := value.(*time.Time); ok {
		if obj == nil {
			out = ""
		} else {
			out = obj.Format(time.RFC3339)
		}
	} else {
		out = fmt.Sprintf("%v", obj)
	}
	return template.HTML(fmt.Sprintf(
		`<input id="%v" type="datetime" name="%v" value="%v"/>`,
		field, field, html.EscapeString(out)))
}

type DateWidget int

func (t DateWidget) HTML(field string, value interface{}) template.HTML {
	var out string
	if obj, ok := value.(time.Time); ok {
		out = obj.Format("2006-01-02")
	} else if obj, ok := value.(*time.Time); ok {
		if obj == nil {
			out = ""
		} else {
			out = obj.Format("2006-01-02")
		}
	} else {
		out = fmt.Sprintf("%v", obj)
	}
	return template.HTML(fmt.Sprintf(
		`<input id="%v" type="date" name="%v" value="%v"/>`,
		field, field, html.EscapeString(out)))
}

type TimeWidget int

func (t TimeWidget) HTML(field string, value interface{}) template.HTML {
	var out string
	if obj, ok := value.(time.Time); ok {
		out = obj.Format("15:04:05")
	} else if obj, ok := value.(*time.Time); ok {
		if obj == nil {
			out = ""
		} else {
			out = obj.Format("15:04:05")
		}
	} else {
		out = fmt.Sprintf("%v", obj)
	}
	return template.HTML(fmt.Sprintf(
		`<input id="%v" type="time" name="%v" value="%v"/>`,
		field, field, html.EscapeString(out)))
}

type Text int

func (t Text) HTML(field string, value interface{}) template.HTML {
	return template.HTML(fmt.Sprintf(
		`<input id="%v" type="text" name="%v" value="%v"/>`,
		field, field, html.EscapeString(
			fmt.Sprintf("%v", value))))
}

type CheckWidget int

func (t CheckWidget) HTML(field string, value interface{}) template.HTML {
	checked := ""
	if val, ok := value.(bool); ok && val {
		checked = ` checked="checked"`
	}
	return template.HTML(fmt.Sprintf(
		`<input id="%v" type="checkbox" name="%v" value="true"%v/>`,
		field, field, checked))
}

type AlohaEditor int

func (t AlohaEditor) HTML(field string, value interface{}) template.HTML {
	return template.HTML(fmt.Sprintf(
		`<textarea class="editor" id="%v" name="%v"/>%v</textarea>`,
		field, field, html.EscapeString(
			fmt.Sprintf("%v", value))))
}

type TextArea int

func (t TextArea) HTML(field string, value interface{}) template.HTML {
	return template.HTML(fmt.Sprintf(
		`<textarea id="%v" name="%v"/>%v</textarea>`,
		field, field, html.EscapeString(
			fmt.Sprintf("%v", value))))
}

// Option of a select widget.
type Option struct {
	Value, Text string
}

// SelectWidget renders a selection field.
type SelectWidget struct {
	Options []Option
}

func (t SelectWidget) HTML(field string, value interface{}) template.HTML {
	var options string
	for _, v := range t.Options {
		selected := ""
		if v.Value == value.(string) {
			selected = " selected"
		}
		options += fmt.Sprintf("<option value=\"%v\"%v>%v</option>\n",
			v.Value, selected, v.Text)
	}
	ret := fmt.Sprintf("<select id=\"%v\" name=\"%v\">\n%v</select>",
		field, field, options)
	return template.HTML(ret)
}

// HiddenWidget renders a hidden input field.
type HiddenWidget int

func (t HiddenWidget) HTML(field string, value interface{}) template.HTML {
	return template.HTML(
		fmt.Sprintf(`<input id="%v" type="hidden" name="%v" value="%v"/>`,
			field, field, value))
}

// PasswordWidget renders a password field.
type PasswordWidget int

func (t PasswordWidget) HTML(field string, value interface{}) template.HTML {
	return template.HTML(
		fmt.Sprintf(`<input id="%v" type="password" name="%v"/>`,
			field, field))
}

// FileWidget renders a file upload field.
type FileWidget int

func (t FileWidget) HTML(field string, value interface{}) template.HTML {
	return template.HTML(
		fmt.Sprintf(`<input id="%v" type="file" name="%v"/>`,
			field, field))
}


*/
