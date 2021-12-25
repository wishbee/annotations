package sample

// Example: You can use wishgen with go generate as shown below.
//go:generate wishgen -dir=.

/* SomeStruct type encapsulate some structure.
 * Below annotation will generate the 'setXXX' methods for all the fields in this struct.
 * @wish:setters
 * Below annotation will generate the 'getXXX' methods for all the fields in this struct.
 * @wish:getters
 */
type SomeStruct struct {
	age int
	name string
}

/*
 * AnotherStruct type encapsulate some structure.
 */
type AnotherStruct struct {
	ID int
	/*
	 * Below annotation will generate the 'setName' method for the 'name' field for this this struct.
	@wish:set
	 */
	name string
}

// SomInterface is an Empty interface
type SomInterface interface {
	MethodA()
}
