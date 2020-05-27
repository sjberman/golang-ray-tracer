# Ray tracer scene

Type: `object`

<i id="#">path: #</i>

&#36;schema: [http://json-schema.org/draft-07/schema#](http://json-schema.org/draft-07/schema#)

**_Properties_**

 - <b id="#/properties/camera">camera</b> `required`
	 - Type: `object`
	 - <i id="#/properties/camera">path: #/properties/camera</i>
	 - **_Properties_**
		 - <b id="#/properties/camera/properties/width">width</b> `required`
			 - _Width of the image._
			 - Type: `integer`
			 - <i id="#/properties/camera/properties/width">path: #/properties/camera/properties/width</i>
		 - <b id="#/properties/camera/properties/height">height</b> `required`
			 - _Height of the image._
			 - Type: `integer`
			 - <i id="#/properties/camera/properties/height">path: #/properties/camera/properties/height</i>
		 - <b id="#/properties/camera/properties/field-of-view">field-of-view</b> `required`
			 - _Field of view in degrees._
			 - Type: `number`
			 - <i id="#/properties/camera/properties/field-of-view">path: #/properties/camera/properties/field-of-view</i>
		 - <b id="#/properties/camera/properties/from">from</b> `required`
			 - _Origin of the camera._
			 - <i id="#/properties/camera/properties/from">path: #/properties/camera/properties/from</i>
			 - &#36;ref: [#/definitions/tuple](#/definitions/tuple)
		 - <b id="#/properties/camera/properties/to">to</b> `required`
			 - _Where the camera looks._
			 - <i id="#/properties/camera/properties/to">path: #/properties/camera/properties/to</i>
			 - &#36;ref: [#/definitions/tuple](#/definitions/tuple)
		 - <b id="#/properties/camera/properties/up">up</b> `required`
			 - _The up direction._
			 - <i id="#/properties/camera/properties/up">path: #/properties/camera/properties/up</i>
			 - &#36;ref: [#/definitions/tuple](#/definitions/tuple)
 - <b id="#/properties/lights">lights</b> `required`
	 - Type: `array`
	 - <i id="#/properties/lights">path: #/properties/lights</i>
		 - **_Items_**
		 - Type: `object`
		 - <i id="#/properties/lights/items">path: #/properties/lights/items</i>
		 - **_Properties_**
			 - <b id="#/properties/lights/items/properties/at">at</b> `required`
				 - _Position of the light._
				 - <i id="#/properties/lights/items/properties/at">path: #/properties/lights/items/properties/at</i>
				 - &#36;ref: [#/definitions/tuple](#/definitions/tuple)
			 - <b id="#/properties/lights/items/properties/intensity">intensity</b> `required`
				 - _Color of the light._
				 - <i id="#/properties/lights/items/properties/intensity">path: #/properties/lights/items/properties/intensity</i>
				 - &#36;ref: [#/definitions/tuple](#/definitions/tuple)
 - <b id="#/properties/shapes">shapes</b>
	 - Type: `array`
	 - <i id="#/properties/shapes">path: #/properties/shapes</i>
		 - **_Items_**
		 - <i id="#/properties/shapes/items">path: #/properties/shapes/items</i>
		 - &#36;ref: [#/definitions/shape](#/definitions/shape)
 - <b id="#/properties/groups">groups</b>
	 - Type: `array`
	 - <i id="#/properties/groups">path: #/properties/groups</i>
		 - **_Items_**
		 - <i id="#/properties/groups/items">path: #/properties/groups/items</i>
		 - &#36;ref: [#/definitions/group](#/definitions/group)
 - <b id="#/properties/csgs">csgs</b>
	 - Type: `array`
	 - <i id="#/properties/csgs">path: #/properties/csgs</i>
		 - **_Items_**
		 - <i id="#/properties/csgs/items">path: #/properties/csgs/items</i>
		 - &#36;ref: [#/definitions/csg](#/definitions/csg)
 - <b id="#/properties/files">files</b>
	 - Type: `array`
	 - <i id="#/properties/files">path: #/properties/files</i>
		 - **_Items_**
		 - Type: `object`
		 - <i id="#/properties/files/items">path: #/properties/files/items</i>
		 - **_Properties_**
			 - <b id="#/properties/files/items/properties/name">name</b> `required`
				 - _Name of the object to be created._
				 - Type: `string`
				 - <i id="#/properties/files/items/properties/name">path: #/properties/files/items/properties/name</i>
			 - <b id="#/properties/files/items/properties/file">file</b> `required`
				 - _OBJ filename to be loaded_
				 - Type: `string`
				 - <i id="#/properties/files/items/properties/file">path: #/properties/files/items/properties/file</i>
			 - <b id="#/properties/files/items/properties/transform">transform</b>
				 - _Ways to transform the object._
				 - <i id="#/properties/files/items/properties/transform">path: #/properties/files/items/properties/transform</i>
				 - &#36;ref: [#/definitions/transform](#/definitions/transform)
			 - <b id="#/properties/files/items/properties/material">material</b>
				 - _Material of the object._
				 - <i id="#/properties/files/items/properties/material">path: #/properties/files/items/properties/material</i>
				 - &#36;ref: [#/definitions/material](#/definitions/material)
