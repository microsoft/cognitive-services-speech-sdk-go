// Copyright (c) Microsoft. All rights reserved.
// Licensed under the MIT license. See LICENSE.md file in the project root for full license information.

package common

// ServicePropertyChannel defines channels used to pass property settings to service.
type ServicePropertyChannel int

const (
	// URIQueryParameter uses URI query parameter to pass property settings to service.
	URIQueryParameter ServicePropertyChannel = 0
)
