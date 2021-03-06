# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: google/cloud/apigeeregistry/v1/admin_models.proto
"""Generated protocol buffer code."""
from google.protobuf import descriptor as _descriptor
from google.protobuf import descriptor_pool as _descriptor_pool
from google.protobuf import message as _message
from google.protobuf import reflection as _reflection
from google.protobuf import symbol_database as _symbol_database
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from google.api import field_behavior_pb2 as google_dot_api_dot_field__behavior__pb2
from google.api import resource_pb2 as google_dot_api_dot_resource__pb2
from google.protobuf import timestamp_pb2 as google_dot_protobuf_dot_timestamp__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n1google/cloud/apigeeregistry/v1/admin_models.proto\x12\x1egoogle.cloud.apigeeregistry.v1\x1a\x1fgoogle/api/field_behavior.proto\x1a\x19google/api/resource.proto\x1a\x1fgoogle/protobuf/timestamp.proto\"\xae\x03\n\tBuildInfo\x12\x12\n\ngo_version\x18\x01 \x01(\t\x12\x0c\n\x04path\x18\x02 \x01(\t\x12>\n\x04main\x18\x03 \x01(\x0b\x32\x30.google.cloud.apigeeregistry.v1.BuildInfo.Module\x12\x46\n\x0c\x64\x65pendencies\x18\x04 \x03(\x0b\x32\x30.google.cloud.apigeeregistry.v1.BuildInfo.Module\x12I\n\x08settings\x18\x05 \x03(\x0b\x32\x37.google.cloud.apigeeregistry.v1.BuildInfo.SettingsEntry\x1a{\n\x06Module\x12\x0c\n\x04path\x18\x01 \x01(\t\x12\x0f\n\x07version\x18\x02 \x01(\t\x12\x0b\n\x03sum\x18\x03 \x01(\t\x12\x45\n\x0breplacement\x18\x04 \x01(\x0b\x32\x30.google.cloud.apigeeregistry.v1.BuildInfo.Module\x1a/\n\rSettingsEntry\x12\x0b\n\x03key\x18\x01 \x01(\t\x12\r\n\x05value\x18\x02 \x01(\t:\x02\x38\x01\"S\n\x06Status\x12\x0f\n\x07message\x18\x01 \x01(\t\x12\x38\n\x05\x62uild\x18\x02 \x01(\x0b\x32).google.cloud.apigeeregistry.v1.BuildInfo\"\x92\x01\n\x07Storage\x12\x13\n\x0b\x64\x65scription\x18\x01 \x01(\t\x12G\n\x0b\x63ollections\x18\x02 \x03(\x0b\x32\x32.google.cloud.apigeeregistry.v1.Storage.Collection\x1a)\n\nCollection\x12\x0c\n\x04name\x18\x01 \x01(\t\x12\r\n\x05\x63ount\x18\x02 \x01(\x03\"\xee\x01\n\x07Project\x12\x0c\n\x04name\x18\x01 \x01(\t\x12\x14\n\x0c\x64isplay_name\x18\x02 \x01(\t\x12\x13\n\x0b\x64\x65scription\x18\x03 \x01(\t\x12\x34\n\x0b\x63reate_time\x18\x04 \x01(\x0b\x32\x1a.google.protobuf.TimestampB\x03\xe0\x41\x03\x12\x34\n\x0bupdate_time\x18\x05 \x01(\x0b\x32\x1a.google.protobuf.TimestampB\x03\xe0\x41\x03:>\xea\x41;\n%apigeeregistry.googleapis.com/Project\x12\x12projects/{project}B\\\n\"com.google.cloud.apigeeregistry.v1B\x10\x41\x64minModelsProtoP\x01Z\"github.com/apigee/registry/rpc;rpcb\x06proto3')



