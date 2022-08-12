package fact

const FactAccountId = "account_id"
const FactAddress = "address"
const FactCountryOfIssuance = "country_of_issuance"
const FactDateOfBirth = "date_of_birth"
const FactDateOfExpiration = "date_of_expiration"
const FactDateOfIssuance = "date_of_issuance"
const FactDocumentNumber = "document_number"
const FactEmailAddress = "email_address"
const FactGivenNames = "given_names"
const FactIssuingAuthority = "issuing_authority"
const FactNationality = "nationality"
const FactNickname = "nickname"
const FactPhoneNumber = "phone_number"
const FactPlaceOfBirth = "place_of_birth"
const FactSelfieVerification = "selfie_verification"
const FactSex = "sex"
const FactSurname = "surname"
const SourceDrivingLicense = "driving_license"
const SourceFacebook = "facebook"
const SourceIdentityCard = "identity_card"
const SourceLinkedin = "linkedin"
const SourceLive = "live"
const SourcePassport = "passport"
const SourceTwitter = "twitter"
const SourceUserSpecified = "user_specified"

var spec = map[string][]string{
	SourceDrivingLicense: {
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
	SourceFacebook: {
		FactAccountId,
		FactNickname,
	},
	SourceIdentityCard: {
		FactDocumentNumber,
		FactSurname,
		FactGivenNames,
		FactDateOfBirth,
		FactDateOfExpiration,
		FactSex,
		FactNationality,
		FactCountryOfIssuance,
	},
	SourceLinkedin: {
		FactAccountId,
		FactNickname,
	},
	SourceLive: {
		FactSelfieVerification,
	},
	SourcePassport: {
		FactDocumentNumber,
		FactSurname,
		FactGivenNames,
		FactDateOfBirth,
		FactDateOfExpiration,
		FactSex,
		FactNationality,
		FactCountryOfIssuance,
	},
	SourceTwitter: {
		FactAccountId,
		FactNickname,
	},
	SourceUserSpecified: {
		FactDocumentNumber,
		FactEmailAddress,
		FactPhoneNumber,
	},
}
