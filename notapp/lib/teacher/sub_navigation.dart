import 'package:flutter/material.dart';
import 'package:notapp/models/json_models.dart';
import 'package:notapp/widgets/widgets.dart';
import 'stand_up_tabs.dart';
import 'package:dio/dio.dart';
import 'package:notapp/models/models.dart';
import 'package:fluttertoast/fluttertoast.dart';
import 'package:web_socket_channel/io.dart';
import 'dart:convert';
import 'package:flushbar/flushbar.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:notapp/main.dart' show NOT_LOGIN;
import 'package:notapp/login/login_main.dart';
import 'dart:math';

const PITCH_ANGLE = 30;
const EYE_CLOSE = 50;

class StandUpClassPage extends StatefulWidget {
  StandUpClassPage({
    this.wsParam,
    this.classResponse,
  });

  final String wsParam;
  final ClassResponse classResponse;

  @override
  State<StatefulWidget> createState() => _StandUpClassPageState(
    wsParam: wsParam,
    classResponse: classResponse,
  );
}

class _StandUpClassPageState extends State<StandUpClassPage> {
  _StandUpClassPageState({
    this.wsParam,
    this.classResponse,
  });

  String wsParam;
  ClassResponse classResponse;
  List<StudentResponse> studentsResponse;
  List<int> studentWarning;
  List<StudentWarningRecordTime> studentWarningRecord = [];
  IOWebSocketChannel _channel;
  FaceCountRecordResponse faceCountRecordResponse;
  int page = 0;

  static const _tabs = <Tab>[
    const Tab(text: "学生信息",),
    const Tab(text: "点名信息",),
    const Tab(text: "学生状态",),
    const Tab(text: "演示文稿",),
  ];

  @override
  void initState() {
    _requestStudents().then((v){
      _wsConnect();
    });
    super.initState();
  }

  @override
  void dispose() {
    if (_channel != null) {
      _channel.sink.close();
    }
    _channel = null;
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return new Scaffold(
      appBar: new AppBar(
        title: new Text("正在上课: ${classResponse.className}", overflow: TextOverflow.ellipsis,),
        elevation: 0.0,
      ),
      body: new NestedScrollView(
        headerSliverBuilder: (BuildContext context, bool innerBoxIsScrolled) {
          return <Widget>[
            new SliverList(delegate: new SliverChildListDelegate([
              buildClassCard(classResponse, context)
            ]))
          ];
        },
        body: _buildTabs(),
      ),
    );
  }

  Widget _buildTabs() {
    return new Padding(
      padding: const EdgeInsets.symmetric(horizontal: 8.0),
      child: new DefaultTabController(
        length: _tabs.length,
        child: new Column(
          children: <Widget>[
            new TabBar(tabs: _tabs, isScrollable: false, labelColor: Theme.of(context).primaryColor,),
            new Expanded(
              child: new TabBarView(children: _tabs.map((tab) {
                switch(tab.text) {
                  case "学生信息":
                    return new StudentsTab(
                      students: studentsResponse,
                      warning: studentWarning,
                    );
                  case "点名信息":
                    return new FaceCountTab(faceCountRecordResponse, studentsResponse);
                  case "学生状态":
                    return new StudentStatusTab(studentWarningRecord);
                  case "演示文稿":
                    return new PDFTab((int p) {
                      _channel?.sink?.add(jsonEncode(StandUpPacket(changePDFPage: p)));
                    }, page);
                }
              }).toList()),
            ),
          ],
        ),
      ),
    );
  }

  Future<void> _requestStudents() async {
    Response response;
    try {
      response = await DioManager.instance.get("/students", data: { "class_id": classResponse.classID });
    } catch(e) {
      print(e);
    }

    StudentsResponse stuResp;
    if (response?.statusCode == 200) {
      Map jsonMap = jsonDecode(response.toString());
      stuResp = StudentsResponse.fromJson(jsonMap);
    } else {
      Fluttertoast.showToast(
        msg: "获取学生信息失败",
        toastLength: Toast.LENGTH_SHORT,
        timeInSecForIos: 1,
        gravity: ToastGravity.CENTER,
      );
      Navigator.pop(context);
    }

    setState(() {
      studentsResponse = stuResp.students;
      studentWarning = List.filled(studentsResponse.length, 0);
    });
    return;
  }

