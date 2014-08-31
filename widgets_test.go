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
	// RenderData is the expected value of the WidgetRenderData Data field
	RenderData interface{}
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
		t.Errorf("AppStruct field is\n%v\nshould be \n%v", test.AppStruct,
			test.FilledValue)
	}
	renderData := form.RenderData()
	expected := WidgetRenderData{
		WidgetBase: WidgetBase{
			Id:          "Id",
			Label:       "Label",
			Description: "Description",
		},
		Data: test.RenderData,
	}
	if len(renderData.Errors) > 0 {
		t.Errorf("RenderData contains general errors: %v", renderData.Errors)
	}
	renderData.Widgets[0].WidgetBase.form = nil
	if !reflect.DeepEqual(renderData.Widgets[0], expected) {
		t.Errorf("RenderData for Widget '%v' =\n%#v,\nexpected\n%#v",
			expected.Id, renderData.Widgets[0], expected)
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
		RenderData: []SelectOption{
			SelectOption{"foo", "Foo", false},
			SelectOption{"bar", "Bar", true},
		},
	})
}

/*

func TestSelectWidget(t *testing.T) {
	widget := SelectWidget{[]Option{
		Option{"foo", "The Foo!"},
		Option{"bar", "The Bar!"}}}
	tests := []struct {
		Name, Value, Expected string
	}{
		{"TestSelect", "", `<select id="TestSelect" name="TestSelect">
<option value="foo">The Foo!</option>
<option value="bar">The Bar!</option>
</select>`},
		{"TestSelect2", "unknown!", `<select id="TestSelect2" name="TestSelect2">
<option value="foo">The Foo!</option>
<option value="bar">The Bar!</option>
</select>`},
		{"TestSelect3", "foo", `<select id="TestSelect3" name="TestSelect3">
<option value="foo" selected>The Foo!</option>
<option value="bar">The Bar!</option>
</select>`},
		{"TestSelect4", "bar", `<select id="TestSelect4" name="TestSelect4">
<option value="foo">The Foo!</option>
<option value="bar" selected>The Bar!</option>
</select>`}}
	for _, v := range tests {
		ret := widget.HTML(v.Name, v.Value)
		if string(ret) != v.Expected {
			t.Errorf(`SelectWidget.HTML("%v", "%v") = "%v", should be "%v".`,
				v.Name, v.Value, ret, v.Expected)
		}
	}
}

func TestHiddenWidget(t *testing.T) {
	widget := new(HiddenWidget)
	ret := widget.HTML("foo", "bar")
	expected := `<input id="foo" type="hidden" name="foo" value="bar"/>`
	if string(ret) != expected {
		t.Errorf(`HiddenWidget.HTML("Foo", "bar") = "%v", should be "%v".`,
			ret, expected)
	}
}

func TestFileWidget(t *testing.T) {
	widget := new(FileWidget)
	ret := widget.HTML("foo", "")
	expected := `<input id="foo" type="file" name="foo"/>`
	if string(ret) != expected {
		t.Errorf(`FileWidget.HTML("Foo", "") = "%v", should be "%v".`,
			ret, expected)
	}
}

func TestPasswordWidget(t *testing.T) {
	widget := new(PasswordWidget)
	ret := widget.HTML("foo", "")
	expected := `<input id="foo" type="password" name="foo"/>`
	if string(ret) != expected {
		t.Errorf(`PasswordWidget.HTML("Foo", "") = "%v", should be "%v".`,
			ret, expected)
	}
}

type TestDateTimeWidgetData struct {
	ID *time.Time
}

func TestDateTimeWidget(t *testing.T) {
	data := TestDateTimeWidgetData{}
	input := `<input id="ID" type="datetime" name="ID" value="2008-09-08T22:47:31-07:00"/>`
	zeroInput := `<input id="ID" type="datetime" name="ID" value=""/>`
	value, err := time.Parse(time.RFC3339, "2008-09-08T22:47:31-07:00")
	if err != nil {
		t.Fatal(err)
	}
	testWidget(t, new(DateTimeWidget), &data, input, zeroInput, &value,
		"2008-09-08T22:47:31-07:00")
}

type TestCheckWidgetData struct {
	ID bool
}

func TestCheckWidget(t *testing.T) {
	data := TestCheckWidgetData{}
	input := `<input id="ID" type="checkbox" name="ID" value="true" checked="checked"/>`
	zeroInput := `<input id="ID" type="checkbox" name="ID" value="true"/>`
	testWidget(t, new(CheckWidget), &data, input, zeroInput, true, "true")
}

/*
func TestDateWidget(t *testing.T) {
	data := TestDateTimeWidgetData{}
	input := `<input id="ID" type="date" name="ID" value="2008-09-08"/>`
	zeroInput := `<input id="ID" type="date" name="ID" value=""/>`
	value, err := time.Parse("2006-01-02", "2008-09-08")
	if err != nil {
		t.Fatal(err)
	}
	testWidget(t, new(DateWidget), &data, input, zeroInput, value, "2008-09-08")
}

func TestTimeWidget(t *testing.T) {
	data := TestDateTimeWidgetData{}
	input := `<input id="ID" type="time" name="ID" value="22:47:31"/>`
	zeroInput := `<input id="ID" type="time" name="ID" value=""/>`
	value, err := time.Parse("15:04:05", "22:47:31")
	if err != nil {
		t.Fatal(err)
	}
	testWidget(t, new(TimeWidget), &data, input, zeroInput, value, "22:47:31")
}
*/
