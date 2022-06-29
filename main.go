package main

import (
	"bufio"
	"encoding/base64"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/smtp"
	"os"
	"time"

	. "github.com/cameronouellette/secretsanta/internal/emails"
	. "github.com/cameronouellette/secretsanta/internal/participant"
	. "github.com/cameronouellette/secretsanta/internal/sender"
)

const (
	// e-mail subject line
	subjectPreface = "Le Père Noël secret!"

	// master list file name
	masterListFilePath = "secret-santa-master-list.txt"
)

// master list mapping participants to their secret santa (i.e. map[person] = person's secret santa)
// - this list will be encoded and stored in a file and should only be read from if someone forgets who they drew and/or accidentally lost the email
var masterListOfSecretSantaSelections = make(map[Participant]Participant)

// participantBlockList is used to block participants from getting each other.
// - couples are hard-coded to start, but this map may be used dynamically to enforce different rules if desired (e.g. two people cannot be each other's secret santa)
var participantBlockList = map[string][]string{

	// --------------- couples ---------------
	"Martine": {"Cameron"},
	"Cameron": {"Martine"},

	"Jan":    {"Pierre"},
	"Pierre": {"Jan"},

	"Annick": {"Phil"},
	"Phil":   {"Annick"},

	// ------ previous year's selections ------

}

// sender contains the information of the email account that will be sending the secret sender selections
var sender = NewSender(SenderName, SenderEmail, SenderPassword)

// SMTP authentication
var auth = smtp.PlainAuth("", sender.GetEmail(), sender.GetPassword(), SMTPHost)

var debug *bool
var attempt *int

func main() {
	debug = flag.Bool("debug", false, "perform a dry-run of the program where no emails are sent and the master list is printed to stdout")
	attempt = flag.Int("attempt", 1, "re-run the secret-santa selection")
	flag.Parse()

	for {
		// construct a list of participants based on the map of names to emails -- this list will represent the metaphorical names in the metaphorical hat
		participants := constructParticipants()

		// construct a list of participants that will represent the secret santas -- the people picking the metaphorical names from the metaphorical hat
		secretSantas := constructParticipants()

		// ensure the master list is reset before entering the selection process
		masterListOfSecretSantaSelections = make(map[Participant]Participant)

		if selectSecretSantas(participants, secretSantas) {
			break
		}

		fmt.Printf("\n\n\n SOMEONE ENDED UP WITHOUT A PARTNER, TRYING AGAIN \n\n\n")
	}

	// print master list if program is run in debug mode
	if *debug {
		prettyPrintMasterList()
	}

	writeMasterListToFile()

	// once everyone is matched, send emails
	for participant, secretSanta := range masterListOfSecretSantaSelections {

		// construct the message body
		body := fmt.Sprintf("Hi %s!\n\n\n", secretSanta.GetName())

		if *attempt > 1 {
			body += fmt.Sprintf("Please ignore the last selection -- it turns out that someone picked someone they weren't supposed to! Time for attempt #%d :P\n\n\n", *attempt)
		}

		body += fmt.Sprintf("This Christmas you will be %s's Secret Santa! :)\n\n"+
			"The spending limit this year is 30$ before tax. Have fun shopping!\n\n\n"+
			"Joyeux Noël!\n"+
			"Le Père Noël",
			participant.GetName())

		// contruct the message with headers and body included
		subject := subjectPreface
		if *attempt > 1 {
			subject += fmt.Sprintf(" Tentative #%d", *attempt)
		}
		message := constructMessage(sender, secretSanta, subject, body)

		// do not send emails if program is run in debug mode, but rather print them to stdout
		if *debug {
			fmt.Printf("-------------------------------\n%s\n-------------------------------\n", message)
		} else {
			// send the message
			fmt.Printf("\n\nSending secret santa e-mail message to %s\n", secretSanta.GetName())

			err := smtp.SendMail(SMTPAddr, auth, sender.GetEmail(), []string{secretSanta.GetEmail()}, []byte(message))
			if err != nil {
				fmt.Println("Message failed to send to ", secretSanta.GetName(), ". Error: ", err)
			}

			fmt.Println("E-mail sent successfully to", secretSanta.GetName(), "!")
		}
	}
}

