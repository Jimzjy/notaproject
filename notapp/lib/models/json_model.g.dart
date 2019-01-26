// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'json_model.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

StandUpPacket _$StandUpPacketFromJson(Map<String, dynamic> json) {
  return StandUpPacket(
      wReadMWriteIndex: json['WReadMWriteIndex'] as int,
      wWriteMReadIndex: json['WWriteMReadIndex'] as int,
      faceCountClose: json['FaceCountClose'] as bool,
      currentPDFPage: json['CurrentPDFPage'] as int,
      changePDFPage: json['ChangePDFPage'] as int);
}

Map<String, dynamic> _$StandUpPacketToJson(StandUpPacket instance) =>
    <String, dynamic>{
      'WReadMWriteIndex': instance.wReadMWriteIndex,
      'WWriteMReadIndex': instance.wWriteMReadIndex,
      'FaceCountClose': instance.faceCountClose,
      'CurrentPDFPage': instance.currentPDFPage,
      'ChangePDFPage': instance.changePDFPage
    };
