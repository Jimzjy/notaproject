import { pathMatchRegexp } from 'utils'
import api from 'api'

const { 
  queryTeachers
} = api

export default {
  namespace: 'teacherDetail',

  state: {
    data: {},
  },

  subscriptions: {
    setup({ dispatch, history }) {
      history.listen(({ pathname }) => {
        const match = pathMatchRegexp('/teachers/:no', pathname)
        if (match) {
          dispatch({ type: 'query', payload: { teacher_no: match[1] } })
        }
      })
    },
  },

  effects: {
    *query({ payload }, { call, put }) {
      const data = yield call(queryTeachers, payload)
      const { success, teachers } = data
      if (success && teachers.length > 0) {
        yield put({
          type: 'querySuccess',
          payload: {
            data: teachers[0],
          },
        })
      } else {
        throw data
      }
    },
  },

  reducers: {
    querySuccess(state, { payload }) {
      const { data } = payload
      return {
        ...state,
        data,
      }
    },
  },
}