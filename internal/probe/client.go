package probe

import (
	"time"
	"context"
)

type Result struct {
	URL        string
	StatusCode int
	Latency    time.Duration
	Err        error
	CheckedAt  time.Time
}

type Client struct{
   httpClient *http.Client
}

func NewClient(timeout time.Duration)*Client{
     return &Client{httpClient:&http.Client{Timeout:timeout}}
}

func (c *Client)Check(ctx context.Context,url string)Result{
	start:=time.Now()
	req,err:=http.NewRequestWithContext(ctx,http.MethodGet,url,nil)
    if err!=nil{
       return Result{URL:url,Latency:time.Since(start),Err :err,CheckedAt:start}
    }
    resp,err:=c.httpClient.Do(req)
    latency:=time.Since(start)
    if err!=nil{
    	return Result{URL:url,Latency:latency,Err :err,CheckedAt:start}
    }
    defer resp.Body.Close()
    if (resp.StatusCode>=300 || resp.StatusCode<200) {
         return Result{URL:url,StatusCode:resp.StatusCode,Latency:latency,Err:fmt.Errorf("unexpected status code %d",resp.StatusCode),CheckedAt:start
         }
    }
    return Result{URL:url,StatusCode:resp.StatusCode,Latency:latency,Err:nil,CheckedAt:start}
}
