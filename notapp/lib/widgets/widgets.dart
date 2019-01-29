import 'package:flutter/material.dart';
import 'package:notapp/models/models.dart';

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

