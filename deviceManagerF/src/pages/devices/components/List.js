import React, { PureComponent } from 'react'
import PropTypes from 'prop-types'
import { Table, Modal, Tag } from 'antd'
import { DropOption } from 'components'
import { router } from 'utils'
import { Color } from 'utils'
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
          onDeleteItem({device_id: record.device_id})
        },
      })
    }
  }

  render() {
    const { onDeleteItem, onEditItem, ...tableProps } = this.props

    const columns = [
      {
        title: '设备ID',
        dataIndex: 'device_id',
        key: 'device_id',
        render: text => <Link to={`devices/${text}`}>{text}</Link>,
      },
      {
        title: '设备IP',
        dataIndex: 'device_path',
        key: 'device_path',
      },
      {
        title: '设备Port',
        dataIndex: 'device_port',
        key: 'device_port',
      },
      {
        title: '摄像头',
        dataIndex: 'camera_ids',
        key: 'camera_ids',
        render: camera_ids => {
          return (
            <span>
              {camera_ids.map(c => <Tag.CheckableTag checked={true} color={Color.blue} key={c} onChange={ () => router.push(`/cameras/${c}`) }>{c}</Tag.CheckableTag>)}
            </span>
          )
        },
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
        rowKey={record => record.device_id}
      />
    )
  }
}

List.propTypes = {
  onDeleteItem: PropTypes.func,
  onEditItem: PropTypes.func,
}

export default List
