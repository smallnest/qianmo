package qianmo

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/net/bpf"
)

func TestParseTcpdumpFitler(t *testing.T) {
	// tcpdump -i any  -dd  tcp port 8080
	data := `
	{ 0x28, 0, 0, 0x0000000e },
	{ 0x15, 0, 6, 0x000086dd },
	{ 0x30, 0, 0, 0x00000016 },
	{ 0x15, 0, 15, 0x00000006 },
	{ 0x28, 0, 0, 0x00000038 },
	{ 0x15, 12, 0, 0x00001f90 },
	{ 0x28, 0, 0, 0x0000003a },
	{ 0x15, 10, 11, 0x00001f90 },
	{ 0x15, 0, 10, 0x00000800 },
	{ 0x30, 0, 0, 0x00000019 },
	{ 0x15, 0, 8, 0x00000006 },
	{ 0x28, 0, 0, 0x00000016 },
	{ 0x45, 6, 0, 0x00001fff },
	{ 0xb1, 0, 0, 0x00000010 },
	{ 0x48, 0, 0, 0x00000010 },
	{ 0x15, 2, 0, 0x00001f90 },
	{ 0x48, 0, 0, 0x00000012 },
	{ 0x15, 0, 1, 0x00001f90 },
	{ 0x6, 0, 0, 0x0000ffff },
	{ 0x6, 0, 0, 0x00000000 },
	`

	raws := ParseTcpdumpFitlerData(data)

	_, ok := bpf.Disassemble(raws)
	require.True(t, ok, "parse tcpdump filter failed")

	s := CreateInstructionsFromData(data)
	t.Logf("instruction: \n%s", s)
}

func TestS2Int(t *testing.T) {
	cases := []struct {
		input string
		want  int
	}{
		{"0x10", 16},
		{"  42  ", 42},
		{"-123", -123},
	}
	for _, c := range cases {
		got := s2int(c.input)
		if got != c.want {
			t.Errorf("s2int(%q) == %d, want %d", c.input, got, c.want)
		}
	}
}

func printInstruction(s interface{}) string {
	v := reflect.ValueOf(s)

	// 打印typename
	t := v.Type()
	str := fmt.Sprintf("bpf.%s{", t.Name())

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

		if t.Name() == "RetConstant" && f.Name == "Val" {
			str += fmt.Sprintf("%s: 0x%x,", f.Name, val)
			continue
		}

		str += fmt.Sprintf("%s: %v,", f.Name, val)
	}

	str = strings.TrimSuffix(str, ",")
	str += fmt.Sprintf("}")

	return str
}

func TestParseExpr(t *testing.T) {
	raws, err := ParseTcpdumpFitlerExpr("tcp port 8080")
	require.NoError(t, err)

	s := CreateInstructionsFromRaws(raws)
	t.Log(s)

	raws, err = ParseTcpdumpFitlerExpr("dst host 8.8.8.8 and icmp")
	require.NoError(t, err)

	s = CreateInstructionsFromRaws(raws)
	t.Log(s)
}
