// Generated by the protocol buffer compiler.  DO NOT EDIT!
// source: tagalias_common.proto

#define INTERNAL_SUPPRESS_PROTOBUF_FIELD_DEPRECATION
#include "tagalias_common.pb.h"

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

namespace {

const ::google::protobuf::Descriptor* UserStuff_descriptor_ = NULL;
const ::google::protobuf::internal::GeneratedMessageReflection*
  UserStuff_reflection_ = NULL;
const ::google::protobuf::Descriptor* MiBaseInfo_descriptor_ = NULL;
const ::google::protobuf::internal::GeneratedMessageReflection*
  MiBaseInfo_reflection_ = NULL;
const ::google::protobuf::EnumDescriptor* ACTION_descriptor_ = NULL;
const ::google::protobuf::EnumDescriptor* QUERY_ACTION_descriptor_ = NULL;
const ::google::protobuf::EnumDescriptor* PLATFORM_descriptor_ = NULL;
const ::google::protobuf::EnumDescriptor* SERVER_TYPE_descriptor_ = NULL;
const ::google::protobuf::EnumDescriptor* STATUS_descriptor_ = NULL;
const ::google::protobuf::EnumDescriptor* DEVSUBTYPE_descriptor_ = NULL;

}  // namespace


void protobuf_AssignDesc_tagalias_5fcommon_2eproto() {
  protobuf_AddDesc_tagalias_5fcommon_2eproto();
  const ::google::protobuf::FileDescriptor* file =
    ::google::protobuf::DescriptorPool::generated_pool()->FindFileByName(
      "tagalias_common.proto");
  GOOGLE_CHECK(file != NULL);
  UserStuff_descriptor_ = file->message_type(0);
  static const int UserStuff_offsets_[3] = {
    GOOGLE_PROTOBUF_GENERATED_MESSAGE_FIELD_OFFSET(UserStuff, uid_),
    GOOGLE_PROTOBUF_GENERATED_MESSAGE_FIELD_OFFSET(UserStuff, type_),
    GOOGLE_PROTOBUF_GENERATED_MESSAGE_FIELD_OFFSET(UserStuff, token_),
  };
  UserStuff_reflection_ =
    new ::google::protobuf::internal::GeneratedMessageReflection(
      UserStuff_descriptor_,
      UserStuff::default_instance_,
      UserStuff_offsets_,
      GOOGLE_PROTOBUF_GENERATED_MESSAGE_FIELD_OFFSET(UserStuff, _has_bits_[0]),
      GOOGLE_PROTOBUF_GENERATED_MESSAGE_FIELD_OFFSET(UserStuff, _unknown_fields_),
      -1,
      ::google::protobuf::DescriptorPool::generated_pool(),
      ::google::protobuf::MessageFactory::generated_factory(),
      sizeof(UserStuff));
  MiBaseInfo_descriptor_ = file->message_type(1);
  static const int MiBaseInfo_offsets_[2] = {
    GOOGLE_PROTOBUF_GENERATED_MESSAGE_FIELD_OFFSET(MiBaseInfo, appid_),
    GOOGLE_PROTOBUF_GENERATED_MESSAGE_FIELD_OFFSET(MiBaseInfo, appsecret_),
  };
  MiBaseInfo_reflection_ =
    new ::google::protobuf::internal::GeneratedMessageReflection(
      MiBaseInfo_descriptor_,
      MiBaseInfo::default_instance_,
      MiBaseInfo_offsets_,
      GOOGLE_PROTOBUF_GENERATED_MESSAGE_FIELD_OFFSET(MiBaseInfo, _has_bits_[0]),
      GOOGLE_PROTOBUF_GENERATED_MESSAGE_FIELD_OFFSET(MiBaseInfo, _unknown_fields_),
      -1,
      ::google::protobuf::DescriptorPool::generated_pool(),
      ::google::protobuf::MessageFactory::generated_factory(),
      sizeof(MiBaseInfo));
  ACTION_descriptor_ = file->enum_type(0);
  QUERY_ACTION_descriptor_ = file->enum_type(1);
  PLATFORM_descriptor_ = file->enum_type(2);
  SERVER_TYPE_descriptor_ = file->enum_type(3);
  STATUS_descriptor_ = file->enum_type(4);
  DEVSUBTYPE_descriptor_ = file->enum_type(5);
}

