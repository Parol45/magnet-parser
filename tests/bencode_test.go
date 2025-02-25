package main

import (
	"fmt"
	"magnet-parser/bencode_converters/json"
	"testing"
)

func TestJsonConversion(t *testing.T) {
	firstBcode := "d3:abc13:zovebathohlov4:ichei-42ee"
	firstInpJson := `{"iche":-42,"abc":"zovebathohlov"}`
	// json keys sorted
	firstOutJson := `{"abc":"zovebathohlov","iche":-42}`
	// ------------------------------------------------------------------------------------------------------
	bcode, err := json.Encode(firstInpJson)
	if err != nil {
		t.Fatal(err)
	}
	if string(bcode) != firstBcode {
		t.Fatal(fmt.Sprintf("Failed at first test: %s != %s", string(bcode), firstBcode))
	}
	// ------------------------------------------------------------------------------------------------------
	js, err := json.Decode([]byte(firstBcode))
	if err != nil {
		t.Fatal(err)
	}
	if js != firstOutJson {
		t.Fatal(fmt.Sprintf("Failed at first test: %s != %s", js, firstOutJson))
	}
	// ------------------------------------------------------------------------------------------------------
	// ------------------------------------------------------------------------------------------------------

	secondBcode := "l4:ichei-42ee"
	secondJson := `["iche",-42]`
	// ------------------------------------------------------------------------------------------------------
	bcode, err = json.Encode(secondJson)
	if err != nil {
		t.Fatal(err)
	}
	if string(bcode) != secondBcode {
		t.Fatal(fmt.Sprintf("Failed at second test: %s != %s", string(bcode), secondBcode))
	}
	// ------------------------------------------------------------------------------------------------------
	js, err = json.Decode([]byte(secondBcode))
	if err != nil {
		t.Fatal(err)
	}
	if js != secondJson {
		t.Fatal(fmt.Sprintf("Failed at second test: %s != %s", js, secondJson))
	}
	// ------------------------------------------------------------------------------------------------------
	// ------------------------------------------------------------------------------------------------------

	complexBcode := "d13:dataTransfersld10:filterName11:min10_max208:gsonType34:OpcSourceFilterStorageDataTransfer9:namespacei0e7:opcName8:OPC_test15:pollingIntervali1000e9:sourceTag26:__dummy_saw300000|1_round110:storageTag27:a__dummy_saw300000|1_round14:type6:DOUBLEee7:filtersd5:emptyde11:min10_max20d10:approxType7:NEAREST22:frozenSignalIntervalMsi10000e8:intervali1e6:jsFunc54:datapoint.Value = datapoint.Value*10; return datapoint8:maxDeltai0e8:maxValuei49e8:minValuei-4e9:timeoutMsi5000eeee"
	complexJsonInp := `{"filters":{"min10_max20":{"timeoutMs":5000,"maxDelta":0,"maxValue":49,"minValue":-4,"interval":1,"approxType":"NEAREST","frozenSignalIntervalMs":10000,"jsFunc":"datapoint.Value = datapoint.Value*10; return datapoint"},"empty":{}},"dataTransfers":[{"gsonType":"OpcSourceFilterStorageDataTransfer","opcName":"OPC_test","namespace":0,"sourceTag":"__dummy_saw300000|1_round1","storageTag":"a__dummy_saw300000|1_round1","filterName":"min10_max20","type":"DOUBLE","pollingInterval":1000}]}`
	// json keys sorted
	complexJsonOut := `{"dataTransfers":[{"filterName":"min10_max20","gsonType":"OpcSourceFilterStorageDataTransfer","namespace":0,"opcName":"OPC_test","pollingInterval":1000,"sourceTag":"__dummy_saw300000|1_round1","storageTag":"a__dummy_saw300000|1_round1","type":"DOUBLE"}],"filters":{"empty":{},"min10_max20":{"approxType":"NEAREST","frozenSignalIntervalMs":10000,"interval":1,"jsFunc":"datapoint.Value = datapoint.Value*10; return datapoint","maxDelta":0,"maxValue":49,"minValue":-4,"timeoutMs":5000}}}`
	// ------------------------------------------------------------------------------------------------------
	bcode, err = json.Encode(complexJsonInp)
	if err != nil {
		t.Fatal(err)
	}
	if string(bcode) != complexBcode {
		t.Fatal(fmt.Sprintf("Failed at complex test: %s != %s", string(bcode), complexBcode))
	}
	// ------------------------------------------------------------------------------------------------------
	js, err = json.Decode([]byte(complexBcode))
	if err != nil {
		t.Fatal(err)
	}
	if js != complexJsonOut {
		t.Fatal(fmt.Sprintf("Failed at complex test: %s != %s", js, complexJsonOut))
	}
	// ------------------------------------------------------------------------------------------------------
	// ------------------------------------------------------------------------------------------------------


}