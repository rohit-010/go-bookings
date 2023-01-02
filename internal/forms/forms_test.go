package forms

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestForms_Valid(t *testing.T) {
	r := httptest.NewRequest("POST", "/test-url", nil)
	form := New(r.PostForm)

	isValid := form.Valid()
	if !isValid {
		t.Error("got invalid when should have been valid form")
	}
}

func TestForm_Required(t *testing.T) {
	r := httptest.NewRequest("POST", "/test-url", nil)
	form := New(r.PostForm)
	form.Required("a", "b", "c")

	if form.Valid() {
		t.Error("form shows valid when required fields are missing")
	}

	postedData := url.Values{}
	postedData.Add("a", "a")
	postedData.Add("b", "a")
	postedData.Add("c", "a")

	r, _ = http.NewRequest("POST", "/test-url", nil)
	r.PostForm = postedData
	form = New(r.PostForm)
	form.Required("a", "b", "c")
	if !form.Valid() {
		t.Error("shows does not have required fields when it does")
	}
}

func TestForm_Has(t *testing.T) {
	postedData := url.Values{}
	form := New(postedData)
	has := form.Has("whatever")
	if has {
		t.Error("forms shows has field when it does not")
	}

	postedData = url.Values{}
	postedData.Add("a", "a")
	form = New(postedData)

	has = form.Has("a")
	if !has {
		t.Error("shows form does not have field when it should")
	}
}

func TestForm_MinLength(t *testing.T) {
	postedData := url.Values{}
	form := New(postedData)

	form.MinLength("xx", 10)
	if form.Valid() {
		t.Error("shows min length for non existed field")
	}

	isError := form.Errors.Get("xx")
	if isError == "" {
		t.Error("should have an error but did not get one")
	}

	postedData = url.Values{}
	postedData.Add("some_field", "some_value")
	form = New(postedData)
	form.MinLength("some_field", 100)
	if form.Valid() {
		t.Error("shows min length of 100 met when input data is shorter")
	}

	postedData = url.Values{}
	postedData.Add("another_field", "abc123")
	form = New(postedData)
	form.MinLength("another_field", 1)
	if !form.Valid() {
		t.Error("shows min length of 1 is not met when it should")
	}

	isError = form.Errors.Get("another_field")
	if isError != "" {
		t.Error("should not have an error but got one")
	}

}

func TestForm_IsEmail(t *testing.T) {
	postedData := url.Values{}
	form := New(postedData)

	form.IsEmail("x")
	if form.Valid() {
		t.Error("form shows valid email for non existent field")
	}

	postedData = url.Values{}
	postedData.Add("email", "email@email.com")
	form = New(postedData)
	form.IsEmail("email")

	if !form.Valid() {
		t.Error("got an invalid email when we should not have")
	}

	postedData = url.Values{}
	postedData.Add("email", "test")
	form = New(postedData)
	form.IsEmail("email")

	if form.Valid() {
		t.Error("got an valid email for invalid email address")
	}
}
