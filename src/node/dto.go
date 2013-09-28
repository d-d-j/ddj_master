package node

type Status struct {
	totalRam, usedRam int32
	gpuTemperature    int32
}

type Element struct {
	series, tag int32
	time        int64
	value       float32
}
