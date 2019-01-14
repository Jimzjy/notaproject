import React, { PureComponent } from 'react'
import PropTypes from 'prop-types'
import { Table, Modal, Tag } from 'antd'
import { DropOption } from 'components'
import Link from 'umi/link'
import styles from './List.less'

const { confirm } = Modal

class List extends PureComponent {
  handleMenuClick = (record, e) => {
    const { onDeleteItem, onEditItem } = this.props

    if (e.key === '1') {
      onEditItem(record)
    } else if (e.key === '2') {
      confirm({
        title: `确认删除所选信息?`,
        onOk() {
          onDeleteItem({classroom_no: record.classroom_no})
        },
      })
    }
  }

  render() {
    const { onDeleteItem, onEditItem, ...tableProps } = this.props

    const columns = [
      {
        title: '教室',
        dataIndex: 'classroom_no',
        key: 'classroom_no',
        render: text => <Link to={`classrooms/${text}`}>{text}</Link>,
      },
      {
        title: '摄像头ID',
        dataIndex: 'camera_id',
        key: 'camera_id',
        render: text => <Link to={`cameras/${text}`}>{text}</Link>,
      },
      {
        title: '操作',
        key: 'operation',
        render: (_, record) => {
          return (
            <DropOption
              onMenuClick={e => this.handleMenuClick(record, e)}
              menuOptions={[
                { key: '1', name: `更新` },
                { key: '2', name: `删除` },
              ]}
            />
          )
        },
      },
    ]

    return (
      <Table
        {...tableProps}
        pagination={{
          ...tableProps.pagination,
          showTotal: total => `共 ${total} 条记录`,
        }}
        className={styles.table}
        bordered
        scroll={{ x: 800 }}
        columns={columns}
        simple
        rowKey={record => record.classroom_no}
      />
    )
  }
}

List.propTypes = {
  onDeleteItem: PropTypes.func,
  onEditItem: PropTypes.func,
}

export default List
