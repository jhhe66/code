// Generated by the protocol buffer compiler.  DO NOT EDIT!
// source: tagalias_common.proto

#ifndef PROTOBUF_tagalias_5fcommon_2eproto__INCLUDED
#define PROTOBUF_tagalias_5fcommon_2eproto__INCLUDED

#include <string>

#include <google/protobuf/stubs/common.h>

#if GOOGLE_PROTOBUF_VERSION < 2005000
#error This file was generated by a newer version of protoc which is
#error incompatible with your Protocol Buffer headers.  Please update
#error your headers.
#endif
#if 2005000 < GOOGLE_PROTOBUF_MIN_PROTOC_VERSION
#error This file was generated by an older version of protoc which is
#error incompatible with your Protocol Buffer headers.  Please
#error regenerate this file with a newer version of protoc.
#endif

#include <google/protobuf/generated_message_util.h>
#include <google/protobuf/message.h>
#include <google/protobuf/repeated_field.h>
#include <google/protobuf/extension_set.h>
#include <google/protobuf/generated_enum_reflection.h>
#include <google/protobuf/unknown_field_set.h>
// @@protoc_insertion_point(includes)

// Internal implementation detail -- do not call these.
void  protobuf_AddDesc_tagalias_5fcommon_2eproto();
void protobuf_AssignDesc_tagalias_5fcommon_2eproto();
void protobuf_ShutdownFile_tagalias_5fcommon_2eproto();

class UserStuff;
class MiBaseInfo;

enum ACTION {
  ADD_USER_TO_TAG = 1,
  DEL_USER_FROM_TAG = 2,
  USER_SET_TAG = 4,
  USER_CLEAR_TAG = 8,
  DEL_TAGS_FROM_APP = 16,
  USER_SET_ALIAS = 32,
  USER_CLEAR_ALIAS = 64,
  DEL_ALIAS_FROM_APP = 128
};
bool ACTION_IsValid(int value);
const ACTION ACTION_MIN = ADD_USER_TO_TAG;
const ACTION ACTION_MAX = DEL_ALIAS_FROM_APP;
const int ACTION_ARRAYSIZE = ACTION_MAX + 1;

const ::google::protobuf::EnumDescriptor* ACTION_descriptor();
inline const ::std::string& ACTION_Name(ACTION value) {
  return ::google::protobuf::internal::NameOfEnum(
    ACTION_descriptor(), value);
}
inline bool ACTION_Parse(
    const ::std::string& name, ACTION* value) {
  return ::google::protobuf::internal::ParseNamedEnum<ACTION>(
    ACTION_descriptor(), name, value);
}
enum QUERY_ACTION {
  GET_TAGS_BY_USER = 1,
  GET_ALIAS_BY_USER = 2,
  GET_USERS_BY_TAG = 4,
  GET_USERS_BY_ALIAS = 8,
  GET_APPKEY_TAGS = 16,
  CHECK_USER_BELONG_TAGS = 32,
  CHECK_USER_BELONG_ALIAS = 64,
  GET_USERS_COUNT_BY_TAGS = 128,
  GET_USERS_COUNT_BY_ALIAS = 256,
  CHECK_TAGS_HAS_USERS = 512,
  CHECK_ALIAS_HAS_USERS = 1024,
  GET_TAGS_COUNT_BY_USER = 2048
};
bool QUERY_ACTION_IsValid(int value);
const QUERY_ACTION QUERY_ACTION_MIN = GET_TAGS_BY_USER;
const QUERY_ACTION QUERY_ACTION_MAX = GET_TAGS_COUNT_BY_USER;
const int QUERY_ACTION_ARRAYSIZE = QUERY_ACTION_MAX + 1;

const ::google::protobuf::EnumDescriptor* QUERY_ACTION_descriptor();
inline const ::std::string& QUERY_ACTION_Name(QUERY_ACTION value) {
  return ::google::protobuf::internal::NameOfEnum(
    QUERY_ACTION_descriptor(), value);
}
inline bool QUERY_ACTION_Parse(
    const ::std::string& name, QUERY_ACTION* value) {
  return ::google::protobuf::internal::ParseNamedEnum<QUERY_ACTION>(
    QUERY_ACTION_descriptor(), name, value);
}
enum PLATFORM {
  ANDROID = 0,
  IPHONE = 1,
  WINPHONE = 2
};
bool PLATFORM_IsValid(int value);
const PLATFORM PLATFORM_MIN = ANDROID;
const PLATFORM PLATFORM_MAX = WINPHONE;
const int PLATFORM_ARRAYSIZE = PLATFORM_MAX + 1;

