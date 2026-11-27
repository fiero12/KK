package storage

import (
	"sync"
	"financialExchange/enum"
	"log/slog"
)

// sync.Mutex полностью блокирует доступ и для чтения и для записи. 
// RWMutex позволяет потокам читать данные одновременно, но запись осуществляется только одним потоком


type Store[K comparable, V any] struct {
	mutex sync.RWMutex
	items map[K]V
	logger *slog.Logger
}

func NewStore[K comparable, V any](capacity int) * Store[K, V]{
	return &Store[K, V]{
		items: make(map[K]V, capacity),
	}
}

func (store *Store[K, V]) Put(key K, value V) {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	store.items[key] = value
}

func (store *Store[K, V]) Get(key K) (V, bool) {
	store.mutex.RLock()
	defer store.mutex.RUnlock()

	value, ok := store.items[key]
	return value, ok
}

func (store *Store[K, V]) Delete(key K, value V) bool {
	store.mutex.Lock()
	defer store.mutex.Unlock()
	
	if _, ok := store.items[key]; !ok {
		return false
	}

	delete(store.items, key)
	return true
}

func convertToOrder(value any) enum.Order {
	return value.(enum.Order)// извлекаем ордер из value 
	//так как value типа any внем может лежать что угодно  ив том числе ордер
}

func (store *Store[K, V]) Update(key K, modif func(V) V) (V, bool) {
	store.mutex.Lock()//ставим блокировку
	defer store.mutex.Unlock()
	// гарантируем то что после выхода из функции блокировка снимается

	value, ok := store.items[key]// достаем из мапы значение
	if !ok {//если ключа нету
		var zero V// создаем переменную типа Value
		return zero, false// возвращаем переменную и фалс 
		//дальше выход из функции
	}

	// var UpdateValue enum.Order
	UpdatedValue := modif(value)
	//создаем переменную UpdateValue и передаем в неё значение из modif 
	store.items[key] = UpdatedValue//
	return UpdatedValue, true// возвращаем updatevalue и true
	//дальше выход из функции
}

type IndexedStorage[K comparable, V any] struct {
	//требует уточнения
	mutex sync.RWMutex
	order []K
	byKey map[K]int	//Позволяет за О(1) узнать порядок ключа в срезе
	data map[K]V	//Хранилище значений по ключу
}

func NewIndexedStorage[K comparable, V any](capacity int) *IndexedStorage[K , V]{
	return &IndexedStorage[K, V]{
		order: make([]K, 0, capacity),
		byKey: make(map[K]int, capacity),
		data: make(map[K]V, capacity),
	}
}

func (storage *IndexedStorage[K, V]) Put(key K, newValue V){
	storage.mutex.Lock()
	defer storage.mutex. Unlock()

	if _, ok := storage.byKey[key]; ok{
		storage.data[key] = newValue
		// fmt.Println("Значение изменено с".value."на")
		return
	}
	storage.byKey[key] = len(storage.order)
	storage.order = append(storage.order, key)
	storage.data[key] = newValue
}

func (storage *IndexedStorage[K, V]) Delete(key K) bool {
	storage.mutex.Lock()
	defer storage.mutex.Unlock()

	indexInSlice, ok := storage.byKey[key]
	if !ok{
		return false
	}
	lastIndex := len(storage.order) - 1

	if indexInSlice != lastIndex {
	lastIndex := storage.order[lastIndex]  //Получает последний элемент среза order
	storage.order[indexInSlice] = lastIndex  //Заменяет удаляемый элемент на последний элемент среза order
	storage.byKey[lastIndex] = indexInSlice  //Обновляет индекс последнего элемента в мапе byKey
	}

	storage.order = storage.order[:lastIndex]  //Удаляет последний элемент из среза order
	delete(storage.byKey, key)//Удаляет ключ из мапы byKey
	delete(storage.data, key)  //Удаляет ключ из мапы data

	return true
}

func (indexStore *IndexedStorage[K, V]) Ordered() []V{
	indexStore.mutex.RLock()
	defer indexStore.mutex.RUnlock()

	out := make([]V, len(indexStore.order))
	for index, value := range indexStore.order{
		out[index] = indexStore.data[value]
	}
	return out
}

type Storage[T comparable] struct{
	storage map[T]struct{}
}//хранение универсальных множеств

func NewStorage[T comparable](capacity int) *Storage[T] {
	return &Storage[T]{
		storage: make(map[T] struct{}, capacity),
	}
}

func (s *Storage[T]) Add(item T) {
	s.storage[item] = struct{}{}
}

func (s *Storage[T]) Remove(value T) {
	delete(s.storage, value)
}

func (s *Storage[T]) Has(value T) bool {
	_, ok := s.storage[value]
	return ok
}

func (s *Storage[T]) Length() int {
	return len(s.storage)
}

func (s *Storage[T]) Items() []T {
	outSlice := make([]T, 0, len(s.storage))
	for value := range s.storage {
		outSlice = append(outSlice, value)
	}
	return outSlice
}
