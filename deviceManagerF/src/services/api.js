export default {
  queryAdminInfo: '/admin',
  logoutAdmin: '/admin/logout',
  loginAdmin: 'POST /admin/login',

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
}