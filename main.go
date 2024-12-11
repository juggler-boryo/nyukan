package main

import (
	"fmt"
	"log"
	"time"

	"encoding/hex"
	"nyukan/lib"
	"nyukan/sound"

	"github.com/ebfe/scard"
)

type CardReader struct {
	idm  string
	card *scard.Card
}

func NewCardReader() (*CardReader, error) {
	ctx, err := scard.EstablishContext()
	if err != nil {
		return nil, err
	}

	readers, err := ctx.ListReaders()
	if err != nil {
		return nil, err
	}

	// Check if there are any readers available
	if len(readers) == 0 {
		return nil, fmt.Errorf("no card readers found")
	}

	card, err := ctx.Connect(readers[0], scard.ShareShared, scard.ProtocolAny)
	if err != nil {
		return nil, err
	}

	return &CardReader{card: card}, nil
}

func (cr *CardReader) ReadID() {
	for {
		status, err := cr.card.Status()
		if err == nil && uint32(status.State)&uint32(scard.StatePresent) != 0 {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}

	command := []byte{0xFF, 0xCA, 0x00, 0x00, 0x00}
	response, err := cr.card.Transmit(command)
	if err != nil {
		log.Printf("Failed to transmit command: %v", err)
		return
	}

	cr.idm = hex.EncodeToString(response[:8])
	fmt.Printf("Detected card IDm: %s\n", cr.idm)
}

func handleNFCCard(suicaID string) error {
	// TODO: play weird sound
	// Get user information
	user, err := lib.FetchUserInfo(suicaID)
	if err != nil {
		return fmt.Errorf("Failed to get user information: %v", err)
	}
	log.Printf("User: %+v", user)

	// // Update entry/exit status
	// if err := lib.UpdateInoutStatus(user.UID, !user.IsIn); err != nil {
	// 	// this means server impl is fucked up
	// 	return fmt.Errorf("Failed to update entry/exit status: %v", err)
	// }

	action := "entered"
	if !user.IsIn {
		action = "exited"
	}

	if action == "entered" {
		sound.PlayIn()
	} else {
		sound.PlayBB()
	}

	log.Printf("%s has %s", user.Username, action)
	return nil
}

func main() {
	log.Println("nyukan: NFC card reading system")
	sound.PlayConnect()

	var cr *CardReader
	var err error
	for {
		cr, err = NewCardReader()
		if err != nil {
			log.Printf("Failed to initialize card reader: %v. Retrying in 5 seconds...", err)
			sound.PlayError()
			time.Sleep(5 * time.Second)
			continue
		}
		break
	}

	for {
		fmt.Println("Waiting for FeliCa card...")
		cr.ReadID()

		sound.PlayTry()
		if err := handleNFCCard(cr.idm); err != nil {
			sound.PlayError()
			msg := "❌未登録のnfcカード: " + cr.idm
			lib.SendMessageToDiscord(lib.GetDiscordChannelID(), msg)
			log.Println(msg)
			log.Println(err)
		}

		time.Sleep(2 * time.Second)
	}
}
