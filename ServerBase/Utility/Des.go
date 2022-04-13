package Utility

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
)

func padding(src []byte,blocksize int) []byte {
	n:=len(src)
	padnum:=blocksize-n%blocksize
	pad:=bytes.Repeat([]byte{byte(padnum)},padnum)
	dst:=append(src,pad...)
	return dst
}

func unpadding(src []byte) []byte {
	n:=len(src)
	unpadnum:=int(src[n-1])
	dst:=src[:n-unpadnum]
	return dst
}

func EncryptDes(src []byte,key []byte) []byte {
	block,_:=des.NewCipher(key)
	src=padding(src,block.BlockSize())
	blockmode:=cipher.NewCBCEncrypter(block,key)
	blockmode.CryptBlocks(src,src)
	return src
}

func DecryptDes(src []byte,key []byte) []byte {
	block,_:=des.NewCipher(key)
	blockmode:=cipher.NewCBCDecrypter(block,key)
	blockmode.CryptBlocks(src,src)
	src=unpadding(src)
	return src
}

