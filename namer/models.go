package namer

type INamer interface {
	Initialize()
	GetName(value interface{}) string
}
