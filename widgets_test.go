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
	"reflect"
	"testing"
	"time"
)

// WidgetTest defines a test for testWidget
type WidgetTest struct {
	// Widget is the widget to test
	Widget Widget
	// AppStruct is the appstruct which will be filled
	AppStruct interface{}
	// UrlValue is the submitted value
	URLValue string
	// FilledValue is the expected value filled into the appstruct
	FilledValue interface{}
	// EmptyValue is the expected value filled into the appstruct if no
	// or an invalid value is submitted.
	EmptyValue interface{}
	// RenderData is the expected value of the WidgetRenderData Data field
	RenderData interface{}
	// Error is the expected error if any
	Error string
	// Template is the expected template Id
	Template string
}

// testWidget performs common tests on the given widget
func testWidget(t *testing.T, test *WidgetTest) {
	form := NewForm(test.AppStruct)
	form.AddWidget(test.Widget, "Id", "Label", "Description")
	urlValues := url.Values{
		"Id": []string{test.URLValue},
	}
	form.Fill(urlValues)
	if !reflect.DeepEqual(test.FilledValue,
		reflect.ValueOf(test.AppStruct).Elem().FieldByName("Id").Interface()) {
		t.Errorf("AppStruct for filled field is\n%v\nshould be \n%v", test.AppStruct,
			test.FilledValue)
	}
	renderData := form.RenderData()
	var errors []string
	if len(test.Error) > 0 {
		errors = append(errors, test.Error)
	}
	expected := WidgetRenderData{
		WidgetBase: WidgetBase{
			Id:          "Id",
			Label:       "Label",
			Description: "Description",
			Errors:      errors,
		},
		Data:     test.RenderData,
		Template: test.Template,
	}
	if len(renderData.Errors) > 0 {
		t.Errorf("RenderData contains general errors: %v", renderData.Errors)
	}
	renderData.Widgets[0].WidgetBase.form = nil
	if !reflect.DeepEqual(renderData.Widgets[0], expected) {
		t.Errorf("RenderData for Widget '%T' =\n%#v,\nexpected\n%#v",
			test.Widget, renderData.Widgets[0], expected)
	}
	form.Fill(nil)
	if !reflect.DeepEqual(test.EmptyValue,
		reflect.ValueOf(test.AppStruct).Elem().FieldByName("Id").Interface()) {
		t.Errorf("AppStruct for missing value field is\n%v\nshould be \n%v",
			test.AppStruct, test.EmptyValue)
	}
}

type TestSelectWidgetData struct {
	Id string
}

func TestSelectWidget(t *testing.T) {
	testWidget(t, &WidgetTest{
		Widget: &SelectWidget{Options: []SelectOption{
			SelectOption{"foo", "Foo", true},
			SelectOption{"bar", "Bar", false},
		}},
		AppStruct:   &TestSelectWidgetData{},
		URLValue:    "bar",
		FilledValue: "bar",
		EmptyValue:  "foo",
		RenderData: []SelectOption{
			SelectOption{"foo", "Foo", false},
			SelectOption{"bar", "Bar", true},
		},
		Template: "select",
	})
}

type TestHiddenWidgetData struct {
	Id string
}

func TestHiddenWidget(t *testing.T) {
	testWidget(t, &WidgetTest{
		Widget:      new(HiddenWidget),
		AppStruct:   &TestHiddenWidgetData{},
		URLValue:    "foo",
		FilledValue: "foo",
		EmptyValue:  "",
		RenderData:  "foo",
		Template:    "hidden",
	})
}

type TestFileWidgetData struct {
	Id string
}

func TestFileWidget(t *testing.T) {
	testWidget(t, &WidgetTest{
		Widget:      new(FileWidget),
		AppStruct:   &TestFileWidgetData{},
		URLValue:    "",
		FilledValue: "",
		EmptyValue:  "",
		RenderData:  "",
		Template:    "file",
	})
}

type TestBoolWidgetData struct {
	Id bool
}

func TestBoolWidget(t *testing.T) {
	testWidget(t, &WidgetTest{
		Widget:      new(BoolWidget),
		AppStruct:   &TestBoolWidgetData{},
		URLValue:    "true",
		FilledValue: true,
		EmptyValue:  false,
		RenderData:  true,
		Template:    "checkbox",
	})
}

type TestTextWidgetData struct {
	Id string
}

func TestTextWidget(t *testing.T) {
	testWidget(t, &WidgetTest{
		Widget:      new(TextWidget),
		AppStruct:   &TestTextWidgetData{},
		URLValue:    "foo",
		FilledValue: "foo",
		EmptyValue:  "",
		RenderData:  "foo",
		Template:    "text",
	})
	testWidget(t, &WidgetTest{
		Widget:      &TextWidget{MinLength: 5, ValidationError: ">=5"},
		AppStruct:   &TestTextWidgetData{},
		URLValue:    "foo",
		FilledValue: "foo",
		EmptyValue:  "",
		RenderData:  "foo",
		Error:       ">=5",
		Template:    "text",
	})
	testWidget(t, &WidgetTest{
		Widget:      &TextWidget{Regexp: `^\w{2}$`, ValidationError: "exactly 2"},
		AppStruct:   &TestTextWidgetData{},
		URLValue:    "fo",
		FilledValue: "fo",
		EmptyValue:  "",
		RenderData:  "fo",
		Template:    "text",
	})
	testWidget(t, &WidgetTest{
		Widget:      &TextWidget{Regexp: `^\w{2}$`, ValidationError: "exactly 2"},
		AppStruct:   &TestTextWidgetData{},
		URLValue:    "foo",
		FilledValue: "foo",
		EmptyValue:  "",
		RenderData:  "foo",
		Error:       "exactly 2",
		Template:    "text",
	})
}

type TestTextAreaWidgetData struct {
	Id string
}

func TestTextAreaWidget(t *testing.T) {
	testWidget(t, &WidgetTest{
		Widget:      new(TextWidget),
		AppStruct:   &TestTextWidgetData{},
		URLValue:    "foo",
		FilledValue: "foo",
		EmptyValue:  "",
		RenderData:  "foo",
		Template:    "text",
	})
	testWidget(t, &WidgetTest{
		Widget:      &TextWidget{MinLength: 5, ValidationError: ">=5"},
		AppStruct:   &TestTextWidgetData{},
		URLValue:    "foo",
		FilledValue: "foo",
		EmptyValue:  "",
		RenderData:  "foo",
		Error:       ">=5",
		Template:    "text",
	})
}

type TestTimeWidgetData struct {
	Id time.Time
}

func TestTimeWidget(t *testing.T) {
	testWidget(t, &WidgetTest{
		Widget:      new(TimeWidget),
		AppStruct:   &TestTimeWidgetData{},
		URLValue:    "1985-04-10T08:10",
		FilledValue: time.Date(1985, time.April, 10, 8, 10, 0, 0, time.UTC),
		EmptyValue:  time.Time{},
		RenderData:  "1985-04-10T08:10:00",
		Template:    "time",
	})
}
