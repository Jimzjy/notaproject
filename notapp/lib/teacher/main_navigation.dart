import 'package:flutter/material.dart';
import 'package:notapp/widgets/widgets.dart';
import 'package:notapp/models/models.dart';
import 'package:dio/dio.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'dart:convert';
import 'package:notapp/models/json_models.dart';
import 'package:fluttertoast/fluttertoast.dart';
import 'package:qrcode_reader/qrcode_reader.dart';
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
  void didChangeDependencies() {
    _requestClasses();
    super.didChangeDependencies();
  }

  @override
  Widget build(BuildContext context) {
    Color primaryColor = Theme.of(context).primaryColor;

    return new Scaffold(
      appBar: new FakeSearchBar(
        left: new IconButton(icon: Icon(IconFontCN.scan, color: Colors.white,), onPressed: _onScanPressed),
        text: "搜索",
      ),
      body: new Container(
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
                          child: new Text("正在上课:  ${currentClass?.className ?? ""}...", style: TextStyle(color: Colors.white, fontSize: 16.0, fontWeight: FontWeight.bold),),
                        ),
                      ),
                      onTap: () {
                        _navigateToStandUpPage(
                            "/stand_up_mobile?class_id=${_standUpStatus.classID}&write_channel_index=${_standUpStatus.wReadMWriteIndex}",
                            currentClass);
                      },
                    ),
                  );
                } else {
                  return new ClassCard(
                    className: _classesResponse.classes[index-1].className,
                    classID: _classesResponse.classes[index-1].classID,
                    classImage: _classesResponse.classes[index-1].classImage,
                    classroomNo: _classesResponse.classes[index-1].classroomNo,
                    itemPressesCallback: () { _navigateToRecordPage(context, index - 1); },
                  );
                }
              } else {
                return new ClassCard(
                  className: _classesResponse.classes[index].className,
                  classID: _classesResponse.classes[index].classID,
                  classImage: _classesResponse.classes[index].classImage,
                  classroomNo: _classesResponse.classes[index].classroomNo,
                  itemPressesCallback: () { _navigateToRecordPage(context, index); },
                );
              }
            },
          ),
          onRefresh: _requestClasses,
        )
      ),
    );
  }
  
  _navigateToStandUpPage(String wsParam, ClassResponse currentClass) {
    Navigator.push(context, new MaterialPageRoute(
      builder: (context) {
        return new StandUpClassPage(
          wsParam: wsParam,
          classResponse: currentClass,
        );
      },
    ));
  }

  _onScanPressed() {
    Future<String> futureString = new QRCodeReader().scan();
    futureString.then((String value) {
      var params = value.split("|");

      ClassResponse currentClass;
      for (var classResponse in _classesResponse.classes) {
        if (classResponse.classID.toString() == params[0]) {
          currentClass = classResponse;
        }
      }

      _navigateToStandUpPage(
          "/stand_up_mobile?class_id=${params[0]}&write_channel_index=${params[1]}",
          currentClass);
    });
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

    ClassesResponse classesResponse;
    StandUpStatus standUpStatus;
    if (response1?.statusCode == 200) {
      Map jsonMap = jsonDecode(response1.toString());
      classesResponse = ClassesResponse.fromJson(jsonMap);
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
      standUpStatus = StandUpStatus.fromJson(jsonMap);
    } else {
      standUpStatus = null;
    }

    setState(() {
      _classesResponse = classesResponse;
      _standUpStatus = standUpStatus;
    });
    return;
  }

  _navigateToRecordPage(BuildContext context, int index) {
    Navigator.push(context, new MaterialPageRoute(
      builder: (context) {
        return new ClassRecordPage(_classesResponse.classes[index].classID);
      },
    ));
  }
}

class HistoryPage extends StatefulWidget {
  @override
  State<StatefulWidget> createState() => _HistoryPageState();
}

class _HistoryPageState extends State<HistoryPage> {
  String _teacherNo = "";
  List<StudentStatusResponse> _studentStatus;

  @override
  void initState() {
    getTeacherNo().then((v){
      _teacherNo = v;
      _requestStudentStatusRecord();
    });
    super.initState();
  }

  @override
  Widget build(BuildContext context) {
    return new Scaffold(
      appBar: new FakeSearchBar(
        text: "搜索",
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
        child: new RefreshIndicator(
          child: ListView.builder(
            itemCount: _studentStatus?.length ?? 0,
            itemBuilder: (BuildContext context, int index) {
              return new HistoryCard(
                studentStatusResponse: _studentStatus[index],
                itemPressesCallback: () {
                  _navigateToClassPage(context, index);
                },
              );
            },
          ),
          onRefresh: _requestStudentStatusRecord,
        ),
      ),
    );
  }

  _navigateToClassPage(BuildContext context, int index) {
    Navigator.push(context, new MaterialPageRoute(
      builder: (context) {
        return new NormalClassPage(
          studentStatusResponse: _studentStatus[index],
        );
      },
    ));
  }

  Future<void> _requestStudentStatusRecord() async {
    Response response;
    try {
      response = await DioManager.instance.get("/student_status_records", data: { "teacher_no": _teacherNo });
    } catch(e) {
      print(e);
    }

    if (response?.statusCode == 200) {
      Map jsonMap = jsonDecode(response.toString());
      setState(() {
        _studentStatus = StudentStatusListResponse.fromJson(jsonMap).studentStatus;
      });
    } else {
      Fluttertoast.showToast(
        msg: "获取记录失败",
        toastLength: Toast.LENGTH_SHORT,
        timeInSecForIos: 1,
        gravity: ToastGravity.CENTER,
      );
    }
    return;
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
  String userNo = prefs.getString('userNo') ?? "";
  return userNo;
}