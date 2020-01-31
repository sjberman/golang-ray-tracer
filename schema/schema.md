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
		 - Type: `object`
		 - <i id="#/properties/shapes/items">path: #/properties/shapes/items</i>
		 - **_Properties_**
			 - <b id="#/properties/shapes/items/properties/name">name</b>
				 - _Name of the shape._
				 - Type: `string`
				 - <i id="#/properties/shapes/items/properties/name">path: #/properties/shapes/items/properties/name</i>
			 - <b id="#/properties/shapes/items/properties/type">type</b> `required`
				 - _Type of shape._
				 - Type: `string`
				 - <i id="#/properties/shapes/items/properties/type">path: #/properties/shapes/items/properties/type</i>
				 - The value is restricted to the following: 
					 1. _"cone"_
					 2. _"cube"_
					 3. _"cylinder"_
					 4. _"plane"_
					 5. _"sphere"_
					 6. _"glassSphere"_
			 - <b id="#/properties/shapes/items/properties/transform">transform</b>
				 - _Ways to transform the shape._
				 - <i id="#/properties/shapes/items/properties/transform">path: #/properties/shapes/items/properties/transform</i>
				 - &#36;ref: [#/definitions/transform](#/definitions/transform)
			 - <b id="#/properties/shapes/items/properties/material">material</b>
				 - _Material of the shape_
				 - <i id="#/properties/shapes/items/properties/material">path: #/properties/shapes/items/properties/material</i>
				 - &#36;ref: [#/definitions/material](#/definitions/material)
			 - <b id="#/properties/shapes/items/properties/closed">closed</b>
				 - _Closed caps for cone or cylinder._
				 - Type: `boolean`
				 - <i id="#/properties/shapes/items/properties/closed">path: #/properties/shapes/items/properties/closed</i>
			 - <b id="#/properties/shapes/items/properties/minimum">minimum</b>
				 - _Minimum value for cone or cylinder._
				 - Type: `number`
				 - <i id="#/properties/shapes/items/properties/minimum">path: #/properties/shapes/items/properties/minimum</i>
			 - <b id="#/properties/shapes/items/properties/maximum">maximum</b>
				 - _Maximum value for cone or cylinder._
				 - Type: `number`
				 - <i id="#/properties/shapes/items/properties/maximum">path: #/properties/shapes/items/properties/maximum</i>
 - <b id="#/properties/groups">groups</b>
	 - Type: `array`
	 - <i id="#/properties/groups">path: #/properties/groups</i>
		 - **_Items_**
		 - Type: `object`
		 - <i id="#/properties/groups/items">path: #/properties/groups/items</i>
		 - **_Properties_**
			 - <b id="#/properties/groups/items/properties/name">name</b>
				 - _Name of the group._
				 - Type: `string`
				 - <i id="#/properties/groups/items/properties/name">path: #/properties/groups/items/properties/name</i>
			 - <b id="#/properties/groups/items/properties/transform">transform</b>
				 - _Ways to transform the group._
				 - <i id="#/properties/groups/items/properties/transform">path: #/properties/groups/items/properties/transform</i>
				 - &#36;ref: [#/definitions/transform](#/definitions/transform)
			 - <b id="#/properties/groups/items/properties/material">material</b>
				 - _Material of the group._
				 - <i id="#/properties/groups/items/properties/material">path: #/properties/groups/items/properties/material</i>
				 - &#36;ref: [#/definitions/material](#/definitions/material)
			 - <b id="#/properties/groups/items/properties/children">children</b> `required`
				 - _Names of objects contained in the group._
				 - Type: `array`
				 - <i id="#/properties/groups/items/properties/children">path: #/properties/groups/items/properties/children</i>
					 - **_Items_**
					 - Type: `string`
					 - <i id="#/properties/groups/items/properties/children/items">path: #/properties/groups/items/properties/children/items</i>
 - <b id="#/properties/csgs">csgs</b>
	 - Type: `array`
	 - <i id="#/properties/csgs">path: #/properties/csgs</i>
		 - **_Items_**
		 - Type: `object`
		 - <i id="#/properties/csgs/items">path: #/properties/csgs/items</i>
		 - **_Properties_**
			 - <b id="#/properties/csgs/items/properties/name">name</b>
				 - _Name of the CSG._
				 - Type: `string`
				 - <i id="#/properties/csgs/items/properties/name">path: #/properties/csgs/items/properties/name</i>
			 - <b id="#/properties/csgs/items/properties/operation">operation</b> `required`
				 - _Type of operation to perform on the objects._
				 - Type: `string`
				 - <i id="#/properties/csgs/items/properties/operation">path: #/properties/csgs/items/properties/operation</i>
				 - The value is restricted to the following: 
					 1. _"union"_
					 2. _"intersection"_
					 3. _"difference"_
			 - <b id="#/properties/csgs/items/properties/leftChild">leftChild</b> `required`
				 - _Name of the left child object._
				 - Type: `string`
				 - <i id="#/properties/csgs/items/properties/leftChild">path: #/properties/csgs/items/properties/leftChild</i>
			 - <b id="#/properties/csgs/items/properties/rightChild">rightChild</b> `required`
				 - _Name of the right child object._
				 - Type: `string`
				 - <i id="#/properties/csgs/items/properties/rightChild">path: #/properties/csgs/items/properties/rightChild</i>
			 - <b id="#/properties/csgs/items/properties/transform">transform</b>
				 - _Ways to transform the CSG._
				 - <i id="#/properties/csgs/items/properties/transform">path: #/properties/csgs/items/properties/transform</i>
				 - &#36;ref: [#/definitions/transform](#/definitions/transform)
			 - <b id="#/properties/csgs/items/properties/material">material</b>
				 - _Material of the CSG._
				 - <i id="#/properties/csgs/items/properties/material">path: #/properties/csgs/items/properties/material</i>
				 - &#36;ref: [#/definitions/material](#/definitions/material)
 - <b id="#/properties/files">files</b>
	 - Type: `array`
	 - <i id="#/properties/files">path: #/properties/files</i>
		 - **_Items_**
		 - Type: `object`
		 - <i id="#/properties/files/items">path: #/properties/files/items</i>
		 - **_Properties_**
			 - <b id="#/properties/files/items/properties/name">name</b>
				 - _Name of the object to be created._
				 - Type: `string`
				 - <i id="#/properties/files/items/properties/name">path: #/properties/files/items/properties/name</i>
			 - <b id="#/properties/files/items/properties/file">file</b> `required`
				 - _OBJ filename to be loaded_
				 - Type: `string`
				 - <i id="#/properties/files/items/properties/file">path: #/properties/files/items/properties/file</i>
# definitions

 - Type: `array`
 - <i id="#/definitions/tuple">path: #/definitions/tuple</i>
 - Item Count: between 3 and 3
	 - **_Items_**
	 - Type: `number`
	 - <i id="#/definitions/tuple/items">path: #/definitions/tuple/items</i>
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