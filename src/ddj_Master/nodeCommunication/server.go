package nodeCommunication


func IOHandler(Query <-chan dto.Query, Result <-chan dto.Result, NodeList *list.List) {
	taskResponse := make(map[int32]chan []dto.Dto)
	for {
		select {
		case query := <-Query:
			log.Debug("Query", query)
			header := query.TaskRequestHeader

			var (
				buf       []byte
				headerBuf []byte
				err       error
			)

			if query.Load != nil {
				buf, err = query.Load.Encode()
				if err != nil {
					log.Error(err)
					continue
				}
			}

			header.LoadSize = (int32)(len(buf))
			headerBuf, err = header.Encode()
			if err != nil {
				log.Error(err)
				continue
			}

			complete := make([]byte, 100)
			copy(complete, headerBuf)
			copy(complete[len(headerBuf):], buf)

			//TODO: Replace this with StoreManager
			if NodeList.Len() == 0 {
				log.Warn("No node connected")
				query.Response <- nil
			} else if query.Code == constants.TASK_SELECT_ALL {
				taskResponse[query.Id] = query.Response
			} else {
				response := make([]dto.Dto, 0)
				query.Response <- response
			}

			for e := NodeList.Front(); e != nil; e = e.Next() {
				Node := e.Value.(Node)
				Node.Incoming <- complete
			}
		case result := <-Result:
			log.Fine("Result: ", result.String(), result.Load)
			ch := taskResponse[result.Id]
			log.Debug("Response channel", ch)
			if ch != nil {
				log.Debug("Pass result data to proper client")
				ch <- result.Load
				delete(taskResponse, result.Id)
			}
		}

	}
}

// Node reading goroutine - reads incoming data from the tcp socket,
// sends it to the Node.Outgoing channel (to be picked up by IOHandler)
func NodeReader(Node *Node) {

	var r dto.Result
	buffer := make([]byte, r.TaskRequestHeader.Size())
	for Node.Read(buffer) {
		log.Debug("NodeReader received data from", Node.Id)
		err := r.DecodeHeader(buffer)
		if err != nil {
			log.Error(err)
		}
		log.Fine("Response header: ", r.TaskRequestHeader)
		if r.LoadSize == 0 {
			r.Load = make([]dto.Dto, 0)
			Node.Outgoing <- r
			continue
		}
		buffer := make([]byte, r.LoadSize)
		Node.Read(buffer)
		if r.Code == constants.TASK_SELECT_ALL {
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
		log.Debug("Send response to IOHandler")
		Node.Outgoing <- r
		buffer = make([]byte, r.TaskRequestHeader.Size())
	}

	log.Info("NodeReader stopped for ", Node.Id)
}

// Node sending goroutine - waits for data to be sent over Node.Incoming
// (from IOHandler), then sends it over the socket
func NodeSender(Node *Node) {
	for {
		select {
		case buffer := <-Node.Incoming:
			log.Debug("NodeSender sending ", buffer, " to ", Node.Id)
			Node.Conn.Write(buffer)
		case <-Node.Quit:
			log.Info("Node ", Node.Id, " quitting")
			Node.Conn.Close()
			break
		}
	}
}
