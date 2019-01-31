import 'package:flutter/material.dart';
import 'package:notapp/models/models.dart';
import 'package:notapp/models/json_models.dart';
import 'package:charts_flutter/flutter.dart' as charts;

class FakeSearchBar extends StatelessWidget implements PreferredSizeWidget {
  FakeSearchBar({this.left, this.right, this.text, this.onSearchBarClick});

  final Widget left;
  final Widget right;
  final String text;
  final Function onSearchBarClick;

  @override
  Size get preferredSize => Size.fromHeight(kToolbarHeight);

  @override
  Widget build(BuildContext context) {
    return new Container(
      color: Theme.of(context).primaryColor,
      child: new SafeArea(
        child: new Row(
          children: _buildBar(context),
        ),
      ),
    );
  }

  List<Widget> _buildBar(BuildContext context) {
    List<Widget> list;
    if (left == null && right == null) {
      list = [
        _buildSearchArea(EdgeInsets.symmetric(vertical: 10.0, horizontal: 24.0)),
      ];
    } else if (right == null) {
      list = [
        _buildLeftRightArea(left, context),
        _buildSearchArea(EdgeInsets.fromLTRB(0.0, 10.0, 24.0, 10.0)),
      ];
    } else if (left == null) {
      list = [
        _buildSearchArea(EdgeInsets.fromLTRB(24.0, 10.0, 0.0, 10.0)),
        _buildLeftRightArea(right, context),
      ];
    } else {
      list = [
        _buildLeftRightArea(left, context),
        _buildSearchArea(EdgeInsets.symmetric(vertical: 10.0)),
        _buildLeftRightArea(right, context),
      ];
    }
    return list;
  }

  Widget _buildLeftRightArea(Widget location, BuildContext context) {
    return Expanded(child: Material(child: location, color: Theme.of(context).primaryColor,));
  }

  Widget _buildSearchArea(EdgeInsetsGeometry margin) {
    return Expanded(
      child: new GestureDetector(
        child: new Container(
          decoration: new BoxDecoration(
            color: Colors.white30,
            borderRadius: BorderRadius.circular(8.0),
          ),
          margin: margin,
          child: new SizedBox.expand(
            child: new Container(
              alignment: Alignment.centerLeft,
              padding: EdgeInsets.only(left: 8.0),
              child: new Row(
                children: <Widget>[
                  new Icon(Icons.search, color: Colors.black26,),
                  new Text(text, style: TextStyle(fontSize: 16.0, color: Colors.black26,)),
                ],
              ),
            ),
          ),
        ),
        onTap: this.onSearchBarClick,
      ),
      flex: 4,
    );
  }
}

class ClassCard extends StatelessWidget {
  ClassCard({ this.className, this.classID, this.classImage, this.classroomNo, this.itemPressesCallback });

  final String className;
  final int classID;
  final String classImage;
  final String classroomNo;
  final GestureTapCallback itemPressesCallback;

  @override
  Widget build(BuildContext context) {
    return new GestureDetector(
      child: new Container(
        height: 96.0,
        margin: const EdgeInsets.symmetric(horizontal: 16.0, vertical: 8.0),
        padding: const EdgeInsets.all(12.0),
        decoration: new BoxDecoration(borderRadius: BorderRadius.circular(8.0), color: Colors.white),
        child: new Row(
          children: <Widget>[
            new Expanded(child: new ClipRRect(
              child: new Image.network("http://$SERVER_ADDRESS/images/$classImage"),
              borderRadius: BorderRadius.circular(12.0),
            ), flex: 1,),
            new Expanded(
              child: new Container(
                margin: const EdgeInsets.only(left: 8.0),
                padding: const EdgeInsets.symmetric(horizontal: 12.0),
                decoration: new BoxDecoration(
                  color: Colors.grey[100],
                  borderRadius: new BorderRadius.circular(8.0),
                ),
                child: new Column(
                  mainAxisAlignment: MainAxisAlignment.center,
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: <Widget>[
                    new Text(className, style: new TextStyle(fontWeight: FontWeight.bold, fontSize: 16.0),),
                    new Padding(padding: const EdgeInsets.only(top: 4.0),
                      child: new Text("教室: $classroomNo", style: new TextStyle(fontSize: 12.0),),
                    ),
                    new Padding(padding: const EdgeInsets.only(top: 4.0),
                      child: new Text("ID: $classID", style: new TextStyle(fontSize: 12.0),),
                    ),
                  ],
                ),
              ),
              flex: 3,),
          ],
        ),
      ),
      onTap: itemPressesCallback,
    );
  }
}

class HistoryCard extends StatelessWidget {
  HistoryCard({this.studentStatusResponse, this.itemPressesCallback});

  final StudentStatusResponse studentStatusResponse;
  final GestureTapCallback itemPressesCallback;

  @override
  Widget build(BuildContext context) {
    return GestureDetector(
      child: new Container(
        height: 88,
        margin: const EdgeInsets.symmetric(horizontal: 16.0, vertical: 4.0),
        padding: const EdgeInsets.all(8.0),
        decoration: new BoxDecoration(borderRadius: BorderRadius.circular(8.0), color: Colors.white),
        child: new Container(
          decoration: new BoxDecoration(
            color: Colors.grey[100],
            borderRadius: new BorderRadius.circular(8.0),
          ),
          padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
          child: new Column(
            children: <Widget>[
              new Expanded(child: new Text(studentStatusResponse.className, style: TextStyle(fontSize: 16.0, fontWeight: FontWeight.bold),)),
              new Expanded(child: new Text("Class ID: " + studentStatusResponse.classID.toString(), style: TextStyle(fontSize: 14.0),)),
              new Row(
                mainAxisAlignment: MainAxisAlignment.end,
                children: <Widget>[
                  new Text(
                    _getTimeString(new DateTime.fromMillisecondsSinceEpoch(studentStatusResponse.updateTime * 1000)),
                    style: TextStyle(fontSize: 12.0),
                  ),
                ],
              )
            ],
            crossAxisAlignment: CrossAxisAlignment.start,
          ),
        ),
      ),
      onTap: itemPressesCallback,
    );
  }

  String _getTimeString(DateTime date) {
    return "${date.month} 月 ${date.day} 日  周${_numToCS(date.weekday)}  ${date.hour}:${date.minute}";
  }

  String _numToCS(int n) {
    String s = "";

    switch (n) {
      case 1:
        s = "一";
        break;
      case 2:
        s = "二";
        break;
      case 3:
        s = "三";
        break;
      case 4:
        s = "四";
        break;
      case 5:
        s = "五";
        break;
      case 6:
        s = "六";
        break;
      case 7:
        s = "日";
        break;
    }
    return s;
  }
}

// https://google.github.io/charts/flutter/example/pie_charts/donut
class DonutPieChart extends StatelessWidget {
  final List<charts.Series> seriesList;
  final bool animate;

  DonutPieChart(this.seriesList, {this.animate = false});

  @override
  Widget build(BuildContext context) {
    return new charts.PieChart(seriesList,
        animate: animate,
        // Configure the width of the pie slices to 60px. The remaining space in
        // the chart will be left as a hole in the center.
        defaultRenderer: new charts.ArcRendererConfig(arcWidth: 30));
  }
}