const ::google::protobuf::EnumDescriptor* PLATFORM_descriptor();
inline const ::std::string& PLATFORM_Name(PLATFORM value) {
  return ::google::protobuf::internal::NameOfEnum(
    PLATFORM_descriptor(), value);
}
inline bool PLATFORM_Parse(
    const ::std::string& name, PLATFORM* value) {
  return ::google::protobuf::internal::ParseNamedEnum<PLATFORM>(
    PLATFORM_descriptor(), name, value);
}
enum SERVER_TYPE {
  INTERFACE_API = 1,
  DEVICE_API = 2
};
bool SERVER_TYPE_IsValid(int value);
const SERVER_TYPE SERVER_TYPE_MIN = INTERFACE_API;
const SERVER_TYPE SERVER_TYPE_MAX = DEVICE_API;
const int SERVER_TYPE_ARRAYSIZE = SERVER_TYPE_MAX + 1;

const ::google::protobuf::EnumDescriptor* SERVER_TYPE_descriptor();
inline const ::std::string& SERVER_TYPE_Name(SERVER_TYPE value) {
  return ::google::protobuf::internal::NameOfEnum(
    SERVER_TYPE_descriptor(), value);
}
inline bool SERVER_TYPE_Parse(
    const ::std::string& name, SERVER_TYPE* value) {
  return ::google::protobuf::internal::ParseNamedEnum<SERVER_TYPE>(
    SERVER_TYPE_descriptor(), name, value);
}
enum STATUS {
  SUCCESS = 0,
  BAD_PROTO = 1,
  NO_APPKEY = 2,
  NO_PLATFORM = 3,
  NO_QUERY_ACTION = 4,
  SERIALIZE_RESP_FAILED = 5,
  SERVER_INTERNAL_ERROR = 6,
  EXCEED_TOTAL_PAGE = 7,
  TOO_MANY_TAGS = 8,
  TOO_MANY_ALIAS = 9,
  NO_UID = 10,
  NO_TAGS = 11,
  NO_ALIAS = 12,
  AUTHORITY_FAILED = 13,
  INVALID_TAGS = 14,
  INVALID_ALIAS = 15,
  INVALID_PAGE = 16,
  STATUS_REDISDOWN = 64,
  TOO_MANY_QUERYS = 65,
  QUERY_TOO_FREQUENCY = 66,
  OTHER_API_ERR = 67,
  STATUS_OTHER_ERR = 128
};
bool STATUS_IsValid(int value);
const STATUS STATUS_MIN = SUCCESS;
const STATUS STATUS_MAX = STATUS_OTHER_ERR;
const int STATUS_ARRAYSIZE = STATUS_MAX + 1;

const ::google::protobuf::EnumDescriptor* STATUS_descriptor();
inline const ::std::string& STATUS_Name(STATUS value) {
  return ::google::protobuf::internal::NameOfEnum(
    STATUS_descriptor(), value);
}
inline bool STATUS_Parse(
    const ::std::string& name, STATUS* value) {
  return ::google::protobuf::internal::ParseNamedEnum<STATUS>(
    STATUS_descriptor(), name, value);
}
enum DEVSUBTYPE {
  NORMAL_TYPE = 1,
  XIAOMI_TYPE = 2,
  HUAWEI_TYPE = 4,
  MEIZU_TYPE = 8
};
bool DEVSUBTYPE_IsValid(int value);
const DEVSUBTYPE DEVSUBTYPE_MIN = NORMAL_TYPE;
const DEVSUBTYPE DEVSUBTYPE_MAX = MEIZU_TYPE;
const int DEVSUBTYPE_ARRAYSIZE = DEVSUBTYPE_MAX + 1;

const ::google::protobuf::EnumDescriptor* DEVSUBTYPE_descriptor();
inline const ::std::string& DEVSUBTYPE_Name(DEVSUBTYPE value) {
  return ::google::protobuf::internal::NameOfEnum(
    DEVSUBTYPE_descriptor(), value);
}
inline bool DEVSUBTYPE_Parse(
    const ::std::string& name, DEVSUBTYPE* value) {
  return ::google::protobuf::internal::ParseNamedEnum<DEVSUBTYPE>(
    DEVSUBTYPE_descriptor(), name, value);
}
// ===================================================================

