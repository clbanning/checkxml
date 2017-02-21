// misc.go - supporting funnctions
// Copyright © 2017 Charles Banning.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package checkxml

import (
	"strings"
)

// List of XML element tags to NOT validate.
var skiptags []string

// SetTagsToIgnore maintains a list of XML element tags in dot-notation
// that should not be validated as exported struct fields.
// NOTE: tags are case sensitive - i.e., "key" != "Key" != "KEY".
//	Dot-notation:
//		- Given XML data:
//			<doc>
//			  <config/>
//			  <data>
//			    <ignore>test</ignore>
//			    <check>this</check>
//			  </data>
//			</doc>
//		- Elements of a XML document are represented as a path from the root:
//		 	"config".
//		- Subelements are represented in a simplier hierarchical manner:
//		 	"data.ignore"
//
func SetTagsToIgnore(s ...string) {
	switch {
	case len(s) == 0 || s[0] == "":
		skiptags = []string{}
	default:
		skiptags = make([]string, len(s))
		copy(skiptags, s)
	}
}

type skipmems struct {
	val   string
	depth int
}

var skipmembers []skipmems

// SetMembersToIgnore creates a list of exported struct member names that should not be checked
// for as tags in the XML-encoded data.  For hierarchical struct members provide the full path for
// the member name using dot-notation. Calling SetMembersToIgnore with no arguments -
// SetMembersToIgnore() - will clear the list.
func SetMembersToIgnore(s ...string) {
	if len(s) == 0 {
		skipmembers = skipmembers[:0]
		return
	}
	skipmembers = make([]skipmems, len(s))
	for i, v := range s {
		skipmembers[i] = skipmems{v, len(strings.Split(v, "."))}
	}
}

// Should we ignore "omitempty" struct tags. By default accept tag.
var omitemptyOK = true

// IgnoreOmitemptyTag determines whether a `xml:",omitempty"` tag is recognized or
// not with respect to the XML data.  By default MissingXMLTags will not include
// in the slice of missing XML tags any struct members that are tagged with "omitempty".
// If the function is toggled or passed the optional argument 'false' then missing
// XML tags may include those XML data tags that correspond to struct members with
// an "omitempty" XML tag.
//
// Calling IgnoreOmitemptyTag with no arguments toggles the handling on/off.  If
// the alternative bool argument is passed, then the argument value determines the
// "omitempty" handling behavior.
func IgnoreOmitemptyTag(ok ...bool) {
	if len(ok) == 0 {
		omitemptyOK = !omitemptyOK
		return
	}
	omitemptyOK = ok[0]
}

// should we try to coerce the values to float64 or bool
var mxjCast bool

// SetMxjCast manages clbanning/mxj decoder flag that causes the mxj.Map values
// to be cast as float64 or bool if possible. The default, SetMxjCast(false), leaves all
// mxj.Map values as string type. Calling SetMxjCast with no arguments - checkxml.SetMxjCast() - 
// will toggle the flag true/false.
// (See github.com/clbanning/mxj documentation of mxj.Map type.)
func SetMxjCast(b ...bool) {
	if len(b) == 0 {
		mxjCast = !mxjCast
	}
	mxjCast = b[0]
}


// HasTags is a convenience function that takes the result slice from MissingTags
// or UnknownTags and returns "true, nil" if the dot-notation 'check' values are
// in the slice.  If one or more of the 'check' values are not in the 'result' slice
// the return value will be "false, []string" where the slice of string values are
// the 'check' values that are not in 'result'.
func HasTags(result []string, check ...string) (bool, []string) {
	r := make(map[string]bool, len(result))
	for _, v := range result {
		r[v] = true
	}
	var missing []string
	for _, v := range check {
		if _, ok := r[v]; ok {
			continue
		}
		missing = append(missing, v)
	}
	if len(missing) > 0 {
		return false, missing
	}
	return true, nil
}
