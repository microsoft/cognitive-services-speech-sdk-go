package audio
// This file defines the proxy functions required to use callbacks

// #include <stdlib.h>
// #include <stdint.h>
// #include <speechapi_c_common.h>
// extern int cgoAudioCallReadCallback(SPXHANDLE handle, uint8_t *buffer, uint32_t size);
// extern void cgoAudioCallGetPropertyCallback(SPXHANDLE handle, int id, uint8_t *value, uint32_t size);
// extern void cgoAudioCallCloseCallback(SPXHANDLE handle);
//
// int cgo_audio_read_callback_wrapper(void *context, uint8_t *buffer, uint32_t size)
// {
//     return cgoAudioCallReadCallback((SPXHANDLE)context, buffer, size);
// }
//
// void cgo_audio_get_property_callback_wrapper(void* context, int id, uint8_t* value, uint32_t size)
// {
//     cgoAudioCallGetPropertyCallback((SPXHANDLE)context, id, value, size);
// }
//
// void cgo_audio_close_callback_wrapper(void *context)
// {
//     cgoAudioCallCloseCallback((SPXHANDLE)context);
// }
import "C"