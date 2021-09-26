#### go 发起 curl 请求
- curl 属于客户端的请求
    ```go
    cli := &http.Client{
      Timeout:6,  // 整个请求周期的时间
      Transport: &http.Transport{
          	ResponseHeaderTimeout: 3, // 等待响应的时间(responseHeader) (响应返回响应头信息)
        	DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
                conn, err := net.DialTimeout(network, addr, time.Second*2)  // 拨号超时 (相当于查找IP 连接服务器)
                if err != nil {
                    return nil, err
                }
                conn.SetDeadline(time.Now().Add(time.Second * this.dialTimeout)) // 读写超时时间
                return conn, nil
            },
      }
    }
    ```
- 设置请求的信息 url header method http协议 等等
    ```go
      req, err := http.NewRequest(method, url, body)  // method-GET|POST url-http://www.xxx.com/xx  body-post的请求体
      req.Header.Set(k, v) // 设置的请求头信息
      req.AddCookie(&http.Cookie{  // 设置cookie
    			Name:  k,
    			Value: v,
    		})
      q := req.URL.Query()
      q.Add(k,v) // 设置  get 请求参数
      req.URL.RawQuery = q.Encode()
    
    ```
    
- 最后发送请求
    ```go
    resp, err := cli.DO(req)  // cli-第一步里的cli  req-第二部里的req
    headerInfo := resp.Header  // 获取头信息
    body, err := ioutil.ReadAll(resp.Body)  // 获取返回的主体信息
    
    ```
    