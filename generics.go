// Go 2
type Map[K, V interface{}] struct { /* ... */ }

// Go 1
type Map struct { /* ... */ }

// Go 2
func NewMap[K, V interface{}]() *Map[K, V]

// Go 1
func NewMap() *Map

// Go 2
func (m *Map) Get[K, V interface{}](key K) (val V)

// Go 1
func (m *Map) Get(key interface{}) (val interface{})

// Go 2
m := go2pkg.NewMap[string, int]()
m.Set("hello", 2)
x := m.Get("hello")

// Go 1
m := go2pkg.NewMap()
m.Set("hello", 2)
x := x.Get("hello").(int)

// Go 2
func Dictionary() *Map[string, string]

// Go 1
func Dictionary() *Map // but it's more than *Map[interface{}, interface{}], see below

// Go 2
dict := go2pkg.Dictionary()
dict.Set("yay", 2) // compile error

// Go 1
dict := go2pkg.Dictionary() // type of dict is *go2pkg.Map
dict.Set("yay", 2)          // runtime error
