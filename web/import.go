package web

import (
	"blockexchange/core"
	"blockexchange/db"
	"blockexchange/types"
	"blockexchange/worldedit"
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io"
	"math/rand"
	"mime/multipart"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/csrf"
)

type ImportModel struct {
	ImportError    error
	ImportedSchema *types.Schema
	FileSize       int
	CSRFField      template.HTML
}

func (ctx *Context) SchemaImport(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	m := &ImportModel{
		CSRFField: csrf.TemplateField(r),
	}

	if r.Method == http.MethodPost {
		handleImport(ctx.Repos, r, c, m)
	}

	ctx.tu.ExecuteTemplate(w, r, "import.html", m)
}

func handleImport(repos *db.Repositories, r *http.Request, c *types.Claims, m *ImportModel) {
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

	var schema *types.Schema
	if strings.HasSuffix(handler.Filename, ".we") {
		schema, err = handleWeImport(repos, c, handler, buf)
	} else if strings.HasSuffix(handler.Filename, ".zip") {
		schema, err = handleBXImport(repos, c, handler, buf)
	} else {
		err = errors.New("unrecognized file extension")
	}

	if err != nil {
		m.ImportError = err
		return
	}

	schema.Complete = true
	err = repos.SchemaRepo.UpdateSchema(schema)
	if err != nil {
		m.ImportError = err
		return
	}

	err = repos.SchemaRepo.CalculateStats(schema.ID)
	if err != nil {
		m.ImportError = err
		return
	}

	_, err = core.UpdatePreview(schema, repos)
	if err != nil {
		m.ImportError = err
		return
	}

	m.FileSize = int(handler.Size)
	m.ImportedSchema, err = repos.SchemaRepo.GetSchemaById(schema.ID)
	if err != nil {
		m.ImportError = err
		return
	}
}

func handleBXImport(repos *db.Repositories, c *types.Claims, handler *multipart.FileHeader, buf *bytes.Buffer) (*types.Schema, error) {
	res, err := core.ImportBXSchema(bytes.NewReader(buf.Bytes()), handler.Size)
	if err != nil {
		return nil, err
	}

	newSchemaName := res.Schema.Name
	searchRes, err := repos.SchemaSearchRepo.Search(&types.SchemaSearchRequest{
		UserName:   &c.Username,
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
	res.Schema.UserID = c.UserID
	res.Schema.Name = newSchemaName

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

func handleWeImport(repos *db.Repositories, c *types.Claims, handler *multipart.FileHeader, buf *bytes.Buffer) (*types.Schema, error) {
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
	searchRes, err := repos.SchemaSearchRepo.Search(&types.SchemaSearchRequest{
		UserName:   &c.Username,
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
		UserID:      c.UserID,
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

	err = repos.SchemaRepo.CreateSchema(schema)
	if err != nil {
		return nil, err
	}

	for _, part := range parts {
		sp, err := part.Convert()
		if err != nil {
			return nil, err
		}
		sp.SchemaID = schema.ID

		err = repos.SchemaPartRepo.CreateOrUpdateSchemaPart(sp)
		if err != nil {
			return nil, err
		}
	}

	for _, modname := range modnames {
		err = repos.SchemaModRepo.CreateSchemaMod(&types.SchemaMod{
			SchemaID: schema.ID,
			ModName:  modname,
		})
		if err != nil {
			return nil, err
		}
	}

	return repos.SchemaRepo.GetSchemaById(schema.ID)
}
