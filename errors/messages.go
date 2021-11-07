package errorTypes

const (
	FetchClientError                = "When fething client"
	ConnectingDbError               = "Connecting Mongo Db"
	ProcessedRecordError            = "When getting processed record"
	RequestParameterError           = "Check your request parameters"
	RequestParameterMinCountError   = "MinCount parameter must be number"
	RequestParameterMaxCountError   = "MaxCount parameter must be number"
	RequestParameterDateTypeError   = "Dates are must be YYYY-MM-DD format"
	MinCountBiggerThanMaxCountError = "MinCount cannot bigger than max count"
	KeyWasNotFound                  = "Key was not found"
)
