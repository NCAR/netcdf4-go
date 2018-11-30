package netcdf4

import "fmt"

//type Set map[interface{}]bool
//
//func NewSet() Set {
//	return make(map[interface{}]bool)
//}
//
//func (s Set) Has(v interface{}) bool {
//	_, ok := s[v]
//	return ok
//}
//
//func (s Set) Add(v interface{}) {
//	s[v] = true
//}
//
//func (s Set) Erase(v interface{}) {
//	delete(s, v)
//}
//
//type Multimap map[interface{}]map[interface{}]bool
//
//func NewMultimap() Multimap {
//	return make(map[interface{}]map[interface{}]bool)
//}
//
//func (m Multimap) Add(key, value interface{}) {
//	if !m.HasKey(key) {
//		m[key] = make(map[interface{}]bool)
//	}
//	m[key][value] = true
//}
//
//func (m Multimap) HasKey(key interface{}) bool {
//	_, ok := m[key]
//	return ok
//}
//
//func (m Multimap) Has(key, value interface{}) bool {
//	if !m.HasKey(key) {
//		return false
//	}
//
//	_, ok := m[key][value]
//	return ok
//}
//
//func (m Multimap) EqualRange(key interface{}) []interface{} {
//	if !m.HasKey(key) {
//		return nil
//	}
//
//	ans := make([]interface{}, len(m[key]))
//	i := 0
//	for v := range m[key] {
//		ans[i] = v
//		i++
//	}
//
//	return ans
//}
//
//func (m Multimap) EraseKey(key interface{}) error {
//	if !m.HasKey(key) {
//		return fmt.Errorf("key %v not present in Multimap", key)
//	}
//
//	delete(m, key)
//	return nil
//}
//
//func (m Multimap) Erase(key, value interface{}) error {
//	if !m.Has(key, value) {
//		return fmt.Errorf("key value pair [%v, %v] not present in Multimap", key, value)
//	}
//
//	delete(m[key], value)
//	if len(m[key]) == 0 {
//		delete(m, key)
//	}
//
//	return nil
//}
//
//func (m Multimap) GetAllPair() ([]interface{}, []interface{}) {
//	keys := make([]interface{}, 0)
//	fields := make([]interface{}, 0)
//	for key := range m {
//		for v := range m[key] {
//			keys = append(keys, key)
//			fields = append(fields, v)
//		}
//	}
//	return keys, fields
//}
//
//func (m Multimap) Length() int {
//	ans := 0
//
//	for _, vm := range m {
//		ans += len(vm)
//	}
//
//	return ans
//}
//
//func (m Multimap) Size() int {
//	return len(m)
//}

///////////////////////////////////////////////////

// type SetG map[*Group]bool

// func NewSetG() SetG {
// 	return make(map[*Group]bool)
// }

// func (s SetG) Has(v *Group) bool {
// 	_, ok := s[v]
// 	return ok
// }

// func (s SetG) Add(v *Group) {
// 	s[v] = true
// }

// func (s SetG) Erase(v *Group) {
// 	delete(s, v)
// }

///////////////////////////////////////////////////
type SetD map[Dim]bool

func NewSetD() SetD {
	return make(map[Dim]bool)
}

func (s SetD) Has(v Dim) bool {
	_, ok := s[v]
	return ok
}

func (s SetD) Add(v Dim) {
	s[v] = true
}

func (s SetD) Erase(v Dim) {
	delete(s, v)
}

///////////////////////////////////////////////////
type SetV map[Var]bool

func NewSetV() SetV {
	return make(map[Var]bool)
}

func (s SetV) Has(v Var) bool {
	_, ok := s[v]
	return ok
}

func (s SetV) Add(v Var) {
	s[v] = true
}

func (s SetV) Erase(v Var) {
	delete(s, v)
}

///////////////////////////////////////////////////
// something is wrong about the general type of Multimap with interface{}, no reason is found, thus, we use this one
type MultimapG map[string]map[*Group]bool

func NewMultimapG() MultimapG {
	return make(map[string]map[*Group]bool)
}

func (m MultimapG) Add(key string, value *Group) {
	if !m.HasKey(key) {
		m[key] = make(map[*Group]bool)
	}
	m[key][value] = true
}

func (m MultimapG) HasKey(key string) bool {
	_, ok := m[key]
	return ok
}

func (m MultimapG) Has(key string, value *Group) bool {
	if !m.HasKey(key) {
		return false
	}

	_, ok := m[key][value]
	return ok
}

func (m MultimapG) EqualRange(key string) []*Group {
	if !m.HasKey(key) {
		return nil
	}

	ans := make([]*Group, len(m[key]))
	i := 0
	for v := range m[key] {
		ans[i] = v
		i++
	}

	return ans
}

