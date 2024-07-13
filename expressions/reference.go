package expressions

import (
	"flim/common"
	"fmt"
)

type ReferenceExpression struct {
	name string
}

func NewReferenceExpression(name string) (ReferenceExpression, error) {
	return ReferenceExpression{name: name}, nil
}

func (e ReferenceExpression) ToString() string {
	return fmt.Sprintf("ReferenceExpression<%s>", e.name)
}

func (e ReferenceExpression) GetTags() map[string]common.Expression {
	return map[string]common.Expression{}
}

func (e ReferenceExpression) ReplaceReferences(tags map[string]common.Expression) (common.Expression, error) {
	if replacement, exists := tags[e.name]; exists {
		return replacement, nil
	} else {
		return nil, fmt.Errorf("could not find tag `%s'", e.name)
	}
}

func (e ReferenceExpression) Evaluate(handlers map[string]common.HandlerFunc) (interface{}, error) {
	return nil, fmt.Errorf("attempted to Evaluate a ReferenceExpression (hint: call ReplaceReferences first)")
}

func (e ReferenceExpression) Serialize(config *common.SerializerConfig, indentLevel int) (string, error) {
	return fmt.Sprintf("&%s", e.name), nil
}
