# 第九章 范型
`Golang` 1.18 已经开始默认支持范型了。通过范型可以写出支持多种类型的函数和类型。目前golang范型支持以下特性：

- 范型函数，支持类型参数列表，如 `func function[T any] (param K){}`
- 范型函数，类型列表中的类型，可以在函数体中使用，定义变量或者使用接口方法
- 范型类型，支持有类型参数， 如`type M[T any]{ value T}`
- 类型参数，必须有类型约束，比如`func F[T constraint](param T){}`
- 类型约束，必须是`interface`类型
- 范型函数，接收的参数必须满足对应的类型约束
- 范型函数/范型类型，在使用或者实例化时，需要明确类型列表中的类型；
- 范型函数/范型类型，大部分时间，可以不用指定类型列表中的类型，编译器会自动判断

范型将有效减少`golang`中空interface的使用，尤其容器类型中的场景。

## 9.1 基本语法

如下所示的代码，定义了一个范型函数，格式为：`func functionName[类型列表] (参数列表) { 函数体 }`

```go
func Print[T any](slice []T) {
	fmt.Print("Generic : ")
	for _, item := range slice {
		fmt.Printf("%v ", item)
	}
	fmt.Println()
}
```

可以使用如下方式使用：

```go
t.Run("test integer slice", func(t *testing.T) {
		Print[int]([]int{1, 2, 3, 4, 5, 6, 7})
		t.Logf("int function address : %p", Print[int])
	})
	t.Run("test string slice", func(t *testing.T) {
		Print([]string{"one", "two", "three", "four", "five", "six", "seven"})
		t.Logf("string function address : %p", Print[string])
	})
	t.Run("test float64 slice", func(t *testing.T) {
		Print[float64]([]float64{1, 2, 3, 4, 5, 6, 7})
		t.Logf("float64 function address : %p", Print[float64])
	})
```


```
t.Run("function address", func(t *testing.T) {
		f1 := Print[int]
		f2 := Print[float64]
		f3 := Print[string]
		//f4 := Print //不能在不指定参数的情况下，索引范型函数

		t.Logf("int function address : %p, float64 function address : %p , string function address : %p\n", f1, f2, f3)
	})
```

会输出结果：

```
Generic : 1 2 3 4 5 6 7
    basic_test.go:8: int function address : 0x10ef360
=== RUN   TestBasic/test_string_slice
Generic : one two three four five six seven
    basic_test.go:12: string function address : 0x10ef3c0
=== RUN   TestBasic/test_float64_slice
Generic : 1 2 3 4 5 6 7
    basic_test.go:16: float64 function address : 0x10ef420
=== RUN   TestBasic/function_address
    basic_test.go:25: int function address : 0x10ef480, float64 function address : 0x10ef4e0 , string function address : 0x10ef540
--- PASS: TestBasic (0.00s)
```

可以看出：
- 相同类型的范型函数是不一样的，类似于每次新建都都新建了一个函数对象
- 可以显示指定类型列表；
- 范型函数，必须指定类型或者编译器可以把所有类型推导出来，否则将不能使用和索引。


## 9.2 依赖类型

如下所示， 可以定义范型类型，格式为：`type 范型类型[类型1 约束，类型2 约束]` 具体的类型如（[]T） ：

```go
type vector[T any] []T

func (v vector[T]) last() (T, error) {
	if len(v) > 0 {
		return v[len(v)-1], nil
	}
	//var zero T 先定义
	//return T{}, errors.New("Empty") //如果是string，int，bool等会有问题
	return *new(T), errors.New("Empty")
}
```

如上所示的代码，对于某种情况，可能依赖于类型列表中的类型创建对象，可以像使用具体类型一样创建对象：
- `var zero T` 预先定义
- `return T{}, errors.New("Empty")` 对于string，int，bool，没有类似的构造函数，会有问题；
- `return *new(T), errors.New("Empty")` 都通用


也可以通过如下方式，直接创建范型类型的对象：

```go
// Zero Value Construction
var vGenInt vector[int]
var vGenStr vector[string]
// Non-Zero Value Construction
vGenInt := vector{10, -1}
vGenStr := vector{"A", "B", string([]byte{0xff})}
```

## 9.3 结构体对象

可以定义一个范型struct，范型内部数据结构。比如如下代码定义了链表节点：

```go
type node[T any] struct {
	data T
	pre  *node[T]
	next *node[T]
}
```
定义链表结构：

```go
type list[T any] struct {
	first *node[T]
	last  *node[T]
}
```

