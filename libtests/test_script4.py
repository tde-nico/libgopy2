def func_list(*args):
	return list(args)

def func_tuple(*args):
	return tuple(args)

def func_dict():
	return {
		'key1': [1, 2, 3],
		'key2': (4, 5.5, None),
		'key3': {'a': 1, 'b': 2},
		'key4': b"Hello, World!",
	}

def func_arg(arg):
	return arg

def func_set():
	return {1, 2, 3, 4}
