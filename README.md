# Secret Santa
Go script that can be used for secret santa selection :) 

## Get ready to run the secret santa selection process

1. Create and populate an `emails.go` file or equivalent based on the example file `emails.go.example`.  
  a) Create a copy of `emails.go.example` and delete the `.example` file extension  
  b) Fill in the sender information using an email account that you don't care about or that was created for this purpose alone. This is important since you have to input the password for this account and if you accidentally commit the password you don't want it to be a big deal.  
  c) Fill in the participants names and emails in the participant map  

**NOTES:**  
- If using a GMail account as your sender, you may have to use an App Password to authenticate rather than the account's password.
- If the account you put in part b) is not a GMail account, you will have to change the SMTP authentication in `emails.go`.
  

## Running the selection process  

Run `go run main.go` from the root directory of the repo to run the selection process and send out emails to all the participants.

### Flags  

| flag | type | default | description |  
| ---- | ---- | ------- | ----------- |  
| `-debug` | bool | false | Performs a dry-run of the program. The emails to be sent will be printed to stdout as well as the master list of secret santa selections |  
| `-attempt` | int | 1 | Denotes the attempt number in the master list, the email subject line and the email message body |

### Blocklist
The application has a hardcoded `participantBlockList` in `main.go`. This can be used to enforce rules such as blocking couples from picking each other or blocking everyone from picking the same person as last year.

## Reading the master list when someone forgets who they selected

1. From the root directory of the repo run `cd read`.
2. Run `go run read.go <insert_name_of_person_who_forgot_here>`.
