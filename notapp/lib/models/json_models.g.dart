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
      requestStartPacket: json['RequestStartPacket'] as bool,
      studentWarningList: json['StudentWarningList'] as String,
      studentWarningRecordList: (json['StudentWarningRecordList'] as List)
          ?.map((e) => e as int)
          ?.toList());
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
      'RequestStartPacket': instance.requestStartPacket,
      'StudentWarningList': instance.studentWarningList,
      'StudentWarningRecordList': instance.studentWarningRecordList
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

TeacherResponse _$TeacherResponseFromJson(Map<String, dynamic> json) {
  return TeacherResponse(
      teacherNo: json['teacher_no'] as String,
      teacherName: json['teacher_name'] as String,
      teacherImage: json['teacher_image'] as String,
      teacherPassword: json['teacher_password'] as String,
      classIDs: (json['class_ids'] as List)?.map((e) => e as int)?.toList());
}

Map<String, dynamic> _$TeacherResponseToJson(TeacherResponse instance) =>
    <String, dynamic>{
      'teacher_no': instance.teacherNo,
      'teacher_name': instance.teacherName,
      'teacher_image': instance.teacherImage,
      'teacher_password': instance.teacherPassword,
      'class_ids': instance.classIDs
    };

TeachersResponse _$TeachersResponseFromJson(Map<String, dynamic> json) {
  return TeachersResponse(
      teachers: (json['teachers'] as List)
          ?.map((e) => e == null
              ? null
              : TeacherResponse.fromJson(e as Map<String, dynamic>))
          ?.toList(),
      total: json['total'] as int);
}

Map<String, dynamic> _$TeachersResponseToJson(TeachersResponse instance) =>
    <String, dynamic>{'teachers': instance.teachers, 'total': instance.total};

Emotion _$EmotionFromJson(Map<String, dynamic> json) {
  return Emotion(
      sadness: (json['sadness'] as num)?.toDouble(),
      neutral: (json['neutral'] as num)?.toDouble(),
      disgust: (json['disgust'] as num)?.toDouble(),
      anger: (json['anger'] as num)?.toDouble(),
      surprise: (json['surprise'] as num)?.toDouble(),
      fear: (json['fear'] as num)?.toDouble(),
      happiness: (json['happiness'] as num)?.toDouble());
}

Map<String, dynamic> _$EmotionToJson(Emotion instance) => <String, dynamic>{
      'sadness': instance.sadness,
      'neutral': instance.neutral,
      'disgust': instance.disgust,
      'anger': instance.anger,
      'surprise': instance.surprise,
      'fear': instance.fear,
      'happiness': instance.happiness
    };

EyeStatus _$EyeStatusFromJson(Map<String, dynamic> json) {
  return EyeStatus(
      noGlassEyeClose: (json['no_glass_eye_close'] as num)?.toDouble(),
      normalGlassEyeClose: (json['normal_glass_eye_close'] as num)?.toDouble());
}

Map<String, dynamic> _$EyeStatusToJson(EyeStatus instance) => <String, dynamic>{
      'no_glass_eye_close': instance.noGlassEyeClose,
      'normal_glass_eye_close': instance.normalGlassEyeClose
    };

HeadPose _$HeadPoseFromJson(Map<String, dynamic> json) {
  return HeadPose(
      yawAngle: (json['yaw_angle'] as num)?.toDouble(),
      pitchAngle: (json['pitch_angle'] as num)?.toDouble(),
      rollAngle: (json['roll_angle'] as num)?.toDouble());
}

Map<String, dynamic> _$HeadPoseToJson(HeadPose instance) => <String, dynamic>{
      'yaw_angle': instance.yawAngle,
      'pitch_angle': instance.pitchAngle,
      'roll_angle': instance.rollAngle
    };

Attributes _$AttributesFromJson(Map<String, dynamic> json) {
  return Attributes(
      emotion: json['emotion'] == null
          ? null
          : Emotion.fromJson(json['emotion'] as Map<String, dynamic>),
      eyeStatus: json['eyestatus'] == null
          ? null
          : EyeStatus.fromJson(json['eyestatus'] as Map<String, dynamic>),
      headPose: json['headpose'] == null
          ? null
          : HeadPose.fromJson(json['headpose'] as Map<String, dynamic>));
}

Map<String, dynamic> _$AttributesToJson(Attributes instance) =>
    <String, dynamic>{
      'emotion': instance.emotion,
      'eyestatus': instance.eyeStatus,
      'headpose': instance.headPose
    };

StudentStatus _$StudentStatusFromJson(Map<String, dynamic> json) {
  return StudentStatus(
      updateTime: json['update_time'] as int,
      studentNo: json['student_no'] as String,
      attributes: json['attributes'] == null
          ? null
          : Attributes.fromJson(json['attributes'] as Map<String, dynamic>));
}

