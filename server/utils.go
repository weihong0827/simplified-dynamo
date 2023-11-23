package main

func IsKeyInRange(key uint32, start uint32, end uint32) bool {
	if start < end {
		return key >= start && key < end
	}
	return key >= start || key < end
}
