package browsers

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/des"
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/asn1"
	"errors"
	"syscall"
	"unsafe"

	"golang.org/x/crypto/pbkdf2"
)

func DPAPI(encryptPass []byte) ([]byte, error) {
	dllCrypt := syscall.NewLazyDLL("Crypt32.dll")
	dllKernel := syscall.NewLazyDLL("Kernel32.dll")
	procDecryptData := dllCrypt.NewProc("CryptUnprotectData")
	procLocalFree := dllKernel.NewProc("LocalFree")

	type dataBlob struct {
		cbData uint32
		pbData *byte
	}

	var outBlob dataBlob
	var newBlob *dataBlob

	if len(encryptPass) == 0 {
		newBlob = &dataBlob{}
	}
	newBlob = &dataBlob{
		pbData: &encryptPass[0],
		cbData: uint32(len(encryptPass)),
	}
	r, _, err := procDecryptData.Call(uintptr(unsafe.Pointer(newBlob)), 0, 0, 0, 0, 0, uintptr(unsafe.Pointer(&outBlob)))
	if r == 0 {
		return nil, err
	}
	defer procLocalFree.Call(uintptr(unsafe.Pointer(outBlob.pbData)))
	d := make([]byte, outBlob.cbData)
	copy(d, (*[1 << 30]byte)(unsafe.Pointer(outBlob.pbData))[:])
	return d, nil
}

type ASN1PBE interface {
	Decrypt(globalSalt, masterPwd []byte) (key []byte, err error)
}

func NewASN1PBE(b []byte) (pbe ASN1PBE, err error) {
	var (
		n nssPBE
		m metaPBE
		l loginPBE
	)
	if _, err := asn1.Unmarshal(b, &n); err == nil {
		return n, nil
	}
	if _, err := asn1.Unmarshal(b, &m); err == nil {
		return m, nil
	}
	if _, err := asn1.Unmarshal(b, &l); err == nil {
		return l, nil
	}
	return nil, errors.New("decode ASN1 data failed")
}

type nssPBE struct {
	AlgoAttr struct {
		asn1.ObjectIdentifier
		SaltAttr struct {
			EntrySalt []byte
			Len       int
		}
	}
	Encrypted []byte
}

func (n nssPBE) Decrypt(globalSalt, masterPwd []byte) (key []byte, err error) {
	hp := sha1.Sum(append(globalSalt, masterPwd...))
	s := append(hp[:], n.salt()...)
	chp := sha1.Sum(s)
	pes := paddingZero(n.salt(), 20)
	tk := hmac.New(sha1.New, chp[:])
	tk.Write(pes)
	pes = append(pes, n.salt()...)
	k1 := hmac.New(sha1.New, chp[:])
	k1.Write(pes)
	tkPlus := append(tk.Sum(nil), n.salt()...)
	k2 := hmac.New(sha1.New, chp[:])
	k2.Write(tkPlus)
	k := append(k1.Sum(nil), k2.Sum(nil)...)
	iv := k[len(k)-8:]
	return des3Decrypt(k[:24], iv, n.encrypted())
}

func (n nssPBE) salt() []byte {
	return n.AlgoAttr.SaltAttr.EntrySalt
}

func (n nssPBE) encrypted() []byte {
	return n.Encrypted
}

type metaPBE struct {
	AlgoAttr  algoAttr
	Encrypted []byte
}

type algoAttr struct {
	asn1.ObjectIdentifier
	Data struct {
		Data struct {
			asn1.ObjectIdentifier
			SlatAttr slatAttr
		}
		IVData ivAttr
	}
}

type ivAttr struct {
	asn1.ObjectIdentifier
	IV []byte
}

type slatAttr struct {
	EntrySalt      []byte
	IterationCount int
	KeySize        int
	Algorithm      struct {
		asn1.ObjectIdentifier
	}
}

func (m metaPBE) Decrypt(globalSalt, _ []byte) (key2 []byte, err error) {
	k := sha1.Sum(globalSalt)
	key := pbkdf2.Key(k[:], m.salt(), m.iterationCount(), m.keySize(), sha256.New)
	iv := append([]byte{4, 14}, m.iv()...)
	return aes128CBCDecrypt(key, iv, m.encrypted())
}

func (m metaPBE) salt() []byte {
	return m.AlgoAttr.Data.Data.SlatAttr.EntrySalt
}

func (m metaPBE) iterationCount() int {
	return m.AlgoAttr.Data.Data.SlatAttr.IterationCount
}

func (m metaPBE) keySize() int {
	return m.AlgoAttr.Data.Data.SlatAttr.KeySize
}

func (m metaPBE) iv() []byte {
	return m.AlgoAttr.Data.IVData.IV
}

func (m metaPBE) encrypted() []byte {
	return m.Encrypted
}

type loginPBE struct {
	CipherText []byte
	Data       struct {
		asn1.ObjectIdentifier
		IV []byte
	}
	Encrypted []byte
}

func (l loginPBE) Decrypt(globalSalt, _ []byte) (key []byte, err error) {
	return des3Decrypt(globalSalt, l.iv(), l.encrypted())
}

func (l loginPBE) iv() []byte {
	return l.Data.IV
}

func (l loginPBE) encrypted() []byte {
	return l.Encrypted
}

func aes128CBCDecrypt(key, iv, encryptPass []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	encryptLen := len(encryptPass)
	if encryptLen < block.BlockSize() {
		return nil, errors.New("length of encrypted password less than block size")
	}

	dst := make([]byte, encryptLen)
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(dst, encryptPass)
	dst = pkcs5UnPadding(dst, block.BlockSize())
	return dst, nil
}

func pkcs5UnPadding(src []byte, blockSize int) []byte {
	n := len(src)
	paddingNum := int(src[n-1])
	if n < paddingNum || paddingNum > blockSize {
		return src
	}
	return src[:n-paddingNum]
}

func des3Decrypt(key, iv []byte, src []byte) ([]byte, error) {
	block, err := des.NewTripleDESCipher(key)
	if err != nil {
		return nil, err
	}
	blockMode := cipher.NewCBCDecrypter(block, iv)
	sq := make([]byte, len(src))
	blockMode.CryptBlocks(sq, src)
	return pkcs5UnPadding(sq, block.BlockSize()), nil
}

func paddingZero(s []byte, l int) []byte {
	h := l - len(s)
	if h <= 0 {
		return s
	}
	for i := len(s); i < l; i++ {
		s = append(s, 0)
	}
	return s
}
