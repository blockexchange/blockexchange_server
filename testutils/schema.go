package testutils

import (
	"blockexchange/db"
	"blockexchange/types"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func CreateSchema(repo *db.SchemaRepository, t *testing.T, user *types.User, schema *types.Schema) *types.Schema {
	if schema == nil {
		schema = &types.Schema{
			Name: CreateName(10),
		}
	}

	schema.UserID = *user.ID
	assert.NoError(t, repo.CreateSchema(schema))
	return schema
}

var partData = "eJzt1zEOwjAQRFEkDkALPUfJITgNytGpUiAldohi2cu+P6Wryfc69nyd2+bZOPdKHpW8k4f/3Pm2tXy1ktHyKv+xsu5/23F9f8Ty/+pugH/zP5p/53+WrM33L375j50zXUf1n/kOsG1uSuM/c9ra53/07J3/6eCO4H/s7PV49Dzgf+ycc8rzHzX85w7/y/sv5yuw/mefqne/0noE/4lz6cytM/rrr7/++uuvv/7666+//vrrrz8AAAAAAPgHPj+2E4Q"
var partMetadata = "eJxljlEKwjAQBe+y30Vs1FL2MmExsS6kSUm2qC25u2mLUPXvMQ+GmaG3QoaEAGcQ7m1MgH50Llfgg7G6p2Fg3y03cQSsVXM4VmDsjUYnaDgKoDr/Mv1guesuUio+ddnfK9Q1YNvucSJvAE9f9iTB2wL/2Ka/BnKATdGU2sSTXSqfpXGVvD5j2kbObzE5TJw"

func CreateSchemaPart(repo *db.SchemaPartRepository, t *testing.T, schema *types.Schema, schemapart *types.SchemaPart) *types.SchemaPart {
	if schemapart == nil {
		schemapart = &types.SchemaPart{
			SchemaID: *schema.ID,
			OffsetX:  0,
			OffsetY:  0,
			OffsetZ:  0,
			Data:     []byte(partData),
			MetaData: []byte(partMetadata),
			Mtime:    time.Now().Unix() * 1000,
		}
	}

	assert.NoError(t, repo.CreateOrUpdateSchemaPart(schemapart))
	return schemapart
}

func CreateSchemaScreenshot(repo db.SchemaScreenshotRepository, t *testing.T, schema *types.Schema, screenshot *types.SchemaScreenshot) *types.SchemaScreenshot {
	if screenshot == nil {
		screenshot = &types.SchemaScreenshot{}
	}

	screenshot.SchemaID = *schema.ID
	assert.NoError(t, repo.Create(screenshot))
	return screenshot
}
