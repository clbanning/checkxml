package checkxml

import (
	"bytes"
	"fmt"
	"testing"
)

func TestMissingXMLTags(t *testing.T) {
	// fmt.Println("===================== TestMissingXMLTags ...")

	type test struct {
		Ok  bool   `xml:"ok"`
		Why string `xml:"why"`
	}
	tv := test{}
	data := []byte(`<doc><ok>true</ok><why>it's a test</why></doc>`)
	mems, _, err := MissingXMLTags(data, tv)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if len(mems) > 0 {
		t.Fatalf(fmt.Sprintf("len(mems) == %d >> %v", len(mems), mems))
	}

	data = []byte(`<doc><ok>true</ok></doc>`)
	check := map[string]bool{"why": true} // use XML tag, not member name
	mems, _, err = MissingXMLTags(data, tv)
	if err != nil {
		t.Fatalf(err.Error())
	}
	results := make(map[string]bool,0)
	for _, v := range mems {
		results[v] = true
		if _, ok := check[v]; !ok {
			t.Fatalf(fmt.Sprintf("missing member not in checklist: %s", v))
		}
	}
	// now check that something didn't get in the result set
	for k := range check {
		if _, ok := results[k]; !ok {
			t.Fatal("unexpected tag in result set:", k)
		}
	}

	data = []byte(`<doc></doc>`)
	check = map[string]bool{"test": true} // the struct type
	mems, _, err = MissingXMLTags(data, tv)
	if err != nil {
		t.Fatalf(err.Error())
	}
	results = make(map[string]bool,0)
	for _, v := range mems {
		results[v] = true
		if _, ok := check[v]; !ok {
			t.Fatalf(fmt.Sprintf("missing member not in checklist: %s", v))
		}
	}
	// now check that something didn't get in the result set
	for k := range check {
		if _, ok := results[k]; !ok {
			t.Fatal("unexpected tag in result set:", k)
		}
	}
}

func TestMissingXMLTagsMap(t *testing.T) {
	// Println("===================== TestMissingXMLTagsMap ...")

	type test struct {
		Ok  bool   `xml:"ok"`
		Why string `xml:"why"`
	}
	tv := test{}
	data := []byte(`<doc><ok>true</ok><why>it's a test</why></doc>`)
	mems, _, _, err := MissingXMLTagsMap(data, tv)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if len(mems) > 0 {
		t.Fatalf(fmt.Sprintf("len(mems) == %d >> %v", len(mems), mems))
	}
}

func TestMissingXMLTagsReader(t *testing.T) {
	// fmt.Println("===================== TestMissingXMLTagsReader ...")

	type test struct {
		Ok  bool   `xml:"ok"`
		Why string `xml:"why"`
	}
	tv := test{}
	data := []byte(`<doc><ok>true</ok><why>it's a test</why></doc>`)
	r := bytes.NewBuffer(data)
	mems, _, err := MissingXMLTagsReader(r, tv)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if len(mems) > 0 {
		t.Fatalf(fmt.Sprintf("len(mems) == %d >> %v", len(mems), mems))
	}

	data = []byte(`<doc><ok>true</ok></doc>`)
	check := map[string]bool{"why": true} // use XML tag, not member name
	r = bytes.NewBuffer(data)
	mems, _, err = MissingXMLTagsReader(r, tv)
	if err != nil {
		t.Fatalf(err.Error())
	}
	results := make(map[string]bool,0)
	for _, v := range mems {
		results[v] = true
		if _, ok := check[v]; !ok {
			t.Fatalf(fmt.Sprintf("missing member not in checklist: %s", v))
		}
	}
	// now check that something didn't get in the result set
	for k := range check {
		if _, ok := results[k]; !ok {
			t.Fatal("unexpected tag in result set:", k)
		}
	}
}

