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
    this.studentWarringList = "",
    this.studentWarningRecordList,
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

  @JsonKey(name: 'StudentWarringList')
  String studentWarringList;

  @JsonKey(name: 'StudentWarringRecordList')
  List<int> studentWarningRecordList;

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

@JsonSerializable()
class Emotion {
  Emotion({this.sadness, this.neutral, this.disgust, this.anger, this.surprise,
    this.fear, this.happiness});

  double sadness;
  double neutral;
  double disgust;
  double anger;
  double surprise;
  double fear;
  double happiness;

  factory Emotion.fromJson(Map<String, dynamic> json) => _$EmotionFromJson(json);
  Map<String, dynamic> toJson() => _$EmotionToJson(this);
}

@JsonSerializable()
class EyeStatus {
  EyeStatus({this.noGlassEyeClose, this.normalGlassEyeClose});

  @JsonKey(name: "no_glass_eye_close")
  double noGlassEyeClose;
  @JsonKey(name: "normal_glass_eye_close")
  double normalGlassEyeClose;

  factory EyeStatus.fromJson(Map<String, dynamic> json) => _$EyeStatusFromJson(json);
  Map<String, dynamic> toJson() => _$EyeStatusToJson(this);
}

@JsonSerializable()
class HeadPose {
  HeadPose({this.yawAngle, this.pitchAngle, this.rollAngle});

  @JsonKey(name: "yaw_angle")
  double yawAngle;
  @JsonKey(name: "pitch_angle")
  double pitchAngle;
  @JsonKey(name: "roll_angle")
  double rollAngle;

  factory HeadPose.fromJson(Map<String, dynamic> json) => _$HeadPoseFromJson(json);
  Map<String, dynamic> toJson() => _$HeadPoseToJson(this);
}

@JsonSerializable()
class Attributes {
  Attributes({this.emotion, this.eyeStatus, this.headPose});

  Emotion emotion;
  @JsonKey(name: "eyes_status")
  EyeStatus eyeStatus;
  @JsonKey(name: "head_pose")
  HeadPose headPose;

  factory Attributes.fromJson(Map<String, dynamic> json) => _$AttributesFromJson(json);
  Map<String, dynamic> toJson() => _$AttributesToJson(this);
}

@JsonSerializable()
class StudentStatus {
  StudentStatus({this.updateTime, this.studentNo, this.attributes});

  @JsonKey(name: "update_time")
  int updateTime;
  @JsonKey(name: "student_no")
  String studentNo;
  Attributes attributes;

  factory StudentStatus.fromJson(Map<String, dynamic> json) => _$StudentStatusFromJson(json);
  Map<String, dynamic> toJson() => _$StudentStatusToJson(this);
}

@JsonSerializable()
class StudentStatusWithPage {
  StudentStatusWithPage({this.pdfPage, this.studentStatus});

  @JsonKey(name: "pdf_page")
  int pdfPage;
  @JsonKey(name: "students_status")
  List<StudentStatus> studentStatus;

  factory StudentStatusWithPage.fromJson(Map<String, dynamic> json) => _$StudentStatusWithPageFromJson(json);
  Map<String, dynamic> toJson() => _$StudentStatusWithPageToJson(this);
}

@JsonSerializable()
class StudentStatusResponse {
  StudentStatusResponse({this.updateTime, this.classID, this.className, this.teacherNo, this.pdf,
    this.faceCountRecordID, this.studentStatus});

  @JsonKey(name: "update_time")
  int updateTime;
  @JsonKey(name: "class_id")
  int classID;
  @JsonKey(name: "class_name")
  String className;
  @JsonKey(name: "teacher_no")
  String teacherNo;
  String pdf;
  @JsonKey(name: "face_count_record_id")
  int faceCountRecordID;
  @JsonKey(name: "student_status")
  List<StudentStatusWithPage> studentStatus;

  factory StudentStatusResponse.fromJson(Map<String, dynamic> json) => _$StudentStatusResponseFromJson(json);
  Map<String, dynamic> toJson() => _$StudentStatusResponseToJson(this);
}

@JsonSerializable()
class StudentStatusListResponse {
  StudentStatusListResponse({this.studentStatus, this.total});

  @JsonKey(name: "student_status")
  List<StudentStatusResponse> studentStatus;
  int total;

  factory StudentStatusListResponse.fromJson(Map<String, dynamic> json) => _$StudentStatusListResponseFromJson(json);
  Map<String, dynamic> toJson() => _$StudentStatusListResponseToJson(this);
}

@JsonSerializable()
class StudentResponse {
  StudentResponse({this.studentNo, this.studentName, this.faceToken,
      this.studentImage, this.studentPassword, this.classIDs});

  @JsonKey(name: "student_no")
  String studentNo;
  @JsonKey(name: "student_name")
  String studentName;
  @JsonKey(name: "face_token")
  String faceToken;
  @JsonKey(name: "student_image")
  String studentImage;
  @JsonKey(name: "student_password")
  String studentPassword;
  @JsonKey(name: "class_ids")
  List<int> classIDs;

  factory StudentResponse.fromJson(Map<String, dynamic> json) => _$StudentResponseFromJson(json);
  Map<String, dynamic> toJson() => _$StudentResponseToJson(this);
}

@JsonSerializable()
class StudentsResponse {
  StudentsResponse({this.students, this.total});

  List<StudentResponse> students;
  int total;

  factory StudentsResponse.fromJson(Map<String, dynamic> json) => _$StudentsResponseFromJson(json);
  Map<String, dynamic> toJson() => _$StudentsResponseToJson(this);
}