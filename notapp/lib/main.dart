import 'package:flutter/material.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'normal/normal_main.dart';
import 'teacher/teacher_main.dart';
import 'login/login_main.dart';

const NOT_LOGIN = "notlogin";
const TEACHER_USER = "teacher";
const NORMAL_USER = "normal";
const TEACHER_APP_COLOR = Colors.blue;
const NORMAL_APP_COLOR = Colors.indigo;

void main() {
  getUserType().then((v) {
    ThemeData theme;
    Widget homePage;

    switch(v) {
      case NORMAL_USER:
        theme = new ThemeData(primarySwatch: NORMAL_APP_COLOR);
        homePage = new NormalPage();
        break;
      case TEACHER_USER:
        theme = new ThemeData(primarySwatch: TEACHER_APP_COLOR);
        homePage = new TeacherPage();
        break;
      case NOT_LOGIN:
        runApp(new LoginApp());
        return;
    }

    runApp(new MyApp(theme: theme, homePage: homePage,));
  });
}

Future getUserType() async {
  SharedPreferences prefs = await SharedPreferences.getInstance();
  String userType = prefs.getString('userType') ?? NOT_LOGIN;
  return userType;
}

class MyApp extends StatelessWidget {
  MyApp({ @required this.theme, @required this.homePage });

  final ThemeData theme;
  final Widget homePage;

  @override
  Widget build(BuildContext context) {
    return new MaterialApp(
      title: 'NOTAPP',
      theme: theme,
      home: homePage,
    );
  }
}