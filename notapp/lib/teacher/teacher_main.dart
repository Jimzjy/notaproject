import 'package:flutter/material.dart';
import 'main_navigation.dart';

class TeacherPage extends StatefulWidget {
  TeacherPage({Key key}) : super(key: key);

  @override
  State<StatefulWidget> createState() => _TeacherPageState();
}

class _TeacherPageState extends State<TeacherPage> {
  var _pageSelected = 0;

  final _bottomItem = <BottomNavigationBarItem>[
    BottomNavigationBarItem(icon: Icon(Icons.group), title: Text('班级')),
    BottomNavigationBarItem(icon: Icon(Icons.bubble_chart), title: Text('记录')),
    BottomNavigationBarItem(icon: Icon(Icons.person), title: Text('我的')),
  ];

  final _navigationPages = <Widget>[
    ClassesPage(),
    HistoryPage(),
    MinePage(),
  ];
  final _pageController = new PageController(initialPage: 0);

  @override
  Widget build(BuildContext context) {
    return new Scaffold(
      body: new PageView(
        physics: NeverScrollableScrollPhysics(),
        controller: _pageController,
        children: _navigationPages,
      ),
      bottomNavigationBar: new BottomNavigationBar(
        currentIndex: _pageSelected,
        items: _bottomItem,
        type: BottomNavigationBarType.fixed,
        onTap: _onPageChanged,
      ),
    );
  }

  @override
  void dispose() {
    _pageController.dispose();
    super.dispose();
  }

  _onPageChanged(int selected) {
    setState(() {
      _pageSelected = selected;
      _pageController.jumpToPage(_pageSelected);
    });
  }
}