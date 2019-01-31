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
  FaceCountTab(this.faceCountRecordResponse);

  final FaceCountRecordResponse faceCountRecordResponse;

  @override
  Widget build(BuildContext context) {
    if (faceCountRecordResponse != null) {
      List<PieChartData> data = [
        new PieChartData(0, faceCountRecordResponse.studentInClassCount),
        new PieChartData(0, faceCountRecordResponse.studentCount - faceCountRecordResponse.studentInClassCount),
      ];

      List<charts.Series> seriesList = [
        new charts.Series<PieChartData, int>(
          id: "FaceCount",
          data: data,
          domainFn: (PieChartData pcd, _) => pcd.domain,
          measureFn: (PieChartData pcd, _) => pcd.data,
        )
      ];

      return new Container(
        padding: const EdgeInsets.all(8),
        child: ListView(
          children: <Widget>[
            new Container(
              padding: const EdgeInsets.all(8),
              child: new Stack(
                alignment: Alignment.center,
                children: <Widget>[
                  new DonutPieChart(seriesList),
                  new Text("${faceCountRecordResponse.studentInClassCount}/${faceCountRecordResponse.studentCount}",
                    style: TextStyle(color: Colors.grey[800], fontWeight: FontWeight.bold),),
                ],
              ),
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
  @override
  Widget build(BuildContext context) {
    // TODO: implement build
    return new Text("pdf");
  }
}