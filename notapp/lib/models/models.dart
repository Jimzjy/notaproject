import 'package:flutter/material.dart';
import 'package:dio/dio.dart';
import 'json_models.dart';

const SERVER_ADDRESS = "192.168.7.119:8000";

class DioManager {
  static Dio get instance => _getInstance();
  static Dio _instance;

  static Dio _getInstance() {
    if (_instance == null) {
      _instance = new Dio();
      _instance.options
          ..baseUrl = "http://" + SERVER_ADDRESS
          ..connectTimeout = 5000
          ..receiveTimeout = 5000;
    }
    return _instance;
  }
}

class IconFontCN {
  static const IconData scan = const IconData(
      0xe624,
      fontFamily: 'IconFontCN',
      matchTextDirection: true
  );
  static const IconData fear = const IconData(
      0xe6e6,
      fontFamily: 'IconFontCN',
      matchTextDirection: true
  );
  static const IconData frown = const IconData(
      0xe77e,
      fontFamily: 'IconFontCN',
      matchTextDirection: true
  );
  static const IconData meh = const IconData(
      0xe780,
      fontFamily: 'IconFontCN',
      matchTextDirection: true
  );
  static const IconData smile = const IconData(
      0xe783,
      fontFamily: 'IconFontCN',
      matchTextDirection: true
  );
  static const IconData surprise = const IconData(
      0xe6e7,
      fontFamily: 'IconFontCN',
      matchTextDirection: true
  );

  static const Emotion = [
    frown,
    meh,
    frown,
    Icons.whatshot,
    surprise,
    fear,
    smile
  ];
}

class PieChartData {
  int domain;
  num data;

  PieChartData(this.domain, this.data);
}

// By page
class StudentStatusSeparate {
  String studentNo;
  List<List<StudentStatusAttributes>> studentStatusAttributesWithPage;

  StudentStatusSeparate(this.studentNo, this.studentStatusAttributesWithPage);
}

class StudentStatusSeparateByTime {
  String studentNo;
  List<StudentStatusAttributes> studentStatusAttributes;

  StudentStatusSeparateByTime(this.studentNo, this.studentStatusAttributes);
}

class StudentStatusAttributes {
  int updateTime;
  Attributes attributes;

  StudentStatusAttributes(this.updateTime, this.attributes);
}

class StudentWarningRecordTime {
  String studentNos;
  DateTime dateTime;

  StudentWarningRecordTime(this.studentNos, this.dateTime);
}