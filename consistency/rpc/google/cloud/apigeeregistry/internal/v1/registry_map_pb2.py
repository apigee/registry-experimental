# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: google/cloud/apigeeregistry/internal/v1/registry_map.proto
"""Generated protocol buffer code."""
from google.protobuf import descriptor as _descriptor
from google.protobuf import descriptor_pool as _descriptor_pool
from google.protobuf import message as _message
from google.protobuf import reflection as _reflection
from google.protobuf import symbol_database as _symbol_database
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()




DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n:google/cloud/apigeeregistry/internal/v1/registry_map.proto\x12\'google.cloud.apigeeregistry.internal.v1\"\x81\x01\n\x03Map\x12J\n\x07\x65ntries\x18\x01 \x03(\x0b\x32\x39.google.cloud.apigeeregistry.internal.v1.Map.EntriesEntry\x1a.\n\x0c\x45ntriesEntry\x12\x0b\n\x03key\x18\x01 \x01(\t\x12\r\n\x05value\x18\x02 \x01(\t:\x02\x38\x01\x42\x65\n+com.google.cloud.apigeeregistry.internal.v1B\x10RegistryMapProtoP\x01Z\"github.com/apigee/registry/rpc;rpcb\x06proto3')



_MAP = DESCRIPTOR.message_types_by_name['Map']
_MAP_ENTRIESENTRY = _MAP.nested_types_by_name['EntriesEntry']
Map = _reflection.GeneratedProtocolMessageType('Map', (_message.Message,), {

  'EntriesEntry' : _reflection.GeneratedProtocolMessageType('EntriesEntry', (_message.Message,), {
    'DESCRIPTOR' : _MAP_ENTRIESENTRY,
    '__module__' : 'google.cloud.apigeeregistry.internal.v1.registry_map_pb2'
    # @@protoc_insertion_point(class_scope:google.cloud.apigeeregistry.internal.v1.Map.EntriesEntry)
    })
  ,
  'DESCRIPTOR' : _MAP,
  '__module__' : 'google.cloud.apigeeregistry.internal.v1.registry_map_pb2'
  # @@protoc_insertion_point(class_scope:google.cloud.apigeeregistry.internal.v1.Map)
  })
_sym_db.RegisterMessage(Map)
_sym_db.RegisterMessage(Map.EntriesEntry)

if _descriptor._USE_C_DESCRIPTORS == False:

  DESCRIPTOR._options = None
  DESCRIPTOR._serialized_options = b'\n+com.google.cloud.apigeeregistry.internal.v1B\020RegistryMapProtoP\001Z\"github.com/apigee/registry/rpc;rpc'
  _MAP_ENTRIESENTRY._options = None
  _MAP_ENTRIESENTRY._serialized_options = b'8\001'
  _MAP._serialized_start=104
  _MAP._serialized_end=233
  _MAP_ENTRIESENTRY._serialized_start=187
  _MAP_ENTRIESENTRY._serialized_end=233
# @@protoc_insertion_point(module_scope)