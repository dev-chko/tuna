package main

import (
	"fmt"
	"bytes"
	"encoding/json"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)
type SmartContract struct{

}
type Tuna struct {
	Vessel string `json:"vessel"`
	Timestamp string `json:"datetime"`
	Location string `json:"location"`
	Holder string `json:"holder`
}

func (s *SmartContract) Init(stub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

func (s *SmartContract) Invoke(stub shim.ChaincodeStubInterface) sc.Response {
	fn, args := stub.GetFunctionAndParameters()
	if fn == "initLedger"{
		return s.initLedger(stub)
	} else if fn == "queryAllTuna"{
		return s.queryAllTuna(stub, args)
	}else if fn == "queryTuna"{
		return s.queryTuna(stub, args)
	}else if fn == "recordTuna"{
		return s.recordTuna(stub, args)
	}else if fn == "changeTunaHolder"{
		return s.changeTunaHolder(stub, args)
	}else {
		return shim.Error("Invaild Function Name")
	}
}

func (s *SmartContract) initLedger(stub shim.ChaincodeStubInterface) sc.Response {
	init_tuna := []Tuna{
		Tuna{Vessel:"923F", Timestamp:"1502403925", Location: "61, -50", Holder: "Miriam"},
		Tuna{Vessel:"245F", Timestamp:"1502452451", Location: "63, -51", Holder: "Carl"},
		Tuna{Vessel:"AS12", Timestamp:"1502235341", Location: "62, -57", Holder: "Sarah"},
		Tuna{Vessel:"652A", Timestamp:"1505776124", Location: "66, -58", Holder: "Carl"},
		Tuna{Vessel:"FA12", Timestamp:"1504567344", Location: "67, -51", Holder: "Sarah"},
		Tuna{Vessel:"775B", Timestamp:"1505619973", Location: "68, -54", Holder: "Miriam"},
		Tuna{Vessel:"241C", Timestamp:"1501232415", Location: "68, -50", Holder: "Miriam"},
		Tuna{Vessel:"B132", Timestamp:"1506612341", Location: "62, -57", Holder: "Sarah"},
		Tuna{Vessel:"A98A", Timestamp:"1502351242", Location: "69, -58", Holder: "Miriam"},
		Tuna{Vessel:"7892", Timestamp:"1505874334", Location: "68, -59", Holder: "Carl"},
	}
	// store to World State with init_tuna
	i := 0
	for i < len(init_tuna){
		fmt.Printf("i is ", i)

		tunaBytes, _ := json.Marshal(init_tuna[i])
		stub.PutState(strconv.Itoa(i+1), tunaBytes)
		fmt.Printf("Added ", init_tuna[i])
		i = i+1
	}
	return shim.Success(nil)
}

func (s *SmartContract) queryTuna(stub shim.ChaincodeStubInterface, args []string) sc.Response{
	//check paramter count -> 1
	if len(args) !=1 {
		return shim.Error("Incorrect number of arguments.")
	}
	//GetState from world State with the first parameter
	tunaBytes, _ := stub.GetState(args[0])
	if tunaBytes == nil {
		return shim.Error("Could not find tuna data")
	} else {
		return shim.Success(tunaBytes)
	}
	//return the result
}

func (s *SmartContract) recordTuna(stub shim.ChaincodeStubInterface, args []string) sc.Response{
	//check paramter 
	if len(args) !=5 {
		return shim.Error("Incorrect number of arguments.")
	}
	//make TUNA object witrh parameters
	var tuna = Tuna{Vessel: args[1], Timestamp: args[2], Location: args[3], Holder: args[4]}
	//convert TUNA object to json byte[]
	//PutState to world state
	//return the Result 
	tunaBytes, _ := json.Marshal(tuna)
	stub.PutState(args[0], tunaBytes)
	err := stub.PutState(args[0], tunaBytes)
	if err != nil{
		return shim.Error(fmt.Sprintf("Failed to recode tuna"))
	}else {
		return shim.Success(nil)
	}
}
func (s *SmartContract) changeTunaHolder(stub shim.ChaincodeStubInterface, args []string) sc.Response{
	//check paramter ->2
	if len(args) !=2 {
		return shim.Error("incorrect parament")
	}
	//GetState from world state
	tunaBytes, _ := stub.GetState(args[0])
	tuna := Tuna{}
	json.Unmarshal(tunaBytes, &tuna)
	//change the holder value of 
	tuna.Holder = args[1]
	//Putstate the changed TUNA object to world state
	tunaBytes, _ = json.Marshal(tuna)
	err := stub.PutState(args[0], tunaBytes)
	//return the result
	if err != nil{
		return shim.Error(fmt.Sprintf("Failed to change tuna", args[0]))
	}else {
		return shim.Success(nil)
	}
}
func (s *SmartContract) queryAllTuna(stub shim.ChaincodeStubInterface, args []string) sc.Response{
	if len(args) !=2{
		return shim.Error("incorrect parament")
	}
	startKey, _ := stub.GetState(args[0])
	endKey, _ := stub.GetState(args[1])
	resultIterator, err := stub.GetStateByRange(string(startKey), string(endKey))

	if err != nil{
		return shim.Error(err.Error())
	}
	defer resultIterator.Close()
	//buffer is a JSON array containing Query Results
	var buffer bytes.Buffer
	buffer.WriteString("[")
	bArrayMemberAlreadyWritten := false
	for resultIterator.HasNext() {
		queryResponse, err := resultIterator.Next()
		if err != nil{
			return shim.Error(err.Error())
		}
	//Add a comma before array members, suppress it for the first array member
	if bArrayMemberAlreadyWritten == true {
		buffer.WriteString(",")
	}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")
		buffer.WriteString(",\"Record\":")
		//Record is a JSON objects so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")
	fmt.Printf("- queryAllTuna :\n$s\n", buffer.String())
	return shim.Success(buffer.Bytes())
}


func main(){
	err := shim.Start(new(SmartContract))
	if err != nil{
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}