class UserStuff : public ::google::protobuf::Message {
 public:
  UserStuff();
  virtual ~UserStuff();

  UserStuff(const UserStuff& from);

  inline UserStuff& operator=(const UserStuff& from) {
    CopyFrom(from);
    return *this;
  }

  inline const ::google::protobuf::UnknownFieldSet& unknown_fields() const {
    return _unknown_fields_;
  }

  inline ::google::protobuf::UnknownFieldSet* mutable_unknown_fields() {
    return &_unknown_fields_;
  }

  static const ::google::protobuf::Descriptor* descriptor();
  static const UserStuff& default_instance();

  void Swap(UserStuff* other);

  // implements Message ----------------------------------------------

  UserStuff* New() const;
  void CopyFrom(const ::google::protobuf::Message& from);
  void MergeFrom(const ::google::protobuf::Message& from);
  void CopyFrom(const UserStuff& from);
  void MergeFrom(const UserStuff& from);
  void Clear();
  bool IsInitialized() const;

  int ByteSize() const;
  bool MergePartialFromCodedStream(
      ::google::protobuf::io::CodedInputStream* input);
  void SerializeWithCachedSizes(
      ::google::protobuf::io::CodedOutputStream* output) const;
  ::google::protobuf::uint8* SerializeWithCachedSizesToArray(::google::protobuf::uint8* output) const;
  int GetCachedSize() const { return _cached_size_; }
  private:
  void SharedCtor();
  void SharedDtor();
  void SetCachedSize(int size) const;
  public:

  ::google::protobuf::Metadata GetMetadata() const;

  // nested types ----------------------------------------------------

  // accessors -------------------------------------------------------

  // optional uint64 uid = 1;
  inline bool has_uid() const;
  inline void clear_uid();
  static const int kUidFieldNumber = 1;
  inline ::google::protobuf::uint64 uid() const;
  inline void set_uid(::google::protobuf::uint64 value);

  // optional .DEVSUBTYPE type = 2 [default = NORMAL_TYPE];
  inline bool has_type() const;
  inline void clear_type();
  static const int kTypeFieldNumber = 2;
  inline ::DEVSUBTYPE type() const;
  inline void set_type(::DEVSUBTYPE value);

  // optional bytes token = 3;
  inline bool has_token() const;
  inline void clear_token();
  static const int kTokenFieldNumber = 3;
  inline const ::std::string& token() const;
  inline void set_token(const ::std::string& value);
  inline void set_token(const char* value);
  inline void set_token(const void* value, size_t size);
  inline ::std::string* mutable_token();
  inline ::std::string* release_token();
  inline void set_allocated_token(::std::string* token);

  // @@protoc_insertion_point(class_scope:UserStuff)
 private:
  inline void set_has_uid();
  inline void clear_has_uid();
  inline void set_has_type();
  inline void clear_has_type();
  inline void set_has_token();
  inline void clear_has_token();

  ::google::protobuf::UnknownFieldSet _unknown_fields_;

  ::google::protobuf::uint64 uid_;
  ::std::string* token_;
  int type_;

  mutable int _cached_size_;
  ::google::protobuf::uint32 _has_bits_[(3 + 31) / 32];

  friend void  protobuf_AddDesc_tagalias_5fcommon_2eproto();
  friend void protobuf_AssignDesc_tagalias_5fcommon_2eproto();
  friend void protobuf_ShutdownFile_tagalias_5fcommon_2eproto();

  void InitAsDefaultInstance();
  static UserStuff* default_instance_;
};
// -------------------------------------------------------------------

class MiBaseInfo : public ::google::protobuf::Message {
 public:
  MiBaseInfo();
  virtual ~MiBaseInfo();

  MiBaseInfo(const MiBaseInfo& from);

  inline MiBaseInfo& operator=(const MiBaseInfo& from) {
    CopyFrom(from);
    return *this;
  }

  inline const ::google::protobuf::UnknownFieldSet& unknown_fields() const {
    return _unknown_fields_;
  }

  inline ::google::protobuf::UnknownFieldSet* mutable_unknown_fields() {
    return &_unknown_fields_;
  }

