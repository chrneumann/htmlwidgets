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

type TestAppDataEmbed struct {
	Title string
}

type TestAppData struct {
	TestAppDataEmbed
	//Name  string
	Age   int
	Extra map[string]interface{}
}

func TestRender(t *testing.T) {
	data := TestAppData{}
	data.Extra = make(map[string]interface{})
	data.Extra["ExtraField"] = ""
	form := NewForm(&data)
	form.AddWidget(new(TextWidget), "Title", "Title", "Your title")
	//TextWidget{Id: "Name", Label: "Name", Description: "Your full name"},
	form.AddWidget(new(IntegerWidget), "Age", "Age", "Years since your birth")
	form.AddWidget(new(BoolWidget), "Extra.ExtraField", "Alive", "Still alive?")

	urlValues := url.Values{
		"Title": []string{""},
		//	"Name":             []string{""},
		"Age":              []string{"14"},
		"Extra.ExtraField": []string{"true"},
	}
	form.Fill(urlValues)
	form.Action = "targetURL"
	renderData := form.RenderData()
	if renderData.Action != "targetURL" {
		t.Errorf(`renderData.Action = %q, should be "targetURL"`, renderData.Action)
	}
	fieldTests := []WidgetRenderData{
		WidgetRenderData{
			WidgetBase: WidgetBase{
				Id:          "Title",
				Label:       "Title",
				Description: "Your title",
			},
			Data: "",
		},
		/*
			{
				Field: "Name",
				Expected: WidgetRenderData{
					Id: "Name",
					Label:       "Your name",
					Description: "Your full name",
					Errors:      []string{"Required!"},
					Data:       "",
		*/
		WidgetRenderData{
			WidgetBase: WidgetBase{
				Id:          "Age",
				Label:       "Age",
				Description: "Years since your birth",
			},
			Data: 14,
		},
		WidgetRenderData{
			WidgetBase: WidgetBase{
				Id:          "Extra.ExtraField",
				Label:       "Alive",
				Description: "Still alive?",
			},
			Data: true,
		},
	}
	for i, test := range fieldTests {
		if len(renderData.Errors) > 0 {
			t.Errorf("RenderData contains general errors: %v", renderData.Errors)
		}
		renderData.Widgets[i].WidgetBase.Form = nil
		if !reflect.DeepEqual(renderData.Widgets[i], test) {
			t.Errorf("RenderData for Widget '%v' =\n%#v,\nexpected\n%#v",
				test.Id, renderData.Widgets[i], test)
		}
	}
}

