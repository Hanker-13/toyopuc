package toyopuc

import (
	"encoding/binary"
	"fmt"
)

// ClientHandler is the interface that groups the Packager and Transporter methods.
type ClientHandler interface {
	Packager
	Transporter
}

type client struct {
	packager    Packager
	transporter Transporter
}

// NewClient creates a new toyopuc client with given backend handler.
func NewClient(handler ClientHandler) Client {
	return &client{packager: handler, transporter: handler}
}

// NewClient2 creates a new toyopuc client with given backend packager and transporter.
func NewClient2(packager Packager, transporter Transporter) Client {
	return &client{packager: packager, transporter: transporter}
}

// ReadSequentialProgramWord
// 顺序程序 读字
//  Function code         : 1 byte (0x18)
func (toyopuc *client) ReadSequentialProgramWord(address, quantity uint16) (results []byte, err error) {
	// 长度校验
	if quantity < 1 || quantity > 0x200 {
		err = fmt.Errorf("toyopuc: quantity '%v' must be between '%v' and '%v',", quantity, 1, 0x200)
		return
	}
	request := ProtocolDataUnit{
		FunctionCode: FunSequentialProgramReadWord,
		Data:         dataBlock(address, quantity),
	}
	response, err := toyopuc.send(&request)
	if err != nil {
		return
	}
	// TODO
	// []byte to []uint16
	// CDAB
	results = response.Data
	return
}

// WriteSequentialProgramWord
// 顺序程序 写字
//  Function code         : 1 byte (0x19)
func (toyopuc *client) WriteSequentialProgramWord(address uint16, value []uint16) (err error) {
	quantity := len(value)
	if quantity < 1 || quantity > 0x200 {
		err = fmt.Errorf("toyopuc: quantity '%v' must be between '%v' and '%v',", quantity, 1, 0x200)
		return
	}
	request := ProtocolDataUnit{
		FunctionCode: FunSequentialProgramWriteWord,
		Data:         dataBlockSuffix(value, address),
	}
	_, err = toyopuc.send(&request)
	if err != nil {
		return
	}

	return
}

// ReadIOWord
// IO寄存器 读字
//  Function code         : 1 byte (0x1C)
func (toyopuc *client) ReadIOWord(address, quantity uint16) (results []byte, err error) {
	// 长度校验
	if quantity < 1 || quantity > 0x200 {
		err = fmt.Errorf("toyopuc: quantity '%v' must be between '%v' and '%v',", quantity, 1, 0x200)
		return
	}
	request := ProtocolDataUnit{
		FunctionCode: FunIOReadWord,
		Data:         dataBlock(address, quantity),
	}
	response, err := toyopuc.send(&request)
	if err != nil {
		return
	}
	// TODO
	// []byte to []uint16
	// CDAB
	results = response.Data
	return
}

// WriteIOWord
// IO寄存器 写字
//  Function code         : 1 byte (0x1D)
func (toyopuc *client) WriteIOWord(address uint16, value []uint16) (err error) {
	quantity := len(value)
	if quantity < 1 || quantity > 0x200 {
		err = fmt.Errorf("toyopuc: quantity '%v' must be between '%v' and '%v',", quantity, 1, 0x200)
		return
	}
	request := ProtocolDataUnit{
		FunctionCode: FunIOWriteWord,
		Data:         dataBlockSuffix(value, address),
	}
	_, err = toyopuc.send(&request)
	if err != nil {
		return
	}
	return
}

// ReadIOByte
// IO寄存器 读字节
//  Function code         : 1 byte (0x1E)
func (toyopuc *client) ReadIOByte(address, quantity uint16) (results []byte, err error) {
	// 长度校验
	if quantity < 1 || quantity > 0x200 {
		err = fmt.Errorf("toyopuc: quantity '%v' must be between '%v' and '%v',", quantity, 1, 0x200)
		return
	}
	request := ProtocolDataUnit{
		FunctionCode: FunIOReadByte,
		Data:         dataBlock(address, quantity),
	}
	response, err := toyopuc.send(&request)
	if err != nil {
		return
	}

	results = response.Data
	return
}

