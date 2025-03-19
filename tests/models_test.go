package tests

import (
	"strings"
	"testing"

	genanki "github.com/npcnixel/genanki-go"
)

func TestNewModel(t *testing.T) {
	modelID := int64(1234567890)
	modelName := "Test Model"

	model := genanki.NewModel(modelID, modelName)

	if model.ID != modelID {
		t.Errorf("Expected model ID %d, got %d", modelID, model.ID)
	}

	if model.Name != modelName {
		t.Errorf("Expected model name %s, got %s", modelName, model.Name)
	}

	if len(model.Fields) != 0 {
		t.Errorf("Expected 0 fields, got %d", len(model.Fields))
	}

	if len(model.Templates) != 0 {
		t.Errorf("Expected 0 templates, got %d", len(model.Templates))
	}
}

func TestNewBasicModel(t *testing.T) {
	modelID := int64(1234567890)
	modelName := "Basic Model"

	basicModel := genanki.NewBasicModel(modelID, modelName)

	if basicModel.ID != modelID {
		t.Errorf("Expected model ID %d, got %d", modelID, basicModel.ID)
	}

	if basicModel.Name != modelName {
		t.Errorf("Expected model name %s, got %s", modelName, basicModel.Name)
	}

	// Basic model should have Front and Back fields
	if len(basicModel.Model.Fields) != 2 {
		t.Errorf("Expected 2 fields, got %d", len(basicModel.Model.Fields))
	} else {
		if basicModel.Model.Fields[0].Name != "Front" {
			t.Errorf("Expected first field name to be 'Front', got '%s'", basicModel.Model.Fields[0].Name)
		}
		if basicModel.Model.Fields[1].Name != "Back" {
			t.Errorf("Expected second field name to be 'Back', got '%s'", basicModel.Model.Fields[1].Name)
		}
	}

	// Basic model should have one template
	if len(basicModel.Model.Templates) != 1 {
		t.Errorf("Expected 1 template, got %d", len(basicModel.Model.Templates))
	} else {
		if basicModel.Model.Templates[0].Name != "Card 1" {
			t.Errorf("Expected template name to be 'Card 1', got '%s'", basicModel.Model.Templates[0].Name)
		}
	}
}

func TestNewClozeModel(t *testing.T) {
	modelID := int64(2345678901)
	modelName := "Cloze Model"

	clozeModel := genanki.NewClozeModel(modelID, modelName)

	if clozeModel.ID != modelID {
		t.Errorf("Expected model ID %d, got %d", modelID, clozeModel.ID)
	}

	if clozeModel.Name != modelName {
		t.Errorf("Expected model name %s, got %s", modelName, clozeModel.Name)
	}

	// Cloze model should have Text and Extra fields
	if len(clozeModel.Model.Fields) != 2 {
		t.Errorf("Expected 2 fields, got %d", len(clozeModel.Model.Fields))
	} else {
		if clozeModel.Model.Fields[0].Name != "Text" {
			t.Errorf("Expected first field name to be 'Text', got '%s'", clozeModel.Model.Fields[0].Name)
		}
		if clozeModel.Model.Fields[1].Name != "Extra" {
			t.Errorf("Expected second field name to be 'Extra', got '%s'", clozeModel.Model.Fields[1].Name)
		}
	}

	// Cloze model should have one template
	if len(clozeModel.Model.Templates) != 1 {
		t.Errorf("Expected 1 template, got %d", len(clozeModel.Model.Templates))
	} else {
		template := clozeModel.Model.Templates[0]
		if template.Name != "Cloze" {
			t.Errorf("Expected template name to be 'Cloze', got '%s'", template.Name)
		}
		if !strings.Contains(template.Qfmt, "{{cloze:Text}}") {
			t.Errorf("Expected Qfmt to contain '{{cloze:Text}}', got '%s'", template.Qfmt)
		}
	}
}

