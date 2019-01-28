import 'package:flutter/material.dart';
import 'package:dio/dio.dart';

const SERVER_ADDRESS = "192.168.6.229:8000";

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
}