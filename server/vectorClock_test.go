package main

import (
	pb "dynamoSimplified/pb"
	"testing"
)

func TestCompareClockSame(t *testing.T) {
	var kv1 *pb.KeyValue
	var kv2 *pb.KeyValue

	kv1 = &pb.KeyValue{
		Key:   "foo",
		Value: "bar",
		VectorClock: &pb.VectorClock{
			Timestamps: map[uint32]*pb.ClockStruct{
				1: &pb.ClockStruct{
					ClokcVal: 1,
				},
			},
		},
	}

	kv2 = &pb.KeyValue{
		Key:   "foo",
		Value: "bar",
		VectorClock: &pb.VectorClock{
			Timestamps: map[uint32]*pb.ClockStruct{
				1: &pb.ClockStruct{
					ClokcVal: 1,
				},
			},
		},
	}

	lt, mt := compareClock(kv1, kv2)

	if lt != false {
		t.Errorf("Expected lt: %v, actual lt: %v", false, lt)
	}

	if mt != false {
		t.Errorf("Expected mt: %v, actual mt: %v", false, mt)
	}

}

func TestCompareClockConcurrent(t *testing.T) {
	var kv1 *pb.KeyValue
	var kv2 *pb.KeyValue

	kv1 = &pb.KeyValue{
		Key:   "foo",
		Value: "bar",
		VectorClock: &pb.VectorClock{
			Timestamps: map[uint32]*pb.ClockStruct{
				1: &pb.ClockStruct{
					ClokcVal: 1,
				},
				2: &pb.ClockStruct{
					ClokcVal: 2,
				},
			},
		},
	}

	kv2 = &pb.KeyValue{
		Key:   "foo",
		Value: "bar",
		VectorClock: &pb.VectorClock{
			Timestamps: map[uint32]*pb.ClockStruct{
				1: &pb.ClockStruct{
					ClokcVal: 2,
				},
				2: &pb.ClockStruct{
					ClokcVal: 1,
				},
			},
		},
	}

	lt, mt := compareClock(kv1, kv2)

	if lt != true {
		t.Errorf("Expected lt: %v, actual lt: %v", true, lt)
	}

	if mt != true {
		t.Errorf("Expected mt: %v, actual mt: %v", true, mt)
	}

}

func TestCompareClockBehind(t *testing.T) {
	var kv1 *pb.KeyValue
	var kv2 *pb.KeyValue

	kv1 = &pb.KeyValue{
		Key:   "foo",
		Value: "bar",
		VectorClock: &pb.VectorClock{
			Timestamps: map[uint32]*pb.ClockStruct{
				1: &pb.ClockStruct{
					ClokcVal: 1,
				},
				2: &pb.ClockStruct{
					ClokcVal: 1,
				},
			},
		},
	}

	kv2 = &pb.KeyValue{
		Key:   "foo",
		Value: "bar",
		VectorClock: &pb.VectorClock{
			Timestamps: map[uint32]*pb.ClockStruct{
				1: &pb.ClockStruct{
					ClokcVal: 2,
				},
				2: &pb.ClockStruct{
					ClokcVal: 1,
				},
			},
		},
	}

	lt, mt := compareClock(kv1, kv2)

	if lt != true {
		t.Errorf("Expected lt: %v, actual lt: %v", true, lt)
	}

	if mt != false {
		t.Errorf("Expected mt: %v, actual mt: %v", false, mt)
	}
}

func TestCompareClockAhead(t *testing.T) {
	var kv1 *pb.KeyValue
	var kv2 *pb.KeyValue

	kv1 = &pb.KeyValue{
		Key:   "foo",
		Value: "bar",
		VectorClock: &pb.VectorClock{
			Timestamps: map[uint32]*pb.ClockStruct{
				1: &pb.ClockStruct{
					ClokcVal: 2,
				},
				2: &pb.ClockStruct{
					ClokcVal: 2,
				},
			},
		},
	}

	kv2 = &pb.KeyValue{
		Key:   "foo",
		Value: "bar",
		VectorClock: &pb.VectorClock{
			Timestamps: map[uint32]*pb.ClockStruct{
				1: &pb.ClockStruct{
					ClokcVal: 2,
				},
				2: &pb.ClockStruct{
					ClokcVal: 1,
				},
			},
		},
	}

	lt, mt := compareClock(kv1, kv2)

	if lt != false {
		t.Errorf("Expected lt: %v, actual lt: %v", false, lt)
	}

	if mt != true {
		t.Errorf("Expected mt: %v, actual mt: %v", true, mt)
	}

}

