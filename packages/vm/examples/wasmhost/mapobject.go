package wasmhost

import (
	"fmt"
)

type MapObject struct {
	vm    *wasmProcessor
	keyId int32
	name  string
}

type ObjFactory func(vm *wasmProcessor) HostObject

func (o *MapObject) checkedObjectId(objId *int32, newObject ObjFactory, typeId int32, expectedTypeId int32) int32 {
	if typeId != expectedTypeId {
		o.error("GetObjectId: Invalid type")
		return 0
	}
	if *objId == 0 {
		*objId = o.vm.TrackObject(newObject(o.vm))
	}
	return *objId
}

func (o *MapObject) error(format string, args ...interface{}) {
	if o.keyId != 0 {
		o.name = o.vm.GetKey(o.keyId)
		o.keyId = 0
	}
	o.vm.SetError(o.name + "." + fmt.Sprintf(format, args...))
}

func (o *MapObject) GetBytes(keyId int32) []byte {
	o.error("GetBytes: Invalid key")
	return []byte(nil)
}

func (o *MapObject) GetInt(keyId int32) int64 {
	o.error("GetInt: Invalid key")
	return 0
}

func (o *MapObject) GetObjectId(keyId int32, typeId int32) int32 {
	o.error("GetObjectId: Invalid key")
	return 0
}

func (o *MapObject) GetString(keyId int32) string {
	o.error("GetString: Invalid key")
	return ""
}

func (o *MapObject) SetBytes(keyId int32, value []byte) {
	o.error("SetBytes: Immutable")
}

func (o *MapObject) SetInt(keyId int32, value int64) {
	o.error("SetInt: Immutable")
}

func (o *MapObject) SetString(keyId int32, value string) {
	o.error("SetString: Immutable")
}
