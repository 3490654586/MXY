package main

import (
	"MXY-WEB/MXY"
	"fmt"
	"net/http"
)

type Mode struct {
	ID int `json:"id"`
	Name string `json:"name"`
}

func main()  {
	 r := MXY.Default()
	 mode := &Mode{}
	 r.POST("/", func(c *MXY.Context) {
	 	c.ShouldBind(mode)
	 	fmt.Println("mod=",mode)
		 c.JSON(http.StatusOK,MXY.H{
		 	"srta":"ok",
		 })
	 })
     r.Run()
}