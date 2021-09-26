### json to struct 
1. 引入包 `encoding/json`
    ```
    // 序列化
    func Marshal(v interface{}) ([]byte, error) {}
    // 反序列化
    func Unmarshal(data []byte, v interface{}) error {}
    
    // 
    func NewDecoder(r io.Reader) *Decoder {}
    Decoder.Decode(&map)
    // NewEncoder 需要传入一个 io
    func NewEncoder(w io.Writer) *Encoder {}
    Encoder.Encode(struct)
    
    ```