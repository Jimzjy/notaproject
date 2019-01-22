export default {
  queryUserInfo: '/user',
  logoutUser: '/user/logout',
  loginUser: 'POST /user/login',

  queryDashboard: '/dashboard',

  queryStudents: '/students',
  createStudent: 'POST /students',
  removeStudents: 'DELETE /students',
  updateStudent: 'PATCH /students',

  queryClasses: '/classes',
  createClass: 'POST /classes',
  removeClasses: 'DELETE /classes',
  updateClass: 'PATCH /classes',

  queryTeachers: '/teachers',
  createTeacher: 'POST /teachers',
  removeTeachers: 'DELETE /teachers',
  updateTeacher: 'PATCH /teachers',

  queryDevices: '/devices',
  createDevice: 'POST /devices',
  removeDevices: 'DELETE /devices',
  updateDevice: 'PATCH /devices',

  queryCameras: '/cameras',
  createCamera: 'POST /cameras',
  removeCameras: 'DELETE /cameras',
  updateCamera: 'PATCH /cameras',

  queryClassrooms: '/classrooms',
  createClassroom: 'POST /classrooms',
  removeClassrooms: 'DELETE /classrooms',
  updateClassroom: 'PATCH /classrooms',

  queryStandupClasses: 'POST /stand_up_classes',
}
