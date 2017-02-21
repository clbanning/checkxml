
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
	result, root, _ := MissingXMLTags([]byte(data), e)
	// result: [elem2.another elem3]
	// root: doc

	result, root, _ = UnknownXMLTags([]byte(data), e)
	// result: [elem2.notes elem4]
	// root: doc

NOTE: this package is dependent upon github.com/clbanning/mxj.

<h4>USAGE</h4>

https://godoc.org/github.com/clbanning/checkxml


<h4>MOTIVATION</h4>

I was having a conversation and the topic of validating streams of XML
data from various sources came up. What the developer was looking for 
was a "signature" that could be used to route the data for reformating
to meet the production application requirements.  Having already done
something similar for JSON with the github.com/clbanning/checkjson package
and having the github.com/clbanning/mxj package available, it was a simple
exercise to do a mashup of the two packages.