// WriteIOByte
// IO寄存器 写字节
//  Function code         : 1 byte (0x1F)
func (toyopuc *client) WriteIOByte(address uint16, value []byte) (err error) {
	quantity := len(value)
	if quantity < 1 || quantity > 0x200 {
		err = fmt.Errorf("toyopuc: quantity '%v' must be between '%v' and '%v',", quantity, 1, 0x200)
		return
	}
	request := ProtocolDataUnit{
		FunctionCode: FunIOWriteByte,
		Data:         dataBlockSuffixByte(value, address),
	}
	r, err := toyopuc.send(&request)
	if err != nil {
		return
	}
	fmt.Println(r)
	return
}

// ReadIOBit
// IO寄存器 读位
//  Function code         : 1 byte (0x20)
func (toyopuc *client) ReadIOBit(address uint16) (results bool, err error) {
	request := ProtocolDataUnit{
		FunctionCode: FunIOReadBit,
		Data:         dataBlock(address),
	}
	response, err := toyopuc.send(&request)
	if err != nil {
		return
	}

	if response.Data[0] == 1 {
		results = true
	}

	return
}

// WriteIOBit
// IO寄存器 写位
//  Function code         : 1 byte (0x21)
func (toyopuc *client) WriteIOBit(address uint16, value byte) (err error) {
	request := ProtocolDataUnit{
		FunctionCode: FunIOWriteBit,
		Data:         dataBlockSuffixBit(value, address),
	}
	r, err := toyopuc.send(&request)
	if err != nil {
		return
	}
	fmt.Println(r)
	return
}

// ReadIOMultipointWord
// IO寄存器 多点读出字
//  Function code         : 1 byte (0x22)
func (toyopuc *client) ReadIOMultipointWord(address []uint16) (results []byte, err error) {
	// 长度校验
	quantity := len(address)
	if quantity < 1 || quantity > 0x80 {
		err = fmt.Errorf("toyopuc: quantity '%v' must be between '%v' and '%v',", quantity, 1, 0x80)
		return
	}
	request := ProtocolDataUnit{
		FunctionCode: FunIOReadMultipointWord,
		Data:         dataBlockList(address),
	}
	response, err := toyopuc.send(&request)
	if err != nil {
		return
	}
	// TODO
	// []byte to []uint16
	// CDAB
	results = response.Data
	return
}

// WriteIOMultipointWord
// IO寄存器 多点写入字
//  Function code         : 1 byte (0x23)
func (toyopuc *client) WriteIOMultipointWord(address, value []uint16) (err error) {
	quantityAddr := len(address)
	quantityVal := len(value)
	if quantityAddr != quantityVal {
		err = fmt.Errorf("toyopuc: the quantity of addresses must be equal to the quantity of values, address quantity: '%v' ,value quantity: '%v'", quantityAddr, quantityVal)
		return
	}
	if quantityAddr < 1 || quantityAddr > 0x80 || quantityVal < 1 || quantityVal > 0x80 {
		err = fmt.Errorf("toyopuc: quantity '%v' must be between '%v' and '%v',", quantityVal, 1, 0x80)
		return
	}
	request := ProtocolDataUnit{
		FunctionCode: FunIOWriteMultipointWord,
		Data:         dataBlockAVList(address, value),
	}
	_, err = toyopuc.send(&request)
	if err != nil {
		return
	}
	return
}

// ReadIOMultipointByte
// IO寄存器 多点读出字节
//  Function code         : 1 byte (0x24)
func (toyopuc *client) ReadIOMultipointByte(address []uint16) (results []byte, err error) {
	// 长度校验
	quantity := len(address)
	if quantity < 1 || quantity > 0x80 {
		err = fmt.Errorf("toyopuc: quantity '%v' must be between '%v' and '%v',", quantity, 1, 0x80)
		return
	}
	request := ProtocolDataUnit{
		FunctionCode: FunIOReadMultipointByte,
		Data:         dataBlockList(address),
	}
	response, err := toyopuc.send(&request)
	if err != nil {
		return
	}

	results = response.Data
	return
}