Map<String, dynamic> _$StudentStatusToJson(StudentStatus instance) =>
    <String, dynamic>{
      'update_time': instance.updateTime,
      'student_no': instance.studentNo,
      'attributes': instance.attributes
    };

StudentStatusWithPage _$StudentStatusWithPageFromJson(
    Map<String, dynamic> json) {
  return StudentStatusWithPage(
      pdfPage: json['pdf_page'] as int,
      studentStatus: (json['students_status'] as List)
          ?.map((e) => e == null
              ? null
              : StudentStatus.fromJson(e as Map<String, dynamic>))
          ?.toList());
}

Map<String, dynamic> _$StudentStatusWithPageToJson(
        StudentStatusWithPage instance) =>
    <String, dynamic>{
      'pdf_page': instance.pdfPage,
      'students_status': instance.studentStatus
    };

StudentWarningRecord _$StudentWarningRecordFromJson(Map<String, dynamic> json) {
  return StudentWarningRecord(
      studentNo: json['student_no'] as String,
      warning: json['warning'] as int,
      lastWarning: json['last_warning'] as bool);
}

Map<String, dynamic> _$StudentWarningRecordToJson(
        StudentWarningRecord instance) =>
    <String, dynamic>{
      'student_no': instance.studentNo,
      'warning': instance.warning,
      'last_warning': instance.lastWarning
    };

StudentStatusResponse _$StudentStatusResponseFromJson(
    Map<String, dynamic> json) {
  return StudentStatusResponse(
      updateTime: json['update_time'] as int,
      classID: json['class_id'] as int,
      className: json['class_name'] as String,
      teacherNo: json['teacher_no'] as String,
      pdf: json['pdf'] as String,
      faceCountRecordID: json['face_count_record_id'] as int,
      studentStatus: (json['student_status'] as List)
          ?.map((e) => e == null
              ? null
              : StudentStatusWithPage.fromJson(e as Map<String, dynamic>))
          ?.toList(),
      studentWarningRecordList: (json['student_warning_record_list'] as List)
          ?.map((e) => e == null
              ? null
              : StudentWarningRecord.fromJson(e as Map<String, dynamic>))
          ?.toList())
    ..pdfPageCount = json['pdf_page_count'] as int;
}

Map<String, dynamic> _$StudentStatusResponseToJson(
        StudentStatusResponse instance) =>
    <String, dynamic>{
      'update_time': instance.updateTime,
      'class_id': instance.classID,
      'class_name': instance.className,
      'teacher_no': instance.teacherNo,
      'pdf': instance.pdf,
      'pdf_page_count': instance.pdfPageCount,
      'face_count_record_id': instance.faceCountRecordID,
      'student_status': instance.studentStatus,
      'student_warning_record_list': instance.studentWarningRecordList
    };

StudentStatusListResponse _$StudentStatusListResponseFromJson(
    Map<String, dynamic> json) {
  return StudentStatusListResponse(
      studentStatus: (json['student_status'] as List)
          ?.map((e) => e == null
              ? null
              : StudentStatusResponse.fromJson(e as Map<String, dynamic>))
          ?.toList(),
      total: json['total'] as int);
}

Map<String, dynamic> _$StudentStatusListResponseToJson(
        StudentStatusListResponse instance) =>
    <String, dynamic>{
      'student_status': instance.studentStatus,
      'total': instance.total
    };

StudentResponse _$StudentResponseFromJson(Map<String, dynamic> json) {
  return StudentResponse(
      studentNo: json['student_no'] as String,
      studentName: json['student_name'] as String,
      faceToken: json['face_token'] as String,
      studentImage: json['student_image'] as String,
      studentPassword: json['student_password'] as String,
      classIDs: (json['class_ids'] as List)?.map((e) => e as int)?.toList());
}

Map<String, dynamic> _$StudentResponseToJson(StudentResponse instance) =>
    <String, dynamic>{
      'student_no': instance.studentNo,
      'student_name': instance.studentName,
      'face_token': instance.faceToken,
      'student_image': instance.studentImage,
      'student_password': instance.studentPassword,
      'class_ids': instance.classIDs
    };

StudentsResponse _$StudentsResponseFromJson(Map<String, dynamic> json) {
  return StudentsResponse(
      students: (json['students'] as List)
          ?.map((e) => e == null
              ? null
              : StudentResponse.fromJson(e as Map<String, dynamic>))
          ?.toList(),
      total: json['total'] as int);
}

Map<String, dynamic> _$StudentsResponseToJson(StudentsResponse instance) =>
    <String, dynamic>{'students': instance.students, 'total': instance.total};
