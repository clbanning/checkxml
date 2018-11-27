package checkxml

import (
	"bytes"
	"encoding/xml"
	// "fmt"
	"testing"
)

func TestUnknownXMLTags(t *testing.T) {
	// fmt.Println("===================== TestUnknownXMLTags ...")

	data := []byte(`
		<Doc>
			<Ok>true</Ok>
			<Why attr="some val">
				<Maybe>true</Maybe>
				<maybenot>false</maybenot>
			</Why>
			<not>I dont't know</not>
		</Doc>`)

	check := map[string]bool{"Why.maybenot": true, "not": true, "Why.-attr": true}
	type test2 struct {
		Maybe bool
	}
	type test struct {
		XMLName xml.Name `xml:"Doc"`
		Ok      bool
		Why     test2
	}

	tv := test{}
	tags, _, err := UnknownXMLTags(data, tv)
	if err != nil {
		t.Fatal(err)
	}
	// fmt.Println(tags)

	// load tags into a map for checking that they're all there
	// and confirm that they were expected in the result set
	results := make(map[string]bool, 0)
	for _, v := range tags {
		results[v] = true
		if _, ok := check[v]; !ok {
			t.Fatal("unknown tag not in slice:", v)
		}
	}
	// now check that something didn't get in the result set
	for k := range check {
		if _, ok := results[k]; !ok {
			t.Fatal("unexpected tag in result set:", k)
		}
	}
}

func TestUnknownXMLTagsWithMemTags(t *testing.T) {
	// fmt.Println("===================== TestUnknownXMLTagsWithMemTags ...")

	data := []byte(`
		<doc>
			<ok>true</ok>
			<why attr="some val">
				<maybe>true</maybe>
				<maybenot>false</maybenot>
			</why>
			<not>I dont't know</not>
		</doc>`)

	check := map[string]bool{"why.maybenot": true, "not": true}
	type test2 struct {
		Val  bool   `xml:"maybe"`
		Attr string `xml:"attr,attr"`
	}
	type test struct {
		Yup bool  `xml:"ok"`
		Why test2 `xml:"why"`
	}

	tv := test{}
	tags, _, err := UnknownXMLTags(data, tv)
	if err != nil {
		t.Fatal(err)
	}
	// fmt.Println(tags)

	// load tags into a map for checking that they're all there
	// and confirm that they were expected in the result set
	results := make(map[string]bool, 0)
	for _, v := range tags {
		results[v] = true
		if _, ok := check[v]; !ok {
			t.Fatal("unknown tag not in slice:", v)
		}
	}
	// now check that something didn't get in the result set
	for k := range check {
		if _, ok := results[k]; !ok {
			t.Fatal("unexpected tag in result set:", k)
		}
	}
}

func TestUnknownXMLTagsReader(t *testing.T) {
	// fmt.Println("===================== TestUnknownXMLTagsReader ...")

	data := []byte(`
		<doc>
			<Ok>true</Ok>
			<Why attr="some val">
				<Maybe>true</Maybe>
				<maybenot>false</maybenot>
			</Why>
			<not>I dont't know</not>
		</doc>`)

	check := map[string]bool{"Why.maybenot": true, "not": true, "Why.-attr": true}
	type test2 struct {
		Maybe bool
	}
	type test struct {
		Ok  bool
		Why test2
	}

	tv := test{}
	buf := bytes.NewBuffer(data)
	tags, _, err := UnknownXMLTagsReader(buf, tv)
	if err != nil {
		t.Fatal(err)
	}

	// load tags into a map for checking that they're all there
	// and confirm that they were expected in the result set
	results := make(map[string]bool, 0)
	for _, v := range tags {
		results[v] = true
		if _, ok := check[v]; !ok {
			t.Fatal("unknown tag not in slice:", v)
		}
	}
	// now check that something didn't get in the result set
	for k := range check {
		if _, ok := results[k]; !ok {
			t.Fatal("unexpected tag in result set:", k)
		}
	}
}

