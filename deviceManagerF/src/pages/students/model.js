import { pageModel } from 'utils/model'
import modelExtend from 'dva-model-extend'
import { pathMatchRegexp } from 'utils'
import { objectToFormData } from 'utils/request'
import { 
  queryStudents,
  createStudent,
  removeStudents,
  updateStudent
} from 'api'

export default modelExtend(pageModel, {
    namespace: 'students',

    state: {
      currentItem: {},
      modalVisible: false,
      modalType: 'create',
      selectedRowKeys: [],
    },

    subscriptions: {
      setup({ dispatch, history }) {
        history.listen(location => {
          if (pathMatchRegexp('/students', location.pathname)) {
            const payload = location.query || { page: 1, pageSize: 10 }
            dispatch({
              type: 'query',
              payload,
            })
          }
        })
      },
    },

    effects: {
      *query({ payload = {} }, { call, put }) {
        const data = yield call(queryStudents, payload)
        if (data.success) {
          yield put({
            type: 'querySuccess',
            payload: {
              list: data.students,
              pagination: {
                current: Number(payload.page) || 1,
                pageSize: Number(payload.pageSize) || 10,
                total: data.total,
              },
            },
          })
        }
      },

      *delete({ payload }, { call, put, select }) {
        const data = yield call(removeStudents, objectToFormData({student_no: payload.student_no}))
        const { selectedRowKeys } = yield select(_ => _.students)
        if (data.success) {
          yield put({
            type: 'updateState',
            payload: {
              selectedRowKeys: selectedRowKeys.filter(_ => _ !== payload.student_no),
            },
          })
        } else {
          throw data
        }
      },
  
      *multiDelete({ payload }, { call, put }) {
        const data = yield call(removeStudents, objectToFormData(payload))
        if (data.success) {
          yield put({ type: 'updateState', payload: { selectedRowKeys: [] } })
        } else {
          throw data
        }
      },
  
      *create({ payload }, { call, put }) {
        const data = yield call(createStudent, objectToFormData(payload))
        if (data.success) {
          yield put({ type: 'hideModal' })
        } else {
          throw data
        }
      },
  
      *update({ payload }, { call, put }) {
        // const no = payload.student_no
        // const data = yield call(updateStudent, {formData: objectToFormData(payload), no})
        const data = yield call(updateStudent, objectToFormData(payload))
        if (data.success) {
          yield put({ type: 'hideModal' })
        } else {
          throw data
        }
      },
    },

    reducers: {
      showModal(state, { payload }) {
        return { ...state, ...payload, modalVisible: true }
      },
  
      hideModal(state) {
        return { ...state, modalVisible: false }
      },
    }
})