namespace {

GOOGLE_PROTOBUF_DECLARE_ONCE(protobuf_AssignDescriptors_once_);
inline void protobuf_AssignDescriptorsOnce() {
  ::google::protobuf::GoogleOnceInit(&protobuf_AssignDescriptors_once_,
                 &protobuf_AssignDesc_tagalias_5fcommon_2eproto);
}

void protobuf_RegisterTypes(const ::std::string&) {
  protobuf_AssignDescriptorsOnce();
  ::google::protobuf::MessageFactory::InternalRegisterGeneratedMessage(
    UserStuff_descriptor_, &UserStuff::default_instance());
  ::google::protobuf::MessageFactory::InternalRegisterGeneratedMessage(
    MiBaseInfo_descriptor_, &MiBaseInfo::default_instance());
}

}  // namespace

void protobuf_ShutdownFile_tagalias_5fcommon_2eproto() {
  delete UserStuff::default_instance_;
  delete UserStuff_reflection_;
  delete MiBaseInfo::default_instance_;
  delete MiBaseInfo_reflection_;
}

void protobuf_AddDesc_tagalias_5fcommon_2eproto() {
  static bool already_here = false;
  if (already_here) return;
  already_here = true;
  GOOGLE_PROTOBUF_VERIFY_VERSION;

  ::google::protobuf::DescriptorPool::InternalAddGeneratedFile(
    "\n\025tagalias_common.proto\"O\n\tUserStuff\022\013\n\003"
    "uid\030\001 \001(\004\022&\n\004type\030\002 \001(\0162\013.DEVSUBTYPE:\013NO"
    "RMAL_TYPE\022\r\n\005token\030\003 \001(\014\".\n\nMiBaseInfo\022\r"
    "\n\005appid\030\001 \001(\014\022\021\n\tappsecret\030\002 \001(\014*\264\001\n\006ACT"
    "ION\022\023\n\017ADD_USER_TO_TAG\020\001\022\025\n\021DEL_USER_FRO"
    "M_TAG\020\002\022\020\n\014USER_SET_TAG\020\004\022\022\n\016USER_CLEAR_"
    "TAG\020\010\022\025\n\021DEL_TAGS_FROM_APP\020\020\022\022\n\016USER_SET"
    "_ALIAS\020 \022\024\n\020USER_CLEAR_ALIAS\020@\022\027\n\022DEL_AL"
    "IAS_FROM_APP\020\200\001*\310\002\n\014QUERY_ACTION\022\024\n\020GET_"
    "TAGS_BY_USER\020\001\022\025\n\021GET_ALIAS_BY_USER\020\002\022\024\n"
    "\020GET_USERS_BY_TAG\020\004\022\026\n\022GET_USERS_BY_ALIA"
    "S\020\010\022\023\n\017GET_APPKEY_TAGS\020\020\022\032\n\026CHECK_USER_B"
    "ELONG_TAGS\020 \022\033\n\027CHECK_USER_BELONG_ALIAS\020"
    "@\022\034\n\027GET_USERS_COUNT_BY_TAGS\020\200\001\022\035\n\030GET_U"
    "SERS_COUNT_BY_ALIAS\020\200\002\022\031\n\024CHECK_TAGS_HAS"
    "_USERS\020\200\004\022\032\n\025CHECK_ALIAS_HAS_USERS\020\200\010\022\033\n"
    "\026GET_TAGS_COUNT_BY_USER\020\200\020*1\n\010PLATFORM\022\013"
    "\n\007ANDROID\020\000\022\n\n\006IPHONE\020\001\022\014\n\010WINPHONE\020\002*0\n"
    "\013SERVER_TYPE\022\021\n\rINTERFACE_API\020\001\022\016\n\nDEVIC"
    "E_API\020\002*\257\003\n\006STATUS\022\013\n\007SUCCESS\020\000\022\r\n\tBAD_P"
    "ROTO\020\001\022\r\n\tNO_APPKEY\020\002\022\017\n\013NO_PLATFORM\020\003\022\023"
    "\n\017NO_QUERY_ACTION\020\004\022\031\n\025SERIALIZE_RESP_FA"
    "ILED\020\005\022\031\n\025SERVER_INTERNAL_ERROR\020\006\022\025\n\021EXC"
    "EED_TOTAL_PAGE\020\007\022\021\n\rTOO_MANY_TAGS\020\010\022\022\n\016T"
    "OO_MANY_ALIAS\020\t\022\n\n\006NO_UID\020\n\022\013\n\007NO_TAGS\020\013"
    "\022\014\n\010NO_ALIAS\020\014\022\024\n\020AUTHORITY_FAILED\020\r\022\020\n\014"
    "INVALID_TAGS\020\016\022\021\n\rINVALID_ALIAS\020\017\022\020\n\014INV"
    "ALID_PAGE\020\020\022\024\n\020STATUS_REDISDOWN\020@\022\023\n\017TOO"
    "_MANY_QUERYS\020A\022\027\n\023QUERY_TOO_FREQUENCY\020B\022"
    "\021\n\rOTHER_API_ERR\020C\022\025\n\020STATUS_OTHER_ERR\020\200"
    "\001*O\n\nDEVSUBTYPE\022\017\n\013NORMAL_TYPE\020\001\022\017\n\013XIAO"
    "MI_TYPE\020\002\022\017\n\013HUAWEI_TYPE\020\004\022\016\n\nMEIZU_TYPE"
    "\020\010", 1282);
  ::google::protobuf::MessageFactory::InternalRegisterGeneratedFile(
    "tagalias_common.proto", &protobuf_RegisterTypes);
  UserStuff::default_instance_ = new UserStuff();
  MiBaseInfo::default_instance_ = new MiBaseInfo();
  UserStuff::default_instance_->InitAsDefaultInstance();
  MiBaseInfo::default_instance_->InitAsDefaultInstance();
  ::google::protobuf::internal::OnShutdown(&protobuf_ShutdownFile_tagalias_5fcommon_2eproto);
}

