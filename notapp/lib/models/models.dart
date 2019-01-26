import 'package:dio/dio.dart';

const SERVER_ADDRESS = "http://192.168.6.229:8000";

class DioManager {
  static Dio get instance => _getInstance();
  static Dio _instance;

  static Dio _getInstance() {
    if (_instance == null) {
      _instance = new Dio();
      _instance.options
          ..baseUrl = SERVER_ADDRESS
          ..connectTimeout = 5000
          ..receiveTimeout = 5000;
    }
    return _instance;
  }
}