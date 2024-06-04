package fact

const FactAddress = "address"
const FactCountryOfIssuance = "country_of_issuance"
const FactDateOfBirth = "date_of_birth"
const FactDateOfExpiration = "date_of_expiration"
const FactDateOfIssuance = "date_of_issuance"
const FactDisplayName = "display_name"
const FactDocumentNumber = "document_number"
const FactEmailAddress = "email_address"
const FactGivenNames = "given_names"
const FactIssuingAuthority = "issuing_authority"
const FactLocation = "location"
const FactNationality = "nationality"
const FactPhoto = "photo"
const FactPlaceOfBirth = "place_of_birth"
const FactSelfieVerification = "selfie_verification"
const FactSex = "sex"
const FactSurname = "surname"
const FactUnverifiedPhoneNumber = "unverified_phone_number"
const SourceDrivingLicense = "driving_license"
const SourceIdentityCard = "identity_card"
const SourceLive = "live"
const SourcePassport = "passport"
const SourceUserSpecified = "user_specified"

var spec = map[string][]string{
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
	SourceLive: []string{
		FactSelfieVerification,
	},
	SourcePassport: []string{
		FactPhoto,
		FactDocumentNumber,
		FactSurname,
		FactGivenNames,
		FactDateOfBirth,
		FactDateOfExpiration,
		FactSex,
		FactNationality,
		FactCountryOfIssuance,
	},
	SourceUserSpecified: []string{
		FactDocumentNumber,
		FactDisplayName,
		FactEmailAddress,
		FactUnverifiedPhoneNumber,
		FactLocation,
	},
}
