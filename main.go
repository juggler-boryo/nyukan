package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
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

func InitCardReader() (*CardReader, error) {
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

func (cr *CardReader) ReadID() error {
	if cr == nil || cr.card == nil {
		return fmt.Errorf("card reader not initialized")
	}

	// Configure for FeliCa cards (similar to Python's 212F setting)
	status, err := cr.card.Status()
	if err != nil {
		return fmt.Errorf("card status error: %v", err)
	}
	log.Println(status)

	// FeliCa polling command (similar to Python's sensf_req)
	// 0x00: system code filter
	// 0xFF: FeliCa system code for all cards
	command := []byte{0xFF, 0xCA, 0x00, 0x00, 0x00}

	response, err := cr.card.Transmit(command)
	if err != nil {
		return fmt.Errorf("transmit error: %v", err)
	}

	cr.idm = hex.EncodeToString(response[:6])
	fmt.Printf("Detected FeliCa card IDm: %s\n", cr.idm)
	return nil
}

func handleNFCCard(suicaID string) error {
	user, err := lib.FetchUserInfo(suicaID)
	if err != nil {
		return fmt.Errorf("Failed to get user information: %v", err)
	}
	log.Printf("User: %+v", user)

	if err := lib.UpdateInoutStatus(user.UID, !user.IsIn); err != nil {
		// this means server impl is fucked up
		return fmt.Errorf("Failed to update entry/exit status: %v", err)
	}

	action := "entered"
	if user.IsIn {
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

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		log.Println("Shutting down gracefully...")
		sound.PlayError()
		os.Exit(0)
	}()

	var cr *CardReader
	for {
		var err error
		cr, err = InitCardReader()
		if err != nil {
			log.Printf("Failed to initialize card reader: %v", err)
			time.Sleep(2 * time.Second)
			continue
		}

		// メインループ
		for {
			fmt.Println("Waiting for FeliCa card...")
			err := cr.ReadID()
			if err != nil {
				if err.Error() == "card status error: scard: Card was removed." {
					break // Break inner loop to reinitialize reader
				}

				if err.Error() == "card reader not initialized" {
					break // Break inner loop to reinitialize reader
				}

				log.Printf("Error reading card: %v. Reinitializing reader...", err)
				break // Break inner loop to reinitialize reader
			}

			sound.PlayTry()

			if err := handleNFCCard(cr.idm); err != nil {
				sound.PlayError()
				msg := "❌未登録のnfcカード: " + cr.idm + "\n ttps://aigrid.vercel.app/profile で登録してください"
				if err := lib.SendMessageToDiscord(lib.GetDiscordChannelID(), msg); err != nil {
					log.Printf("Failed to send Discord message: %v", err)
				}
				log.Println(msg)
				log.Println(err)
				time.Sleep(2 * time.Second)
				continue
			}

			time.Sleep(100 * time.Millisecond)
		}

		// Add small delay before trying to reinitialize
		time.Sleep(1 * time.Second)
	}
}
