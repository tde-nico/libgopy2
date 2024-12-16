package libgopy

/*
#cgo CFLAGS: -I/usr/include/python3.10 -I/usr/include/python3.10
#cgo LDFLAGS: -L/usr/lib/python3.10/config-3.10-x86_64-linux-gnu -L/usr/lib/x86_64-linux-gnu -lcrypt -ldl -lm -lm -lpython3.10
#include <Python.h>
*/
import "C"
import (
	"fmt"
	"strings"
	"unsafe"
)

var funcs map[string]*C.PyObject

func init() {
	C.Py_Initialize()
	cbytes := C.CString("import sys; sys.path.insert(0, '.')")
	defer C.free(unsafe.Pointer(cbytes))
	C.PyRun_SimpleString(cbytes)

	Load("libtests.test_script2")
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

	funcs = make(map[string]*C.PyObject)
	for i := range int(size) {
		item := C.PyList_GetItem(dir, C.Py_ssize_t(i))
		if item == nil {
			return fmt.Errorf("failed to get item from dir of module: %s", module)
		}

		itemStr := C.PyUnicode_AsUTF8(item)
		if itemStr == nil {
			return fmt.Errorf("failed to convert item to string from dir of module: %s", module)
		}

		if itemStr == nil || C.GoString(itemStr)[0] == '_' {
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
		case []uint8:
			cstr := C.CBytes(v)
			defer C.free(unsafe.Pointer(cstr))
			obj = C.PyBytes_FromStringAndSize((*C.char)(unsafe.Pointer(cstr)), C.Py_ssize_t(len(v)))
		case string:
			cstr := C.CString(v)
			defer C.free(unsafe.Pointer(cstr))
			obj = C.PyUnicode_FromString(cstr)
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
	default:
		return nil, fmt.Errorf("unsupported type: %s", typeGoStr)
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
