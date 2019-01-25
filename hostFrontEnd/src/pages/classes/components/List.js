import React, { PureComponent } from 'react'
import PropTypes from 'prop-types'
import { Table, Modal, Avatar, Tag } from 'antd'
import { DropOption } from 'components'
import Link from 'umi/link'
import styles from './List.less'
import { apiPrefix } from 'utils/config'

const { confirm } = Modal

class List extends PureComponent {
  handleMenuClick = (record, e) => {
    const { onDeleteItem, onEditItem } = this.props

    if (e.key === '1') {
      onEditItem(record)
    } else if (e.key === '2') {
      confirm({
        title: `确认删除所选班级信息?`,
        onOk() {
          onDeleteItem({ class_id: record.class_id })
        },
      })
    }
  }

  render() {
    const { onDeleteItem, onEditItem, ...tableProps } = this.props

    const columns = [
      {
        title: '班级图',
        dataIndex: 'class_image',
        key: 'class_image',
        width: 72,
        fixed: 'left',
        render: text => (
          <Avatar
            style={{ marginLeft: 8 }}
            src={`${apiPrefix}/images/${text}`}
          />
        ),
      },
      {
        title: '班级号',
        dataIndex: 'class_id',
        key: 'class_id',
        render: text => <Link to={`classes/${text}`}>{text}</Link>,
      },
      {
        title: '班级名',
        dataIndex: 'class_name',
        key: 'class_name',
      },
      {
        title: '班级人数',
        dataIndex: 'face_count',
        key: 'face_count',
      },
      {
        title: '教室号',
        dataIndex: 'classroom_no',
        key: 'classroom_no',
      },
      {
        title: 'FaceSet Token',
        dataIndex: 'faceset_token',
        key: 'faceset_token',
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
        rowKey={record => record.class_id}
      />
    )
  }
}

List.propTypes = {
  onDeleteItem: PropTypes.func,
  onEditItem: PropTypes.func,
}

export default List
