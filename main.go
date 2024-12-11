package main

import (
	"fmt"
	"log"

	"nyukan/lib"
)

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

	// TODO: play success sound

	action := "entered"
	if !user.IsIn {
		action = "exited"
	}
	log.Printf("%s has %s", user.Username, action)
	return nil
}

func main() {
	log.Println("nyukan: NFC card reading system")

	suicaID := "test_unregistered_suica_id"
	if err := handleNFCCard(suicaID); err != nil {
		// TODO: play error sound
		msg := "❌未登録のnfcカード: " + suicaID
		lib.SendMessageToDiscord(lib.GetDiscordChannelID(), msg)
	}
}
