package participant

// Participant holds information about a participant who is partaking in the secret santa draw
type Participant struct {
	Name  string
	Email string
}

func (p *Participant) GetName() string {
	if p != nil {
		return p.Name
	}

	return ""
}

func (p *Participant) GetEmail() string {
	if p != nil {
		return p.Email
	}

	return ""
}

func NewParticipant(name, email string) Participant {
	return Participant{Name: name, Email: email}
}
