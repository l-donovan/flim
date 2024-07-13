package common

type HandlerFunc func(data interface{}) (interface{}, error)

type Expression interface {
	ToString() string
	GetTags() map[string]Expression
	ReplaceReferences(map[string]Expression) (Expression, error)
	Evaluate(map[string]HandlerFunc) (interface{}, error)
	Serialize(config *SerializerConfig, indentLevel int) (string, error)
}
