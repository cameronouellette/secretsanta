package sender

// Sender holds information about the account that will send emails to everyone
type Sender struct {
	Name     string
	Email    string
	Password string
}

func (s *Sender) GetName() string {
	if s != nil {
		return s.Name
	}

	return ""
}

func (s *Sender) GetEmail() string {
	if s != nil {
		return s.Email
	}

	return ""
}

func (s *Sender) GetPassword() string {
	if s != nil {
		return s.Password
	}

	return ""
}

// NewSender creates a new sender with the given fields
func NewSender(name, email, password string) Sender {
	return Sender{Name: name, Email: email, Password: password}
}
