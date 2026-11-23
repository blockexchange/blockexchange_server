package bx_test

import (
	"archive/zip"
	"blockexchange/schematic/bx"
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createZip(t *testing.T, files map[string]string) []byte {
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	for name, content := range files {
		w, err := zw.Create(name)
		require.NoError(t, err)
		_, err = w.Write([]byte(content))
		require.NoError(t, err)
	}
	require.NoError(t, zw.Close())
	return buf.Bytes()
}

func TestBXImportRejectsInvalidContent(t *testing.T) {
	validSchema := `{"size_x":16,"size_y":16,"size_z":16}`

	cases := map[string][]byte{
		"not a zip":       []byte("hello"),
		"missing schema":  createZip(t, map[string]string{"mods.json": "[]"}),
		"missing mods":    createZip(t, map[string]string{"schema.json": validSchema}),
		"invalid json":    createZip(t, map[string]string{"schema.json": "{", "mods.json": "[]"}),
		"huge size":       createZip(t, map[string]string{"schema.json": `{"size_x":2000000000,"size_y":2000000000,"size_z":2000000000}`, "mods.json": "[]"}),
		"size too large":  createZip(t, map[string]string{"schema.json": `{"size_x":501,"size_y":1,"size_z":1}`, "mods.json": "[]"}),
		"negative size":   createZip(t, map[string]string{"schema.json": `{"size_x":-1,"size_y":1,"size_z":1}`, "mods.json": "[]"}),
		"zero size":       createZip(t, map[string]string{"schema.json": `{"size_x":0,"size_y":1,"size_z":1}`, "mods.json": "[]"}),
		"invalid modname": createZip(t, map[string]string{"schema.json": validSchema, "mods.json": `["../etc"]`}),
		"invalid dependency type": createZip(t, map[string]string{
			"schema.json": validSchema, "mods.json": "[]",
			"dependencies.json": `[{"mod_type":"texturepack","git_url":"https://example.com/x.git"}]`,
		}),
		"dependency file url": createZip(t, map[string]string{
			"schema.json": validSchema, "mods.json": "[]",
			"dependencies.json": `[{"mod_type":"mod","git_url":"file:///etc"}]`,
		}),
		"dependency version option": createZip(t, map[string]string{
			"schema.json": validSchema, "mods.json": "[]",
			"dependencies.json": `[{"mod_type":"mod","git_url":"https://example.com/x.git","version":"--upload-pack=x"}]`,
		}),
		"too many dependencies": createZip(t, map[string]string{
			"schema.json": validSchema, "mods.json": "[]",
			"dependencies.json": "[" + strings.TrimSuffix(strings.Repeat(`{"mod_type":"mod","git_url":"https://example.com/x.git"},`, bx.MaxDependencies+1), ",") + "]",
		}),
		"invalid dependencies json": createZip(t, map[string]string{
			"schema.json": validSchema, "mods.json": "[]", "dependencies.json": `{"mod_type":"mod"}`,
		}),
		"offset mismatch": createZip(t, map[string]string{
			"schema.json":           validSchema,
			"mods.json":             "[]",
			"schemapart_0_0_0.json": `{"offset_x":4000,"offset_y":0,"offset_z":0}`,
		}),
		// highly compressible file exceeding the per-file limit
		"zip bomb": createZip(t, map[string]string{
			"schema.json": `{"size_x":16,"size_y":16,"size_z":16,"description":"` + strings.Repeat("a", bx.MaxFileSize) + `"}`,
			"mods.json":   "[]",
		}),
	}

	for name, data := range cases {
		t.Run(name, func(t *testing.T) {
			start := time.Now()
			_, err := bx.ParseBXContent(data)
			assert.ErrorIs(t, err, bx.ErrInvalidBX)
			assert.Less(t, time.Since(start), 5*time.Second)
		})
	}
}

func TestBXImportMinimal(t *testing.T) {
	data := createZip(t, map[string]string{
		"schema.json": `{"name":"test","size_x":17,"size_y":1,"size_z":1}`,
		"mods.json":   `["default"]`,
	})
	res, err := bx.ParseBXContent(data)
	require.NoError(t, err)
	assert.Equal(t, []string{"default"}, res.Mods)
	assert.Equal(t, 0, len(res.Parts))
}
