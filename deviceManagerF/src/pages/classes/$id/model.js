import { pathMatchRegexp } from 'utils'
import api from 'api'

const { 
  queryClasses
} = api

export default {
  namespace: 'classDetail',

  state: {
    data: {},
  },

  subscriptions: {
    setup({ dispatch, history }) {
      history.listen(({ pathname }) => {
        const match = pathMatchRegexp('/classes/:id', pathname)
        if (match) {
          dispatch({ type: 'query', payload: { class_id: match[1] } })
        }
      })
    },
  },

  effects: {
    *query({ payload }, { call, put }) {
      const data = yield call(queryClasses, payload)
      const { success, classes } = data
      if (success && classes.length > 0) {
        yield put({
          type: 'querySuccess',
          payload: {
            data: classes[0],
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