import { parse } from 'qs'
import modelExtend from 'dva-model-extend'
import { pathMatchRegexp } from 'utils'
import { model } from 'utils/model'
import { Color } from 'utils'
import api from 'api'

const { queryDashboard } = api

const numbersTemplate = [
  {
    icon: 'deployment-unit',
    color: Color.blue,
    title: '计算设备',
    number: 0,
  },
  {
    icon: 'video-camera',
    color: Color.green,
    title: '摄像头',
    number: 0,
  },
  {
    icon: 'meh',
    color: Color.yellow,
    title: '警告',
    number: 0,
  },
  {
    icon: 'frown',
    color: Color.red,
    title: '错误',
    number: 0,
  },
]

const systemStatsTemplate = {
  cpuUsage: 0,
  memUsage: 0,
  systemStats: [],
}

export default modelExtend(model, {
  namespace: 'dashboard',
  state: {
    numbers: numbersTemplate,
    systemStats: systemStatsTemplate,
  },
  subscriptions: {
    setup({ dispatch, history }) {
      history.listen(({ pathname }) => {
        if (
          pathMatchRegexp('/dashboard', pathname) ||
          pathMatchRegexp('/', pathname)
        ) {
          dispatch({ type: 'query' })
        }
      })
    },
  },
  effects: {
    *query({ payload }, { call, put }) {
      const data = yield call(queryDashboard, parse(payload))
      const { number_card, system_stats, ...dashBoardData } = data

      dashBoardData.numbers = numbersTemplate
      dashBoardData.numbers[0].number = number_card.devices
      dashBoardData.numbers[1].number = number_card.cameras

      dashBoardData.systemStats = systemStatsTemplate
      dashBoardData.systemStats.cpuUsage = system_stats[0].cpu_used
      dashBoardData.systemStats.memUsage = system_stats[0].mem_used
      dashBoardData.systemStats.systemStats = system_stats.reverse()

      yield put({
        type: 'updateState',
        payload: dashBoardData,
      })
    },
  },
})
