import { pageModel } from 'utils/model'
import modelExtend from 'dva-model-extend'
import { pathMatchRegexp } from 'utils'
import { objectToFormData } from 'utils/request'
import api from 'api'

const { queryClasses, createClass, removeClasses, updateClass } = api

export default modelExtend(pageModel, {
  namespace: 'classes',

  state: {
    currentItem: {},
    modalVisible: false,
    modalType: 'create',
    selectedRowKeys: [],
  },

  subscriptions: {
    setup({ dispatch, history }) {
      history.listen(location => {
        if (pathMatchRegexp('/classes', location.pathname)) {
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
      const data = yield call(queryClasses, payload)
      if (data.success) {
        yield put({
          type: 'querySuccess',
          payload: {
            list: data.classes,
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
      const data = yield call(
        removeClasses,
        objectToFormData({ class_id: payload.class_id })
      )
      const { selectedRowKeys } = yield select(_ => _.classes)
      if (data.success) {
        yield put({
          type: 'updateState',
          payload: {
            selectedRowKeys: selectedRowKeys.filter(
              _ => _ !== payload.class_id
            ),
          },
        })
      } else {
        throw data
      }
    },

    *multiDelete({ payload }, { call, put }) {
      const data = yield call(removeClasses, objectToFormData(payload))
      if (data.success) {
        yield put({ type: 'updateState', payload: { selectedRowKeys: [] } })
      } else {
        throw data
      }
    },

    *create({ payload }, { call, put }) {
      const data = yield call(createClass, objectToFormData(payload))
      if (data.success) {
        yield put({ type: 'hideModal' })
      } else {
        throw data
      }
    },

    *update({ payload }, { call, put }) {
      const data = yield call(updateClass, objectToFormData(payload))
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
  },
})
