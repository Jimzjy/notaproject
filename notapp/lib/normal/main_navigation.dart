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
  String _studentNo = "";
  ClassesResponse _classesResponse;

  @override
  void initState() {
    getStudentNo().then((v) {
      _studentNo = v;
      _requestClasses();
    });

    super.initState();
  }

  @override
  Widget build(BuildContext context) {
    Color primaryColor = Theme.of(context).primaryColor;

    return new Scaffold(
      appBar: new FakeSearchBar(
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
              itemCount: _classesResponse?.total ?? 0,
              itemBuilder: (BuildContext context, int index) {
                return new ClassCard(
                  className: _classesResponse.classes[index].className,
                  classID: _classesResponse.classes[index].classID,
                  classImage: _classesResponse.classes[index].classImage,
                  classroomNo: _classesResponse.classes[index].classroomNo,
                  itemPressesCallback: () { _navigateToRecordPage(context, index); },
                );
              },
            ),
            onRefresh: _requestClasses,
          )
      ),
    );
  }

  Future<void> _requestClasses() async {
    Response response1;
    try {
      response1 = await dio.get("/classes", data: { "student_no": _studentNo});
    } catch (err) {
      print(err);
    }

    ClassesResponse classesResponse;
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

    setState(() {
      _classesResponse = classesResponse;
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

class MinePage extends StatefulWidget {
  @override
  State<StatefulWidget> createState() => _MinePageState();
}

class _MinePageState extends State {
  StudentResponse _studentResponse = new StudentResponse(classIDs: [0]);

  @override
  void initState() {
    getStudentNo().then((v) {
      _requestStudent(v);
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
                  child: _studentResponse.studentImage != "" ? new Image.network("http://$SERVER_ADDRESS/images/${_studentResponse.studentImage}") : null,
                )
            ),
            new Text(_studentResponse.studentName, style: TextStyle(fontWeight: FontWeight.bold, fontSize: 16.0, color: Colors.white),),
            new Text(_studentResponse.studentNo, style: TextStyle(fontSize: 12.0, color: Colors.white),)
          ],
        ),
      ),
    );
  }

  Future<void> _requestStudent(String studentNo) async {
    Response response;
    try {
      response = await DioManager.instance.get("/students", data: { "student_no": studentNo });
    } catch(e) {
      print(e);
    }

    if (response?.statusCode == 200) {
      Map jsonMap = jsonDecode(response.toString());
      StudentsResponse stuResponse = StudentsResponse.fromJson(jsonMap);
      setState(() {
        _studentResponse = stuResponse.students[0];
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

Future<String> getStudentNo() async {
  SharedPreferences prefs = await SharedPreferences.getInstance();
  String userNo = prefs.getString('userNo') ?? "";
  return userNo;
}