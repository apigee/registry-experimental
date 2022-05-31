# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: google/cloud/apigeeregistry/applications/v1alpha1/registry_styleguide.proto
"""Generated protocol buffer code."""
from google.protobuf import descriptor as _descriptor
from google.protobuf import descriptor_pool as _descriptor_pool
from google.protobuf import message as _message
from google.protobuf import reflection as _reflection
from google.protobuf import symbol_database as _symbol_database
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from google.api import field_behavior_pb2 as google_dot_api_dot_field__behavior__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\nKgoogle/cloud/apigeeregistry/applications/v1alpha1/registry_styleguide.proto\x12\x31google.cloud.apigeeregistry.applications.v1alpha1\x1a\x1fgoogle/api/field_behavior.proto\"\xef\x01\n\nStyleGuide\x12\x0f\n\x02id\x18\x01 \x01(\tB\x03\xe0\x41\x02\x12\x19\n\x0c\x64isplay_name\x18\x02 \x01(\tB\x03\xe0\x41\x02\x12\x17\n\nmime_types\x18\x03 \x03(\tB\x03\xe0\x41\x02\x12P\n\nguidelines\x18\x04 \x03(\x0b\x32<.google.cloud.apigeeregistry.applications.v1alpha1.Guideline\x12J\n\x07linters\x18\x05 \x03(\x0b\x32\x39.google.cloud.apigeeregistry.applications.v1alpha1.Linter\"\xc3\x02\n\tGuideline\x12\x0f\n\x02id\x18\x01 \x01(\tB\x03\xe0\x41\x02\x12\x19\n\x0c\x64isplay_name\x18\x02 \x01(\tB\x03\xe0\x41\x02\x12\x13\n\x0b\x64\x65scription\x18\x03 \x01(\t\x12\x46\n\x05rules\x18\x04 \x03(\x0b\x32\x37.google.cloud.apigeeregistry.applications.v1alpha1.Rule\x12S\n\x06status\x18\x05 \x01(\x0e\x32\x43.google.cloud.apigeeregistry.applications.v1alpha1.Guideline.Status\"X\n\x06Status\x12\x16\n\x12STATUS_UNSPECIFIED\x10\x00\x12\x0c\n\x08PROPOSED\x10\x01\x12\n\n\x06\x41\x43TIVE\x10\x02\x12\x0e\n\nDEPRECATED\x10\x03\x12\x0c\n\x08\x44ISABLED\x10\x04\"\xac\x02\n\x04Rule\x12\x0f\n\x02id\x18\x01 \x01(\tB\x03\xe0\x41\x02\x12\x14\n\x0c\x64isplay_name\x18\x02 \x01(\t\x12\x13\n\x0b\x64\x65scription\x18\x03 \x01(\t\x12\x13\n\x06linter\x18\x04 \x01(\tB\x03\xe0\x41\x02\x12\x1c\n\x0flinter_rulename\x18\x05 \x01(\tB\x03\xe0\x41\x02\x12R\n\x08severity\x18\x06 \x01(\x0e\x32@.google.cloud.apigeeregistry.applications.v1alpha1.Rule.Severity\x12\x0f\n\x07\x64oc_uri\x18\x07 \x01(\t\"P\n\x08Severity\x12\x18\n\x14SEVERITY_UNSPECIFIED\x10\x00\x12\t\n\x05\x45RROR\x10\x01\x12\x0b\n\x07WARNING\x10\x02\x12\x08\n\x04INFO\x10\x03\x12\x08\n\x04HINT\x10\x04\"-\n\x06Linter\x12\x11\n\x04name\x18\x01 \x01(\tB\x03\xe0\x41\x02\x12\x10\n\x03uri\x18\x02 \x01(\tB\x03\xe0\x41\x02\x42v\n5com.google.cloud.apigeeregistry.applications.v1alpha1B\x17RegistryStyleGuideProtoP\x01Z\"github.com/apigee/registry/rpc;rpcb\x06proto3')



