import json
object_orient='''{
        person1:{
            name:"Alice",
            welcome:"hello"+self.name + "!",
            },
        person2:self.person1{name:"bob"},
}'''

print(object_orient)
import _jsonnet
print(_jsonnet.evaluate_snippet("snippet",object_orient))


function='''
local Person(name="Alice")={
name:name,
welcome:'Hello'+name+"!",
};
{person1:Person(),
person2:Person('bob'),
}
'''

print(_jsonnet.evaluate_snippet("snippet",function))
