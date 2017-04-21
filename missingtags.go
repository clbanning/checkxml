// missingtags.go - check XML data for tags that are missing from struct definition.
// Copyright © 2017 Charles Banning.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package checkxml

import (
	"io"
	"reflect"
	"strings"

	"github.com/clbanning/mxj"
)

// MissingXMLTags returns a slice of members of val, the struct definition, that will NOT be set
// by unmarshaling the XML-encoded data; rather, they will assume their initialization or default
// values. For nested structs, member labels are the dot-notation hierachical
// path for the missing XML tag.
// The XML root tag for XML data, 'b', that was scanned is also returned.
// Specific struct members can be igored when scanning the XML object by declaring them using
// SetMembersToIgnore().
//
//	Examples:
//		data1 := `<doc>
//		            <e1>test</e1>
//		          </doc>`
//		type doc1 struct {
//			E1 string `xml:"e1"`
//			E2 string `xml:"e2"`
//		}
//
//		doc1 := doc1{}
//		tags, _, _ := MissingXMLTags([]byte(data1), doc2)
//		fmt.Println(tags) // prints: [e2]
//
//		data2 := `<doc>
//		            <e1>test</e1>
//		            <subdoc>
//		              <e1>test</e1>
//		            </subdoc>
//		          </doc>
//		type subdoc struct {
//			E1 string `xml:"e1"`
//			E2 string `xml:"e2"`
//		}
//		type doc2 struct {
//			E1  string `xml:"e1"`
//			E2  string `xml:"e2"`
//			Sub subdoc `xml:"subdoc"`
//		}
//
//		doc2 := doc2{}
//		tags, _, _ := MissingXMLTags([]byte(data2), doc2)
//		fmt.Println(tags) // prints: [e2 subdoc.e2]
//
// By default missing tags in the XML data that are associated with struct members that
// have XML tags "-" and "omitempty" are not included in the returned slice.
// IgnoreOmitemptyTag(false) can be called to override the handling of "omitempty"
// tags - this might be useful if you want to find the "omitempty" members that
// are not set by decoding the XML data..
//
// NOTE: dot-notation XML tag values returned by MissingXMLTags use the
// struct member `xml` tag or the public field name if there is no `xml` tag.
// This allows the members of the returned slice to be used to directly manipulate a mxj.Map
// representation of the XML data if it is available.
// (See github.com/clbanning/mxj documentation of mxj.Map type.)
func MissingXMLTags(b []byte, val interface{}) ([]string, string, error) {
	var s []string

	m, err := mxj.NewMapXml(b)
	if err != nil {
		return nil, "", err
	}
	// strip off the root value
	var root string
	var v interface{}
	for _, v = range m {
		break
	}

	vv, ok := v.(map[string]interface{})
	if !ok {
		if _, ok = v.([]interface{}); !ok {
			// return the name of the value passed if not a map[string]interface{} value
			s = append(s, reflect.ValueOf(val).Type().Name())
			return s, root, nil
		}
	}

	checkMembers(vv, reflect.ValueOf(val), &s, "")
	return s, root, nil
}

// MissingXMLTagsMap returns the mxj.Map - map[string]interface{} - representation of the XML data
// and the XML root tag in addition to the missing XML tags.
// (See github.com/clbanning/mxj documentation of mxj.Map type.)
func MissingXMLTagsMap(b []byte, val interface{}) ([]string, mxj.Map, string, error) {
	var s []string

	m, err := mxj.NewMapXml(b, mxjCast)
	if err != nil {
		return nil, m, "", err
	}
	// strip off the root value
	var root string
	var v interface{}
	for root, v = range m {
		break
	}

	vv, ok := v.(map[string]interface{})
	if !ok {
		if _, ok = v.([]interface{}); !ok {
			// return the name of the value passed if not a map[string]interface{} value
			s = append(s, reflect.ValueOf(val).Type().Name())
			return s, m, root, nil
		}
	}

	checkMembers(vv, reflect.ValueOf(val), &s, "")
	return s, m, root, nil
}

