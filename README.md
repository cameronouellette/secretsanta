# Secret Santa
Go script that can be used for secret santa selection :) 

## Running the secret santa selection process

1. Create and populate an `emails.go` file or equivalent based on the example file `emails.go.example`.  
  a) Delete the `.example` file extension  
  b) Fill in the sender information using an email account that you don't care about or that was created for this purpose alone. This is important since you have to input the password for this account and if you accidentally commit the password you don't want it to be a big deal.  
  c) Fill in the participants names and emails in the participant map  

**NOTE:** If the account you put in part b) is not a GMail account, you will have to change the SMTP authentication in `emails.go`!
  
2. Run `go run main.go` from the root directory of the repo.


## Reading the master list when someone forgets who they selected

1. From the root directory of the repo run `cd read`.
2. Run `go run read.go <insert_name_here>`.
