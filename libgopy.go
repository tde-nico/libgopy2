package libgopy2

/*
#cgo pkg-config: python-3.12-embed
#include <Python.h>
*/
import "C"
import (
	"fmt"
	"strings"
	"unsafe"
)

var funcs map[string]*C.PyObject

func Init() {
	C.Py_Initialize()
	cbytes := C.CString("import sys; sys.path.insert(0, '.')")
	defer C.free(unsafe.Pointer(cbytes))
	C.PyRun_SimpleString(cbytes)

	funcs = make(map[string]*C.PyObject)
}

func Finalize() {
	C.Py_Finalize()
}

func Load(module string) error {
	cstr := C.CString(module)
	defer C.free(unsafe.Pointer(cstr))

	name := C.PyUnicode_FromString(cstr)
	if name == nil {
		return fmt.Errorf("failed to convert \"%s\" into string", module)
	}

	loadModule := C.PyImport_Import(name)
	if loadModule == nil {
		return fmt.Errorf("failed to load module: %s", module)
	}

	dir := C.PyObject_Dir(loadModule)
	if dir == nil {
		return fmt.Errorf("failed to get dir of module: %s", module)
	}
	defer C.Py_DecRef(dir)

	size := C.PyList_Size(dir)
	if size == -1 {
		return fmt.Errorf("failed to get size of dir of module: %s", module)
	}

	for i := range int(size) {
		item := C.PyList_GetItem(dir, C.Py_ssize_t(i))
		if item == nil {
			return fmt.Errorf("failed to get item from dir of module: %s", module)
		}

		itemStr := C.PyUnicode_AsUTF8(item)
		if itemStr == nil {
			return fmt.Errorf("failed to convert item to string from dir of module: %s", module)
		}

		if C.GoString(itemStr)[0] == '_' {
			continue
		}

		f := C.PyObject_GetAttrString(loadModule, itemStr)
		if f == nil {
			return fmt.Errorf("failed to get attribute from module: %s", module)
		}

		funcs[C.GoString(itemStr)] = f
	}

	return nil
}

func setupArgs(args []any) (*C.PyObject, error) {
	pyargs := C.PyTuple_New(C.Py_ssize_t(len(args)))
	if pyargs == nil {
		return nil, fmt.Errorf("failed to create tuple")
	}

	for i, arg := range args {
		var obj *C.PyObject

		switch v := arg.(type) {
		case int64:
			obj = C.PyLong_FromLongLong(C.longlong(v))
		case int32:
			obj = C.PyLong_FromLong(C.long(v))
		case int16:
			obj = C.PyLong_FromLong(C.long(v))
		case int8:
			obj = C.PyLong_FromLong(C.long(v))
		case int:
			obj = C.PyLong_FromLong(C.long(v))
		case uint64:
			obj = C.PyLong_FromUnsignedLongLong(C.ulonglong(v))
		case uint32:
			obj = C.PyLong_FromUnsignedLong(C.ulong(v))
		case uint16:
			obj = C.PyLong_FromUnsignedLong(C.ulong(v))
		case uint8:
			obj = C.PyLong_FromUnsignedLong(C.ulong(v))
		case uint:
			obj = C.PyLong_FromUnsignedLong(C.ulong(v))
		case float64:
			obj = C.PyFloat_FromDouble(C.double(v))
		case float32:
			obj = C.PyFloat_FromDouble(C.double(v))
		case bool:
			if v {
				obj = C.Py_True
			} else {
				obj = C.Py_False
			}
			C.Py_IncRef(obj)
		case nil:
			obj = C.Py_None
			C.Py_IncRef(obj)
		case []uint8:
			cstr := C.CBytes(v)
			defer C.free(unsafe.Pointer(cstr))
			obj = C.PyBytes_FromStringAndSize((*C.char)(unsafe.Pointer(cstr)), C.Py_ssize_t(len(v)))
		case string:
			cstr := C.CString(v)
			defer C.free(unsafe.Pointer(cstr))
			obj = C.PyUnicode_FromString(cstr)
		case []any:
			list := C.PyList_New(C.Py_ssize_t(len(v)))
			if list == nil {
				return nil, fmt.Errorf("failed to create list")
			}
			for i, item := range v {
				pyItem, err := setupArgs([]any{item})
				if err != nil {
					return nil, err
				}
				C.PyList_SetItem(list, C.Py_ssize_t(i), C.PyTuple_GetItem(pyItem, 0))
			}
			obj = list
		case map[any]any:
			dict := C.PyDict_New()
			if dict == nil {
				return nil, fmt.Errorf("failed to create dict")
			}
			for key, value := range v {
				pyKey, err := setupArgs([]any{key})
				if err != nil {
					return nil, err
				}
				pyValue, err := setupArgs([]any{value})
				if err != nil {
					return nil, err
				}
				C.PyDict_SetItem(dict, C.PyTuple_GetItem(pyKey, 0), C.PyTuple_GetItem(pyValue, 0))
			}
			obj = dict
		case *C.PyObject:
			obj = v
		default:
			return nil, fmt.Errorf("unknown type: %T", v)
		}

		if obj == nil {
			return nil, fmt.Errorf("failed to create object")
		}

		C.PyTuple_SetItem(pyargs, C.Py_ssize_t(i), obj)
	}

	return pyargs, nil
}