func TestUnknownXMLTagsToIgnore(t *testing.T) {
	// fmt.Println("===================== TestUnknownXMLTagsToIgnore ...")

	data := []byte(`
		<doc>
			<Ok>true</Ok>
			<Why attr="some val">
				<Maybe>true</Maybe>
				<maybenot>false</maybenot>
			</Why>
			<not>I dont't know</not>
		</doc>`)

	SetTagsToIgnore("Why.maybenot", "not", "Why.-attr")
	defer SetTagsToIgnore()

	type test2 struct {
		Maybe bool
	}
	type test struct {
		Ok  bool
		Why test2
	}

	tv := test{}
	tags, _, err := UnknownXMLTags(data, tv)
	if err != nil {
		t.Fatal(err)
	}
	if len(tags) != 0 {
		t.Fatal("tags:", tags)
	}
}

func TestUnknownXMLTagsWithIgnoreTag(t *testing.T) {
	// fmt.Println("===================== TestUnknownXMLTagsWithIgnoreTag ...")

	data := []byte(`
		<doc>
			<Ok>true</Ok>
			<Why attr="some val">
				<Maybe>true</Maybe>
				<maybenot>false</maybenot>
			</Why>
			<not>I dont't know</not>
		</doc>`)

	check := map[string]bool{"Why.maybenot": true, "not": true, "Why.-attr": true}
	type test2 struct {
		Maybe bool `xml:"-"`
	}
	type test struct {
		Ok  bool
		Why test2
	}

	tv := test{}
	tags, _, err := UnknownXMLTags(data, tv)
	if err != nil {
		t.Fatal(err)
	}

	// load tags into a map for checking that they're all there
	// and confirm that they were expected in the result set
	results := make(map[string]bool, 0)
	for _, v := range tags {
		results[v] = true
		if _, ok := check[v]; !ok {
			t.Fatal("unknown tag not in slice:", v)
		}
	}
	// now check that something didn't get in the result set
	for k := range check {
		if _, ok := results[k]; !ok {
			t.Fatal("unexpected tag in result set:", k)
		}
	}
}

func TestUnknownXMLTagsWithSubelemTag(t *testing.T) {
	// fmt.Println("===================== TestUnknownXMLTagsWithSubelemTag ...")

	data := []byte(`
		<doc>
			<Ok>true</Ok>
			<Why attr="some val">
				<Maybe>true</Maybe>
				<maybenot>false</maybenot>
			</Why>
			<not>I dont't know</not>
		</doc>`)

	check := map[string]bool{"Why.maybenot": true, "not": true, "Why.-attr": true}
	type test2 struct {
		Maybe bool `xml:"-"`
	}
	type test struct {
		Ok  bool
		Why test2 `xml:"Why>Maybe"`
	}

	tv := test{}
	tags, _, err := UnknownXMLTags(data, tv)
	if err != nil {
		t.Fatal(err)
	}

	// load tags into a map for checking that they're all there
	// and confirm that they were expected in the result set
	results := make(map[string]bool, 0)
	for _, v := range tags {
		results[v] = true
		if _, ok := check[v]; !ok {
			t.Fatal("unknown tag not in slice:", v)
		}
	}
	// now check that something didn't get in the result set
	for k := range check {
		if _, ok := results[k]; !ok {
			t.Fatal("unexpected tag in result set:", k)
		}
	}
}

// ===================== 11/27/18: handle single member slices correctly =============
// thanks to: zhengfang.sun sunsun314 (github)

func TestUnknownTagsSingletonList(t *testing.T) {
	var b = []byte(`<yy>
	<xx>1</xx>
</yy>`)

	var d = []byte(`<yy>
	<zz>1</zz>
</yy>`)

	type Y struct {
		Property []string `xml:"xx"`
	}

	y := Y{}
	tags, _, err := UnknownXMLTags(b, y)
	if err != nil {
		t.Fatal("err on b:", err.Error())
	} else if len(tags) != 0 {
		t.Fatal("reported tags for b", tags)
	}

	y = Y{}
	tags, _, err = UnknownXMLTags(d, y)
	if err != nil {
		t.Fatal(err.Error())
	} else if tags[0] != "zz" {
		t.Fatal("didn't report 'zz' for d:", tags)
	}
}
