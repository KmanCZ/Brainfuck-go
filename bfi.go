package main

import (
    "fmt"
    "bufio"
    "os"
    "errors"
)

const (
    INC_POS = iota
    DEC_POS
    INC_VAL
    DEC_VAL
    OUTPUT
    INPUT
    LOOP_START
    LOOP_END
)

const dataSize = 30000

type instruction struct {
    operation uint8
    operand uint
}

func check(err error) {
    if err != nil {
        panic(err)
    }
}

func interpret(program *[]instruction) {
    data := make([]byte, dataSize)
    var dataPtr uint

    reader := bufio.NewReader(os.Stdin)

    for pc := 0; pc < len(*program); pc++ {
        switch (*program)[pc].operation {
            case INC_POS:
                dataPtr++
            case DEC_POS:
                dataPtr--
            case INC_VAL:
                data[dataPtr]++
            case DEC_VAL:
                data[dataPtr]--
            case OUTPUT:
                fmt.Printf("%c", data[dataPtr])
            case INPUT:
                inputByte, _ := reader.ReadByte()
                data[dataPtr] = inputByte
            case LOOP_START:
                if data[dataPtr] == 0 {
                    pc = int((*program)[pc].operand)
                }
            case LOOP_END:
                if data[dataPtr] != 0 {
                    pc = int((*program)[pc].operand)
                }
            default:
                panic("Unknown operation!")
        }
    }
}

func compile(path string) (*[]instruction, error) {
    f, err := os.Open(path)
    check(err)
    defer f.Close()
 
    scaner := bufio.NewScanner(f)
    var pc, jmpPc uint 
    program := make([]instruction, 0)
    jmpStack := make([]uint, 0)

    for scaner.Scan() {
        for _, symbol := range scaner.Text() {
            switch symbol {
                case '>':
                    program = append(program, instruction{INC_POS, 0})
                case '<':
                    program = append(program, instruction{DEC_POS, 0})
                case '+':
                    program = append(program, instruction{INC_VAL, 0})
                case '-':
                    program = append(program, instruction{DEC_VAL, 0})
                case '.':
                    program = append(program, instruction{OUTPUT, 0})
                case ',':
                    program = append(program, instruction{INPUT, 0})
                case '[':
                    program = append(program, instruction{LOOP_START, 0})
                    jmpStack = append(jmpStack, pc)
                case ']':
                    if len(jmpStack) == 0 {
                        return nil, errors.New("Compilation error")
                    }
                    jmpPc = jmpStack[len(jmpStack)-1]
                    jmpStack = jmpStack[:len(jmpStack)-1]
                    program = append(program, instruction{LOOP_END, jmpPc})
                    program[jmpPc].operand = pc
                default:
                    continue
            }
            pc++
        }
    }
    
    return &program, nil
}

func main() {
    path := os.Args[1]
    program, err := compile(path)
    check(err)
    interpret(program)
}
