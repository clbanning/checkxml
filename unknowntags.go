// unknowntags.go - identify XML data elements that will not be unmarshalled to a struct field
// Copyright © 2017 Charles Banning.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package checkxml

import (
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/clbanning/mxj"
)

// UnknownXMLTags returns a slice of the tags for XML data elements or attributes
// that will not be decoded to a member of 'val', which is of type struct along with
// the XML data root tag.
// For complex elements the tags are reported using dot-notation.
// Attribute tags are prepended with a hyphen symbol, "-", the clbanning/mxj
// package convention.
//	Examples:
//		data1 := `<doc>
//		            <e1>test</e1>
//		            <e2>more</e2>
//		          </doc>`
//		type doc1 struct {
//			E1 string `xml:"e1"`
//		}
//
//		doc1 := doc1{}
//		tags, _, _ := UnknownXMLTags([]byte(data1), doc2)
//		fmt.Println(tags) // prints: [e2]
//
//		data2 := `<doc>
//		            <e1>test</e1>
//		            <subdoc>
//		              <e1>test</e1>
//		              <e2>more</e2>
//		            </subdoc>
//		          </doc>
//		type subdoc struct {
//			E1 string `xml:"e1"`
//		}
//		type doc2 struct {
//			E1  string `xml:"e1"`
//			Sub subdoc `xml:"subdoc"`
//		}
//
//		doc2 := doc2{}
//		tags, _, _ := UnknownXMLTags([]byte(data2), doc2)
//		fmt.Println(tags) // prints: [subdoc.e2]
//
//		data3 := `<doc>
//		            <e1 attr="something">test</e1>
//		            <e2>more</e2>
//		          </doc>`
//		type doc3 struct {
//			E1 string `xml:"e1"`
//		}
//
//		doc3 := doc3{}
//		tags, _, _ := UnknownXMLTags([]byte(data3), doc3)
//		fmt.Println(tags) // prints: [e1.-attr e2]
//
// If a struct member has a "-" XML tag and a corresponding tag with the
// struct member name occurs in the XML data, the XML tag will not be included
// as part of the unknown tags since it is "known" as part of the struct
// definition even if it won't be decoded by the encoding/xml package.
//
// NOTE: dot-notation XML tag values returned by UnknownXMLTags use the
// struct member `xml` tag or the public field name if there is no `xml` tag.
// This allows the members of the returned slice to be directly used with
// the mxj package if the mxj.Map representation of the XML data is available..
// (See github.com/clbanning/mxj documentation of mxj.Map type.)
//	Example - print out XML data tags and values that will not be decoded to the struct "myStruct":
//	import "github.com/clbanning/mxj"
//	...
//		tags, root, m, err := UnknownXMLTagsMap(xmlData, myStruct)
//		if err != nil {
//		   // handle error
//		}
//		for _, tag := range tags {
//		   fmt.Printf("%s: %#v\n", tag, m.ValuesForPath(root+"."+tag))
//		}
func UnknownXMLTags(b []byte, val interface{}) ([]string, string, error) {
	var s []string

	m, err := mxj.NewMapXml(b)
	if err != nil {
		return nil, "", err
	}
	// strip the root tag and seed 'key'
	var root string
	var v interface{}
	for root, v = range m {
		break
	}

	if _, ok := v.(map[string]interface{}); !ok {
		if _, ok = v.([]interface{}); !ok {
			// nothing to work with just return the root key
			return s, root, fmt.Errorf("no elements")
		}
	}

	checkAllTags(v, reflect.ValueOf(val), &s, "")
	return s, root, nil
}

// UnknownXMLTagsMap returns the mxj.Map - map[string]interface{} - representation
// of the XML data in addition to the unknown XML tags and the XML data root tag.
// (See github.com/clbanning/mxj documentation of mxj.Map type.)
func UnknownXMLTagsMap(b []byte, val interface{}) ([]string, string, mxj.Map, error) {
	var s []string

	m, err := mxj.NewMapXml(b, mxjCast)
	if err != nil {
		return nil, "", nil, err
	}
	// strip the root tag and seed 'key'
	var root string
	var v interface{}
	for root, v = range m {
		break
	}
	if _, ok := v.(map[string]interface{}); !ok {
		if _, ok = v.([]interface{}); !ok {
			// nothing to work with just return the root key
			return s, root, m, fmt.Errorf("no elements")
		}
	}

	checkAllTags(v, reflect.ValueOf(val), &s, "")
	return s, root, m, nil
}

