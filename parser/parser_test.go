package parser

import (
	"reflect"
	"testing"

	"github.com/EVFUBS/AlphaLang/ast"
)

func TestParser_ParseStringLiteral(t *testing.T) {
	tests := []struct {
		name string
		p    *Parser
		want *ast.StringLiteral
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.ParseStringLiteral(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parser.ParseStringLiteral() = %v, want %v", got, tt.want)
			}
		})
	}
}
