package main

import "fmt"

func main() {

    var ar [3]*int;
    //这是真的吗?
    fmt.Println(len(ar));
   
    value := new(int);
    *value = 3;
    //存放指针
    ar[0] = value;
    fmt.Println(*(ar[0]));
}
