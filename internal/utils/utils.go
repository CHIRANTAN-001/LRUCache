package utils

type KeyValue struct {
	Key   string
	Value string
}

type KeyValueList struct {
	Items KeyValue
	Next  *KeyValueList
	Prev  *KeyValueList
}

func (l *KeyValueList) Add(key, value string) {
	
}