  _requestFaceCountRecord(int recordID) async {
    Response response;
    try {
      response = await DioManager.instance.get("/face_count_record", data: { "record_id": recordID });
    } catch(e) {
      print(e);
    }

    FaceCountRecordResponse faceCRR;
    if (response?.statusCode == 200) {
      Map jsonMap = jsonDecode(response.toString());
      faceCRR = FaceCountRecordResponse.fromJson(jsonMap);

      setState(() {
        faceCountRecordResponse = faceCRR;
      });
    }
  }

  _wsConnect() async {
    try {
      _channel = new IOWebSocketChannel.connect("ws://" + SERVER_ADDRESS + wsParam);
      _channel.stream.listen((message) {
        Map jsonMap = jsonDecode(message.toString());
        StandUpPacket sup = StandUpPacket.fromJson(jsonMap);

        setState(() {
          _handleStudentWarringChange(sup);
        });

        _handleFaceCount(sup);
        _handleStudentWarringNotification(sup);
        _handlePageChange(sup);
        _handleSayGoodbye(sup);
      });
    } catch(e) {
      print(e);
    }
  }

  _handleStudentWarringChange(StandUpPacket sup) {
    if (sup.studentWarningRecordList != null && sup.studentWarningRecordList.length == studentWarning.length) {
        studentWarning = sup.studentWarningRecordList;
    }

    if (sup.studentWarningList != null && sup.studentWarningList.length > 0) {
      studentWarningRecord.add(StudentWarningRecordTime(sup.studentWarningList, DateTime.now()));
    }
  }

  _handleStudentWarringNotification(StandUpPacket sup) {
    if (sup.studentWarningList != null && sup.studentWarningList.length > 0) {
      Flushbar(flushbarPosition: FlushbarPosition.TOP,)
        ..title = "以下学生可能未认真听课(学号): "
        ..message = sup.studentWarningList
        ..icon = Icon(Icons.error_outline, color: Theme.of(context).primaryColor, size: 28,)
        ..duration = Duration(seconds: 4)
        ..show(context);
    }
  }

  _handleFaceCount(StandUpPacket sup) {
    if (sup.faceCountClose) {
      _requestFaceCountRecord(sup.faceCountRecordID);
    }
  }

  _handlePageChange(StandUpPacket sup) {
    if (sup.currentPDFPage > 0) {
      page = sup.currentPDFPage;
    }
  }

  _handleSayGoodbye(StandUpPacket sup) {
    if (sup.sayGoodbye) {
      if (_channel != null) {
        _channel.sink.close();
      }
      _channel = null;

      _goodByeDialog();
    }
  }

  Future<void> _goodByeDialog() async {
    return showDialog<void>(
      context: context,
      barrierDismissible: false,
      builder: (BuildContext context) {
        return AlertDialog(
          content: Text('已下课, 是否留在当前页面?'),
          actions: <Widget>[
            FlatButton(
              child: Text('否'),
              onPressed: () {
                Navigator.pop(context);
                Navigator.pop(this.context);
              },
            ),
            FlatButton(
              child: Text('是'),
              onPressed: () {
                Navigator.pop(context);
              },
            ),
          ],
        );
      },
    );
  }
}

class SettingPage extends StatelessWidget {
  Dio dio = DioManager.instance;

