package dtos

type StructValidator interface {
	ValidateStruct(interface{}) error
	Engine() interface{}
}
