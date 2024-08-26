package core

import (
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
	Parts      chan *types.SchemaPart
}

type SchemaClient struct {
	opts *SchemaClientOpts
	pos1 *mt.Pos
	pos2 *mt.Pos
}

func NewSchemaClient(opts *SchemaClientOpts) *SchemaClient {
	pos1 := mt.NewPos(opts.Pull.PosX, opts.Pull.PosY, opts.Pull.PosZ)
	size := mt.NewPos(opts.Schema.SizeX, opts.Schema.SizeY, opts.Schema.SizeZ)
	pos2 := pos1.Add(size.Add(mt.NewPos(-1, -1, -1)))

	return &SchemaClient{
		opts: opts,
		pos1: pos1,
		pos2: pos2,
	}
}

func (sc *SchemaClient) blockDataHandler(ch chan commands.Command) {
	var ser_ver = uint8(28)
	id_node_mapping := map[int]string{}

	for o := range ch {
		switch cmd := o.(type) {
		case *commands.ServerHello:
			fmt.Printf("Version: %d\n", cmd.SerializationVersion)
			ser_ver = cmd.SerializationVersion
		case *commands.ServerNodeDefinitions:
			for _, ndef := range cmd.Definitions {
				id_node_mapping[int(ndef.ID)] = ndef.Name
			}
			fmt.Printf("Mapped %d nodedefs\n", len(cmd.Definitions))
		case *commands.ServerBlockData:
			mb_pos1 := cmd.Pos.Multiply(16)
			mb_pos2 := mb_pos1.Add(mt.NewPos(15, 15, 15))
			fmt.Printf("Blockdata: %d bytes, pos1: %s, pos2: %s\n", len(cmd.BlockData), mb_pos1, mb_pos2)

			mb, err := mapparser.ParseNetwork(ser_ver, cmd.BlockData)
			if err != nil {
				fmt.Printf("map parse error @ %v: %v\n", cmd.Pos, err)
			}
			if mb != nil {
				fmt.Printf("BlockMapping @ %v: %v\n", cmd.Pos, mb.BlockMapping)
			}
		}
	}
}

func (sc *SchemaClient) Run() error {

	client := commandclient.NewCommandClient(sc.opts.Pull.Hostname, sc.opts.Pull.Port)
	//go commandclient.DebugHandler(client)

	ch := make(chan commands.Command, 100)
	client.AddListener(ch)
	defer client.RemoveListener(ch)
	go sc.blockDataHandler(ch)

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

	time.Sleep(5 * time.Second)

	err = client.Disconnect()
	if err != nil {
		return fmt.Errorf("disconnect error: %v", err)
	}

	return nil
}
