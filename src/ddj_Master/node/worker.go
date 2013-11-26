package node

// TODO: MOVE FROM HERE
/*
if r.Type == common.TASK_SELECT_ALL {
	length := int(r.LoadSize / 24)
	load := make([]dto.Dto, length)
	for i := 0; i < length; i++ {
		var e dto.Element
		err = e.Decode(buffer[i*(e.Size()+4):])
		if err != nil {
			log.Error(err)
			continue
		}
		load[i] = &e
	}
	r.Load = load
}
*/
