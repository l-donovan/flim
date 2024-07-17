package expressions

import (
	"github.com/l-donovan/flim/common"
	"fmt"
	"strings"
)

type IntegerLiteralExpression struct {
	val int64
}

func NewIntegerLiteralExpression(val int64) (IntegerLiteralExpression, error) {
	return IntegerLiteralExpression{val: val}, nil
}

func (e IntegerLiteralExpression) ToString() string {
	return fmt.Sprintf("IntegerLiteralExpression<%d>", e.val)
}

func (e IntegerLiteralExpression) GetTags() map[string]common.Expression {
	return map[string]common.Expression{}
}

func (e IntegerLiteralExpression) ReplaceReferences(tags map[string]common.Expression) (common.Expression, error) {
	return e, nil
}

func (e IntegerLiteralExpression) Evaluate(handlers map[string]common.HandlerFunc) (interface{}, error) {
	return e.val, nil
}

func (e IntegerLiteralExpression) Serialize(config *common.SerializerConfig, indentLevel int) (string, error) {
	return fmt.Sprintf("%d", e.val), nil
}

type FloatLiteralExpression struct {
	val float64
}

func NewFloatLiteralExpression(val float64) (FloatLiteralExpression, error) {
	return FloatLiteralExpression{val: val}, nil
}

func (e FloatLiteralExpression) ToString() string {
	return fmt.Sprintf("FloatLiteralExpression<%f>", e.val)
}

func (e FloatLiteralExpression) GetTags() map[string]common.Expression {
	return map[string]common.Expression{}
}

func (e FloatLiteralExpression) ReplaceReferences(tags map[string]common.Expression) (common.Expression, error) {
	return e, nil
}

func (e FloatLiteralExpression) Evaluate(handlers map[string]common.HandlerFunc) (interface{}, error) {
	return e.val, nil
}

func (e FloatLiteralExpression) Serialize(config *common.SerializerConfig, indentLevel int) (string, error) {
	return fmt.Sprintf("%f", e.val), nil
}

type BooleanLiteralExpression struct {
	val bool
}

func NewBooleanLiteralExpression(val bool) (BooleanLiteralExpression, error) {
	return BooleanLiteralExpression{val: val}, nil
}

func (e BooleanLiteralExpression) ToString() string {
	return fmt.Sprintf("BooleanLiteralExpression<%t>", e.val)
}

func (e BooleanLiteralExpression) ReplaceReferences(tags map[string]common.Expression) (common.Expression, error) {
	return e, nil
}

func (e BooleanLiteralExpression) GetTags() map[string]common.Expression {
	return map[string]common.Expression{}
}

func (e BooleanLiteralExpression) Evaluate(handlers map[string]common.HandlerFunc) (interface{}, error) {
	return e.val, nil
}

func (e BooleanLiteralExpression) Serialize(config *common.SerializerConfig, indentLevel int) (string, error) {
	return fmt.Sprintf("%t", e.val), nil
}

type StringLiteralExpression struct {
	val string
}

func NewStringLiteralExpression(val string) (StringLiteralExpression, error) {
	return StringLiteralExpression{val: val}, nil
}

func (e StringLiteralExpression) ToString() string {
	return fmt.Sprintf("StringLiteralExpression<%s>", e.val)
}

func (e StringLiteralExpression) ReplaceReferences(tags map[string]common.Expression) (common.Expression, error) {
	return e, nil
}

func (e StringLiteralExpression) GetTags() map[string]common.Expression {
	return map[string]common.Expression{}
}

func (e StringLiteralExpression) Evaluate(handlers map[string]common.HandlerFunc) (interface{}, error) {
	return e.val, nil
}

func (e StringLiteralExpression) Serialize(config *common.SerializerConfig, indentLevel int) (string, error) {
	// TODO Escape more things
	replacer := strings.NewReplacer(
		"\n", "\\n",
	)
	escapedVal := replacer.Replace(e.val)
	return fmt.Sprintf("\"%s\"", escapedVal), nil
}

type NullLiteralExpression struct{}

func NewNullLiteralExpression() (NullLiteralExpression, error) {
	return NullLiteralExpression{}, nil
}

func (e NullLiteralExpression) ToString() string {
	return "NullLiteralExpression<>"
}

func (e NullLiteralExpression) ReplaceReferences(tags map[string]common.Expression) (common.Expression, error) {
	return e, nil
}

func (e NullLiteralExpression) GetTags() map[string]common.Expression {
	return map[string]common.Expression{}
}

func (e NullLiteralExpression) Evaluate(handlers map[string]common.HandlerFunc) (interface{}, error) {
	return nil, nil
}

func (e NullLiteralExpression) Serialize(config *common.SerializerConfig, indentLevel int) (string, error) {
	return "null", nil
}
