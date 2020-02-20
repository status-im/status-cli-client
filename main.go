package main

import (
	"context"
	"crypto/ecdsa"
	"database/sql"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	crypto "github.com/ethereum/go-ethereum/crypto"
	gethbridge "github.com/status-im/status-go/eth-node/bridge/geth"
	gonode "github.com/status-im/status-go/node"
	params "github.com/status-im/status-go/params"
	status "github.com/status-im/status-go/protocol"
	protobuf "github.com/status-im/status-go/protocol/protobuf"
)

func exitErr(err error) {
	fmt.Println(err)
	os.Exit(1)
}

type keysGetter struct {
	privateKey *ecdsa.PrivateKey
}

// flag variables
var port, timeout int
var addr, chatName, keyHex, dataDir, message, ensName string

func flagsInit() {
	flag.IntVar(&port, "port", 30303, "Listening port for Whisper node thread.")
	flag.IntVar(&timeout, "timeout", 500, "Timeout for message delivery in milliseconds.")
	flag.StringVar(&addr, "addr", "0.0.0.0", "Listening address for Whisper node thread.")
	flag.StringVar(&chatName, "chat", "whatever", "Name of public chat to send to.")
	flag.StringVar(&keyHex, "key", "", "Hex private key for your Status identity.")
	flag.StringVar(&dataDir, "data", "/tmp/status-cli-client", "Location for Status data.")
	flag.StringVar(&message, "message", "TEST", "Message to send to the public channel.")
	flag.StringVar(&ensName, "ens", "", "ENS name to send with the message.")
	flag.Parse()
}

func main() {
	flagsInit()

	// Generate private key if it's not provided through a CLI flag
	var privateKey *ecdsa.PrivateKey
	if keyHex == "" {
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
	fmt.Printf("Private Key: %#x\n", crypto.FromECDSA(privateKey))

	nodeConfig, err := generateConfig(dataDir, addr, port)
	if err != nil {
		exitErr(errors.Wrap(err, "failed generate config"))
	}

	statusNode := gonode.New()
	accsMgr, _ := statusNode.AccountManager()

	if err := statusNode.Start(nodeConfig, accsMgr); err != nil {
		exitErr(errors.Wrap(err, "failed to start node"))
	}

	// Using an in-memory SQLite DB since we have nothing worth preserving
	db, err := sql.Open("sqlite3", "file:mem?mode=memory&cache=shared")
	if err != nil {
		exitErr(errors.Wrap(err, "failed to open sqlite database"))
	}

	// Create a custom logger to suppress Status logs
	logger := zap.NewNop()

	options := []status.Option{
		status.WithDatabase(db),
		status.WithCustomLogger(logger),
	}

	var instID string = uuid.New().String()

	messenger, err := status.NewMessenger(
		privateKey,
		gethbridge.NewNodeBridge(statusNode.GethNode()),
		instID,
		options...,
	)
	if err != nil {
		exitErr(errors.Wrap(err, "failed to create Messenger"))
	}
	if err := messenger.Start(); err != nil {
		exitErr(errors.Wrap(err, "failed to start Messenger"))
	}
	if err := messenger.Init(); err != nil {
		exitErr(errors.Wrap(err, "failed to init Messenger"))
	}

	// Join the channel
	chat := status.CreatePublicChat(chatName, messenger.Timesource())
	messenger.Join(chat)
	messenger.SaveChat(&chat)

	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(timeout)*time.Millisecond)
	defer cancel()

	// Crate the Status chat message
	var statusMsg status.Message
	statusMsg.Text = message
	statusMsg.ChatId = chatName
	statusMsg.EnsName = ensName
	statusMsg.ContentType = protobuf.ChatMessage_TEXT_PLAIN
	fmt.Println("Destination:", chatName)
	fmt.Println("Message:", message)

	resp, err := messenger.SendChatMessage(ctx, &statusMsg)
	if err != nil {
		exitErr(errors.Wrap(err, "failed to send message"))
	}

	// keccak256(compressedAuthorPubKey, data)
	fmt.Printf("Message ID: %v\n", resp.Messages[0].ID)

	// FIXME this is an ugly hack, wait for delivery event properly
	time.Sleep(time.Duration(timeout+100) * time.Millisecond)
}

// Generate a sane configuration for a Status Node
func generateConfig(dataDir, addr string, port int) (*params.NodeConfig, error) {
	options := []params.Option{
		params.WithFleet(params.FleetProd),
		withListenAddr(addr, port),
	}

	var configFiles []string
	config, err := params.NewNodeConfigWithDefaultsAndFiles(
		dataDir,
		params.MainNetworkID,
		options,
		configFiles,
	)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func withListenAddr(addr string, port int) params.Option {
	return func(c *params.NodeConfig) error {
		c.ListenAddr = fmt.Sprintf("%s:%d", addr, port)
		return nil
	}
}
