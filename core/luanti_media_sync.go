package core

import (
	"blockexchange/types"
	"context"
	"crypto/sha1"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"

	"cirello.io/pglock"
	"github.com/sirupsen/logrus"

	"github.com/luanti-golang/luanti-client/client"
	"github.com/luanti-golang/luanti-client/protocol"
	"github.com/luanti-golang/luanti-client/viewerdef"
)

// upper bound for a whole sync run (connect, login and media download)
const luantiSyncTimeout = 5 * time.Minute

// timeout for connecting and logging in
const luantiConnectTimeout = 30 * time.Second

var ErrLuantiSyncRunning = errors.New("a sync for this connection is already running")

// media types that are stored, the same as the in-game mod uploads (textures and models)
var luantiMediaExtensions = map[string]bool{
	".png": true, ".jpg": true, ".jpeg": true, ".bmp": true, ".tga": true,
	".obj": true, ".b3d": true, ".x": true, ".gltf": true, ".glb": true,
}

type luantiMediaResult struct {
	Nodedefs   []*types.Nodedefinition
	Mediafiles []*types.Mediafile
	Size       int64
	Skipped    int // media files dropped because of a hash mismatch
}

// connects to the luanti server and collects the node definitions and media files it sends while joining
func fetchLuantiMedia(ctx context.Context, conn *types.LuantiMediaConnection) (*luantiMediaResult, error) {
	var (
		mu         sync.Mutex
		announced  = map[string][20]byte{}
		media      = map[string][]byte{}
		bunches    = map[uint16]bool{}
		skipped    = 0
		mediaDone  = make(chan struct{})
		doneOnce   sync.Once
		closeMedia = func() { doneOnce.Do(func() { close(mediaDone) }) }
	)

	log := logrus.WithField("connection", conn.ID)
	c := client.New(client.Config{
		Address:        net.JoinHostPort(conn.Hostname, strconv.Itoa(conn.Port)),
		Username:       conn.Username,
		Password:       conn.Password,
		Register:       true,
		RequestMedia:   true,
		ConnectTimeout: luantiConnectTimeout,
		Logf: func(format string, args ...any) {
			log.Debugf("luanti-client: "+format, args...)
		},
	})

	client.Handle(c, func(p *protocol.ToClientAnnounceMedia) {
		mu.Lock()
		for _, f := range p.Files {
			announced[f.Name] = f.SHA1
		}
		mu.Unlock()
		if len(p.Files) == 0 {
			// nothing to download
			closeMedia()
		}
	})
	client.Handle(c, func(p *protocol.ToClientMedia) {
		mu.Lock()
		for _, f := range p.Files {
			if !luantiMediaExtensions[strings.ToLower(path.Ext(f.Name))] {
				continue
			}
			if sum, ok := announced[f.Name]; !ok || sum != sha1.Sum(f.Data) {
				skipped++
				continue
			}
			media[f.Name] = f.Data
		}
		bunches[p.BunchIndex] = true
		done := len(bunches) >= int(p.NumBunches)
		mu.Unlock()
		if done {
			closeMedia()
		}
	})

	err := c.Connect(ctx)
	if err != nil {
		return nil, fmt.Errorf("connect: %w", err)
	}
	defer c.Close()

	// Connect may return before all media bunches arrived, wait for them explicitly
	select {
	case <-mediaDone:
	case <-c.Done():
		return nil, fmt.Errorf("connection closed while downloading media: %v", c.Wait())
	case <-ctx.Done():
		return nil, fmt.Errorf("media download: %w", ctx.Err())
	}

	defs := c.NodeDefs()
	if len(defs) == 0 {
		return nil, errors.New("no node definitions received")
	}

	res := &luantiMediaResult{}
	modnames := map[string]bool{}
	for name, def := range viewerdef.Convert(defs) {
		modname, _, found := strings.Cut(name, ":")
		if !found {
			// built-in nodes (air, ignore, unknown)
			continue
		}
		data, err := json.Marshal(def)
		if err != nil {
			return nil, fmt.Errorf("nodedef '%s': %w", name, err)
		}
		modnames[modname] = true
		res.Nodedefs = append(res.Nodedefs, &types.Nodedefinition{
			Name:       name,
			ModName:    modname,
			Definition: string(data),
		})
	}

	mu.Lock()
	defer mu.Unlock()
	for name, data := range media {
		res.Mediafiles = append(res.Mediafiles, &types.Mediafile{
			Name:    name,
			ModName: guessMediaModName(name, modnames),
			Data:    data,
		})
		res.Size += int64(len(data))
	}
	res.Skipped = skipped

	return res, nil
}

