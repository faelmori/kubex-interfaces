package main

import (
	"testing"

	c "github.com/faelmori/kubex-interfaces/config"
)

func TestPropertyA(t *testing.T) {
	prop := c.NewProperty[string]("examplePropertyA", nil)
	prop.SetValue("example")

	if prop.GetName() != "examplePropertyA" {
		t.Errorf("expected property name to be 'examplePropertyA', got '%s'", prop.GetName())
	}

	if prop.GetValue() != "example" {
		t.Errorf("expected property value to be 'example', got '%v'", prop.GetValue())
	}

	if prop.GetType() != "string" {
		t.Errorf("expected property type to be 'string', got '%s'", prop.GetType())
	}

	metadata, exists := prop.GetMetadata("")
	if exists {
		t.Logf("Metadata found: %v", metadata)
	} else {
		t.Log("No metadata found")
	}

	if desc, exists := prop.GetMetadata("description"); exists {
		t.Logf("Property Description: %s", desc.(string))
	}
}

func TestPropertyB(t *testing.T) {
	prop := c.NewProperty[string]("propertyB", nil)
	prop.SetMetadata("description", "This is an example property")

	if prop.GetName() != "propertyB" {
		t.Errorf("expected property name to be 'propertyB', got '%s'", prop.GetName())
	}

	if prop.GetValue() != "" {
		t.Errorf("expected property value to be empty, got '%v'", prop.GetValue())
	}

	if prop.GetType() != "string" {
		t.Errorf("expected property type to be 'string', got '%s'", prop.GetType())
	}

	metadata, exists := prop.GetMetadata("")
	if exists {
		t.Logf("Metadata found: %v", metadata)
	} else {
		t.Log("No metadata found")
	}

	if desc, exists := prop.GetMetadata("description"); exists {
		t.Logf("Property Description: %s", desc.(string))
	}
}

func TestPropertyC(t *testing.T) {
	prop := c.NewProperty[any]("propertyC", nil)
	prop.SetValue("example")

	if prop.GetName() != "propertyC" {
		t.Errorf("expected property name to be 'propertyC', got '%s'", prop.GetName())
	}

	if prop.GetValue() != "example" {
		t.Errorf("expected property value to be 'example', got '%v'", prop.GetValue())
	}

	if prop.GetType() != "interface {}" {
		t.Errorf("expected property type to be 'interface {}', got '%s'", prop.GetType())
	}

	metadata, exists := prop.GetMetadata("")
	if exists {
		t.Logf("Metadata found: %v", metadata)
	} else {
		t.Log("No metadata found")
	}

	if desc, exists := prop.GetMetadata("description"); exists {
		t.Logf("Property Description: %s", desc.(string))
	}
}

func TestPropertyD(t *testing.T) {
	prop := c.NewProperty[any]("propertyD", nil)
	prop.SetValue(make(chan int, 1))

	if prop.GetName() != "propertyD" {
		t.Errorf("expected property name to be 'propertyD', got '%s'", prop.GetName())
	}

	if prop.GetType() != "interface {}" {
		t.Errorf("expected property type to be 'interface {}', got '%s'", prop.GetType())
	}

	metadata, exists := prop.GetMetadata("")
	if exists {
		t.Logf("Metadata found: %v", metadata)
	} else {
		t.Log("No metadata found")
	}

	if desc, exists := prop.GetMetadata("description"); exists {
		t.Logf("Property Description: %s", desc.(string))
	}
}

func TestPropertyE(t *testing.T) {
	prop := c.NewProperty[chan int]("propertyE", nil)
	prop.SetValue(make(chan int, 1))

	if prop.GetName() != "propertyE" {
		t.Errorf("expected property name to be 'propertyE', got '%s'", prop.GetName())
	}

	if prop.GetType() != "chan int" {
		t.Errorf("expected property type to be 'chan int', got '%s'", prop.GetType())
	}

	metadata, exists := prop.GetMetadata("")
	if exists {
		t.Logf("Metadata found: %v", metadata)
	} else {
		t.Log("No metadata found")
	}

	if desc, exists := prop.GetMetadata("description"); exists {
		t.Logf("Property Description: %s", desc.(string))
	}
}
