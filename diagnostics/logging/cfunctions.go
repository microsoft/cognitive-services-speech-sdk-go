// Copyright (c) Microsoft. All rights reserved.
// Licensed under the MIT license. See LICENSE.md file in the project root for full license information.

package logging

// #include <stdlib.h>
// #include <speechapi_c_diagnostics.h>
//
// extern void goEventLoggerCallback(const char* logLine);
//
// static void cgo_event_logger_callback(const char* logLine)
// {
//     goEventLoggerCallback(logLine);
// }
//
// uintptr_t cgo_register_event_callback()
// {
//     return (uintptr_t)diagnostics_logmessage_set_callback(cgo_event_logger_callback);
// }
//
// uintptr_t cgo_unregister_event_callback()
// {
//     return (uintptr_t)diagnostics_logmessage_set_callback((DIAGNOSTICS_CALLBACK_FUNC)0);
// }
import "C"
