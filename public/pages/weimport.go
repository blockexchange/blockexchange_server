package pages

import (
	"blockexchange/colormapping"
	"blockexchange/controller"
	"blockexchange/core"
	"blockexchange/types"
	"blockexchange/worldedit"
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

type WeImportModel struct {
	ImportError    error
	ImportedSchema *types.Schema
	FileSize       int
}

func handleWeImportPost(rc *controller.RenderContext, m *WeImportModel) {
	r := rc.Request()
	err := r.ParseMultipartForm(1000 * 1000 * 10)
	if err != nil {
		m.ImportError = err
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		m.ImportError = err
		return
	}
	defer file.Close()

	buf := bytes.NewBuffer([]byte{})
	_, err = io.Copy(buf, file)
	if err != nil {
		m.ImportError = err
		return
	}

	entries, modnames, err := worldedit.Parse(buf.Bytes())
	if err != nil {
		m.ImportError = err
		return
	}
	max_x, max_y, max_z := worldedit.GetBoundaries(entries)

	parts, err := worldedit.Import(entries)
	if err != nil {
		m.ImportError = err
		return
	}

	newSchemaName := strings.Split(handler.Filename, ".")[0]
	searchRes, err := rc.Repositories().SchemaSearchRepo.Search(&types.SchemaSearchRequest{
		UserName:   &rc.Claims().Username,
		SchemaName: &newSchemaName,
	}, 1, 0)
	if err != nil {
		m.ImportError = err
		return
	}
	if len(searchRes) > 0 {
		newSchemaName = fmt.Sprintf("%s_%d", newSchemaName, rand.Int())
	}

	schema := &types.Schema{
		Created:     time.Now().Unix() * 1000,
		UserID:      rc.Claims().UserID,
		Name:        newSchemaName,
		Mtime:       time.Now().Unix() * 1000,
		Description: "WE Import",
		Complete:    false,
		SizeX:       max_x + 1,
		SizeY:       max_y + 1,
		SizeZ:       max_z + 1,
		TotalParts:  len(parts),
		License:     "CC0",
	}

	err = rc.Repositories().SchemaRepo.CreateSchema(schema)
	if err != nil {
		m.ImportError = err
		return
	}

	for _, part := range parts {
		sp, err := part.Convert()
		if err != nil {
			m.ImportError = err
			return
		}
		sp.SchemaID = schema.ID

		err = rc.Repositories().SchemaPartRepo.CreateOrUpdateSchemaPart(sp)
		if err != nil {
			m.ImportError = err
			return
		}
	}

	for _, modname := range modnames {
		err = rc.Repositories().SchemaModRepo.CreateSchemaMod(&types.SchemaMod{
			SchemaID: schema.ID,
			ModName:  modname,
		})
		if err != nil {
			m.ImportError = err
			return
		}
	}

	schema.Complete = true
	err = rc.Repositories().SchemaRepo.UpdateSchema(schema)
	if err != nil {
		m.ImportError = err
		return
	}

	err = rc.Repositories().SchemaRepo.CalculateStats(schema.ID)
	if err != nil {
		m.ImportError = err
		return
	}

	cm := colormapping.NewColorMapping()
	err = cm.LoadDefaults()
	if err != nil {
		m.ImportError = err
		return
	}

	_, err = core.UpdatePreview(schema, rc.Repositories(), cm)
	if err != nil {
		m.ImportError = err
		return
	}

	m.FileSize = int(handler.Size)
	m.ImportedSchema, m.ImportError = rc.Repositories().SchemaRepo.GetSchemaById(schema.ID)
}

func WeImport(rc *controller.RenderContext) error {
	m := &WeImportModel{}

	if rc.Request().Method == http.MethodPost {
		handleWeImportPost(rc, m)
	}

	return rc.Render("pages/weimport.html", m)
}
