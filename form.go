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
	"strconv"
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
	Widgets   []Widget
	widgetMap map[string]Widget
	data      interface{}
	errors    map[string][]string
	// Action defines the action parameter of the HTML form
	Action string
}

// WidgetById returns the widget with the given id.
func (f Form) WidgetById(id string) Widget {
	return f.widgetMap[id]
}

// AddWidget adds a new widget to the form and sets the given attributes.
//
// It returns the added widget
func (f *Form) AddWidget(widget Widget, id, label, description string) Widget {
	base := widget.Base()
	if base == nil {
		*base = WidgetBase{}
	}
	base.Id = id
	base.Label = label
	base.Description = description
	base.form = f
	f.Widgets = append(f.Widgets, widget)
	f.widgetMap[id] = widget
	return widget
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
	form := Form{
		data:      data,
		Widgets:   make([]Widget, 0),
		widgetMap: make(map[string]Widget),
		errors:    make(map[string][]string, 0)}
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
		if _, ok := widget.(*FileWidget); ok {
			renderData.EncTypeAttr = `enctype="multipart/form-data"`
		}
		widgetRenderData := widget.GetRenderData()
		widgetRenderData.Errors = append(widgetRenderData.Errors,
			f.errors[widget.Base().Id]...)
		renderData.Widgets = append(renderData.Widgets, widgetRenderData)
	}
	renderData.Errors = f.errors[""]
	return
}

// AddError adds an error to a widget's error list.
//
// To add global form errors, use an empty string as the widget's name.
func (f *Form) AddError(widgetId string, error string) {
	f.errors[widgetId] = append(f.errors[widgetId], error)
}

// getNestedField searches for the given nested field in the given data
func (f Form) getNestedField(field string) (reflect.Value, error) {
	return f.findNestedField(field, nil, false)
}

// findNestedField searches for the given field in the form data.
//
// If setValue is given, it will be set to the field.
// If remove is given, the value will be removed from its parent slice
// or map.
func (f *Form) findNestedField(field string, setValue interface{}, remove bool) (
	reflect.Value, error) {
	parts := strings.Split(field, ".")
	value := reflect.ValueOf(f.data)
	var lastMapValue reflect.Value
	var lastMapIndex reflect.Value
	for len(parts) != 0 {
		setIt := len(parts) == 1 && setValue != nil
		removeIt := len(parts) == 1 && remove
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
			if removeIt {
				value.SetMapIndex(reflect.ValueOf(part), reflect.Value{})
				return reflect.Value{}, nil
			}
			lastMapValue = value
			lastMapIndex = reflect.ValueOf(part)
			value = value.MapIndex(reflect.ValueOf(part))
		case reflect.Slice:
			index, err := strconv.Atoi(part)
			if err != nil {
				return reflect.Value{},
					fmt.Errorf("Form: Expected index, got %q in field id %q", part, index)
			}
			if removeIt {
				sliceSetvalue := reflect.AppendSlice(
					value.Slice(0, index),
					value.Slice(index+1, value.Len()))
				if value.CanSet() {
					value.Set(sliceSetvalue)
				} else {
					lastMapValue.SetMapIndex(lastMapIndex, sliceSetvalue)
					value = lastMapValue.MapIndex(lastMapIndex).Elem()
				}
				return reflect.Value{}, nil
			}
			if value.Len() == index {
				sliceSetvalue := reflect.Append(value, reflect.New(value.Type().Elem()).Elem())
				if value.CanSet() {
					value.Set(sliceSetvalue)
				} else {
					lastMapValue.SetMapIndex(lastMapIndex, sliceSetvalue)
					value = lastMapValue.MapIndex(lastMapIndex).Elem()
				}
			}
			value = value.Index(index)
		default:
			return reflect.Value{},
				fmt.Errorf("form: Can't find field %q in data, reflect type: %v",
					field, value.Type().Kind())
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
// Returns true iff the form validates and there are none of the known
// "htmlwidgets-action--*" parameters present.
func (f *Form) Fill(values url.Values) bool {
	ret := true
	for _, widget := range f.Widgets {
		if ok := widget.Fill(values); !ok {
			ret = false
		}
	}
	return ret
}
