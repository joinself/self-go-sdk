package fact

var sourceDefinition = []byte(`{
	"sources": {
		"user_specified": [
			"document_number",
			"display_name",
			"email_address",
			"phone_number"
		],
		"passport": [
			"document_number",
			"surname",
			"given_names",
			"date_of_birth",
			"date_of_expiration",
			"sex",
			"nationality",
			"country_of_issuance"
		],
		"driving_license": [
			"document_number",
			"surname",
			"given_names",
			"date_of_birth",
			"date_of_issuance",
			"date_of_expiration",
			"address",
			"issuing_authority",
			"place_of_birth"
		],
		"identity_card": [
			"document_number",
			"surname",
			"given_names",
			"date_of_birth",
			"date_of_expiration",
			"sex",
			"nationality",
			"country_of_issuance"
		],
		"twitter": [
			"account_id",
			"nickname"
		],
		"linkedin": [
			"account_id",
			"nickname"
		],
		"facebook": [
			"account_id",
			"nickname"
		],
		"live": [
			"selfie_verification"
		]
	}
}`)

const SourcePassport = "passport"
const FactDocumentNumber = "document_number"
const FactSurname = "surname"
const FactGivenNames = "given_names"
const FactDateOfBirth = "date_of_birth"
const FactDateOfExpiration = "date_of_expiration"
const FactSex = "sex"
const FactNationality = "nationality"
const FactCountryOfIssuance = "country_of_issuance"
const SourceDrivingLicense = "driving_license"
const FactDateOfIssuance = "date_of_issuance"
const FactAddress = "address"
const FactIssuingAuthority = "issuing_authority"
const FactPlaceOfBirth = "place_of_birth"
const SourceIdentityCard = "identity_card"
const SourceTwitter = "twitter"
const FactAccountId = "account_id"
const FactNickname = "nickname"
const SourceLinkedin = "linkedin"
const SourceFacebook = "facebook"
const SourceLive = "live"
const FactSelfieVerification = "selfie_verification"
const SourceUserSpecified = "user_specified"
const FactDisplayName = "display_name"
const FactEmailAddress = "email_address"
const FactPhoneNumber = "phone_number"

var spec = map[string][]string{
	SourcePassport: []string{
		FactDocumentNumber,
		FactSurname,
		FactGivenNames,
		FactDateOfBirth,
		FactDateOfExpiration,
		FactSex,
		FactNationality,
		FactCountryOfIssuance,
	},
	SourceDrivingLicense: []string{
		FactDocumentNumber,
		FactSurname,
		FactGivenNames,
		FactDateOfBirth,
		FactDateOfIssuance,
		FactDateOfExpiration,
		FactAddress,
		FactIssuingAuthority,
		FactPlaceOfBirth,
	},
	SourceIdentityCard: []string{
		FactDocumentNumber,
		FactSurname,
		FactGivenNames,
		FactDateOfBirth,
		FactDateOfExpiration,
		FactSex,
		FactNationality,
		FactCountryOfIssuance,
	},
	SourceTwitter: []string{
		FactAccountId,
		FactNickname,
	},
	SourceLinkedin: []string{
		FactAccountId,
		FactNickname,
	},
	SourceFacebook: []string{
		FactAccountId,
		FactNickname,
	},
	SourceLive: []string{
		FactSelfieVerification,
	},
	SourceUserSpecified: []string{
		FactDocumentNumber,
		FactDisplayName,
		FactEmailAddress,
		FactPhoneNumber,
	},
}
