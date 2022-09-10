# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: google/cloud/apigeeregistry/v1/apihub/extensions.proto
"""Generated protocol buffer code."""
from google.protobuf import descriptor as _descriptor
from google.protobuf import descriptor_pool as _descriptor_pool
from google.protobuf import message as _message
from google.protobuf import reflection as _reflection
from google.protobuf import symbol_database as _symbol_database
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from google.api import field_behavior_pb2 as google_dot_api_dot_field__behavior__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n6google/cloud/apigeeregistry/v1/apihub/extensions.proto\x12%google.cloud.apigeeregistry.v1.apihub\x1a\x1fgoogle/api/field_behavior.proto\"\xbc\x02\n\x14\x41piSpecExtensionList\x12\n\n\x02id\x18\x01 \x01(\t\x12\x0c\n\x04kind\x18\x02 \x01(\t\x12\x14\n\x0c\x64isplay_name\x18\x03 \x01(\t\x12\x13\n\x0b\x64\x65scription\x18\x04 \x01(\t\x12`\n\nextensions\x18\x05 \x03(\x0b\x32L.google.cloud.apigeeregistry.v1.apihub.ApiSpecExtensionList.ApiSpecExtension\x1a}\n\x10\x41piSpecExtension\x12\x0f\n\x02id\x18\x01 \x01(\tB\x03\xe0\x41\x02\x12\x19\n\x0c\x64isplay_name\x18\x02 \x01(\tB\x03\xe0\x41\x02\x12\x13\n\x0b\x64\x65scription\x18\x03 \x01(\t\x12\x0e\n\x06\x66ilter\x18\x04 \x01(\t\x12\x18\n\x0buri_pattern\x18\x05 \x01(\tB\x03\xe0\x41\x02\x42\x62\n)com.google.cloud.apigeeregistry.v1.apihubB\x0f\x45xtensionsProtoP\x01Z\"github.com/apigee/registry/rpc;rpcb\x06proto3')



_APISPECEXTENSIONLIST = DESCRIPTOR.message_types_by_name['ApiSpecExtensionList']
_APISPECEXTENSIONLIST_APISPECEXTENSION = _APISPECEXTENSIONLIST.nested_types_by_name['ApiSpecExtension']
ApiSpecExtensionList = _reflection.GeneratedProtocolMessageType('ApiSpecExtensionList', (_message.Message,), {

  'ApiSpecExtension' : _reflection.GeneratedProtocolMessageType('ApiSpecExtension', (_message.Message,), {
    'DESCRIPTOR' : _APISPECEXTENSIONLIST_APISPECEXTENSION,
    '__module__' : 'google.cloud.apigeeregistry.v1.apihub.extensions_pb2'
    # @@protoc_insertion_point(class_scope:google.cloud.apigeeregistry.v1.apihub.ApiSpecExtensionList.ApiSpecExtension)
    })
  ,
  'DESCRIPTOR' : _APISPECEXTENSIONLIST,
  '__module__' : 'google.cloud.apigeeregistry.v1.apihub.extensions_pb2'
  # @@protoc_insertion_point(class_scope:google.cloud.apigeeregistry.v1.apihub.ApiSpecExtensionList)
  })
_sym_db.RegisterMessage(ApiSpecExtensionList)
_sym_db.RegisterMessage(ApiSpecExtensionList.ApiSpecExtension)

if _descriptor._USE_C_DESCRIPTORS == False:

  DESCRIPTOR._options = None
  DESCRIPTOR._serialized_options = b'\n)com.google.cloud.apigeeregistry.v1.apihubB\017ExtensionsProtoP\001Z\"github.com/apigee/registry/rpc;rpc'
  _APISPECEXTENSIONLIST_APISPECEXTENSION.fields_by_name['id']._options = None
  _APISPECEXTENSIONLIST_APISPECEXTENSION.fields_by_name['id']._serialized_options = b'\340A\002'
  _APISPECEXTENSIONLIST_APISPECEXTENSION.fields_by_name['display_name']._options = None
  _APISPECEXTENSIONLIST_APISPECEXTENSION.fields_by_name['display_name']._serialized_options = b'\340A\002'
  _APISPECEXTENSIONLIST_APISPECEXTENSION.fields_by_name['uri_pattern']._options = None
  _APISPECEXTENSIONLIST_APISPECEXTENSION.fields_by_name['uri_pattern']._serialized_options = b'\340A\002'
  _APISPECEXTENSIONLIST._serialized_start=131
  _APISPECEXTENSIONLIST._serialized_end=447
  _APISPECEXTENSIONLIST_APISPECEXTENSION._serialized_start=322
  _APISPECEXTENSIONLIST_APISPECEXTENSION._serialized_end=447
# @@protoc_insertion_point(module_scope)
