from libtests.test_script3 import TEST
import random

LOCAL_GLOBAL_TEST_VAR = 10

def func6(*args):
	print(args)
	return 3.14

def func7(*args):
	print(args)
	print(LOCAL_GLOBAL_TEST_VAR)
	print(TEST)
	if args:
		return args[0]

def func8(*args):
	if args:
		ret = args[0][::-1]
		return ret

def func9(*args):
	return random.choice(list(args)+[[1,"ciao",[b's', 3.14]]])