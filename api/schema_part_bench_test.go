package api_test

import (
	"blockexchange/api"
	"blockexchange/core"
	"blockexchange/db"
	"blockexchange/parser"
	"blockexchange/testutils"
	"blockexchange/types"
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
)

// createBenchPart creates a cubic schemapart with the given edge length,
// either filled with one node ("uniform") or random nodes, params and metadata ("random")
func createBenchPart(b *testing.B, edge int, random bool) *types.SchemaPart {
	return createBenchPartSeed(b, edge, random, 1)
}

// createBenchPartSeed is createBenchPart with a configurable random seed
func createBenchPartSeed(b *testing.B, edge int, random bool, seed int64) *types.SchemaPart {
	rnd := rand.New(rand.NewSource(seed))
	volume := edge * edge * edge

	psp := &parser.ParsedSchemaPart{
		NodeIDS:        make([]int16, volume),
		Param1:         make([]byte, volume),
		Param2:         make([]byte, volume),
		NodeNameLookup: map[int16]string{},
		Meta: &parser.SchemaPartMetadata{
			NodeMapping: map[string]int16{},
			Size:        parser.SchemaPartSize{X: edge, Y: edge, Z: edge},
			Metadata: &parser.Metadata{
				Meta:   map[string]*parser.MetadataEntry{},
				Timers: map[string]*parser.MetadataTimer{},
			},
		},
	}

	nodecount := 1
	if random {
		nodecount = 64
	}
	for i := 0; i < nodecount; i++ {
		name := fmt.Sprintf("mymod:node_%d", i)
		psp.NodeNameLookup[int16(i)] = name
		psp.Meta.NodeMapping[name] = int16(i)
	}

	if random {
		for i := 0; i < volume; i++ {
			psp.NodeIDS[i] = int16(rnd.Intn(nodecount))
			psp.Param1[i] = byte(rnd.Intn(256))
			psp.Param2[i] = byte(rnd.Intn(256))
			if rnd.Intn(50) == 0 {
				x, y, z := psp.GetPos(i)
				psp.Meta.Metadata.Meta[psp.Meta.Metadata.GetKey(x, y, z)] = &parser.MetadataEntry{
					Fields:      map[string]string{"infotext": fmt.Sprintf("node %d", rnd.Int())},
					Inventories: map[string][]string{"main": {"default:dirt 99", "", "default:stone 12"}},
				}
			}
		}
	}

	sp, err := psp.Convert()
	if err != nil {
		b.Fatal(err)
	}
	return sp
}

func BenchmarkCreateSchemaPart(b *testing.B) {
	if !testutils.IsDatabaseAvailable() {
		b.SkipNow()
	}
	repos, err := db.Init()
	if err != nil {
		b.Fatal(err)
	}
	cfg := types.CreateConfig()
	a, _, err := api.NewApi(repos, cfg)
	if err != nil {
		b.Fatal(err)
	}

	user := &types.User{
		UID:  uuid.NewString(),
		Name: core.CreateToken(10),
		Type: types.UserTypeLocal,
		Role: types.UserRoleDefault,
	}
	if err = repos.UserRepo.CreateUser(user); err != nil {
		b.Fatal(err)
	}
	token, err := core.New(cfg, nil).CreateJWT(user, []types.JWTPermission{types.JWTPermissionUpload}, time.Hour)
	if err != nil {
		b.Fatal(err)
	}

	for _, edge := range []int{1, 8, 16} {
		for _, content := range []string{"uniform", "random"} {
			for _, mode := range []string{"insert", "update"} {
				name := fmt.Sprintf("size=%dx%dx%d/content=%s/mode=%s", edge, edge, edge, content, mode)
				b.Run(name, func(b *testing.B) {
					schema := &types.Schema{
						UID:     uuid.NewString(),
						UserUID: user.UID,
						Name:    core.CreateToken(10),
						SizeX:   2000,
						SizeY:   2000,
						SizeZ:   2000,
					}
					if err := repos.SchemaRepo.CreateSchema(schema); err != nil {
						b.Fatal(err)
					}

					sp := createBenchPart(b, edge, content == "random")
					sp.SchemaUID = schema.UID

					// pre-encode the bodies, insert mode uses a new offset for every request
					bodies := make([][]byte, b.N)
					for i := range bodies {
						if mode == "insert" {
							sp.OffsetX = (i % 125) * 16
							sp.OffsetY = ((i / 125) % 125) * 16
							sp.OffsetZ = (i / (125 * 125)) * 16
						}
						bodies[i], err = json.Marshal(sp)
						if err != nil {
							b.Fatal(err)
						}
					}

					handler := a.Secure(a.CreateSchemaPart)
					b.ResetTimer()

					for i := 0; i < b.N; i++ {
						r := httptest.NewRequest("POST", "http://", bytes.NewReader(bodies[i]))
						r.Header.Set("Authorization", token)
						w := httptest.NewRecorder()
						handler(w, r)
						if w.Code != 200 {
							b.Fatalf("unexpected status %d: %s", w.Code, w.Body.String())
						}
					}
					b.ReportMetric(float64(len(bodies[0])), "req-bytes")
				})
			}
		}
	}
}

type patchBenchCase struct {
	edge     int
	content  string  // "uniform" or "random", used for the stored part and the patch
	filter   string  // name of the filter applied to the selected nodes
	coverage float64 // fraction of nodes selected by the patchfilter
}