然后可以定义方法，格式 `func (l 范型类型[类型列表]) 方法名称 (函数列表) 返回值` ：
```go
func (l *list[T]) add(data T) *node[T] {
	n := node[T]{
		data: data,
		pre:  l.last,
		next: nil,
	}
	if l.first == nil {
		l.first = &n
		l.last = &n
		return &n
	}
	l.last.next = &n
	l.last = &n
	return &n
}
```

可以按照如下方式使用：

```go
type user struct {
	name string
}

func TestAdd(t *testing.T) {
	var users list[user]
	n1 := users.add(user{name: "weenxin"})
	n2 := users.add(user{name: "stone"})
	t.Log(n1.data.name, n2.data.name)

	var pUsers list[*user]
	n3 := pUsers.add(&user{name: "zhansan"})
	n4 := pUsers.add(&user{name: "lisi"})
	t.Log(n3.data.name, n4.data.name)
}
```
## 9.4 行为约束

范型类型或者函数可以约束类型列表中的类型需要具备某种能力。

如下所示的代码，表示了基于多态定义一个方法。

```go
type User struct {
	name string
}

func (u User) String() string {
	return u.name
}

type Stringer interface {
	String() string
}

func Concrete(u User) { //使用具体类型
	u.String()
}

func Polymorphic(s Stringer) { //使用多态类型
	s.String()
}
```

如下所示的代码定义了`stringfy`只能接收满足`fmt.Stringer`接口的类型。

```go

func stringfy[T fmt.Stringer](slice []T) []T {
	result := make([]T, len(slice))
	for index := range slice {
		result[index] = slice[index]
	}
	return result
}
```

由于`User`类型实现了`fmt.Stringer`的接口，因此可以按照如下方式使用：

```go
func TestStringfy(t *testing.T) {
	user := []User{{name: "weenxin"}, {name: "stone"}, {name: "zhangsan"}, {name: "lisi"}}
	values := stringfy(user)
	t.Log(values)
}
```

## 9.5 类型约束

可以指定类型列表中的类型只能是其中一种。比如如下所示代码，定义了一个 `Addable`类型

```go
type Addable interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | uintptr | float32 | float64 | complex64 | complex128 | string
}
```
可以按照如下方式使用定义的类型。
```go
func Add[T Addable](v1 T, v2 T) T {
	return v1 + v2
}
```

按照如下方式使用定义的函数：
```go
func TestAdd(t *testing.T) {
	t.Run("add interge", func(t *testing.T) {
	    t.Log(Add[float64](1, 2.1))
		t.Log(Add(1, 2))
		//t.Log(Add(1, 2.1))  default type float64 of 2.1 does not match inferred type int for T
	})
}
```

如上所示的代码：
- 可以在范型函数中指定具体类型;
- 或者可以通过编译器推断类型，如果类型不能完全匹配将会编译出错;


也可以自定义类型和接口约束一起定义，如下所示的代码：

```go

type Person struct{
	name string
}
func (p Person) matcher(p2 Person) bool {
	return p.name == p2.name
}

type Other struct{
	name string
}
func (o Other) matcher(o2 Other) bool {
	return o.name == o2.name
}

type Food struct {
	name string
}

type matcher[T any] interface{
	Person | Food
	matcher(v T) bool
}

func match[T matcher[T]] (list []T , find T) int {
	for i, v := range  list {
		if v.matcher(find) {
			return i
		}
	}
	return -1
}
```

可以用如下方式使用：

```go

func TestMatch(t * testing.T) {
    //由于Person即实现了interface，又在列表中
	t.Run("test person", func(t *testing.T) {
		peoples := []Person{{name:"weenxin"},{name:"stone"}}
		index := match(peoples,Person{name: "stone"})
		t.Logf("find : %v ", index)
	})

	// Other 不在matcher列表中，所以编译出错
	//t.Run("test others",func(t *testing.T) {
	//	others := []Other{{name:"weenxin"},{name:"stone"}}
	//	index := match(others,Other{name: "stone"})
	//	t.Logf("find : %v ", index)
	//})

	//Food 没有实现方法，所以编译出错
	//t.Run("test foods",func(t *testing.T) {
	//	foods := []Food{{name:"weenxin"},{name:"stone"}}
	//	index := match(foods,Food{name: "stone"})
	//	t.Logf("find : %v ", index)
	//})
}
```

## 9.6 多类型参数

`golang`的范型并不局限在类型列表中只使用一种类型。如下所示的代码 ：

```go
func Print[L fmt.Stringer , V any ](labels []L, values []V)  {
	if len(labels) != len(values) {
		panic("labels and values should be equal")
	}
	for i, v := range values {
		fmt.Printf("%s = %v" , labels[i], v )
	}
}
```

