package bpf

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"golang.org/x/net/bpf"
)

var jumpTests = map[bpf.JumpTest]string{
	bpf.JumpEqual:          "JumpEqual",
	bpf.JumpNotEqual:       "JumpNotEqual",
	bpf.JumpGreaterThan:    "JumpGreaterThan",
	bpf.JumpLessThan:       "JumpLessThan",
	bpf.JumpGreaterOrEqual: "JumpGreaterOrEqual",
	bpf.JumpLessOrEqual:    "JumpLessOrEqual",
	bpf.JumpBitsSet:        "JumpBitsSet",
	bpf.JumpBitsNotSet:     "JumpBitsNotSet",
}

// ParseTcpdumpFitler parses tcpdump filter to bpf.RawInstruction.
// Example:
// tcpdump -i eth0 -dd 'tcp and port 80'
func ParseTcpdumpFitlerData(data string) (raws []bpf.RawInstruction) {
	lines := strings.Split(data, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		line = strings.TrimPrefix(line, "{")
		line = strings.TrimSuffix(line, " },")
		items := strings.Split(line, ",")
		// assert len(items) == 4

		raw := bpf.RawInstruction{
			Op: uint16(s2int(items[0])),
			Jt: uint8(s2int(items[1])),
			Jf: uint8(s2int(items[2])),
			K:  uint32(s2int(items[3])),
		}

		raws = append(raws, raw)
	}

	return raws
}

// ParseTcpdumpFitlerExpr parses tcpdump filter to bpf.RawInstruction.
func ParseTcpdumpFitlerExpr(linkType layers.LinkType, expr string) (raws []bpf.RawInstruction, err error) {
	insts, err := pcap.CompileBPFFilter(linkType, 1, expr)
	if err != nil {
		return nil, err
	}

	raws = make([]bpf.RawInstruction, 0, len(insts))
	for _, inst := range insts {
		raws = append(raws, bpf.RawInstruction{
			Op: inst.Code,
			Jt: inst.Jt,
			Jf: inst.Jf,
			K:  inst.K,
		})
	}

	return raws, nil
}

// CreateInstructionsFromData creates bpf.Instruction from tcpdump filter.
func CreateInstructionsFromData(data string) string {
	raws := ParseTcpdumpFitlerData(data)
	insts, _ := bpf.Disassemble(raws)

	var filter = "var filter = []bpf.Instruction {\n"

	var instStrs = make([]string, 0, len(insts))
	for _, inst := range insts {
		instStr := createInstruction(inst)
		instStrs = append(instStrs, instStr)

		filter += fmt.Sprintf("\t%s,\n", instStr)
	}
	filter += "}"

	return filter
}

// CreateInstructionsFromExpr creates bpf.Instruction from tcpdump filter expression.
func CreateInstructionsFromExpr(linkType layers.LinkType, expr string) string {
	raws, _ := ParseTcpdumpFitlerExpr(linkType, expr)

	return CreateInstructionsFromRaws(raws)
}

// CreateInstructionsFromRaws creates bpf.Instruction from bpf.RawInstruction.
func CreateInstructionsFromRaws(raws []bpf.RawInstruction) string {
	insts, _ := bpf.Disassemble(raws)

	var filter = "var filter = []bpf.Instruction {\n"

	var instStrs = make([]string, 0, len(insts))
	for _, inst := range insts {
		instStr := createInstruction(inst)
		instStrs = append(instStrs, instStr)

		filter += fmt.Sprintf("\t%s,\n", instStr)
	}
	filter += "}"

	return filter
}

// createInstruction creates bpf.Instruction from bpf.RawInstruction.
func createInstruction(s bpf.Instruction) string {
	v := reflect.ValueOf(s)

	// 打印typename
	t := v.Type()
	str := fmt.Sprintf("bpf.%s{", t.Name())

	var comment string
	// 遍历结构体字段
	for i := 0; i < v.NumField(); i++ {
		// 获取每个字段的结构体Field
		f := t.Field(i)
		val := v.Field(i).Interface()

		// 拼接字段名和值到字符串
		if f.Name == "SkipTrue" || f.Name == "SkipFalse" {
			if v, ok := val.(uint8); ok && v == 0 {
				continue
			}
		}

		// 拼接字段名和值到字符串
		if t.Name() == "JumpIf" && f.Name == "Val" {
			if v, ok := val.(uint32); ok && v > 15 {
				comment += fmt.Sprintf(" %d = 0x%x,", val, val)
			}
		}

		if t.Name() == "JumpIf" && f.Name == "Cond" {
			if v, ok := val.(bpf.JumpTest); ok {
				str += fmt.Sprintf("%s: bpf.%s,", f.Name, jumpTests[bpf.JumpTest(v)])
				continue
			}
		}

		if t.Name() == "RetConstant" && f.Name == "Val" {
			str += fmt.Sprintf("%s: 0x%x,", f.Name, val)
			continue
		}

		str += fmt.Sprintf("%s: %v,", f.Name, val)
	}

	str = strings.TrimSuffix(str, ",")
	str += fmt.Sprintf("}")

	if comment != "" {
		comment = strings.TrimSuffix(comment, ",")
		str += fmt.Sprintf(" // %s", comment)
	}

	return str
}

func s2int(s string) int {
	s = strings.TrimSpace(s)
	if strings.HasPrefix(s, "0x") {
		result, _ := strconv.ParseInt(s, 0, 64)
		return int(result)
	}

	result, _ := strconv.ParseInt(s, 10, 64)
	return int(result)
}
