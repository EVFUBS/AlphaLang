println("test")

var test = "Hello World"
println(test)

var IncrementTest = 1
IncrementTest += 2
println(IncrementTest)

var a = 15
println(a)

print(10 + 15)

func testFunc(a, b) {
    println(a + b)
}

for var x = 0; x < 10; x+=1 {
    println("hello")
}

testFunc(90,10)

var hello = "Hello"
var world = "World"

var array = [1,2,3,4]
println(array)
var last = pop(array)
println(array)

println(hello + " " + world)

func returnTest(num) {
    return num + 1
}

var num1 = rand(1,100)
println(num1)

var retest = returnTest(10)
println(retest)

var indexTest = [1,2,3,4]
println(indexTest[2])

var mapTest = {"test1": 1, "test2": 2}

println(mapTest["test1"])
append(mapTest, "test3", 3)
println(mapTest)

var userInput = input("Enter your name")
print("this is your name ", userInput)

var appendTest = [1,2,3]
append(appendTest, 4)
println(appendTest)



var num2 = rand(101, 200)
println(num2)