// WriteIOMultipointByte
// IO寄存器 多点写入字节
//  Function code         : 1 byte (0x25)
func (toyopuc *client) WriteIOMultipointByte(address []uint16, value []byte) (err error) {
	quantityAddr := len(address)
	quantityVal := len(value)
	if quantityAddr != quantityVal {
		err = fmt.Errorf("toyopuc: the quantity of addresses must be equal to the quantity of values, address quantity: '%v' ,value quantity: '%v'", quantityAddr, quantityVal)
		return
	}
	if quantityAddr < 1 || quantityAddr > 0x80 || quantityVal < 1 || quantityVal > 0x80 {
		err = fmt.Errorf("toyopuc: quantity '%v' must be between '%v' and '%v',", quantityVal, 1, 0x80)
		return
	}
	request := ProtocolDataUnit{
		FunctionCode: FunIOWriteMultipointByte,
		Data:         dataBlockAVByteList(address, value),
	}
	_, err = toyopuc.send(&request)
	if err != nil {
		return
	}
	return
}

// ReadIOMultipointBit
// IO寄存器 多点读出位
//  Function code         : 1 byte (0x26)
func (toyopuc *client) ReadIOMultipointBit(address []uint16) (results []bool, err error) {
	// 长度校验
	quantity := len(address)
	if quantity < 1 || quantity > 0x80 {
		err = fmt.Errorf("toyopuc: quantity '%v' must be between '%v' and '%v',", quantity, 1, 0x80)
		return
	}
	request := ProtocolDataUnit{
		FunctionCode: FunIOReadMultipointBit,
		Data:         dataBlockList(address),
	}
	response, err := toyopuc.send(&request)
	if err != nil {
		return
	}

	results = make([]bool, len(response.Data))
	for k, v := range response.Data {
		if v == 1 {
			results[k] = true
		}
	}
	return
}

// WriteIOMultipointBit
// IO寄存器 多点写入位
//  Function code         : 1 byte (0x27)
func (toyopuc *client) WriteIOMultipointBit(address []uint16, value []byte) (err error) {
	quantityAddr := len(address)
	quantityVal := len(value)
	if quantityAddr != quantityVal {
		err = fmt.Errorf("toyopuc: the quantity of addresses must be equal to the quantity of values, address quantity: '%v' ,value quantity: '%v'", quantityAddr, quantityVal)
		return
	}
	if quantityAddr < 1 || quantityAddr > 0x80 || quantityVal < 1 || quantityVal > 0x80 {
		err = fmt.Errorf("toyopuc: quantity '%v' must be between '%v' and '%v',", quantityVal, 1, 0x80)
		return
	}
	request := ProtocolDataUnit{
		FunctionCode: FunIOWriteMultipointBit,
		Data:         dataBlockAVByteList(address, value),
	}
	_, err = toyopuc.send(&request)
	if err != nil {
		return
	}
	return
}

// ReadProgramExpansionWord
// 程序扩展 读字
//  Function code         : 1 byte (0x90)
func (toyopuc *client) ReadProgramExpansionWord(no byte, address, quantity uint16) (results []byte, err error) {
	// 长度校验
	if quantity < 1 || quantity > 0x200 {
		err = fmt.Errorf("toyopuc: quantity '%v' must be between '%v' and '%v',", quantity, 1, 0x200)
		return
	}
	request := ProtocolDataUnit{
		FunctionCode: FunProgramExpansionReadWord,
		Data:         dataBlockExpansion(no, address, quantity),
	}
	response, err := toyopuc.send(&request)
	if err != nil {
		return
	}
	// TODO
	// []byte to []uint16
	// CDAB
	results = response.Data
	return
}

