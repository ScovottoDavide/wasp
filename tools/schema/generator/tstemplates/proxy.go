package tstemplates

var proxyTs = map[string]string{
	// *******************************
	"proxyContainers": `
$#if array typedefProxyArray
$#if map typedefProxyMap
`,
	// *******************************
	"proxyMethods": `
$#set varID sc.idxMap[sc.Idx$Kind$FldName]
$#if core setCoreVarID
$#if array proxyArray proxyMethods2
`,
	// *******************************
	"proxyMethods2": `
$#if map proxyMap proxyMethods3
`,
	// *******************************
	"proxyMethods3": `
$#if basetype proxyBaseType proxyNewType
`,
	// *******************************
	"setCoreVarID": `
$#set varID wasmlib.Key32.fromString(sc.$Kind$FldName)
`,
	// *******************************
	"proxyArray": `

    $fldName(): sc.ArrayOf$mut$FldType {
		let arrID = wasmlib.getObjectID(this.mapID, $varID, $arrayTypeID|$FldTypeID);
		return new sc.ArrayOf$mut$FldType(arrID);
	}
`,
	// *******************************
	"proxyMap": `
$#if this proxyMapThis proxyMapOther
`,
	// *******************************
	"proxyMapThis": `

    $fldName(): sc.Map$FldMapKey$+To$mut$FldType {
		return new sc.Map$FldMapKey$+To$mut$FldType(this.mapID);
	}
`,
	// *******************************
	"proxyMapOther": `

    $fldName(): sc.Map$FldMapKey$+To$mut$FldType {
		let mapID = wasmlib.getObjectID(this.mapID, $varID, wasmlib.TYPE_MAP);
		return new sc.Map$FldMapKey$+To$mut$FldType(mapID);
	}
`,
	// *******************************
	"proxyBaseType": `

    $fldName(): wasmlib.Sc$mut$FldType {
		return new wasmlib.Sc$mut$FldType(this.mapID, $varID);
	}
`,
	// *******************************
	"proxyNewType": `

    $fldName(): sc.$mut$FldType {
		return new sc.$mut$FldType(this.mapID, $varID);
	}
`,
}