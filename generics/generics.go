package generics

import "fmt"

/*
TODO Go 泛型：本教程介绍了 Go 语言中泛型的基础知识。借助泛型，你可以声明和使用函数或类型，这些函数或类型被编写为与调用代码提供的一组类型中的任何一个一起工作
TODO 在本教程中，你将声明两个简单的非泛型函数，然后在一个泛型函数中捕获相同的逻辑
TODO 将逐步完成以下部分：
- 为你的代码创建一个文件夹 generics
- 添加非泛型函数 SumInts、SumFloats
- 添加一个泛型函数以处理多种类型 SumIntsOrFloats
- 在调用泛型函数时移除类型参数，编译器可以推断
- 声明类型约束，SumNumbers 可以通过捕获整数和浮点数的并集到一个可重用的类型约束中（例如从其他代码中重用）来进一步简化该函数

TODO 1、添加非泛型函数：添加两个函数，每个函数将 map 的值相加并返回总和，声明了两个函数而不是一个，因为你使用了两种不同类型的映射：一种存储 int 64 值，另一种存储 float 64 值
*/

// 声明 两个函数将 map 的值相加并返回总和
// TODO SumInts 将字符串映射到 int64 值
func SumInts(m map[string]int64) int64 {
	var s int64 = 0
	// 遍历 map[string]int64 获取 k, v
	for _, v := range m {
		s += v
	}
	return s
}

// TODO SumFloats 将字符串映射到 float64 值
func SumFloats(m map[string]float64) float64 {
	var s float64 = 0
	// 遍历 map[string]float64 获取 k, v
	for _, v := range m {
		s += v
	}
	return s
}

// TODO 2、使用泛型，可以在这里编写一个函数而不是两个。接下来，将为包含整数或浮点值的映射添加一个单一的泛型函数
// TODO 将添加一个通用函数，该函数可以接收一个包含整数或浮点值的映射，有效地用一个函数替换你刚刚编写的两个函数

// TODO 为了支持这两种类型的值，那个单一函数将需要一种方法来声明它支持哪些类型。另一方面，调用代码将需要一种方法来指定它是使用整数映射还是浮点数映射进行调用
// TODO 为了支持这一点，将编写一个函数，该函数除了普通函数参数外还声明类型参数。这些类型参数使函数具有通用性，使其能够处理不同类型的参数。你将使用类型实参和普通函数参数调用该函数
// TODO 每个类型参数都有一个类型约束，它对于类型参数而言就像是一种元类型。每个类型约束指定了调用代码可以为相应类型参数使用的允许类型参数
// TODO 虽然类型参数的约束通常表示一组类型，但在编译时，类型参数代表单个类型 —— 即调用代码作为类型参数提供的类型。如果类型参数的约束不允许类型参数的类型，则代码将无法编译
// TODO 类型参数必须支持泛型代码对其执行的所有操作。例如，如果你的函数代码试图对一个类型参数执行字符串操作（比如索引），而该类型参数的约束包括数字类型，那么代码将无法编译

// TODO 3、使用一个允许整型或浮点型的约束，创建通用函数，SumIntsOrFloats 对映射 m 的值进行求和。它支持 int64 和 float64 作为映射值的类型
// TODO 声明一个带有两个类型参数（在方括号内）的 SumIntsOrFloats 函数，K 和 V 以及一个使用类型参数的参数类型为 map[K]V 的 m，该函数返回一个类型为 V 的值
// TODO 为 K 类型参数指定类型约束 comparable 专门针对此类情况，comparable 约束在 Go 中是预先声明的。允许任何其值可以用作比较运算符 == 和 != 的操作数的类型。Go 要求映射键是可比较的。因此，将 K 声明为 comparable 是必要的，这样你就可以将 K 用作映射变量中的键。它还确保调用代码为映射键使用允许的类型
// TODO 为 V 类型参数指定一个约束，该约束是两种类型的联合：int64 和 float64 使用 | 指定这两种类型的联合，这意味着此约束允许任一种类型。编译器将允许任一种类型作为调用代码中的参数
// TODO 指定 m 参数的类型为 map[K]V 其中 K 和 V 是已经为类型参数指定的类型。请注意，我们知道 map[K]V 是有效的映射类型，因为 K 是可比较类型。如果我们没有声明 K 为可比较类型，编译器将拒绝引用 map[K]V
func SumIntsOrFloats[K comparable, V int64 | float64](m map[K]V) V {
	var s V
	// 遍历 map[K]V 获取 k, v
	for _, v := range m {
		s += v
	}
	return s
}

