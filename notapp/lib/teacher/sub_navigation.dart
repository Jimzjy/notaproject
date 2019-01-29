import 'package:flutter/material.dart';
import 'package:notapp/models/json_models.dart';
import 'package:notapp/widgets/widgets.dart';
import 'stand_up_tabs.dart';
import 'package:web_socket_channel/io.dart';
import 'dart:convert';

class StandUpClassPage extends StatefulWidget {
  StandUpClassPage({
    this.classID,
    this.wReadMWriteIndex,
    this.classResponse,
  });

  final int classID;
  final int wReadMWriteIndex;
  final ClassResponse classResponse;

  @override
  State<StatefulWidget> createState() => _StandUpClassPageState(
    classID: classID,
    wReadMWriteIndex: wReadMWriteIndex,
    classResponse: classResponse,
  );
}

class _StandUpClassPageState extends State<StandUpClassPage> {
  _StandUpClassPageState({
    this.classID,
    this.wReadMWriteIndex,
    this.classResponse,
  });

  int classID;
  int wReadMWriteIndex;
  ClassResponse classResponse;

  static const _tabs = <Tab>[
    const Tab(text: "学生信息",),
    const Tab(text: "点名信息",),
    const Tab(text: "学生状态",),
    const Tab(text: "演示文稿",),
  ];

  final _tabContents = <Widget>[
    new StudentsTab(),
    new FaceCountTab(),
    new StudentStatusTab(),
    new PDFPage(),
  ];

  @override
  Widget build(BuildContext context) {
    return new Scaffold(
      appBar: new AppBar(
        title: new Text(classResponse.className),
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
              child: new TabBarView(children: _tabContents),
            ),
          ],
        ),
      ),
    );
  }
}

class SettingPage extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    // TODO: implement build
    return new Text("setting");
  }
}