_STYLEGUIDE = DESCRIPTOR.message_types_by_name['StyleGuide']
_GUIDELINE = DESCRIPTOR.message_types_by_name['Guideline']
_RULE = DESCRIPTOR.message_types_by_name['Rule']
_LINTER = DESCRIPTOR.message_types_by_name['Linter']
_GUIDELINE_STATUS = _GUIDELINE.enum_types_by_name['Status']
_RULE_SEVERITY = _RULE.enum_types_by_name['Severity']
StyleGuide = _reflection.GeneratedProtocolMessageType('StyleGuide', (_message.Message,), {
  'DESCRIPTOR' : _STYLEGUIDE,
  '__module__' : 'google.cloud.apigeeregistry.applications.v1alpha1.registry_styleguide_pb2'
  # @@protoc_insertion_point(class_scope:google.cloud.apigeeregistry.applications.v1alpha1.StyleGuide)
  })
_sym_db.RegisterMessage(StyleGuide)

Guideline = _reflection.GeneratedProtocolMessageType('Guideline', (_message.Message,), {
  'DESCRIPTOR' : _GUIDELINE,
  '__module__' : 'google.cloud.apigeeregistry.applications.v1alpha1.registry_styleguide_pb2'
  # @@protoc_insertion_point(class_scope:google.cloud.apigeeregistry.applications.v1alpha1.Guideline)
  })
_sym_db.RegisterMessage(Guideline)

Rule = _reflection.GeneratedProtocolMessageType('Rule', (_message.Message,), {
  'DESCRIPTOR' : _RULE,
  '__module__' : 'google.cloud.apigeeregistry.applications.v1alpha1.registry_styleguide_pb2'
  # @@protoc_insertion_point(class_scope:google.cloud.apigeeregistry.applications.v1alpha1.Rule)
  })
_sym_db.RegisterMessage(Rule)

Linter = _reflection.GeneratedProtocolMessageType('Linter', (_message.Message,), {
  'DESCRIPTOR' : _LINTER,
  '__module__' : 'google.cloud.apigeeregistry.applications.v1alpha1.registry_styleguide_pb2'
  # @@protoc_insertion_point(class_scope:google.cloud.apigeeregistry.applications.v1alpha1.Linter)
  })
_sym_db.RegisterMessage(Linter)

if _descriptor._USE_C_DESCRIPTORS == False:

  DESCRIPTOR._options = None
  DESCRIPTOR._serialized_options = b'\n5com.google.cloud.apigeeregistry.applications.v1alpha1B\027RegistryStyleGuideProtoP\001Z\"github.com/apigee/registry/rpc;rpc'
  _STYLEGUIDE.fields_by_name['id']._options = None
  _STYLEGUIDE.fields_by_name['id']._serialized_options = b'\340A\002'
  _STYLEGUIDE.fields_by_name['display_name']._options = None
  _STYLEGUIDE.fields_by_name['display_name']._serialized_options = b'\340A\002'
  _STYLEGUIDE.fields_by_name['mime_types']._options = None
  _STYLEGUIDE.fields_by_name['mime_types']._serialized_options = b'\340A\002'
  _GUIDELINE.fields_by_name['id']._options = None
  _GUIDELINE.fields_by_name['id']._serialized_options = b'\340A\002'
  _GUIDELINE.fields_by_name['display_name']._options = None
  _GUIDELINE.fields_by_name['display_name']._serialized_options = b'\340A\002'
  _RULE.fields_by_name['id']._options = None
  _RULE.fields_by_name['id']._serialized_options = b'\340A\002'
  _RULE.fields_by_name['linter']._options = None
  _RULE.fields_by_name['linter']._serialized_options = b'\340A\002'
  _RULE.fields_by_name['linter_rulename']._options = None
  _RULE.fields_by_name['linter_rulename']._serialized_options = b'\340A\002'
  _LINTER.fields_by_name['name']._options = None
  _LINTER.fields_by_name['name']._serialized_options = b'\340A\002'
  _LINTER.fields_by_name['uri']._options = None
  _LINTER.fields_by_name['uri']._serialized_options = b'\340A\002'
  _STYLEGUIDE._serialized_start=164
  _STYLEGUIDE._serialized_end=403
  _GUIDELINE._serialized_start=406
  _GUIDELINE._serialized_end=729
  _GUIDELINE_STATUS._serialized_start=641
  _GUIDELINE_STATUS._serialized_end=729
  _RULE._serialized_start=732
  _RULE._serialized_end=1032
  _RULE_SEVERITY._serialized_start=952
  _RULE_SEVERITY._serialized_end=1032
  _LINTER._serialized_start=1034
  _LINTER._serialized_end=1079
# @@protoc_insertion_point(module_scope)