// TODO 4、声明类型约束：将先前定义的约束移入其自己的接口中，以便可以在多个地方重复使用它。以这种方式声明约束有助于简化代码，例如当约束更复杂时
// TODO 将类型约束声明为一个接口。该约束允许实现此接口的任何类型。例如，如果声明一个具有三个方法的类型约束接口，然后在泛型函数中使用带有类型参数的该接口，用于调用函数的类型实参必须具有所有那些方法
// TODO 约束接口也可以引用特定类型，正如你将在本节中看到的那样
// TODO 声明 Number 接口类型以用作类型约束
type Number interface {
	// TODO 在接口内部声明 int64 和 float64 的联合，本质上，正在将联合类型从函数声明移到一个新的类型约束中。这样，当你想将类型参数约束为 int64 或 float64 时，你可以使用这个 Number 类型约束，而不是写出 int64 | float64
	int64 | float64
}

// TODO 声明一个泛型函数，该函数具有与你之前声明的泛型函数相同的逻辑，但使用新的接口类型而不是联合类型作为类型约束。和之前一样，你将类型参数用于参数和返回类型
func SumNumbers[K comparable, V Number](m map[K]V) V {
	var s V
	for _, v := range m {
		s += v
	}
	return s
}

func SumPrint() {
	// TODO 在 Go 语言中，切片（slice） 和 映射（map） 是两种非常重要的数据结构。它们的声明和使用方式各有特点
	// TODO 切片： 是对数组的抽象，它是一个动态数组，长度可以动态调整。切片的底层实际上是一个指向数组的指针、长度和容量
	// TODO 映射：是一个键值对集合，类似于其他语言中的字典或哈希表。Go 中的 map 是无序的，键必须是可比较的类型（如 int、string 等），值可以是任意类型
	/*
		切片（Slice）
		- 动态长度，可以通过 append 扩展
		- 底层是数组，具有容量限制
		- 支持切片操作（如 arr[low:high]）
		- 常见操作：
			初始化：[]Type{...} 或 make([]Type, length, capacity)
			添加元素：slice = append(slice, element)
			获取长度：len(slice)
			获取容量：cap(slice)

		映射（Map）
		- 键值对集合，键必须唯一
		- 无序存储
		- 常见操作：
			初始化：map[KeyType]ValueType{} 或 make(map[KeyType]ValueType)
			添加或修改元素：map[key] = value
			访问元素：value := map[key]
			删除元素：delete(map, key)
			检查键是否存在：value, exists := map[key]
	*/

	// TODO 初始化一个 float64 值映射和一个 int64 值映射，每个映射都有两个条目
	// 初始化一个用于整数值的映射map
	ints := map[string]int64{
		"first":  34,
		"second": 12,
	}
	// 初始化一个用于存储浮点值的映射map
	floats := map[string]float64{
		"first":  35.98,
		"second": 26.99,
	}

	// TODO 调用前面声明的两个函数来求每个 map 值的和，打印结果
	// fmt.Printf("SumInts: %d, SumFloats: %f\n", SumInts(ints), SumFloats(floats))
	fmt.Printf("Non-Generic Sums: %v and %v\n", SumInts(ints), SumFloats(floats))

	// TODO 调用刚刚声明的泛型函数 SumIntsOrFloats，传递创建的每个映射，指定类型参数（方括号中的类型名称），以明确在调用的函数中应替换类型参数的类型，通常可以省略函数调用中的类型参数。Go 通常可以从代码中推断出它们
	// 在每次调用中，编译器都将类型参数替换为该调用中指定的具体类型，在调用您编写的泛型函数时，您指定了类型参数，这些参数告诉编译器使用什么类型来代替函数的类型参数。因为编译器可以推断它们
	fmt.Printf("Generic Sums: %v and %v\n", SumIntsOrFloats[string, int64](ints), SumIntsOrFloats[string, float64](floats))
	// TODO 调用泛型函数时移除类型参数，当 Go 编译器可以推断出你想要使用的类型时，你可以在调用代码中省略类型参数。编译器从函数参数的类型推断类型参数，
	// TODO 请注意，这并不总是可行的。例如，如果您需要调用一个没有参数的泛型函数，则需要在函数调用中包含类型参数
	fmt.Printf("Generic Sums, type parameters inferred: %v and %v\n", SumIntsOrFloats(ints), SumIntsOrFloats(floats))
	// TODO 调用 SumNumbers 函数处理每个映射，并打印每个映射值的总和，与上一行一样，在调用泛型函数时省略类型参数（方括号中的类型名称）。Go 编译器可以从其他参数中推断出类型参数
	fmt.Printf("Generic Sums with Constraint: %v and %v\n", SumNumbers(ints), SumNumbers(floats))
}
