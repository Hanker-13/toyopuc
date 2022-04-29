package toyopuc

type Client interface {
	// 程序顺序
	// 程序顺序 字读出
	// result CDAB
	ReadSequentialProgramWord(address, quantity uint16) (results []byte, err error)
	// 程序顺序 字写入
	// CDAB
	WriteSequentialProgramWord(address uint16, value []uint16) (err error)
	// I/O寄存器
	// I/O寄存器字读出
	// result CDAB
	ReadIOWord(address, quantity uint16) (results []byte, err error)
	// I/O寄存器 字写入
	// CDAB
	WriteIOWord(address uint16, value []uint16) (err error)
	// I/O寄存器 字节读出
	ReadIOByte(address, quantity uint16) (results []byte, err error)
	// I/O寄存器 字节写入
	WriteIOByte(address uint16, value []byte) (err error)
	// I/O寄存器 位读出
	// 一次只能读一位
	ReadIOBit(address uint16) (result bool, err error)
	// I/O寄存器 位写入
	// 一次只能写一位
	// 1 ON 0 OFF
	WriteIOBit(address uint16, value byte) (err error)
	// I/O寄存器 多点字读出
	// CDAB
	ReadIOMultipointWord(address []uint16) (results []byte, err error)
	// I/O寄存器 多点字写入
	// CDAB
	WriteIOMultipointWord(address, value []uint16) (err error)
	// I/O寄存器 多点字节读出
	ReadIOMultipointByte(address []uint16) (results []byte, err error)
	// I/O寄存器 多点字节写入
	WriteIOMultipointByte(address []uint16, value []byte) (err error)
	// I/O寄存器 多点位读出
	ReadIOMultipointBit(address []uint16) (results []bool, err error)
	// I/O寄存器 多点位写入
	// 1 ON 0 OFF
	WriteIOMultipointBit(address []uint16, value []byte) (err error)

	// 程序扩展
	// 程序扩展 字读出
	// no 01 PRG1 02 PRG2 03 PRG3
	// CDAB
	ReadProgramExpansionWord(no byte, address, quantity uint16) (results []byte, err error)
	// 程序扩展 字写入
	// CDAB
	WriteProgramExpansionWord(no byte, address uint16, value []uint16) (err error)

	// 数据扩展
	// 数据扩展 字读出
	// no 00 扩展位区域（包含ES\EN\H） 01 PRG.1 02 PRG.2 03 PRG.3 08 扩展寄存器区域（U）
	// CDAB
	ReadDataExpansionWord(no byte, address, quantity uint16) (results []byte, err error)
	// 数据扩展 字写入
	// CDAB
	WriteDataExpansionWord(no byte, address uint16, value []uint16) (err error)

	// 数据扩展 字节读出
	ReadDataExpansionByte(no byte, address, quantity uint16) (results []byte, err error)
	// 数据扩展 字节写入
	WriteDataExpansionByte(no byte, address uint16, value []byte) (err error)

	// 数据扩展 多点读出
	ReadDataExpansionMultipoint(numBit, numByte, numWord byte, bitNo []byte, bitAddr []uint16, bytesNo []byte, bytesAddr []uint16, wordNo []byte, wordAddr []uint16) (results []byte, err error)
	// TODO
	// 数据扩展 多点写入
}
