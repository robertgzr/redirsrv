package kiwi

type StringValue string

func (v StringValue) MarshalBinary() ([]byte, error) {
	return []byte(v), nil
}
func (v *StringValue) UnmarshalBinary(data []byte) error {
	*v = StringValue(string(data))
	return nil
}

type ByteValue []byte

func (v ByteValue) MarshalBinary() ([]byte, error) {
	return []byte(v), nil
}
func (v *ByteValue) UnmarshalBinary(data []byte) error {
	*v = ByteValue(data)
	return nil
}