func parsePyObject(obj *C.PyObject) (any, error) {
	t := C.PyObject_Type(obj)
	typeStr := C.PyObject_Str(t)
	defer C.Py_DecRef(t)
	defer C.Py_DecRef(typeStr)

	cTypeStr := C.PyUnicode_AsUTF8(typeStr)
	if cTypeStr == nil {
		return nil, fmt.Errorf("failed to convert type to string")
	}

	typeGoStr := strings.SplitN(C.GoString(cTypeStr), "'", 3)[1]

	switch typeGoStr {
	case "int":
		return int64(C.PyLong_AsLong(obj)), nil
	case "float":
		return float64(C.PyFloat_AsDouble(obj)), nil
	case "bool":
		if obj == C.Py_True {
			return true, nil
		} else {
			return false, nil
		}
	case "str":
		return C.GoString(C.PyUnicode_AsUTF8(obj)), nil
	case "bytes":
		return C.GoBytes(unsafe.Pointer(C.PyBytes_AsString(obj)), C.int(C.PyBytes_Size(obj))), nil
	case "list":
		size := C.PyList_Size(obj)
		list := make([]any, size)
		for i := range int(size) {
			item := C.PyList_GetItem(obj, C.Py_ssize_t(i))
			parsedItem, err := parsePyObject(item)
			if err != nil {
				return nil, err
			}
			list[i] = parsedItem
		}
		return list, nil
	case "tuple":
		size := C.PyTuple_Size(obj)
		tuple := make([]any, size)
		for i := range int(size) {
			item := C.PyTuple_GetItem(obj, C.Py_ssize_t(i))
			parsedItem, err := parsePyObject(item)
			if err != nil {
				return nil, err
			}
			tuple[i] = parsedItem
		}
		return tuple, nil
	case "dict":
		size := C.PyDict_Size(obj)
		dict := make(map[any]any)
		keys := C.PyDict_Keys(obj)
		for i := range int(size) {
			key := C.PyList_GetItem(keys, C.Py_ssize_t(i))
			value := C.PyDict_GetItem(obj, key)
			parsedKey, err := parsePyObject(key)
			if err != nil {
				return nil, err
			}
			parsedValue, err := parsePyObject(value)
			if err != nil {
				return nil, err
			}
			dict[parsedKey] = parsedValue
		}
		return dict, nil
	case "set":
		size := C.PySet_Size(obj)
		set := make(map[any]bool, size)
		iter := C.PyObject_GetIter(obj)
		for range int(size) {
			item := C.PyIter_Next(iter)
			if item == nil {
				break
			}
			parsedItem, err := parsePyObject(item)
			if err != nil {
				return nil, err
			}
			set[parsedItem] = true
		}
		return set, nil
	case "NoneType":
		return nil, nil
	default:
		return obj, nil
	}
}

func Call(name string, args ...any) (any, error) {
	f, ok := funcs[name]
	if !ok {
		return nil, fmt.Errorf("function not found: %s", name)
	}

	pyArgs, err := setupArgs(args)
	if err != nil {
		return nil, err
	}

	res := C.PyObject_CallObject(f, pyArgs)
	if res == nil {
		return nil, fmt.Errorf("failed to call function: %s", name)
	}

	ret, err := parsePyObject(res)
	if err != nil {
		return nil, err
	}

	return ret, nil
}