func (m MultimapG) EraseKey(key string) error {
	if !m.HasKey(key) {
		return fmt.Errorf("key %v not present in MultimapG", key)
	}

	delete(m, key)
	return nil
}

func (m MultimapG) Erase(key string, value *Group) error {
	if !m.Has(key, value) {
		return fmt.Errorf("key value pair [%v, %v] not present in MultimapG", key, value)
	}

	delete(m[key], value)
	if len(m[key]) == 0 {
		delete(m, key)
	}

	return nil
}

func (m MultimapG) GetAllPair() ([]string, []*Group) {
	keys := make([]string, 0)
	fields := make([]Group, 0)
	for key := range m {
		for v := range m[key] {
			keys = append(keys, key)
			fields = append(fields, v)
		}
	}
	return keys, fields
}

func (m MultimapG) Length() int {
	ans := 0

	for _, vm := range m {
		ans += len(vm)
	}

	return ans
}

func (m MultimapG) Size() int {
	return len(m)
}

///////////////////////////////////////////////////
// something is wrong about the general type of Multimap with interface{}, no reason is found, thus, we use this one
type MultimapD map[string]map[Dim]bool

func NewMultimapD() MultimapD {
	return make(map[string]map[Dim]bool)
}

func (m MultimapD) Add(key string, value Dim) {
	if !m.HasKey(key) {
		m[key] = make(map[Dim]bool)
	}
	m[key][value] = true
}

func (m MultimapD) HasKey(key string) bool {
	_, ok := m[key]
	return ok
}

func (m MultimapD) Has(key string, value Dim) bool {
	if !m.HasKey(key) {
		return false
	}

	_, ok := m[key][value]
	return ok
}

func (m MultimapD) EqualRange(key string) []Dim {
	if !m.HasKey(key) {
		return nil
	}

	ans := make([]Dim, len(m[key]))
	i := 0
	for v := range m[key] {
		ans[i] = v
		i++
	}

	return ans
}

func (m MultimapD) EraseKey(key string) error {
	if !m.HasKey(key) {
		return fmt.Errorf("key %v not present in MultimapD", key)
	}

	delete(m, key)
	return nil
}

func (m MultimapD) Erase(key string, value Dim) error {
	if !m.Has(key, value) {
		return fmt.Errorf("key value pair [%v, %v] not present in MultimapD", key, value)
	}

	delete(m[key], value)
	if len(m[key]) == 0 {
		delete(m, key)
	}

	return nil
}

func (m MultimapD) GetAllPair() ([]string, []Dim) {
	keys := make([]string, 0)
	fields := make([]Dim, 0)
	for key := range m {
		for v := range m[key] {
			keys = append(keys, key)
			fields = append(fields, v)
		}
	}
	return keys, fields
}

func (m MultimapD) Length() int {
	ans := 0

	for _, vm := range m {
		ans += len(vm)
	}

	return ans
}

func (m MultimapD) Size() int {
	return len(m)
}

///////////////////////////////////////////////////
// something is wrong about the general type of Multimap with interface{}, no reason is found, thus, we use this one
type MultimapV map[string]map[Var]bool

func NewMultimapV() MultimapV {
	return make(map[string]map[Var]bool)
}

func (m MultimapV) Add(key string, value Var) {
	if !m.HasKey(key) {
		m[key] = make(map[Var]bool)
	}
	m[key][value] = true
}

func (m MultimapV) HasKey(key string) bool {
	_, ok := m[key]
	return ok
}

func (m MultimapV) Has(key string, value Var) bool {
	if !m.HasKey(key) {
		return false
	}

	_, ok := m[key][value]
	return ok
}

func (m MultimapV) EqualRange(key string) []Var {
	if !m.HasKey(key) {
		return nil
	}

	ans := make([]Var, len(m[key]))
	i := 0
	for v := range m[key] {
		ans[i] = v
		i++
	}

	return ans
}

func (m MultimapV) EraseKey(key string) error {
	if !m.HasKey(key) {
		return fmt.Errorf("key %v not present in MultimapV", key)
	}

	delete(m, key)
	return nil
}

func (m MultimapV) Erase(key string, value Var) error {
	if !m.Has(key, value) {
		return fmt.Errorf("key value pair [%v, %v] not present in MultimapV", key, value)
	}

	delete(m[key], value)
	if len(m[key]) == 0 {
		delete(m, key)
	}

	return nil
}

func (m MultimapV) GetAllPair() ([]string, []Var) {
	keys := make([]string, 0)
	fields := make([]Var, 0)
	for key := range m {
		for v := range m[key] {
			keys = append(keys, key)
			fields = append(fields, v)
		}
	}
	return keys, fields
}

func (m MultimapV) Length() int {
	ans := 0

	for _, vm := range m {
		ans += len(vm)
	}

	return ans
}

func (m MultimapV) Size() int {
	return len(m)
}
