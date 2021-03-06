/* global window */

import { router } from 'utils'
import { stringify } from 'qs'
import store from 'store'
import { ROLE_TYPE, ROUTE_LIST } from 'utils/constant'
import { queryLayout, pathMatchRegexp } from 'utils'
import { CANCEL_REQUEST_MESSAGE } from 'utils/constant'
import config from 'config'
import api from 'api'

const { logoutUser, queryUserInfo } = api

const UserPermission = {
  ADMIN: {
    visit: [
      '1',
      '5',
      '6',
      '51',
      '52',
      '53',
      '511',
      '521',
      '531',
      '61',
      '62',
      '63',
      '611',
      '621',
      '631',
    ],
    role: ROLE_TYPE.ADMIN,
  },
  TEACHER: {
    visit: ['7', '71'],
    role: ROLE_TYPE.TEACHER,
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
    notifications: [],
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
      const { success, user } = yield call(queryUserInfo, payload)
      const { locationPathname } = yield select(_ => _.app)

      if (success && user) {
        const list = ROUTE_LIST
        let permissions = {}
        let routeList = list
        let routeTo = '/dashboard'
        switch (user.permissions) {
          case UserPermission.ADMIN.role:
            permissions.visit = UserPermission.ADMIN.visit
            break
          case UserPermission.TEACHER.role:
            permissions.visit = UserPermission.TEACHER.visit
            routeTo = '/standup'
            break
          default:
            permissions.visit = UserPermission.ADMIN.visit
        }
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
        yield put({
          type: 'updateState',
          payload: {
            user,
            permissions,
            routeList,
          },
        })
        if (pathMatchRegexp(['/', '/login'], window.location.pathname)) {
          router.push({
            pathname: routeTo,
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
      const data = yield call(logoutUser)
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
