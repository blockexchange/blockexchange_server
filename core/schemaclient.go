package core

import (
	"blockexchange/parser"
	"blockexchange/types"
	"fmt"
	"time"

	"github.com/minetest-go/mapparser"
	"github.com/minetest-go/minetest_client/commandclient"
	"github.com/minetest-go/minetest_client/commands"
	mt "github.com/minetest-go/types"
)

type SchemaClientOpts struct {
	Pull       *types.SchematicPull
	PullClient *types.SchematicPullClient
	Schema     *types.Schema
	SetNode    func(pos *mt.Pos, node *mt.Node) error
	SetMeta    func(pos *mt.Pos, md *parser.MetadataEntry) error
}

type SchemaClient struct {
	opts            *SchemaClientOpts
	origin          *mt.Pos
	area            *mt.Area
	id_node_mapping map[int]string
}

func NewSchemaClient(opts *SchemaClientOpts) *SchemaClient {
	origin := mt.NewPos(opts.Pull.PosX, opts.Pull.PosY, opts.Pull.PosZ)
	size := mt.NewPos(opts.Schema.SizeX, opts.Schema.SizeY, opts.Schema.SizeZ)
	pos2 := origin.Add(size.Add(mt.NewPos(-1, -1, -1)))

	return &SchemaClient{
		opts:            opts,
		origin:          origin,
		area:            mt.NewArea(origin, pos2),
		id_node_mapping: map[int]string{},
	}
}

func (sc *SchemaClient) applyBlockChanges(rel_pos1 *mt.Pos, mb *mt.MapBlock) error {

	for i := 0; i < 4096; i++ {
		pos := mt.NewPosFromIndex(i)

		node, err := mb.GetNode(pos)
		if err != nil {
			return fmt.Errorf("getnode error: %v", err)
		}

		if node != nil {
			s_pos := rel_pos1.Add(pos)
			sc.opts.SetNode(s_pos, node)
			// TODO: metadata
		}
	}

	return nil
}

func (sc *SchemaClient) blockDataHandler(ch chan commands.Command, errchan chan error) {
	var ser_ver = uint8(28)

	for o := range ch {
		switch cmd := o.(type) {
		case *commands.ServerHello:
			fmt.Printf("Version: %d\n", cmd.SerializationVersion)
			ser_ver = cmd.SerializationVersion
		case *commands.ServerNodeDefinitions:
			for _, ndef := range cmd.Definitions {
				sc.id_node_mapping[int(ndef.ID)] = ndef.Name
			}
			fmt.Printf("Mapped %d nodedefs\n", len(cmd.Definitions))
		case *commands.ServerBlockData:
			pos1 := cmd.Pos.Multiply(16)
			pos2 := pos1.Add(mt.NewPos(15, 15, 15))
			area := mt.NewArea(pos1, pos2)
			fmt.Printf("Blockdata: %d bytes, pos1: %s, pos2: %s\n", len(cmd.BlockData), pos1, pos2)

			if !area.Intersects(sc.area) {
				// no valid data
				continue
			}

			mb, err := mapparser.ParseNetwork(ser_ver, cmd.BlockData)
			if err != nil {
				fmt.Printf("map parse error @ %v: %v\n", cmd.Pos, err)
			}
			if mb != nil {
				fmt.Printf("BlockMapping @ %v: %v\n", cmd.Pos, mb.BlockMapping)
			}

			rel_pos1 := pos1.Subtract(sc.origin)
			err = sc.applyBlockChanges(rel_pos1, mb)
			if err != nil {
				errchan <- fmt.Errorf("apply block changes error: %v", err)
				return
			}
		}
	}
}

func (sc *SchemaClient) Run() error {

	client := commandclient.NewCommandClient(sc.opts.Pull.Hostname, sc.opts.Pull.Port)
	//go commandclient.DebugHandler(client)
	errchan := make(chan error)

	ch := make(chan commands.Command, 100)
	client.AddListener(ch)
	defer client.RemoveListener(ch)
	go sc.blockDataHandler(ch, errchan)

	err := client.Connect()
	if err != nil {
		return fmt.Errorf("connect error: %v", err)
	}

	err = commandclient.Init(client, sc.opts.PullClient.Username)
	if err != nil {
		return fmt.Errorf("init error: %v", err)
	}

	err = commandclient.Login(client, sc.opts.PullClient.Username, sc.opts.PullClient.Password, false)
	if err != nil {
		return fmt.Errorf("login error: %v", err)
	}

	go commandclient.ClientReady(client)

	select {
	case <-time.After(5 * time.Second):
		break
	case err = <-errchan:
		return fmt.Errorf("blockhandler error: %v", err)
	}

	err = client.Disconnect()
	if err != nil {
		return fmt.Errorf("disconnect error: %v", err)
	}

	return nil
}
