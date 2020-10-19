package serializer

import (
	"fmt"
	"io/ioutil"

	"google.golang.org/protobuf/proto"
)

//WriteProtobufToJSONFile writes protocol buffer message to JSON file
func WriteProtobufToJSONFile(message proto.Message, fileName string) error {
	data, err := ProtobufToJSON(message)
	if err != nil {
		return fmt.Errorf("Cannot marshal proto buffer message to JSON: %v", err)
	}

	err = ioutil.WriteFile(fileName, []byte(data), 0644)
	if err != nil {
		return fmt.Errorf("Cannot write JSON data to file: %v", err)
	}

	return nil
}

//WriteProtobufToBinaryFile writes protocol buffer message to binary file
func WriteProtobufToBinaryFile(message proto.Message, fileName string) error {
	data, err := proto.Marshal(message)
	if err != nil {
		return fmt.Errorf("Cannot marshal proto message to binary: %v", err)
	}

	err = ioutil.WriteFile(fileName, []byte(data), 0644)
	if err != nil {
		return fmt.Errorf("Cannot write binary data to file: %v", err)
	}

	return nil
}

//ReadProtobufFromBinaryFile reads protocol buffer message from binary file
func ReadProtobufFromBinaryFile(fileName string, message proto.Message) error {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return fmt.Errorf("Cannot read file: %v", err)
	}

	err = proto.Unmarshal(data, message)
	if err != nil {
		return fmt.Errorf("Cannot unmarshal binary to proto message: %v", err)
	}

	return nil
}
