package transformers

import (
	"errors"
)

type Transformation interface {
	GetInTypes() []string
	GetOutType() string
	GetName() string
	DoTransformation(interface{}) interface{}
	Validate() error
}

type Transformer struct {
	Name            string
	Transformations []Transformation
}

func (t *Transformer) Validate() error {
	// validate config
	for _, transformation := range t.Transformations {
		err := transformation.Validate()
		if err != nil {
			log.Errorf("invald transformation configuration")
			return errors.New("invalid transformation configuration")
		}
	}

	// validate type logix
	varType := "string"
	formerName := "collector origin"
	for _, transformation := range t.Transformations {
		if !InStringList(varType, transformation.GetInTypes()) {
			log.Errorf(
				"variable type '%s' from '%s' not matching transformer type '%s' of '%s'",
				varType,
				formerName,
				transformation.GetInTypes(),
				transformation.GetName(),
			)
			return errors.New("transformation types not matching")
		}
		varType = transformation.GetOutType()
		formerName = transformation.GetName()
	}
	return nil
}

func (t *Transformer) DoTransformations(data interface{}) interface{} {
	for _, transformation := range t.Transformations {
		data = transformation.DoTransformation(data)
	}
	return nil
}

func NewTransformation(cnf map[string]interface{}) Transformation {
	return TransformationRegister[cnf["action"].(string)](cnf)
}

func NewTransformer(cnf map[string]interface{}) Transformer {
	transformer := Transformer{}
	transformer.Name = cnf["name"].(string)
	for _, transformationCnf := range cnf["transformations"].([]map[string]interface{}) {
		transformer.Transformations = append(transformer.Transformations, NewTransformation(transformationCnf))
	}
	return transformer
}
