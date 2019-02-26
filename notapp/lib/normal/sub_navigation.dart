import 'package:flutter/material.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:notapp/main.dart' show NOT_LOGIN;
import 'package:notapp/login/login_main.dart';
import 'package:notapp/models/json_models.dart';
import 'package:notapp/models/models.dart';
import 'package:notapp/widgets/widgets.dart';
import 'package:dio/dio.dart';
import 'package:fluttertoast/fluttertoast.dart';
import 'package:notapp/teacher/sub_navigation.dart' show EMOTION_TEXT, PITCH_ANGLE, EYE_CLOSE;
import 'main_navigation.dart' show getStudentNo;
import 'dart:convert';
import 'dart:math';

class SettingPage extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return new Scaffold(
      appBar: new AppBar(
        title: new Text("设置"),
        elevation: 0.0,
      ),
      body: new Container(
        color: Theme.of(context).primaryColor,
        child: new Center(
          child: new MaterialButton(
            color: Colors.white,
            onPressed: () => _onLogoutButtonPressed(context),
            child: new Text("退出登录"),
          ),
        ),
      ),
    );
  }

  _onLogoutButtonPressed(BuildContext context) async {
    SharedPreferences prefs = await SharedPreferences.getInstance();
    await prefs.setString('userType', NOT_LOGIN);

    Navigator.pushReplacement(context, new MaterialPageRoute(builder: (context) {
      return new LoginApp();
    }));
  }
}

class ClassRecordPage extends StatefulWidget {
  ClassRecordPage(this.classID);

  final int classID;

  @override
  State<StatefulWidget> createState() => _ClassRecordPage(classID);
}

class _ClassRecordPage extends State<ClassRecordPage> {
  _ClassRecordPage(this.classID);

  int classID;
  List<StudentStatusResponse> _studentStatus;

  @override
  void initState() {
    _requestStudentStatusRecord();
    super.initState();
  }

  @override
  Widget build(BuildContext context) {
    return new Scaffold(
      appBar: new AppBar(
        title: new Text("记录", overflow: TextOverflow.ellipsis,),
        elevation: 0.0,
      ),
      body: new Container(
        decoration: new BoxDecoration(color: Theme.of(context).primaryColor),
        child: new ListView.builder(
          itemCount: _studentStatus?.length ?? 0,
          itemBuilder: (BuildContext context, int index) {
            return new HistoryCard(
              studentStatusResponse: _studentStatus[index],
              itemPressesCallback: () {
                _navigateToStudentStatusPage(context, index);
              },
            );
          },
        ),
      ),
    );
  }