// WriteProgramExpansionWord
// 程序扩展 字写入
//  Function code         : 1 byte (0x91)
func (toyopuc *client) WriteProgramExpansionWord(no byte, address uint16, value []uint16) (err error) {
	quantity := len(value)
	if quantity < 1 || quantity > 0x200 {
		err = fmt.Errorf("toyopuc: quantity '%v' must be between '%v' and '%v',", quantity, 1, 0x200)
		return
	}
	request := ProtocolDataUnit{
		FunctionCode: FunProgramExpansionWriteWord,
		Data:         dataBlockExpansionSuffix(no, value, address),
	}
	_, err = toyopuc.send(&request)
	if err != nil {
		return
	}
	return
}

// ReadDataExpansionWord
// 数据扩展 读字
//  Function code         : 1 byte (0x94)
func (toyopuc *client) ReadDataExpansionWord(no byte, address, quantity uint16) (results []byte, err error) {
	// 长度校验
	if quantity < 1 || quantity > 0x200 {
		err = fmt.Errorf("toyopuc: quantity '%v' must be between '%v' and '%v',", quantity, 1, 0x200)
		return
	}
	request := ProtocolDataUnit{
		FunctionCode: FunDataExpansionReadWord,
		Data:         dataBlockExpansion(no, address, quantity),
	}
	response, err := toyopuc.send(&request)
	if err != nil {
		return
	}
	// TODO
	// []byte to []uint16
	// CDAB
	results = response.Data
	return
}

// WriteDataExpansionWord
// 数据扩展 字写入
//  Function code         : 1 byte (0x95)
func (toyopuc *client) WriteDataExpansionWord(no byte, address uint16, value []uint16) (err error) {
	quantity := len(value)
	if quantity < 1 || quantity > 0x200 {
		err = fmt.Errorf("toyopuc: quantity '%v' must be between '%v' and '%v',", quantity, 1, 0x200)
		return
	}
	request := ProtocolDataUnit{
		FunctionCode: FunDateExpansionWriteWord,
		Data:         dataBlockExpansionSuffix(no, value, address),
	}
	_, err = toyopuc.send(&request)
	if err != nil {
		return
	}
	return
}

// ReadDataExpansionByte
// 数据扩展 读字节
//  Function code         : 1 byte (0x96)
func (toyopuc *client) ReadDataExpansionByte(no byte, address, quantity uint16) (results []byte, err error) {
	// 长度校验
	if quantity < 1 || quantity > 0x400 {
		err = fmt.Errorf("toyopuc: quantity '%v' must be between '%v' and '%v',", quantity, 1, 0x400)
		return
	}
	request := ProtocolDataUnit{
		FunctionCode: FunDataExpansionReadByte,
		Data:         dataBlockExpansion(no, address, quantity),
	}
	response, err := toyopuc.send(&request)
	if err != nil {
		return
	}

	results = response.Data
	return
}

// WriteDataExpansionByte
// 数据扩展 写字节
//  Function code         : 1 byte (0x97)
func (toyopuc *client) WriteDataExpansionByte(no byte, address uint16, value []byte) (err error) {
	quantity := len(value)
	if quantity < 1 || quantity > 0x400 {
		err = fmt.Errorf("toyopuc: quantity '%v' must be between '%v' and '%v',", quantity, 1, 0x400)
		return
	}
	request := ProtocolDataUnit{
		FunctionCode: FunDataExpansionWriteByte,
		Data:         dataBlockExpansionSuffixByte(no, value, address),
	}
	r, err := toyopuc.send(&request)
	if err != nil {
		return
	}
	fmt.Println(r)
	return
}

