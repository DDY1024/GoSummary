package main

// 运行时类型检查: type-switch

func main() {
	var x interface{}

	// type-switch 类型检查
	switch x.(type) {
	case int:
	default:
	}
}
