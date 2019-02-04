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
  List<String> studentWarningRecord;
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
              _buildClassCard()
            ]))
          ];
        },
        body: _buildTabs(),
      ),
    );
  }
  
  Widget _buildClassCard() {
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
        print(message);

        Map jsonMap = jsonDecode(message.toString());
        StandUpPacket sup = StandUpPacket.fromJson(jsonMap);

        setState(() {
          _handleStudentWarringChange(sup);
        });

        _handleFaceCount(sup);
        _handleStudentWarringNotification(sup);
        _handlePageChange(sup);
      });
    } catch(e) {
      print(e);
    }
  }

  _handleStudentWarringChange(StandUpPacket sup) {
    if (sup.studentWarningRecordList != null && sup.studentWarningRecordList.length == studentWarning.length) {
        studentWarning = sup.studentWarningRecordList;
    }

    if (sup.studentWarringList != null && sup.studentWarringList.length > 0) {
      studentWarningRecord.add(sup.studentWarringList);
    }
  }

  _handleStudentWarringNotification(StandUpPacket sup) {
    if (sup.studentWarringList != null && sup.studentWarringList.length > 0) {
      Flushbar()
        ..title = "以下学生可能未认真听课(学号): "
        ..message = sup.studentWarringList
        ..icon = Icon(Icons.error_outline, color: Theme.of(context).primaryColor, size: 28,)
        ..duration = Duration(seconds: 3)
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
}

class SettingPage extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return new Scaffold(
      body: new Text("setting"),
    );
  }
}

class NormalClassPage extends StatefulWidget {
  NormalClassPage({this.classResponse});

  final ClassResponse classResponse;

  @override
  State<StatefulWidget> createState() {
    return _NormalClassPageState(classResponse: this.classResponse);
  }
}

class _NormalClassPageState extends State<NormalClassPage> {
  _NormalClassPageState({this.classResponse});

  ClassResponse classResponse;

  @override
  Widget build(BuildContext context) {
    return new Scaffold(
      appBar: new AppBar(
        title: new Text(classResponse.className, overflow: TextOverflow.ellipsis,),
        elevation: 0.0,
      ),
      body: new NestedScrollView(
        headerSliverBuilder: (BuildContext context, bool innerBoxIsScrolled) {
          return <Widget>[
            new SliverList(delegate: new SliverChildListDelegate([
              _buildClassCard()
            ]))
          ];
        },
        body: _buildTabs(),
      ),
    );
  }
}