# annotations
Golang annotations by Wish

Use annotations in golang comments for struct or field to generate getter/setter for the fields in your structs.

Format of specifying annotation: ```@wish:yyy``` where yyy is the supported annotation. 

Currently following annotations are supported.
- get
- set

Above annotations are supported at field level.

- getter
- setter

Above annotations are supported at struct level.

Example usage:
```
/* 
 * SomeStruct type encapsulate some structure.
 * Below annotation will generate the 'setXXX' methods for all the fields in this struct.
 * @wish:setters
 * Below annotation will generate the 'getXXX' methods for all the fields in this struct.
 * @wish:getters
 */
type SomeStruct struct {
	age int
	name string
}
```