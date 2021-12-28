# annotations
Golang annotations by Wish

Express annotations in golang comments as wishes to generate code for fields/struct. 

Format of specifying annotation: ```@wish:yyy:payload``` 

where:
 - ```yyy``` is the supported 'wish'. Wish implements the code generation logic, depending on the type of wish user wants to express in the comment. 
 - ```payload``` is the extra information string that can be passed in the comment depending on wish which supports them.

Currently following annotations are supported.
- getter (To generate the getters for the fields of the struct. The getters does not have 'Get' prefix in the getter name.)
- setter (To generate the setters for the fields of the struct)
- fluent (To generate the setters for the fields of the struct. Each setter returns the receiver as return value so that method chaining can be performed)

Above annotations are supported at struct level.

Example usage for wishes ```@wish:setter``` and ```@wish:getter```:
```
/* 
 * SomeStruct type encapsulate some structure.
 * Below annotation will generate the 'setXXX' methods for all the fields in this struct.
 * @wish:setter
 * Below annotation will generate the 'getXXX' methods for all the fields in this struct.
 * @wish:getter
 */
type SomeStruct struct {
	age int
	name string
}
```

The wishes expressed above ```@wish:setter``` and ```@wish:getter``` would generate below code when processed by ```wishgen```

```
func (s *SomeStruct)SetAge(v int) {
	s.age = v
}
func (s *SomeStruct)SetName(v string) {
	s.name = v
}
func (s *SomeStruct)Age() int {
	return s.age
}
func (s *SomeStruct)Name() string {
	return s.name
}
```

Example usage for wish ```@wish:fluent```:
```
/*
	@wish:fluent
 */
type SomeTypeForFluent struct {
	age int
	name string
}
```

The wish expressed above ```@wish:fluent```  would generate below code when processed by ```wishgen``` 
```
func (s *SomeTypeForFluent)SetAge(v int) *SomeTypeForFluent {
	s.age = v
	return s
}
func (s *SomeTypeForFluent)SetName(v string) *SomeTypeForFluent {
	s.name = v
	return s
}
```
