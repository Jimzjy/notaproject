import React, { PureComponent } from 'react'
import PropTypes from 'prop-types'
import { Form, Input, Modal } from 'antd'
import Selects from './Selects'
import PicturesWall from './PicturesWall'

const FormItem = Form.Item

const formItemLayout = {
  labelCol: {
    span: 6,
  },
  wrapperCol: {
    span: 14,
  },
}

@Form.create()
class UserModal extends PureComponent {
  handleOk = () => {
    const { item = {}, onOk, form } = this.props
    const { validateFields, getFieldsValue } = form

    validateFields(errors => {
      if (errors) {
        return
      }
      const data = {
        ...getFieldsValue(),
        key: item.key,
      }
      onOk(data)
    })
  }

  render() {
    const { item = {}, onOk, form, type, ...modalProps } = this.props
    const { getFieldDecorator } = form

    return (
      <Modal {...modalProps} onOk={this.handleOk}>
        <Form layout="horizontal">
          <FormItem label={`学生姓名`} hasFeedback {...formItemLayout}>
            {getFieldDecorator('student_name', {
              initialValue: item.student_name,
              rules: [
                {
                  required: true,
                },
              ],
            })(<Input />)}
          </FormItem>
          <FormItem label={`学号`} hasFeedback {...formItemLayout}>
            {getFieldDecorator('student_no', {
              initialValue: item.student_no,
              rules: [
                {
                  required: true,
                },
              ],
            })(<Input disabled={ type === 'create' ? false : true } />)}
          </FormItem>
          <FormItem label={'更新密码'} hasFeedback {...formItemLayout}>
            {getFieldDecorator('student_password', {
              rules: [

              ],
            })(<Input />)}
          </FormItem>
          {/* <FormItem label={`班级`} hasFeedback {...formItemLayout}>
            {getFieldDecorator('class_ids', {
              initialValue: item.class_ids,
              valuePropName: "classIds",
            })(<Selects totalClasses={item.totalClasses || []}/>)}
          </FormItem> */}
          <FormItem label={`头像图片`} hasFeedback {...formItemLayout}>
            {getFieldDecorator('student_image', {
              initialValue: item.student_image,
              valuePropName: "image",
            })(<PicturesWall />)}
          </FormItem>
        </Form>
      </Modal>
    )
  }
}

UserModal.propTypes = {
  item: PropTypes.object,
  onOk: PropTypes.func,
}

export default UserModal
