//
// Copyright (c) Microsoft. All rights reserved.
// Licensed under the MIT license. See LICENSE.md file in the project root for full license information.
//

#pragma once
#include <speechapi_c_common.h>
#include <speechapi_c.h>

namespace Microsoft {
namespace CognitiveServices {
namespace Speech {
namespace Impl {

// Callback function type for binding events
typedef void (*PBINDING_CALLBACK_FUNC)(SPXEVENTHANDLE hevent, SPXEVENTHANDLE hresult, void* pvContext); //take handle to hevent and hresponse, instantiate as null

// Function to register a callback for HTTP requests
SPXAPI register_send_callback(PBINDING_CALLBACK_FUNC pCallback, void* pvContext);

// Function for C# to call when streaming data is available
SPXAPI process_streaming_data(SPXEVENTHANDLE hevent, const uint8_t* data, size_t size);

// Function to get the property bag from event args
SPXAPI http_eventargs_get_property_bag(SPXEVENTHANDLE hevent, SPXPROPERTYBAGHANDLE* hpropbag);

}
}
}
} // Microsoft::CognitiveServices::Speech::Impl

