package render

import (
	"blockexchange/types"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

type MockSchemaPartRepository struct {
}

func (r *MockSchemaPartRepository) GetBySchemaIDAndOffset(schema_id int64, offset_x, offset_y, offset_z int) (*types.SchemaPart, error) {
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
	cm, err := GetColorMapping()
	assert.NoError(t, err)

	renderer := NewISORenderer(repo.GetBySchemaIDAndOffset, cm)
	schema := types.Schema{
		SizeX: 32,
		SizeY: 32,
		SizeZ: 32,
		ID:    0,
	}
	png, err := renderer.RenderIsometricPreview(&schema)
	assert.NoError(t, err)
	assert.NotNil(t, png)

	file, err := os.Create("/tmp/render.png")
	assert.NoError(t, err)

	file.Write(png)
	fmt.Println(file.Name())
}
