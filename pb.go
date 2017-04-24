package main

/*

package main;
	
enum FOO { X = 17; };
	
message Test {
	  required string label = 1;
	  optional int32 type = 2 [default=77];
	  repeated int64 reps = 3;
	  optional group OptionalGroup = 4 {
	    required string RequiredField = 5;
	  }
}

*/

import (
  "log"
  "github.com/golang/protobuf/proto"
)

func main() {

  test := &Test{
      Label: proto.String("hello"),
      Type:  proto.Int32(17),
      Optionalgroup: &Test_OptionalGroup{
          RequiredField: proto.String("good bye"),
      },
  }
  data, err := proto.Marshal(test)
  if err != nil {
      log.Fatal("marshaling error: ", err)
  }
  newTest := &Test{}
  err = proto.Unmarshal(data, newTest)
  if err != nil {
      log.Fatal("unmarshaling error: ", err)
  }
  // Now test and newTest contain the same data.
  if test.GetLabel() != newTest.GetLabel() {
      log.Fatalf("data mismatch %q != %q", test.GetLabel(), newTest.GetLabel())
  }

  log.Printf("Unmarshalled to: %+v", newTest)

}