var patchBenchFilters = map[string]types.PatchFilterType{
	"noop":     types.PatchNoOp,
	"param2":   types.PatchReplaceParam2,
	"nodename": types.PatchReplaceNodename,
	"all":      types.PatchReplaceAll,
}

func (c patchBenchCase) String() string {
	return fmt.Sprintf("size=%dx%dx%d/content=%s/filter=%s/coverage=%g%%", c.edge, c.edge, c.edge, c.content, c.filter, c.coverage*100)
}

func patchBenchCases() []patchBenchCase {
	cases := []patchBenchCase{}
	for _, edge := range []int{1, 8, 16} {
		for _, content := range []string{"uniform", "random"} {
			for _, filter := range []string{"noop", "param2", "nodename", "all"} {
				cases = append(cases, patchBenchCase{edge, content, filter, 1})
			}
		}
	}
	// partial patches on the largest part
	for _, coverage := range []float64{0.01, 0.1} {
		cases = append(cases, patchBenchCase{16, "random", "all", coverage})
	}
	return cases
}

// build returns the stored part and the patch request for this case
func (c patchBenchCase) build(b *testing.B) (*types.SchemaPart, *types.PatchedSchemapart) {
	random := c.content == "random"
	stored := createBenchPartSeed(b, c.edge, random, 1)
	patchpart := createBenchPartSeed(b, c.edge, random, 2)
	if !random {
		// uniform patch with a different node than the stored part
		psp, err := parser.ParseSchemaPart(patchpart)
		if err != nil {
			b.Fatal(err)
		}
		psp.NodeNameLookup = map[int16]string{0: "mymod:other"}
		psp.Meta.NodeMapping = map[string]int16{"mymod:other": 0}
		for i := range psp.Param2 {
			psp.Param2[i] = 3
		}
		patchpart, err = psp.Convert()
		if err != nil {
			b.Fatal(err)
		}
	}

	rnd := rand.New(rand.NewSource(3))
	volume := c.edge * c.edge * c.edge
	pf := make([]types.PatchFilterType, volume)
	for i := range pf {
		if rnd.Float64() < c.coverage {
			pf[i] = patchBenchFilters[c.filter]
		}
	}
	return stored, &types.PatchedSchemapart{SchemaPart: patchpart, PatchFilter: pf}
}

func BenchmarkPatchSchemaPart(b *testing.B) {
	if !testutils.IsDatabaseAvailable() {
		b.SkipNow()
	}
	repos, err := db.Init()
	if err != nil {
		b.Fatal(err)
	}
	cfg := types.CreateConfig()
	a, _, err := api.NewApi(repos, cfg)
	if err != nil {
		b.Fatal(err)
	}

	user := &types.User{
		UID:  uuid.NewString(),
		Name: core.CreateToken(10),
		Type: types.UserTypeLocal,
		Role: types.UserRoleDefault,
	}
	if err = repos.UserRepo.CreateUser(user); err != nil {
		b.Fatal(err)
	}
	token, err := core.New(cfg, nil).CreateJWT(user, []types.JWTPermission{types.JWTPermissionUpload}, time.Hour)
	if err != nil {
		b.Fatal(err)
	}

	for _, c := range patchBenchCases() {
		b.Run(c.String(), func(b *testing.B) {
			schema := &types.Schema{
				UID:     uuid.NewString(),
				UserUID: user.UID,
				Name:    core.CreateToken(10),
				SizeX:   16,
				SizeY:   16,
				SizeZ:   16,
			}
			if err := repos.SchemaRepo.CreateSchema(schema); err != nil {
				b.Fatal(err)
			}

			stored, patch := c.build(b)
			stored.SchemaUID = schema.UID
			patch.SchemaUID = schema.UID
			if err := repos.SchemaPartRepo.CreateOrUpdateSchemaPart(stored); err != nil {
				b.Fatal(err)
			}

			body, err := json.Marshal(patch)
			if err != nil {
				b.Fatal(err)
			}

			handler := a.Secure(a.PatchSchemaPart)
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				r := httptest.NewRequest("PATCH", "http://", bytes.NewReader(body))
				r.Header.Set("Authorization", token)
				w := httptest.NewRecorder()
				handler(w, r)
				if w.Code != 200 {
					b.Fatalf("unexpected status %d: %s", w.Code, w.Body.String())
				}
			}
			b.ReportMetric(float64(len(body)), "req-bytes")
		})
	}
}

// BenchmarkPatchSchemaPartCPU measures the in-memory part of the PATCH handler
// (parse stored part, parse patch, apply, convert) without http, json and database
func BenchmarkPatchSchemaPartCPU(b *testing.B) {
	for _, c := range patchBenchCases() {
		b.Run(c.String(), func(b *testing.B) {
			stored, patch := c.build(b)
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				parsed, err := parser.ParseSchemaPart(stored)
				if err != nil {
					b.Fatal(err)
				}
				parsedPatch, err := parser.ParseSchemaPart(patch.SchemaPart)
				if err != nil {
					b.Fatal(err)
				}
				if err = parsed.Patch(parsedPatch, patch.PatchFilter); err != nil {
					b.Fatal(err)
				}
				if _, err = parsed.Convert(); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