// ReadDataExpansionMultipoint
// 数据扩展 读多点
//  Function code         : 1 byte (0x98)
func (toyopuc *client) ReadDataExpansionMultipoint(numBit, numByte, numWord byte, bitNo []byte, bitAddr []uint16, bytesNo []byte, bytesAddr []uint16, wordNo []byte, wordAddr []uint16) (results []byte, err error) {
	quantity := numBit + numByte + numWord
	if quantity < 1 || quantity > 176 {
		err = fmt.Errorf("toyopuc: address quantity '%v' must be between '%v' and '%v',", quantity, 1, 176)
		return
	}
	dataQuantity := numBit/8 + numByte + numWord*2
	if dataQuantity < 1 || dataQuantity > 128 {
		err = fmt.Errorf("toyopuc: data quantity '%v' must be between '%v' and '%v',", dataQuantity, 1, 128)
		return
	}
	request := ProtocolDataUnit{
		FunctionCode: FunDataExpansionReadMultipoint,
		Data:         dataBlockExpansionSuffixMultipoint(numBit, numByte, numWord, bitNo, bitAddr, bytesNo, bytesAddr, wordNo, wordAddr),
	}
	response, err := toyopuc.send(&request)
	if err != nil {
		return
	}
	results = response.Data
	return
}

// Helpers

// send sends request and checks possible exception in the response.
// 发送请求并检查响应中可能出现的异常
func (toyopuc *client) send(request *ProtocolDataUnit) (response *ProtocolDataUnit, err error) {
	aduRequest, err := toyopuc.packager.Encode(request)
	if err != nil {
		return
	}
	aduResponse, err := toyopuc.transporter.Send(aduRequest)
	if err != nil {
		return
	}
	// 校验
	if err = toyopuc.packager.Verify(aduRequest, aduResponse); err != nil {
		return
	}
	response, err = toyopuc.packager.Decode(aduResponse)
	if err != nil {
		return
	}
	// Check correct function code returned (exception)
	if response.FunctionCode != request.FunctionCode {
		err = fmt.Errorf("toyopuc: function code error")

	}
	// Check error code
	if response.Response != 0x00 {
		err = responseError(response)
		return
	}
	return
}

// dataBlock creates a sequence of uint16 data.
// 数据块创建uint16数据序列 address + quantity ([]byte)
// 顺序程序以及IO寄存器 读字用
func dataBlock(value ...uint16) []byte {
	data := make([]byte, 2*len(value))
	for i, v := range value {
		binary.LittleEndian.PutUint16(data[i*2:], v)
	}
	return data
}

// dataBlockSuffix 创建一个uint16数据序列 address + value ([]byte)
// 顺序程序以及IO寄存器 写字用
func dataBlockSuffix(suffix []uint16, address uint16) []byte {
	// 编码指令 CDAB
	data := make([]byte, 2+len(suffix)*2)
	binary.LittleEndian.PutUint16(data[0:], address)

	suffixList := make([]byte, 2*len(suffix))
	for k, v := range suffix {
		binary.LittleEndian.PutUint16(suffixList[k*2:], v)
	}
	copy(data[2:], suffixList)
	return data
}

// dataBlockSuffixByte 创建一个uint16数据序列 address + value ([]byte)
// IO寄存器 写字节用
func dataBlockSuffixByte(suffix []byte, address uint16) []byte {
	// 编码指令 CDAB
	data := make([]byte, 2+len(suffix))
	binary.LittleEndian.PutUint16(data[0:], address)
	copy(data[2:], suffix)
	return data
}

// dataBlockSuffixBit 创建一个uint16数据序列 address + value ([]byte)
// IO寄存器 写位用
func dataBlockSuffixBit(suffix byte, address uint16) []byte {
	// 编码指令 CDAB
	data := make([]byte, 3)
	binary.LittleEndian.PutUint16(data[0:], address)
	data[2] = suffix
	return data
}

// 数据块创建uint16数据序列 address ([]byte)
// IO寄存器 多点读出字、多点读出字节、多点读出位用
func dataBlockList(value []uint16) []byte {
	data := make([]byte, 2*len(value))
	for i, v := range value {
		binary.LittleEndian.PutUint16(data[i*2:], v)
	}
	return data
}