// ================= io.Reader functions ...

// MissingXMLTagsReader consumes the XML data from an io.Reader and returns the XML tags
// that are missing with respect to the struct 'val' and the XML root tag.
func MissingXMLTagsReader(r io.Reader, val interface{}) ([]string, string, error) {
	var s []string

	m, err := mxj.NewMapXmlReader(r)
	if err != nil {
		return nil, "", err
	}
	// strip off the root value
	var root string
	var v interface{}
	for root, v = range m {
		break
	}

	vv, ok := v.(map[string]interface{})
	if !ok {
		if _, ok = v.([]interface{}); !ok {
			// return the name of the value passed if not a map[string]interface{} value
			s = append(s, reflect.ValueOf(val).Type().Name())
			return s, root, nil
		}
	}

	checkMembers(vv, reflect.ValueOf(val), &s, "")
	return s, root, nil
}

// MissingXMLTagsReaderMap consumes the XML data from an io.Reader and returns the
// mxj.Map - map[string]interface{} - representation of the XML data and the root
// XML tag in addition to the missing XML tags.
// (See github.com/clbanning/mxj documentation of mxj.Map type.)
func MissingXMLTagsReaderMap(r io.Reader, val interface{}) ([]string, mxj.Map, string, error) {
	var s []string

	m, err := mxj.NewMapXmlReader(r, mxjCast)
	if err != nil {
		return nil, m, "", err
	}
	// strip off the root value
	var root string
	var v interface{}
	for root, v = range m {
		break
	}

	vv, ok := v.(map[string]interface{})
	if !ok {
		if _, ok = v.([]interface{}); !ok {
			// return the name of the value passed if not a map[string]interface{} value
			s = append(s, reflect.ValueOf(val).Type().Name())
			return s, m, root, nil
		}
	}

	checkMembers(vv, reflect.ValueOf(val), &s, "")
	return s, m, root, nil
}

// MissingXMLTagsReaderMapRaw consumes the XML data from an io.Reader and returns
// the mxj.Map - map[string]interface{} - representation of the XML data and the raw XML data
// that was read from the io.Reader in addition to the missing XML tags.
// (See github.com/clbanning/mxj documentation of mxj.Map type.)
func MissingXMLTagsReaderMapRaw(r io.Reader, val interface{}) ([]string, mxj.Map, string, []byte, error) {
	var s []string

	m, raw, err := mxj.NewMapXmlReaderRaw(r, mxjCast)
	if err != nil {
		return nil, m, "", raw, err
	}
	// strip off the root value
	var root string
	var v interface{}
	for root, v = range m {
		break
	}

	vv, ok := v.(map[string]interface{})
	if !ok {
		if _, ok = v.([]interface{}); !ok {
			// return the name of the value passed if not a map[string]interface{} value
			s = append(s, reflect.ValueOf(val).Type().Name())
			return s, m, root, raw, nil
		}
	}

	checkMembers(vv, reflect.ValueOf(val), &s, "")
	return s, m, root, raw, nil
}

// ================== where the work is done ...

