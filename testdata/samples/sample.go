package sample

// Example: You can use wishgen with go generate as shown below.
//go:generate wishgen -dir=.

/* SomeStruct type encapsulate some structure.
 * Below annotation will generate the 'setXXX' methods for all the fields in this struct.
 * @wish:setter
 * Below annotation will generate the 'XXX' methods for all the fields in this struct.
 * These methods will simply return the value of the field.
 * @wish:getter
 */
type SomeStruct struct {
	age  int
	name string
}

/*
 * AnotherStruct type encapsulate some structure.
 */
type AnotherStruct struct {
	ID int
	/*
		 * Below annotation will generate the 'setName' method for the 'name' field for this this struct.
					@wish:setter
	*/
	name string
}

// SomInterface is an Empty interface
type SomInterface interface {
	MethodA()
}

/*
	@wish:fluent
*/
type SomeTypeForFluent struct {
	age  int
	name string
}

// Once above wish is fulfilled then below method should compile without error.
func fluentTypeUse() {
	s := &SomeTypeForFluent{}
	s.SetName("Bill").SetAge(35)
}
