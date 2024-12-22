class TestScript:
	def __init__(self, value=None):
		self.value = value
	
	def __repr__(self):
		return f"TestScript({self.value})"


def create_pyobj(value):
	return TestScript(value)

def test_pyobj(obj):
	print(obj)
	return obj.value