// ================= io.Reader functions ...

// UnknownXMLTagsReader consumes the XML data from an io.Reader and returns
// the XML tags that are unknown with respect to the struct 'val' and the XML data
// root tag.
func UnknownXMLTagsReader(r io.Reader, val interface{}) ([]string, string, error) {
	var s []string

	m, err := mxj.NewMapXmlReader(r)
	if err != nil {
		return nil, "", err
	}
	// strip the root tag and seed 'key'
	var root string
	var v interface{}
	for root, v = range m {
		break // just for safety
	}

	if _, ok := v.(map[string]interface{}); !ok {
		if _, ok = v.([]interface{}); !ok {
			// nothing to work with just return the root key
			return s, root, fmt.Errorf("no elements")
		}
	}

	checkAllTags(v, reflect.ValueOf(val), &s, "")
	return s, root, nil
}

// UnknownXMLTagsReaderMap consumes the XML data from an io.Reader and returns
// the mxj.Map - map[string]interface{} - representation of the XML data in addition
// to the unknown XML tags and the XML data root tag. 
// (See github.com/clbanning/mxj documentation of mxj.Map type.)
func UnknownXMLTagsReaderMap(r io.Reader, val interface{}) ([]string, string, mxj.Map, error) {
	var s []string

	m, err := mxj.NewMapXmlReader(r, mxjCast)
	if err != nil {
		return nil, "", m, err
	}
	// strip the root tag and seed 'key'
	var root string
	var v interface{}
	for root, v = range m {
		break // just for safety
	}

	if _, ok := v.(map[string]interface{}); !ok {
		if _, ok = v.([]interface{}); !ok {
			// nothing to work with just return the root key
			return s, root, m, fmt.Errorf("no elements")
		}
	}

	checkAllTags(v, reflect.ValueOf(val), &s, "")
	return s, root, m, nil
}

// UnknownXMLTagsReaderMapRaw consumes the XML data from an io.Reader and returns
// the raw XML data that was processed in addition to the unknown element tags,
// the mxj.Map - map[string]interface{} - representation of the XML data, and the XML
// data root tag.
// (See github.com/clbanning/mxj documentation of mxj.Map type.)
func UnknownXMLTagsReaderMapRaw(r io.Reader, val interface{}) ([]string, string, mxj.Map, []byte, error) {
	var s []string

	m, raw, err := mxj.NewMapXmlReaderRaw(r, mxjCast)
	if err != nil {
		return nil, "", m, raw, err
	}
	// strip the root tag and seed 'key'
	var root string
	var v interface{}
	for root, v = range m {
		break // just for safety
	}
	if _, ok := v.(map[string]interface{}); !ok {
		if _, ok = v.([]interface{}); !ok {
			// nothing to work with just return the root key
			return s, root, m, raw, fmt.Errorf("no elements")
		}
	}
	checkAllTags(v, reflect.ValueOf(val), &s, "")
	return s, root, m, raw, nil
}

// ================== where the work is done ...