// cmem is the parent struct member for nested structs
func checkMembers(mv interface{}, val reflect.Value, s *[]string, cmem string) {
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
		// Slice may be nil, so create a Value of it's type
		// 'mv' must be of type []interface{}. If it isn't coerce it.
		sval := reflect.New(tval)
		var slice []interface{}
		var ok bool
		slice, ok = mv.([]interface{})
		if !ok {
			slice = []interface{}{mv}
		}
		// 2.1. Check members of XML list array.
		//      This forces all of them to be regular and w/o typos in key labels.
		for _, sl := range slice {
			checkMembers(sl, sval, s, cmem)
		}
		return // done with reflect.Slice value
	}

	// 3a. Ignore anything that's not a struct.
	if typ.Kind() != reflect.Struct {
		return // just ignore it - don't look for k:v pairs
	}
	// 3b. map value must represent k:v pairs
	mm, ok := mv.(map[string]interface{})
	if !ok {
		*s = append(*s, cmem+typ.Name())
		return
	}
	// 3c. Coerce keys to lower case.
	mkeys := make(map[string]interface{}, len(mm))
	for k, v := range mm {
		mkeys[k] = v
	}

	// 4. Build the list of struct field name:value
	//    We make every key (field) label look like an exported label - "Fieldname".
	//    If there is a XML tag it is used instead of the field label, and saved to
	//    insure that the spec'd tag matches the XML tag exactly.
	type fieldSpec struct {
		name      string
		val       reflect.Value
		tag       []string
		omitempty bool
	}
	fieldCnt := val.NumField()
	var fields []*fieldSpec // use a list so members are in sequence
	var attr bool
	var tagvals string
	var tags []string
	var tag []string
	var oempty bool
	for i := 0; i < fieldCnt; i++ {
		if len(typ.Field(i).PkgPath) > 0 {
			continue // field is NOT exported
		}
		// Ignore xml.Name type fields - they don't appear in the map mm.
		// The root label is handed in as "key" in the initial call.
		if typ.Field(i).Type.Name() == "Name" && typ.Field(i).Type.PkgPath() == "encoding/xml" {
			continue
		}
		tagvals = typ.Field(i).Tag.Get("xml")
		tags = strings.Split(tagvals, ",")
		tag = strings.Split(tags[0], ">")
		// Fields with "-" may or maynot be in the the XML data.
		// don't even bother to check that the Field occurs.
		if tag[0] == "-" {
			continue
		}
		// Scan rest of tags for "omitempty" and "attr".
		// If omitempty occurs we will allow it to occur or not
		// unless the omitemptyOK flag is false, then we strictly
		// check that the kwy is there.
		oempty, attr = false, false
		for _, v := range tags[1:] {
			if v == "omitempty" {
				oempty = true
			}
			if v == "attr" {
				attr = true
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
				fields = append(fields, &fieldSpec{typ.Field(i).Name, val.Field(i), tag, oempty})
			} else {
				fields = append(fields, &fieldSpec{typ.Field(i).Name, val.Field(i), tag, oempty})
			}
		case true:
			if tag[0] == "" {
				fields = append(fields, &fieldSpec{"-" + typ.Field(i).Name, val.Field(i), tag, oempty})
			} else {
				tag[0] = "-" + tag[0]
				fields = append(fields, &fieldSpec{"-" + typ.Field(i).Name, val.Field(i), tag, oempty})
			}
		}
	}

	// 5. check that field names/tags have corresponding map key
	// var ok bool
	var v interface{}
	// var err error
	cmemdepth := 1
	if len(cmem) > 0 {
		cmemdepth = len(strings.Split(cmem, ".")) + 1 // struct hierarchy
	}
	var fn string
	for _, field := range fields {
		// see if we should use XML tag to lookup map key
		if len(field.tag[0]) > 0 {
			fn = field.tag[0]
		} else {
			fn = field.name
		}
		for _, sm := range skipmembers {
			// skip any skipmembers values that aren't at same depth
			if cmemdepth != sm.depth {
				continue
			}
			if len(cmem) > 0 {
				if cmem+"."+fn == sm.val {
					goto next
				}
			} else if fn == sm.val {
				goto next
			}
		}
		v, ok = mkeys[fn]
		// If map key is missing, then record it
		// if there's no omitempty tag or we're ignoring  omitempty tag.
		if !ok && (!field.omitempty || !omitemptyOK) {
			if len(cmem) > 0 {
				// *s = append(*s, cmem+"."+field.name)
				*s = append(*s, cmem+"."+fn)
			} else {
				// *s = append(*s, field.name)
				*s = append(*s, fn)
			}
		}
		if len(cmem) > 0 {
			checkMembers(v, field.val, s, cmem+"."+fn)
		} else {
			checkMembers(v, field.val, s, fn)
		}
	next:
	}
}
