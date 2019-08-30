package main

import (
	"os"
	"fmt"
	"time"
	"flag"
	"context"
	"strings"
	"database/sql"
	"crypto/ecdsa"
	"github.com/pkg/errors"
	"github.com/google/uuid"

	"github.com/ethereum/go-ethereum/crypto"
	status "github.com/status-im/status-protocol-go"
	params "github.com/status-im/status-go/params"
	gonode "github.com/status-im/status-go/node"
)

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

// flag variables
var port int
var addr, chatName, keyHex, dataDir, message string

func flagsInit() {
	flag.IntVar(&port, "port", 30303, "Listening port for Whisper node thread.")
	flag.StringVar(&addr, "addr", "0.0.0.0", "Listening address for Whisper node thread.")
	flag.StringVar(&chatName, "chat", "whatever", "Name of public chat to send to.")
	flag.StringVar(&keyHex, "key", "", "Hex private key for your Status identity.")
	flag.StringVar(&dataDir, "data", "/tmp/status-cli-client", "Location for Status data.")
	flag.StringVar(&message, "message", "TEST", "Message to send to the public channel.")
	flag.Parse()
}

func main() {
	flagsInit()

	// Generate private key if it's not provided through a CLI flag
	var privateKey *ecdsa.PrivateKey
	if (keyHex == "") {
		if key, err := crypto.GenerateKey(); err != nil {
			exitErr(err)
		} else {
			privateKey = key
		}
	} else {
		if key, err := crypto.HexToECDSA(strings.TrimPrefix(keyHex, "0x")); err != nil {
			exitErr(err)
		} else {
			privateKey = key
		}
	}
	fmt.Printf("Using private key: %#x\n", crypto.FromECDSA(privateKey))

	var configFiles []string
	// create a new status-go config
	config, err := params.NewNodeConfigWithDefaultsAndFiles(
		dataDir,
		params.MainNetworkID,
		[]params.Option{
			params.WithFleet(params.FleetBeta),
			withListenAddr(fmt.Sprintf("%s:%d", addr, port)),
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

	// Using an in-memory SQLite DB since we have nothing worth preserving
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

	payload := []byte(message)
	chat := status.Chat{ID: chatName, Name: chatName}
	// TODO error handling
	messenger.Send(ctx, chat, payload)

	// FIXME this is an ugly hack, wait for delivery event properly
	time.Sleep(500 * time.Millisecond)
}
