import 'package:flutter/material.dart';
import 'package:notapp/widgets/widgets.dart';
import 'package:notapp/models/models.dart';
import 'package:dio/dio.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'dart:convert';
import 'package:notapp/models/json_models.dart';
import 'package:fluttertoast/fluttertoast.dart';
import 'sub_navigation.dart';

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
    getTeacherNo().then((v) {
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
                ClassResponse currentClass;

                for (var classResponse in _classesResponse.classes) {
                  if (classResponse.classID == _standUpStatus.classID) {
                    currentClass = classResponse;
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
                          child: new Text("正在上课:  ${currentClass.className}...", style: TextStyle(color: Colors.white, fontSize: 16.0, fontWeight: FontWeight.bold),),
                        ),
                      ),
                      onTap: () {
                        Navigator.push(context, new MaterialPageRoute(
                          builder: (context) {
                            return new StandUpClassPage(
                              classID: _standUpStatus.classID,
                              wReadMWriteIndex: _standUpStatus.wReadMWriteIndex,
                              classResponse: currentClass,
                            );
                          },
                        ));
                      },
                    ),
                  );
                } else {
                  return new ClassCard(
                    className: _classesResponse.classes[index-1].className,
                    classID: _classesResponse.classes[index-1].classID,
                    classImage: _classesResponse.classes[index-1].classImage,
                    classroomNo: _classesResponse.classes[index-1].classroomNo,
                    itemPressesCallback: () {},
                  );
                }
              } else {
                return new ClassCard(
                  className: _classesResponse.classes[index].className,
                  classID: _classesResponse.classes[index].classID,
                  classImage: _classesResponse.classes[index].classImage,
                  classroomNo: _classesResponse.classes[index].classroomNo,
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

class MinePage extends StatefulWidget {
  @override
  State<StatefulWidget> createState() => _MinePageState();
}

class _MinePageState extends State {
  TeacherResponse _teacherResponse = new TeacherResponse(classIDs: [0]);

  @override
  void initState() {
    getTeacherNo().then((v) {
      _requestTeacher(v);
    });
    super.initState();
  }

  @override
  Widget build(BuildContext context) {
    return new Scaffold(
      appBar: new AppBar(
        elevation: 0,
        actions: <Widget>[
          new IconButton(icon: Icon(Icons.settings), onPressed: (){
            Navigator.push(context, new MaterialPageRoute(builder: (context) {
              return new SettingPage();
            }));
          })
        ],
      ),
      body: new Container(
        decoration: new BoxDecoration(
          color: Theme.of(context).primaryColor,
          gradient: LinearGradient(
            colors: [Theme.of(context).primaryColor, Colors.white],
            begin: Alignment.topCenter,
            end: Alignment.bottomCenter,
          ),
        ),
        alignment: Alignment.center,
        padding: const EdgeInsets.only(top: 16.0),
        child: Column(
          children: <Widget>[
            new Container(
              margin: const EdgeInsets.only(bottom: 12.0),
              height: 108,
              width: 108,
              decoration: new BoxDecoration(
                borderRadius: new BorderRadius.circular(54.0),
                color: Colors.white,
              ),
              child: new ClipOval(
                child: _teacherResponse.teacherImage != "" ? new Image.network("http://$SERVER_ADDRESS/images/${_teacherResponse.teacherImage}") : null,
              )
            ),
            new Text(_teacherResponse.teacherName, style: TextStyle(fontWeight: FontWeight.bold, fontSize: 16.0, color: Colors.white),),
            new Text(_teacherResponse.teacherNo, style: TextStyle(fontSize: 12.0, color: Colors.white),)
          ],
        ),
      ),
    );
  }

  Future<void> _requestTeacher(String teacherNo) async {
    Response response;
    try {
      response = await DioManager.instance.get("/teachers", data: { "teacher_no": teacherNo });
    } catch(e) {
      print(e);
    }

    if (response?.statusCode == 200) {
      Map jsonMap = jsonDecode(response.toString());
      TeachersResponse tsResponse = TeachersResponse.fromJson(jsonMap);
      setState(() {
        _teacherResponse = tsResponse.teachers[0];
      });
    } else {
      Fluttertoast.showToast(
        msg: "获取信息失败",
        toastLength: Toast.LENGTH_SHORT,
        timeInSecForIos: 1,
        gravity: ToastGravity.CENTER,
      );
    }
  }
}

Future<String> getTeacherNo() async {
  SharedPreferences prefs = await SharedPreferences.getInstance();
  String userType = prefs.getString('teacherNo') ?? "";
  return userType;
}