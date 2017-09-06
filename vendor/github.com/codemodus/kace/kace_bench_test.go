package kace

import (
	"testing"
)

func BenchmarkCamel4(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_ = Camel("this_is_a_test")
	}
}

func BenchmarkPascal4(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_ = Pascal("this_is_a_test")
	}
}

func BenchmarkSnake4(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_ = Snake("ThisIsATest")
	}
}

func BenchmarkSnakeUpper4(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_ = SnakeUpper("ThisIsATest")
	}
}

func BenchmarkKebab4(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_ = Kebab("ThisIsATest")
	}
}

func BenchmarkKebabUpper4(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_ = KebabUpper("ThisIsATest")
	}
}
