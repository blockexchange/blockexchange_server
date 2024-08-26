package core

import (
	"blockexchange/types"
	"fmt"
	"os"
	"time"

	"github.com/minetest-go/mapparser"
	"github.com/minetest-go/minetest_client/commandclient"
	"github.com/minetest-go/minetest_client/commands"
)

func blockDataHandler(ch chan commands.Command) {
	var ser_ver = uint8(28)

	for o := range ch {
		switch cmd := o.(type) {
		case *commands.ServerHello:
			fmt.Printf("Version: %d\n", cmd.SerializationVersion)
			ser_ver = cmd.SerializationVersion
		case *commands.ServerBlockData:
			fmt.Printf("Blockdata: %d bytes, pos: %v\n", len(cmd.BlockData), cmd.Pos)

			offsetData := make([]byte, len(cmd.BlockData)+1)
			offsetData[0] = ser_ver
			copy(offsetData[1:], cmd.BlockData)

			f, err := os.CreateTemp(os.TempDir(), "mapblock")
			if err != nil {
				fmt.Printf("create temp error: %v\n", err)
				return
			}
			f.Write(offsetData)
			f.Close()

			mb, err := mapparser.Parse(offsetData)
			if err != nil {
				fmt.Printf("map parse error @ %v: %v\n", cmd.Pos, err)
			}
			if mb != nil {
				fmt.Printf("BlockMapping @ %v: %v\n", cmd.Pos, mb.BlockMapping)
			}
		}
	}
}

func SchemaClient(p *types.SchematicPull, pc *types.SchematicPullClient) error {

	client := commandclient.NewCommandClient(p.Hostname, p.Port)
	//go commandclient.DebugHandler(client)

	ch := make(chan commands.Command, 100)
	client.AddListener(ch)
	defer client.RemoveListener(ch)
	go blockDataHandler(ch)

	err := client.Connect()
	if err != nil {
		return fmt.Errorf("connect error: %v", err)
	}

	err = commandclient.Init(client, pc.Username)
	if err != nil {
		return fmt.Errorf("init error: %v", err)
	}

	err = commandclient.Login(client, pc.Username, pc.Password)
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
