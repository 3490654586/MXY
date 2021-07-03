package main

import (
	"fmt"
	"runtime"
	"time"
)

type Mode struct {
	ID int `json:"id"`
	Name string `json:"name"`
}

func main()  {
	 //r := MXY.Default()
	 //mode := &Mode{}
	 //r.POST("/", func(c *MXY.Context) {
	 //	c.ShouldBind(mode)
	 //	fmt.Println("mod=",mode)
		// c.JSON(http.StatusOK,MXY.H{
		// 	"srta":"ok",
		// })
	 //})
     //r.Run()

     runtime.GOMAXPROCS(2)

     go func() {
     	fmt.Println("1")
	 }()

	go func() {
		fmt.Println("2")
	}()

	go func() {
		fmt.Println("3")
	}()


	time.Sleep(time.Second)


	go func() {
		fmt.Println("4")
	}()

}