func TestAddField(t *testing.T) {
	model := genanki.NewModel(1234567890, "Test Model")

	field1 := genanki.Field{Name: "Field1", Font: "Arial", Size: 20}
	model.AddField(field1)

	if len(model.Fields) != 1 {
		t.Errorf("Expected 1 field, got %d", len(model.Fields))
	}

	if model.Fields[0].Name != "Field1" {
		t.Errorf("Expected field name to be 'Field1', got '%s'", model.Fields[0].Name)
	}

	if model.Fields[0].Ord != 0 {
		t.Errorf("Expected field ord to be 0, got %d", model.Fields[0].Ord)
	}

	field2 := genanki.Field{Name: "Field2", Font: "Arial", Size: 20}
	model.AddField(field2)

	if len(model.Fields) != 2 {
		t.Errorf("Expected 2 fields, got %d", len(model.Fields))
	}

	if model.Fields[1].Name != "Field2" {
		t.Errorf("Expected field name to be 'Field2', got '%s'", model.Fields[1].Name)
	}

	if model.Fields[1].Ord != 1 {
		t.Errorf("Expected field ord to be 1, got %d", model.Fields[1].Ord)
	}
}

func TestAddTemplate(t *testing.T) {
	model := genanki.NewModel(1234567890, "Test Model")

	template1 := genanki.Template{
		Name: "Template1",
		Qfmt: "{{Front}}",
		Afmt: "{{Back}}",
	}
	model.AddTemplate(template1)

	if len(model.Templates) != 1 {
		t.Errorf("Expected 1 template, got %d", len(model.Templates))
	}

	if model.Templates[0].Name != "Template1" {
		t.Errorf("Expected template name to be 'Template1', got '%s'", model.Templates[0].Name)
	}

	if model.Templates[0].Ord != 0 {
		t.Errorf("Expected template ord to be 0, got %d", model.Templates[0].Ord)
	}

	template2 := genanki.Template{
		Name: "Template2",
		Qfmt: "{{Front}} - Question",
		Afmt: "{{Front}} - Answer",
	}
	model.AddTemplate(template2)

	if len(model.Templates) != 2 {
		t.Errorf("Expected 2 templates, got %d", len(model.Templates))
	}

	if model.Templates[1].Name != "Template2" {
		t.Errorf("Expected template name to be 'Template2', got '%s'", model.Templates[1].Name)
	}

	if model.Templates[1].Ord != 1 {
		t.Errorf("Expected template ord to be 1, got %d", model.Templates[1].Ord)
	}
}

func TestSetCSS(t *testing.T) {
	model := genanki.NewModel(1234567890, "Test Model")

	customCSS := `.card { font-family: Arial; font-size: 20px; }`
	model.SetCSS(customCSS)

	if model.CSS != customCSS {
		t.Errorf("Expected CSS to be '%s', got '%s'", customCSS, model.CSS)
	}
}

// Test the convenience functions that use standard IDs
func TestConvenienceFunctions(t *testing.T) {
	// Test StandardBasicModel
	basicModel := genanki.StandardBasicModel("Standard Basic")
	if basicModel.ID != 1607392319 {
		t.Errorf("Expected StandardBasicModel to use ID 1607392319, got %d", basicModel.ID)
	}
	if basicModel.Name != "Standard Basic" {
		t.Errorf("Expected name to be 'Standard Basic', got %s", basicModel.Name)
	}
	if len(basicModel.Fields) != 2 || basicModel.Fields[0].Name != "Front" || basicModel.Fields[1].Name != "Back" {
		t.Errorf("StandardBasicModel has incorrect fields")
	}

	// Test StandardClozeModel
	clozeModel := genanki.StandardClozeModel("Standard Cloze")
	if clozeModel.ID != 1122334455 {
		t.Errorf("Expected StandardClozeModel to use ID 1122334455, got %d", clozeModel.ID)
	}
	if clozeModel.Name != "Standard Cloze" {
		t.Errorf("Expected name to be 'Standard Cloze', got %s", clozeModel.Name)
	}
	if len(clozeModel.Fields) != 2 || clozeModel.Fields[0].Name != "Text" || clozeModel.Fields[1].Name != "Extra" {
		t.Errorf("StandardClozeModel has incorrect fields")
	}

	// Test StandardDeck
	deck := genanki.StandardDeck("Standard Deck", "Standard Deck Description")
	if deck.ID != 1347639657110 {
		t.Errorf("Expected StandardDeck to use ID 1347639657110, got %d", deck.ID)
	}
	if deck.Name != "Standard Deck" {
		t.Errorf("Expected name to be 'Standard Deck', got %s", deck.Name)
	}
	if deck.Desc != "Standard Deck Description" {
		t.Errorf("Expected description to match, got %s", deck.Desc)
	}
}
