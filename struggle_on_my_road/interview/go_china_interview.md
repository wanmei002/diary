### ��golang�������������ܽ�
1. defer panic 
    ```go
    package main
     import (
         "fmt"
    )
     func main() {
         defer_call()
     }
    
    func defer_call() {
        defer func() { fmt.Println("��ӡǰ") }()
        defer func() { fmt.Println("��ӡ��") }()
        defer func() { fmt.Println("��ӡ��") }()
    
        panic("�����쳣")
     
        defer func(){ fmt.Println("�쳣��") }()
    }
    ```
    > ����: ����Ľ����  ��ӡ�� - ��ӡ�� - ��ӡǰ - �����쳣�� �ڴ����쳣֮ǰ�Ѿ�����ջ�ڵ� defer �ᴥ��, ���� �쳣(panic) ���defer �ǲ������� ,
    ����û���ü�ѹ��ջ��
    
2. for range point ����
    ```go
    slic := []int{0,1,2,3}
    m := make(map[int]*int)
    for key, val := range slic {
        m[key] = &val
    }
    
    for key,val := range m {
        fmt.Println(key, "->", *val)
    }
    // �����:
    0->3
    1->3
    2->3
    3->3
    ```
    > ����: for range ������ key, val �Ǿֲ�����, range ��ʱ���� slice map ��ļ�ֵ ��ֵ�� key, val, 
    ���� key val �ĵ�ַ�ǲ���� ���ֻ�ǵ�ַ��洢��ֵ