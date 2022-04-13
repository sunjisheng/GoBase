package Utility

import (
	"../Log"
	"../Math"
	"fmt"
	"math/rand"
	"net"
	"os"
	"path/filepath"
	"runtime/debug"
	"strconv"
	"strings"
	"syscall"
)

func StartRoutine(fn func()) {
	defer func(){
		if r := recover();r!=nil{
			Log.WriteLog(Log.Log_Level_Error, "Catch Panic: %s", r)
			Log.WriteLog(Log.Log_Level_Error, "%s", string(debug.Stack()))
		}
	}()
	fn()
}

func StartRoutine_Arg(fn func(arg interface{}), arg interface{}) {
	defer func(){
		if r := recover();r!=nil{
			Log.WriteLog(Log.Log_Level_Error, "Catch Panic: %s", r)
			Log.WriteLog(Log.Log_Level_Error, "%s", string(debug.Stack()))
		}
	}()
	fn(arg)
}

func ReadUInt32(buf []byte, offset uint32) uint32 {
	return uint32(buf[offset]<<24) | uint32(buf[offset+1])<<16 | uint32(buf[offset+2])<<8 | uint32(buf[offset+3])
}

func ReadUInt32_Little(buf []byte, offset uint32) uint32 {
	return uint32(buf[offset]) | uint32(buf[offset+1])<<8 | uint32(buf[offset+2])<<16 | uint32(buf[offset+3])<<24
}

func WriteUInt32(buf []byte, offset uint32, value uint32){
	buf[offset] = byte(value >>24)
	buf[offset+1] = byte(value >>16)
	buf[offset+2] = byte(value >>8)
	buf[offset+3] = byte(value)
}

func ReadUInt64(buf []byte, offset uint32) uint64 {
	return uint64(buf[offset])<<56 | uint64(buf[offset+1])<<48 | uint64(buf[offset+2])<<40 | uint64(buf[offset+3])<<32 |
		uint64(buf[offset+4])<<24| uint64(buf[offset+5])<<16| uint64(buf[offset+6])<<8| uint64(buf[offset+7])
}

func WriteUInt64(buf []byte, offset uint32, value uint64){
	buf[offset] = byte(value>>56)
	buf[offset+1] = byte(value >>48)
	buf[offset+2] = byte(value >>40)
	buf[offset+3] = byte(value >>32)
	buf[offset+4] = byte(value >>24)
	buf[offset+5] = byte(value >>16)
	buf[offset+6] = byte(value >>8)
	buf[offset+7] = byte(value)
}

func GetCurrentPath() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return ""
	}
	return dir
}

func GetIPPort(sa syscall.Sockaddr) (string,uint32){
	var ip string
	var port uint32
	switch inst:=sa.(type){
	case *syscall.SockaddrInet4:
		ip =net.IP(inst.Addr[:]).String()
		port = uint32(inst.Port)
	}
	return ip,port
}

func SplitIPPort(addr string) (string,uint32) {
	var ip string = ""
	var port uint32 = 0
	arr := strings.Split(addr, ":")
	if len(arr) == 2 {
		ip = arr[0]
		_port,_ := strconv.Atoi(arr[1])
		port = uint32(_port)
	}
	return ip,port
}

func MakeWord(low, high uint8) uint32 {
	var ret uint16 = uint16(high)<<8 + uint16(low)
	return uint32(ret)
}

//字符串IP转字节数组
func Inet_Addr(ipaddr string) [4]byte {
	var (
		ips = strings.Split(ipaddr, ".")
		ip  [4]uint64
		ret [4]byte
	)
	for i := 0; i < 4; i++ {
		ip[i], _ = strconv.ParseUint(ips[i], 10, 8)
	}
	for i := 0; i < 4; i++ {
		ret[i] = byte(ip[i])
	}
	return ret
}

//字符串哈希
func HashStr(str string) uint64 {
	var hash uint64 = 5381
	for i := 0; i < len(str); i++ {
		hash = ((hash << 5) + hash) + uint64(str[i])
	}
	return hash
}

//浮点数保留指定位数的小数
func FloatRound(f float32, n int) float32 {
	format := "%." + strconv.Itoa(n) + "f"
	res, _ := strconv.ParseFloat(fmt.Sprintf(format, f), 64)
	return float32(res)
}

//切分字符串为整数数组
func SplitInt(str string, sep string) []uint32 {
	strarr := strings.Split(str, sep)
	result := make([]uint32, len(strarr))
	for i := 0; i < len(strarr); i++ {
		value, err := strconv.Atoi(strarr[i])
		if err != nil {
			break
		}
		result[i] = uint32(value)
	}
	return result
}

//多次随机
func MultiRandom (max int, count int, repeat bool) []int{
	result := make([]int, count)
	if repeat {
		for i := 0; i < count; i++ {
			result[i] = rand.Intn(int(max))
		}
		return result
	} else {
		arr := rand.Perm(int(max))
		return arr[0:count]
	}
}

//判断整数是否在slice里
func Exist(arr []uint32, i uint32) bool{
	for _, v := range arr {
		if v == i {
			return true
		}
	}
	return false
}

//字符串转32位整数
func Str2UInt32(str string) uint32 {
	i, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}
	return uint32(i)
}

//字符串转64位整数
func Str2UInt64(str string) uint64 {
	i, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}
	return uint64(i)
}

//字符串转Float
func Str2Float(str string) float32 {
	f, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return 0
	}
	return float32(f)
}

//字符串转Vector3
func Str2Vector3(str string, sep string) *Math.Vector3 {
	ve := new(Math.Vector3)
	arr := strings.Split(str, sep)
	if len(arr) != 3 {
		return ve
	}
	ve.X = Str2Float(arr[0])
	ve.Y = Str2Float(arr[1])
	ve.Z = Str2Float(arr[2])
	return ve
}

