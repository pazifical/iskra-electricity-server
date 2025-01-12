package sml

import (
	"io/ioutil"
	"strconv"
	"strings"
	"testing"
)

func importReadoutWithPin() string {
	filePath := "./testdata/readout_with_pin.txt"
	readout, err := ioutil.ReadFile(filePath)

	if err != nil {
		panic(err)
	}

	return string(readout)
}

// TODO: Implement Tests
func TestParseSMLCorrectly(t *testing.T) {
	// Arrange
	readout := importReadoutWithPin()
	slicedReadout := sliceReadout(readout)
	chunkIndexes := getChunkStartingIndexes(slicedReadout)
	decimals, _ := processChunk(slicedReadout[chunkIndexes[0]:chunkIndexes[1]])

	// Act
	got, _ := extractUsageFromArray(decimals)
	want := 7408.2524

	// Assert
	if got != want {
		t.Errorf("got %f; want %f", got, want)
	}
}

func TestChunkIndexesWithRealData(t *testing.T) {
	// Arrange
	readout := importReadoutWithPin()
	slicedReadout := sliceReadout(readout)

	// Act
	got := getChunkStartingIndexes(slicedReadout)
	want0 := 0
	want1 := 384
	// fmt.Printf("\nChunk 1: %v", slicedReadout[got[0]:got[1]])
	// fmt.Printf("\nChunk 2: %v", slicedReadout[got[1]:got[2]])

	// Assert
	if got[0] != want0 {
		t.Errorf("got %d; want %d", got[0], want0)
	}
	if got[1] != want1 {
		t.Errorf("got %d; want %d", got[1], want1)
	}
}

func TestSliceReadout(t *testing.T) {
	readout := "asd98 ab cd\nef e10000 a1"

	got := sliceReadout(readout)
	want := []string{"ab", "cd", "ef", "a1"}

	if len(got) != len(want) {
		t.Errorf("got %d; want %d", len(got), len(want))
	}
	for i := range got {
		if got[i] != want[i] {
			t.Errorf("got %s; want %s", got[i], want[i])
		}
	}
}

func TestGetChunkStartingIndexes(t *testing.T) {
	slicedReadout := strings.Split("1b 1b 1b 1b 01 01 01 01 00 00 1b 1b 1b 1b 01 01 01 01", " ")

	got := getChunkStartingIndexes(slicedReadout)
	want := []int{0, 10}

	if len(got) != len(want) {
		t.Errorf("got %d; want %d", len(got), len(want))
	}
	for i := range got {
		if got[i] != want[i] {
			t.Errorf("got %d; want %d", got[i], want[i])
		}
	}
}

// TODO: Implement Test
func TestProcessChunk(t *testing.T) {
	// 	// Arrange
	// 	readout := importReadoutWithPin()
	// 	sliced := sliceReadout(readout)
	// 	chunkIndexes := getChunkStartingIndexes(sliced)
	// 	chunk := sliced[chunkIndexes[0]:chunkIndexes[1]]

	// 	// Act
	// 	_ = processChunk(chunk)

	// 	// Assert
	// 	if chunk[0] != "1b" {
	// 		t.Errorf("got %s; want %s", chunk[0], "1b")
	// 	}
}

func TestHexRuneToInt(t *testing.T) {
	// Act
	want := 0
	got, _ := hexToInt("0")
	// Assert
	if got != want {
		t.Errorf("got %d; want %d", got, want)
	}

	// Act
	want = 11
	got, _ = hexToInt("b")
	// Assert
	if got != want {
		t.Errorf("got %d; want %d", got, want)
	}

	// Act
	want = 15
	got, _ = hexToInt("f")
	// Assert
	if got != want {
		t.Errorf("got %d; want %d", got, want)
	}

	// Act
	_, err := hexToInt("g")
	// Assert
	if err == nil {
		t.Errorf("hexVal should fail for g")
	}
}

func TestHandleIntegerHexBlock630102(t *testing.T) {
	// Arrange
	chunk := []string{
		"63",
		"01",
		"02",
	}
	hexVal := chunk[0]

	// Act
	wantDecimal := 258
	wantLength := 3
	gotDecimal, gotLenth, ok := handleIntegerHexBlock(hexVal, chunk)

	// Assert
	if !ok {
		t.Errorf("TestHandleIntegerHexBlock should return ok.")
	}
	if wantLength != gotLenth {
		t.Errorf("got %d; want %d", gotLenth, wantLength)
	}
	if wantDecimal != int(gotDecimal) {
		t.Errorf("got %d; want %d", gotDecimal, wantDecimal)
	}
}