// Force AddDescriptors() to be called at static initialization time.
struct StaticDescriptorInitializer_tagalias_5fcommon_2eproto {
  StaticDescriptorInitializer_tagalias_5fcommon_2eproto() {
    protobuf_AddDesc_tagalias_5fcommon_2eproto();
  }
} static_descriptor_initializer_tagalias_5fcommon_2eproto_;
const ::google::protobuf::EnumDescriptor* ACTION_descriptor() {
  protobuf_AssignDescriptorsOnce();
  return ACTION_descriptor_;
}
bool ACTION_IsValid(int value) {
  switch(value) {
    case 1:
    case 2:
    case 4:
    case 8:
    case 16:
    case 32:
    case 64:
    case 128:
      return true;
    default:
      return false;
  }
}

const ::google::protobuf::EnumDescriptor* QUERY_ACTION_descriptor() {
  protobuf_AssignDescriptorsOnce();
  return QUERY_ACTION_descriptor_;
}
bool QUERY_ACTION_IsValid(int value) {
  switch(value) {
    case 1:
    case 2:
    case 4:
    case 8:
    case 16:
    case 32:
    case 64:
    case 128:
    case 256:
    case 512:
    case 1024:
    case 2048:
      return true;
    default:
      return false;
  }
}

const ::google::protobuf::EnumDescriptor* PLATFORM_descriptor() {
  protobuf_AssignDescriptorsOnce();
  return PLATFORM_descriptor_;
}
bool PLATFORM_IsValid(int value) {
  switch(value) {
    case 0:
    case 1:
    case 2:
      return true;
    default:
      return false;
  }
}

const ::google::protobuf::EnumDescriptor* SERVER_TYPE_descriptor() {
  protobuf_AssignDescriptorsOnce();
  return SERVER_TYPE_descriptor_;
}
bool SERVER_TYPE_IsValid(int value) {
  switch(value) {
    case 1:
    case 2:
      return true;
    default:
      return false;
  }
}

const ::google::protobuf::EnumDescriptor* STATUS_descriptor() {
  protobuf_AssignDescriptorsOnce();
  return STATUS_descriptor_;
}
bool STATUS_IsValid(int value) {
  switch(value) {
    case 0:
    case 1:
    case 2:
    case 3:
    case 4:
    case 5:
    case 6:
    case 7:
    case 8:
    case 9:
    case 10:
    case 11:
    case 12:
    case 13:
    case 14:
    case 15:
    case 16:
    case 64:
    case 65:
    case 66:
    case 67:
    case 128:
      return true;
    default:
      return false;
  }
}

