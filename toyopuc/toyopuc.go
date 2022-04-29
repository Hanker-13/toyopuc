package toyopuc

/*
Package toyopuc provides a client for TOYOPUC.
*/

import (
	"fmt"
)

// 帧格式
const (
	// 请求帧
	RequestFTByte = 0x00
	// 响应帧
	ResponseFTByte = 0x80
)

// 功能码
const (
	FunSequentialProgramReadWord  = 0x18
	FunSequentialProgramWriteWord = 0x19
	FunIOReadWord                 = 0x1C
	FunIOWriteWord                = 0x1D
	FunIOReadByte                 = 0x1E
	FunIOWriteByte                = 0x1F
	FunIOReadBit                  = 0x20
	FunIOWriteBit                 = 0x21
	FunIOReadMultipointWord       = 0x22
	FunIOWriteMultipointWord      = 0x23
	FunIOReadMultipointByte       = 0x24
	FunIOWriteMultipointByte      = 0x25
	FunIOReadMultipointBit        = 0x26
	FunIOWriteMultipointBit       = 0x27

	// 扩展
	FunProgramExpansionReadWord    = 0x90
	FunProgramExpansionWriteWord   = 0x91
	FunDataExpansionReadWord       = 0x94
	FunDateExpansionWriteWord      = 0x95
	FunDataExpansionReadByte       = 0x96
	FunDataExpansionWriteByte      = 0x97
	FunDataExpansionReadMultipoint = 0x98

	// TODO
	FunDataExpansionWriteMultipoint = 0x99
)

// 错误码
const (
	ExceptionCodeHardwareAbnormalityOfCPU                    = 0x11
	ExceptionCodeIllegalENQ                                  = 0x20
	ExceptionCodeAbnormalTransmissionQuantity                = 0x21
	ExceptionCodeIllegalCommandCode                          = 0x23
	ExceptionCodeIllegalSubcommandCode                       = 0x24
	ExceptionCodeIllegalDataByteInCommandFormat              = 0x25
	ExceptionCodeIllegalNumberOfFunctionCallOperands         = 0x26
	ExceptionCodeWriteForbiddenInArea                        = 0x31
	ExceptionCodeDisableCommandWithStopDuration              = 0x32
	ExceptionCodeDebugFunctionWithNotDebugMode               = 0x33
	ExceptionCodeNoAccessByAccessProhibitionSetting          = 0x34
	ExceptionCodeCannotExecWithoutPermission                 = 0x35
	ExceptionCodeCannotExecWithoutPermissionByOtherDeviceSet = 0x36
	ExceptionCodeNoResetAfterWriteIOParams                   = 0x39
	ExceptionCodeUnenforceableCommandWithSeriousFault        = 0x3C
	ExceptionCodeConflictWithOtherCommand                    = 0x3D
	ExceptionCodeConnotExecWithReset                         = 0x3E
	ExceptionCodeConnotExecWithStopStatus                    = 0x3F
	ExceptionCodeAddressNotInRange                           = 0x40
	ExceptionCodeNumOutOfRange                               = 0x41
	ExceptionCodeDataOtherThanSpecified                      = 0x42
	ExceptionCodeErrorInOperandFunctionCall                  = 0x43
	ExceptionCodeCommandWithOutTimerCounter                  = 0x52
	ExceptionCodeNoAnswer                                    = 0x66
	ExceptionCodeCommandNotUsed                              = 0x70
	ExceptionCodeDataModeNoAnswer                            = 0x72
	ExceptionCodeCommandCannotProceed                        = 0x73
)

// toyopucError implements error interface.
// 实现了错误接口
type toyopucError struct {
	FunctionCode  byte
	ExceptionCode byte
}

