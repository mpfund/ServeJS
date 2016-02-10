package htmlcheck

import (
	"encoding/json"
	"io/ioutil"

	"github.com/BlackEspresso/htmlcheck"

	"./../pluginbase"
	"github.com/robertkrimen/otto"
)

func InitPlugin() *pluginbase.Plugin {

	p1 := pluginbase.Plugin{
		Name: "htmlcheck",
		Init: registerVM,
	}

	return &p1
}

func registerVM(vm *otto.Otto) {
	validater := htmlcheck.Validator{}
	obj, _ := vm.Object("({})")
	vm.Set("htmlcheck", obj)
	obj.Set("loadTags", func(c otto.FunctionCall) otto.Value {
		path, err := c.Argument(0).ToString()
		tags, err := LoadTagsFromFile(path)
		if err == nil {
			validater.AddValidTags(tags)
			return otto.TrueValue()
		}
		return otto.FalseValue()
	})
	obj.Set("validate", func(c otto.FunctionCall) otto.Value {
		htmltext, _ := c.Argument(0).ToString()
		errors := validater.ValidateHtmlString(htmltext)
		objs, _ := vm.ToValue(errors)
		return objs
	})
}

func LoadTagsFromFile(path string) ([]htmlcheck.ValidTag, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return []htmlcheck.ValidTag{}, err
	}

	var validTags []htmlcheck.ValidTag
	err = json.Unmarshal(content, &validTags)
	if err != nil {
		return []htmlcheck.ValidTag{}, err
	}

	return validTags, nil
}
