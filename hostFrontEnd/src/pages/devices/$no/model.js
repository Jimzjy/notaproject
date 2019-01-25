import { pathMatchRegexp } from 'utils'
import api from 'api'

const { 
  queryDevices
} = api

export default {
  namespace: 'deviceDetail',

  state: {
    data: {},
  },

  subscriptions: {
    setup({ dispatch, history }) {
      history.listen(({ pathname }) => {
        const match = pathMatchRegexp('/devices/:id', pathname)
        if (match) {
          dispatch({ type: 'query', payload: { device_id: match[1] } })
        }
      })
    },
  },

  effects: {
    *query({ payload }, { call, put }) {
      const data = yield call(queryDevices, payload)
      const { success, devices } = data
      if (success && devices.length > 0) {
        yield put({
          type: 'querySuccess',
          payload: {
            data: devices[0],
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