// Error converts known modbus exception code to error message.
// 错误 将已知的modbus异常代码转换为错误消息
func (e *toyopucError) Error() string {
	var name string
	switch e.ExceptionCode {
	case ExceptionCodeHardwareAbnormalityOfCPU:
		name = "hardware abnormality of CPU"
	case ExceptionCodeIllegalENQ:
		name = "ENQ connot be 5"
	case ExceptionCodeAbnormalTransmissionQuantity:
		name = "abnormal transmission quantity"
	case ExceptionCodeIllegalCommandCode:
		name = "illegal command code"
	case ExceptionCodeIllegalSubcommandCode:
		name = "illegal subcommand code"
	case ExceptionCodeIllegalDataByteInCommandFormat:
		name = "Illegal data byte in command format"
	case ExceptionCodeIllegalNumberOfFunctionCallOperands:
		name = "number of illegal function call operands"
	case ExceptionCodeWriteForbiddenInArea:
		name = "an attempt was made to write in an area where writing is prohibited in a sequential operation"
	case ExceptionCodeDisableCommandWithStopDuration:
		name = "during the stop duration, a disable command is sent"
	case ExceptionCodeDebugFunctionWithNotDebugMode:
		name = "although not in debug mode, a debug function call was attempted"
	case ExceptionCodeNoAccessByAccessProhibitionSetting:
		name = "access is prohibited through access prohibition setting"
	case ExceptionCodeCannotExecWithoutPermission:
		name = "because the execution permission is set, it cannot be executed"
	case ExceptionCodeCannotExecWithoutPermissionByOtherDeviceSet:
		name = "since the execution permission is set through other devices, it cannot be executed"
	case ExceptionCodeNoResetAfterWriteIOParams:
		name = "after I/O point parameters and I / O allocated point parameters are written, try to start scanning without resetting"
	case ExceptionCodeUnenforceableCommandWithSeriousFault:
		name = "during a serious failure, an unenforceable command was issued"
	case ExceptionCodeConflictWithOtherCommand:
		name = "the command cannot be executed because a reset is in progress"
	case ExceptionCodeConnotExecWithReset:
		name = "in the command execution of other factors, the processing will conflict, so it is not executable"
	case ExceptionCodeConnotExecWithStopStatus:
		name = "the command cannot be executed because it is stopped"
	case ExceptionCodeAddressNotInRange:
		name = "the address is not in the range due to reading and writing commands, or the address + data quantity of the command deviates from the address range"
	case ExceptionCodeNumOutOfRange:
		name = "the number of words or bytes is out of range"
	case ExceptionCodeErrorInOperandFunctionCall:
		name = "there is an error in the operand of the function call"
	case ExceptionCodeCommandWithOutTimerCounter:
		name = "although the timer and counter are not used, the commands of reading and writing the set value and current value are still sent"
	case ExceptionCodeNoAnswer:
		name = "there is no response from the link module of the link No. and station number specified by the relay command. (the specified link module does not exist or the power supply is off, or the line is abnormal, etc.)"
	case ExceptionCodeCommandNotUsed:
		name = "the module of link No. specified by the relay command cannot be used (the error of the specified link No. or the exception of the link module)"
	case ExceptionCodeDataModeNoAnswer:
		name = "there is no response from the link module of the link No. and station number specified by the relay command. (the specified link module does not exist or the power supply is off, or the line is abnormal, etc.)"
	case ExceptionCodeCommandCannotProceed:
		name = "since multiple relay commands are repeatedly sent to the same link module in the CPU module, the command processing cannot be carried out. (please send the command again)"
	default:
		name = "unknown"
	}
	return fmt.Sprintf("toyopuc: exception '%v' (%s), function '%v'", e.ExceptionCode, name, e.FunctionCode)
}

// ProtocolDataUnit (PDU) is independent of underlying communication layers.
// 独立于底层通信层
type ProtocolDataUnit struct {
	FunctionCode byte
	Data         []byte
	Response     byte
}

// Packager specifies the communication layer.
// 通信层
type Packager interface {
	Encode(pdu *ProtocolDataUnit) (adu []byte, err error)
	Decode(adu []byte) (pdu *ProtocolDataUnit, err error)
	Verify(aduRequest []byte, aduResponse []byte) (err error)
}

// Transporter specifies the transport layer.
// 传输层
type Transporter interface {
	Send(aduRequest []byte) (aduResponse []byte, err error)
}
