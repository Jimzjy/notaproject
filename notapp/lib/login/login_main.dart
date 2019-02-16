import 'package:flutter/material.dart';
import 'package:dio/dio.dart';
import 'package:notapp/main.dart';
import 'package:notapp/models/models.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:notapp/normal/normal_main.dart';
import 'package:notapp/teacher/teacher_main.dart';
import 'package:fluttertoast/fluttertoast.dart';

class LoginApp extends StatefulWidget {

  @override
  State<StatefulWidget> createState() => _LoginAppState();
}

class _LoginAppState extends State<LoginApp> {
  Color primaryColor = TEACHER_APP_COLOR;

  @override
  Widget build(BuildContext context) {
    return new MaterialApp(
      title: 'NOTAPP',
      theme: new ThemeData(
          primarySwatch: primaryColor,
      ),
      home: LoginPage(primaryColorChangeCallBack: (color) {
        setState(() {
          primaryColor = color;
        });
      },),
    );
  }
}

class LoginPage extends StatefulWidget {
  LoginPage({Key key, @required this.primaryColorChangeCallBack}) : super(key: key);

  final PrimaryColorChangeCallBack primaryColorChangeCallBack;

  @override
  State<StatefulWidget> createState() => _LoginPageState(primaryColorChangeCallBack: primaryColorChangeCallBack);
}

class _LoginPageState extends State<LoginPage> {
  _LoginPageState({ @required this.primaryColorChangeCallBack });

  TextEditingController _usernameController = new TextEditingController();
  TextEditingController _pwdController = new TextEditingController();
  GlobalKey _formKey= new GlobalKey<FormState>();
  String _userType = TEACHER_USER;
  PrimaryColorChangeCallBack primaryColorChangeCallBack;


  @override
  Widget build(BuildContext context) {
    Color pColor = Theme.of(context).primaryColor;

    return new Scaffold(
      appBar: new AppBar(
        title: new Text("登录", style: TextStyle(color: Colors.white),),
      ),
      body: new ListView(
        scrollDirection: Axis.vertical,
        children: <Widget>[
          new Padding(
            padding: const EdgeInsets.symmetric(vertical: 48.0),
            child: new Image.asset("assets/logo.png", height: 120,),
          ),
          new Row(
            mainAxisAlignment: MainAxisAlignment.center,
            children: <Widget>[
              new Container(
                padding: const EdgeInsets.symmetric(horizontal: 8.0),
                margin: const EdgeInsets.symmetric(horizontal: 8.0),
                child: new ActionChip(
                  label: new Text("教师", style: TextStyle(color: _userType == TEACHER_USER ? Colors.white : pColor, fontSize: 14.0),),
                  backgroundColor: _userType == TEACHER_USER ? pColor : null,
                  onPressed: () {
                    onUserTypeSelectPressed(TEACHER_USER);
                  },
                ),
              ),
              new Container(
                padding: const EdgeInsets.symmetric(horizontal: 8.0),
                margin: const EdgeInsets.symmetric(horizontal: 8.0),
                child: new ActionChip(
                  label: new Text("普通", style: TextStyle(color: _userType == NORMAL_USER ? Colors.white : pColor, fontSize: 14.0),),
                  backgroundColor: _userType == NORMAL_USER ? pColor : null,
                  onPressed: () {
                    onUserTypeSelectPressed(NORMAL_USER);
                  },
                ),
              ),
            ],
          ),
          new Container(
            margin: const EdgeInsets.fromLTRB(24.0, 0.0, 24.0, 0.0),
            child: new Form(
              key: _formKey,
              autovalidate: true,
              child: new Column(
                children: <Widget>[
                  TextFormField(
                      controller: _usernameController,
                      decoration: InputDecoration(
                          labelText: "用户名",
                          hintText: "用户名",
                          icon: Icon(Icons.person)
                      ),
                      validator: (v) {
                        return v
                            .trim()
                            .length > 0 ? null : "用户名不能为空";
                      }
                  ),
                  TextFormField(
                      controller: _pwdController,
                      decoration: InputDecoration(
                          labelText: "密码",
                          hintText: "登录密码",
                          icon: Icon(Icons.lock)
                      ),
                      obscureText: true,
                      validator: (v) {
                        return v
                            .trim()
                            .length > 0 ? null : "密码不能为空";
                      }
                  ),
                ],
              ),
            ),
          ),
          new Container(
            margin: const EdgeInsets.symmetric(horizontal: 32.0, vertical: 20.0),
            child: new RaisedButton(
              padding: const EdgeInsets.all(12.0),
              shape: new RoundedRectangleBorder(borderRadius: BorderRadius.circular(12.0)),
              onPressed: () => onLoginButtonPressed(context),
              child: new Text("登录", style: new TextStyle(fontSize: 16.0, color: Colors.white),),
              color: pColor,
            ),
          )
        ],
      ),
    );
  }

  onLoginButtonPressed(BuildContext context) async {
    if((_formKey.currentState as FormState).validate()){
      Dio dio = DioManager.instance;

      FormData formData = new FormData.from({
        "username": _usernameController.text,
        "password": _pwdController.text,
        "user_type": _userType,
      });

      Response response;
      try {
        response = await dio.post("/user/mobile_login", data: formData);
      } catch(err) {
        print(err);
      }

      if (response?.statusCode == 200) {
        SharedPreferences prefs = await SharedPreferences.getInstance();
        await prefs.setString('userType', _userType);
        await prefs.setString('userNo', _usernameController.text);

        Navigator.pushReplacement(context, new MaterialPageRoute(
          builder: (context) {
            switch(_userType) {
              case TEACHER_USER:
                return new TeacherPage();
                break;
              case NORMAL_USER:
                return new NormalPage();
                break;
              default:
                return new TeacherPage();
            }
          },
          maintainState: false,
        ));
      } else {
        Fluttertoast.showToast(
          msg: "用户名或密码错误!",
          toastLength: Toast.LENGTH_SHORT,
          timeInSecForIos: 1,
        );
      }
    }
  }

  onUserTypeSelectPressed(String type) {
    _userType = type;
    primaryColorChangeCallBack(_userType == TEACHER_USER ? TEACHER_APP_COLOR : NORMAL_APP_COLOR);
  }

}

typedef PrimaryColorChangeCallBack = void Function(Color color);