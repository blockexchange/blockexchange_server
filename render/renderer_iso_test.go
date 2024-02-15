package render

import (
	"blockexchange/types"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/minetest-go/colormapping"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

type MockSchemaPartRepository struct {
}

func (r *MockSchemaPartRepository) GetBySchemaUIDAndOffset(schema_uid string, offset_x, offset_y, offset_z int) (*types.SchemaPart, error) {
	logrus.WithFields(logrus.Fields{
		"offset_x": offset_x,
		"offset_y": offset_y,
		"offset_z": offset_z,
	}).Trace("MockSchemaPartRepository::GetBySchemaIDAndOffset")

	f, err := os.Open(fmt.Sprintf("testdata/%d_%d_%d.json", offset_x, offset_y, offset_z))
	if err != nil || f == nil {
		return nil, nil
	}
	part := types.SchemaPart{}
	err = json.NewDecoder(f).Decode(&part)
	return &part, err
}

func TestISORenderer(t *testing.T) {
	logrus.SetLevel(logrus.TraceLevel)

	repo := MockSchemaPartRepository{}
	cm := colormapping.NewColorMapping()
	assert.NoError(t, cm.LoadDefaults())

	renderer := NewISORenderer(repo.GetBySchemaUIDAndOffset, cm)
	schema := types.Schema{
		SizeX: 32,
		SizeY: 32,
		SizeZ: 32,
		UID:   uuid.NewString(),
	}
	png, err := renderer.RenderIsometricPreview(&schema)
	assert.NoError(t, err)
	assert.NotNil(t, png)

	file, err := os.Create("/tmp/render.png")
	assert.NoError(t, err)

	file.Write(png)
	fmt.Println(file.Name())
}