// the server does not tell which mod a media file belongs to, by convention
// the file name is prefixed with the mod name ("default_stone.png")
func guessMediaModName(filename string, modnames map[string]bool) string {
	base := strings.TrimSuffix(filename, path.Ext(filename))
	// longest matching prefix wins ("mesecons_wires_x.png" -> "mesecons_wires")
	for i := len(base); i > 0; i-- {
		if (i == len(base) || base[i] == '_') && modnames[base[:i]] {
			return base[:i]
		}
	}
	return ""
}

func (c *Core) syncLuantiMedia(conn *types.LuantiMediaConnection) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), luantiSyncTimeout)
	defer cancel()

	res, err := fetchLuantiMedia(ctx, conn)
	if err != nil {
		return "", err
	}

	err = c.repos.MediaRepo.UpsertNodedefinitions(res.Nodedefs)
	if err != nil {
		return "", fmt.Errorf("store node definitions: %w", err)
	}
	err = c.repos.MediaRepo.UpsertMediafiles(res.Mediafiles)
	if err != nil {
		return "", fmt.Errorf("store media files: %w", err)
	}

	msg := fmt.Sprintf("fetched %d node definitions and %d media files (%.1f MB)",
		len(res.Nodedefs), len(res.Mediafiles), float64(res.Size)/1024/1024)
	if res.Skipped > 0 {
		msg += fmt.Sprintf(", skipped %d media files with a hash mismatch", res.Skipped)
	}
	return msg, nil
}

// fetches and stores the node definitions and media files of the connection and writes the result back to it
func (c *Core) SyncLuantiMediaConnection(conn *types.LuantiMediaConnection) error {
	lock, err := c.repos.PGLock.Acquire("luanti-media-sync-"+conn.ID, pglock.FailIfLocked())
	if errors.Is(err, pglock.ErrNotAcquired) {
		return ErrLuantiSyncRunning
	}
	if err != nil {
		return fmt.Errorf("acquire lock: %w", err)
	}
	defer lock.Close()

	log := logrus.WithFields(logrus.Fields{"connection": conn.ID, "hostname": conn.Hostname, "port": conn.Port})
	log.Info("syncing luanti media")

	msg, syncErr := c.syncLuantiMedia(conn)
	if syncErr != nil {
		msg = syncErr.Error()
		log.WithError(syncErr).Warn("luanti media sync failed")
	} else {
		log.Info(msg)
	}

	conn.LastUpdatedTimestamp = time.Now().UnixMilli()
	conn.LastUpdatedSuccess = syncErr == nil
	conn.LastUpdatedMessage = msg
	return c.repos.LuantiMediaConnectionRepo.UpdateStatus(conn.ID, conn.LastUpdatedTimestamp, conn.LastUpdatedSuccess, conn.LastUpdatedMessage)
}

// syncs every connection whose last run is older than the given interval
func (c *Core) SyncDueLuantiMediaConnections(interval time.Duration) {
	list, err := c.repos.LuantiMediaConnectionRepo.GetAll()
	if err != nil {
		logrus.WithError(err).Error("list luanti media connections")
		return
	}
	for _, conn := range list {
		if time.Since(time.UnixMilli(conn.LastUpdatedTimestamp)) < interval {
			continue
		}
		err = c.SyncLuantiMediaConnection(conn)
		if err != nil && !errors.Is(err, ErrLuantiSyncRunning) {
			logrus.WithError(err).WithField("connection", conn.ID).Error("luanti media sync")
		}
	}
}