/*

func TestMapRender(t *testing.T) {
	data := make(map[string]interface{})
	data["Name"] = new(string)
	data["Age"] = new(int)
	data["Foo"] = map[string]string{
		"Bar": "ee"}

	form := NewForm(data, []Field{
		Field{"Name", "Your name", "Your full name", Required("Req!"), nil},
		Field{"Age", "Your age", "Years since your birth.", Required("Req!"), nil},
		Field{"Foo.Bar", "Bar", "Some foo's bar.", Required("Req!"), nil},
	})
	vals := url.Values{
		"Name":    []string{""},
		"Age":     []string{"14"},
		"Foo.Bar": []string{"Bla"},
	}
	form.Fill(vals)
	renderData := form.RenderData()
	fieldTests := []struct {
		Field    string
		Expected WidgetRenderData
	}{
		{
			Field: "Name",
			Expected: WidgetRenderData{
				Label:       "Your name",
				Description: "Your full name",
				Errors:      []string{"Req!"},
				Input:       `<input id="Name" type="text" name="Name" value=""/>`}},
		{
			Field: "AGE",
			Expected: WidgetRenderData{
				Label:       "Your age",
				Description: "Years since your birth.",
				Errors:      nil,
				Input:       `<input id="Age" type="text" name="Age" value="14"/>`}},
		{
			Field: "Foo.Bar",
			Expected: WidgetRenderData{
				Label:       "Bar",
				Description: "Some foo's bar.",
				Errors:      nil,
				Input:       `<input id="Foo.Bar" type="text" name="Foo.Bar" value="Bla"/>`}},
	}
	for i, test := range fieldTests {
		if len(renderData.Errors) > 0 {
			t.Errorf("RenderData contains general errors: %v", renderData.Errors)
		}
		if !reflect.DeepEqual(renderData.Fields[i], test.Expected) {
			t.Errorf("RenderData for Field '%v' =\n%v,\nexpected\n%v",
				test.Field, renderData.Fields[i], test.Expected)
		}
	}
}

func TestAddError(t *testing.T) {
	data := TestData{}
	form := NewForm(&data, []Field{
		Field{"Name", "Your name", "Your full name", Required("Req!"), nil},
		Field{"Age", "Your age", "Years since your birth.", Required("Req!"), nil}})
	form.AddError("Name", "Foo")
	form.AddError("", "Bar")
	renderData := form.RenderData()
	if len(renderData.Fields[0].Errors) != 1 ||
		renderData.Fields[0].Errors[0] != "Foo" {
		t.Errorf(`Field "Name" should have error "Foo"`)
	}
	if len(renderData.Errors) != 1 ||
		renderData.Errors[0] != "Bar" {
		t.Errorf(`Missing global error "Bar"`)
	}
	if len(renderData.Fields[1].Errors) != 0 {
		t.Errorf(`Field "Bar" should have no errors`)
	}
}

type TestDataEncTypeAttr struct {
	Name string
	File string
}

func TestEncTypeAttr(t *testing.T) {
	data := TestDataEncTypeAttr{}
	vals := url.Values{
		"Name": []string{""}}
	fieldTests := []struct {
		Form    *Form
		EncType string
	}{
		{
			Form: NewForm(&data, []Field{
				Field{"Name", "Your name", "Your full name", Required("Req!"),
					nil},
				Field{"File", "File Dummy", "", nil, nil}}),
			EncType: ""},
		{
			Form: NewForm(&data, []Field{
				Field{"Name", "Your name", "Your full name", Required("Req!"), nil},
				Field{"File", "File!", "", nil, new(FileWidget)}}),
			EncType: `enctype="multipart/form-data"`}}

	for i, v := range fieldTests {
		v.Form.Fill(vals)
		renderData := v.Form.RenderData()
		if string(renderData.EncTypeAttr) != v.EncType {
			t.Errorf("Test %v: RenderData.EncTypeAttr is %q, should be %q", i,
				renderData.EncTypeAttr, v.EncType)
		}
	}
}

func TestFill(t *testing.T) {
	data := TestData{}
	data.Extra = make(map[string]interface{}, 0)
	data.Extra["Number"] = new(int)
	form := NewForm(&data, []Field{
		Field{"Name", "Your name", "Your full name", Required("Req!"), nil},
		Field{"Age", "Your age", "Years since your birth.", Required("Req!"), nil},
		Field{"Extra.Number", "Number", "", nil, nil},
	})
	vals := url.Values{
		"Name":         []string{"Foo"},
		"Age":          []string{"14"},
		"Foo":          []string{"noting here"},
		"Extra.Number": []string{"10"},
	}
	expected := TestData{Name: "Foo", Age: 14}
	expected.Extra = make(map[string]interface{}, 0)
	number := 10
	expected.Extra["Number"] = number
	if !form.Fill(vals) {
		t.Errorf("form.Fill(..) returns false, should be true. Errors: %v",
			form.RenderData().Errors)
	}
	if !reflect.DeepEqual(expected, data) {
		t.Errorf("Filled data should be %v, is %v", expected, data)
	}
	vals["Name"] = []string{""}
	data.Name = ""
	if form.Fill(vals) {
		t.Errorf("form.Fill(..) returns true, should be false.")
	}
}

func TestRequire(t *testing.T) {
	invalid, valid := "", "foo"
	validator := Required("Req!")
	err := validator(valid)
	if err != nil {
		t.Errorf("require(%v) = %v, want %v", valid, err, nil)
	}
	err = validator(invalid)
	if err == nil {
		t.Errorf("require(%v) = %v, want %v", invalid, err, "'Required.'")
	}
}

func TestRegex(t *testing.T) {
	tests := []struct {
		Exp    string
		String string
		Valid  bool
	}{
		{"^[\\w]+$", "", false},
		{"^[\\w]+$", "foobar", true},
		{"", "", true},
		{"", "foobar", true},
		{"^[^!]+$", "foobar", true},
		{"^[^!]+$", "foo!bar", false}}

	for _, v := range tests {
		ret := Regex(v.Exp, "damn!")(v.String)
		if (ret == nil && !v.Valid) || (ret != nil && v.Valid) {
			t.Errorf(`Regex("%v")("%v") = %v, this is wrong!`, v.Exp, v.String,
				ret)
		}
	}
}

func TestAnd(t *testing.T) {
	tests := []struct {
		String     string
		Validators []Validator
		Valid      bool
	}{
		{"Hey! 1", []Validator{Required("Req!")}, true},
		{"", []Validator{Required("Req!")}, false},
		{"Hey! 2", []Validator{Required("Req!"), Regex("Oink", "No way!")}, false},
		{"Hey! 3", []Validator{Required("Req!"), Regex("Hey", "No way!")}, true}}
	for _, v := range tests {
		ret := And(v.Validators...)(v.String)
		if (ret == nil && !v.Valid) || (ret != nil && v.Valid) {
			t.Errorf(`And(...)("%v") = %v, this is wrong!`, v.String, ret)
		}
	}
}

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

// testWidget performs common test on the given widget
//
// data is the form data struct.
// input is the rendered field with the given urlValue filled into the
// data struct.
// zeroInput is the rendered field with zero valued data struct field.
// value is the expected value in the data struct after filling the
// given input into it.
// urlValue is a HTTP parameter value to be filled into the data struct.
func testWidget(t *testing.T, widget Widget, data interface{}, input,
	zeroInput string, value interface{}, urlValue string) {
	form := NewForm(data, []Field{Field{"ID", "T", "H", nil, widget}})
	renderData := form.RenderData()
	if renderData.Fields[0].Input != template.HTML(zeroInput) {
		t.Errorf("Input field for zero value is\n%v\nshould be \n%v",
			renderData.Fields[0].Input, zeroInput)
	}
	vals := url.Values{"ID": []string{urlValue}}
	form.Fill(vals)
	renderData = form.RenderData()
	if renderData.Fields[0].Input != template.HTML(input) {
		t.Errorf("Input field is\n%v\nshould be \n%v",
			renderData.Fields[0].Input, input)
	}
	if !reflect.DeepEqual(value,
		reflect.ValueOf(data).Elem().FieldByName("ID").Interface()) {
		t.Errorf("Data is\n%v\nshould be \n%v", data, value)
	}
	form.Fill(url.Values{})
	renderData = form.RenderData()
	if renderData.Fields[0].Input != template.HTML(zeroInput) {
		t.Errorf("Input field for empty value is\n%v\nshould be \n%v",
			renderData.Fields[0].Input, zeroInput)
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
