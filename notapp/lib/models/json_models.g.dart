// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'json_models.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

StandUpPacket _$StandUpPacketFromJson(Map<String, dynamic> json) {
  return StandUpPacket(
      wReadMWriteIndex: json['WReadMWriteIndex'] as int,
      wWriteMReadIndex: json['WWriteMReadIndex'] as int,
      faceCountClose: json['FaceCountClose'] as bool,
      faceCountRecordID: json['FaceCountRecordID'] as int,
      currentPDFPage: json['CurrentPDFPage'] as int,
      changePDFPage: json['ChangePDFPage'] as int,
      pdfUrl: json['PDFUrl'] as String,
      requestStartPacket: json['RequestStartPacket'] as bool);
}

Map<String, dynamic> _$StandUpPacketToJson(StandUpPacket instance) =>
    <String, dynamic>{
      'WReadMWriteIndex': instance.wReadMWriteIndex,
      'WWriteMReadIndex': instance.wWriteMReadIndex,
      'FaceCountClose': instance.faceCountClose,
      'FaceCountRecordID': instance.faceCountRecordID,
      'CurrentPDFPage': instance.currentPDFPage,
      'ChangePDFPage': instance.changePDFPage,
      'PDFUrl': instance.pdfUrl,
      'RequestStartPacket': instance.requestStartPacket
    };

ClassResponse _$ClassResponseFromJson(Map<String, dynamic> json) {
  return ClassResponse(
      classID: json['class_id'] as int,
      className: json['class_name'] as String,
      classImage: json['class_image'] as String,
      faceCount: json['face_count'] as int,
      faceSetToken: json['faceset_token'] as String,
      classroomNo: json['classroom_no'] as String,
      studentNos:
          (json['student_nos'] as List)?.map((e) => e as String)?.toList(),
      teacherNos:
          (json['teacher_nos'] as List)?.map((e) => e as String)?.toList());
}

Map<String, dynamic> _$ClassResponseToJson(ClassResponse instance) =>
    <String, dynamic>{
      'class_id': instance.classID,
      'class_name': instance.className,
      'class_image': instance.classImage,
      'face_count': instance.faceCount,
      'faceset_token': instance.faceSetToken,
      'classroom_no': instance.classroomNo,
      'student_nos': instance.studentNos,
      'teacher_nos': instance.teacherNos
    };

ClassesResponse _$ClassesResponseFromJson(Map<String, dynamic> json) {
  return ClassesResponse(
      classes: (json['classes'] as List)
          ?.map((e) => e == null
              ? null
              : ClassResponse.fromJson(e as Map<String, dynamic>))
          ?.toList(),
      total: json['total'] as int);
}

Map<String, dynamic> _$ClassesResponseToJson(ClassesResponse instance) =>
    <String, dynamic>{'classes': instance.classes, 'total': instance.total};

StandUpStatus _$StandUpStatusFromJson(Map<String, dynamic> json) {
  return StandUpStatus(
      classID: json['class_id'] as int,
      teacherNo: json['teacher_no'] as String,
      wReadMWriteIndex: json['WReadMWriteIndex'] as int,
      wWriteMReadIndex: json['WWriteMReadIndex'] as int);
}

Map<String, dynamic> _$StandUpStatusToJson(StandUpStatus instance) =>
    <String, dynamic>{
      'class_id': instance.classID,
      'teacher_no': instance.teacherNo,
      'WReadMWriteIndex': instance.wReadMWriteIndex,
      'WWriteMReadIndex': instance.wWriteMReadIndex
    };

FaceRectangle _$FaceRectangleFromJson(Map<String, dynamic> json) {
  return FaceRectangle(
      width: json['width'] as int,
      top: json['top'] as int,
      left: json['left'] as int,
      height: json['height'] as int);
}

Map<String, dynamic> _$FaceRectangleToJson(FaceRectangle instance) =>
    <String, dynamic>{
      'width': instance.width,
      'top': instance.top,
      'left': instance.left,
      'height': instance.height
    };

FaceRectToken _$FaceRectTokenFromJson(Map<String, dynamic> json) {
  return FaceRectToken(
      faceRectangle: json['face_rectangle'] == null
          ? null
          : FaceRectangle.fromJson(
              json['face_rectangle'] as Map<String, dynamic>),
      faceToken: json['face_token'] as String);
}

Map<String, dynamic> _$FaceRectTokenToJson(FaceRectToken instance) =>
    <String, dynamic>{
      'face_rectangle': instance.faceRectangle,
      'face_token': instance.faceToken
    };

FaceCountRecordResponse _$FaceCountRecordResponseFromJson(
    Map<String, dynamic> json) {
  return FaceCountRecordResponse(
      faceRectTokens: (json['face_rect_tokens'] as List)
          ?.map((e) => e == null
              ? null
              : FaceRectToken.fromJson(e as Map<String, dynamic>))
          ?.toList(),
      studentCount: json['student_count'] as int,
      studentInClassCount: json['student_in_class_count'] as int,
      studentNotInClass: (json['student_not_in_class'] as List)
          ?.map((e) => e as String)
          ?.toList());
}

Map<String, dynamic> _$FaceCountRecordResponseToJson(
        FaceCountRecordResponse instance) =>
    <String, dynamic>{
      'face_rect_tokens': instance.faceRectTokens,
      'student_count': instance.studentCount,
      'student_in_class_count': instance.studentInClassCount,
      'student_not_in_class': instance.studentNotInClass
    };
