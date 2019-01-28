import 'package:flutter/material.dart';
import 'package:notapp/widgets/widgets.dart';
import 'package:notapp/models/models.dart';
import 'package:dio/dio.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'dart:convert';
import 'package:notapp/models/json_models.dart';
import 'package:fluttertoast/fluttertoast.dart';

class ClassesPage extends StatefulWidget {

  @override
  State<StatefulWidget> createState() => _ClassPageState();
}

class _ClassPageState extends State<ClassesPage> {
  Dio dio = DioManager.instance;
  String _teacherNo = "";
  ClassesResponse _classesResponse;
  StandUpStatus _standUpStatus;

  @override
  void initState() {
    _getTeacherNo().then((v) {
      _teacherNo = v;
      _requestClasses();
    });

    super.initState();
  }

  @override
  Widget build(BuildContext context) {
    Color primaryColor = Theme.of(context).primaryColor;

    return new Scaffold(
      appBar: new FakeSearchBar(
        left: new IconButton(icon: Icon(IconFontCN.scan, color: Colors.white,), onPressed: () {}),
        text: "搜索",
      ),
      body: new Container(
        constraints: BoxConstraints(minHeight: double.infinity, minWidth: double.infinity),
        decoration: BoxDecoration(
          gradient: LinearGradient(
            colors: [primaryColor, Colors.white],
            begin: Alignment.topCenter,
            end: Alignment.bottomCenter,
          ),
        ),
        child: new RefreshIndicator(
          child: new ListView.builder(
            itemCount: _getClassListSize(),
            itemBuilder: (BuildContext context, int index) {
              if (_standUpStatus != null) {
                var className = "...";

                for (var classResponse in _classesResponse.classes) {
                  if (classResponse.classID == _standUpStatus.classID) {
                    className = classResponse.className;
                  }
                }

                if (index == 0) {
                  return new Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 16.0, vertical: 8.0),
                    child: new GestureDetector(
                      child: new Container(
                        decoration: BoxDecoration(
                          gradient: LinearGradient(
                            colors: [Colors.orange, Colors.amber],
                          ),
                          borderRadius: BorderRadius.circular(12.0),
                        ),
                        child: new Padding(
                          padding: const EdgeInsets.all(16.0),
                          child: new Text("正在上课:  $className...", style: TextStyle(color: Colors.white, fontSize: 16.0, fontWeight: FontWeight.bold),),
                        ),
                      ),
                      onTap: () {

                      },
                    ),
                  );
                } else {
                  return new ClassCard(
                    className: _classesResponse.classes[index-1].className,
                    classID: _classesResponse.classes[index-1].classID,
                    classImage: _classesResponse.classes[index-1].classImage,
                    itemPressesCallback: () {},
                  );
                }
              } else {
                return new ClassCard(
                  className: _classesResponse.classes[index].className,
                  classID: _classesResponse.classes[index].classID,
                  classImage: _classesResponse.classes[index].classImage,
                  itemPressesCallback: () {},
                );
              }
            },
          ),
          onRefresh: _requestClasses,
        )
      ),
    );
  }

  int _getClassListSize() {
    int size = 0;

    if (_standUpStatus != null) {
      size = 1 + _classesResponse?.total ?? 0;
    } else {
      size = _classesResponse?.total ?? 0;
    }

    return size;
  }

  Future<String> _getTeacherNo() async {
    SharedPreferences prefs = await SharedPreferences.getInstance();
    String userType = prefs.getString('teacherNo') ?? "";
    return userType;
  }

  Future<void> _requestClasses() async {
    Response response1;
    Response response2;
    try {
      response1 = await dio.get("/classes", data: { "teacher_no": _teacherNo});
      response2 = await dio.get("/current_stand_up", data: { "teacher_no": _teacherNo});
    } catch (err) {
      print(err);
    }

    if (response1?.statusCode == 200) {
      Map jsonMap = jsonDecode(response1.toString());
      _classesResponse = ClassesResponse.fromJson(jsonMap);
    } else {
      Fluttertoast.showToast(
        msg: "获取课程列表失败",
        toastLength: Toast.LENGTH_SHORT,
        timeInSecForIos: 1,
        gravity: ToastGravity.CENTER,
      );
    }

    if (response2?.statusCode == 200) {
      Map jsonMap = jsonDecode(response2.toString());
      _standUpStatus = StandUpStatus.fromJson(jsonMap);
    } else {
      _standUpStatus = null;
    }

    setState(() {});
    return;
  }
}

class HistoryPage extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    // TODO: implement build
    return new Text("History");
  }
}

class MinePage extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    // TODO: implement build
    return new Text("Mine");
  }
}