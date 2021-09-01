// Copyright (c) Microsoft. All rights reserved.
// Licensed under the MIT license. See LICENSE.md file in the project root for full license information.

package common

// #include <stdlib.h>
// #include <speechapi_c_common.h>
// #include <speechapi_c_property_bag.h>
import "C"
import "unsafe"

// PropertyCollection is a class to retrieve or set a property value from a property collection.
type PropertyCollection struct {
	handle C.SPXHANDLE
}

// GetProperty returns value of a property.
// If the property value is not defined, the specified default value is returned.
func (properties PropertyCollection) GetProperty(id PropertyID, defaultValue string) string {
	defValue := C.CString(defaultValue)
	defer C.free(unsafe.Pointer(defValue))
	value := C.property_bag_get_string(properties.handle, (C.int)(id), nil, defValue)
	goValue := C.GoString(value)
	C.property_bag_free_string(value)
	return goValue
}

// GetPropertyByString returns value of a property.
// If the property value is not defined, the specified default value is returned.
func (properties PropertyCollection) GetPropertyByString(name string, defaultValue string) string {
	defValue := C.CString(defaultValue)
	defer C.free(unsafe.Pointer(defValue))
	n := C.CString(name)
	defer C.free(unsafe.Pointer(n))
	value := C.property_bag_get_string(properties.handle, -1, n, defValue)
	goValue := C.GoString(value)
	C.property_bag_free_string(value)
	return goValue
}

// SetProperty sets the value of a property.
func (properties PropertyCollection) SetProperty(id PropertyID, value string) error {
	v := C.CString(value)
	defer C.free(unsafe.Pointer(v))
	ret := uintptr(C.property_bag_set_string(properties.handle, (C.int)(id), nil, v))
	if ret != C.SPX_NOERROR {
		return NewCarbonError(ret)
	}
	return nil
}

// SetPropertyByString sets the value of a property.
func (properties PropertyCollection) SetPropertyByString(name string, value string) error {
	n := C.CString(name)
	defer C.free(unsafe.Pointer(n))
	v := C.CString(value)
	defer C.free(unsafe.Pointer(v))
	ret := uintptr(C.property_bag_set_string(properties.handle, -1, n, v))
	if ret != C.SPX_NOERROR {
		return NewCarbonError(ret)
	}
	return nil
}

// Close disposes the associated resources.
func (properties PropertyCollection) Close() {
	C.property_bag_release(properties.handle)
}

// NewPropertyCollectionFromHandle creates a PropertyCollection from a handle (for internal use)
func NewPropertyCollectionFromHandle(handle SPXHandle) *PropertyCollection {
	propertyCollection := new(PropertyCollection)
	propertyCollection.handle = uintptr2handle(handle)
	return propertyCollection
}
