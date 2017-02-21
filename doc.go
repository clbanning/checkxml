// misc.go - supporting funnctions
// Copyright © 2017 Charles Banning.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package checkxml provides functions for validating XML data against a struct definition.

The MissingXMLTags functions decode XML data and return a slice of struct fields that will
not be set using the encoding/xml Unmarshal function. The UnknownXMLTags functions decode
XML data and return a slice of XML elements that will not be decoded using the encoding/xml
Unmarshal function for the specified struct definition.

	Example:
	
	data := `<doc>
	           <elem1>a simple element</elem1>
	           <elem2>
	             <subelem>something more complex</subelem>
	             <notes>take a look at this</notes>
	           </elem2>
	           <elem4>extraneous</elem4>
	         </doc>`

	type sub struct {
		Subelem string `xml:"subelem,omitempty"`
		Another string `xml:"another"`
	}
	type elem struct {
		Elem1 string `xml:"elem1"`
		Elem2 sub    `xml:"elem2"`
		Elem3 bool   `xml:"elem3"`
	}

	e := new(elem)
	result, _, _ := MissingXMLTags([]byte(data), e)
	// result: [elem3 elem2.another]

	result, _, _ = UnknownTags([]byte(data), e)
	// result: [elem4 elem2.notes]

NOTE: this package is dependent upon github.com/clbanning/mxj.
*/
package checkxml