// 数据块创建uint16数据序列 address + value ([]byte)
// IO寄存器 多点写入字用
func dataBlockAVList(address, value []uint16) []byte {
	data := make([]byte, 2*len(value)*2)
	for i, v := range address {
		binary.LittleEndian.PutUint16(data[i*4:], v)
		binary.LittleEndian.PutUint16(data[i*4+2:], value[i]) // CDAB
	}
	return data
}

// 数据块创建uint16数据序列 address + value ([]byte)
// IO寄存器 多点写入字节用
func dataBlockAVByteList(address []uint16, value []byte) []byte {
	data := make([]byte, 3*len(value))
	for i, v := range address {
		binary.LittleEndian.PutUint16(data[i*3:], v)
		data[i*3+2] = value[i]
	}
	return data
}

// 数据块创建uint16数据序列 no +address + quantity ([]byte)
// 程序扩展\数据扩展 读字\字节用
func dataBlockExpansion(no byte, value ...uint16) []byte {
	data := make([]byte, 2*len(value)+1)
	data[0] = no
	for i, v := range value {
		binary.LittleEndian.PutUint16(data[i*2+1:], v)
	}
	return data
}

// dataBlockExpansionSuffix 创建一个uint16数据序列 no + address + value ([]byte)
// 程序扩展/数据扩展 写字用
func dataBlockExpansionSuffix(no byte, suffix []uint16, address uint16) []byte {
	// 编码指令 CDAB
	data := make([]byte, 2+len(suffix)*2+1)
	data[0] = no
	binary.LittleEndian.PutUint16(data[1:], address)
	suffixList := make([]byte, 2*len(suffix))
	for k, v := range suffix {
		binary.LittleEndian.PutUint16(suffixList[k*2:], v)
	}
	copy(data[3:], suffixList)
	return data
}

// dataBlockExpansionSuffixByte 创建一个uint16数据序列 no + address + value ([]byte)
// 数据扩展 写字节用
func dataBlockExpansionSuffixByte(no byte, suffix []byte, address uint16) []byte {
	// 编码指令 CDAB
	data := make([]byte, 2+len(suffix)+1)
	data[0] = no
	binary.LittleEndian.PutUint16(data[1:], address)
	copy(data[3:], suffix)
	return data
}

// dataBlockExpansionSuffixByte 创建一个uint16数据序列 no + address + value ([]byte)
// 数据扩展 读多点
func dataBlockExpansionSuffixMultipoint(numBit, numByte, numWord byte, bitNo []byte, bitAddr []uint16, bytesNo []byte, bytesAddr []uint16, wordNo []byte, wordAddr []uint16) []byte {
	// 编码指令 CDAB
	data := make([]byte, 3+len(bitNo)+2*len(bitAddr)+len(bytesNo)+2*len(bytesAddr)+len(wordNo)+2*len(wordAddr))
	data[0] = numBit
	data[1] = numByte
	data[2] = numWord
	suffixBit := make([]byte, 3*len(bitNo))
	for k, v := range bitNo {
		suffixBit[k*3] = v
		binary.LittleEndian.PutUint16(suffixBit[k*3+1:], bitAddr[k])
	}
	suffixByte := make([]byte, 3*len(bytesNo))
	for k, v := range bytesNo {
		suffixByte[k*3] = v
		binary.LittleEndian.PutUint16(suffixByte[k*3+1:], bytesAddr[k])
	}
	suffixWord := make([]byte, 3*len(wordNo))
	for k, v := range wordNo {
		suffixWord[k*3] = v
		binary.LittleEndian.PutUint16(suffixWord[k*3+1:], wordAddr[k])
	}
	copy(data[3:], suffixBit)
	copy(data[3+len(suffixBit):], suffixByte)
	copy(data[3+len(suffixBit)+len(suffixByte):], suffixWord)
	return data
}

// 错误
func responseError(response *ProtocolDataUnit) error {
	toyopucError := &toyopucError{FunctionCode: response.FunctionCode}
	if response.Response != 0 {
		toyopucError.ExceptionCode = response.Response
	}
	return toyopucError
}
