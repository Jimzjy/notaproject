import { pathMatchRegexp } from 'utils'
import api from 'api'

const { 
  queryDevices,
  queryDeviceStats,
} = api

const systemStatsTemplate = {
  cpuUsage: 0,
  memUsage: 0,
  systemStats: [],
}

export default {
  namespace: 'deviceDetail',

  state: {
    data: {},
    systemStats: systemStatsTemplate,
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
      const stats = yield call(queryDeviceStats, payload)
      const { devices } = data
      const success1 = data.success
      const success2 = stats.success

      const _stats = systemStatsTemplate
      _stats.cpuUsage = stats.device_stats[0].cpu_used
      _stats.memUsage = stats.device_stats[0].mem_used
      _stats.systemStats = stats.device_stats.reverse()

      if (success1 && success2 && devices.length > 0) {
        yield put({
          type: 'querySuccess',
          payload: {
            data: devices[0],
            systemStats: _stats,
          },
        })
      } else {
        throw data
      }
    },
  },

  reducers: {
    querySuccess(state, { payload }) {
      const { data, systemStats } = payload
      return {
        ...state,
        data,
        systemStats,
      }
    },
  },
}