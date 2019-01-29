import 'package:json_annotation/json_annotation.dart';

part 'json_models.g.dart';

@JsonSerializable()
class StandUpPacket {
  StandUpPacket({
    this.wReadMWriteIndex = 0,
    this.wWriteMReadIndex = 0,
    this.faceCountClose = false,
    this.faceCountRecordID = 0,
    this.currentPDFPage = 0,
    this.changePDFPage = 0,
    this.pdfUrl = "",
    this.requestStartPacket = false,
  });

  @JsonKey(name: 'WReadMWriteIndex')
  int wReadMWriteIndex;

  @JsonKey(name: 'WWriteMReadIndex')
  int wWriteMReadIndex;

  @JsonKey(name: 'FaceCountClose')
  bool faceCountClose;

  @JsonKey(name: 'FaceCountRecordID')
  int faceCountRecordID;

  @JsonKey(name: 'CurrentPDFPage')
  int currentPDFPage;

  @JsonKey(name: 'ChangePDFPage')
  int changePDFPage;

  @JsonKey(name: 'PDFUrl')
  String pdfUrl;

  @JsonKey(name: 'RequestStartPacket')
  bool requestStartPacket;

  factory StandUpPacket.fromJson(Map<String, dynamic> json) => _$StandUpPacketFromJson(json);
  Map<String, dynamic> toJson() => _$StandUpPacketToJson(this);
}

@JsonSerializable()
class ClassResponse {
  ClassResponse({this.classID, this.className, this.classImage, this.faceCount,
    this.faceSetToken, this.classroomNo, this.studentNos, this.teacherNos});

  @JsonKey(name: 'class_id')
  int classID;
  @JsonKey(name: 'class_name')
  String className;
  @JsonKey(name: 'class_image')
  String classImage;
  @JsonKey(name: 'face_count')
  int faceCount;
  @JsonKey(name: 'faceset_token')
  String faceSetToken;
  @JsonKey(name: 'classroom_no')
  String classroomNo;
  @JsonKey(name: 'student_nos')
  List<String> studentNos;
  @JsonKey(name: 'teacher_nos')
  List<String> teacherNos;

  factory ClassResponse.fromJson(Map<String, dynamic> json) => _$ClassResponseFromJson(json);
  Map<String, dynamic> toJson() => _$ClassResponseToJson(this);
}

@JsonSerializable()
class ClassesResponse {
  ClassesResponse({this.classes, this.total});

  @JsonKey(name: 'classes')
  List<ClassResponse> classes;
  @JsonKey(name: 'total')
  int total;

  factory ClassesResponse.fromJson(Map<String, dynamic> json) => _$ClassesResponseFromJson(json);
  Map<String, dynamic> toJson() => _$ClassesResponseToJson(this);
}

@JsonSerializable()
class StandUpStatus {
  StandUpStatus({this.classID, this.teacherNo, this.wReadMWriteIndex,
      this.wWriteMReadIndex});

  @JsonKey(name: 'class_id')
  int classID;

  @JsonKey(name: 'teacher_no')
  String teacherNo;

  @JsonKey(name: 'WReadMWriteIndex')
  int wReadMWriteIndex;

  @JsonKey(name: 'WWriteMReadIndex')
  int wWriteMReadIndex;

  factory StandUpStatus.fromJson(Map<String, dynamic> json) => _$StandUpStatusFromJson(json);
  Map<String, dynamic> toJson() => _$StandUpStatusToJson(this);
}

@JsonSerializable()
class FaceRectangle {
  FaceRectangle({this.width, this.top, this.left, this.height});

  int width;
  int top;
  int left;
  int height;

  factory FaceRectangle.fromJson(Map<String, dynamic> json) => _$FaceRectangleFromJson(json);
  Map<String, dynamic> toJson() => _$FaceRectangleToJson(this);
}

@JsonSerializable()
class FaceRectToken {
  FaceRectToken({this.faceRectangle, this.faceToken});

  @JsonKey(name: "face_rectangle")
  FaceRectangle faceRectangle;
  @JsonKey(name: "face_token")
  String faceToken;

  factory FaceRectToken.fromJson(Map<String, dynamic> json) => _$FaceRectTokenFromJson(json);
  Map<String, dynamic> toJson() => _$FaceRectTokenToJson(this);
}

@JsonSerializable()
class FaceCountRecordResponse {
  FaceCountRecordResponse({this.faceRectTokens, this.studentCount,
      this.studentInClassCount, this.studentNotInClass});

  @JsonKey(name: "face_rect_tokens")
  List<FaceRectToken> faceRectTokens;
  @JsonKey(name: "student_count")
  int studentCount;
  @JsonKey(name: "student_in_class_count")
  int studentInClassCount;
  @JsonKey(name: "student_not_in_class")
  List<String> studentNotInClass;

  factory FaceCountRecordResponse.fromJson(Map<String, dynamic> json) => _$FaceCountRecordResponseFromJson(json);
  Map<String, dynamic> toJson() => _$FaceCountRecordResponseToJson(this);
}

@JsonSerializable()
class TeacherResponse {
  TeacherResponse({
    this.teacherNo = "0",
    this.teacherName = "0",
    this.teacherImage = "",
    this.teacherPassword = "",
    this.classIDs,
  });

  @JsonKey(name: "teacher_no")
  String teacherNo;
  @JsonKey(name: "teacher_name")
  String teacherName;
  @JsonKey(name: "teacher_image")
  String teacherImage;
  @JsonKey(name: "teacher_password")
  String teacherPassword;
  @JsonKey(name: "class_ids")
  List<int> classIDs;

  factory TeacherResponse.fromJson(Map<String, dynamic> json) => _$TeacherResponseFromJson(json);
  Map<String, dynamic> toJson() => _$TeacherResponseToJson(this);
}

@JsonSerializable()
class TeachersResponse {
  TeachersResponse({this.teachers, this.total});

  List<TeacherResponse> teachers;
  int total;

  factory TeachersResponse.fromJson(Map<String, dynamic> json) => _$TeachersResponseFromJson(json);
  Map<String, dynamic> toJson() => _$TeachersResponseToJson(this);
}