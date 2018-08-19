package util

import (
	"github.com/SpaceHexagon/ecs/object"
)

func CopyObject(valueNode object.Object) object.Object {
	switch valueNode.Type() {
	case "boolean":
		return &object.Boolean{Value: valueNode.(*object.Boolean).Value}
	case "int":
		return &object.Integer{Value: valueNode.(*object.Integer).Value}
	case "float":
		return &object.Float{Value: valueNode.(*object.Float).Value}
	case "string":
		return &object.String{Value: valueNode.(*object.String).Value}
	case "array":
		return &object.Array{Elements: valueNode.(*object.Array).Elements}
	case "hash":
		return CopyHashMap(valueNode)
	case "function":
	case "BUILTIN":
		return valueNode
	default:
		return &object.Null{}
	}
	return &object.Null{}
}

func CopyHashMap(data object.Object) object.Object {
	pairData := data.(*object.Hash).Pairs

	pairs := make(map[object.HashKey]object.HashPair)

	for key, pair := range pairData {
		valueNode := pair.Value
		keyNode := pair.Key
		isStatic := pair.Modifiers != nil && hasModifier(pair.Modifiers, 1)

		var (
			NewValue object.Object
			newPair  object.HashPair
		)

		if isStatic {
			pairs[key] = pair
		} else {
			NewValue = CopyObject(valueNode)
			newPair = object.HashPair{Key: keyNode, Value: NewValue}
			if pair.Modifiers != nil {
				newPair.Modifiers = pair.Modifiers
			}
			pairs[key] = newPair
		}
	}

	return &object.Hash{Pairs: pairs}
}

func hasModifier(modifiers []int64, modifier int64) bool {
	for mod := range modifiers {
		if modifiers[mod] == modifier {
			return true
		}
	}
	return false
}

func MakeBuiltinClass(className string, fields []StringObjectPair) object.Hash {
	instance := MakeBuiltinInterface(fields)

	// instance.Constructor = instance.Pairs.Get(&object.String(className).HashKey())
	// instance.className = className
	// instance.Pairs.builtin = &object.HashPair{Key: strBuiltin.HashKey(), Value: TRUE}
	return instance
}

type StringObjectPair struct {
	name string
	obj  object.Object
}

func MakeBuiltinInterface(methods []StringObjectPair) object.Hash {
	pairs := make(map[object.HashKey]object.HashPair)
	for _, v := range methods {
		key := &object.String{Value: v.name}
		pairs[key.HashKey()] = object.HashPair{
			Key:   key,
			Value: v.obj,
		}
	}

	return object.Hash{Pairs: pairs}
}

// func addMethod (allMethods, methodName string, contextName string, builtinFn object.Builtin) {
// 	allMethods = append(allMethods, &{[methodName]: &object.HashPair{
// 		Key: new object.String(methodName),
// 		Value: new object.Builtin(builtinFn, contextName)
// 	}});
// }

func NativeListToArray(items []interface{}) object.Array {
	var (
		elements []object.Object
	)
	for _, element := range items {
		switch element.(type) {
		case string:
			elements = append(elements, &object.String{Value: element.(string)})
		case int64:
			elements = append(elements, &object.Integer{Value: element.(int64)})
		case float64:
			elements = append(elements, &object.Float{Value: element.(float64)})
		case bool:
			elements = append(elements, &object.Boolean{Value: element.(bool)})
		// case []interface{}:
		// 	elements = append(elements, (nativeListToArray(element))
		// case interface{}:
		// 	elements = append(elements, nativeObjToMap(element.(map[string]interface{})).(interface{}(object.Object).(type)))
		default:
			elements = append(elements, &object.Null{})
		}
	}

	return object.Array{Elements: elements} //obj
}

// func nativeObjToMap (obj: {[key: string]: any} = {}): object.Hash => {
func NativeObjToMap(obj map[string]interface{}) object.Hash {
	newMap := object.Hash{Pairs: nil}

	for objectKey, data := range obj {
		var (
			value object.Object
		)

		switch data.(type) {
		case string:
			value = &object.String{Value: data.(string)}
			break
		case int64:
			value = &object.Integer{Value: data.(int64)}
			break
		case float64:
			value = &object.Float{Value: data.(float64)}
		case bool:
			value = &object.Boolean{Value: data.(bool)}
			break
		// case interface{}:
		// 	value = nativeObjToMap(data)
		// 	break
		// case :
		// 	console.log("native function", data)
		// 	// need to figure this out
		// 	// new object.Builtin(builtinFn, contextName)
		// 	break
		default:

		}
		key := &object.String{Value: objectKey}
		newMap.Pairs[key.HashKey()] = object.HashPair{
			Key:   key,
			Value: value,
		}
	}

	return newMap
}
