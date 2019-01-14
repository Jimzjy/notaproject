import { pageModel } from 'utils/model'
import modelExtend from 'dva-model-extend'
import { pathMatchRegexp } from 'utils'
import { objectToFormData } from 'utils/request'
import { 
  queryCameras,
  createCamera,
  removeCameras,
  updateCamera
} from 'api'

export default modelExtend(pageModel, {
    namespace: 'cameras',

    state: {
      currentItem: {},
      modalVisible: false,
      modalType: 'create',
      selectedRowKeys: [],
    },

    subscriptions: {
      setup({ dispatch, history }) {
        history.listen(location => {
          if (pathMatchRegexp('/cameras', location.pathname)) {
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
        const data = yield call(queryCameras, payload)
        if (data.success) {
          yield put({
            type: 'querySuccess',
            payload: {
              list: data.cameras,
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
        const data = yield call(removeCameras, objectToFormData({camera_id: payload.camera_id}))
        const { selectedRowKeys } = yield select(_ => _.cameras)
        if (data.success) {
          yield put({
            type: 'updateState',
            payload: {
              selectedRowKeys: selectedRowKeys.filter(_ => _ !== payload.camera_id),
            },
          })
        } else {
          throw data
        }
      },
  
      *multiDelete({ payload }, { call, put }) {
        const data = yield call(removeCameras, objectToFormData(payload))
        if (data.success) {
          yield put({ type: 'updateState', payload: { selectedRowKeys: [] } })
        } else {
          throw data
        }
      },
  
      *create({ payload }, { call, put }) {
        const data = yield call(createCamera, objectToFormData(payload))
        if (data.success) {
          yield put({ type: 'hideModal' })
        } else {
          throw data
        }
      },
  
      *update({ payload }, { call, put }) {
        const data = yield call(updateCamera, objectToFormData(payload))
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