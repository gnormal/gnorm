package kace_test

import (
	"fmt"

	"github.com/codemodus/kace"
)

func Example() {
	s := "this is a test sql."

	fmt.Println(kace.Camel(s))
	fmt.Println(kace.Pascal(s))

	fmt.Println(kace.Snake(s))
	fmt.Println(kace.SnakeUpper(s))

	fmt.Println(kace.Kebab(s))
	fmt.Println(kace.KebabUpper(s))

	customInitialisms := map[string]bool{
		"THIS": true,
	}
	k, err := kace.New(customInitialisms)
	if err != nil {
		// handle error
	}

	fmt.Println(k.Camel(s))
	fmt.Println(k.Pascal(s))

	fmt.Println(k.Snake(s))
	fmt.Println(k.SnakeUpper(s))

	fmt.Println(k.Kebab(s))
	fmt.Println(k.KebabUpper(s))

	// Output:
	// thisIsATestSQL
	// ThisIsATestSQL
	// this_is_a_test_sql
	// THIS_IS_A_TEST_SQL
	// this-is-a-test-sql
	// THIS-IS-A-TEST-SQL
	// thisIsATestSql
	// THISIsATestSql
	// this_is_a_test_sql
	// THIS_IS_A_TEST_SQL
	// this-is-a-test-sql
	// THIS-IS-A-TEST-SQL
}
