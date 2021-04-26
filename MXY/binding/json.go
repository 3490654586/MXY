package binding

import (
	"encoding/json"
	"fmt"
	"net/http"
)


type jsonBinding struct {}

func (jsonBinding)Name() string  {
return "jsom"
}

func (jsonBinding) Bind(req *http.Request, obj interface{}) error{
	if req == nil  ||req.Body ==nil{
		return fmt.Errorf("请求主体空")
	}
	decoder :=json.NewDecoder(req.Body)
    err :=	decoder.Decode(obj)
	if err != nil {
		fmt.Println(err)
	}
   return nil
}

