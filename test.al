func main () {
    var run = true
    var num = rand(1,100)
    println(num)
    while run == true {
        var test = input("enter something")
        var intTest = int(test)
        if intTest == num {
            println("correct")
            run = false
        } elif intTest < num {
            println("higher")
        } else {
            println("lower")
        }
    }    

    func playAgain() {
        var playAgain = input("would you like to play again? Y/N")
        if playAgain == "Y" {
            main()
        } elif playAgain == "N" {
            println("thanks for playing")
        } else {
            println("enter valid input")
            playAgain()
        }
    }
    playAgain()
}

main()