  _requestStudentStatusRecord() async {
    Response response;
    try {
      response = await DioManager.instance.get("/student_status_records", data: { "class_id": classID });
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

  _navigateToStudentStatusPage(BuildContext context, int index) {
    Navigator.push(context, new MaterialPageRoute(
      builder: (context) {
        return new StudentStatusPage(_studentStatus[index]);
      },
    ));
  }
}

class StudentStatusPage extends StatefulWidget {
  StudentStatusPage(this.studentStatus);

  final StudentStatusResponse studentStatus;

  @override
  State<StatefulWidget> createState() => _StudentStatusPageState();
}

class _StudentStatusPageState extends State<StudentStatusPage> with SingleTickerProviderStateMixin {

  StudentStatusSeparate studentStatusSeparate;
  StudentResponse studentResponse;
  StudentStatusSeparateByTime studentStatusSeparateByTime;
  int currentStatusCount = 0;
  int cubeAxis = 0;
  double cubeAngle = 0.0;
  Animation<double> animation;
  AnimationController controller;
  bool first = true;

  // front left right top bottom back
  final List<Container> studentHeadCube = [];
  final cubeSize = 80.0;
  final cubeRadius = 40.0;

  @override
  void initState() {
    List<List<StudentStatusAttributes>> separateList = [];
    List<StudentStatusAttributes> byTimeList = [];

    controller = new AnimationController(
        duration: const Duration(milliseconds: 500), vsync: this);
    animation = CurvedAnimation(parent: controller, curve: Curves.easeInOut);
    animation.addListener(() {
      setState(() {});
    });

    getStudentNo().then((no) {
      _requestStudent(no).then((v) {
        studentResponse = v;

        for (var i = 0; i < widget.studentStatus.pdfPageCount; i++) {
          separateList.add([]);
        }

        widget.studentStatus.studentStatus.forEach((StudentStatusWithPage _studentStatusWithPage) {
          int _page = _studentStatusWithPage.pdfPage - 1;

          for (var j = 0; j < (_studentStatusWithPage?.studentStatus?.length ?? 0); j++) {
            var _status = _studentStatusWithPage.studentStatus[j];
            if (_status.studentNo != studentResponse.studentNo) {
              continue;
            }

            separateList[_page].add(StudentStatusAttributes(_status.updateTime, _status.attributes));
            byTimeList.add(StudentStatusAttributes(_status.updateTime, _status.attributes));
          }
        });

        studentStatusSeparate = StudentStatusSeparate(v.studentNo, separateList);
        studentStatusSeparateByTime = StudentStatusSeparateByTime(v.studentNo, byTimeList);

        if (studentStatusSeparateByTime.studentStatusAttributes.length <= currentStatusCount) {
          setState(() {});
          return;
        }
        studentHeadCube.add(Container(
          height: cubeSize,
          width: cubeSize,
          decoration: BoxDecoration(
              color: Colors.blue[300],
              border: Border.all(color: Colors.indigo[500])
          ),
        ));
        studentHeadCube.add(Container());
        _changeCubeState();
        setState(() {});
      });
    });
    super.initState();
  }

  @override
  void dispose() {
    controller.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    if (studentResponse == null || studentStatusSeparateByTime == null || studentStatusSeparate == null) {
      return new Scaffold();
    }

    return new Scaffold(
        appBar: new AppBar(
          title: new Text("学生信息"),
          elevation: 0.0,
        ),
        body: Container(
          color: Theme.of(context).primaryColor,
          child: NestedScrollView(
            headerSliverBuilder: (BuildContext context, bool innerBoxIsScrolled) {
              return <Widget>[
                new SliverList(delegate: new SliverChildListDelegate([
                  _buildStudentCard(),
                  _buildStudentStatusCard(),
                ]))
              ];
            },
            body: _buildPageStatusList(),
          ),
        )
    );
  }

  Future<StudentResponse> _requestStudent(String studentNo) async {
    Response response;
    try {
      response = await DioManager.instance.get("/students", data: { "student_no": studentNo });
    } catch(e) {
      print(e);
    }

    StudentResponse stuResponse;
    if (response?.statusCode == 200) {
      Map jsonMap = jsonDecode(response.toString());
      stuResponse = StudentsResponse.fromJson(jsonMap).students[0];
    } else {
      Fluttertoast.showToast(
        msg: "获取信息失败",
        toastLength: Toast.LENGTH_SHORT,
        timeInSecForIos: 1,
        gravity: ToastGravity.CENTER,
      );
    }
    return stuResponse;
  }

  Widget _buildStudentCard() {
    return new Card(
      margin: const EdgeInsets.symmetric(vertical: 4.0, horizontal: 16.0),
      child: new Container(
        padding: const EdgeInsets.all(12.0),
        child: new Row(
          crossAxisAlignment: CrossAxisAlignment.center,
          children: <Widget>[
            new Container(
              margin: const EdgeInsets.only(left: 8.0, right: 24.0),
              height: 60,
              width: 60,
              child: CircleAvatar(backgroundImage: Image.network("http://$SERVER_ADDRESS/images/${studentResponse.studentImage}").image,),
            ),
            new Expanded(
              child: new Column(
                mainAxisAlignment: MainAxisAlignment.start,
                crossAxisAlignment: CrossAxisAlignment.start,
                children: <Widget>[
                  new Text("学号: ${studentResponse.studentNo}", style: TextStyle(fontSize: 16.0, fontWeight: FontWeight.bold),),
                  new Padding(padding: const EdgeInsets.only(top: 6.0), child: new Text("姓名: ${studentResponse.studentName}", style: TextStyle(fontSize: 14.0),),)
                ],
              ),
            )
          ],
        ),
      ),
    );
  }

  Container _buildCubeFace(IconData iconData) {
    return new Container(
      height: cubeSize,
      width: cubeSize,
      decoration: BoxDecoration(
          color: Colors.blue[300],
          border: Border.all(color: Colors.indigo[500])
      ),
      alignment: Alignment.center,
      child: new Icon(iconData, size: 40, color: Colors.white,),
    );
  }

  Widget _buildStudentStatusCard() {
    if (studentStatusSeparateByTime.studentStatusAttributes.length <= currentStatusCount) {
      return new Container();
    }

    var _date = DateTime.fromMillisecondsSinceEpoch(studentStatusSeparateByTime.studentStatusAttributes[currentStatusCount].updateTime * 1000);
    var angle = first ? cubeAngle : cubeAngle * animation.value;
    var cosAR = cubeRadius * cos(angle);
    var sinAR = cubeRadius * sin(angle);

    return new Card(
      margin: const EdgeInsets.symmetric(vertical: 4.0, horizontal: 16.0),
      child: new Container(
          height: 200,
          padding: const EdgeInsets.all(12.0),
          child: new Stack(
            children: <Widget>[
              new Center(child: new Stack(
                children: <Widget>[
                  _buildCubeFaceSlice(cubeAxis, angle, cosAR, sinAR),
                  _buildCubeSlice(cubeAxis, angle, cosAR, sinAR),
                ],
              ),),
              new Positioned(child: Container(
                alignment: Alignment.center,
                child: IconButton(icon: Icon(Icons.navigate_before), onPressed: () => _changeStudentStatus(-1)),
              ), left: 0, top: 0, bottom: 0,),
              new Positioned(child: Container(
                alignment: Alignment.center,
                child: IconButton(icon: Icon(Icons.navigate_next), onPressed: () => _changeStudentStatus(1)),
              ), right: 0, top: 0, bottom: 0,),
              new Positioned(child: new Text("${_date.hour}:${_date.minute}:${_date.second}", textAlign: TextAlign.center, style: TextStyle(fontWeight: FontWeight.bold),),
                bottom: 0, left: 0, right: 0,),
            ],
          )
      ),
    );
  }

  // axis: x y
  Transform _buildCubeSlice(int axis, double angle, double cosAR, double sinAR) {
    Matrix4 _transform = Matrix4.identity();
    var _cosAR = angle > 0 ? cosAR : -cosAR;
    var _sinAR = angle > 0 ? sinAR : -sinAR;

    switch(axis) {
      case 0:
        _transform
          ..setEntry(3, 2, 0.001)
          ..translate(0.0, -_cosAR, -_sinAR)
          ..rotateX(angle + pi / 2);
        break;
      case 1:
        _transform
          ..setEntry(3, 2, 0.001)
          ..translate(_cosAR, 0.0, -_sinAR)
          ..rotateY(angle + pi / 2);
        break;
    }

    return Transform(
      alignment: Alignment.center,
      transform: _transform,
      child: Center(
        child: studentHeadCube[0],
      ),
    );
  }

  Transform _buildCubeFaceSlice(int axis, double angle, double cosAR, double sinAR) {
    Matrix4 _transform = Matrix4.identity();

    switch(axis) {
      case 0:
        _transform
          ..setEntry(3, 2, 0.001)
          ..translate(0.0, sinAR, -cosAR)
          ..rotateX(angle);
        break;
      case 1:
        _transform
          ..setEntry(3, 2, 0.001)
          ..translate(-sinAR, 0.0, -cosAR)
          ..rotateY(angle);
        break;
    }

    return Transform(
      alignment: Alignment.center,
      transform: _transform,
      child: Center(
        child: studentHeadCube[1],
      ),
    );
  }

  _changeCubeState() {
    var pitch = studentStatusSeparateByTime.studentStatusAttributes[currentStatusCount].attributes.headPose.pitchAngle;
    var yaw = studentStatusSeparateByTime.studentStatusAttributes[currentStatusCount].attributes.headPose.yawAngle;

    if (pitch.abs() > yaw.abs()) {
      cubeAngle = pitch / 180 * pi;
      cubeAxis = 0;
    } else {
      cubeAngle = yaw / 180 * pi;
      cubeAxis = 1;
    }

    var emotion = studentStatusSeparateByTime.studentStatusAttributes[currentStatusCount].attributes.emotion.getEmotion();

    studentHeadCube[1] = _buildCubeFace(IconFontCN.Emotion[emotion]);
  }

  // stepDir: next > 0, back < 0
  _changeStudentStatus(int stepDir) {
    if (stepDir > 0) {
      if (currentStatusCount < studentStatusSeparateByTime.studentStatusAttributes.length - 1) {
        currentStatusCount += 1;
      } else {
        currentStatusCount = 0;
      }
    } else if (stepDir < 0) {
      if (currentStatusCount > 0) {
        currentStatusCount -= 1;
      } else {
        currentStatusCount = studentStatusSeparateByTime.studentStatusAttributes.length - 1;
      }
    }

    _changeCubeState();
    first = false;
    controller.forward(from: 0.0);
  }

  List<Widget> _buildPageStatus(int index) {
    var list = <Widget>[];

    var page = studentStatusSeparate.studentStatusAttributesWithPage[index];
    if (page.length == 0) {
      return list;
    }

    list.add(Container(
      height: 24,
      decoration: BoxDecoration(color: Colors.grey[100]),
      padding: const EdgeInsets.symmetric(horizontal: 24.0),
      alignment: Alignment.centerLeft,
      child: Text("第 ${index + 1} 页", style: TextStyle(fontWeight: FontWeight.bold),),
    ));

    for (var i = 0; i < page.length; i++) {
      var _date = DateTime.fromMillisecondsSinceEpoch(page[i].updateTime * 1000);

      list.add(Container(
        height: 48,
        margin: const EdgeInsets.only(top: 8.0),
        child: new Row(
          crossAxisAlignment: CrossAxisAlignment.center,
          children: <Widget>[
            Expanded(child: Icon(IconFontCN.Emotion[page[i].attributes.emotion.getEmotion()], size: 30,), flex: 1,),
            Expanded(child: _buildStatusCards(page[i].attributes), flex: 3,),
            Expanded(child: Text("${_date.hour}:${_date.minute}:${_date.second}"), flex: 1,),
          ],
        ),
      ));
    }

    return list;
  }

  Widget _buildStatusTag(String text, Color color) {
    return Container(
      padding: const EdgeInsets.all(4.0),
      margin: const EdgeInsets.only(right: 6.0),
      decoration: BoxDecoration(color: color, borderRadius: BorderRadius.all(Radius.circular(4.0))),
      child: Text(text, style: TextStyle(fontSize: 13.0, color: Colors.white),),
    );
  }

  Row _buildStatusCards(Attributes attributes) {
    var list = <Widget>[];

    list.add(_buildStatusTag(EMOTION_TEXT[attributes.emotion.getEmotion()], Colors.blue[400]));
    if (attributes.headPose.pitchAngle > PITCH_ANGLE) {
      list.add(_buildStatusTag("低头", Colors.amber[700]));
    }
    if (attributes.eyesStatus.leftEyeStatus.noGlassEyeClose > EYE_CLOSE ||
        attributes.eyesStatus.leftEyeStatus.normalGlassEyeClose > EYE_CLOSE ||
        attributes.eyesStatus.rightEyeStatus.noGlassEyeClose > EYE_CLOSE ||
        attributes.eyesStatus.rightEyeStatus.normalGlassEyeClose > EYE_CLOSE) {
      list.add(_buildStatusTag("闭眼", Colors.red[300]));
    }

    return Row(
      children: list,
    );
  }

  Widget _buildPageStatusList() {
    return Card(
      margin: const EdgeInsets.symmetric(vertical: 4.0, horizontal: 16.0),
      child: Container(
        padding: const EdgeInsets.all(8.0),
        child: ListView.builder(itemBuilder: (BuildContext context, int index) {
          return Column(
            children: _buildPageStatus(index),
          );
        }, itemCount: studentStatusSeparate.studentStatusAttributesWithPage.length,),
      ),
    );
  }
}