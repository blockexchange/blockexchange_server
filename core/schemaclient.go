package core

import (
	"blockexchange/types"
	"fmt"
	"time"

	"github.com/minetest-go/mapparser"
	"github.com/minetest-go/minetest_client/commandclient"
	"github.com/minetest-go/minetest_client/commands"
)

type SchemaClientOpts struct {
	Pull       *types.SchematicPull
	PullClient *types.SchematicPullClient
	Schema     *types.Schema
	Parts      chan *types.SchemaPart
}

type SchemaClient struct {
	opts *SchemaClientOpts
}

func NewSchemaClient(opts *SchemaClientOpts) *SchemaClient {
	return &SchemaClient{opts: opts}
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
			node_offset := cmd.Pos.Multiply(16)
			fmt.Printf("Blockdata: %d bytes, node_offset: %v\n", len(cmd.BlockData), node_offset)

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
