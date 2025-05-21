package utils

// import "fmt"

func Dot(a []float32, b []float32) float32 {
	if len(a) != len(b) || len(a) == 0 {
		return 0
	}
	var sum float32
	for i := range a {
		sum += a[i] * b[i]
	}
	// fmt.Printf("The dot product is :%v", sum)
	return sum
}
