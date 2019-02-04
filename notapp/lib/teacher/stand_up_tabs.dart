import 'package:flutter/material.dart';
import 'package:notapp/models/json_models.dart';
import 'package:notapp/models/models.dart';
import 'package:notapp/widgets/widgets.dart';
import 'package:charts_flutter/flutter.dart' as charts;

class StudentsTab extends StatelessWidget {
  StudentsTab({this.students, this.warning});

  final List<StudentResponse> students;
  final List<int> warning;

  @override
  Widget build(BuildContext context) {
    return new Container(
      child: new ListView.builder(
        itemCount: students?.length ?? 0,
        itemBuilder: (BuildContext context, int index) {
          return ListTile(
            leading: CircleAvatar(
              backgroundImage: Image.network("http://$SERVER_ADDRESS/images/${students[index].studentImage}").image,
            ),
            title: new Text(students[index].studentName),
            trailing: new Text("警告: ${warning[index]}", style: TextStyle(color: warning[index] > 0 ? Colors.red[400] : Colors.green[400]),),
          );
        },
      ),
    );
  }
}

class FaceCountTab extends StatelessWidget {
  FaceCountTab(this.faceCountRecordResponse, this.students);

  final List<StudentResponse> students;
  final FaceCountRecordResponse faceCountRecordResponse;

  @override
  Widget build(BuildContext context) {
    if (faceCountRecordResponse != null) {
      List<PieChartData> data = [
        new PieChartData(0, faceCountRecordResponse.studentInClassCount),
        new PieChartData(1, faceCountRecordResponse.studentNotInClass.length),
      ];

      List<charts.Series> seriesList = [
        new charts.Series<PieChartData, int>(
          id: "FaceCount",
          data: data,
          domainFn: (PieChartData pcd, _) => pcd.domain,
          measureFn: (PieChartData pcd, _) => pcd.data,
        )
      ];

      List<StudentFaceCountResultCard> studentNotIn = [];
      for (var sni in faceCountRecordResponse.studentNotInClass) {
        for (var s in students) {
          if (sni == s.studentNo) {
            studentNotIn.add(new StudentFaceCountResultCard(s.studentName, s.studentImage));
            break;
          }
        }
      }

      return new Container(
        padding: const EdgeInsets.all(8),
        child: ListView(
          children: <Widget>[
            new Container(
              height: 260,
              child: new Stack(
                alignment: Alignment.center,
                children: <Widget>[
                  new DonutPieChart(seriesList),
                  new Text("${faceCountRecordResponse.studentInClassCount} / ${faceCountRecordResponse.studentCount}",
                    style: TextStyle(color: Colors.grey[800], fontWeight: FontWeight.bold, fontSize: 18),),
                ],
              ),
            ),
            new Container(
              alignment: Alignment.center,
              padding: const EdgeInsets.all(12.0),
              margin: const EdgeInsets.symmetric(vertical: 4.0),
              decoration: new BoxDecoration(color: Colors.grey[200], borderRadius: new BorderRadius.circular(8.0)),
              child: new Text("未到学生", style: TextStyle(fontSize: 16, fontWeight: FontWeight.bold),),
            ),
            new Wrap(
              alignment: WrapAlignment.center,
              children: studentNotIn,
            )
          ],
        )
      );
    } else {
      return new Container();
    }
  }
}

class StudentStatusTab extends StatelessWidget {
  StudentStatusTab(this.studentWarningRecord);

  final List<String> studentWarningRecord;

  @override
  Widget build(BuildContext context) {
    return ListView.builder(
      itemCount: studentWarningRecord?.length ?? 0,
      itemBuilder: (BuildContext context, int index) {
        return new Card(
          margin: const EdgeInsets.fromLTRB(8.0, 4.0, 8.0, 4.0),
          child: new Text(studentWarningRecord[index]),
        );
      },
    );
  }
}

class PDFTab extends StatelessWidget {
  PDFTab(this.pageCallBack, this.page);

  final PageCallBack pageCallBack;
  final int page;

  @override
  Widget build(BuildContext context) {
    return new Column(
      mainAxisAlignment: MainAxisAlignment.center,
      crossAxisAlignment: CrossAxisAlignment.center,
      children: <Widget>[
        new Container(
          height: 150,
          width: 200,
          child: new Text(page.toString(), style: TextStyle(fontSize: 18, color: Colors.white, fontWeight: FontWeight.bold),),
        ),
        new Row(
          mainAxisAlignment: MainAxisAlignment.spaceAround,
          children: <Widget>[
            new IconButton(icon: Icon(Icons.navigate_before), onPressed: () => pageCallBack(1)),
            new IconButton(icon: Icon(Icons.navigate_next), onPressed: () => pageCallBack(-1))
          ],
        ),
      ],
    );
  }
}

typedef PageCallBack = Function(int page);