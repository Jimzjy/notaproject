import React, { PureComponent } from 'react'
import PropTypes from 'prop-types'
import { Table, Modal, Avatar, Tag } from 'antd'
import { DropOption } from 'components'
import { router } from 'utils'
import { Color } from 'utils'
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
        title: `确认删除所选学生信息?`,
        onOk() {
          onDeleteItem({student_no: record.student_no})
        },
      })
    }
  }

  render() {
    const { onDeleteItem, onEditItem, ...tableProps } = this.props

    const columns = [
      {
        title: '头像',
        dataIndex: 'student_image',
        key: 'student_image',
        width: 72,
        fixed: 'left',
        render: text => <Avatar style={{ marginLeft: 8 }} src={`${apiPrefix}/images/${text}`} />,
      },
      {
        title: '学号',
        dataIndex: 'student_no',
        key: 'student_no',
        render: text => <Link to={`students/${text}`}>{text}</Link>,
      },
      {
        title: '姓名',
        dataIndex: 'student_name',
        key: 'student_name',
      },
      {
        title: '班级',
        dataIndex: 'class_ids',
        key: 'class_ids',
        render: class_ids => {
          return (
            <span>
              {class_ids.map(c => <Tag color={Color.blue} key={c} onChange={ router.push(`classes/${c}`) }>{c}</Tag>)}
            </span>
          )
        },
      },
      {
        title: 'Face Token',
        dataIndex: 'face_token',
        key: 'face_token',
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
        rowKey={record => record.student_no}
      />
    )
  }
}

List.propTypes = {
  onDeleteItem: PropTypes.func,
  onEditItem: PropTypes.func,
}

export default List
