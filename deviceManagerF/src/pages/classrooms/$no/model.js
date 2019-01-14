import { pathMatchRegexp } from 'utils'
import { queryClassrooms } from 'api'

export default {
  namespace: 'classroomDetail',

  state: {
    data: {},
  },

  subscriptions: {
    setup({ dispatch, history }) {
      history.listen(({ pathname }) => {
        const match = pathMatchRegexp('/classrooms/:no', pathname)
        if (match) {
          dispatch({ type: 'query', payload: { classroom_no: match[1] } })
        }
      })
    },
  },

  effects: {
    *query({ payload }, { call, put }) {
      const data = yield call(queryClassrooms, payload)
      const { success, classroom } = data
      if (success && classroom.length > 0) {
        yield put({
          type: 'querySuccess',
          payload: {
            data: classroom[0],
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