package dtos

type UpdateUserProfile struct {
	FirstName          string            `json:"first_name"`
	LastName           string            `json:"last_name"`
	DateOfBirth        string            `json:"date_of_birth" validate:"datetime=2006-01-02"`
	Gender             string            `json:"gender"`
	PhoneNumber        string            `json:"phone_number"`
	Address            string            `json:"address"`
	City               string            `json:"city"`
	State              string            `json:"state"`
	Country            string            `json:"country"`
	PostalCode         string            `json:"postal_code"`
	Institution        string            `json:"institution"`
	Department         string            `json:"department"`
	JobTitle           string            `json:"job_title"`
	EducationLevel     string            `json:"education_level"`
	FieldOfExpertise   string            `json:"field_of_expertise"`
	YearsOfExperience  int               `json:"years_of_experience"`
	ProfilePicture     string            `json:"profile_picture"`
	Biography          string            `json:"biography"`
	GitHubProfile      string            `json:"github_profile"`
	Goals              string            `json:"goals"`
	Skills             []string          `json:"skills"`
	Interests          []string          `json:"interests"`
	Languages          []string          `json:"languages"`
	Certifications     []string          `json:"certifications"`
	Awards             []string          `json:"awards"`
	Publications       []string          `json:"publications"`
	Projects           []string          `json:"projects"`
	Hobbies            []string          `json:"hobbies"`
	SocialMediaHandles map[string]string `json:"social_media_handles"`
	Preferences        map[string]string `json:"preferences"`
	EmailPreferences   map[string]bool   `json:"email_preferences"`
	ContactPreferences map[string]bool   `json:"contact_preferences"`
	SecurityQuestions  map[string]string `json:"security_questions"`
	TwoFactorEnabled   bool              `json:"two_factor_enabled"`
	Timezone           string            `json:"timezone"`
	PreferredLanguage  string            `json:"preferred_language"`
	SolvedCount        int               `json:"solved_count"`
	MaxDifficulty      string            `json:"max_difficulty"`
	Reputation         float64           `json:"reputation"`
}
