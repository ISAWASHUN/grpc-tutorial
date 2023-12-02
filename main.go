package main

import (
	"fmt"
	"gRPC-tutorial/pb"
	"log"

	"github.com/golang/protobuf/jsonpb"
)

func main() {
	employee := &pb.Employee {
		Id: 1,
		Name: "Suzuki",
		Email: "test@test.com",
		Occupation: pb.Occupation_ENGINEER,
		PhoneNumber: []string{"080-1234^5678", "090-1234-5678"},
		Project: map[string]*pb.Company_Project{"ProjectX": {}},
		Profile: &pb.Employee_Text{
			Text: "My name is Suzuki",
		},
		Birthday: &pb.Date{
			Year: 2000,
			Month: 1,
			Day: 1,
		},
	}

	// binData, err := proto.Marshal(employee)
	// if err != nil {
	// 	log.Fatalln("Cannot serialize", err)
	// }

	// if err := ioutil.WriteFile("test.bin", binData, 0666); err != nil {
	// 	log.Fatalln("Cannot write", err)
	// }

	// in, err := ioutil.ReadFile("test.bin")
	// if err != nil {
	// 	log.Fatalln("Cannot read file", err)
	// }

	// readEmployee := &pb.Employee{}

	// err = proto.Unmarshal(in, readEmployee)
	// if err != nil {
	// 	log.Fatalln("Cannot deserialize", err)
	// }

	// fmt.Println(readEmployee)

	m := jsonpb.Marshaler{}
	out , err := m.MarshalToString(employee)
	if err != nil {
		log.Fatalln("Cannot marshal to json", err)
	}

	// fmt.Println(out)

	readEmployee := &pb.Employee{}
	if err := jsonpb.UnmarshalString(out, readEmployee); err != nil {
		log.Fatalln("Cannot unmarshal to json", err)
	}

	fmt.Println(readEmployee)
}