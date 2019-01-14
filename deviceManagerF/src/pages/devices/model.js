import { pageModel } from 'utils/model'
import modelExtend from 'dva-model-extend'
import { pathMatchRegexp } from 'utils'
import { objectToFormData } from 'utils/request'
import { 
  queryDevices,
  createDevice,
  removeDevices,
  updateDevice
} from 'api'

export default modelExtend(pageModel, {
    namespace: 'devices',

    state: {
      currentItem: {},
      modalVisible: false,
      modalType: 'create',
      selectedRowKeys: [],
    },

    subscriptions: {
      setup({ dispatch, history }) {
        history.listen(location => {
          if (pathMatchRegexp('/devices', location.pathname)) {
            const payload = location.query || { page: 1, pageSize: 10 }
            dispatch({
              type: 'query',
              payload,
            })
          }
        })
      },
    },

    effects: {
      *query({ payload = {} }, { call, put }) {
        const data = yield call(queryDevices, payload)
        if (data.success) {
          yield put({
            type: 'querySuccess',
            payload: {
              list: data.devices,
              pagination: {
                current: Number(payload.page) || 1,
                pageSize: Number(payload.pageSize) || 10,
                total: data.total,
              },
            },
          })
        }
      },

      *delete({ payload }, { call, put, select }) {
        const data = yield call(removeDevices, objectToFormData({device_id: payload.device_id}))
        const { selectedRowKeys } = yield select(_ => _.devices)
        if (data.success) {
          yield put({
            type: 'updateState',
            payload: {
              selectedRowKeys: selectedRowKeys.filter(_ => _ !== payload.device_id),
            },
          })
        } else {
          throw data
        }
      },
  
      *multiDelete({ payload }, { call, put }) {
        const data = yield call(removeDevices, objectToFormData(payload))
        if (data.success) {
          yield put({ type: 'updateState', payload: { selectedRowKeys: [] } })
        } else {
          throw data
        }
      },
  
      *create({ payload }, { call, put }) {
        const data = yield call(createDevice, objectToFormData(payload))
        if (data.success) {
          yield put({ type: 'hideModal' })
        } else {
          throw data
        }
      },
  
      *update({ payload }, { call, put }) {
        const data = yield call(updateDevice, objectToFormData(payload))
        if (data.success) {
          yield put({ type: 'hideModal' })
        } else {
          throw data
        }
      },
    },

    reducers: {
      showModal(state, { payload }) {
        return { ...state, ...payload, modalVisible: true }
      },
  
      hideModal(state) {
        return { ...state, modalVisible: false }
      },
    }
})