const ::google::protobuf::EnumDescriptor* DEVSUBTYPE_descriptor() {
  protobuf_AssignDescriptorsOnce();
  return DEVSUBTYPE_descriptor_;
}
bool DEVSUBTYPE_IsValid(int value) {
  switch(value) {
    case 1:
    case 2:
    case 4:
    case 8:
      return true;
    default:
      return false;
  }
}


// ===================================================================

#ifndef _MSC_VER
const int UserStuff::kUidFieldNumber;
const int UserStuff::kTypeFieldNumber;
const int UserStuff::kTokenFieldNumber;
#endif  // !_MSC_VER

UserStuff::UserStuff()
  : ::google::protobuf::Message() {
  SharedCtor();
}

void UserStuff::InitAsDefaultInstance() {
}

UserStuff::UserStuff(const UserStuff& from)
  : ::google::protobuf::Message() {
  SharedCtor();
  MergeFrom(from);
}

void UserStuff::SharedCtor() {
  _cached_size_ = 0;
  uid_ = GOOGLE_ULONGLONG(0);
  type_ = 1;
  token_ = const_cast< ::std::string*>(&::google::protobuf::internal::kEmptyString);
  ::memset(_has_bits_, 0, sizeof(_has_bits_));
}

UserStuff::~UserStuff() {
  SharedDtor();
}

void UserStuff::SharedDtor() {
  if (token_ != &::google::protobuf::internal::kEmptyString) {
    delete token_;
  }
  if (this != default_instance_) {
  }
}

void UserStuff::SetCachedSize(int size) const {
  GOOGLE_SAFE_CONCURRENT_WRITES_BEGIN();
  _cached_size_ = size;
  GOOGLE_SAFE_CONCURRENT_WRITES_END();
}
const ::google::protobuf::Descriptor* UserStuff::descriptor() {
  protobuf_AssignDescriptorsOnce();
  return UserStuff_descriptor_;
}

const UserStuff& UserStuff::default_instance() {
  if (default_instance_ == NULL) protobuf_AddDesc_tagalias_5fcommon_2eproto();
  return *default_instance_;
}

UserStuff* UserStuff::default_instance_ = NULL;

UserStuff* UserStuff::New() const {
  return new UserStuff;
}

void UserStuff::Clear() {
  if (_has_bits_[0 / 32] & (0xffu << (0 % 32))) {
    uid_ = GOOGLE_ULONGLONG(0);
    type_ = 1;
    if (has_token()) {
      if (token_ != &::google::protobuf::internal::kEmptyString) {
        token_->clear();
      }
    }
  }
  ::memset(_has_bits_, 0, sizeof(_has_bits_));
  mutable_unknown_fields()->Clear();
}

