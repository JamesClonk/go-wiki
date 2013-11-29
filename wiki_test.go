package main

import "testing"
import "os"
import "strings"
import "html/template"
import "net/url"
import "net/http"
import "net/http/httptest"

var testpage = "Unittest_SamplePage"

func cleanup(t *testing.T) {
	filename := dataDirectory + "/" + testpage + ".txt"
	os.Remove(filename)
	_, err := os.Open(filename)
	if err == nil {
		t.Errorf("The test file [%v] was not properly deleted!", filename)
		return
	}
}

func Test__setup(t *testing.T) {
	isLogging = false
	checkDirectories()
	cleanup(t)
}

func Test_getFilename(t *testing.T) {
	expected := dataDirectory + "/hello.txt"
	filename := getFilename("hello")
	if filename != expected {
		t.Errorf("filename is not as expected: [%v], instead of [%v]", filename, expected)
	}
}

func Test_savePage(t *testing.T) {
	title := testpage
	p1 := &Page{Title: title, Body: template.HTML("Hello World!")}

	err := p1.save()
	if err != nil {
		t.Error(err)
	}
}

func Test_loadPage(t *testing.T) {
	title := testpage
	p2, err := loadPage(title)
	if err != nil {
		t.Error(err)
		return
	}

	expected := title
	if p2.Title != expected {
		t.Errorf("loaded page title is not as expected: [%v], instead of [%v]", p2.Title, expected)
		return
	}

	expected = "Hello World!"
	if string(p2.Body) != expected {
		t.Errorf("loaded page body is not as expected: [%v], instead of [%v]", p2, expected)
		return
	}
}

func Test_getTitle(t *testing.T) {
	request, err := http.NewRequest("GET", "http://localhost:8008/view/Invalid!!"+testpage, nil)
	if err != nil {
		t.Error(err)
		return
	}
	response := httptest.NewRecorder()

	getTitle(response, request)
	code := response.Code
	if code != 404 {
		t.Errorf("Test_getTitle() response code was [%v], but expected [%v]", code, 404)
	}

	request, err = http.NewRequest("GET", "http://localhost:8008/view/"+testpage, nil)
	if err != nil {
		t.Error(err)
		return
	}
	response = httptest.NewRecorder()

	title, err := getTitle(response, request)
	if err != nil {
		t.Error(err)
		return
	}
	if title != testpage {
		t.Errorf("Test_getTitle() return value was [%v], but expected [%v]", title, testpage)
		return
	}
	code = response.Code
	if code != 200 {
		t.Errorf("Test_getTitle() response code was [%v], but expected [%v]", code, 200)
	}
}

func Test_indexHandler(t *testing.T) {
	request, err := http.NewRequest("GET", "http://localhost:8008/", nil)
	if err != nil {
		t.Error(err)
		return
	}
	response := httptest.NewRecorder()

	indexHandler(response, request)
	code := response.Code
	header := response.Header()

	if code != 302 {
		t.Errorf("indexHandler() response code was [%v], but expected [%v]", code, 302)
	}

	location := header.Get("Location")
	expected := "/view/" + defaultPage
	if location != expected { // HTTP redirect
		t.Errorf("indexHandler() response location header was [%v], but expected [%v]", location, expected)
	}
}

func Test_viewHandler(t *testing.T) {
	request, err := http.NewRequest("GET", "http://localhost:8008/view/"+testpage, nil)
	if err != nil {
		t.Error(err)
		return
	}
	response := httptest.NewRecorder()

	viewHandler(response, request, testpage)
	code := response.Code
	body := response.Body.String()

	if code != 200 {
		t.Errorf("viewHandler() response code was [%v], but expected [%v]", code, 200)
	}

	shouldContain := "<a href=\"/edit/" + testpage + "\">edit</a>"
	if !strings.Contains(body, shouldContain) {
		t.Errorf("viewHandler() response body was [%v], but expected it to contain [%v]", body, shouldContain)
	}

	shouldContain = "<div>Hello World!</div>"
	if !strings.Contains(body, shouldContain) {
		t.Errorf("viewHandler() response body was [%v], but expected it to contain [%v]", body, shouldContain)
	}
}

func Test_editHandler(t *testing.T) {
	request, err := http.NewRequest("GET", "http://localhost:8008/edit/"+testpage, nil)
	if err != nil {
		t.Error(err)
		return
	}
	response := httptest.NewRecorder()

	editHandler(response, request, testpage)
	code := response.Code
	body := response.Body.String()

	if code != 200 {
		t.Errorf("editHandler() response code was [%v], but expected [%v]", code, 200)
	}

	shouldContain := "<h1>Editing " + testpage + "</h1>"
	if !strings.Contains(body, shouldContain) {
		t.Errorf("editHandler() response body was [%v], but expected it to contain [%v]", body, shouldContain)
	}

	shouldContain = "<form action=\"/save/" + testpage + "\" method=\"POST\">"
	if !strings.Contains(body, shouldContain) {
		t.Errorf("editHandler() response body was [%v], but expected it to contain [%v]", body, shouldContain)
	}

	shouldContain = "<textarea name=\"body\" rows=\"20\" cols=\"80\">Hello World!</textarea>"
	if !strings.Contains(body, shouldContain) {
		t.Errorf("editHandler() response body was [%v], but expected it to contain [%v]", body, shouldContain)
	}
}

func Test_saveHandler(t *testing.T) {
	request, err := http.NewRequest("POST", "http://localhost:8008/save/"+testpage, nil)
	if err != nil {
		t.Error(err)
		return
	}
	request.PostForm = url.Values{"body": {"Hallo Welt!"}}
	response := httptest.NewRecorder()

	saveHandler(response, request, testpage)
	code := response.Code
	header := response.Header()

	if code != 302 { // HTTP redirect
		t.Errorf("saveHandler() response code was [%v], but expected [%v]", code, 302)
	}
	location := header.Get("Location")
	expected := "/view/" + testpage
	if location != expected { // HTTP redirect
		t.Errorf("saveHandler() response location header was [%v], but expected [%v]", location, expected)
	}

	// check by requesting /view/ again to see if it reflects the changes made
	request, err = http.NewRequest("GET", "http://localhost:8008/view/"+testpage, nil)
	if err != nil {
		t.Error(err)
		return
	}
	response = httptest.NewRecorder()

	viewHandler(response, request, testpage)
	code = response.Code
	body := response.Body.String()

	if code != 200 {
		t.Errorf("saveHandler() -> viewHandler() response code was [%v], but expected [%v]", code, 200)
	}
	shouldContain := "Hallo Welt!"
	if !strings.Contains(body, shouldContain) {
		t.Errorf("saveHandler() -> viewHandler() response body was [%v], but expected it to contain [%v]", body, shouldContain)
	}
}

func Test_makeHandler(t *testing.T) {
	var test string
	testFunc := func(w http.ResponseWriter, r *http.Request, title string) {
		test = title
	}
	makeHandlerFunc := makeHandler(testFunc)

	request, err := http.NewRequest("GET", "http://localhost:8008/view/"+testpage, nil)
	if err != nil {
		t.Error(err)
		return
	}
	response := httptest.NewRecorder()

	makeHandlerFunc(response, request)

	if test != testpage {
		t.Errorf("makeHandler() test was [%v], but expected [%v]", test, testpage)
	}
}

func Test__teardown(t *testing.T) {
	cleanup(t)
}
