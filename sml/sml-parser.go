package sml

import (
	"encoding/hex"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/pazifical/iskra-electricity-server/internal/logging"
	"github.com/pazifical/iskra-electricity-server/internal/types"
)

func ParseSensorReadout(readout string) (types.EnergyReading, error) {
	slicedReadout := sliceReadout(readout)
	chunkIndexes := getChunkStartingIndexes(slicedReadout)

	if len(chunkIndexes) < 2 {
		return types.EnergyReading{}, fmt.Errorf("parsing sml: len(chunkIndexes) < 2")
	}

	decimals, ok := processChunk(slicedReadout[chunkIndexes[0]:chunkIndexes[1]])
	if !ok {
		logging.Error(fmt.Sprintf("ERROR: ParseSML failed for readout: %v", slicedReadout))
	}

	kwhValue, err := extractUsageFromArray(decimals)
	if err != nil {
		return types.EnergyReading{}, fmt.Errorf("parsing sml: %w", err)
	}

	return types.EnergyReading{
		Type:  types.Electricity,
		Time:  time.Now(),
		Unit:  "kWh",
		Value: kwhValue,
	}, nil
}

// Removing non hex values and line breaks to create a slice of hex values
// TODO: use a regex to do the same
func sliceReadout(readout string) []string {
	readout = strings.Replace(readout, "\n", " ", -1)
	splitted := strings.Split(readout, " ")
	cleaned := []string{}
	for _, s := range splitted {
		if len(s) == 2 {
			cleaned = append(cleaned, s)
		}
	}
	return cleaned
}

func getChunkStartingIndexes(slicedReadout []string) []int {
	chunkStart := "1b 1b 1b 1b 01 01 01 01"
	var chunkStartIndexes []int

	for i := 0; i < len(slicedReadout)-7; i++ {
		joined := strings.Join(slicedReadout[i:(i+8)], " ")
		if joined == chunkStart {
			chunkStartIndexes = append(chunkStartIndexes, i)
		}
	}
	return chunkStartIndexes
}

func processChunk(chunk []string) ([]int64, bool) {
	var decimals []int64

	strippedChunk := chunk[8:(len(chunk) - 8)]

	var i int
	breakOut := false

	for i < len(strippedChunk) {
		// Stop the processing under specific conditions
		if breakOut {
			break
		}

		hexVal := strippedChunk[i]

		if hexVal == "00" { // TODO: Not yet sure why that is. End of a block?
			log.Printf("WARNING: Encountered 00 hex value. End of a block?")
			i += 1
			continue
		}

		switch rune(hexVal[0]) {
		case '0': // Octat
			listLength, ok := handleOctatHexBlock(hexVal, strippedChunk[i:])
			if ok {
				i += listLength
			} else {
				log.Printf("WARNING: Cannot handle octat hex block: %s", hexVal)
				breakOut = true
			}

		case '4': // Boolean. Not sure if and how that could be handled
			continue

		case '5': // Integer (handled the same as an unsigned integer)
			fallthrough

		case '6': // Unsigned integer
			decimal, listLength, ok := handleIntegerHexBlock(hexVal, strippedChunk[i:])
			if ok {
				decimals = append(decimals, decimal)
				i += listLength
			} else {
				log.Printf("WARNING: Cannot handle (unsigned) integer hex block: %s", hexVal)
				breakOut = true
			}

		case '7':
			_, err := hexToInt(string(hexVal[1]))
			if err == nil {
				// TODO: Maybe read the block names?
				i += 1
			} else {
				log.Printf("ERROR: _ Length problem: %s", hexVal)
				breakOut = true
			}

		case '8': // TODO: Parse the block after 8x
			log.Printf("WARNING: Hex starting value of 8 encountered. Cannot be parsed yet and will be skipped.")
			breakOut = true

		default:
			breakOut = true
		}
	}
	return decimals, true
}

// Processing a block of hex values that contain an (unsigned) Integer value
func handleIntegerHexBlock(hexVal string, chunk []string) (int64, int, bool) {
	listLength, err := hexToInt(string(hexVal[1]))
	if err != nil {
		logging.Error(fmt.Sprintf("length problem: %s : %v", hexVal, err))
		return 0, 0, false
	}
	hexNumber := strings.Join(chunk[1:listLength], " ")

	decimal, ok := hexToDecimal(hexNumber)
	if !ok {
		logging.Warning(fmt.Sprintf("conversion to decimal does not work: %s", hexNumber))
		return 0, listLength, false
	}
	return decimal, listLength, true
}

func handleOctatHexBlock(hexVal string, chunk []string) (int, bool) {
	listLength, err := hexToInt(string(hexVal[1]))
	if err != nil {
		logging.Error(fmt.Sprintf("length problem: %s : %v", hexVal, err))
		return 0, false
	}

	if listLength == 0 {
		logging.Error(fmt.Sprintf("length problem: List length: %d\n", listLength))
		return 0, false
	}
	if listLength == 1 {
		return listLength, true
	}

	hexOctat := strings.Join(chunk[1:listLength], "")

	_, err = hex.DecodeString(hexOctat)
	if err != nil {
		fmt.Println(err)
		return listLength, false
	}
	return listLength, true
}

func hexToInt(hex string) (int, error) {
	converted, err := strconv.ParseInt(hex, 16, 64)
	if err != nil {
		return 0, err
	}
	return int(converted), nil
}

// Extracting the WattHours from the processed readout
// TODO: This is not really robust and only works with the current readout.
// What does the 255 30 actually mean?
func extractUsageFromArray(decimals []int64) (float64, error) {
	for i := 2; i < len(decimals); i++ {
		if decimals[i-1] == 255 && decimals[i-2] == 30 {
			wattHours := float64(decimals[i]) / 10
			return wattHours / 1000, nil
		}
	}
	return 0, fmt.Errorf("cannot extract energy usage from array '%v'", decimals)
}

func hexToDecimal(hexValue string) (int64, bool) {
	hexValueCleaned := strings.Replace(hexValue, " ", "", -1)
	decimal, err := strconv.ParseInt(hexValueCleaned, 16, 64)
	if err != nil {
		logging.Error(err.Error())
		return 0, false
	}
	return decimal, true
}
