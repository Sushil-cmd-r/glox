package vm

type Opcode = byte

const (
	OpReturn Opcode = iota
	OpConstant
	OpAdd
	OpSub
	OpMul
	OpDiv
	OpNegate
	OpNot
	OpNil
	OpPop
	OpPrint
	OpCall
	OpDefineGlobal
	OpSetGlobal
	OpGetGlobal
	OpGetLocal
	OpSetLocal
)

var opcodes = [...]string{
	OpReturn:       "OpReturn",
	OpConstant:     "OpConstant",
	OpAdd:          "OpAdd",
	OpSub:          "OpSub",
	OpMul:          "OpMul",
	OpDiv:          "OpDiv",
	OpNegate:       "OpNegate",
	OpNot:          "OpNot",
	OpNil:          "OpNil",
	OpPop:          "OpPop",
	OpPrint:        "OpPrint",
	OpCall:         "OpCall",
	OpDefineGlobal: "OpDefineGlobal",
	OpSetGlobal:    "OpSetGlobal",
	OpGetGlobal:    "OpGetGlobal",
	OpGetLocal:     "OpGetLocal",
	OpSetLocal:     "OpSetLocal",
}
