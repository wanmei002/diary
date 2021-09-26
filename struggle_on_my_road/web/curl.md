#### go ���� curl ����
- curl ���ڿͻ��˵�����
    ```go
    cli := &http.Client{
      Timeout:6,  // �����������ڵ�ʱ��
      Transport: &http.Transport{
          	ResponseHeaderTimeout: 3, // �ȴ���Ӧ��ʱ��(responseHeader) (��Ӧ������Ӧͷ��Ϣ)
        	DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
                conn, err := net.DialTimeout(network, addr, time.Second*2)  // ���ų�ʱ (�൱�ڲ���IP ���ӷ�����)
                if err != nil {
                    return nil, err
                }
                conn.SetDeadline(time.Now().Add(time.Second * this.dialTimeout)) // ��д��ʱʱ��
                return conn, nil
            },
      }
    }
    ```
- �����������Ϣ url header method httpЭ�� �ȵ�
    ```go
      req, err := http.NewRequest(method, url, body)  // method-GET|POST url-http://www.xxx.com/xx  body-post��������
      req.Header.Set(k, v) // ���õ�����ͷ��Ϣ
      req.AddCookie(&http.Cookie{  // ����cookie
    			Name:  k,
    			Value: v,
    		})
      q := req.URL.Query()
      q.Add(k,v) // ����  get �������
      req.URL.RawQuery = q.Encode()
    
    ```
    
- ���������
    ```go
    resp, err := cli.DO(req)  // cli-��һ�����cli  req-�ڶ������req
    headerInfo := resp.Header  // ��ȡͷ��Ϣ
    body, err := ioutil.ReadAll(resp.Body)  // ��ȡ���ص�������Ϣ
    
    ```
    