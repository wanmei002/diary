##### Go�����Ŀ�
 - ��ռCPU�������� Goroutine ����
 ```go
 func main() {
 	runtime.GOMAXPROCS(1)
 	
 	go func(){
 		for i:=0; i < 10; i++ {
 			fmt.Println(i)
 		}
 	}()
 	// for{} ռ��CPU ��������Goroutine �͵ò���ִ��
 	// Ҫ������������ ��� for { runtime.Gosched() }
 	// ���� for ���� select
 	for {} 
 }
 ``` 
 
 - ��ѭ���ڲ�ִ�� defer ���
 ```go
 func main(){
 	for i:=0; i<5; i++ {
 		f,err := os.Open("/path/to/file")
 		if err != nil {
 			fmt.Println(err)
 		}
 		// defer �ں����˳�ʱ��ִ�� �� for ִ�� defer �ᵼ����Դ�ӳ��ͷ�
 		defer f.Close()
 	}
 }
 // ����ķ��������� for �й���һ���ֲ�����, �ھֲ������ڲ�ִ�� defer :
 func main(){
 	for i:=0; i<5; i++ {
 		func(){
 			f, err := os.Open("/path/to/file")
 			if err != nil {
 				fmt.Println(err)
 			}
 			defer f.Close()
 		}()
 	}
 }
 ```
 
 - ��Ƭ�ᵼ�������ײ����鱻�������ײ������޷��ͷ��ڴ棬����ײ�����ϴ����ڴ�����ܴ��ѹ��
 
 - ��ֹ main �����˳��ķ���
 ```go
 func main() {
 	defer func(){ for{} }()
 	// �����������
 	defer func(){ select {} }()
 	// �����������
 	defer func(){ <-make(chan bool) }()
 }
 ```
 
 
 
 
 
 