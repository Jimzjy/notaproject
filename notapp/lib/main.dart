import 'package:flutter/material.dart';
import 'package:web_socket_channel/io.dart';
import 'package:qrcode_reader/qrcode_reader.dart';
import 'package:json_annotation/json_annotation.dart';
import 'dart:convert';

part 'main.g.dart';

void main() => runApp(new MyApp());

class MyApp extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return new MaterialApp(
      title: 'NOTAPP',
      theme: new ThemeData(
        primarySwatch: Colors.blue,
      ),
      home: new WebSocketRoute(),
    );
  }
}

class WebSocketRoute extends StatefulWidget {
  @override
  _WebSocketRouteState createState() => new _WebSocketRouteState();
}

class _WebSocketRouteState extends State<WebSocketRoute> {
  IOWebSocketChannel _channel;
  String _text = "";

  @override
  Widget build(BuildContext context) {
    return new Scaffold(
      appBar: new AppBar(
        actions: <Widget>[
          new IconButton(icon: new Icon(Icons.camera_alt), onPressed: _onScanPressed)
        ],
      ),
      body: new Padding(
        padding: const EdgeInsets.all(8.0),
        child: new Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: <Widget>[
            new Expanded(child: new SingleChildScrollView(
              child: _buildStreamBuild(),
            )),
            new Padding(padding: const EdgeInsets.all(24.0),
              child: new Row(
                mainAxisAlignment: MainAxisAlignment.spaceEvenly,
                children: <Widget>[
                  new RaisedButton(onPressed: _onBackButtonPressed, child: new Text("上一页"), color: Theme.of(context).primaryColor,),
                  new RaisedButton(onPressed: _onNextButtonPressed, child: new Text("下一页"), color: Theme.of(context).primaryColor,),
                ],
              ),
            )
          ],
        ),
      ),
    );
  }

  @override
  void dispose() {
    if (_channel != null) {
      _channel.sink.close();
    }
    _channel = null;
    super.dispose();
  }

  _onScanPressed() {
    Future<String> futureString = new QRCodeReader().scan();
    futureString.then((String value) {
      _connectToServer(value);
    });
  }

  _connectToServer(String param) {
    if (_channel != null) {
      return;
    }

    _channel = new IOWebSocketChannel.connect('ws://192.168.50.188:8000' + param);
    setState(() {});
  }

  _onNextButtonPressed() {
    _channel.sink.add(jsonEncode(StandUpPacket(changePDFPage: 1)));
  }

  _onBackButtonPressed() {
    _channel.sink.add(jsonEncode(StandUpPacket(changePDFPage: -1)));
  }

  Widget _buildStreamBuild() {
    if (_channel == null) {
      return new Text(_text);
    }

    return new StreamBuilder(
      stream: _channel.stream,
      builder: (context, snapshot) {
        if (snapshot.hasError) {
          _text = "连接失败...";
        } else if (snapshot.hasData) {
          _text += _analyzeData(snapshot.data + "");
        }
        return new Padding(
          padding: const EdgeInsets.symmetric(vertical: 24.0),
          child: new Text(_text),
        );
      },
    );
  }
  
  String _analyzeData(String json) {
    String text = "";
    Map jsonMap = jsonDecode(json);
    StandUpPacket sup = StandUpPacket.fromJson(jsonMap);

    if (sup.currentPDFPage > 0) {
      text += "切换到第 ${sup.currentPDFPage} 页\n";
    }

    if (sup.faceCountClose) {
      text += "点名完成!\n";
    }

    return text;
  }
}

@JsonSerializable()
class StandUpPacket {
  StandUpPacket({
    this.wReadMWriteIndex = 0,
    this.wWriteMReadIndex = 0,
    this.faceCountClose = false,
    this.currentPDFPage = 0,
    this.changePDFPage = 0,
  });

  @JsonKey(name: 'WReadMWriteIndex')
  int wReadMWriteIndex;

  @JsonKey(name: 'WWriteMReadIndex')
  int wWriteMReadIndex;

  @JsonKey(name: 'FaceCountClose')
  bool faceCountClose;

  @JsonKey(name: 'CurrentPDFPage')
  int currentPDFPage;

  @JsonKey(name: 'ChangePDFPage')
  int changePDFPage;

  factory StandUpPacket.fromJson(Map<String, dynamic> json) => _$StandUpPacketFromJson(json);
  Map<String, dynamic> toJson() => _$StandUpPacketToJson(this);
}