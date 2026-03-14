package identity

//go:generate go-enum -f=$GOFILE --marshal

// Gender represents the gender of a human identity.
// ENUM(
// Unspecified="GENDER_UNSPECIFIED"
// Female="GENDER_FEMALE"
// Male="GENDER_MALE"
// Diverse="GENDER_DIVERSE"
// )
type Gender string

type Email struct {
	Address    string
	IsVerified bool
}

type Human struct {
	GivenName         string
	FamilyName        string
	NickName          *string
	DisplayName       string
	Gender            Gender
	PreferredLanguage *string
}

type Identity struct {
	ID       int64
	Human    *Human
	Email    Email
	IsActive bool
}
