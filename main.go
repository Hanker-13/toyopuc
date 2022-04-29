package main

import (
	"time"
	"toyopuc/log"
	"toyopuc/toyopuc"
)

func main() {
	log.Debug("started...")

	// 初始化连接
	handler := toyopuc.NewTCPClientHandler("127.0.0.1:9991")
	handler.Timeout = time.Duration(3) * time.Second
	// handler.SlaveId = 0xFF
	err := handler.Connect()
	if err != nil {
		return
	}

	client := toyopuc.NewClient(handler)
	// Panic处理
	// 防止过程中掉线导致的程序崩溃
	// defer func() {
	// 	if err := recover(); err != nil {
	// 		log.Error("modbus tru comm error, ", err)
	// 	}
	// }()

	result, err := client.ReadDataExpansionMultipoint(1, 1, 1, []byte{112}, []uint16{3072}, []byte{0}, []uint16{8192}, []byte{8}, []uint16{0})
	// err = client.WriteDataExpansionByte(byte(0), uint16(8192), []byte{0x12, 0x34})
	if err != nil {
		log.Error(err)
	}
	log.Debug("result: ", result)
}