func TestCompareClockMoreClock(t *testing.T) {
	var kv1 *pb.KeyValue
	var kv2 *pb.KeyValue

	kv1 = &pb.KeyValue{
		Key:   "foo",
		Value: "bar",
		VectorClock: &pb.VectorClock{
			Timestamps: map[uint32]*pb.ClockStruct{
				1: &pb.ClockStruct{
					ClokcVal: 2,
				},
				2: &pb.ClockStruct{
					ClokcVal: 2,
				},
			},
		},
	}

	kv2 = &pb.KeyValue{
		Key:   "foo",
		Value: "bar",
		VectorClock: &pb.VectorClock{
			Timestamps: map[uint32]*pb.ClockStruct{
				1: &pb.ClockStruct{
					ClokcVal: 2,
				},
			},
		},
	}

	lt, mt := compareClock(kv1, kv2)

	if lt != false {
		t.Errorf("Expected lt: %v, actual lt: %v", false, lt)
	}

	if mt != true {
		t.Errorf("Expected mt: %v, actual mt: %v", true, mt)
	}

}

func TestDeleteDataStart(t *testing.T) {
	var data []*pb.KeyValue
	var i int

	// the data has len of 3
	data = make([]*pb.KeyValue, 0)
	data = append(data, &pb.KeyValue{
		Key:   "foo",
		Value: "bar",
		VectorClock: &pb.VectorClock{
			Timestamps: map[uint32]*pb.ClockStruct{
				1: &pb.ClockStruct{
					ClokcVal: 1,
				},
				2: &pb.ClockStruct{
					ClokcVal: 2,
				},
			},
		},
	})
	data = append(data, &pb.KeyValue{
		Key:   "foo",
		Value: "bar",
		VectorClock: &pb.VectorClock{
			Timestamps: map[uint32]*pb.ClockStruct{
				1: &pb.ClockStruct{
					ClokcVal: 2,
				},
				2: &pb.ClockStruct{
					ClokcVal: 1,
				},
			},
		},
	})
	data = append(data, &pb.KeyValue{
		Key:   "foo",
		Value: "bar",
		VectorClock: &pb.VectorClock{
			Timestamps: map[uint32]*pb.ClockStruct{
				1: &pb.ClockStruct{
					ClokcVal: 2,
				},
				2: &pb.ClockStruct{
					ClokcVal: 2,
				},
			},
		},
	})

	i = 0

	data = delete_data(i, data)

	if len(data) != 2 {
		t.Errorf("Expected len: %v, actual len: %v", 2, len(data))
	}

	if data[0].VectorClock.Timestamps[1].ClokcVal != 2 {
		t.Errorf("Expected clock: %v, actual clock: %v", 2, data[0].VectorClock.Timestamps[1].ClokcVal)
	}

	if data[0].VectorClock.Timestamps[2].ClokcVal != 1 {
		t.Errorf("Expected clock: %v, actual clock: %v", 2, data[1].VectorClock.Timestamps[1].ClokcVal)
	}

}

