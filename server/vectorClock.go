package main

import (
	pb "dynamoSimplified/pb"
)

// Find out the concurrent key value pair from the vector clocks
// Return the concurrent key value pair, or value array
func CompareVectorClocks(data []*pb.KeyValue) []*pb.KeyValue {
	if len(data) == 1 {
		return data
	}

	var current *pb.KeyValue
	i := 0

	for i < len(data) {
		if data[i].Value == "Read Failed!" {
			data = delete_data(i, data)
			continue
		}

		if current == nil {
			current = data[i]
			data = delete_data(i, data)
			continue
		}

		lt, mt := compareClock(current, data[i])

		// if lt and mt are true, then the two clocks are concurrent
		// if lt is true, then the current clock is behind the data[i] clock
		// if mt is true, then the current clock is ahead of the data[i] clock
		// if lt and mt are false, then the two clocks are equal
		if lt && mt {
			i++
		} else if lt {
			current = data[i]
			data = delete_data(i, data)
		} else if mt {
			data = delete_data(i, data)
		} else {
			data = delete_data(i, data)
		}
	}

	data = append(data, current)
	return data
}

func compareClock(KV1 *pb.KeyValue, KV2 *pb.KeyValue) (bool, bool) {
	lt := false                                            //check if there are clocks1 less than clocks2
	mt := false                                            //check if there are clocks1 greater than clocks2
	for node, clock1 := range KV1.VectorClock.Timestamps { //extract all nodes in compare clock
		clock2, found := KV2.VectorClock.Timestamps[node] //check against the clocks in the next machine
		if !found {
			clock2 = &pb.ClockStruct{
				ClokcVal: 0,
			}
		}
		if clock1.ClokcVal < clock2.ClokcVal {
			lt = true
		} else if clock1.ClokcVal > clock2.ClokcVal {
			mt = true
		}
	}
	return lt, mt
}

func delete_data(i int, data []*pb.KeyValue) []*pb.KeyValue {
	newData := make([]*pb.KeyValue, 0)
	if i == 0 {
		newData = data[1:]
	} else if i == len(data)-1 {
		newData = data[:len(data)-1]
	} else {
		newData = append(data[:i], data[i+1:]...)
	}
	return newData
}