func checkAllTags(mv interface{}, val reflect.Value, s *[]string, key string) {
	var tkey string

	// 1. Convert any pointer value.
	if val.Kind() == reflect.Ptr {
		val = reflect.Indirect(val)
	}
	// zero Value?
	if !val.IsValid() {
		return
	}
	typ := val.Type()

	// 2. If its a slice then 'mv' should hold a []interface{} value.
	//    Loop through the members of 'mv' and see that they are valid relative
	//    to the <T> of val []<T>.
	if typ.Kind() == reflect.Slice {
		tval := typ.Elem()
		if tval.Kind() == reflect.Ptr {
			tval = tval.Elem()
		}
		// slice may be nil, so create a Value of it's type
		// 'mv' must be of type []interface{}
		sval := reflect.New(tval)
		slice, ok := mv.([]interface{})
		if !ok {
			*s = append(*s, key)
			return
		}
		// 2.1. Check members of XML data
		//      This forces all of them to be regular and w/o typos in key labels.
		for _, sl := range slice {
			checkAllTags(sl, sval, s, key) // all list elements have same tag
		}
		return
	}

	// 3a. Ignore anything that's not a struct.
	if typ.Kind() != reflect.Struct {
		return // just ignore it - don't look for k:v pairs
	}
	// 3b. map value must represent k:v pairs
	mm, ok := mv.(map[string]interface{})
	if !ok {
		*s = append(*s, key)
	}

	// 4. Build the map of struct field name:fieldSpec
	//    We make every key (field) label look like an exported label - "Fieldname".
	//    If there is a XML tag it is used instead of the field label, and saved to
	//    insure that the spec'd tag matches the XML tag exactly.
	type fieldSpec struct {
		val reflect.Value
		tag []string // tag may be a path
	}
	fieldCnt := val.NumField()
	fields := make(map[string]*fieldSpec, fieldCnt)
	for i := 0; i < fieldCnt; i++ {
		if len(typ.Field(i).PkgPath) > 0 {
			continue // field is NOT exported
		}
		// Ignore xml.Name type fields - they don't appear in the map mm.
		// The root label is handed in as "key" in the initial call.
		if typ.Field(i).Type.Name() == "Name" && typ.Field(i).Type.PkgPath() == "encoding/xml" {
			continue
		}
		// Get tag and attr info from member spec
		// A go xml tag may be a single label, e.g., "elem",
		// or it may be a path to a subelement, e.g., "elem>sub>stuff",
		// see: https://golang.org/pkg/encoding/xml/#example_Unmarshal.
		// We just ignore the rest of the path for now - see discussion below in #5.
		attr := false
		tagvals := typ.Field(i).Tag.Get("xml")
		tags := strings.Split(tagvals, ",")
		tag := strings.Split(tags[0], ">")
		// Fields with "-" might, validly, be there
		// so allow the field name to be included.
		if tag[0] == "-" {
			tag = []string{""}
		}
		// See if struct member is an attribute value.
		for _, v := range tags[1:] {
			if v == "attr" {
				attr = true
				break
			}
		}
		// If attr==true then the mm key will be prepended with "-"
		// so the Field name and the 'tag' value must be prepended with "-"
		// to match the decoded value.
		// NOTE: the xml decoder requires that elem/attr tags match exactly
		// the public member name or its xml tag label; unlike json decoder
		// there is no coersion of lower case element tags to public
		// member names.
		switch attr {
		case false:
			if tag[0] == "" {
				fields[typ.Field(i).Name] = &fieldSpec{val.Field(i), tag}
			} else {
				fields[tag[0]] = &fieldSpec{val.Field(i), tag}
			}
		case true:
			if tag[0] == "" {
				fields["-"+typ.Field(i).Name] = &fieldSpec{val.Field(i), tag}
			} else {
				fields["-"+tag[0]] = &fieldSpec{val.Field(i), tag}
			}
		}
	}

	// 5. check that map keys correspond to exported field names
	//    We handle the keys in the map literally, unlike for encoding/json.
	var spec *fieldSpec
	for k, m := range mm {
		for _, sk := range skiptags {
			if key == "" && k == sk {
				goto next
			} else if key != "" && key+"."+k == sk {
				goto next
			}
		}
		// used for !ok and recursion on checkAllTags
		if key == "" {
			tkey = k
		} else {
			tkey = key + "." + k
		}
		spec, ok = fields[k]
		if !ok {
			*s = append(*s, tkey)
			continue
		}
		// todo(clb): resolve how to handle subelement xml tags.
		// Do we even need to for unknown tags? -
		// Perhaps not, as the decoder must be able to walk the path per the struct
		// definition.  MissingXMLTags() can be used to see if the desired path can
		// be walked in the XML data, that result can then be used to see if the
		// desired subelement path is in the XML data. Something like:
		// 	subelemtag = "doc.elem.text" // we've replace ">" with "."
		// 	mems, _ := MissingXMLTags(...)
		// 	for _, v := range mems {
		// 		if subelemtag == v {
		// 			fmt.Println("subelement xml tag does not exist in XML data:", subelemtag)
		// 		}
		// 	}
		//
		checkAllTags(m, spec.val, s, tkey)
	next:
	}

	return
}
