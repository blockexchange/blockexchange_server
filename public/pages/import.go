package pages

import (
	"blockexchange/colormapping"
	"blockexchange/controller"
	"blockexchange/core"
	"blockexchange/types"
	"blockexchange/worldedit"
	"bytes"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"mime/multipart"
	"net/http"
	"strings"
	"time"
)

type ImportModel struct {
	ImportError    error
	ImportedSchema *types.Schema
	FileSize       int
}

func handleImport(rc *controller.RenderContext, m *ImportModel) {
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

	t := r.FormValue("type")
	var schema *types.Schema
	switch t {
	case "worldedit":
		schema, err = handleWeImport(rc, handler, buf)
	case "bx":
		schema, err = handleBXImport(rc, handler, buf)
	default:
		err = errors.New("unrecognized format: " + t)
	}

	if err != nil {
		m.ImportError = err
		return
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
	m.ImportedSchema, err = rc.Repositories().SchemaRepo.GetSchemaById(schema.ID)
	if err != nil {
		m.ImportError = err
		return
	}
}

func handleBXImport(rc *controller.RenderContext, handler *multipart.FileHeader, buf *bytes.Buffer) (*types.Schema, error) {
	res, err := core.ImportBXSchema(bytes.NewReader(buf.Bytes()), handler.Size)
	if err != nil {
		return nil, err
	}

	newSchemaName := res.Schema.Name
	searchRes, err := rc.Repositories().SchemaSearchRepo.Search(&types.SchemaSearchRequest{
		UserName:   &rc.Claims().Username,
		SchemaName: &newSchemaName,
	}, 1, 0)
	if err != nil {
		return nil, err
	}
	if newSchemaName == "" || !core.ValidateName(newSchemaName) {
		// invalid or no name
		newSchemaName = "import"
	}
	if len(searchRes) > 0 {
		newSchemaName = fmt.Sprintf("%s_%d", newSchemaName, rand.Int())
	}

	res.Schema.Created = time.Now().Unix() * 1000
	res.Schema.Mtime = time.Now().Unix() * 1000
	res.Schema.UserID = rc.Claims().UserID
	res.Schema.Name = newSchemaName

	repos := rc.Repositories()

	err = repos.SchemaRepo.CreateSchema(res.Schema)
	if err != nil {
		return nil, err
	}

	for _, modname := range res.Mods {
		err = repos.SchemaModRepo.CreateSchemaMod(&types.SchemaMod{
			SchemaID: res.Schema.ID,
			ModName:  modname,
		})
		if err != nil {
			return nil, err
		}
	}

	for _, part := range res.Parts {
		part.SchemaID = res.Schema.ID
		part.Mtime = res.Schema.Mtime
		err = repos.SchemaPartRepo.CreateOrUpdateSchemaPart(part)
		if err != nil {
			return nil, err
		}
	}

	return res.Schema, nil
}

func handleWeImport(rc *controller.RenderContext, handler *multipart.FileHeader, buf *bytes.Buffer) (*types.Schema, error) {
	entries, modnames, err := worldedit.Parse(buf.Bytes())
	if err != nil {
		return nil, err
	}
	max_x, max_y, max_z := worldedit.GetBoundaries(entries)

	parts, err := worldedit.Import(entries)
	if err != nil {
		return nil, err
	}

	newSchemaName := strings.Split(handler.Filename, ".")[0]
	searchRes, err := rc.Repositories().SchemaSearchRepo.Search(&types.SchemaSearchRequest{
		UserName:   &rc.Claims().Username,
		SchemaName: &newSchemaName,
	}, 1, 0)
	if err != nil {
		return nil, err
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
		return nil, err
	}

	for _, part := range parts {
		sp, err := part.Convert()
		if err != nil {
			return nil, err
		}
		sp.SchemaID = schema.ID

		err = rc.Repositories().SchemaPartRepo.CreateOrUpdateSchemaPart(sp)
		if err != nil {
			return nil, err
		}
	}

	for _, modname := range modnames {
		err = rc.Repositories().SchemaModRepo.CreateSchemaMod(&types.SchemaMod{
			SchemaID: schema.ID,
			ModName:  modname,
		})
		if err != nil {
			return nil, err
		}
	}

	return rc.Repositories().SchemaRepo.GetSchemaById(schema.ID)
}

func SchemaImport(rc *controller.RenderContext) error {
	m := &ImportModel{}

	if rc.Request().Method == http.MethodPost {
		handleImport(rc, m)
	}

	return rc.Render("pages/import.html", m)
}
