import { pathMatchRegexp } from 'utils'
import api from 'api'

const { queryCameras } = api

export default {
  namespace: 'cameraDetail',

  state: {
    data: {},
  },

  subscriptions: {
    setup({ dispatch, history }) {
      history.listen(({ pathname }) => {
        const match = pathMatchRegexp('/cameras/:id', pathname)
        if (match) {
          dispatch({ type: 'query', payload: { camera_id: match[1] } })
        }
      })
    },
  },

  effects: {
    *query({ payload }, { call, put }) {
      const data = yield call(queryCameras, payload)
      const { success, cameras } = data
      if (success && cameras.length > 0) {
        yield put({
          type: 'querySuccess',
          payload: {
            data: cameras[0],
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