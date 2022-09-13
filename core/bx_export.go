package core

import (
	"archive/zip"
	"blockexchange/types"
	"bytes"
	"encoding/json"
	"io"
	"strconv"
	"time"
)

func ExportBXSchema(w io.Writer, schema *types.Schema, mods []types.SchemaMod, it types.SchemaPartIterator) error {

	archive := zip.NewWriter(w)
	defer archive.Close()

	schema_data, err := json.Marshal(schema)
	if err != nil {
		return err
	}

	err = addDataToZip(archive, "schema.json", schema_data)
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

	err = addDataToZip(archive, "mods.json", mods_data)
	if err != nil {
		return err
	}

	for {
		schemapart, err := it()
		if err != nil {
			return err
		}
		if schemapart == nil {
			// done
			break
		}

		schemapart_data, err := json.Marshal(schemapart)
		if err != nil {
			return err
		}

		partfilename := "schemapart_" + strconv.Itoa(schemapart.OffsetX) + "_" +
			strconv.Itoa(schemapart.OffsetY) + "_" +
			strconv.Itoa(schemapart.OffsetZ) + ".json"

		err = addDataToZip(archive, partfilename, schemapart_data)
		if err != nil {
			return err
		}
	}

	return nil
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
