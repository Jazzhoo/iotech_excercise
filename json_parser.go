package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

//input data structure

type Devices struct {
	Devices []Device `json:"devices"`
}

type Device struct {
	Name      string `json:"name"`
	Type      string `json:"type"`
	Info      string `json:"info"`
	Value     string `json:"value"`
	Timestamp string `json:"timestamp"`
}

//output data structure

type OutputStruct struct {
	ValueTotal int32  `json:"ValueTotal"`
	UUIDS      []UUID `json:"UUIDS"`
}

type UUID struct {
	UUID string `json:"uuid"`
}

// helper function to extract uuid from string
func uuidExtractor(inputString string) string {
	strArr := strings.Split(inputString, " ")
	var uuid string
	for _, strItem := range strArr {
		if strings.Contains(strItem, "uuid") {
			uuid = strings.Replace(strItem, "uuid:", "", 1)
			uuid = strings.Trim(uuid, ",")
			break
		}
	}
	return uuid
}

func main() {
	fileContent, err := os.Open("./input.json")
	if err != nil {
		log.Fatal("Failed: Error with opening file", err)
	} else {
		fmt.Println("Success: File open")
	}
	defer fileContent.Close()

	jsonRawData, _ := io.ReadAll(fileContent)
	var devices Devices
	var outStct OutputStruct
	outStct.ValueTotal = 0

	err = json.Unmarshal([]byte(jsonRawData), &devices)
	if err != nil {
		log.Fatal("Failed: Unmarshal data failed", err)
	} else {
		fmt.Println("Success: Unmarshal data")
	}
	for _, dev := range devices.Devices {
		decoded_value, _ := base64.StdEncoding.DecodeString(dev.Value)
		dv_int, _ := strconv.Atoi(string(decoded_value))
		ts_int, _ := strconv.Atoi(dev.Timestamp)
		if int64(ts_int) < time.Now().Unix() {
		} else {
			outStct.ValueTotal += int32(dv_int)
			item := UUID{uuidExtractor(dev.Info)}
			outStct.UUIDS = append(outStct.UUIDS, item)
		}
	}

	outStr, err := json.MarshalIndent(outStct, "", " ")
	if err != nil {
		log.Fatal("Failed: Marshal data failed", err)
	} else {
		fmt.Println("Success: Marshal data")
	}
	f, err := os.Create("./output.json")
	if err != nil {
		log.Fatal("Failed: Output file creation error", err)
	} else {
		fmt.Println("Success: Output file created")
	}
	_, err = f.Write(outStr)
	if err != nil {
		log.Fatal("Output data write error", err)
	} else {
		fmt.Println("Success: Output data written")
	}

}
