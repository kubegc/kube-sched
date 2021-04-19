package algorithm


var SingleAlgorithmBuilder = map[string]SingleScheduleAlgorithm{}
var BatchAlgorithmBuilder = map[string]BatchScheduleAlgorithm{}
func init() {
	RegisterBatch(&BatchFairScheduleAlgorithm{})
	RegisterBatch(&BatchTradeScheduleAlgorithm{})
	RegisterBatch(&BatchThroughputScheduleAlgorithm{})
}

func RegisterSingle(algorithm SingleScheduleAlgorithm) {
	SingleAlgorithmBuilder[algorithm.Name()] = algorithm
}

func RegisterBatch(algorithm BatchScheduleAlgorithm) {
	BatchAlgorithmBuilder[algorithm.Name()] = algorithm
}

func GetSingleScheduleAlgorithm(name string) SingleScheduleAlgorithm {
	return SingleAlgorithmBuilder[name]
}
func GetBatchScheduleAlgorithm(name string) BatchScheduleAlgorithm {
	return BatchAlgorithmBuilder[name]
}






