package auth

type ChulaSSOCredential struct {
	UID         string   `json:"uid"`
	Username    string   `json:"username"`
	Gecos       string   `json:"gecos"`
	Email       string   `json:"email"`
	Disable     bool     `json:"disable"`
	Roles       []string `json:"roles"`
	Firstname   string   `json:"firstname"`
	Lastname    string   `json:"lastname"`
	FirstnameTH string   `json:"firstnameth"`
	LastnameTH  string   `json:"lastnameth"`
	Ouid        string   `json:"ouid"`
}
