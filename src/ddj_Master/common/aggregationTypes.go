package common

//Aggregation Type
const (
	// Values aggregation
	AGGREGATION_NONE         int32 = 0
	AGGREGATION_ADD          int32 = 1
	AGGREGATION_MIN          int32 = 2
	AGGREGATION_MAX          int32 = 3
	AGGREGATION_AVERAGE      int32 = 4
	AGGREGATION_STDDEVIATION int32 = 5
	AGGREGATION_COUNT        int32 = 6
	AGGREGATION_VARIANCE     int32 = 7
	AGGREGATION_DIFFERENTIAL int32 = 8
	AGGREGATION_INTEGRAL     int32 = 9
	// Series aggregation
	AGGREGATION_NONE_SERIES         int32 = 10
	AGGREGATION_ADD_SERIES          int32 = 11
	AGGREGATION_MIN_SERIES          int32 = 12
	AGGREGATION_MAX_SERIES          int32 = 13
	AGGREGATION_AVERAGE_SERIES      int32 = 14
	AGGREGATION_STDDEVIATION_SERIES int32 = 15
	AGGREGATION_COUNT_SERIES        int32 = 16
	AGGREGATION_VARIANCE_SERIES     int32 = 17
	AGGREGATION_DIFFERENTIAL_SERIES int32 = 18
	AGGREGATION_INTEGRAL_SERIES     int32 = 19
)