func TestDeleteDataMiddle(t *testing.T) {
	var data []*pb.KeyValue
	var i int

	// the data has len of 3
	data = make([]*pb.KeyValue, 0)
	data = append(data, &pb.KeyValue{
		Key:   "foo",
		Value: "bar",
		VectorClock: &pb.VectorClock{
			Timestamps: map[uint32]*pb.ClockStruct{
				1: &pb.ClockStruct{
					ClokcVal: 1,
				},
				2: &pb.ClockStruct{
					ClokcVal: 2,
				},
			},
		},
	})
	data = append(data, &pb.KeyValue{
		Key:   "foo",
		Value: "bar",
		VectorClock: &pb.VectorClock{
			Timestamps: map[uint32]*pb.ClockStruct{
				1: &pb.ClockStruct{
					ClokcVal: 2,
				},
				2: &pb.ClockStruct{
					ClokcVal: 1,
				},
			},
		},
	})
	data = append(data, &pb.KeyValue{
		Key:   "foo",
		Value: "bar",
		VectorClock: &pb.VectorClock{
			Timestamps: map[uint32]*pb.ClockStruct{
				1: &pb.ClockStruct{
					ClokcVal: 2,
				},
				2: &pb.ClockStruct{
					ClokcVal: 2,
				},
			},
		},
	})

	i = 1

	data = delete_data(i, data)

	if len(data) != 2 {
		t.Errorf("Expected len: %v, actual len: %v", 2, len(data))
	}

	if data[1].VectorClock.Timestamps[1].ClokcVal != 2 {
		t.Errorf("Expected clock: %v, actual clock: %v", 2, data[0].VectorClock.Timestamps[1].ClokcVal)
	}

	if data[1].VectorClock.Timestamps[2].ClokcVal != 2 {
		t.Errorf("Expected clock: %v, actual clock: %v", 2, data[1].VectorClock.Timestamps[1].ClokcVal)
	}

}

func TestDeleteDataEnd(t *testing.T) {
	var data []*pb.KeyValue
	var i int

	// the data has len of 3
	data = make([]*pb.KeyValue, 0)
	data = append(data, &pb.KeyValue{
		Key:   "foo",
		Value: "bar",
		VectorClock: &pb.VectorClock{
			Timestamps: map[uint32]*pb.ClockStruct{
				1: &pb.ClockStruct{
					ClokcVal: 1,
				},
				2: &pb.ClockStruct{
					ClokcVal: 2,
				},
			},
		},
	})
	data = append(data, &pb.KeyValue{
		Key:   "foo",
		Value: "bar",
		VectorClock: &pb.VectorClock{
			Timestamps: map[uint32]*pb.ClockStruct{
				1: &pb.ClockStruct{
					ClokcVal: 2,
				},
				2: &pb.ClockStruct{
					ClokcVal: 1,
				},
			},
		},
	})
	data = append(data, &pb.KeyValue{
		Key:   "foo",
		Value: "bar",
		VectorClock: &pb.VectorClock{
			Timestamps: map[uint32]*pb.ClockStruct{
				1: &pb.ClockStruct{
					ClokcVal: 2,
				},
				2: &pb.ClockStruct{
					ClokcVal: 2,
				},
			},
		},
	})

	i = 2

	data = delete_data(i, data)

	if len(data) != 2 {
		t.Errorf("Expected len: %v, actual len: %v", 2, len(data))
	}

	if data[1].VectorClock.Timestamps[1].ClokcVal != 2 {
		t.Errorf("Expected clock: %v, actual clock: %v", 2, data[0].VectorClock.Timestamps[1].ClokcVal)
	}

	if data[1].VectorClock.Timestamps[2].ClokcVal != 1 {
		t.Errorf("Expected clock: %v, actual clock: %v", 2, data[1].VectorClock.Timestamps[1].ClokcVal)
	}

}

