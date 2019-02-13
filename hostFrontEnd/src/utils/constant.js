export const ROLE_TYPE = {
  ADMIN: 'admin',
  TEACHER: 'teacher',
}

export const CANCEL_REQUEST_MESSAGE = 'cancle request'

export const ROUTE_LIST = [
  {
    id: '1',
    icon: 'dashboard',
    name: 'Dashboard',
    zhName: '仪表盘',
    route: '/dashboard',
  },
  {
    id: '5',
    breadcrumbParentId: '1',
    name: 'Device Management',
    zhName: '设备管理',
    icon: 'cluster',
  },
  {
    id: '51',
    breadcrumbParentId: '5',
    menuParentId: '5',
    name: 'Devices',
    zhName: '设备',
    icon: 'deployment-unit',
    route: '/devices',
  },
  {
    id: '52',
    breadcrumbParentId: '5',
    menuParentId: '5',
    name: 'Cameras',
    zhName: '摄像头',
    icon: 'video-camera',
    route: '/cameras',
  },
  {
    id: '53',
    breadcrumbParentId: '5',
    menuParentId: '5',
    name: 'Classrooms',
    zhName: '教室',
    icon: 'read',
    route: '/classrooms',
  },
  {
    id: '511',
    breadcrumbParentId: '51',
    menuParentId: '-1',
    name: 'Device Detail',
    zhName: '设备详情',
    route: '/devices/:id',
  },
  {
    id: '521',
    breadcrumbParentId: '52',
    menuParentId: '-1',
    name: 'Camera Detail',
    zhName: '摄像头详情',
    route: '/cameras/:id',
  },
  {
    id: '531',
    breadcrumbParentId: '53',
    menuParentId: '-1',
    name: 'Classrooms Detail',
    zhName: '教室详情',
    route: '/classrooms/:id',
  },
  {
    id: '6',
    breadcrumbParentId: '1',
    name: 'Students Management',
    zhName: '学生管理',
    icon: 'calendar',
  },
  {
    id: '61',
    breadcrumbParentId: '6',
    menuParentId: '6',
    name: 'Classes',
    zhName: '班级',
    icon: 'team',
    route: '/classes',
  },
  {
    id: '62',
    breadcrumbParentId: '6',
    menuParentId: '6',
    name: 'Students',
    zhName: '学生',
    icon: 'user',
    route: '/students',
  },
  {
    id: '63',
    breadcrumbParentId: '6',
    menuParentId: '6',
    name: 'Teachers',
    zhName: '教师',
    icon: 'solution',
    route: '/teachers',
  },
  {
    id: '611',
    breadcrumbParentId: '61',
    menuParentId: '-1',
    name: 'Class Detail',
    zhName: '班级详情',
    route: '/classes/:no',
  },
  {
    id: '621',
    breadcrumbParentId: '62',
    menuParentId: '-1',
    name: 'Student Detail',
    zhName: '学生详情',
    route: '/students/:no',
  },
  {
    id: '631',
    breadcrumbParentId: '63',
    menuParentId: '-1',
    name: 'Teacher Detail',
    zhName: '教师详情',
    route: '/teachers/:no',
  },
  {
    id: '7',
    icon: 'experiment',
    name: 'StandUp',
    zhName: '上课',
    route: '/standup',
  },
  {
    id: '71',
    breadcrumbParentId: '7',
    menuParentId: '-1',
    name: 'Standup Detail',
    zhName: '上课详情',
    route: '/standup/:id',
  },
]