package main

import (
	"os"
	"fmt"
	"time"
	"context"
	"database/sql"
	"crypto/ecdsa"
	"github.com/pkg/errors"
	"github.com/google/uuid"

	"github.com/ethereum/go-ethereum/crypto"
	status "github.com/status-im/status-protocol-go"
	params "github.com/status-im/status-go/params"
	gonode "github.com/status-im/status-go/node"
)

// just some random key I generated from status-console-client
// TODO generat key on the fly or take it from cli argumens
const keyHex string = "4f7db2c72e3dd604b2be4258680844f1b66c911ab13701f5f33f5f5c03103221"
const fleet string = params.FleetBeta
const dataDir string = "/tmp/status-cli-client"
const listenAddr string = ":30303"
const chatName string = "whatever"

func exitErr(err error) {
	fmt.Println(err)
	os.Exit(1)
}

func withListenAddr(listenAddr string) params.Option {
	return func(c *params.NodeConfig) error {
		c.ListenAddr = listenAddr
		return nil
	}
}

type keysGetter struct {
	privateKey *ecdsa.PrivateKey
}

func (k keysGetter) PrivateKey() (*ecdsa.PrivateKey, error) {
	return k.privateKey, nil
}

func main() {
	var configFiles []string
	// create a new status-go config
	config, err := params.NewNodeConfigWithDefaultsAndFiles(
		dataDir,
		params.MainNetworkID,
		[]params.Option{
			params.WithFleet(fleet),
			withListenAddr(listenAddr),
		},
		configFiles,
	)
	if err != nil {
		exitErr(err)
	}

	statusNode := gonode.New()

	accsMgr, _ := statusNode.AccountManager()

	if err := statusNode.Start(config, accsMgr); err != nil {
		exitErr(errors.Wrap(err, "failed to start node"))
	}

	shhService, err := statusNode.WhisperService()
	if err != nil {
		exitErr(errors.Wrap(err, "failed to get Whisper service"))
	}

	var instID string = uuid.New().String()

	privateKey, err := crypto.HexToECDSA(keyHex)
	if err != nil {
		exitErr(err)
	}

	db, _ := sql.Open("sqlite3", "file:mem?mode=memory&cache=shared")
	options := []status.Option{
		status.WithDatabase(db),
	}

	messenger, err := status.NewMessenger(
		privateKey,
		shhService,
		instID,
		options...,
	)
	if err != nil {
		exitErr(errors.Wrap(err, "failed to create Messenger"))
	}
	
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// TODO parametrize the message to be sent
	data := []byte("THIS IS A TEST MESSAGE")
	chat := status.Chat{ID: chatName, Name: chatName}
	// TODO error handling
	messenger.Send(ctx, chat, data)
}