func TestCompareVectorClocksOneData(t *testing.T) {
	var data []*pb.KeyValue

	// the data has len of 1
	data = make([]*pb.KeyValue, 0)
	data = append(data, &pb.KeyValue{
		Key:   "foo",
		Value: "bar",
		VectorClock: &pb.VectorClock{
			Timestamps: map[uint32]*pb.ClockStruct{
				1: &pb.ClockStruct{
					ClokcVal: 1,
				},
				2: &pb.ClockStruct{
					ClokcVal: 1,
				},
			},
		},
	})

	newData := CompareVectorClocks(data)

	if len(newData) != 1 {
		t.Errorf("Expected len: %v, actual len: %v", 2, len(data))
	}

	if newData[0].VectorClock.Timestamps[1].ClokcVal != 1 {
		t.Errorf("Expected clock: %v, actual clock: %v", 1, newData[0].VectorClock.Timestamps[1].ClokcVal)
	}

	if newData[0].VectorClock.Timestamps[2].ClokcVal != 1 {
		t.Errorf("Expected clock: %v, actual clock: %v", 1, newData[0].VectorClock.Timestamps[2].ClokcVal)
	}

}

func TestCompareVectorClocksSame(t *testing.T) {
	var data []*pb.KeyValue

	// the data has len of 3
	data = make([]*pb.KeyValue, 0)
	data = append(data, &pb.KeyValue{
		Key:   "foo",
		Value: "bar",
		VectorClock: &pb.VectorClock{
			Timestamps: map[uint32]*pb.ClockStruct{
				1: &pb.ClockStruct{
					ClokcVal: 1,
				},
				2: &pb.ClockStruct{
					ClokcVal: 1,
				},
			},
		},
	})
	data = append(data, &pb.KeyValue{
		Key:   "foo",
		Value: "bar",
		VectorClock: &pb.VectorClock{
			Timestamps: map[uint32]*pb.ClockStruct{
				1: &pb.ClockStruct{
					ClokcVal: 1,
				},
				2: &pb.ClockStruct{
					ClokcVal: 1,
				},
			},
		},
	})
	data = append(data, &pb.KeyValue{
		Key:   "foo",
		Value: "bar",
		VectorClock: &pb.VectorClock{
			Timestamps: map[uint32]*pb.ClockStruct{
				1: &pb.ClockStruct{
					ClokcVal: 1,
				},
				2: &pb.ClockStruct{
					ClokcVal: 1,
				},
			},
		},
	})

	data = CompareVectorClocks(data)

	if len(data) != 1 {
		t.Errorf("Expected len: %v, actual len: %v", 1, len(data))
	}

	if data[0].VectorClock.Timestamps[1].ClokcVal != 1 {
		t.Errorf("Expected clock: %v, actual clock: %v", 1, data[0].VectorClock.Timestamps[1].ClokcVal)
	}

	if data[0].VectorClock.Timestamps[2].ClokcVal != 1 {
		t.Errorf("Expected clock: %v, actual clock: %v", 1, data[0].VectorClock.Timestamps[2].ClokcVal)
	}

}

func TestCompareVectorClocksConcurrent(t *testing.T) {
	var data []*pb.KeyValue

	// the data has len of 2
	data = make([]*pb.KeyValue, 0)
	data = append(data, &pb.KeyValue{
		Key:   "foo",
		Value: "bar",
		VectorClock: &pb.VectorClock{
			Timestamps: map[uint32]*pb.ClockStruct{
				1: &pb.ClockStruct{
					ClokcVal: 2,
				},
				2: &pb.ClockStruct{
					ClokcVal: 1,
				},
			},
		},
	})
	data = append(data, &pb.KeyValue{
		Key:   "foo",
		Value: "bar",
		VectorClock: &pb.VectorClock{
			Timestamps: map[uint32]*pb.ClockStruct{
				1: &pb.ClockStruct{
					ClokcVal: 1,
				},
				2: &pb.ClockStruct{
					ClokcVal: 2,
				},
			},
		},
	})

	newData := CompareVectorClocks(data)

	if len(newData) != 2 {
		t.Errorf("Expected len: %v, actual len: %v", 2, len(data))
	}
}

