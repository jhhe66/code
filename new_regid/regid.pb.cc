// Generated by the protocol buffer compiler.  DO NOT EDIT!
// source: regid.proto

#define INTERNAL_SUPPRESS_PROTOBUF_FIELD_DEPRECATION
#include "regid.pb.h"

#include <algorithm>

#include <google/protobuf/stubs/common.h>
#include <google/protobuf/stubs/once.h>
#include <google/protobuf/io/coded_stream.h>
#include <google/protobuf/wire_format_lite_inl.h>
#include <google/protobuf/descriptor.h>
#include <google/protobuf/generated_message_reflection.h>
#include <google/protobuf/reflection_ops.h>
#include <google/protobuf/wire_format.h>
// @@protoc_insertion_point(includes)

namespace JPush {

namespace {

const ::google::protobuf::Descriptor* regid_descriptor_ = NULL;
const ::google::protobuf::internal::GeneratedMessageReflection*
  regid_reflection_ = NULL;

}  // namespace


void protobuf_AssignDesc_regid_2eproto() {
  protobuf_AddDesc_regid_2eproto();
  const ::google::protobuf::FileDescriptor* file =
    ::google::protobuf::DescriptorPool::generated_pool()->FindFileByName(
      "regid.proto");
  GOOGLE_CHECK(file != NULL);
  regid_descriptor_ = file->message_type(0);
  static const int regid_offsets_[4] = {
    GOOGLE_PROTOBUF_GENERATED_MESSAGE_FIELD_OFFSET(regid, version_),
    GOOGLE_PROTOBUF_GENERATED_MESSAGE_FIELD_OFFSET(regid, appkey_),
    GOOGLE_PROTOBUF_GENERATED_MESSAGE_FIELD_OFFSET(regid, uid_),
    GOOGLE_PROTOBUF_GENERATED_MESSAGE_FIELD_OFFSET(regid, platform_),
  };
  regid_reflection_ =
    new ::google::protobuf::internal::GeneratedMessageReflection(
      regid_descriptor_,
      regid::default_instance_,
      regid_offsets_,
      GOOGLE_PROTOBUF_GENERATED_MESSAGE_FIELD_OFFSET(regid, _has_bits_[0]),
      GOOGLE_PROTOBUF_GENERATED_MESSAGE_FIELD_OFFSET(regid, _unknown_fields_),
      -1,
      ::google::protobuf::DescriptorPool::generated_pool(),
      ::google::protobuf::MessageFactory::generated_factory(),
      sizeof(regid));
}

namespace {

GOOGLE_PROTOBUF_DECLARE_ONCE(protobuf_AssignDescriptors_once_);
inline void protobuf_AssignDescriptorsOnce() {
  ::google::protobuf::GoogleOnceInit(&protobuf_AssignDescriptors_once_,
                 &protobuf_AssignDesc_regid_2eproto);
}

void protobuf_RegisterTypes(const ::std::string&) {
  protobuf_AssignDescriptorsOnce();
  ::google::protobuf::MessageFactory::InternalRegisterGeneratedMessage(
    regid_descriptor_, &regid::default_instance());
}

}  // namespace

void protobuf_ShutdownFile_regid_2eproto() {
  delete regid::default_instance_;
  delete regid_reflection_;
}

void protobuf_AddDesc_regid_2eproto() {
  static bool already_here = false;
  if (already_here) return;
  already_here = true;
  GOOGLE_PROTOBUF_VERIFY_VERSION;

  ::google::protobuf::DescriptorPool::InternalAddGeneratedFile(
    "\n\013regid.proto\022\005JPush\"G\n\005regid\022\017\n\007version"
    "\030\001 \002(\006\022\016\n\006appkey\030\002 \002(\014\022\013\n\003uid\030\003 \002(\006\022\020\n\010p"
    "latform\030\004 \002(\006", 93);
  ::google::protobuf::MessageFactory::InternalRegisterGeneratedFile(
    "regid.proto", &protobuf_RegisterTypes);
  regid::default_instance_ = new regid();
  regid::default_instance_->InitAsDefaultInstance();
  ::google::protobuf::internal::OnShutdown(&protobuf_ShutdownFile_regid_2eproto);
}

// Force AddDescriptors() to be called at static initialization time.
struct StaticDescriptorInitializer_regid_2eproto {
  StaticDescriptorInitializer_regid_2eproto() {
    protobuf_AddDesc_regid_2eproto();
  }
} static_descriptor_initializer_regid_2eproto_;

// ===================================================================

#ifndef _MSC_VER
const int regid::kVersionFieldNumber;
const int regid::kAppkeyFieldNumber;
const int regid::kUidFieldNumber;
const int regid::kPlatformFieldNumber;
#endif  // !_MSC_VER

regid::regid()
  : ::google::protobuf::Message() {
  SharedCtor();
}

void regid::InitAsDefaultInstance() {
}

regid::regid(const regid& from)
  : ::google::protobuf::Message() {
  SharedCtor();
  MergeFrom(from);
}

void regid::SharedCtor() {
  _cached_size_ = 0;
  version_ = GOOGLE_ULONGLONG(0);
  appkey_ = const_cast< ::std::string*>(&::google::protobuf::internal::kEmptyString);
  uid_ = GOOGLE_ULONGLONG(0);
  platform_ = GOOGLE_ULONGLONG(0);
  ::memset(_has_bits_, 0, sizeof(_has_bits_));
}

regid::~regid() {
  SharedDtor();
}

void regid::SharedDtor() {
  if (appkey_ != &::google::protobuf::internal::kEmptyString) {
    delete appkey_;
  }
  if (this != default_instance_) {
  }
}

void regid::SetCachedSize(int size) const {
  GOOGLE_SAFE_CONCURRENT_WRITES_BEGIN();
  _cached_size_ = size;
  GOOGLE_SAFE_CONCURRENT_WRITES_END();
}
const ::google::protobuf::Descriptor* regid::descriptor() {
  protobuf_AssignDescriptorsOnce();
  return regid_descriptor_;
}

const regid& regid::default_instance() {
  if (default_instance_ == NULL) protobuf_AddDesc_regid_2eproto();
  return *default_instance_;
}

regid* regid::default_instance_ = NULL;

regid* regid::New() const {
  return new regid;
}

void regid::Clear() {
  if (_has_bits_[0 / 32] & (0xffu << (0 % 32))) {
    version_ = GOOGLE_ULONGLONG(0);
    if (has_appkey()) {
      if (appkey_ != &::google::protobuf::internal::kEmptyString) {
        appkey_->clear();
      }
    }
    uid_ = GOOGLE_ULONGLONG(0);
    platform_ = GOOGLE_ULONGLONG(0);
  }
  ::memset(_has_bits_, 0, sizeof(_has_bits_));
  mutable_unknown_fields()->Clear();
}

bool regid::MergePartialFromCodedStream(
    ::google::protobuf::io::CodedInputStream* input) {
#define DO_(EXPRESSION) if (!(EXPRESSION)) return false
  ::google::protobuf::uint32 tag;
  while ((tag = input->ReadTag()) != 0) {
    switch (::google::protobuf::internal::WireFormatLite::GetTagFieldNumber(tag)) {
      // required fixed64 version = 1;
      case 1: {
        if (::google::protobuf::internal::WireFormatLite::GetTagWireType(tag) ==
            ::google::protobuf::internal::WireFormatLite::WIRETYPE_FIXED64) {
          DO_((::google::protobuf::internal::WireFormatLite::ReadPrimitive<
                   ::google::protobuf::uint64, ::google::protobuf::internal::WireFormatLite::TYPE_FIXED64>(
                 input, &version_)));
          set_has_version();
        } else {
          goto handle_uninterpreted;
        }
        if (input->ExpectTag(18)) goto parse_appkey;
        break;
      }

      // required bytes appkey = 2;
      case 2: {
        if (::google::protobuf::internal::WireFormatLite::GetTagWireType(tag) ==
            ::google::protobuf::internal::WireFormatLite::WIRETYPE_LENGTH_DELIMITED) {
         parse_appkey:
          DO_(::google::protobuf::internal::WireFormatLite::ReadBytes(
                input, this->mutable_appkey()));
        } else {
          goto handle_uninterpreted;
        }
        if (input->ExpectTag(25)) goto parse_uid;
        break;
      }

      // required fixed64 uid = 3;
      case 3: {
        if (::google::protobuf::internal::WireFormatLite::GetTagWireType(tag) ==
            ::google::protobuf::internal::WireFormatLite::WIRETYPE_FIXED64) {
         parse_uid:
          DO_((::google::protobuf::internal::WireFormatLite::ReadPrimitive<
                   ::google::protobuf::uint64, ::google::protobuf::internal::WireFormatLite::TYPE_FIXED64>(
                 input, &uid_)));
          set_has_uid();
        } else {
          goto handle_uninterpreted;
        }
        if (input->ExpectTag(33)) goto parse_platform;
        break;
      }

      // required fixed64 platform = 4;
      case 4: {
        if (::google::protobuf::internal::WireFormatLite::GetTagWireType(tag) ==
            ::google::protobuf::internal::WireFormatLite::WIRETYPE_FIXED64) {
         parse_platform:
          DO_((::google::protobuf::internal::WireFormatLite::ReadPrimitive<
                   ::google::protobuf::uint64, ::google::protobuf::internal::WireFormatLite::TYPE_FIXED64>(
                 input, &platform_)));
          set_has_platform();
        } else {
          goto handle_uninterpreted;
        }
        if (input->ExpectAtEnd()) return true;
        break;
      }

      default: {
      handle_uninterpreted:
        if (::google::protobuf::internal::WireFormatLite::GetTagWireType(tag) ==
            ::google::protobuf::internal::WireFormatLite::WIRETYPE_END_GROUP) {
          return true;
        }
        DO_(::google::protobuf::internal::WireFormat::SkipField(
              input, tag, mutable_unknown_fields()));
        break;
      }
    }
  }
  return true;
#undef DO_
}

void regid::SerializeWithCachedSizes(
    ::google::protobuf::io::CodedOutputStream* output) const {
  // required fixed64 version = 1;
  if (has_version()) {
    ::google::protobuf::internal::WireFormatLite::WriteFixed64(1, this->version(), output);
  }

  // required bytes appkey = 2;
  if (has_appkey()) {
    ::google::protobuf::internal::WireFormatLite::WriteBytes(
      2, this->appkey(), output);
  }

  // required fixed64 uid = 3;
  if (has_uid()) {
    ::google::protobuf::internal::WireFormatLite::WriteFixed64(3, this->uid(), output);
  }

  // required fixed64 platform = 4;
  if (has_platform()) {
    ::google::protobuf::internal::WireFormatLite::WriteFixed64(4, this->platform(), output);
  }

  if (!unknown_fields().empty()) {
    ::google::protobuf::internal::WireFormat::SerializeUnknownFields(
        unknown_fields(), output);
  }
}

::google::protobuf::uint8* regid::SerializeWithCachedSizesToArray(
    ::google::protobuf::uint8* target) const {
  // required fixed64 version = 1;
  if (has_version()) {
    target = ::google::protobuf::internal::WireFormatLite::WriteFixed64ToArray(1, this->version(), target);
  }

  // required bytes appkey = 2;
  if (has_appkey()) {
    target =
      ::google::protobuf::internal::WireFormatLite::WriteBytesToArray(
        2, this->appkey(), target);
  }

  // required fixed64 uid = 3;
  if (has_uid()) {
    target = ::google::protobuf::internal::WireFormatLite::WriteFixed64ToArray(3, this->uid(), target);
  }

  // required fixed64 platform = 4;
  if (has_platform()) {
    target = ::google::protobuf::internal::WireFormatLite::WriteFixed64ToArray(4, this->platform(), target);
  }

  if (!unknown_fields().empty()) {
    target = ::google::protobuf::internal::WireFormat::SerializeUnknownFieldsToArray(
        unknown_fields(), target);
  }
  return target;
}

int regid::ByteSize() const {
  int total_size = 0;

  if (_has_bits_[0 / 32] & (0xffu << (0 % 32))) {
    // required fixed64 version = 1;
    if (has_version()) {
      total_size += 1 + 8;
    }

    // required bytes appkey = 2;
    if (has_appkey()) {
      total_size += 1 +
        ::google::protobuf::internal::WireFormatLite::BytesSize(
          this->appkey());
    }

    // required fixed64 uid = 3;
    if (has_uid()) {
      total_size += 1 + 8;
    }

    // required fixed64 platform = 4;
    if (has_platform()) {
      total_size += 1 + 8;
    }

  }
  if (!unknown_fields().empty()) {
    total_size +=
      ::google::protobuf::internal::WireFormat::ComputeUnknownFieldsSize(
        unknown_fields());
  }
  GOOGLE_SAFE_CONCURRENT_WRITES_BEGIN();
  _cached_size_ = total_size;
  GOOGLE_SAFE_CONCURRENT_WRITES_END();
  return total_size;
}

void regid::MergeFrom(const ::google::protobuf::Message& from) {
  GOOGLE_CHECK_NE(&from, this);
  const regid* source =
    ::google::protobuf::internal::dynamic_cast_if_available<const regid*>(
      &from);
  if (source == NULL) {
    ::google::protobuf::internal::ReflectionOps::Merge(from, this);
  } else {
    MergeFrom(*source);
  }
}

void regid::MergeFrom(const regid& from) {
  GOOGLE_CHECK_NE(&from, this);
  if (from._has_bits_[0 / 32] & (0xffu << (0 % 32))) {
    if (from.has_version()) {
      set_version(from.version());
    }
    if (from.has_appkey()) {
      set_appkey(from.appkey());
    }
    if (from.has_uid()) {
      set_uid(from.uid());
    }
    if (from.has_platform()) {
      set_platform(from.platform());
    }
  }
  mutable_unknown_fields()->MergeFrom(from.unknown_fields());
}

void regid::CopyFrom(const ::google::protobuf::Message& from) {
  if (&from == this) return;
  Clear();
  MergeFrom(from);
}

void regid::CopyFrom(const regid& from) {
  if (&from == this) return;
  Clear();
  MergeFrom(from);
}

bool regid::IsInitialized() const {
  if ((_has_bits_[0] & 0x0000000f) != 0x0000000f) return false;

  return true;
}

void regid::Swap(regid* other) {
  if (other != this) {
    std::swap(version_, other->version_);
    std::swap(appkey_, other->appkey_);
    std::swap(uid_, other->uid_);
    std::swap(platform_, other->platform_);
    std::swap(_has_bits_[0], other->_has_bits_[0]);
    _unknown_fields_.Swap(&other->_unknown_fields_);
    std::swap(_cached_size_, other->_cached_size_);
  }
}

::google::protobuf::Metadata regid::GetMetadata() const {
  protobuf_AssignDescriptorsOnce();
  ::google::protobuf::Metadata metadata;
  metadata.descriptor = regid_descriptor_;
  metadata.reflection = regid_reflection_;
  return metadata;
}


// @@protoc_insertion_point(namespace_scope)

}  // namespace JPush

// @@protoc_insertion_point(global_scope)