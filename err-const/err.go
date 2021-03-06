package err

import "fmt"

var (
	fEf = fmt.Errorf
)

var (
	FOR_TEST = fEf("FOR TESTING")
	NO_ERROR = fEf("NO ERROR, FOR ERROR CHAN GO THROUGH")

	NOT_IMPLEMENTED = fEf("SRC NOT IMPLEMENTED")
	NOT_SUPPORTED   = fEf("NOT SUPPORTED")

	EXTERNAL           = fEf("EXTERNAL ERR")
	NET_TIMEOUT        = fEf("NET TIME OUT")
	NET_NO_RESPONSE    = fEf("NO NET RESPONSE")
	HTTP_REQBODY_EMPTY = fEf("EMPTY HTTP REQUEST BODY")

	INTERNAL          = fEf("INTERNAL ERR")
	INTERNAL_SCAN_ERR = fEf("INTERNAL ERR - SCAN")
	INTERNAL_INIT_ERR = fEf("INTERNAL ERR - INIT")
	INTERNAL_DEADLOCK = fEf("INTERNAL ERR - DEADLOCK")

	CLI_ARG_ERR        = fEf("CLI ARGUMENT ERR")
	CLI_FLAG_ERR       = fEf("CLI FLAG ERR")
	CLI_SUBCMD_ERR     = fEf("CLI SUB COMMAND ERR")
	CLI_SUBCMD_UNKNOWN = fEf("CLI SUB COMMAND UNKNOWN")

	CFG_INIT_ERR     = fEf("CONFIG INIT ERR")
	CFG_SIGN_MISSING = fEf("CONFIG KEY SIGN MISSING")
	SRC_SIGN_MISSING = fEf("SOURCE FILE KEY SIGN MISSING")

	PARAM_INVALID            = fEf("INVALID PARAM(S)")
	PARAM_NOT_SUPPORTED      = fEf("NOT SUPPORTED PARAM")
	PARAM_INVALID_PTR        = fEf("INVALID POINTER PARAM")
	PARAM_INVALID_STRUCT     = fEf("INVALID STRUCT PARAM")
	PARAM_INVALID_STRUCT_PTR = fEf("INVALID STRUCT POINTER PARAM")
	PARAM_INVALID_SLICE      = fEf("INVALID SLICE PARAM")
	PARAM_INVALID_INDEX      = fEf("INVALID INDEX PARAM")
	PARAM_INVALID_MAP        = fEf("INVALID MAP PARAM")
	PARAM_INVALID_JSON       = fEf("INVALID JSON PARAM")
	PARAM_INVALID_XML        = fEf("INVALID XML PARAM")
	PARAM_INVALID_CSV        = fEf("INVALID CSV PARAM")
	PARAM_INVALID_FMT        = fEf("INVALID PARAM FORMAT")

	VAR_INVALID    = fEf("INVALID VARIABLE(S)")
	DIR_NOT_FOUND  = fEf("DIRECTORY NOT FOUND")
	FILE_NOT_FOUND = fEf("FILE NOT FOUND")
	FILE_EMPTY     = fEf("FILE EMPTY OR READ ERR")

	XML_INVALID             = fEf("INVALID XML")
	XML_NOT_FMT             = fEf("XML NOT FORMATTED")
	JSON_INVALID            = fEf("INVALID JSON")
	JSON_NOT_FMT            = fEf("JSON NOT FORMATTED")
	JSON_ARRAY_INVALID      = fEf("INVALID JSON ARRAY")
	JSON_ARRAY_NOT_FMT      = fEf("JSON ARRAY NOT FORMATTED")
	CSV_COLUMN_HEADER_EMPTY = fEf("CSV EMPTY HEADER")
	UUID_INVALID            = fEf("NOT STANDARD UUID (8-4-4-4-12)")

	NUM_MUST_POS          = fEf("NUMBER MUST BE POSITIVE")
	NUM_MUST_NOT_POS      = fEf("NUMBER MUST NOT BE POSITIVE")
	NUM_MUST_NEG          = fEf("NUMBER MUST BE NEGATIVE")
	NUM_MUST_NOT_NEG      = fEf("NUMBER MUST NOT BE NEGATIVE")
	STR_EMPTY             = fEf("EMPTY STRING")
	STR_BLANK             = fEf("BLANK STRING")
	MAP_INVALID           = fEf("NOT A MAP")
	MAPS_DIF_KEY_TYPE     = fEf("DIFFERENT KEY TYPE")
	MAPS_DIF_VALUE_TYPE   = fEf("DIFFERENT VALUE TYPE")
	SLICE_INVALID         = fEf("NOT A SLICE")
	SLICE_INCORRECT_COUNT = fEf("INCORRECT ELEMENT COUNT")
	SLICES_DIF_LEN        = fEf("DIFFERENT LENGTH")
)
