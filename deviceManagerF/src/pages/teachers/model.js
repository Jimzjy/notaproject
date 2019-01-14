import { pageModel } from 'utils/model'
import modelExtend from 'dva-model-extend'
import { pathMatchRegexp } from 'utils'
import { objectToFormData } from 'utils/request'
import { 
  queryTeachers,
  createTeacher,
  removeTeachers,
  updateTeacher
} from 'api'

export default modelExtend(pageModel, {
    namespace: 'teachers',

    state: {
      currentItem: {},
      modalVisible: false,
      modalType: 'create',
      selectedRowKeys: [],
    },

    subscriptions: {
      setup({ dispatch, history }) {
        history.listen(location => {
          if (pathMatchRegexp('/teachers', location.pathname)) {
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
        const data = yield call(queryTeachers, payload)
        if (data.success) {
          yield put({
            type: 'querySuccess',
            payload: {
              list: data.teachers,
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
        const data = yield call(removeTeachers, objectToFormData({teacher_no: payload.teacher_no}))
        const { selectedRowKeys } = yield select(_ => _.teachers)
        if (data.success) {
          yield put({
            type: 'updateState',
            payload: {
              selectedRowKeys: selectedRowKeys.filter(_ => _ !== payload.teacher_no),
            },
          })
        } else {
          throw data
        }
      },
  
      *multiDelete({ payload }, { call, put }) {
        const data = yield call(removeTeachers, objectToFormData(payload))
        if (data.success) {
          yield put({ type: 'updateState', payload: { selectedRowKeys: [] } })
        } else {
          throw data
        }
      },
  
      *create({ payload }, { call, put }) {
        const data = yield call(createTeacher, objectToFormData(payload))
        if (data.success) {
          yield put({ type: 'hideModal' })
        } else {
          throw data
        }
      },
  
      *update({ payload }, { call, put }) {
        const data = yield call(updateTeacher, objectToFormData(payload))
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