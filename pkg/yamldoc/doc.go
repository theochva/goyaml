/*
	Package yamldoc provides a simple wrapper around a YAML document's contents.

	This package is intended to be used when you need to quickly read/write
	a handful values from some YAML content contained in a string or a byte array.

	You basically avoid having to go through the process of "modeling" the YAML
	document with a struct or having to dig through a map map[interface{}]interface{}
	to read/write a handful of values.
*/
package yamldoc
