package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/smtp"
	"os"
	"strings"

	"github.com/cameronouellette/secretsanta/internal/emails"
	"github.com/cameronouellette/secretsanta/internal/participant"
	"github.com/cameronouellette/secretsanta/internal/sender"
)

const (
	// e-mail subject line
	subject = "Le Père Noël secret!"

	// master list file name
	masterListFilePath = "../secret-santa-master-list.txt"
)

// santa contains the information of the email account that will be sending the secret santa selections
var santa = sender.NewSender(emails.SenderName, emails.SenderEmail, emails.SenderPassword)

// SMTP authentication
var auth = smtp.PlainAuth("", santa.GetEmail(), santa.GetPassword(), emails.SMTPHost)

func main() {
	args := os.Args[1:]

	if len(args) != 1 {
		log.Fatal("Please pass in exactly one person who forgot who they were shopping for.")
	}

	readParticipantFromFile(args[0])
}

func readParticipantFromFile(secretSantaName string) {
	// open the master list file and defer closing it
	file, err := os.Open(masterListFilePath)
	if err != nil {
		log.Fatal("could not open the master list file at the given path! Error:", err)
	}
	defer file.Close()

	// read the entire file, which is base64 encoded
	encodedMasterList, err := os.ReadFile(masterListFilePath)
	if err != nil {
		log.Fatal("could not read the master list file at the given path! Error:", err)
	}

	// base64 decode the file
	decodedMasterList, err := base64.StdEncoding.DecodeString(string(encodedMasterList))
	if err != nil {
		log.Fatal("could not decode the encoded master list! Error:", err)
	}

	// re-construct the master list based on the contents of the file, return early if you find the secret santa we're looking for
	listOfAllSecretSantaNames := make([]string, 0)
	lines := strings.Split(string(decodedMasterList), "\n")
	for _, line := range lines {

		pair := strings.Split(line, " : ")
		if len(pair) != 2 {
			continue
		}

		secretSantaInfo := strings.Split(pair[0], ",")
		if len(secretSantaInfo) != 2 {
			continue
		}
		secretSanta := participant.NewParticipant(secretSantaInfo[0], secretSantaInfo[1])

		participantInfo := strings.Split(pair[1], ",")
		if len(participantInfo) != 2 {
			continue
		}
		participant := participant.NewParticipant(participantInfo[0], participantInfo[1])

		if secretSanta.GetName() == secretSantaName {
			// we found the secret santa pairing!
			body := fmt.Sprintf("Hi %s!\n\n\n"+
				"This Christmas you will be %s's Secret Santa! :)\n\n"+
				"The spending limit this year is 30$ before tax. Have fun shopping!\n\n\n"+
				"Joyeux Noël!\n"+
				"Le Père Noël",
				secretSanta.GetName(),
				participant.GetName())

			// contruct the message with headers and body included
			message := constructMessage(santa, secretSanta, subject, body)

			// send the message
			fmt.Printf("\n\nSending secret santa e-mail message to %s\n", secretSanta.GetName())
			// UNCOMMENT FOR DEBUGGING
			// fmt.Printf("-------------------------------\n%s\n-------------------------------\n", message)
			err := smtp.SendMail(emails.SMTPAddr, auth, santa.GetEmail(), []string{secretSanta.GetEmail()}, []byte(message))
			if err != nil {
				fmt.Println("Message failed to send to ", secretSanta.GetName(), ". Error: ", err)
			}
			fmt.Println("E-mail sent successfully to", secretSanta.GetName(), "!")
			return
		}

		listOfAllSecretSantaNames = append(listOfAllSecretSantaNames, secretSanta.GetName())
	}

	fmt.Println("Sorry, but the master file did not contain the entered name as a secret santa! Are you sure you typed it correctly (first letter capitalised, spelling, etc.)?")
	fmt.Println("Here is a list of all the names on the master list :")
	for secretSanta := range listOfAllSecretSantaNames {
		fmt.Println(secretSanta)
	}
}

func constructMessage(sender sender.Sender, participant participant.Participant, subject, body string) (message string) {
	// construct SMTP formatted message
	message = "From:" + sender.GetName() + "\n" +
		"To:" + participant.GetEmail() + "\n" +
		"Subject:" + subject + "\n\n" +
		body

	return
}
