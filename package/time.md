### time ��
 - ʱ���ʽ `2016-01-02 15:04:05`
 
#### ��ȡ time ���� Time ����
 - ��ȡ��ǰʱ�� `time.Now()`
 
 - ��ȡָ��ʱ�� `func ParseInLocation(layout, value string, loc *Location)(Time, error)`
    ```go
    time.ParseInLocation("2006-01-02 15:04:05", "2020-06-06 10:59:59", time.Local)
    ```
    
 - ��ȡʱ���
    ```go
    fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
    ```
    
 - ʱ���
    ```go
    dt, _ := time.ParseDuration("1m50s") // �����s ����� 60 h �����24 ��������
    now := time.Now()
    newTime := now.Add(dt)
    fmt.Println(newTime) // 1��50����ʱ��
    ```
    
 - ʱ���
    ```go
    now := time.Now()
    now.Sub(Time) // ��������һ��ʱ����� ���� Duration ����
    now.Add(Duration) // ʱ����С��λ
    ```
    
 - �Ƚ�����ʱ���
    ```go
    now := time.Now()
    now.After(Time) // ������һ�� Time �ṹ���ʵ�� Before Ҳ��һ����
    now.Before(Time)
    ```
    
#### ��ʱ��صķ���
 - `time.After(d Duration)` ��ʾ����ʱ��֮��, ������ȡ�� channel ����֮ǰ������, ����������Լ���ִ��
    `After`ͨ���������������ʱ���� 
    ```go
    select {
    case m := <-c:
 	    fmt.Println("hello world")
    case <- time.After(5 * time.Minute) :
 	    fmt.Println("timed out")
    }
    ```
 - `time.Sleep(d Duration)`  ��������ָ����ʱ��, Ȼ��ż���ִ��
 
 - `time.Tick(d Duration) <-chan Time` �÷��� `time.After` ���, �������Ǳ�ʾÿ������ʱ��֮��, ��һ���ظ��Ĺ���(���Ե�����ʹ��)�������� `After` һ��
 
 - `ticker := time.NewTicker(1 * time.Second)` �� Tick �÷�һ��, ���ǿ��Ե��� `ticker.Stop()` ��ֹͣ��ʱ  
 
 