# definitions

 - Type: `array`
 - <i id="#/definitions/tuple">path: #/definitions/tuple</i>
 - Item Count: between 3 and 3
	 - **_Items_**
	 - Type: `number`
	 - <i id="#/definitions/tuple/items">path: #/definitions/tuple/items</i>
 - _A 3 dimensional shape._
 - Type: `object`
 - <i id="#/definitions/shape">path: #/definitions/shape</i>
 - **_Properties_**
	 - <b id="#/definitions/shape/properties/name">name</b> `required`
		 - _Name of the shape._
		 - Type: `string`
		 - <i id="#/definitions/shape/properties/name">path: #/definitions/shape/properties/name</i>
	 - <b id="#/definitions/shape/properties/type">type</b> `required`
		 - _Type of shape._
		 - Type: `string`
		 - <i id="#/definitions/shape/properties/type">path: #/definitions/shape/properties/type</i>
		 - The value is restricted to the following: 
			 1. _"cone"_
			 2. _"cube"_
			 3. _"cylinder"_
			 4. _"plane"_
			 5. _"sphere"_
			 6. _"glassSphere"_
	 - <b id="#/definitions/shape/properties/transform">transform</b>
		 - _Ways to transform the shape._
		 - <i id="#/definitions/shape/properties/transform">path: #/definitions/shape/properties/transform</i>
		 - &#36;ref: [#/definitions/transform](#/definitions/transform)
	 - <b id="#/definitions/shape/properties/material">material</b>
		 - _Material of the shape_
		 - <i id="#/definitions/shape/properties/material">path: #/definitions/shape/properties/material</i>
		 - &#36;ref: [#/definitions/material](#/definitions/material)
	 - <b id="#/definitions/shape/properties/closed">closed</b>
		 - _Closed caps for cone or cylinder._
		 - Type: `boolean`
		 - <i id="#/definitions/shape/properties/closed">path: #/definitions/shape/properties/closed</i>
	 - <b id="#/definitions/shape/properties/minimum">minimum</b>
		 - _Minimum value for cone or cylinder._
		 - Type: `number`
		 - <i id="#/definitions/shape/properties/minimum">path: #/definitions/shape/properties/minimum</i>
	 - <b id="#/definitions/shape/properties/maximum">maximum</b>
		 - _Maximum value for cone or cylinder._
		 - Type: `number`
		 - <i id="#/definitions/shape/properties/maximum">path: #/definitions/shape/properties/maximum</i>
	 - <b id="#/definitions/shape/properties/inherits">inherits</b>
		 - _Inherits the properties from another shape._
		 - Type: `string`
		 - <i id="#/definitions/shape/properties/inherits">path: #/definitions/shape/properties/inherits</i>
 - _A group of objects._
 - Type: `object`
 - <i id="#/definitions/group">path: #/definitions/group</i>
 - **_Properties_**
	 - <b id="#/definitions/group/properties/name">name</b> `required`
		 - _Name of the group._
		 - Type: `string`
		 - <i id="#/definitions/group/properties/name">path: #/definitions/group/properties/name</i>
	 - <b id="#/definitions/group/properties/transform">transform</b>
		 - _Ways to transform the group._
		 - <i id="#/definitions/group/properties/transform">path: #/definitions/group/properties/transform</i>
		 - &#36;ref: [#/definitions/transform](#/definitions/transform)
	 - <b id="#/definitions/group/properties/material">material</b>
		 - _Material of the group._
		 - <i id="#/definitions/group/properties/material">path: #/definitions/group/properties/material</i>
		 - &#36;ref: [#/definitions/material](#/definitions/material)
	 - <b id="#/definitions/group/properties/children">children</b> `required`
		 - _Objects contained in the group._
		 - Type: `array`
		 - <i id="#/definitions/group/properties/children">path: #/definitions/group/properties/children</i>
			 - **_Items_**
			 - <i id="#/definitions/group/properties/children/items">path: #/definitions/group/properties/children/items</i>
			 - &#36;ref: [#/definitions/objectShell](#/definitions/objectShell)
 - _Constructive Solid Geometry (combination of two shapes)._
 - Type: `object`
 - <i id="#/definitions/csg">path: #/definitions/csg</i>
 - **_Properties_**
	 - <b id="#/definitions/csg/properties/name">name</b> `required`
		 - _Name of the CSG._
		 - Type: `string`
		 - <i id="#/definitions/csg/properties/name">path: #/definitions/csg/properties/name</i>
	 - <b id="#/definitions/csg/properties/operation">operation</b> `required`
		 - _Type of operation to perform on the objects._
		 - Type: `string`
		 - <i id="#/definitions/csg/properties/operation">path: #/definitions/csg/properties/operation</i>
		 - The value is restricted to the following: 
			 1. _"union"_
			 2. _"intersection"_
			 3. _"difference"_
	 - <b id="#/definitions/csg/properties/leftChild">leftChild</b> `required`
		 - _Left child object._
		 - <i id="#/definitions/csg/properties/leftChild">path: #/definitions/csg/properties/leftChild</i>
		 - &#36;ref: [#/definitions/objectShell](#/definitions/objectShell)
	 - <b id="#/definitions/csg/properties/rightChild">rightChild</b> `required`
		 - _Right child object._
		 - <i id="#/definitions/csg/properties/rightChild">path: #/definitions/csg/properties/rightChild</i>
		 - &#36;ref: [#/definitions/objectShell](#/definitions/objectShell)
	 - <b id="#/definitions/csg/properties/transform">transform</b>
		 - _Ways to transform the CSG._
		 - <i id="#/definitions/csg/properties/transform">path: #/definitions/csg/properties/transform</i>
		 - &#36;ref: [#/definitions/transform](#/definitions/transform)
	 - <b id="#/definitions/csg/properties/material">material</b>
		 - _Material of the CSG._
		 - <i id="#/definitions/csg/properties/material">path: #/definitions/csg/properties/material</i>
		 - &#36;ref: [#/definitions/material](#/definitions/material)
 - Type: `object`
 - <i id="#/definitions/objectShell">path: #/definitions/objectShell</i>
 - **_Properties_**
	 - <b id="#/definitions/objectShell/properties/name">name</b> `required`
		 - Type: `string`
		 - <i id="#/definitions/objectShell/properties/name">path: #/definitions/objectShell/properties/name</i>
	 - <b id="#/definitions/objectShell/properties/transform">transform</b>
		 - <i id="#/definitions/objectShell/properties/transform">path: #/definitions/objectShell/properties/transform</i>
		 - &#36;ref: [#/definitions/transform](#/definitions/transform)
	 - <b id="#/definitions/objectShell/properties/material">material</b>
		 - <i id="#/definitions/objectShell/properties/material">path: #/definitions/objectShell/properties/material</i>
		 - &#36;ref: [#/definitions/material](#/definitions/material)
 - Type: `array`
 - <i id="#/definitions/transform">path: #/definitions/transform</i>
	 - **_Items_**
	 - Type: `object`
	 - <i id="#/definitions/transform/items">path: #/definitions/transform/items</i>
	 - **_Properties_**
		 - <b id="#/definitions/transform/items/properties/type">type</b> `required`
			 - _Type of transform to be performed._
			 - Type: `string`
			 - <i id="#/definitions/transform/items/properties/type">path: #/definitions/transform/items/properties/type</i>
			 - The value is restricted to the following: 
				 1. _"translate"_
				 2. _"scale"_
				 3. _"rotate"_
				 4. _"shear"_
		 - <b id="#/definitions/transform/items/properties/values">values</b> `required`
			 - _x,y,z values for the transformation._
			 - Type: `array`
			 - <i id="#/definitions/transform/items/properties/values">path: #/definitions/transform/items/properties/values</i>
				 - **_Items_**
				 - Type: `number`
				 - <i id="#/definitions/transform/items/properties/values/items">path: #/definitions/transform/items/properties/values/items</i>
 - Type: `object`
 - <i id="#/definitions/material">path: #/definitions/material</i>
 - **_Properties_**
	 - <b id="#/definitions/material/properties/color">color</b>
		 - <i id="#/definitions/material/properties/color">path: #/definitions/material/properties/color</i>
		 - &#36;ref: [#/definitions/tuple](#/definitions/tuple)
	 - <b id="#/definitions/material/properties/pattern">pattern</b>
		 - Type: `object`
		 - <i id="#/definitions/material/properties/pattern">path: #/definitions/material/properties/pattern</i>
		 - **_Properties_**
			 - <b id="#/definitions/material/properties/pattern/properties/type">type</b> `required`
				 - Type: `string`
				 - <i id="#/definitions/material/properties/pattern/properties/type">path: #/definitions/material/properties/pattern/properties/type</i>
				 - The value is restricted to the following: 
					 1. _"checker"_
					 2. _"gradient"_
					 3. _"ring"_
					 4. _"stripe"_
			 - <b id="#/definitions/material/properties/pattern/properties/color1">color1</b> `required`
				 - <i id="#/definitions/material/properties/pattern/properties/color1">path: #/definitions/material/properties/pattern/properties/color1</i>
				 - &#36;ref: [#/definitions/tuple](#/definitions/tuple)
			 - <b id="#/definitions/material/properties/pattern/properties/color2">color2</b> `required`
				 - <i id="#/definitions/material/properties/pattern/properties/color2">path: #/definitions/material/properties/pattern/properties/color2</i>
				 - &#36;ref: [#/definitions/tuple](#/definitions/tuple)
			 - <b id="#/definitions/material/properties/pattern/properties/transform">transform</b>
				 - <i id="#/definitions/material/properties/pattern/properties/transform">path: #/definitions/material/properties/pattern/properties/transform</i>
				 - &#36;ref: [#/definitions/transform](#/definitions/transform)
	 - <b id="#/definitions/material/properties/ambient">ambient</b>
		 - Type: `number`
		 - <i id="#/definitions/material/properties/ambient">path: #/definitions/material/properties/ambient</i>
	 - <b id="#/definitions/material/properties/diffuse">diffuse</b>
		 - Type: `number`
		 - <i id="#/definitions/material/properties/diffuse">path: #/definitions/material/properties/diffuse</i>
	 - <b id="#/definitions/material/properties/specular">specular</b>
		 - Type: `number`
		 - <i id="#/definitions/material/properties/specular">path: #/definitions/material/properties/specular</i>
	 - <b id="#/definitions/material/properties/shininess">shininess</b>
		 - Type: `number`
		 - <i id="#/definitions/material/properties/shininess">path: #/definitions/material/properties/shininess</i>
	 - <b id="#/definitions/material/properties/reflective">reflective</b>
		 - Type: `number`
		 - <i id="#/definitions/material/properties/reflective">path: #/definitions/material/properties/reflective</i>
	 - <b id="#/definitions/material/properties/transparency">transparency</b>
		 - Type: `number`
		 - <i id="#/definitions/material/properties/transparency">path: #/definitions/material/properties/transparency</i>
	 - <b id="#/definitions/material/properties/refractiveIndex">refractiveIndex</b>
		 - Type: `number`
		 - <i id="#/definitions/material/properties/refractiveIndex">path: #/definitions/material/properties/refractiveIndex</i>
	 - <b id="#/definitions/material/properties/shadow">shadow</b>
		 - Type: `boolean`
		 - <i id="#/definitions/material/properties/shadow">path: #/definitions/material/properties/shadow</i>

_Generated with [json-schema-md-doc](https://brianwendt.github.io/json-schema-md-doc/)_