func TestHandleIntegerHexBlock530102(t *testing.T) {
	// Arrange
	chunk := []string{
		"53",
		"01",
		"02",
	}
	hexVal := chunk[0]

	// Act
	wantDecimal := 258
	wantLength := 3
	gotDecimal, gotLenth, ok := handleIntegerHexBlock(hexVal, chunk)

	// Assert
	if !ok {
		t.Errorf("TestHandleIntegerHexBlock should return ok.")
	}
	if wantLength != gotLenth {
		t.Errorf("got %d; want %d", gotLenth, wantLength)
	}
	if wantDecimal != int(gotDecimal) {
		t.Errorf("got %d; want %d", gotDecimal, wantDecimal)
	}
}

func TestHandleIntegerHexBlockFailedDecoding(t *testing.T) {
	// Arrange
	chunk := []string{
		"03",
		"zz",
		"zz",
	}
	hexVal := chunk[0]

	// Act
	_, gotLength, ok := handleIntegerHexBlock(hexVal, chunk)
	wantLength := 3

	// Assert
	if ok {
		t.Errorf("TestHandleIntegerHexBlock with %v should not return ok.", chunk)
	}
	if wantLength != gotLength {
		t.Errorf("got %d; want %d", gotLength, wantLength)
	}
}

func TestHandleIntegerHexBlockFailedLength(t *testing.T) {
	// Arrange
	chunk := []string{
		"0z",
		"01",
		"02",
	}
	hexVal := chunk[0]

	// Act
	_, gotLenth, ok := handleIntegerHexBlock(hexVal, chunk)
	wantLength := 0

	// Assert
	if ok {
		t.Errorf("TestHandleIntegerIntegerBlock should not return ok.")
	}
	if wantLength != gotLenth {
		t.Errorf("got %d; want %d", gotLenth, wantLength)
	}
}

func TestHandleOctatHexBlock030102(t *testing.T) {
	// Arrange
	chunk := []string{
		"03",
		"01",
		"02",
	}
	hexVal := chunk[0]

	// Act
	wantLength := 3
	gotLenth, ok := handleOctatHexBlock(hexVal, chunk)

	// Assert
	if !ok {
		t.Errorf("TestHandleIntegerOctatBlock should return ok.")
	}
	if wantLength != gotLenth {
		t.Errorf("got %d; want %d", gotLenth, wantLength)
	}
}

func TestHandleOctatHexBlockFailedDecode(t *testing.T) {
	// Arrange
	chunk := []string{
		"03",
		"zz",
		"zz",
	}
	hexVal := chunk[0]

	// Act
	gotLenth, ok := handleOctatHexBlock(hexVal, chunk)
	wantLength := 3

	// Assert
	if ok {
		t.Errorf("TestHandleIntegerOctatBlock should not return ok.")
	}
	if wantLength != gotLenth {
		t.Errorf("got %d; want %d", gotLenth, wantLength)
	}
}

func TestHandleOctatHexBlockFailedLength(t *testing.T) {
	// Arrange
	chunk := []string{
		"0z",
		"01",
		"02",
	}
	hexVal := chunk[0]

	// Act
	gotLenth, ok := handleOctatHexBlock(hexVal, chunk)
	wantLength := 0

	// Assert
	if ok {
		t.Errorf("TestHandleIntegerOctatBlock should not return ok.")
	}
	if wantLength != gotLenth {
		t.Errorf("got %d; want %d", gotLenth, wantLength)
	}
}

func TestHexToDecimalSuccess(t *testing.T) {
	hex := "1b ff"

	got, _ := hexToDecimal(hex)
	var want int64 = 7167

	if got != want {
		t.Errorf("got != want; got %d, want %d", got, want)
	}
}

func TestHexToDecimalSuccessBlanks(t *testing.T) {
	hex := "1b    Ff"

	got, _ := hexToDecimal(hex)
	var want int64 = 7167

	if got != want {
		t.Errorf("got != want; got %d, want %d", got, want)
	}
}

func TestHexToDecimalFails(t *testing.T) {
	hex := "1b    Fz"

	_, ok := hexToDecimal(hex)

	if ok {
		t.Errorf("Should not be ok")
	}
}

func TestTest(t *testing.T) {
	hex := "a"
	want := 10
	got, err := strconv.ParseInt(hex, 16, 64)
	if err != nil {
		t.Errorf(err.Error())
	}
	if int(got) != want {
		t.Errorf("%d != %d", got, want)
	}

}
