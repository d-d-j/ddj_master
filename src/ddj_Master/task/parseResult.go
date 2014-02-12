package task

import (
	log "code.google.com/p/log4go"
	"ddj_Master/common"
	"ddj_Master/dto"
	"ddj_Master/reduce"
)

func parseResults(results []*dto.Result, query *dto.Query) []reduce.Aggregates {
	aggregationType := query.AggregationType
	switch aggregationType {
	case common.AGGREGATION_ADD:
		var value dto.Value
		return parseResultsUsingCreator(results, value)
	case common.AGGREGATION_AVERAGE:
		return parseResultsUsingCreator(results, dto.AverageElement{})
	case common.AGGREGATION_STDDEVIATION:
		return parseResultsUsingCreator(results, dto.VarianceElement{})
	case common.AGGREGATION_VARIANCE:
		return parseResultsUsingCreator(results, dto.VarianceElement{})
	case common.AGGREGATION_INTEGRAL:
		return parseResultsUsingCreator(results, dto.IntegralElement{})
	case common.AGGREGATION_HISTOGRAM_BY_TIME:
		return parseResultsToHistograms(results, int(query.AdditionalData.GetBucketCount()))
	case common.AGGREGATION_HISTOGRAM_BY_VALUE:
		return parseResultsToHistograms(results, int(query.AdditionalData.GetBucketCount()))
	case common.AGGREGATION_SERIES_SUM:
		return parseResultsToInterpolateElement(results, int(query.AdditionalData.GetBucketCount()))
	default:
		return parseResultsUsingCreator(results, dto.Element{})
	}
}

func parseResultsUsingCreator(results []*dto.Result, creator dto.Creator) []reduce.Aggregates {
	elementSize := creator.Size()
	resultsCount := len(results)
	elements := make([]reduce.Aggregates, 0, resultsCount)

	for _, result := range results {
		length := len(result.Data) / elementSize
		for j := 0; j < length; j++ {
			e, err := creator.Create(result.Data[j*elementSize:])
			if err != nil {
				log.Error("Problem with parsing data", err)
				continue
			}
			elements = append(elements, e)
		}
	}
	return elements
}

func parseResultsToHistograms(results []*dto.Result, bucketCount int) []reduce.Aggregates {

	histograms := make([]reduce.Aggregates, 0, len(results)*bucketCount)

	for _, result := range results {

		var combinedHistogramFromOneNodeButAllGpus dto.Histogram
		err := combinedHistogramFromOneNodeButAllGpus.Decode(result.Data)
		if err != nil {
			log.Error("Problem with parsing data", err)
			continue
		}

		for index := 0; index < len(combinedHistogramFromOneNodeButAllGpus.Data); index += bucketCount {
			var histogramForOneGpu dto.Histogram
			histogramForOneGpu.Data = combinedHistogramFromOneNodeButAllGpus.Data[index : index+bucketCount]
			histograms = append(histograms, &histogramForOneGpu)
		}
	}

	return histograms
}

func parseResultsToInterpolateElement(results []*dto.Result, samplesCount int) []reduce.Aggregates {

	histograms := make([]reduce.Aggregates, 0, len(results)*samplesCount)

	for _, result := range results {

		var combinedInterpolateElementFromOneNodeButAllGpus dto.InterpolateElement
		err := combinedInterpolateElementFromOneNodeButAllGpus.Decode(result.Data)
		if err != nil {
			log.Error("Problem with parsing data", err)
			continue
		}

		for index := 0; index < len(combinedInterpolateElementFromOneNodeButAllGpus.Data); index += samplesCount {
			var histogramForOneGpu dto.InterpolateElement
			histogramForOneGpu.Data = combinedInterpolateElementFromOneNodeButAllGpus.Data[index : index+samplesCount]
			histograms = append(histograms, &histogramForOneGpu)
		}
	}

	return histograms
}
