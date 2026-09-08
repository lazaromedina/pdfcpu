package api

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

// Validación pero generando un PDF diminuto en memoria con /Annots null.
func TestSplitRawWithNullAnnots(t *testing.T) {
	outDir := t.TempDir()

	pdfBytes := buildNullAnnotsPDF()

	ctx, err := ReadContext(bytes.NewReader(pdfBytes), model.NewDefaultConfiguration())
	if err != nil {
		t.Fatalf("ReadContext failed: %v", err)
	}
	if err := ctx.EnsurePageCount(); err != nil {
		t.Fatalf("EnsurePageCount failed: %v", err)
	}
	pageIndRef, err := ctx.PageDictIndRef(1)
	if err != nil {
		t.Fatalf("PageDictIndRef failed: %v", err)
	}
	pageDict, err := ctx.DereferenceDict(*pageIndRef)
	if err != nil {
		t.Fatalf("DereferenceDict failed: %v", err)
	}
	if _, found := pageDict.Find("Annots"); !found {
		t.Fatalf("expected Annots entry in generated PDF")
	}

	if err := Split(bytes.NewReader(pdfBytes), outDir, "mini.pdf", 1, model.NewDefaultConfiguration()); err != nil {
		t.Fatalf("Split returned an error: %v", err)
	}
}

func buildNullAnnotsPDF() []byte {
	var buf bytes.Buffer
	buf.WriteString("%PDF-1.4\n")

	var offsets []int
	writeObj := func(content string) {
		offsets = append(offsets, buf.Len())
		buf.WriteString(content)
	}

	writeObj("1 0 obj\n<< /Type /Catalog /Pages 2 0 R >>\nendobj\n\n")
	writeObj("2 0 obj\n<< /Type /Pages /Count 1 /Kids [3 0 R] >>\nendobj\n\n")
	writeObj("3 0 obj\n<< /Type /Page /Parent 2 0 R /MediaBox [0 0 200 200] /Resources << >> /Contents 4 0 R /Annots 5 0 R >>\nendobj\n\n")
	writeObj("4 0 obj\n<< /Length 0 >>\nstream\n\nendstream\nendobj\n\n")
	writeObj("5 0 obj\nnull\nendobj\n\n")

	xrefPos := buf.Len()
	buf.WriteString("xref\n")
	buf.WriteString(fmt.Sprintf("0 %d\n", len(offsets)+1))
	buf.WriteString("0000000000 65535 f \n")
	for _, off := range offsets {
		buf.WriteString(fmt.Sprintf("%010d 00000 n \n", off))
	}
	buf.WriteString("trailer\n")
	buf.WriteString(fmt.Sprintf("<< /Size %d /Root 1 0 R >>\n", len(offsets)+1))
	buf.WriteString("startxref\n")
	buf.WriteString(fmt.Sprintf("%d\n", xrefPos))
	buf.WriteString("%%EOF\n")

	return buf.Bytes()
}
