import React, { PureComponent } from 'react'
import PropTypes from 'prop-types'
import { router } from 'utils'
import { connect } from 'dva'
import { Row, Col, Button, Popconfirm } from 'antd'
import { Page } from 'components'
import { stringify } from 'qs'
import List from './components/List'
import Filter from './components/Filter'
import Modal from './components/Modal'

@connect(({ classes, loading }) => ({ classes, loading }))
class Classes extends PureComponent {
  render() {
    const { location, dispatch, classes, loading } = this.props
    const { query, pathname } = location
    const {
      list,
      pagination,
      currentItem,
      modalVisible,
      modalType,
      selectedRowKeys,
    } = classes

    const handleRefresh = newQuery => {
      router.push({
        pathname,
        search: stringify(
          {
            ...query,
            ...newQuery,
          },
          { arrayFormat: 'repeat' }
        ),
      })
    }

    const modalProps = {
      item: modalType === 'create' ? {} : currentItem,
      type: modalType,
      visible: modalVisible,
      maskClosable: false,
      confirmLoading: loading.effects[`classes/${modalType}`],
      title: `${
        modalType === 'create' ? `创建班级信息` : `更新班级信息`
      }`,
      centered: true,
      onOk(data) {
        dispatch({
          type: `classes/${modalType}`,
          payload: data,
        }).then(() => {
          handleRefresh()
        })
      },
      onCancel() {
        dispatch({
          type: 'classes/hideModal',
        })
      },
    }

    const listProps = {
      dataSource: list,
      loading: loading.effects['classes/query'],
      pagination,
      onChange(page) {
        handleRefresh({
          page: page.current,
          pageSize: page.pageSize,
        })
      },
      onDeleteItem(id) {
        dispatch({
          type: 'classes/delete',
          payload: id,
        }).then(() => {
          handleRefresh({
            page:
              list.length === 1 && pagination.current > 1
                ? pagination.current - 1
                : pagination.current,
          })
        })
      },
      onEditItem(item) {
        dispatch({
          type: 'classes/showModal',
          payload: {
            modalType: 'update',
            currentItem: item,
          },
        })
      },
      rowSelection: {
        selectedRowKeys,
        onChange: keys => {
          dispatch({
            type: 'classes/updateState',
            payload: {
              selectedRowKeys: keys,
            },
          })
        },
      },
    }

    const filterProps = {
      filter: {
        ...query,
      },
      onFilterChange(value) {
        handleRefresh({
          ...value,
          page: 1,
        })
      },
      onAdd() {
        dispatch({
          type: 'classes/showModal',
          payload: {
            modalType: 'create',
          },
        })
      },
    }

    const handleDeleteItems = () => {
      dispatch({
        type: 'classes/multiDelete',
        payload: {
          class_ids: selectedRowKeys,
        },
      }).then(() => {
        handleRefresh({
          page:
            list.length === selectedRowKeys.length && pagination.current > 1
              ? pagination.current - 1
              : pagination.current,
        })
      })
    }

    return (
      <Page inner>
        <Filter {...filterProps} />
        {selectedRowKeys.length > 0 && (
          <Row style={{ marginBottom: 24, textAlign: 'right', fontSize: 13 }}>
            <Col>
              {`已选择 ${selectedRowKeys.length} 项 `}
              <Popconfirm
                title="确认删除选择项?"
                placement="left"
                onConfirm={handleDeleteItems}
              >
                <Button type="primary" style={{ marginLeft: 8 }}>
                  删除
                </Button>
              </Popconfirm>
            </Col>
          </Row>
        )}
        <List {...listProps} />
        {modalVisible && <Modal {...modalProps} />}
      </Page>
    )
  }
}

Classes.propTypes = {
  classes: PropTypes.object,
  location: PropTypes.object,
  dispatch: PropTypes.func,
  loading: PropTypes.object,
}

export default Classes