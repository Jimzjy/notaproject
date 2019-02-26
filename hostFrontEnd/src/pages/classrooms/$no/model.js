import { pathMatchRegexp } from 'utils'
import api from 'api'

const { 
  queryClassrooms,
  queryClassroomStats
} = api

export default {
  namespace: 'classroomDetail',

  state: {
    data: {},
    stats: {},
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
      const stats = yield call(queryClassroomStats, payload)
      const { classrooms } = data
      const success1 = data.success
      const success2 = stats.success

      if (success1 && success2 && classrooms.length > 0) {
        yield put({
          type: 'querySuccess',
          payload: {
            data: classrooms[0],
            stats: stats,
          },
        })
      } else {
        throw data
      }
    },
  },

  reducers: {
    querySuccess(state, { payload }) {
      const { data, stats } = payload
      return {
        ...state,
        data,
        stats,
      }
    },
  },
}