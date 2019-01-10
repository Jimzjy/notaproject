export const ROLE_TYPE = {
  DEFAULT: 'admin',
  DEVELOPER: 'developer',
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
    id: '2',
    breadcrumbParentId: '1',
    name: 'Request',
    zhName: 'Request',
    icon: 'api',
    route: '/request',
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
    route: '/user',
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
    id: '621',
    breadcrumbParentId: '62',
    menuParentId: '-1',
    name: 'Student Detail',
    zhName: '学生详情',
    route: '/students/:no',
  },
]