func TestMissingXMLTagsReaderMap(t *testing.T) {
	// fmt.Println("===================== TestMissingXMLTagsReaderMap ...")

	type test struct {
		Ok  bool   `xml:"ok"`
		Why string `xml:"why"`
	}
	tv := test{}
	data := []byte(`<doc><ok>true</ok><why>it's a test</why></doc>`)
	r := bytes.NewBuffer(data)
	mems, _, _, err := MissingXMLTagsReaderMap(r, tv)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if len(mems) > 0 {
		t.Fatalf(fmt.Sprintf("len(mems) == %d >> %v", len(mems), mems))
	}
}

func TestMissingXMLTagsReaderMapRaw(t *testing.T) {
	// fmt.Println("===================== TestMissingXMLTagsReaderMapRaw ...")

	type test struct {
		Ok  bool   `xml:"ok"`
		Why string `xml:"why"`
	}
	tv := test{}
	s := "<doc><ok>true</ok><why>it's a test</why></doc>"
	data := []byte(s)
	r := bytes.NewBuffer(data)
	mems, _, _, raw, err := MissingXMLTagsReaderMapRaw(r, tv)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if string(raw) != string(data) {
		t.Fatalf("raw not data:\n%T, %d> %s\n%T, %d> %s", raw, len(raw), string(raw), data, len(data), string(data))
	}
	if len(mems) > 0 {
		t.Fatalf(fmt.Sprintf("len(mems) == %d >> %v", len(mems), mems))
	}
}

func TestMissingXMLTagsSubElements(t *testing.T) {
	// fmt.Println("===================== TestMissingXMLTagsSubElements ...")
	type test3 struct {
		Something string `xml:"something"`
		Else      string `xml:"else"`
	}

	type test2 struct {
		Why     string
		Not     string
		Another test3
	}

	type test struct {
		Ok   bool   `xml:"ok"`
		Why  string `xml:"why"`
		More test2  `xml:"more"`
	}
	tv := test{}
	data := []byte(`<doc>
		<ok>true</ok>
		<why>it's a test</why>
		<more>
			<Why>again</Why>
			<Not>no more</Not>
			<Another>
				<something>a thing</something>
				<else>ok</else>
			</Another>
		</more>
	</doc>`)
	mems, _, err := MissingXMLTags(data, tv)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if len(mems) > 0 {
		t.Fatalf(fmt.Sprintf("len(mems) == %d >> %v", len(mems), mems))
	}

	data = []byte(`<doc>
		<ok>true</ok>
		<more>
			<Why>again</Why>
			<Another>
				<else>ok"</else>
			</Another>
		</more>
	</doc>`)
	// use XML tag values instead of struct member names where appropriate
	check := map[string]bool{"why": true, "more.Not": true, "more.Another.something": true}
	mems, _, err = MissingXMLTags(data, tv)
	if err != nil {
		t.Fatalf(err.Error())
	}
	results := make(map[string]bool, 0)
	for _, v := range mems {
		results[v] = true
		if _, ok := check[v]; !ok {
			t.Fatalf(fmt.Sprintf("missing member not in checklist: %s", v))
		}
	}
	// now check that something didn't get in the result set
	for k := range check {
		if _, ok := results[k]; !ok {
			t.Fatal("unexpected tag in result set:", k)
		}
	}
}

func TestMissingXMLTagsSkipMems(t *testing.T) {
	// fmt.Println("===================== TestMissingXMLTagsSkipMems ...")

	type test3 struct {
		Something string `xml:"something"`
		Else      string `xml:"else"`
	}

	type test2 struct {
		Why     string
		Not     string
		Another test3
	}

	type test struct {
		Ok   bool   `xml:"ok"`
		Why  string `xml:"why"`
		More test2  `xml:"more"`
	}
	tv := test{}
	data := []byte(`<doc>
		<ok>true</ok>
		<more>
			<Why>again</Why>
			<Another>
				<else>ok</else>
			</Another>
		</more>
	</doc>`)

	SetMembersToIgnore("why", "more.Not", "more.Another.something")
	defer SetMembersToIgnore()

	mems, _, err := MissingXMLTags(data, tv)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if len(mems) != 0 {
		t.Fatalf(fmt.Sprintf("missing mems: %d - %#v", len(mems), mems))
	}
}
