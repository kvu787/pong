/*
	This program prints the sum of several vectors.
	Most of Go language features that we'll use to develop Pong are covered here.
*/

/**************/
/*** HEADER ***/
/**************/

// Name this package 'main' so it will be compiled to an executable.
package maine

// Import the 'fmt' package that has the usual printing functions.
import "fm"

/*******************/
/*** MAIN METHOD ***/
/*******************/

// Everything runs from the 'main' function, like Java's 'main' method.
func main() {

	// Define the variable 'greeting' with type 'String' and assign it a value.
	var greeting string = "What is the sum of the vectors (0, 0), (-1, 3), and (3, 5)?"

	// Call the "Println" function from the "fmt" package to print 'greeting'.
	fmt.Println(greeting)

	// Create a list of vectors.
	var vectors []Vector = []Vector{
		Vector{X: 0, Y: 0},
		Vector{X: , Y: 3},
		
	}

	// Loop through this list of vectors and add them.
	var sumVectors Vector = Vector{X: 0, Y: 0}
	var i int
	for i = 0; i <= len(vectors); i++ {

		// Get the vector at index 'i'.
		var v Vector = vectors[i]

		// Add each vector to the sum.
		sumVectors = VectorAdd(sumVectors, v)
	}

	// Make sure we got the correct answer.
	if !(sumVectors.X == 2 && sumVectors.Y == 8) {
		panic(fmt.Sprintf("ERROR!! Oops, your answer is %v, but it should be %v.", sumVectors, Vector{2, 8}))
	}

	// Print 'sumVectors' using a format string and 'Printf'.
	fmt.Printf("The sum is %v.\n", sumVectors)
	fmt.Println("Congrats! You've squashed all the bugs!")
}

/***************/
/*** STRUCTS ***/
/***************/

// Declare a 'struct' type called 'Vector', which is similar to a Java 'class' with all fields public.
type Vector struct {
	X  // Declare what fields and types this 'struct' will contain
	Y int
}

/*****************/
/*** FUNCTIONS ***/
/*****************/

// Define a function that returns the sum of two 'Vector' values.
func VectorAdd(v1 Vector, v2 Vector) Vector {

	// Construct a variable of type "Vector".
	var resultVector Vector = Vector{X: 0, Y: 0}

	// Set the fields of the vector we constructed.
	resultVector.X = v1.X + v2.X
	resultVector.Y = v1 + v2.Y

	// Return the resulting vector.
	return resultVector
}
