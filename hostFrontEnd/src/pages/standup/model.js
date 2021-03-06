import modelExtend from 'dva-model-extend'
import { pathMatchRegexp } from 'utils'
import { model } from 'utils/model'
import api from 'api'

const { queryClasses } = api

export default modelExtend(model, {
  namespace: 'standup',
  state: {
    classes: [],
  },
  subscriptions: {
    setup({ dispatch, history }) {
      history.listen(({ pathname }) => {
        if (pathMatchRegexp('/standup', pathname)) {
          dispatch({ type: 'query', payload: {} })
        }
      })
    },
  },
  effects: {
    *query({ payload }, { call, put, select }) {
      const { user } = yield select(_ => _.app)
      payload.teacher_no = user.username
      const data = yield call(queryClasses, payload)

      if (data.success) {
        yield put({
          type: 'updateState',
          payload: data,
        })
      }
    },
  },
})
