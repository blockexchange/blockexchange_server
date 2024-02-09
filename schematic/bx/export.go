package bx

import (
	"archive/zip"
	"blockexchange/types"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"time"
)

type Exporter struct {
	w          *zip.Writer
	intialized bool
}

func NewExporter(w io.Writer) *Exporter {
	archive := zip.NewWriter(w)

	return &Exporter{w: archive, intialized: false}
}

func (e *Exporter) ExportMetadata(schema *types.Schema, mods []*types.SchemaMod) error {
	schema_data, err := json.Marshal(schema)
	if err != nil {
		return err
	}

	err = addDataToZip(e.w, "schema.json", schema_data)
	if err != nil {
		return err
	}

	modlist := []string{}
	for _, mod := range mods {
		modlist = append(modlist, mod.ModName)
	}

	mods_data, err := json.Marshal(modlist)
	if err != nil {
		return err
	}

	err = addDataToZip(e.w, "mods.json", mods_data)
	if err != nil {
		return err
	}

	return nil
}

func (e *Exporter) Export(schemapart *types.SchemaPart) error {
	schemapart_data, err := json.Marshal(schemapart)
	if err != nil {
		return err
	}

	return addDataToZip(e.w, formatSchemapartFilename(schemapart), schemapart_data)
}

func (e *Exporter) Close() error {
	return e.w.Close()
}

func formatSchemapartFilename(schemapart *types.SchemaPart) string {
	return fmt.Sprintf("schemapart_%d_%d_%d.json", schemapart.OffsetX, schemapart.OffsetY, schemapart.OffsetZ)
}

func addDataToZip(archive *zip.Writer, filename string, data []byte) error {
	header := zip.FileHeader{
		Name:               filename,
		Modified:           time.Now(),
		UncompressedSize64: uint64(len(data)),
		Method:             zip.Deflate,
	}

	writer, err := archive.CreateHeader(&header)
	if err != nil {
		return err
	}

	_, err = io.Copy(writer, bytes.NewReader(data))
	return err
}
