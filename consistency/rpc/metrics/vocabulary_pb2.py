# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: metrics/vocabulary.proto
"""Generated protocol buffer code."""
from google.protobuf import descriptor as _descriptor
from google.protobuf import descriptor_pool as _descriptor_pool
from google.protobuf import message as _message
from google.protobuf import reflection as _reflection
from google.protobuf import symbol_database as _symbol_database
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()




DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\x18metrics/vocabulary.proto\x12\x12gnostic.metrics.v1\"(\n\tWordCount\x12\x0c\n\x04word\x18\x01 \x01(\t\x12\r\n\x05\x63ount\x18\x02 \x01(\x05\"\xe3\x01\n\nVocabulary\x12\x0c\n\x04name\x18\x01 \x01(\t\x12.\n\x07schemas\x18\x02 \x03(\x0b\x32\x1d.gnostic.metrics.v1.WordCount\x12\x31\n\nproperties\x18\x03 \x03(\x0b\x32\x1d.gnostic.metrics.v1.WordCount\x12\x31\n\noperations\x18\x04 \x03(\x0b\x32\x1d.gnostic.metrics.v1.WordCount\x12\x31\n\nparameters\x18\x05 \x03(\x0b\x32\x1d.gnostic.metrics.v1.WordCount\"F\n\x0eVocabularyList\x12\x34\n\x0cvocabularies\x18\x01 \x03(\x0b\x32\x1e.gnostic.metrics.v1.Vocabulary\"\xb5\x01\n\x07Version\x12\x0c\n\x04name\x18\x01 \x01(\t\x12\x16\n\x0enew_term_count\x18\x02 \x01(\x05\x12\x31\n\tnew_terms\x18\x03 \x01(\x0b\x32\x1e.gnostic.metrics.v1.Vocabulary\x12\x1a\n\x12\x64\x65leted_term_count\x18\x04 \x01(\x05\x12\x35\n\rdeleted_terms\x18\x05 \x01(\x0b\x32\x1e.gnostic.metrics.v1.Vocabulary\"M\n\x0eVersionHistory\x12\x0c\n\x04name\x18\x01 \x01(\t\x12-\n\x08versions\x18\x02 \x03(\x0b\x32\x1b.gnostic.metrics.v1.VersionB\x1eZ\x1c./metrics;gnostic_metrics_v1b\x06proto3')



_WORDCOUNT = DESCRIPTOR.message_types_by_name['WordCount']
_VOCABULARY = DESCRIPTOR.message_types_by_name['Vocabulary']
_VOCABULARYLIST = DESCRIPTOR.message_types_by_name['VocabularyList']
_VERSION = DESCRIPTOR.message_types_by_name['Version']
_VERSIONHISTORY = DESCRIPTOR.message_types_by_name['VersionHistory']
WordCount = _reflection.GeneratedProtocolMessageType('WordCount', (_message.Message,), {
  'DESCRIPTOR' : _WORDCOUNT,
  '__module__' : 'metrics.vocabulary_pb2'
  # @@protoc_insertion_point(class_scope:gnostic.metrics.v1.WordCount)
  })
_sym_db.RegisterMessage(WordCount)

Vocabulary = _reflection.GeneratedProtocolMessageType('Vocabulary', (_message.Message,), {
  'DESCRIPTOR' : _VOCABULARY,
  '__module__' : 'metrics.vocabulary_pb2'
  # @@protoc_insertion_point(class_scope:gnostic.metrics.v1.Vocabulary)
  })
_sym_db.RegisterMessage(Vocabulary)

VocabularyList = _reflection.GeneratedProtocolMessageType('VocabularyList', (_message.Message,), {
  'DESCRIPTOR' : _VOCABULARYLIST,
  '__module__' : 'metrics.vocabulary_pb2'
  # @@protoc_insertion_point(class_scope:gnostic.metrics.v1.VocabularyList)
  })
_sym_db.RegisterMessage(VocabularyList)

Version = _reflection.GeneratedProtocolMessageType('Version', (_message.Message,), {
  'DESCRIPTOR' : _VERSION,
  '__module__' : 'metrics.vocabulary_pb2'
  # @@protoc_insertion_point(class_scope:gnostic.metrics.v1.Version)
  })
_sym_db.RegisterMessage(Version)

VersionHistory = _reflection.GeneratedProtocolMessageType('VersionHistory', (_message.Message,), {
  'DESCRIPTOR' : _VERSIONHISTORY,
  '__module__' : 'metrics.vocabulary_pb2'
  # @@protoc_insertion_point(class_scope:gnostic.metrics.v1.VersionHistory)
  })
_sym_db.RegisterMessage(VersionHistory)

if _descriptor._USE_C_DESCRIPTORS == False:

  DESCRIPTOR._options = None
  DESCRIPTOR._serialized_options = b'Z\034./metrics;gnostic_metrics_v1'
  _WORDCOUNT._serialized_start=48
  _WORDCOUNT._serialized_end=88
  _VOCABULARY._serialized_start=91
  _VOCABULARY._serialized_end=318
  _VOCABULARYLIST._serialized_start=320
  _VOCABULARYLIST._serialized_end=390
  _VERSION._serialized_start=393
  _VERSION._serialized_end=574
  _VERSIONHISTORY._serialized_start=576
  _VERSIONHISTORY._serialized_end=653
# @@protoc_insertion_point(module_scope)
