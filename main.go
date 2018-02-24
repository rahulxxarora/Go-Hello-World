package main

func main() {
    a := App{} 
    a.Initialize("root", "root", "go_todo")
    a.Run(":8080")
}