  static const ::google::protobuf::Descriptor* descriptor();
  static const MiBaseInfo& default_instance();

  void Swap(MiBaseInfo* other);

  // implements Message ----------------------------------------------

  MiBaseInfo* New() const;
  void CopyFrom(const ::google::protobuf::Message& from);
  void MergeFrom(const ::google::protobuf::Message& from);
  void CopyFrom(const MiBaseInfo& from);
  void MergeFrom(const MiBaseInfo& from);
  void Clear();
  bool IsInitialized() const;

  int ByteSize() const;
  bool MergePartialFromCodedStream(
      ::google::protobuf::io::CodedInputStream* input);
  void SerializeWithCachedSizes(
      ::google::protobuf::io::CodedOutputStream* output) const;
  ::google::protobuf::uint8* SerializeWithCachedSizesToArray(::google::protobuf::uint8* output) const;
  int GetCachedSize() const { return _cached_size_; }
  private:
  void SharedCtor();
  void SharedDtor();
  void SetCachedSize(int size) const;
  public:

  ::google::protobuf::Metadata GetMetadata() const;

  // nested types ----------------------------------------------------

  // accessors -------------------------------------------------------

  // optional bytes appid = 1;
  inline bool has_appid() const;
  inline void clear_appid();
  static const int kAppidFieldNumber = 1;
  inline const ::std::string& appid() const;
  inline void set_appid(const ::std::string& value);
  inline void set_appid(const char* value);
  inline void set_appid(const void* value, size_t size);
  inline ::std::string* mutable_appid();
  inline ::std::string* release_appid();
  inline void set_allocated_appid(::std::string* appid);

  // optional bytes appsecret = 2;
  inline bool has_appsecret() const;
  inline void clear_appsecret();
  static const int kAppsecretFieldNumber = 2;
  inline const ::std::string& appsecret() const;
  inline void set_appsecret(const ::std::string& value);
  inline void set_appsecret(const char* value);
  inline void set_appsecret(const void* value, size_t size);
  inline ::std::string* mutable_appsecret();
  inline ::std::string* release_appsecret();
  inline void set_allocated_appsecret(::std::string* appsecret);

  // @@protoc_insertion_point(class_scope:MiBaseInfo)
 private:
  inline void set_has_appid();
  inline void clear_has_appid();
  inline void set_has_appsecret();
  inline void clear_has_appsecret();

  ::google::protobuf::UnknownFieldSet _unknown_fields_;

  ::std::string* appid_;
  ::std::string* appsecret_;

  mutable int _cached_size_;
  ::google::protobuf::uint32 _has_bits_[(2 + 31) / 32];

  friend void  protobuf_AddDesc_tagalias_5fcommon_2eproto();
  friend void protobuf_AssignDesc_tagalias_5fcommon_2eproto();
  friend void protobuf_ShutdownFile_tagalias_5fcommon_2eproto();

  void InitAsDefaultInstance();
  static MiBaseInfo* default_instance_;
};
// ===================================================================


// ===================================================================

// UserStuff

// optional uint64 uid = 1;
inline bool UserStuff::has_uid() const {
  return (_has_bits_[0] & 0x00000001u) != 0;
}
inline void UserStuff::set_has_uid() {
  _has_bits_[0] |= 0x00000001u;
}
inline void UserStuff::clear_has_uid() {
  _has_bits_[0] &= ~0x00000001u;
}
inline void UserStuff::clear_uid() {
  uid_ = GOOGLE_ULONGLONG(0);
  clear_has_uid();
}
inline ::google::protobuf::uint64 UserStuff::uid() const {
  return uid_;
}
inline void UserStuff::set_uid(::google::protobuf::uint64 value) {
  set_has_uid();
  uid_ = value;
}

// optional .DEVSUBTYPE type = 2 [default = NORMAL_TYPE];
inline bool UserStuff::has_type() const {
  return (_has_bits_[0] & 0x00000002u) != 0;
}
inline void UserStuff::set_has_type() {
  _has_bits_[0] |= 0x00000002u;
}
inline void UserStuff::clear_has_type() {
  _has_bits_[0] &= ~0x00000002u;
}
inline void UserStuff::clear_type() {
  type_ = 1;
  clear_has_type();
}
inline ::DEVSUBTYPE UserStuff::type() const {
  return static_cast< ::DEVSUBTYPE >(type_);
}
inline void UserStuff::set_type(::DEVSUBTYPE value) {
  assert(::DEVSUBTYPE_IsValid(value));
  set_has_type();
  type_ = value;
}