// returns true if everyone is matched, false otherwise
func selectSecretSantas(participants, secretSantas []Participant) bool {

	// selection process -- select secret santas and send out messages
	for _, secretSanta := range secretSantas {

		// choose the participant that secretSanta will be shopping for
		participant := selectParticipant(&participants, secretSanta)

		// if participant is empty then not everyone was matched so return false
		if (participant == Participant{}) {
			return false
		}

		// add to master list
		masterListOfSecretSantaSelections[participant] = secretSanta

		// add the secretSanta to the selected participant's blocklist -- UNCOMMENT if participants are not allowed to shop for each other
		//participantBlockList[participant.name] = append(participantBlockList[participant.name], secretSanta.name)

	}

	return true
}

func constructMessage(sender Sender, participant Participant, subject, body string) (message string) {
	// construct SMTP formatted message
	message = "From:" + sender.GetName() + "\n" +
		"To:" + participant.GetEmail() + "\n" +
		"Subject:" + subject + "\n\n" +
		body

	return
}

func shuffleParticipants(participants *[]Participant) {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(*participants), func(i, j int) { (*participants)[i], (*participants)[j] = (*participants)[j], (*participants)[i] })
}

func selectParticipant(participants *[]Participant, secretSanta Participant) (selectedParticipant Participant) {

	// shuffle list of participants before looping through the participants -- this will increase the chances that everyone has a partner in the end
	shuffleParticipants(participants)

	// cycle through the list of shuffled participants and select the first one that's allowed
	for i, participant := range *participants {

		// cannot be your own secret santa
		if participant == secretSanta {
			continue
		}

		// one person cannot have multiple secret santa's
		if _, found := masterListOfSecretSantaSelections[participant]; found {
			continue
		}

		// check the participant against the secret santa's blocklist
		if isBlocklisted(secretSanta, participant) {
			continue
		}

		// select this participant and then remove them from the selection pool (their name has been drawn from the metaphorical hat, if you will)
		selectedParticipant = participant

		(*participants) = append((*participants)[:i], (*participants)[i+1:]...)
		break
	}

	return
}

// returns true if the participants is on the secretSanta's blocklist, false otherwise
func isBlocklisted(secretSanta Participant, participant Participant) bool {
	if blockedParticipants, found := participantBlockList[secretSanta.GetName()]; found {
		for _, blockedParticipant := range blockedParticipants {
			if blockedParticipant == participant.GetName() {
				return true
			}
		}
	}

	return false
}

func constructParticipants() (participants []Participant) {
	for name, email := range ParticipantMap {
		participants = append(participants, NewParticipant(name, email))
	}

	return
}

func prettyPrintMasterList() {
	fmt.Println("The following is a list of \"secret santa : person they're shopping for\" :")

	fmt.Println(constructMasterListPrettyPrintString())
}

func writeMasterListToFile() {
	// create our master list file and defer closing it
	file, err := os.Create(masterListFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// create a buffered writer and defer flushing it
	writer := bufio.NewWriter(file)
	defer writer.Flush()

	// write the base64-encoded master list to avoid accidental peeking
	encodedMasterList := base64.StdEncoding.EncodeToString([]byte(constructMasterListPrettyPrintString()))
	_, err = writer.WriteString(encodedMasterList)
	if err != nil {
		log.Fatal("Error writing the encoded master list to the given path! Error:", err)
	}
}

func constructMasterListPrettyPrintString() string {
	prettyPrintString := fmt.Sprintf("Attempt #%d\n", *attempt)
	for participant, secretSanta := range masterListOfSecretSantaSelections {
		prettyPrintString += secretSanta.GetName() + "," + secretSanta.GetEmail() + " : " + participant.GetName() + "," + participant.GetEmail() + "\n"
	}

	return prettyPrintString
}
