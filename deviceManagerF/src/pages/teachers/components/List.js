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
        title: `确认删除所选教师信息?`,
        onOk() {
          onDeleteItem({teacher_no: record.teacher_no})
        },
      })
    }
  }

  render() {
    const { onDeleteItem, onEditItem, ...tableProps } = this.props

    const columns = [
      {
        title: '头像',
        dataIndex: 'teacher_image',
        key: 'teacher_image',
        width: 72,
        fixed: 'left',
        render: text => <Avatar style={{ marginLeft: 8 }} src={`${apiPrefix}/images/${text}`} />,
      },
      {
        title: '教师号',
        dataIndex: 'teacher_no',
        key: 'teacher_no',
        render: text => <Link to={`teachers/${text}`}>{text}</Link>,
      },
      {
        title: '姓名',
        dataIndex: 'teacher_name',
        key: 'teacher_name',
      },
      {
        title: '班级',
        dataIndex: 'class_ids',
        key: 'class_ids',
        render: class_ids => {
          return (
            <span>
              {class_ids.map(c => <Tag.CheckableTag checked={true} color={Color.blue} key={c} onChange={ () => router.push(`/classes/${c}`) }>{c}</Tag.CheckableTag>)}
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
        rowKey={record => record.teacher_no}
      />
    )
  }
}

List.propTypes = {
  onDeleteItem: PropTypes.func,
  onEditItem: PropTypes.func,
}

export default List
