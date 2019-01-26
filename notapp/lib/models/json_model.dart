import 'package:json_annotation/json_annotation.dart';

part 'package:notapp/models/json_model.g.dart';

@JsonSerializable()
class StandUpPacket {
  StandUpPacket({
    this.wReadMWriteIndex = 0,
    this.wWriteMReadIndex = 0,
    this.faceCountClose = false,
    this.currentPDFPage = 0,
    this.changePDFPage = 0,
  });

  @JsonKey(name: 'WReadMWriteIndex')
  int wReadMWriteIndex;

  @JsonKey(name: 'WWriteMReadIndex')
  int wWriteMReadIndex;

  @JsonKey(name: 'FaceCountClose')
  bool faceCountClose;

  @JsonKey(name: 'CurrentPDFPage')
  int currentPDFPage;

  @JsonKey(name: 'ChangePDFPage')
  int changePDFPage;

  factory StandUpPacket.fromJson(Map<String, dynamic> json) => _$StandUpPacketFromJson(json);
  Map<String, dynamic> toJson() => _$StandUpPacketToJson(this);
}