func TestCompareVectorClocksBehind(t *testing.T) {
	var data []*pb.KeyValue

	// the data has len of 2
	data = make([]*pb.KeyValue, 0)
	data = append(data, &pb.KeyValue{
		Key:   "foo",
		Value: "bar",
		VectorClock: &pb.VectorClock{
			Timestamps: map[uint32]*pb.ClockStruct{
				1: &pb.ClockStruct{
					ClokcVal: 1,
				},
				2: &pb.ClockStruct{
					ClokcVal: 1,
				},
			},
		},
	})
	data = append(data, &pb.KeyValue{
		Key:   "foo",
		Value: "bar",
		VectorClock: &pb.VectorClock{
			Timestamps: map[uint32]*pb.ClockStruct{
				1: &pb.ClockStruct{
					ClokcVal: 1,
				},
				2: &pb.ClockStruct{
					ClokcVal: 2,
				},
			},
		},
	})

	newData := CompareVectorClocks(data)

	if len(newData) != 1 {
		t.Errorf("Expected len: %v, actual len: %v", 1, len(data))
	}

	if newData[0].VectorClock.Timestamps[1].ClokcVal != 1 {
		t.Errorf("Expected clock: %v, actual clock: %v", 1, newData[0].VectorClock.Timestamps[1].ClokcVal)
	}

	if newData[0].VectorClock.Timestamps[2].ClokcVal != 2 {
		t.Errorf("Expected clock: %v, actual clock: %v", 2, newData[0].VectorClock.Timestamps[2].ClokcVal)
	}

}

func TestCompareVectorClocksAhead(t *testing.T) {
	var data []*pb.KeyValue

	// the data has len of 2
	data = make([]*pb.KeyValue, 0)
	data = append(data, &pb.KeyValue{
		Key:   "foo",
		Value: "bar",
		VectorClock: &pb.VectorClock{
			Timestamps: map[uint32]*pb.ClockStruct{
				1: &pb.ClockStruct{
					ClokcVal: 2,
				},
				2: &pb.ClockStruct{
					ClokcVal: 1,
				},
			},
		},
	})
	data = append(data, &pb.KeyValue{
		Key:   "foo",
		Value: "bar",
		VectorClock: &pb.VectorClock{
			Timestamps: map[uint32]*pb.ClockStruct{
				1: &pb.ClockStruct{
					ClokcVal: 1,
				},
				2: &pb.ClockStruct{
					ClokcVal: 1,
				},
			},
		},
	})

	newData := CompareVectorClocks(data)

	if len(newData) != 1 {
		t.Errorf("Expected len: %v, actual len: %v", 1, len(data))
	}

	if newData[0].VectorClock.Timestamps[1].ClokcVal != 2 {
		t.Errorf("Expected clock: %v, actual clock: %v", 2, newData[0].VectorClock.Timestamps[1].ClokcVal)
	}

	if newData[0].VectorClock.Timestamps[2].ClokcVal != 1 {
		t.Errorf("Expected clock: %v, actual clock: %v", 1, newData[0].VectorClock.Timestamps[2].ClokcVal)
	}

}

func TestCompareVectorClocksReadFailed(t *testing.T) {
	var data []*pb.KeyValue

	// the data has len of 2 with one read failed
	data = make([]*pb.KeyValue, 0)
	data = append(data, &pb.KeyValue{
		Key:   "foo",
		Value: "bar",
		VectorClock: &pb.VectorClock{
			Timestamps: map[uint32]*pb.ClockStruct{
				1: &pb.ClockStruct{
					ClokcVal: 2,
				},
				2: &pb.ClockStruct{
					ClokcVal: 1,
				},
			},
		},
	})
	data = append(data, &pb.KeyValue{
		Key:   "Read",
		Value: "Read Failed!",
	})

	newData := CompareVectorClocks(data)

	if len(newData) != 1 {
		t.Errorf("Expected len: %v, actual len: %v", 1, len(data))
	}

	if newData[0].VectorClock.Timestamps[1].ClokcVal != 2 {
		t.Errorf("Expected clock: %v, actual clock: %v", 2, newData[0].VectorClock.Timestamps[1].ClokcVal)
	}

	if newData[0].VectorClock.Timestamps[2].ClokcVal != 1 {
		t.Errorf("Expected clock: %v, actual clock: %v", 1, newData[0].VectorClock.Timestamps[2].ClokcVal)
	}

}
