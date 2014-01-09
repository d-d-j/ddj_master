package common

//Aggregation Type
const (
	// Values aggregation
	AGGREGATION_NONE               int32 = 0
	AGGREGATION_ADD                int32 = 1
	AGGREGATION_MIN                int32 = 2
	AGGREGATION_MAX                int32 = 3
	AGGREGATION_AVERAGE            int32 = 4
	AGGREGATION_STDDEVIATION       int32 = 5
	AGGREGATION_VARIANCE           int32 = 6
	AGGREGATION_DIFFERENTIAL       int32 = 7
	AGGREGATION_INTEGRAL           int32 = 8
	AGGREGATION_HISTOGRAM_BY_VALUE int32 = 9
	AGGREGATION_HISTOGRAM_BY_TIME  int32 = 10
	AGGREGATION_SERIES_SUM         int32 = 13
)
