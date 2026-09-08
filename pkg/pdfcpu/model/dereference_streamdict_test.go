package model

import (
	"testing"

	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
)

func TestDereferenceDictAcceptsStreamDict(t *testing.T) {
	xRefTable := newXRefTable(NewDefaultConfiguration())

	dict := types.Dict{"Type": types.Name("Example")}
	got, err := xRefTable.DereferenceDict(types.StreamDict{Dict: dict})
	if err != nil {
		t.Fatalf("DereferenceDict returned an error: %v", err)
	}

	if got.NameEntry("Type") == nil || *got.NameEntry("Type") != "Example" {
		t.Fatalf("unexpected dict: %v", got)
	}
}
