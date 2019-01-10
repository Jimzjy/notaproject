import { pathMatchRegexp } from 'utils'
import { queryStudents } from 'api'

export default {
  namespace: 'studentDetail',

  state: {
    data: {},
  },

  subscriptions: {
    setup({ dispatch, history }) {
      history.listen(({ pathname }) => {
        const match = pathMatchRegexp('/students/:no', pathname)
        if (match) {
          dispatch({ type: 'query', payload: { student_no: match[1] } })
        }
      })
    },
  },

  effects: {
    *query({ payload }, { call, put }) {
      const data = yield call(queryStudents, payload)
      const { success, students } = data
      if (success && students.length > 0) {
        yield put({
          type: 'querySuccess',
          payload: {
            data: students[0],
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