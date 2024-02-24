package main

import (
	"ZINX/myDemo/protobufDemo/pb"
	"fmt"
	"google.golang.org/protobuf/proto"
)

func main() {
	person := &pb.Person{
		Name:   "Liudehua",
		Age:    114514,
		Emails: []string{"liu@gmail.com", "dehua@163.com"},
		Phones: []*pb.PhoneNumber{
			&pb.PhoneNumber{
				Number: "1314",
				Type:   pb.PhoneType_MOBILE,
			},
			&pb.PhoneNumber{
				Number: "666",
				Type:   pb.PhoneType_HOME,
			},
			&pb.PhoneNumber{
				Number: "999",
				Type:   pb.PhoneType_WORK,
			},
		},
	}
	// 编码
	// 将person对象，就是protobuf的message进行序列化，得到一个二进制文件
	data, err := proto.Marshal(person)
	// data是要传输的数据，对端需要按照Message Person的格式进行解析
	if err != nil {
		fmt.Println("Marshal err:", err)
	}

	// 解码
	newPerson := &pb.Person{}
	err = proto.Unmarshal(data, newPerson)
	if err != nil {
		fmt.Println("Unmarshal err", err)
	}
	fmt.Println("原数据:", person)
	fmt.Println("解码后的数据:", newPerson)

}
