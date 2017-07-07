package test

import (
	`greentea/prebuilder/builder`
	"testing"
)

func TestBuilder(t *testing.T) {
	builder.ParseFile(`model.go`)
}