  @override
  Widget build(BuildContext context) {
    return new Scaffold(
      appBar: new AppBar(
        title: new Text("设置"),
        elevation: 0.0,
      ),
      body: new Container(
        decoration: BoxDecoration(
          color: Theme.of(context).primaryColor
        ),
        child: new Center(
          child: new Column(
            mainAxisAlignment: MainAxisAlignment.center,
            crossAxisAlignment: CrossAxisAlignment.center,
            children: <Widget>[
              new MaterialButton(
                color: Colors.white,
                onPressed: () => _clearError(),
                child: new Text("清除错误"),
              ),
              new MaterialButton(
                color: Colors.white,
                onPressed: () => _onLogoutButtonPressed(context),
                child: new Text("退出登录"),
              ),
            ],
          ),
        )
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

  _clearError() async {
    Response response;

    try {
      response = await dio.get("/clear");
    } catch (err) {
      print(err);
    }
  }
}

class NormalClassPage extends StatefulWidget {
  NormalClassPage({this.studentStatusResponse});

  final StudentStatusResponse studentStatusResponse;

  @override
  State<StatefulWidget> createState() {
    return _NormalClassPageState(studentStatusResponse: studentStatusResponse);
  }
}

class _NormalClassPageState extends State<NormalClassPage> {
  _NormalClassPageState({this.studentStatusResponse});

  StudentStatusResponse studentStatusResponse;
  ClassResponse classResponse;
  FaceCountRecordResponse faceCountRecordResponse;
  List<StudentResponse> studentsResponse;
  List<String> studentNoSort;
  List<StudentStatusSeparate> studentStatusSeparateList = [];
  List<StudentStatusSeparateByTime> studentStatusSeparateListByTime = [];

  static const _tabs = <Tab>[
    const Tab(text: "学生信息",),
    const Tab(text: "点名信息",),
  ];

  @override
  void initState() {
    _requestClass(studentStatusResponse.classID).then((_) {
      _requestFaceCountRecord(studentStatusResponse.faceCountRecordID);
      _requestStudents(classResponse.classID).then((_studentResponse) {
        if (_studentResponse == null) {
          return;
        }

        studentNoSort = [];
        _studentResponse.forEach((StudentResponse _stuResp) {
          studentNoSort.add(_stuResp.studentNo);

          List<List<StudentStatusAttributes>> studentStatusAttributesWithPage = [];
          for (var i = 0; i < studentStatusResponse.pdfPageCount; i++) {
            studentStatusAttributesWithPage.add([]);
          }
          studentStatusSeparateList.add(StudentStatusSeparate(_stuResp.studentNo, studentStatusAttributesWithPage));
          studentStatusSeparateListByTime.add(StudentStatusSeparateByTime(_stuResp.studentNo, []));
        });

        studentStatusResponse.studentStatus.forEach((StudentStatusWithPage _studentStatusWithPage) {
          int _page = _studentStatusWithPage.pdfPage - 1;

          _studentStatusWithPage.studentStatus?.forEach((StudentStatus _studentStatus) {
            int _index = studentNoSort.indexOf(_studentStatus.studentNo);

            studentStatusSeparateList[_index]
                .studentStatusAttributesWithPage[_page]
                .add(StudentStatusAttributes(_studentStatus.updateTime, _studentStatus.attributes));
            studentStatusSeparateListByTime[_index]
                .studentStatusAttributes
                .add(StudentStatusAttributes(_studentStatus.updateTime, _studentStatus.attributes));
          });
        });
      });
    });
    super.initState();
  }

  @override
  Widget build(BuildContext context) {
    return new Scaffold(
      appBar: new AppBar(
        title: new Text(studentStatusResponse.className, overflow: TextOverflow.ellipsis,),
        elevation: 0.0,
      ),
      body: new NestedScrollView(
        headerSliverBuilder: (BuildContext context, bool innerBoxIsScrolled) {
          return <Widget>[
            new SliverList(delegate: new SliverChildListDelegate([
              buildClassCard(classResponse, context)
            ]))
          ];
        },
        body: _buildTabs(),
      ),
    );
  }

  Widget _buildTabs() {
    return new Padding(
      padding: const EdgeInsets.symmetric(horizontal: 8.0),
      child: new DefaultTabController(
        length: _tabs.length,
        child: new Column(
          children: <Widget>[
            new TabBar(tabs: _tabs, isScrollable: false, labelColor: Theme.of(context).primaryColor,),
            new Expanded(
              child: new TabBarView(children: _tabs.map((tab) {
                switch(tab.text) {
                  case "学生信息":
                    return new RecordStudentsTab(
                      students: studentsResponse,
                      warningRecord: studentStatusResponse.studentWarningRecordList,
                      tapCallBack: (int index) {
                        _navigateToStudentStatusPage(studentStatusSeparateList[index], studentsResponse[index], studentStatusSeparateListByTime[index]);
                      },
                    );
                  case "点名信息":
                    return new FaceCountTab(faceCountRecordResponse, studentsResponse);
                }
              }).toList()),
            ),
          ],
        ),
      ),
    );
  }

  _navigateToStudentStatusPage(StudentStatusSeparate studentStatusSeparate, StudentResponse studentResponse, StudentStatusSeparateByTime studentStatusSeparateByTime) {
    Navigator.push(context, new MaterialPageRoute(builder: (context) {
      return new StudentStatusPage(studentStatusSeparate, studentResponse, studentStatusSeparateByTime);
    }));
  }

  _requestClass(int classID) async {
    Response response;
    try {
      response = await DioManager.instance.get("/classes", data: { "class_id": classID });
    } catch(e) {
      print(e);
    }

    ClassesResponse clResp;
    if (response?.statusCode == 200) {
      Map jsonMap = jsonDecode(response.toString());
      clResp = ClassesResponse.fromJson(jsonMap);
      setState(() {
        classResponse = clResp.classes[0];
      });
    } else {
      Fluttertoast.showToast(
        msg: "获取学生信息失败",
        toastLength: Toast.LENGTH_SHORT,
        timeInSecForIos: 1,
        gravity: ToastGravity.CENTER,
      );
      Navigator.pop(context);
    }
  }

  Future<List<StudentResponse>> _requestStudents(int classID) async {
    Response response;
    try {
      response = await DioManager.instance.get("/students", data: { "class_id": classID });
    } catch(e) {
      print(e);
    }

    StudentsResponse stuResp;
    if (response?.statusCode == 200) {
      Map jsonMap = jsonDecode(response.toString());
      stuResp = StudentsResponse.fromJson(jsonMap);
      setState(() {
        studentsResponse = stuResp.students;
      });

      return stuResp.students;
    }
    return null;
  }

  _requestFaceCountRecord(int recordID) async {
    Response response;
    try {
      response = await DioManager.instance.get("/face_count_record", data: { "record_id": recordID });
    } catch(e) {
      print(e);
    }

    FaceCountRecordResponse faceCRR;
    if (response?.statusCode == 200) {
      Map jsonMap = jsonDecode(response.toString());
      faceCRR = FaceCountRecordResponse.fromJson(jsonMap);

      setState(() {
        faceCountRecordResponse = faceCRR;
      });
    }
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
                _navigateToClassPage(context, index);
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

  _navigateToClassPage(BuildContext context, int index) {
    Navigator.push(context, new MaterialPageRoute(
      builder: (context) {
        return new NormalClassPage(
          studentStatusResponse: _studentStatus[index],
        );
      },
    ));
  }
}

Widget buildClassCard(ClassResponse classResponse, BuildContext context) {
  if (classResponse == null) {
    return new Container();
  }

  return new Container(
    decoration: new BoxDecoration(
      color: Theme.of(context).primaryColor,
    ),
    padding: const EdgeInsets.fromLTRB(0, 4.0, 0, 8.0),
    child: new ClassCard(
      className: classResponse.className,
      classroomNo: classResponse.classroomNo,
      classImage: classResponse.classImage,
      classID: classResponse.classID,
    ),
  );
}

class StudentStatusPage extends StatefulWidget {
  StudentStatusPage(this.studentStatusSeparate, this.studentResponse, this.studentStatusSeparateByTime);

  final StudentStatusSeparate studentStatusSeparate;
  final StudentResponse studentResponse;
  final StudentStatusSeparateByTime studentStatusSeparateByTime;

  @override
  State<StatefulWidget> createState() => _StudentStatusPageState(studentStatusSeparate, studentResponse, studentStatusSeparateByTime);
}

class _StudentStatusPageState extends State<StudentStatusPage> with SingleTickerProviderStateMixin {
  _StudentStatusPageState(this.studentStatusSeparate, this.studentResponse, this.studentStatusSeparateByTime);

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
    controller = new AnimationController(
        duration: const Duration(milliseconds: 500), vsync: this);
    animation = CurvedAnimation(parent: controller, curve: Curves.easeInOut);
    animation.addListener(() {
      setState(() {});
    });

    if (studentStatusSeparateByTime.studentStatusAttributes.length <= currentStatusCount) {
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
    super.initState();
  }

  @override
  void dispose() {
    controller.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
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

const EMOTION_TEXT = [
  "伤心",
  "平常",
  "厌恶",
  "愤怒",
  "惊讶",
  "害怕",
  "开心",
];