// optional bytes token = 3;
inline bool UserStuff::has_token() const {
  return (_has_bits_[0] & 0x00000004u) != 0;
}
inline void UserStuff::set_has_token() {
  _has_bits_[0] |= 0x00000004u;
}
inline void UserStuff::clear_has_token() {
  _has_bits_[0] &= ~0x00000004u;
}
inline void UserStuff::clear_token() {
  if (token_ != &::google::protobuf::internal::kEmptyString) {
    token_->clear();
  }
  clear_has_token();
}
inline const ::std::string& UserStuff::token() const {
  return *token_;
}
inline void UserStuff::set_token(const ::std::string& value) {
  set_has_token();
  if (token_ == &::google::protobuf::internal::kEmptyString) {
    token_ = new ::std::string;
  }
  token_->assign(value);
}
inline void UserStuff::set_token(const char* value) {
  set_has_token();
  if (token_ == &::google::protobuf::internal::kEmptyString) {
    token_ = new ::std::string;
  }
  token_->assign(value);
}
inline void UserStuff::set_token(const void* value, size_t size) {
  set_has_token();
  if (token_ == &::google::protobuf::internal::kEmptyString) {
    token_ = new ::std::string;
  }
  token_->assign(reinterpret_cast<const char*>(value), size);
}
inline ::std::string* UserStuff::mutable_token() {
  set_has_token();
  if (token_ == &::google::protobuf::internal::kEmptyString) {
    token_ = new ::std::string;
  }
  return token_;
}
inline ::std::string* UserStuff::release_token() {
  clear_has_token();
  if (token_ == &::google::protobuf::internal::kEmptyString) {
    return NULL;
  } else {
    ::std::string* temp = token_;
    token_ = const_cast< ::std::string*>(&::google::protobuf::internal::kEmptyString);
    return temp;
  }
}
inline void UserStuff::set_allocated_token(::std::string* token) {
  if (token_ != &::google::protobuf::internal::kEmptyString) {
    delete token_;
  }
  if (token) {
    set_has_token();
    token_ = token;
  } else {
    clear_has_token();
    token_ = const_cast< ::std::string*>(&::google::protobuf::internal::kEmptyString);
  }
}

// -------------------------------------------------------------------

// MiBaseInfo

// optional bytes appid = 1;
inline bool MiBaseInfo::has_appid() const {
  return (_has_bits_[0] & 0x00000001u) != 0;
}
inline void MiBaseInfo::set_has_appid() {
  _has_bits_[0] |= 0x00000001u;
}
inline void MiBaseInfo::clear_has_appid() {
  _has_bits_[0] &= ~0x00000001u;
}
inline void MiBaseInfo::clear_appid() {
  if (appid_ != &::google::protobuf::internal::kEmptyString) {
    appid_->clear();
  }
  clear_has_appid();
}
inline const ::std::string& MiBaseInfo::appid() const {
  return *appid_;
}
inline void MiBaseInfo::set_appid(const ::std::string& value) {
  set_has_appid();
  if (appid_ == &::google::protobuf::internal::kEmptyString) {
    appid_ = new ::std::string;
  }
  appid_->assign(value);
}
inline void MiBaseInfo::set_appid(const char* value) {
  set_has_appid();
  if (appid_ == &::google::protobuf::internal::kEmptyString) {
    appid_ = new ::std::string;
  }
  appid_->assign(value);
}
inline void MiBaseInfo::set_appid(const void* value, size_t size) {
  set_has_appid();
  if (appid_ == &::google::protobuf::internal::kEmptyString) {
    appid_ = new ::std::string;
  }
  appid_->assign(reinterpret_cast<const char*>(value), size);
}
inline ::std::string* MiBaseInfo::mutable_appid() {
  set_has_appid();
  if (appid_ == &::google::protobuf::internal::kEmptyString) {
    appid_ = new ::std::string;
  }
  return appid_;
}
inline ::std::string* MiBaseInfo::release_appid() {
  clear_has_appid();
  if (appid_ == &::google::protobuf::internal::kEmptyString) {
    return NULL;
  } else {
    ::std::string* temp = appid_;
    appid_ = const_cast< ::std::string*>(&::google::protobuf::internal::kEmptyString);
    return temp;
  }
}
inline void MiBaseInfo::set_allocated_appid(::std::string* appid) {
  if (appid_ != &::google::protobuf::internal::kEmptyString) {
    delete appid_;
  }
  if (appid) {
    set_has_appid();
    appid_ = appid;
  } else {
    clear_has_appid();
    appid_ = const_cast< ::std::string*>(&::google::protobuf::internal::kEmptyString);
  }
}

