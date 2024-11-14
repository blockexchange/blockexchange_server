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

	fmt.Printf("Origin: %v, pos2: %v, size: %v\n", origin, pos2, size)

	return &SchemaClient{
		opts:   opts,
		origin: origin,
		area:   mt.NewArea(origin, pos2),
		id_node_mapping: map[int]string{
			CONTENT_AIR:     "air",
			CONTENT_IGNORE:  "ignore",
			CONTENT_UNKNOWN: "unknown",
		},
	}
}

func (sc *SchemaClient) applyBlockChanges(mb_startpos *mt.Pos, mb *mt.MapBlock) error {

	for i := 0; i < 4096; i++ {
		pos := mt.NewPosFromIndex(i)
		pos1 := mb_startpos.Add(pos)
		if !pos1.IsWithin(sc.area.Pos1, sc.area.Pos2) {
			// outside schematic
			continue
		}

		nodeid := mb.ContentId[i]

		if nodeid == CONTENT_IGNORE || nodeid == CONTENT_UNKNOWN {
			// don't map those
			continue
		}

		// relative position inside the schematic
		s_pos := pos1.Subtract(sc.origin)
		node := &mt.Node{
			Pos:    s_pos,
			Name:   sc.id_node_mapping[nodeid],
			Param1: mb.Param1[i],
			Param2: mb.Param2[i],
		}
		sc.opts.SetNode(s_pos, node)

		md := &parser.MetadataEntry{
			Inventories: mb.Inventory[i],
			Fields:      mb.Fields[i],
		}
		sc.opts.SetMeta(pos, md)
	}

	return nil
}

func (sc *SchemaClient) blockDataHandler(client *commandclient.CommandClient, errchan chan error) {
	ch := client.CommandChannel()
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
				errchan <- fmt.Errorf("map parse error @ %v: %v", cmd.Pos, err)
				return
			}

			err = sc.applyBlockChanges(pos1, mb)
			if err != nil {
				errchan <- fmt.Errorf("apply block changes error: %v", err)
				return
			}
		}
	}
}

func (sc *SchemaClient) Run() error {

	client := commandclient.NewCommandClient(sc.opts.Pull.Hostname, sc.opts.Pull.Port)
	errchan := make(chan error)

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

	err = clientReady(client)
	if err != nil {
		return fmt.Errorf("clientready error: %v", err)
	}

	go sc.blockDataHandler(client, errchan)

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

func clientReady(client *commandclient.CommandClient) error {
	for o := range client.CommandChannel() {
		switch cmd := o.(type) {
		case *commands.ServerCSMRestrictionFlags:
			err := client.SendCommand(commands.NewClientReady(5, 5, 5, "bx-bot", 4))
			if err != nil {
				return fmt.Errorf("send command error: %v", err)
			}

			ppos := commands.NewClientPlayerPos()
			err = client.SendCommand(ppos)
			if err != nil {
				return fmt.Errorf("send command error: %v", err)
			}

		case *commands.ServerMovePlayer:
			ppos := commands.NewClientPlayerPos()
			ppos.Pitch = uint32(cmd.Pitch)
			ppos.Yaw = uint32(cmd.Yaw)
			ppos.RequestViewRange = 15
			ppos.FOV = 149
			ppos.PosX = uint32(cmd.X * 10)
			ppos.PosY = uint32(cmd.Y * 10)
			ppos.PosZ = uint32(cmd.Z * 10)
			err := client.SendCommand(ppos)
			if err != nil {
				return fmt.Errorf("send command error: %v", err)
			}
			return nil
		}
	}
	return nil
}
