import 'package:flutter/material.dart';

class NormalPage extends StatefulWidget {
  NormalPage({Key key}) : super(key: key);

  @override
  State<StatefulWidget> createState() => _NormalPageState();
}

class _NormalPageState extends State<NormalPage> {

  @override
  Widget build(BuildContext context) {
    return new Scaffold(
      appBar: new AppBar(),
    );
  }
}