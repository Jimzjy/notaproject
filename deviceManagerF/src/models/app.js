/* global window */

import { router } from 'utils'
import { stringify } from 'qs'
import store from 'store'
import { ROLE_TYPE, ROUTE_LIST } from 'utils/constant'
import { queryLayout, pathMatchRegexp } from 'utils'
import { CANCEL_REQUEST_MESSAGE } from 'utils/constant'
import { logoutAdmin, queryAdminInfo } from 'api'
import config from 'config'

const UserPermission = {
  DEFAULT: {
    visit: ['1', '5', '6', '51', '52', '53', '61', '62', '621'],
    role: ROLE_TYPE.DEFAULT,
  },
  DEVELOPER: {
    role: ROLE_TYPE.DEVELOPER,
  },
}

export default {
  namespace: 'app',
  state: {
    user: {},
    permissions: {
      visit: [],
    },
    routeList: [
      {
        id: '1',
        icon: 'dashboard',
        name: 'Dashboard',
        zhName: '仪表盘',
        route: '/dashboard',
      },
    ],
    locationPathname: '',
    locationQuery: {},
    theme: store.get('theme') || 'light',
    collapsed: store.get('collapsed') || false,
    notifications: [
    ],
  },
  subscriptions: {
    setupHistory({ dispatch, history }) {
      history.listen(location => {
        dispatch({
          type: 'updateState',
          payload: {
            locationPathname: location.pathname,
            locationQuery: location.query,
          },
        })
      })
    },

    setupRequestCancel({ history }) {
      history.listen(() => {
        const { cancelRequest = new Map() } = window

        cancelRequest.forEach((value, key) => {
          if (value.pathname !== window.location.pathname) {
            value.cancel(CANCEL_REQUEST_MESSAGE)
            cancelRequest.delete(key)
          }
        })
      })
    },

    setup({ dispatch }) {
      dispatch({ type: 'query' })
    },
  },
  effects: {
    *query({ payload }, { call, put, select }) {
      const { success, user } = yield call(queryAdminInfo, payload)
      const { locationPathname } = yield select(_ => _.app)

      if (success && user) {
        const list = ROUTE_LIST
        let permissions = {}
        let routeList = list
        if (user.permissions === UserPermission.DEVELOPER.role) {
          permissions.visit = list.map(item => item.id)
        } else {
          permissions.visit = UserPermission.DEFAULT.visit
          routeList = list.filter(item => {
            const cases = [
              permissions.visit.includes(item.id),
              item.mpid
                ? permissions.visit.includes(item.mpid) || item.mpid === '-1'
                : true,
              item.bpid ? permissions.visit.includes(item.bpid) : true,
            ]
            return cases.every(_ => _)
          })
        }
        yield put({
          type: 'updateState',
          payload: {
            user,
            permissions,
            routeList,
          },
        })
        if (pathMatchRegexp('/login', window.location.pathname)) {
          router.push({
            pathname: '/dashboard',
          })
        }
      } else if (queryLayout(config.layouts, locationPathname) !== 'public') {
        router.push({
          pathname: '/login',
          search: stringify({
            from: locationPathname,
          }),
        })
      }
    },

    *signOut({ payload }, { call, put }) {
      const data = yield call(logoutAdmin)
      if (data.success) {
        yield put({
          type: 'updateState',
          payload: {
            user: {},
            permissions: { visit: [] },
            menu: [],
          },
        })
        yield put({ type: 'query' })
      } else {
        throw data
      }
    },
  },
  reducers: {
    updateState(state, { payload }) {
      return {
        ...state,
        ...payload,
      }
    },

    handleThemeChange(state, { payload }) {
      store.set('theme', payload)
      state.theme = payload
    },

    handleCollapseChange(state, { payload }) {
      store.set('collapsed', payload)
      state.collapsed = payload
    },

    allNotificationsRead(state) {
      state.notifications = []
    },
  },
}
