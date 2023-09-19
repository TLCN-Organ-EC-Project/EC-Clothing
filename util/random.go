package util

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"
const numbers = "0123456789"
const codes = "ABCDEFGHIJKLMNOPQRSTUVWXYZ123456789"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func RandomFloat(min, max int64) float64 {
	// Tạo một số ngẫu nhiên kiểu int64 trong khoảng từ min đến max.
	randomInt := RandomInt(min, max)

	// Chuyển đổi số nguyên thành float64 và chia cho một giá trị ngẫu nhiên lớn hơn 1 để có số thập phân ngẫu nhiên trong khoảng từ 0 đến 1.
	randomFloat := float64(randomInt) / float64(max)

	// Tạo số float64 trong khoảng từ min đến max.
	return float64(min) + randomFloat*(float64(max)-float64(min))
}

func RandomString(n int, x string) string {
	var sb strings.Builder
	k := len(x)

	for i := 0; i < n; i++ {
		c := x[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

func RandomOwner() string {
	return RandomString(6, alphabet)
}

func RandomPhoneNo() string {
	return RandomString(10, numbers)
}

func RandonEmail() string {
	return fmt.Sprintf("%s@email.com", RandomString(6, alphabet))
}

func RandomProvince() int64 {
	return RandomInt(1, 63)
}

func RandomResetPasswordToken() string {
	return RandomString(5, numbers)
}

func RandomSize() string {
	sizes := []string{S, M, L, XL, XXL, OVERSIZE}
	n := len(sizes)
	return sizes[rand.Intn(n)]
}

func RandomOrderCode() string {
	return RandomString(18, codes)
}