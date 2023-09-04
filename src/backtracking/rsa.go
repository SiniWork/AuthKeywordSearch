/*
author: siyu
date: 2023/03/26
*/

package backtracking

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"github.com/ethereum/go-ethereum/rlp"
	"math/big"
)

type Accumulator struct {
	N *big.Int    // 模数 (saved on chain)
	G *big.Int   // 底数 (saved on chain)
	Ac *big.Int   // RSA累加器中的累积值 (saved on chain)
	pi *big.Int  // RSA累加器中元素的乘积
	Keywords []string  // hash values of keywords
}

func (acc *Accumulator) Initial(keys []string) {
	privKey, _ := rsa.GenerateKey(rand.Reader, 2048)
	acc.N = privKey.PublicKey.N
	acc.G = big.NewInt(37)
	acc.Ac = big.NewInt(1)
	acc.pi = big.NewInt(1)
	for _, key := range keys {
		acc.Keywords = append(acc.Keywords, key)
		acc.Ac.Mul(acc.Ac, Str2BInt(key))
		acc.pi.Mul(acc.pi, Str2BInt(key))
	}
	acc.Ac.Exp(acc.G, acc.pi, acc.N)
}

func (acc *Accumulator) AddEle(hs []byte) {
	h := new(big.Int).SetBytes(hs)
	hPri := GeneratePrime(h)
	acc.pi.Mul(acc.pi, hPri)
	acc.Ac.Exp(acc.Ac, hPri, acc.N)
}

func (acc *Accumulator) DelEle(hs []byte) {
	h := new(big.Int).SetBytes(hs)
	hPri := GeneratePrime(h)
	acc.pi.Div(acc.pi, hPri)
	acc.Ac.Exp(acc.G, acc.pi, acc.N)
}

func (acc *Accumulator) NonMemberProof(k string) (*big.Int, *big.Int){
	// ax + by = 1, a: pi, b: ai, a > b
	// return x, y
	b := Str2BInt(k)
	x, y, _ := ExtEu(acc.pi, b)

	d := new(big.Int).Exp(acc.G, y, acc.N)

	return x, d
}

func ExtEu(a, b *big.Int) (*big.Int, *big.Int, *big.Int) {
	// a > b
	zero := big.NewInt(0)
	one := big.NewInt(1)

	if b.Cmp(zero) == 0 || new(big.Int).Mod(a, b).Cmp(zero) == 0{
		return zero, one, b
	}

	x1, y1, gcd := ExtEu(b, new(big.Int).Mod(a, b))
	x0 := y1
	ty := new(big.Int).Div(a, b)
	y0 := new(big.Int).Sub(x1, ty.Mul(ty, y1))

	return x0, y0, gcd
}

func Test(a *big.Int) {
	if a.Cmp(big.NewInt(0)) == 0 {
		return
	}
	Test(new(big.Int).Sub(a, big.NewInt(1)))
	fmt.Println(a)
}


// 生成不小于给定整数的素数
func GeneratePrime(n *big.Int) *big.Int {
	for i := n; ; i.Add(i, big.NewInt(1)) {
		if isPrime(i) {
			return i
		}
	}
}

// 判断一个数是否为素数
func isPrime(n *big.Int) bool {
	if n.Cmp(big.NewInt(2)) < 0 {
		return false
	}
	if n.Cmp(big.NewInt(2)) == 0 {
		return true
	}
	if n.Bit(0) == 0 {
		return false
	}
	for i := big.NewInt(3); i.Mul(i, i).Cmp(n) <= 0; i.Add(i, big.NewInt(2)) {
		if new(big.Int).Mod(n, i).Cmp(big.NewInt(0)) == 0 {
			return false
		}
	}
	return true
}

func Str2BInt(str string) *big.Int {
	raw := []interface{}{str}
	rlp, _ := rlp.EncodeToBytes(raw)
	h := new(big.Int).SetBytes(Hash(rlp))
	return GeneratePrime(h)
}
