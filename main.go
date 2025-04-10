package main

import (
	"fmt"
	c "github.com/faelmori/kubex-interfaces/config"
)

func testPropertyA() {
	prop := c.NewProperty[string]("examplePropertyA", nil)
	prop.SetValue("example")

	// Print the property name and value
	println("Property Name:", prop.GetName())
	println("Property Value:", prop.GetValue())
	println("Property Type:", prop.GetType())

	// Print the metadata
	var metadata interface{}
	var exists bool
	if metadata, exists = prop.GetMetadata(""); exists {
		println("Metadata found")
	} else {
		println("No metadata found")
	}
	println(fmt.Sprintf("Metadata: %v", metadata))

	// Print metadata
	if desc, exists := prop.GetMetadata("description"); exists {
		println("Property Description:", desc.(string))
	}
}

func testPropertyB() {
	prop := c.NewProperty[string]("propertyB", nil)
	prop.SetMetadata("description", "This is an example property")
	// Print the property name and value
	println("Property Name:", prop.GetName())
	println("Property Value:", prop.GetValue())
	println("Property Type:", prop.GetType())

	// Print the metadata
	var metadata interface{}
	var exists bool
	if metadata, exists = prop.GetMetadata(""); exists {
		println("Metadata found")
	} else {
		println("No metadata found")
	}
	println(fmt.Sprintf("Metadata: %v", metadata))

	// Print metadata
	if desc, exists := prop.GetMetadata("description"); exists {
		println("Property Description:", desc.(string))
	}
}

func testPropertyC() {
	prop := c.NewProperty[any]("propertyC", nil)
	prop.SetValue("example")

	// Print the property name and value
	println("Property Name:", prop.GetName())
	println("Property Value:", prop.GetValue())
	println("Property Type:", prop.GetType())

	// Print the metadata
	var metadata interface{}
	var exists bool
	if metadata, exists = prop.GetMetadata(""); exists {
		println("Metadata found")
	} else {
		println("No metadata found")
	}
	println(fmt.Sprintf("Metadata: %v", metadata))

	// Print metadata
	if desc, exists := prop.GetMetadata("description"); exists {
		println("Property Description:", desc.(string))
	}
}

func testPropertyD() {
	prop := c.NewProperty[any]("propertyD", nil)
	prop.SetValue(make(chan int, 1))

	// Print the property name and value
	println("Property Name:", prop.GetName())
	println("Property Value:", prop.GetValue())
	println("Property Type:", prop.GetType())

	// Print the metadata
	var metadata interface{}
	var exists bool
	if metadata, exists = prop.GetMetadata(""); exists {
		println("Metadata found")
	} else {
		println("No metadata found")
	}
	println(fmt.Sprintf("Metadata: %v", metadata))

	// Print metadata
	if desc, exists := prop.GetMetadata("description"); exists {
		println("Property Description:", desc.(string))
	}
}

func testPropertyE() {
	prop := c.NewProperty[chan int]("propertyE", nil)
	prop.SetValue(make(chan int, 1))

	// Print the property name and value
	println("Property Name:", prop.GetName())
	println("Property Value:", prop.GetValue())
	println("Property Type:", prop.GetType())

	// Print the metadata
	var metadata interface{}
	var exists bool
	if metadata, exists = prop.GetMetadata(""); exists {
		println("Metadata found")
	} else {
		println("No metadata found")
	}
	println(fmt.Sprintf("Metadata: %v", metadata))

	// Print metadata
	if desc, exists := prop.GetMetadata("description"); exists {
		println("Property Description:", desc.(string))
	}
}

func main() {
	testPropertyA()
	println("======================================")
	testPropertyB()
	println("======================================")
	testPropertyC()
	println("======================================")
	testPropertyD()
	println("======================================")
	testPropertyE()
}