// optional bytes appsecret = 2;
inline bool MiBaseInfo::has_appsecret() const {
  return (_has_bits_[0] & 0x00000002u) != 0;
}
inline void MiBaseInfo::set_has_appsecret() {
  _has_bits_[0] |= 0x00000002u;
}
inline void MiBaseInfo::clear_has_appsecret() {
  _has_bits_[0] &= ~0x00000002u;
}
inline void MiBaseInfo::clear_appsecret() {
  if (appsecret_ != &::google::protobuf::internal::kEmptyString) {
    appsecret_->clear();
  }
  clear_has_appsecret();
}
inline const ::std::string& MiBaseInfo::appsecret() const {
  return *appsecret_;
}
inline void MiBaseInfo::set_appsecret(const ::std::string& value) {
  set_has_appsecret();
  if (appsecret_ == &::google::protobuf::internal::kEmptyString) {
    appsecret_ = new ::std::string;
  }
  appsecret_->assign(value);
}
inline void MiBaseInfo::set_appsecret(const char* value) {
  set_has_appsecret();
  if (appsecret_ == &::google::protobuf::internal::kEmptyString) {
    appsecret_ = new ::std::string;
  }
  appsecret_->assign(value);
}
inline void MiBaseInfo::set_appsecret(const void* value, size_t size) {
  set_has_appsecret();
  if (appsecret_ == &::google::protobuf::internal::kEmptyString) {
    appsecret_ = new ::std::string;
  }
  appsecret_->assign(reinterpret_cast<const char*>(value), size);
}
inline ::std::string* MiBaseInfo::mutable_appsecret() {
  set_has_appsecret();
  if (appsecret_ == &::google::protobuf::internal::kEmptyString) {
    appsecret_ = new ::std::string;
  }
  return appsecret_;
}
inline ::std::string* MiBaseInfo::release_appsecret() {
  clear_has_appsecret();
  if (appsecret_ == &::google::protobuf::internal::kEmptyString) {
    return NULL;
  } else {
    ::std::string* temp = appsecret_;
    appsecret_ = const_cast< ::std::string*>(&::google::protobuf::internal::kEmptyString);
    return temp;
  }
}
inline void MiBaseInfo::set_allocated_appsecret(::std::string* appsecret) {
  if (appsecret_ != &::google::protobuf::internal::kEmptyString) {
    delete appsecret_;
  }
  if (appsecret) {
    set_has_appsecret();
    appsecret_ = appsecret;
  } else {
    clear_has_appsecret();
    appsecret_ = const_cast< ::std::string*>(&::google::protobuf::internal::kEmptyString);
  }
}


// @@protoc_insertion_point(namespace_scope)

#ifndef SWIG
namespace google {
namespace protobuf {

template <>
inline const EnumDescriptor* GetEnumDescriptor< ::ACTION>() {
  return ::ACTION_descriptor();
}
template <>
inline const EnumDescriptor* GetEnumDescriptor< ::QUERY_ACTION>() {
  return ::QUERY_ACTION_descriptor();
}
template <>
inline const EnumDescriptor* GetEnumDescriptor< ::PLATFORM>() {
  return ::PLATFORM_descriptor();
}
template <>
inline const EnumDescriptor* GetEnumDescriptor< ::SERVER_TYPE>() {
  return ::SERVER_TYPE_descriptor();
}
template <>
inline const EnumDescriptor* GetEnumDescriptor< ::STATUS>() {
  return ::STATUS_descriptor();
}
template <>
inline const EnumDescriptor* GetEnumDescriptor< ::DEVSUBTYPE>() {
  return ::DEVSUBTYPE_descriptor();
}

}  // namespace google
}  // namespace protobuf
#endif  // SWIG

// @@protoc_insertion_point(global_scope)

#endif  // PROTOBUF_tagalias_5fcommon_2eproto__INCLUDED