可以用如下方式调用范型方法：

```go

type User struct {
	name string
}

func (u User) String() string{
	return u.name
}

func TestPrint(t *testing.T) {
	t.Run("label is string ,value is integer", func(t *testing.T) {
		labels := []User{{"id"},{"age"}}
		values := []int{1,35}
		Print(labels,values)
	})
}
```

## 9.7 访问字段

目前`golang`对于字段的访问还没有做到如作者所说的能力，具体可以见[issue](https://github.com/golang/go/issues/48522)

官方推荐的字段访问方式，为使用接口，如下所示：

```go
type UserAble interface{
	GetName() string
	GetEmail() string
	SetId(int64)
}

type User struct {
	Id int64
	Name string
	Email string
}

func (u User) GetName() string {
	return u.Name
}

func (u User) GetEmail() string {
	return  u.Email
}

func (u* User) SetId( id int64) {
	u.Id = id
}


type Customer struct {
	Id int64
	Name string
	Email string
}

func (u Customer) GetName() string {
	return u.Name
}

func (u Customer) GetEmail() string {
	return  u.Email
}

func (u* Customer) SetId( id int64) {
	u.Id = id
}

type Entities interface{
	*User | *Customer
	UserAble
}

func InsertData[T Entities](entity T) (T ,error) {
	fmt.Printf("Insert data: name :%s , email :%s \n" ,entity.GetName(),entity.GetEmail())
	entity.SetId(1000)
	return entity,nil
}
```

可以通过如下方式，使用范型方法：

```go
func TestInsertData(t *testing.T) {
	t.Run("insert user", func(t *testing.T) {
		u := &User{
			Id:    0,
			Name:  "weenxin",
			Email: "weenxin@123.com",
		}
		u , _ = InsertData(u)
		t.Logf("after insert data id : %d",u.Id)
	})
	t.Run("insert customer", func(t *testing.T) {
		c := &Customer{
			Id:    0,
			Name:  "weenxin",
			Email: "weenxin@123.com",
		}
		c , _ = InsertData(c)
		t.Logf("after insert data id : %d",c.Id)
	})
}
```

## 9.8 数组限制

slice是golang中使用相对频繁的数据结构了，有时可能你需要对一个不同类型的slice进行类似的操作。那么你可以使用数组限制来提取数据中元素的类型。

比如，先定义一个处理函数:
```go
type operateFunc[T any] func(item T) T
```

然后定义使用方法：

```go
func Every[T any](items []T ,operator operateFunc[T]) []T {
	for index, item := range items {
		items[index] = operator(item)
	}
	return items
}
```

如下方式调用方法：
```go
func double(value int) int {
	return value * 2
}

func TestEvery(t *testing.T) {
	t.Run("test integer", func(t *testing.T) {
		items := []int{1,2,3,4,5,6}
		//提取了items中元素中的类型
		Every(items,double)
		t.Logf("values : %v", items)
	})
}
```

假如我们有如下代码：

```go
type Numbers []int

func Double(n Numbers) Numbers {
    fn := func(n int) int {
        return 2 * n
    }
    numbers := operate(n, fn)
    fmt.Printf("%T", numbers)
    return numbers
}
```

此时我们定义了一个`Numbers`的类型，可不可以定一个支持所有元素类型的slice呢 ？

```go
//定义接收所有类型的slice
type Slice[T any] interface{
    ~[]T
}

func EverySlice[S Slice[T] , T any](s S, operator operateFunc[T]) S {
	for index, item := range s {
		s[index] = operator(item)
	}
	return s
}
```

如下方式调用：

```go
func double(value int) int {
	return value * 2
}
func TestEvery(t *testing.T) {
	t.Run("test slice function", func(t *testing.T) {
		items := []int{1,2,3,4,5,6}
		EverySlice(items,double)
		t.Logf("values : %v", items)
	})
}
```

## 9.9 Channels

也可以定义个接收任何类型的chan。如下所示的代码：

```go
//定义一个处理方法
type workFunc[Result any] func(ctx context.Context) Result

//开始工作，返回一个chan，用户可以往里面扔任务
func DoWork[Result any](ctx context.Context, work workFunc[Result]) chan Result {
	ch :=  make(chan Result,1)
	go func() {
		ch <- work(ctx)
		fmt.Println("done work")
	}()
	return ch
}
```

可以如下方式使用：

```go
//处理方法
func work(ctx context.Context) int {
	time.Sleep(time.Duration(rand.Intn(200) )* time.Millisecond)
	return 100
}

//开始工作
func TestDoWork(t *testing.T) {
	duration := time.Duration(rand.Intn(150))  * time.Millisecond
	ctx , cancel := context.WithTimeout(context.Background(),duration)

	defer cancel()
	//开启处理协程
	ch := DoWork(ctx,work)

    //等待处理结束，或者timeout
	select{
	case data := <-ch:
		t.Logf("got data %v", data)
	case <-ctx.Done():
		t.Log("timeout ")
	}
}
```


也可以开一个pooling, 可以见如下所示代码：

```go
//定义一个处理函数
type workInput[Input any , Result any ] func (input Input) Result
```

```go
//定义处理池
func PoolWork[Input any, Result any]  (size int , work workInput[Input , Result]) (chan Input, func()) {

	var wg sync.WaitGroup
	wg.Add(size)

	ch := make(chan Input)

	for i := 0 ; i < size ; i++ {
	//开启工作协程
		go func() {
			defer wg.Done()
			for input := range ch {
				result := work(input)
				fmt.Println("pollWork :", result)
			}
		}()
	}

    //返回一个结束函数
	cancel:= func() {
		close(ch)
		wg.Wait()
	}
	return ch,cancel
}
```

如下方式使用：

```go

//定义处理函数，输入是一个int，输出也是一个int
func double(value int) int {
	return value *2
}

//开始处理工作
func TestPoolingWork(t *testing.T) {

    //开启处理池
	ch , cancel := PoolWork(4,double)
	//塞工作
	for index := 0 ; index < 100; index ++ {
		ch<-index
	}
	//没有工作了，等待处理结束
	cancel()
}
```

## 9.10 哈希表

哈希表，是一个很常用的container类型，`key`和`value`都是范型中最常见的用例。设计一个hash表，需要考虑如下因素；
- 基于Key的HashFunction
- 基于Key查询Value
- 插入Key,Value

代码的可重用代码非常多。

```go
//定义一个hash函数
type hashFunction[K comparable] func(key K, buckets int) int
```

```go
//定义KeyValue存储对象
type KeyValuePair [K comparable, V any] struct{
	key K
	value V
}
```

```go
//定义表结构
type Table[K comparable, V any] struct {
	hashFunc hashFunction[K]
	buckets int
	data [][]KeyValuePair[K,V]
}
```
新建函数
```go
func New[K comparable, V any](buckets int, function hashFunction[K]) *Table[K,V]{
	return &Table[K,V]{
		hashFunc: function,
		buckets:  buckets,
		data:     make([][]KeyValuePair[K,V], buckets),
	}
}
```

插入数据

```go
func (t * Table[K,V]) Insert(key K, value V){
	bucket := t.hashFunc(key,t.buckets)
	for index, theKey := range t.data[bucket] {
		if key == theKey.key {
			t.data[bucket][index].value = value
			return
		}
	}
	pair := KeyValuePair[K,V] {
		key: key,
		value: value,
	}
	t.data[bucket] = append(t.data[bucket], pair)
}
```

查询数据

```go
func (t * Table[K,V]) Get(key K) (V, bool) {
	bucket := t.hashFunc(key,t.buckets)
	for index, theKey := range t.data[bucket] {
		if key == theKey.key {
			return t.data[bucket][index].value, true
		}
	}
	var zero  V
	return zero, false
}
```

可以按照如下方式使用：

```go
func TestHash(t *testing.T) {
    //定义hash函数
	function := func(data string, buckets int) int {
		h := fnv.New32()
		h.Write([]byte(data))
		return int(h.Sum32())%buckets
	}
	//定义hash函数
	function2 := func(data int, buckets int) int {
		return data%buckets
	}
	//创建表，由于function中只有key的类型，所以需要显式指定
	t1 := New[string,int](100,function)
	//创建表，由于function中只有key的类型，所以需要显式指定
	t2 := New[int,string](100,function2)

	values := map[string]int{"one":1,"two":2,"three":3,"four":4}
	//插入数据
	for key,value := range values {
		t1.Insert(key,value)
		t2.Insert(value,key)
	}
	//获取数据，一套hashtable代码支持多种类型
	for key,value := range values {
		intValue , exists := t1.Get(key)
		if !exists {
			t.Log("not exists")
		}else{
			t.Logf("getting data: %s : %d",key, intValue)
		}
		stringValue , exists := t2.Get(value)
		if !exists {
			t.Log("not exists")
		}else{
			t.Logf("getting data: %d : %s",value, stringValue)
		}
	}
}
```






















