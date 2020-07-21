// Copyright (c) 2020, RetailNext, Inc.
// This material contains trade secrets and confidential information of
// RetailNext, Inc.  Any use, reproduction, disclosure or dissemination
// is strictly prohibited without the explicit written permission
// of RetailNext, Inc.
// All rights reserved.

package easypost

const (
	AddressParametersInvalidCharacter ErrorCode = "ADDRESS.PARAMETERS.INVALID_CHARACTER" // The parameters passed contained an invalid character
	AddressParametersInvalid ErrorCode = "ADDRESS.PARAMETERS.INVALID" // The parameters passed to create an Address were missing or invalid
	AddressCountryInvalid ErrorCode = "ADDRESS.COUNTRY.INVALID" // Invalid 'country', please provide a 2 character ISO country code
	AddressNotFound ErrorCode = "ADDRESS.VERIFICATION.NOT_FOUND" // Address Not Found
	AddressVerificationFailure ErrorCode = "ADDRESS.VERIFICATION.FAILURE" // The address was unable to be verified
	AddressVerifyUnavailable ErrorCode = "ADDRESS.VERIFY.UNAVAILABLE" // Address verification is not available. Please try again.
	AddressVerificationInvalid ErrorCode = "ADDRESS.VERIFICATION.INVALID" // One of the verifications selected is invalid.
	AddressVerifyFailure ErrorCode = "ADDRESS.VERIFY.FAILURE" // Unable to verify address.
	AddressVerifyCarrierInvalid ErrorCode = "ADDRESS.VERIFY.CARRIER_INVALID" // Unable to verify address using provided carrier.
	AddressVerifyUpstreamUnavailable ErrorCode = "ADDRESS.VERIFY.UPSTREAM_UNAVAILABLE" // Address verification is not available due to an upstream service not responding. Please try again.
	AddressVerifyOnlyUS ErrorCode = "ADDRESS.VERIFY.ONLY_US" // USPS can only validate US addresses.
	AddressVerifyInternationalNotEnabled ErrorCode = "ADDRESS.VERIFY.INTL_NOT_ENABLED" // International Verification not enabled on this account. Please contact support@easypost.com.
	AddressVerifyMissingStreet ErrorCode = "ADDRESS.VERIFY.MISSING_STREET" // Insufficient address data provided. A street must be provided.
	AddressVerifyMissingCityStateZip ErrorCode = "ADDRESS.VERIFY.MISSING_CITY_STATE_ZIP" // Insufficient address data provided. A city and state or a zip must be provided.

	AddressVerificationCountryUnsupported ErrorCode = "E.COUNTRY.UNSUPPORTED" // This country is not supported in this mode. Try in production mode.
	AddressVerificationEngineUnavailable ErrorCode = "E.ENGINE.UNAVAILABLE" // No verification engine is available to service this country. Please try again later.
	AddressVerificationQueryUnanswerable ErrorCode = "E.QUERY.UNANSWERABLE" // We can not provide enough data in this country to satisfy this request. Please try again later.
	AddressVerificationNotfound ErrorCode = "E.ADDRESS.NOT_FOUND" // Address not found.
	AddressVerificationSecondaryInformationInvalid ErrorCode = "E.SECONDARY_INFORMATION.INVALID" // Invalid secondary information(Apt/Suite#). Please note, this can show up as an error on a successful request.
	AddressVerificationSecondaryInformationMissing ErrorCode= "E.SECONDARY_INFORMATION.MISSING" // Missing secondary information(Apt/Suite#). Please note, this can show up as an error on a successful request.
	AddressVerificationHouseNumberMissing ErrorCode= "E.HOUSE_NUMBER.MISSING" // House number is missing.
	AddressVerificationHouseNumberInvalid ErrorCode= "E.HOUSE_NUMBER.INVALID" // House number is invalid.
	AddressVerificationStreetMissing ErrorCode = "E.STREET.MISSING" // Street is missing.
	AddressVerificationStreetInvalid ErrorCode = "E.STREET.INVALID" // Street is invalid.
	AddressVerificationBoxNumberMissing ErrorCode= "E.BOX_NUMBER.MISSING" // Box number is missing.
	AddressVerificationBoxNumberInvalid ErrorCode= "E.BOX_NUMBER.INVALID" // Box number is invalid.
	AddressVerificationAddressInvalid ErrorCode= "E.ADDRESS.INVALID" // Invalid city/state/ZIP.
	AddressVerificationZipNotFound ErrorCode= "E.ZIP.NOT_FOUND" // Zip not found.
	AddressVerificationZipInvalid ErrorCode= "E.ZIP.INVALID" // Zip invalid.
	AddressVerificationZip4NotFound ErrorCode= "E.ZIP.PLUS4.NOT_FOUND" // Zip + 4 not found.
	AddressVerificationAddressMultiple ErrorCode = "E.ADDRESS.MULTIPLE" // Multiple addresses were returned with the same zip.
	AddressVerificationAddressInsufficient ErrorCode = "E.ADDRESS.INSUFFICIENT" // Insufficient/incorrect address data.
	AddressVerificationAddressDual ErrorCode = "E.ADDRESS.DUAL" // Dual address.
	AddressVerificationStreetMagnet ErrorCode = "E.STREET.MAGNET" // Multiple response due to magnet street syndrome.
	AddressVerificationCityStateInvalid ErrorCode = "E.CITY_STATE.INVALID" // Unverifiable city / state.
	AddressVerificationStateInvalid ErrorCode = "E.STATE.INVALID" // Invalid State.
	AddressVerificationDeliveryInvalid ErrorCode = "E.ADDRESS.DELIVERY.INVALID" // Invalid delivery address.
	AddressVerificationTimedOut ErrorCode = "E.TIMED_OUT" // Exceeded max timeout.
	AddressVerificationTimeZoneUnavailable ErrorCode = "E.TIME_ZONE.UNAVAILABLE" // The time zone service is currently unavailable.
	AddressVerificationPOBoxInternational ErrorCode = "E.PO_BOX.INTERNATIONAL" // Cannot verify international PO Box. Please note, this can show up as an error on a successful request.
)