bool UserStuff::MergePartialFromCodedStream(
    ::google::protobuf::io::CodedInputStream* input) {
#define DO_(EXPRESSION) if (!(EXPRESSION)) return false
  ::google::protobuf::uint32 tag;
  while ((tag = input->ReadTag()) != 0) {
    switch (::google::protobuf::internal::WireFormatLite::GetTagFieldNumber(tag)) {
      // optional uint64 uid = 1;
      case 1: {
        if (::google::protobuf::internal::WireFormatLite::GetTagWireType(tag) ==
            ::google::protobuf::internal::WireFormatLite::WIRETYPE_VARINT) {
          DO_((::google::protobuf::internal::WireFormatLite::ReadPrimitive<
                   ::google::protobuf::uint64, ::google::protobuf::internal::WireFormatLite::TYPE_UINT64>(
                 input, &uid_)));
          set_has_uid();
        } else {
          goto handle_uninterpreted;
        }
        if (input->ExpectTag(16)) goto parse_type;
        break;
      }

      // optional .DEVSUBTYPE type = 2 [default = NORMAL_TYPE];
      case 2: {
        if (::google::protobuf::internal::WireFormatLite::GetTagWireType(tag) ==
            ::google::protobuf::internal::WireFormatLite::WIRETYPE_VARINT) {
         parse_type:
          int value;
          DO_((::google::protobuf::internal::WireFormatLite::ReadPrimitive<
                   int, ::google::protobuf::internal::WireFormatLite::TYPE_ENUM>(
                 input, &value)));
          if (::DEVSUBTYPE_IsValid(value)) {
            set_type(static_cast< ::DEVSUBTYPE >(value));
          } else {
            mutable_unknown_fields()->AddVarint(2, value);
          }
        } else {
          goto handle_uninterpreted;
        }
        if (input->ExpectTag(26)) goto parse_token;
        break;
      }

      // optional bytes token = 3;
      case 3: {
        if (::google::protobuf::internal::WireFormatLite::GetTagWireType(tag) ==
            ::google::protobuf::internal::WireFormatLite::WIRETYPE_LENGTH_DELIMITED) {
         parse_token:
          DO_(::google::protobuf::internal::WireFormatLite::ReadBytes(
                input, this->mutable_token()));
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

void UserStuff::SerializeWithCachedSizes(
    ::google::protobuf::io::CodedOutputStream* output) const {
  // optional uint64 uid = 1;
  if (has_uid()) {
    ::google::protobuf::internal::WireFormatLite::WriteUInt64(1, this->uid(), output);
  }

  // optional .DEVSUBTYPE type = 2 [default = NORMAL_TYPE];
  if (has_type()) {
    ::google::protobuf::internal::WireFormatLite::WriteEnum(
      2, this->type(), output);
  }

  // optional bytes token = 3;
  if (has_token()) {
    ::google::protobuf::internal::WireFormatLite::WriteBytes(
      3, this->token(), output);
  }

  if (!unknown_fields().empty()) {
    ::google::protobuf::internal::WireFormat::SerializeUnknownFields(
        unknown_fields(), output);
  }
}

::google::protobuf::uint8* UserStuff::SerializeWithCachedSizesToArray(
    ::google::protobuf::uint8* target) const {
  // optional uint64 uid = 1;
  if (has_uid()) {
    target = ::google::protobuf::internal::WireFormatLite::WriteUInt64ToArray(1, this->uid(), target);
  }

  // optional .DEVSUBTYPE type = 2 [default = NORMAL_TYPE];
  if (has_type()) {
    target = ::google::protobuf::internal::WireFormatLite::WriteEnumToArray(
      2, this->type(), target);
  }

  // optional bytes token = 3;
  if (has_token()) {
    target =
      ::google::protobuf::internal::WireFormatLite::WriteBytesToArray(
        3, this->token(), target);
  }

  if (!unknown_fields().empty()) {
    target = ::google::protobuf::internal::WireFormat::SerializeUnknownFieldsToArray(
        unknown_fields(), target);
  }
  return target;
}

int UserStuff::ByteSize() const {
  int total_size = 0;

  if (_has_bits_[0 / 32] & (0xffu << (0 % 32))) {
    // optional uint64 uid = 1;
    if (has_uid()) {
      total_size += 1 +
        ::google::protobuf::internal::WireFormatLite::UInt64Size(
          this->uid());
    }

    // optional .DEVSUBTYPE type = 2 [default = NORMAL_TYPE];
    if (has_type()) {
      total_size += 1 +
        ::google::protobuf::internal::WireFormatLite::EnumSize(this->type());
    }

    // optional bytes token = 3;
    if (has_token()) {
      total_size += 1 +
        ::google::protobuf::internal::WireFormatLite::BytesSize(
          this->token());
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

void UserStuff::MergeFrom(const ::google::protobuf::Message& from) {
  GOOGLE_CHECK_NE(&from, this);
  const UserStuff* source =
    ::google::protobuf::internal::dynamic_cast_if_available<const UserStuff*>(
      &from);
  if (source == NULL) {
    ::google::protobuf::internal::ReflectionOps::Merge(from, this);
  } else {
    MergeFrom(*source);
  }
}

void UserStuff::MergeFrom(const UserStuff& from) {
  GOOGLE_CHECK_NE(&from, this);
  if (from._has_bits_[0 / 32] & (0xffu << (0 % 32))) {
    if (from.has_uid()) {
      set_uid(from.uid());
    }
    if (from.has_type()) {
      set_type(from.type());
    }
    if (from.has_token()) {
      set_token(from.token());
    }
  }
  mutable_unknown_fields()->MergeFrom(from.unknown_fields());
}

void UserStuff::CopyFrom(const ::google::protobuf::Message& from) {
  if (&from == this) return;
  Clear();
  MergeFrom(from);
}

void UserStuff::CopyFrom(const UserStuff& from) {
  if (&from == this) return;
  Clear();
  MergeFrom(from);
}

bool UserStuff::IsInitialized() const {

  return true;
}

void UserStuff::Swap(UserStuff* other) {
  if (other != this) {
    std::swap(uid_, other->uid_);
    std::swap(type_, other->type_);
    std::swap(token_, other->token_);
    std::swap(_has_bits_[0], other->_has_bits_[0]);
    _unknown_fields_.Swap(&other->_unknown_fields_);
    std::swap(_cached_size_, other->_cached_size_);
  }
}

::google::protobuf::Metadata UserStuff::GetMetadata() const {
  protobuf_AssignDescriptorsOnce();
  ::google::protobuf::Metadata metadata;
  metadata.descriptor = UserStuff_descriptor_;
  metadata.reflection = UserStuff_reflection_;
  return metadata;
}


// ===================================================================

#ifndef _MSC_VER
const int MiBaseInfo::kAppidFieldNumber;
const int MiBaseInfo::kAppsecretFieldNumber;
#endif  // !_MSC_VER

MiBaseInfo::MiBaseInfo()
  : ::google::protobuf::Message() {
  SharedCtor();
}

void MiBaseInfo::InitAsDefaultInstance() {
}

MiBaseInfo::MiBaseInfo(const MiBaseInfo& from)
  : ::google::protobuf::Message() {
  SharedCtor();
  MergeFrom(from);
}

void MiBaseInfo::SharedCtor() {
  _cached_size_ = 0;
  appid_ = const_cast< ::std::string*>(&::google::protobuf::internal::kEmptyString);
  appsecret_ = const_cast< ::std::string*>(&::google::protobuf::internal::kEmptyString);
  ::memset(_has_bits_, 0, sizeof(_has_bits_));
}

MiBaseInfo::~MiBaseInfo() {
  SharedDtor();
}

void MiBaseInfo::SharedDtor() {
  if (appid_ != &::google::protobuf::internal::kEmptyString) {
    delete appid_;
  }
  if (appsecret_ != &::google::protobuf::internal::kEmptyString) {
    delete appsecret_;
  }
  if (this != default_instance_) {
  }
}

void MiBaseInfo::SetCachedSize(int size) const {
  GOOGLE_SAFE_CONCURRENT_WRITES_BEGIN();
  _cached_size_ = size;
  GOOGLE_SAFE_CONCURRENT_WRITES_END();
}
const ::google::protobuf::Descriptor* MiBaseInfo::descriptor() {
  protobuf_AssignDescriptorsOnce();
  return MiBaseInfo_descriptor_;
}

const MiBaseInfo& MiBaseInfo::default_instance() {
  if (default_instance_ == NULL) protobuf_AddDesc_tagalias_5fcommon_2eproto();
  return *default_instance_;
}

MiBaseInfo* MiBaseInfo::default_instance_ = NULL;

MiBaseInfo* MiBaseInfo::New() const {
  return new MiBaseInfo;
}

void MiBaseInfo::Clear() {
  if (_has_bits_[0 / 32] & (0xffu << (0 % 32))) {
    if (has_appid()) {
      if (appid_ != &::google::protobuf::internal::kEmptyString) {
        appid_->clear();
      }
    }
    if (has_appsecret()) {
      if (appsecret_ != &::google::protobuf::internal::kEmptyString) {
        appsecret_->clear();
      }
    }
  }
  ::memset(_has_bits_, 0, sizeof(_has_bits_));
  mutable_unknown_fields()->Clear();
}

bool MiBaseInfo::MergePartialFromCodedStream(
    ::google::protobuf::io::CodedInputStream* input) {
#define DO_(EXPRESSION) if (!(EXPRESSION)) return false
  ::google::protobuf::uint32 tag;
  while ((tag = input->ReadTag()) != 0) {
    switch (::google::protobuf::internal::WireFormatLite::GetTagFieldNumber(tag)) {
      // optional bytes appid = 1;
      case 1: {
        if (::google::protobuf::internal::WireFormatLite::GetTagWireType(tag) ==
            ::google::protobuf::internal::WireFormatLite::WIRETYPE_LENGTH_DELIMITED) {
          DO_(::google::protobuf::internal::WireFormatLite::ReadBytes(
                input, this->mutable_appid()));
        } else {
          goto handle_uninterpreted;
        }
        if (input->ExpectTag(18)) goto parse_appsecret;
        break;
      }

      // optional bytes appsecret = 2;
      case 2: {
        if (::google::protobuf::internal::WireFormatLite::GetTagWireType(tag) ==
            ::google::protobuf::internal::WireFormatLite::WIRETYPE_LENGTH_DELIMITED) {
         parse_appsecret:
          DO_(::google::protobuf::internal::WireFormatLite::ReadBytes(
                input, this->mutable_appsecret()));
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

void MiBaseInfo::SerializeWithCachedSizes(
    ::google::protobuf::io::CodedOutputStream* output) const {
  // optional bytes appid = 1;
  if (has_appid()) {
    ::google::protobuf::internal::WireFormatLite::WriteBytes(
      1, this->appid(), output);
  }

  // optional bytes appsecret = 2;
  if (has_appsecret()) {
    ::google::protobuf::internal::WireFormatLite::WriteBytes(
      2, this->appsecret(), output);
  }

  if (!unknown_fields().empty()) {
    ::google::protobuf::internal::WireFormat::SerializeUnknownFields(
        unknown_fields(), output);
  }
}

::google::protobuf::uint8* MiBaseInfo::SerializeWithCachedSizesToArray(
    ::google::protobuf::uint8* target) const {
  // optional bytes appid = 1;
  if (has_appid()) {
    target =
      ::google::protobuf::internal::WireFormatLite::WriteBytesToArray(
        1, this->appid(), target);
  }

  // optional bytes appsecret = 2;
  if (has_appsecret()) {
    target =
      ::google::protobuf::internal::WireFormatLite::WriteBytesToArray(
        2, this->appsecret(), target);
  }

  if (!unknown_fields().empty()) {
    target = ::google::protobuf::internal::WireFormat::SerializeUnknownFieldsToArray(
        unknown_fields(), target);
  }
  return target;
}

int MiBaseInfo::ByteSize() const {
  int total_size = 0;

  if (_has_bits_[0 / 32] & (0xffu << (0 % 32))) {
    // optional bytes appid = 1;
    if (has_appid()) {
      total_size += 1 +
        ::google::protobuf::internal::WireFormatLite::BytesSize(
          this->appid());
    }

    // optional bytes appsecret = 2;
    if (has_appsecret()) {
      total_size += 1 +
        ::google::protobuf::internal::WireFormatLite::BytesSize(
          this->appsecret());
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

void MiBaseInfo::MergeFrom(const ::google::protobuf::Message& from) {
  GOOGLE_CHECK_NE(&from, this);
  const MiBaseInfo* source =
    ::google::protobuf::internal::dynamic_cast_if_available<const MiBaseInfo*>(
      &from);
  if (source == NULL) {
    ::google::protobuf::internal::ReflectionOps::Merge(from, this);
  } else {
    MergeFrom(*source);
  }
}

void MiBaseInfo::MergeFrom(const MiBaseInfo& from) {
  GOOGLE_CHECK_NE(&from, this);
  if (from._has_bits_[0 / 32] & (0xffu << (0 % 32))) {
    if (from.has_appid()) {
      set_appid(from.appid());
    }
    if (from.has_appsecret()) {
      set_appsecret(from.appsecret());
    }
  }
  mutable_unknown_fields()->MergeFrom(from.unknown_fields());
}

void MiBaseInfo::CopyFrom(const ::google::protobuf::Message& from) {
  if (&from == this) return;
  Clear();
  MergeFrom(from);
}

void MiBaseInfo::CopyFrom(const MiBaseInfo& from) {
  if (&from == this) return;
  Clear();
  MergeFrom(from);
}

bool MiBaseInfo::IsInitialized() const {

  return true;
}

void MiBaseInfo::Swap(MiBaseInfo* other) {
  if (other != this) {
    std::swap(appid_, other->appid_);
    std::swap(appsecret_, other->appsecret_);
    std::swap(_has_bits_[0], other->_has_bits_[0]);
    _unknown_fields_.Swap(&other->_unknown_fields_);
    std::swap(_cached_size_, other->_cached_size_);
  }
}

::google::protobuf::Metadata MiBaseInfo::GetMetadata() const {
  protobuf_AssignDescriptorsOnce();
  ::google::protobuf::Metadata metadata;
  metadata.descriptor = MiBaseInfo_descriptor_;
  metadata.reflection = MiBaseInfo_reflection_;
  return metadata;
}


// @@protoc_insertion_point(namespace_scope)

// @@protoc_insertion_point(global_scope)
