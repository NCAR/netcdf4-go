package netcdf4

/* This part is mainly copied from suyash with MIT License
   MIT License

   Copyright (c) 2016 Suyash

   Permission is hereby granted, free of charge, to any person obtaining a copy
   of this software and associated documentation files (the "Software"), to deal
   in the Software without restriction, including without limitation the rights
   to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
   copies of the Software, and to permit persons to whom the Software is
   furnished to do so, subject to the following conditions:

   The above copyright notice and this permission notice shall be included in all
   copies or substantial portions of the Software.

   THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
   IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
   FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
   AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
   LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
   OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
   SOFTWARE.*/

import "fmt"

type Set map[interface{}]bool

func NewSet() Set {
	return make(map[interface{}]bool)
}

func (s Set) Has(v interface{}) bool {
	_, ok := s[v]
	return ok
}

func (s Set) Add(v interface{}) {
	s[v] = true
}

func (s Set) Erase(v interface{}) {
	delete(s, v)
}

type Multimap map[interface{}]map[interface{}]bool

func NewMultimap() Multimap {
	return make(map[interface{}]map[interface{}]bool)
}

func (m Multimap) Add(key, value interface{}) {
	if !m.HasKey(key) {
		m[key] = make(map[interface{}]bool)
	}
	m[key][value] = true
}

func (m Multimap) HasKey(key interface{}) bool {
	_, ok := m[key]
	return ok
}

func (m Multimap) Has(key, value interface{}) bool {
	if !m.HasKey(key) {
		return false
	}

	_, ok := m[key][value]
	return ok
}

func (m Multimap) EqualRange(key interface{}) []interface{} {
	if !m.HasKey(key) {
		return nil
	}

	ans := make([]interface{}, len(m[key]))
	i := 0
	for v := range m[key] {
		ans[i] = v
		i++
	}

	return ans
}

func (m Multimap) EraseKey(key interface{}) error {
	if !m.HasKey(key) {
		return fmt.Errorf("key %v not present in Multimap", key)
	}

	delete(m, key)
	return nil
}

func (m Multimap) Erase(key, value interface{}) error {
	if !m.Has(key, value) {
		return fmt.Errorf("key value pair [%v, %v] not present in Multimap", key, value)
	}

	delete(m[key], value)
	if len(m[key]) == 0 {
		delete(m, key)
	}

	return nil
}

func (m Multimap) GetAllPair() ([]interface{}, []interface{}) {
	keys := make([]interface{}, 0)
	fields := make([]interface{}, 0)
	for key := range m {
		for v := range m[key] {
			keys = append(keys, key)
			fields = append(fields, v)
		}
	}
	return keys, fields
}

func (m Multimap) Length() int {
	ans := 0

	for _, vm := range m {
		ans += len(vm)
	}

	return ans
}

func (m Multimap) Size() int {
	return len(m)
}



type MultimapG map[string]map[Group]bool

func NewMultimapG() MultimapG {
	return make(map[interface{}]map[interface{}]bool)
}

func (m Multimap) Add(key, value interface{}) {
	if !m.HasKey(key) {
		m[key] = make(map[interface{}]bool)
	}
	m[key][value] = true
}

func (m Multimap) HasKey(key interface{}) bool {
	_, ok := m[key]
	return ok
}

func (m Multimap) Has(key, value interface{}) bool {
	if !m.HasKey(key) {
		return false
	}

	_, ok := m[key][value]
	return ok
}

func (m Multimap) EqualRange(key interface{}) []interface{} {
	if !m.HasKey(key) {
		return nil
	}

	ans := make([]interface{}, len(m[key]))
	i := 0
	for v := range m[key] {
		ans[i] = v
		i++
	}

	return ans
}

func (m Multimap) EraseKey(key interface{}) error {
	if !m.HasKey(key) {
		return fmt.Errorf("key %v not present in Multimap", key)
	}

	delete(m, key)
	return nil
}

func (m Multimap) Erase(key, value interface{}) error {
	if !m.Has(key, value) {
		return fmt.Errorf("key value pair [%v, %v] not present in Multimap", key, value)
	}

	delete(m[key], value)
	if len(m[key]) == 0 {
		delete(m, key)
	}

	return nil
}

func (m Multimap) GetAllPair() ([]interface{}, []interface{}) {
	keys := make([]interface{}, 0)
	fields := make([]interface{}, 0)
	for key := range m {
		for v := range m[key] {
			keys = append(keys, key)
			fields = append(fields, v)
		}
	}
	return keys, fields
}

func (m Multimap) Length() int {
	ans := 0

	for _, vm := range m {
		ans += len(vm)
	}

	return ans
}

func (m Multimap) Size() int {
	return len(m)
}