_BUILDINFO = DESCRIPTOR.message_types_by_name['BuildInfo']
_BUILDINFO_MODULE = _BUILDINFO.nested_types_by_name['Module']
_BUILDINFO_SETTINGSENTRY = _BUILDINFO.nested_types_by_name['SettingsEntry']
_STATUS = DESCRIPTOR.message_types_by_name['Status']
_STORAGE = DESCRIPTOR.message_types_by_name['Storage']
_STORAGE_COLLECTION = _STORAGE.nested_types_by_name['Collection']
_PROJECT = DESCRIPTOR.message_types_by_name['Project']
BuildInfo = _reflection.GeneratedProtocolMessageType('BuildInfo', (_message.Message,), {

  'Module' : _reflection.GeneratedProtocolMessageType('Module', (_message.Message,), {
    'DESCRIPTOR' : _BUILDINFO_MODULE,
    '__module__' : 'google.cloud.apigeeregistry.v1.admin_models_pb2'
    # @@protoc_insertion_point(class_scope:google.cloud.apigeeregistry.v1.BuildInfo.Module)
    })
  ,

  'SettingsEntry' : _reflection.GeneratedProtocolMessageType('SettingsEntry', (_message.Message,), {
    'DESCRIPTOR' : _BUILDINFO_SETTINGSENTRY,
    '__module__' : 'google.cloud.apigeeregistry.v1.admin_models_pb2'
    # @@protoc_insertion_point(class_scope:google.cloud.apigeeregistry.v1.BuildInfo.SettingsEntry)
    })
  ,
  'DESCRIPTOR' : _BUILDINFO,
  '__module__' : 'google.cloud.apigeeregistry.v1.admin_models_pb2'
  # @@protoc_insertion_point(class_scope:google.cloud.apigeeregistry.v1.BuildInfo)
  })
_sym_db.RegisterMessage(BuildInfo)
_sym_db.RegisterMessage(BuildInfo.Module)
_sym_db.RegisterMessage(BuildInfo.SettingsEntry)

Status = _reflection.GeneratedProtocolMessageType('Status', (_message.Message,), {
  'DESCRIPTOR' : _STATUS,
  '__module__' : 'google.cloud.apigeeregistry.v1.admin_models_pb2'
  # @@protoc_insertion_point(class_scope:google.cloud.apigeeregistry.v1.Status)
  })
_sym_db.RegisterMessage(Status)

Storage = _reflection.GeneratedProtocolMessageType('Storage', (_message.Message,), {

  'Collection' : _reflection.GeneratedProtocolMessageType('Collection', (_message.Message,), {
    'DESCRIPTOR' : _STORAGE_COLLECTION,
    '__module__' : 'google.cloud.apigeeregistry.v1.admin_models_pb2'
    # @@protoc_insertion_point(class_scope:google.cloud.apigeeregistry.v1.Storage.Collection)
    })
  ,
  'DESCRIPTOR' : _STORAGE,
  '__module__' : 'google.cloud.apigeeregistry.v1.admin_models_pb2'
  # @@protoc_insertion_point(class_scope:google.cloud.apigeeregistry.v1.Storage)
  })
_sym_db.RegisterMessage(Storage)
_sym_db.RegisterMessage(Storage.Collection)

Project = _reflection.GeneratedProtocolMessageType('Project', (_message.Message,), {
  'DESCRIPTOR' : _PROJECT,
  '__module__' : 'google.cloud.apigeeregistry.v1.admin_models_pb2'
  # @@protoc_insertion_point(class_scope:google.cloud.apigeeregistry.v1.Project)
  })
_sym_db.RegisterMessage(Project)

if _descriptor._USE_C_DESCRIPTORS == False:

  DESCRIPTOR._options = None
  DESCRIPTOR._serialized_options = b'\n\"com.google.cloud.apigeeregistry.v1B\020AdminModelsProtoP\001Z\"github.com/apigee/registry/rpc;rpc'
  _BUILDINFO_SETTINGSENTRY._options = None
  _BUILDINFO_SETTINGSENTRY._serialized_options = b'8\001'
  _PROJECT.fields_by_name['create_time']._options = None
  _PROJECT.fields_by_name['create_time']._serialized_options = b'\340A\003'
  _PROJECT.fields_by_name['update_time']._options = None
  _PROJECT.fields_by_name['update_time']._serialized_options = b'\340A\003'
  _PROJECT._options = None
  _PROJECT._serialized_options = b'\352A;\n%apigeeregistry.googleapis.com/Project\022\022projects/{project}'
  _BUILDINFO._serialized_start=179
  _BUILDINFO._serialized_end=609
  _BUILDINFO_MODULE._serialized_start=437
  _BUILDINFO_MODULE._serialized_end=560
  _BUILDINFO_SETTINGSENTRY._serialized_start=562
  _BUILDINFO_SETTINGSENTRY._serialized_end=609
  _STATUS._serialized_start=611
  _STATUS._serialized_end=694
  _STORAGE._serialized_start=697
  _STORAGE._serialized_end=843
  _STORAGE_COLLECTION._serialized_start=802
  _STORAGE_COLLECTION._serialized_end=843
  _PROJECT._serialized_start=846
  _PROJECT._serialized_end=1084
# @@protoc_insertion_point(module_scope)
